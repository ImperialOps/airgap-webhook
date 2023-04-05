# Microservice for Kubernetes Image Inventory :smile:

![Release](https://github.com/ImperialOps/airgap-webhook/workflows/release/badge.svg?branch=main)

This microservice acts as a [validating admission webhook](https://kubernetes.io/docs/reference/access-authn-authz/admission-controllers/) in a Kubernetes cluster. It maintains a stateful inventory of all images that are deployed to the cluster.

## Installation

Use the official [helm chart](https://github.com/ImperialOps/helm-charts).

## Configuration

The microservice can be configured using the following environment variables:

- NAMESPACE: The Kubernetes namespace to watch for image deployments. Defaults to all the namespaces.

## Usage

To use the microservice, deploy it as a validating webhook in your Kubernetes cluster. The webhook will be invoked whenever a new image deployment is created in the specified namespace. The webhook will validate all requests but also keep inventory of every image deployed.

You can view the inventory of deployed images by: TODO.

## Contributing

Contributions are welcome! To contribute, please fork the repository and submit a pull request.

## License

This microservice is licensed under the MIT License. See the LICENSE file for more information.
