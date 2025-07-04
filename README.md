# Provider Tailscale

`provider-tailscale` is a [Crossplane](https://crossplane.io/) provider that
is built using [Upjet](https://github.com/crossplane/upjet) code
generation tools and exposes XRM-conformant managed resources for the
Tailscale API.

## Getting Started

Install the provider by using the following command after changing the image tag
to the [latest release](https://marketplace.upbound.io/providers/supahlab/provider-tailscale):

```console
crossplane xpkg install provider xpkg.crossplane.io/supahlab/provider-tailscale:v0.1.0
```

Alternatively, you can use declarative installation:

```console
cat <<EOF | kubectl apply -f -
apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-tailscale
spec:
  package: xpkg.crossplane.io/supahlab/provider-tailscale:v0.1.0
EOF
```

Notice that in this example Provider resource is referencing ControllerConfig with debug enabled.

You can see the API reference at [doc.crds.dev/github.com/supahlab/provider-tailscale](https://doc.crds.dev/github.com/supahlab/provider-tailscale).

## Developing

Run code-generation pipeline:

```console
go run cmd/generator/main.go "$PWD"
```

Run against a Kubernetes cluster:

```console
make run
```

Build, push, and install:

```console
make all
```

Build binary:

```console
make build
```

## Report a Bug

For filing bugs, suggesting improvements, or requesting new features, please
open an [issue](https://github.com/supahlab/provider-tailscale/issues).
