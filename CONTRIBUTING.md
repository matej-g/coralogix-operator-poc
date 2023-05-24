# Contributing to Coralogix Operator

Thank you for your interest in contributing to the Coralogix Operator. We welcome your contributions. Here you'll find information to help you get started with operator development.

You are contributing under the terms and conditions of the [Contributor License Agreement](LICENSE). [For signing](https://cla-assistant.io/coralogix/coralogix-operator-poc).

Building the Operator
---------------------

### Requirements

- You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind)
- [Go](https://golang.org/doc/install) 1.20.x (for building the operator locally)
- [kuttl](https://kuttl.dev/) (for running e2e tests).
- [kubebuilder](https://book-v1.book.kubebuilder.io/getting_started/installation_and_setup.html) for developing the operator.

### Steps

Clone the repository locally.

```sh
$ git clone git@github.com:coralogix/coralogix-operator-poc
```

Navigate to the operator directory and install it.
This command will build the operator and install the CRDs located at [crd directory](./config/crd) into the K8s cluster
specified in ~/.kube/config.

```sh
$ make install
```

Add the region and api-key as environment variables (or later as flags).

```sh
$ export CORALOGIX_API_KEY="xxx-xxx-xxx"
$ export CORALOGIX_ENV = "EUROPE2"
```

Run the operator locally
```sh
$ go run main.go
```
Or with `regin` and `api-key` flags
```sh
$ go run main.go -region EUROPE2 -api-key xxx-xxx-xxx
```

Running examples
---------------------
It's possible to use one of the samples in the [sample directory](./config/samples) or creating your own sample file.
Then apply it -

```sh
$ kubectl apply -f config/samples/alerts/standard_alert.yaml
```

Getting the resource status

```sh
$ kubectl get alerts.coralogix.coralogix standard-alert-example -oyaml
```

Destroying the resource.

```sh
$ kubectl delete alerts.coralogix.coralogix standard-alert-example
```

Developing
---------------------
We use [kubebuilder](https://book.kubebuilder.io/) for developing the operator.
When creating or updating CRDs remember to run 
```sh
make manifests
````

Tests
---------------------
We use [kuttl](https://kuttl.dev/) for end-to-end tests.
The test files are located at [./tests/e2e/](./tests/e2e)
In order to run the full suite of `kuttl` tests, run -
```sh
$ make e2e
````

*Note:* `kuttl` tests create real resources and in a case of failure some resources may not be removed.

Releases
---------------------
To determine the release convention we use [semantic-release](.releaserc.json) -
```sh
"releaseRules": [
          {"message": "major*", "release": "major"},
          {"message": "minor*", "release": "minor"},
          {"message": "patch*", "release": "patch"}
        ]
````