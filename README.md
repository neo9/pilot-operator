# Pilot Operator

Kubernetes operator to manage and deploy Java, NodeJS and web applications.

[API documentation](doc/api.md)

Work in progress.

## Requirements

- Kubernetes 1.12.x (other versions are not tested yet)


## Pilot operator deployment

```bash
git clone https://github.com/neo9/pilot-operator
cd pilot-operator/deploy
kubectl -n [namespace] apply -f ./crds
kubectl -n [namespace] apply -f ./
```

## Examples

### Application


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
