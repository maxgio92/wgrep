# wgrep: like grep but for web sites

`wgrep` (world wide web grep) search for patterns in a web site directory hierarchy over HTTP, through hypertext references.

By default, it searches for patterns inside paragraphs.

The tool is inspired by GNU `grep(1)` and `wget(1)`.

### Usage

```
wgrep PATTERN URL [flags]
```

For details please read the CLI [documentation](./docs/wgrep.md).

### In action

```shell
$ wgrep -r -i kubernetes https://blog.maxgio.me
https://blog.maxgio.me/posts/k8s-stride-05-denial-of-service/:
Users that are authorized to make patch requests to the Kubernetes API server can send a specially crafted patch of type json-patch (e.g. kubectl patch - type json or Content-Type: application/json-patch+json) that consumes excessive resources while processing, causing a denial of service on the API server.

https://blog.maxgio.me/posts/stride-threat-modeling-kubernetes-elevation-of-privileges/: Hello everyone, a long time has passed after the 5th part of this journey through STRIDE thread modeling in Kubernetes has been published.
If you recall well, STRIDE is a model of threats for identifying security threats, by providing a mnemonic for security threats in six categories:

https://blog.maxgio.me/posts/stride-threat-modeling-kubernetes-elevation-of-privileges/:
In Kubernetes Role-Based Access Control authorizes or not access to Kubernetes resources through roles, but we also have underlying infrastructure resources, and Kubernetes provides primitives to authorize workload to access operating system resources, like Linux namespaces.
...
```
