/*
Copyright © 2022 Andrew Feller <andyfeller@github.com>

*/
package cmd

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/spf13/cobra"
)

var (
	owner        string
	repoLimits   []string
	repoExcludes []string
	outputFile   string
)

type repositoryQuery struct {
	Repository struct {
		Name  string
		Owner struct {
			Login string
		}
		DependencyGraphManifests struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []struct {
				Filename     string
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
				} `graphql:"dependencies(first: 100, after: $dependencyCursor)"`
				DependenciesCount int
				ExceedsMaxSize    bool
				Parseable         bool
			}
		} `graphql:"dependencyGraphManifests(first: 10, after: $manifestCursor)"`
	} `graphql:"repository(owner: $owner, name: $name)"`
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gh-dependency-report [flags] owner [repo ...]",
	Short: "Generate report of repository manifests and dependencies discovered through the dependency graph",
	Long:  "Generate report of repository manifests and dependencies discovered through the dependency graph",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := gh.GQLClient(&api.ClientOptions{
			Headers: map[string]string{
				"Accept": "application/vnd.github.hawkgirl-preview+json",
			},
		})

		if err != nil {
			return err
		}

		owner = args[0]
		repos := []string{}

		if len(args) > 1 {
			repos = args[1:]
		}

		if len(repos) <= 0 {
			// TODO: Get list of repositories by owner if repoLimits is empty and exclude repoExcludes
			return errors.New("Need to implement logic to retrieve owner repositories")
		}

		if len(repoExcludes) > 0 {
			return errors.New("Need to implement logic to exclude repositories from report")
		}

		csvWriterOutput := os.Stdout

		if len(outputFile) > 0 {
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

		for _, repo := range repos {
			query := new(repositoryQuery)
			hasNextPage := true
			manifestCursor := ""
			dependencyCursor := ""

			for hasNextPage {
				variables := map[string]interface{}{
					"owner":            graphql.String(owner),
					"name":             graphql.String(repo),
					"manifestCursor":   graphql.String(manifestCursor),
					"dependencyCursor": graphql.String(dependencyCursor),
				}

				err = client.Query("repoDependencies", query, variables)

				if err != nil {
					return err
				}

				fmt.Println(query.Repository.Name, query.Repository.Owner)
				repository := query.Repository
				hasNextPage = repository.DependencyGraphManifests.PageInfo.HasNextPage
				manifestCursor = repository.DependencyGraphManifests.PageInfo.EndCursor

				for _, manifest := range repository.DependencyGraphManifests.Nodes {
					fmt.Println(manifest)
					hasNextPage = hasNextPage || manifest.Dependencies.PageInfo.HasNextPage
					dependencyCursor = manifest.Dependencies.PageInfo.EndCursor

					for _, dependency := range manifest.Dependencies.Nodes {
						csvWriter.Write([]string{
							repository.Owner.Login,
							repository.Name,
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
				}
			}
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
}
