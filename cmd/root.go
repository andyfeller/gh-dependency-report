/*
Copyright © 2022 Andrew Feller <andyfeller@github.com>

*/
package cmd

import (
	"encoding/csv"
	"errors"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var (
	owner        string
	repoExcludes []string
	outputFile   string
	client       api.GQLClient
	sugar        *zap.SugaredLogger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-dependency-report [flags] owner [repo ...]",
	Short: "Generate report of repository manifests and dependencies discovered through the dependency graph",
	Long:  "Generate report of repository manifests and dependencies discovered through the dependency graph",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		var err error

		client, err = gh.GQLClient(&api.ClientOptions{
			Headers: map[string]string{
				"Accept": "application/vnd.github.hawkgirl-preview+json",
			},
		})

		if err != nil {
			return err
		}

		owner = args[0]

		// Resolve repositories in scope of report
		var repos []string

		if len(args) > 1 {
			repos = args[1:]
		} else {
			repos = make([]string, 0, 100) // Struggle for initial slice length given potential growth for large organizations
			var reposCursor *string

			for {
				reposQuery, err := getRepos(owner, reposCursor)

				if err != nil {
					return err
				}

				for _, repo := range reposQuery.RepositoryOwner.Repositories.Nodes {
					repos = append(repos, repo.Name)
				}

				reposCursor = &reposQuery.RepositoryOwner.Repositories.PageInfo.EndCursor

				if !reposQuery.RepositoryOwner.Repositories.PageInfo.HasNextPage {
					break
				}
			}
		}

		sort.Strings(repos)

		if len(repoExcludes) > 0 {
			sort.Strings(repoExcludes)
			sugar.Debugf("Excluding repos", "repos", repoExcludes)

			for _, repoExclude := range repoExcludes {
				for i, repo := range repos {
					if repoExclude == repo {
						repos = append(repos[:i], repos[i+1:]...)
					}
				}
			}
		}

		if len(repos) <= 0 {
			return errors.New("No repositories to report on")
		}

		sugar.Infof("Processing repos: %s", repos)

		// Prepare writer for outputting report
		csvWriterOutput := os.Stdout

		if len(outputFile) > 0 {
			sugar.Debugf("Setting up output file \"%s\"", outputFile)

			if _, err := os.Stat(outputFile); errors.Is(err, os.ErrExist) {
				return err
			}

			output, err := os.OpenFile(outputFile, os.O_WRONLY|os.O_CREATE, 0644)

			if err != nil {
				return err
			}

			csvWriterOutput = output
		}

		csvWriter := csv.NewWriter(csvWriterOutput)

		csvWriter.Write([]string{
			"Owner",
			"Repo",
			"Manifest",
			"Exceeds Max Size",
			"Parseable",
			"Package Manager",
			"Dependency",
			"Has Dependencies?",
			"Requirements",
			"License",
			"License Url",
		})

		// Retrieve data and produce report
		var backoffQueue []manifestBackoff

		for _, repo := range repos {
			var manifestCursor *string
			sugar.Debugf("Processing %s/%s", owner, repo)

			for {
				manifestsQuery, err := getManifests(owner, repo, manifestCursor)

				if err != nil {
					wtf := err.Error()
					if strings.Contains(wtf, "Message: loading") {
						backoffQueue = append(backoffQueue, manifestBackoff{
							Owner:          owner,
							RepositoryName: repo,
							EndCursor:      manifestCursor,
						})
						break
					} else {
						return err
					}
				}

				for _, manifest := range manifestsQuery.Repository.DependencyGraphManifests.Nodes {
					var dependencyCursor *string
					sugar.Debugf("Processing %s/%s > %s", owner, repo, manifest.Filename)

					for {
						dependenciesQuery, err := getDependencies(manifest.Id, dependencyCursor)

						if err != nil {
							sugar.Warnf("Error processing %s/%s > %s: %s", owner, repo, manifest.Filename, err)
							break
						}

						for _, dependency := range dependenciesQuery.Node.DependencyGraphManifest.Dependencies.Nodes {
							sugar.Debugf("Processing %s/%s > %s > %s", owner, repo, manifest.Filename, dependency.PackageName)

							csvWriter.Write([]string{
								owner,
								repo,
								manifest.Filename,
								strconv.FormatBool(manifest.ExceedsMaxSize),
								strconv.FormatBool(manifest.Parseable),
								dependency.PackageManager,
								dependency.PackageName,
								strconv.FormatBool(dependency.HasDependencies),
								dependency.Requirements,
								dependency.Repository.LicenseInfo.SpdxId,
								dependency.Repository.LicenseInfo.Url,
							})
						}

						dependencyCursor = &dependenciesQuery.Node.DependencyGraphManifest.Dependencies.PageInfo.EndCursor

						if !dependenciesQuery.Node.DependencyGraphManifest.Dependencies.PageInfo.HasNextPage {
							break
						}
					}
				}

				manifestCursor = &manifestsQuery.Repository.DependencyGraphManifests.PageInfo.EndCursor

				if !manifestsQuery.Repository.DependencyGraphManifests.PageInfo.HasNextPage {
					break
				}
			}
		}

		if len(backoffQueue) > 0 {
			sugar.Debugf("Reconciling back off queue: %s", backoffQueue)
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringSliceVarP(&repoExcludes, "exclude", "e", []string{}, "Repositories to exclude from report")
	rootCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Name of file to write CSV report, defaults to stdout")

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	sugar = logger.Sugar()
}

type dependenciesQuery struct {
	Node struct {
		DependencyGraphManifest struct {
			Dependencies struct {
				PageInfo struct {
					HasNextPage bool
					EndCursor   string
				}
				Nodes []struct {
					HasDependencies bool
					PackageManager  string
					PackageName     string
					Repository      struct {
						LicenseInfo struct {
							SpdxId string
							Url    string
						}
					}
					Requirements string
				}
				TotalCount int
			} `graphql:"dependencies(first: 100, after: $endCursor)"`
		} `graphql:"... on DependencyGraphManifest"`
	} `graphql:"node(id: $id)"`
}

func getDependencies(id string, endCursor *string) (*dependenciesQuery, error) {
	query := new(dependenciesQuery)
	variables := map[string]interface{}{
		"id":        graphql.ID(id),
		"endCursor": (*graphql.String)(endCursor),
	}

	err := client.Query("getDependencies", query, variables)

	return query, err
}

type manifestsQuery struct {
	Repository struct {
		DependencyGraphManifests struct {
			Nodes []struct {
				DependenciesCount int
				ExceedsMaxSize    bool
				Filename          string
				Id                string
				Parseable         bool
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			TotalCount int
		} `graphql:"dependencyGraphManifests(first: 10, after: $endCursor, withDependencies: true)"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}

func getManifests(owner string, repo string, endCursor *string) (*manifestsQuery, error) {
	query := new(manifestsQuery)
	variables := map[string]interface{}{
		"owner":     graphql.String(owner),
		"repo":      graphql.String(repo),
		"endCursor": (*graphql.String)(endCursor),
	}

	err := client.Query("getManifests", query, variables)

	return query, err
}

type reposQuery struct {
	RepositoryOwner struct {
		Repositories struct {
			Nodes []struct {
				Name string
			}
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
		} `graphql:"repositories(first: 100, after: $endCursor, ownerAffiliations: [OWNER])"`
	} `graphql:"repositoryOwner(login: $owner)"`
}

func getRepos(owner string, endCursor *string) (*reposQuery, error) {
	query := new(reposQuery)
	variables := map[string]interface{}{
		"owner":     graphql.String(owner),
		"endCursor": (*graphql.String)(endCursor),
	}

	err := client.Query("getRepos", query, variables)

	return query, err
}

type manifestBackoff struct {
	Owner          string
	RepositoryName string
	EndCursor      *string
}
