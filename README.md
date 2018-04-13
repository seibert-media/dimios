# Kubernetes Deploy
[![Go Report Card](https://goreportcard.com/badge/github.com/seibert-media/k8s-deploy)](https://goreportcard.com/report/github.com/seibert-media/k8s-deploy)
[![Codacy Badge](https://api.codacy.com/project/badge/Grade/513590eff4e54095a25b66bf65bd1323)](https://www.codacy.com/app/kwiesmueller/k8s-deploy?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=seibert-media/k8s-deploy&amp;utm_campaign=Badge_Grade)
[![Build Status](https://travis-ci.org/seibert-media/k8s-deploy.svg?branch=master)](https://travis-ci.org/seibert-media/k8s-deploy)
[![Docker Repository on Quay](https://quay.io/repository/seibertmedia/k8s-deploy/status "Docker Repository on Quay")](https://quay.io/repository/seibertmedia/k8s-deploy)

Kubernetes Deploy is an application to keep a kubernetes cluster on a clean state based on provided manifests.
Those manifests get automatically applied and no longer existing resources are being removed by the application.

## Requirements

To run the tool using our [Makefile](Makefile) you need to set a few local environment variables.
This can be done manually or by creating a file called `.env` in the projects root directory:

```bash
export NAMESPACE=k8s-deploy
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

**-sentryDsn** _string_
: Sentry DSN key

**-debug**
: Enable debug mode

**-v** _value_
: Log level for glog V logs

**-logtostderr**
: Log to standard error instead of files

**-version**
: Show version info (default true)

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

