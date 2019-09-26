# Kubernetes Pilot Operator

[![Build Status](https://travis-ci.org/neo9/pilot-operator.svg?branch=master)](https://travis-ci.org/neo9/pilot-operator)

Kubernetes operator to manage and deploy Java, NodeJS and web applications.

[API documentation](doc/api.md)

[Docker image on Docker Hub](https://hub.docker.com/r/neo9sas/pilot-operator)

## Requirements

- Kubernetes 1.14.x / 1.13.x (other versions are not tested yet)
- Go 1.12.x
- Operator SDK v0.10.0

## Local

Build

```bash
go mod vendor -v
operator-sdk build image_name
```

Tests

```bash
operator-sdk test local ./test/e2e --up-local --namespace default
```

This will use the default ~/.kube/config configuration or you can
modify the env variable `KUBECONFIG`.

## Motivations

Ansible has a push only strategy and does not sync well resources.  If you
delete a deployment from the Ansible vars, you should also remember to delete
it manually from the cluster. Therefore, the architecture is not immutable.

We also found limitations with Kustomize which does not handle loops and
complex logic because it was not designed for it.  We used Helm instead but
we've had issues syncing the versions from our CI/CD pipelines which only
updates the tag and our source value file was not being synced. We did not want
to manage the charts directly in each repository to avoid misconfiguration and
syncing work accross our dozen microservices. We also had to duplicate the
main charts to have default values between NodeJS, Java and Web applications.
With Helm 2, Tiller also had lots of security issues and deploying multiple
Tiller instances made the workflow more complex.

Pilot operator is designed to fix the limitations we've had and to manage
deployments, hpa and services with simpler pipelines, logic and more flexible
configuration files. Getting back to the core of what makes Kubernetes great
and avoiding pushing the limits of Helm and Kustomize too far just because it
did not meet our needs.


## Pilot operator deployment

With kubectl:

```bash
git clone https://github.com/neo9/pilot-operator
cd pilot-operator/deploy
kubectl -n [namespace] apply -f ./crds
kubectl -n [namespace] apply -f ./
```

With Helm:

```bash
helm repo add n9 https://n9-charts.storage.googleapis.com
helm upgrade -i \
  --namespace integration \
  pilot n9/pilot-operator

# If you are deploying a second instance, crd are already created
# You can use --no-crd-hook instead
helm install --no-crd-hook \
  --namespace another-namespace \
  --name pilot n9/pilot-operator
```


## Examples

### Application

[API documentation](doc/api.md)

```yaml
apiVersion: pilot.neo9.fr/v1alpha1
kind: Application
metadata:
  name: my-api
  namespace: integration
spec:
  type: java
  repository: eu.gcr.io/example/my-api
  initTag: v0.0.1
  service:
    targetPort: 8080
  pod:
    secrets:
      - name: my-api
        key: MONGO_URI

```
