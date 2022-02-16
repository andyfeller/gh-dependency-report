package cmd

import (
	"bytes"
	"fmt"
	"testing"
)

type testGetter struct {
	stubs map[string][]interface{}
}

func newTestGetter() *testGetter {
	return &testGetter{
		stubs: map[string][]interface{}{},
	}
}

func (t *testGetter) Stub(queryName string, outQuery interface{}) {
	t.stubs[queryName] = append(t.stubs[queryName], outQuery)
}

func (t *testGetter) GetRepos(owner string, endCursor *string) (*reposQuery, error) {
	result := t.query("GetRepos")
	query := result.(reposQuery)
	return &query, nil
}

func (t *testGetter) GetManifests(owner string, repo string, endCursor *string) (*manifestsQuery, error) {
	result := t.query("GetManifests")
	query := result.(manifestsQuery)
	return &query, nil
}

func (t *testGetter) GetDependencies(id string, endCursor *string) (*dependenciesQuery, error) {
	result := t.query("GetDependencies")
	query := result.(dependenciesQuery)
	return &query, nil
}

func (t *testGetter) query(name string) interface{} {
	stubs, ok := t.stubs[name]
	if !ok || len(stubs) == 0 {
		panic(fmt.Sprintf("no stub for query: %s", name))
	}
	query := stubs[0]

	t.stubs[name] = t.stubs[name][1:]

	return query
}

// TODO these three stub helpers are arbitrary: they could take parameters and
// generate structs or you could make a new helper for each test case. It's
// entirely up to you.
func stubbedRepoList(names ...string) reposQuery {
	q := reposQuery{}
	for _, name := range names {
		q.RepositoryOwner.Repositories.Nodes = append(q.RepositoryOwner.Repositories.Nodes, struct{ Name string }{name})
	}

	return q
}

func stubbedManifests() manifestsQuery {
	q := manifestsQuery{}
	q.Repository.DependencyGraphManifests.Nodes = []manifest{
		{
			DependenciesCount: 1,
			Filename:          "cool.pkg",
			Id:                "123",
			Parseable:         false,
		},
	}
	return q
}

func stubbedDependencies() dependenciesQuery {
	q := dependenciesQuery{}
	q.Node.DependencyGraphManifest.Dependencies.Nodes = []dependency{
		{
			PackageManager: "Great manager",
			PackageName:    "Great package",
			Repository: struct {
				LicenseInfo licenseInfo
			}{
				LicenseInfo: licenseInfo{"123", "https://cool.license"},
			},
		},
	}
	q.Node.DependencyGraphManifest.Dependencies.TotalCount = 1
	return q
}

func TestCommand(t *testing.T) {
	cases := []struct {
		name         string
		owner        string
		repos        []string
		repoExcludes []string
		stubs        func(g *testGetter)
		wantOut      func(t *testing.T, output *bytes.Buffer)
	}{
		{
			name:         "example test",
			owner:        "OWNER",
			repos:        []string{},
			repoExcludes: []string{},
			stubs: func(g *testGetter) {
				g.Stub("GetRepos", stubbedRepoList("REPO"))
				g.Stub("GetManifests", stubbedManifests())
				g.Stub("GetDependencies", stubbedDependencies())
			},
			wantOut: func(t *testing.T, output *bytes.Buffer) {
				expected := "Owner,Repo,Manifest,Exceeds Max Size,Parseable,Package Manager,Dependency,Has Dependencies?,Requirements,License,License Url\nOWNER,REPO,cool.pkg,false,false,Great manager,Great package,false,,123,https://cool.license\n"
				if output.String() != expected {
					t.Errorf("\nExpected: %s\nGot:      %s", expected, output.String())
				}
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			client := newTestGetter()
			out := &bytes.Buffer{}
			if tt.stubs != nil {
				tt.stubs(client)
			}

			err := runCmd(tt.owner, tt.repos, tt.repoExcludes, client, out)
			if err != nil {
				t.Errorf("Did not expect error; got %s", err.Error())
			}

			tt.wantOut(t, out)
		})
	}
}
