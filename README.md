# Dimios Kubernetes Manager
[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/dimios)](https://goreportcard.com/report/github.com/seibert-media/dimios)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/513590eff4e54095a25b66bf65bd1323)](https://www.codacy.com/app/seibert-media/dimios?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=seibert-media/dimios&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/seibert-media/dimios.svg?branch=master)](https://travis-ci.org/seibert-media/dimios)
[![Docker Repository on Quay](https://quay.io/repository/seibertmedia/dimios/status "Docker Repository on Quay")](https://quay.io/repository/seibertmedia/dimios)

Dimios is an application to keep a kubernetes cluster on a clean state based on provided manifests.
Those manifests get automatically applied and no longer existing resources are being removed by the application.

## Requirements

To run the tool using our [Makefile](Makefile) you need to set a few local environment variables.
This can be done manually or by creating a file called `.env` in the projects root directory:

```bash
export NAMESPACE=dimios
export MANIFEST_DIR=~/git/k8s/smedia-kubernetes/
```

If you want to use [TeamVault](https://github.com/trehn/teamvault) as passwords backend inside templates,
you need to have a file called `~/.teamvault-sm.json` containing your credentials:

```json
{
  "url": "https://teamvault.example.com",
  "user": "username",
  "pass": "password"
}
```

## Options

**-dir** _string_
: Path to the template/manifest directory

**-kubeconfig** _string_
: (optional) absolute path to the kubeconfig file (default "~/.kube/config")

**-staging**
: If set the application will run in no-op mode not doing any changes to the cluster

**-namespaces** _string_
: List of kubernetes namespaces separated by comma

**-teamvault-config** _string_
: Teamvault config path

**-teamvault-pass** _string_
: Teamvault password (if `teamvault-config` is unset)

**-teamvault-url** _string_
: Teamvault url (if `teamvault-config` is unset)

**-teamvault-user** _string_
: Teamvault user (if `teamvault-config` is unset)

**-v** _value_
: Log level for glog V logs

**-logtostderr**
: Log to standard error instead of files

**-version**
: Show version info (default false)

**-webhook**
: Start as webserver and sync on each call

**-port** _int_
: Port webserver listens on

**-whitelist** _string_
: List of kind to manange separated by comma

## Dependencies
All dependencies inside this project are being managed by [dep](https://github.com/golang/dep) and are checked in.
After pulling the repository, it should not be required to do any further preparations aside from `make deps` to prepare the dev tools.

If new dependencies get added while coding, make sure to add them using `dep ensure --add "importpath"`.

## Contributing
Feedback and contributions are highly welcome. Feel free to file issues or pull requests.

## Related Projects

* [github.com/heptio/ark](https://github.com/heptio/ark)
* [github.com/hasura/gitkube](https://github.com/hasura/gitkube)

## Attributions

* [Kolide for providing `kit`](https://github.com/kolide/kit)

## Examples

Simple dryrun

```
dimios \
-staging \
-dir ~/manifests \
-namespaces download \
-teamvault-config ~/.teamvault.json \
-kubeconfig ~/.kube/config \
-logtostderr \
-v=1
```

Manange only Services and Deployment with dryrun

```
dimios \
-staging \
-dir ~/manifests \
-namespaces download \
-teamvault-config ~/.teamvault.json \
-kubeconfig ~/.kube/config \
-logtostderr \
-v=1 \
-whitelist Deployment,Service
```


Start as webserver on port 9999 with dryrun

```
dimios \
-staging \
-dir ~/manifests \
-namespaces download \
-teamvault-config ~/.teamvault.json \
-kubeconfig ~/.kube/config \
-logtostderr \
-v=1 \
-port=9999 \
-webhook
```
