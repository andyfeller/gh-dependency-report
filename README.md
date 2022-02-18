# gh-dependency-report

A `gh` extension to generate report of repository manifests and dependencies discovered through GitHub's [software supply chain](https://docs.github.com/en/code-security/supply-chain-security) capabilities.

![Demo of gh-dependency-report extension](https://user-images.githubusercontent.com/2089743/154634826-716abba3-f139-4b7a-a106-01c0ab5b68c4.gif)

## Quickstart

1. `gh extension install andyfeller/gh-dependency-report`
1. `gh dependency-report $(whoami)`
1. Profit! :moneybag: :money_with_wings: :money_mouth_face: :money_with_wings: :moneybag:

## Usage

Pulling [manifests](https://docs.github.com/en/graphql/reference/objects#dependencygraphmanifest) and [dependencies](https://docs.github.com/en/graphql/reference/objects#dependencygraphdependency) including [license info](https://docs.github.com/en/graphql/reference/objects#license) around [repositories](https://docs.github.com/en/graphql/reference/objects#repository) from [GitHub's GraphQL API](https://docs.github.com/en/graphql/reference/).  This is only works for repositories that have [enabled the dependency graph feature](https://docs.github.com/en/code-security/supply-chain-security/understanding-your-software-supply-chain/about-the-dependency-graph#enabling-the-dependency-graph).

The result is a CSV that companies and individuals can use to attest to software licenses in use, making the jobs of platform engineering, legal, security, and auditors easier.

```shell
 $ gh dependency-report --help

Generate report of repository manifests and dependencies discovered through the dependency graph

Usage:
  gh-dependency-report [flags] owner [repo ...]

Flags:
  -d, --debug                Whether to debug logging
  -e, --exclude strings      Repositories to exclude from report
  -h, --help                 help for gh-dependency-report
  -o, --output-file string   Name of file to write CSV report (default "report-20220216081518.csv")
```

The resulting CSV file contains the most common information used for these purposes:

<dl>
  <dt><code>Owner</code></dt>
  <dd>Login name of the organization or user that owns the repository</dd>
  <dd>
    Examples:
    <ul>
      <li><code>andyfeller</code></li>
      <li><code>github</code></li>
      <li><code>cli</code></li>
    </ul>
  </dd>

  <dt><code>Repo</code></dt>
  <dd>Name of the repository containing the manifest; does not duplicate owner information</dd>
  <dd>
    Examples:
    <ul>
      <li><code>gh-dependency-report</code> <em>(for <code>andyfeller/gh-dependency-report</code>)</em></li>
      <li><code>codeql</code> <em>(for <code>github/codeql</code>)</em></li>
      <li><code>cli</code> <em>(for <code>cli/cli</code>)</em></li>
    </ul>
  </dd>

  <dt><code>Manifest</code></dt>
  <dd>Fully qualified manifest filename</dd>
  <dd>
    Examples:
      <li><code>go.mod</code></li>
      <li><code>.github/workflows/release.yml</code></li>
      <li><code>package.json</code></li>
  </dd>

  <dt><code>Exceeds Max Size</code></dt>
  <dd>Is the manifest too big to parse?</dd>

  <dt><code>Parseable</code></dt>
  <dd>Were we able to parse the manifest?</dd>

  <dt><code>Package Manager</code></dt>
  <dd>The dependency package manager.</dd>
  <dd>
    Examples:
      <li><code>ACTIONS</code></li>
      <li><code>COMPOSER</code></li>
      <li><code>GO</code></li>
      <li><code>MAVEN</code></li>
      <li><code>NPM</code></li>
      <li><code>NUGET</code></li>
      <li><code>PIP</code></li>
      <li><code>RUBYGEMS</code></li>
  </dd>

  <dt><code>Dependency</code></dt>
  <dd>
    The name of the package in the canonical form used by the package manager.  This may differ from the original textual form (see packageLabel), for example in a package manager that uses case-insensitive comparisons.
  </dd>
  <dd>
    Examples:
      <li><code>actions/checkout</code> <em>(actions)</em></li>
      <li><code>github.com/spf13/cobra</code> <em>(go)</em></li>
      <li><code>@actions/core</code> <em>(npm)</em></li>
  </dd>

  <dt><code>Has Dependencies?</code></dt>
  <dd>Does the dependency itself have dependencies?</dd>

  <dt><code>Requirements</code></dt>
  <dd>The dependency version requirements.</dd>

  <dt><code>License</code></dt>
  <dd>Short identifier specified by <a href="https://spdx.org/licenses">https://spdx.org/licenses</a>.</dd>

  <dt><code>License Url</code></dt>
  <dd>URL to the license on <a href="https://choosealicense.com">https://choosealicense.com</a>.</dd>
</dl>

### Example Report

The following is an example of a report generated around my own personal repositories:

<details>
  <summary>Example report on <code>andyfeller</code></summary>

  ```
  Owner,Repo,Manifest,Exceeds Max Size,Parseable,Package Manager,Dependency,Has Dependencies?,Requirements,License,License Url
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/cli/go-gh,true,= 0.0.2-0.20211206104242-8180ab76d996,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/cli/safeexec,false,= 1.0.0,,
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/cli/shurcooL-graphql,true,= 0.0.1,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/henvic/httpretty,false,= 0.0.6,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/inconshreveable/mousetrap,false,= 1.0.0,,
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/spf13/cobra,true,= 1.3.0,Apache-2.0,http://choosealicense.com/licenses/apache-2.0/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,github.com/spf13/pflag,false,= 1.0.5,,
  andyfeller,gh-dependency-report,go.mod,false,true,GO,go.uber.org/atomic,true,= 1.9.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,go.uber.org/multierr,true,= 1.7.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,go.uber.org/zap,true,= 1.20.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.mod,false,true,GO,golang.org/x/net,false,= 0.0.0-20211112202133-69e39bad7dc2,,
  andyfeller,gh-dependency-report,go.mod,false,true,GO,gopkg.in/yaml.v3,true,= 3.0.0-20210107192922-496545a6307b,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/benbjohnson/clock,false,= v1.1.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/cli/go-gh,true,= v0.0.2-0.20211206104242-8180ab76d996,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/cli/safeexec,false,= v1.0.0,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/cli/shurcooL-graphql,true,= v0.0.1,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/davecgh/go-spew,false,= v1.1.1,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/henvic/httpretty,false,= v0.0.6,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/inconshreveable/mousetrap,false,= v1.0.0,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/kr/pretty,true,= v0.2.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/kr/text,true,= v0.1.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/MakeNowJust/heredoc,false,= v1.0.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/pkg/errors,false,= v0.8.1,BSD-2-Clause,http://choosealicense.com/licenses/bsd-2-clause/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/pmezard/go-difflib,false,= v1.0.0,NOASSERTION,http://choosealicense.com/licenses/other/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/spf13/cobra,true,= v1.3.0,Apache-2.0,http://choosealicense.com/licenses/apache-2.0/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/spf13/pflag,false,= v1.0.5,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,github.com/stretchr/testify,true,= v1.7.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,go.uber.org/atomic,true,= v1.9.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,go.uber.org/goleak,true,= v1.1.11,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,go.uber.org/multierr,true,= v1.7.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,go.uber.org/zap,true,= v1.20.0,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,go.sum,false,true,GO,golang.org/x/net,false,= v0.0.0-20211112202133-69e39bad7dc2,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,gopkg.in/check.v1,true,= v1.0.0-20190902080502-41f04d3bba15,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,gopkg.in/yaml.v2,true,= v2.4.0,,
  andyfeller,gh-dependency-report,go.sum,false,true,GO,gopkg.in/yaml.v3,true,= v3.0.0-20210107192922-496545a6307b,,
  andyfeller,gh-dependency-report,.github/workflows/release.yml,false,true,ACTIONS,actions/checkout,false,= 2,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,.github/workflows/release.yml,false,true,ACTIONS,cli/gh-extension-precompile,false,= 1,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,.github/workflows/release.yml,false,true,ACTIONS,actions/checkout,false,= 2,MIT,http://choosealicense.com/licenses/mit/
  andyfeller,gh-dependency-report,.github/workflows/release.yml,false,true,ACTIONS,cli/gh-extension-precompile,false,= 1,MIT,http://choosealicense.com/licenses/mit/
  ```
</details>


## Setup

Like any other `gh` CLI extension, `gh-dependency-report` is trivial to install or upgrade and works on most operating systems:

- **Installation**

  ```shell
  gh extension install andyfeller/gh-dependency-report
  ```
  
  _For more information: [`gh extension install`](https://cli.github.com/manual/gh_extension_install)_

- **Upgrade**

  ```shell
  gh extension upgrade gh-dependency-report
  ```

  _For more information: [`gh extension upgrade`](https://cli.github.com/manual/gh_extension_upgrade)_
