# GoPull

GoPull is a REST API written in go that runs a deployment script after receiving a webhook request from Bitbucket, GitHub, or GitLab. All local repositories configured to work with GoPull must be placed in the same directory known as `REPOS_DIR`. In order for this API to access the remote repository, an SSH access key must be added to the remote. A webhook can then be configured to fire on any events supported by the remote (code push by default). When the API receives a request, it checks which branch was affected by the change. If this branch matches the configured branch of the local repository, then the corresponding deployment script is run.

This API features a command line tool for easily adding new repositories or configuring current local repositories. In order for the CLI tool to work, the `.env` file must be set up to provide paths to the local installation of GoPull and `REPOS_DIR`.

## Installation

GoPull can be installed using `go get`:

```
go get -u github.com/mattmeyers/gopull
```

The repository can also be cloned and the binaries manually built.
