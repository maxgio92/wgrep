[![Latest release](https://img.shields.io/github/v/release/maxgio92/wgrep?style=for-the-badge)](https://github.com/maxgio92/wgrep/releases/latest)
[![License](https://img.shields.io/github/license/maxgio92/wgrep?style=for-the-badge)](COPYING)
![Go version](https://img.shields.io/github/go-mod/go-version/maxgio92/wgrep?style=for-the-badge)

# wgrep: like grep but for web sites

`wgrep` (world wide web grep) search for patterns in a web site directory hierarchy over HTTP, through hypertext references.

By default, it searches for patterns inside paragraphs.

The tool is inspired by GNU `grep(1)` and `wget(1)`.

### Usage

```
wgrep PATTERN URL [flags]
```

For details please read the CLI [documentation](./docs/wgrep.md).

#### Recursive search

```shell
wgrep --recursive|-r PATTERN URL
```

#### Case insensitive patterns

```shell
wgrep --ignore-case|-i PATTERN URL
```

#### Include specific locations only

While referenced locations that have a host name different from the one specified in the `URL` argument are skipped by default, it's possible to include only locations of which HTTP path follows a specific pattern.

Similarly to how `grep` allows with the `--include` flag to include specific locations in the search, it's possible to filter the pages by URL when recursively look for a pattern.
The include location filter pattern supports **regular expressions** in the [Go flavor](https://pkg.go.dev/regexp/syntax).

```shell
wgrep -r --include "my-section\/.+" PATTERN URL
```

#### Search on specific HTML elements

By default, the element filter is set to "p", as standard paragraphs are represented in HTML. However this filter can be customized with the `--element`|`-e` flag:

```shell
wgrep --element|-e "article" PATTERN URL
```

The element filter supports [GoQuery](https://github.com/PuerkitoBio/goquery) patterns.
For example, this allows to select elements based on class attributes:

```shell
wgrep -e ".my-class" PATTERN URL
```

For more information about the selector syntax please refer to the [GoQuery](https://github.com/PuerkitoBio/goquery) documentation.

### In action

```shell
$ wgrep --include "posts\/" -ri kubernetes https://blog.maxgio.me
https://blog.maxgio.me/posts/k8s-stride-05-denial-of-service/:
Users that are authorized to make patch requests to the Kubernetes API server can send a specially crafted patch of type json-patch (e.g. kubectl patch - type json or Content-Type: application/json-patch+json) that consumes excessive resources while processing, causing a denial of service on the API server.

https://blog.maxgio.me/posts/stride-threat-modeling-kubernetes-elevation-of-privileges/: Hello everyone, a long time has passed after the 5th part of this journey through STRIDE thread modeling in Kubernetes has been published.
If you recall well, STRIDE is a model of threats for identifying security threats, by providing a mnemonic for security threats in six categories:

https://blog.maxgio.me/posts/stride-threat-modeling-kubernetes-elevation-of-privileges/:
In Kubernetes Role-Based Access Control authorizes or not access to Kubernetes resources through roles, but we also have underlying infrastructure resources, and Kubernetes provides primitives to authorize workload to access operating system resources, like Linux namespaces.
...
```
