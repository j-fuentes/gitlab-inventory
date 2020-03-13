# Gitlab project inventary

This generates a list of all the projects in your Gitlab instance.

## Install

```
$ go get github.com/j-fuentes/gitlab-inventory
$ go install -i github.com/j-fuentes/gitlab-inventory
```

## Usage

Generate an access token for the Gitlab API and execute the following:

```
$ export GITLAB_URL="https://gitlab.mydomain.com"
$ export GITLAB_TOKEN="mytoken"
$ gitlab-inventory
```
