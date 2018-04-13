# Kubernetes Deploy

Deploys all manifest files to Kubernetes.

## Options

**--namespaces** _NAMESPACES_
: list of managed namespaces seperated by comma.  


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

## Related Projects

* [github.com/heptio/ark](https://github.com/heptio/ark)
* [github.com/hasura/gitkube](https://github.com/hasura/gitkube)

