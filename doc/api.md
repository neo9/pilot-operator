# API

This Document documents the types introduced by the Pilot Operator 

## Application

Application describes a common JAVA, NodeJS or Web application

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata | Standard object’s metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#objectmeta-v1-meta) | false |
| spec | Specification of the desired behavior of the Alertmanager cluster. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status | [ApplicationSpec](#applicationspec) | true |
| status | Most recent observed status of the Alertmanager cluster. Read-only. Not included when requesting from the apiserver, only from the Prometheus Operator API itself. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#spec-and-status | *[ApplicationStatus](#applicationstatus) | false |

## ApplicationSpec

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| type | One of `java`, `nodejs` or `web` that boostraps default resources and config | string | true |
| repository | Image without tag | string | true |
| initTag | First deployment tag. Default `latest` | string | false |
| tag | Deployment tag which overrides `initTag`. Used by CI/CD pipelines | string | false |
| service | Application service configuration | *[ApplicationService](#applicationService) | false |
| resources | Application resources configuration | *[ApplicationResources](#applicationResources) | false |
| healthCheck | Application health check configuration | *[ApplicationHealthCheck](#applicationHealthCheck) | false |
| pod | Application pod configuration | *[ApplicationPod](#applicationPod) | false |

## ApplicationService

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| targetPort | Pod http port pod to expose. Default `80` | int32 | false |
| port | Service http port .Default `80` | int32 | false |

## ApplicationResources

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| memory | Pod memory request and limit | string | false |
| cpu | Pod CPU request and no CPU limit (CFS Kernel issue that causes unwanted throttling) | string | false |

Default CPU resource values:
- `5m` for web applications
- `10m` for others

Default memory resource values:
- `20M` for web applications
- `186M` for NodeJS applications
- `300M` for Java applications

## ApplicationPod

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| secrets | Pod secrets binding | *[][ApplicationSecret](#applicationSecret) | false |

## ApplicationHealthCheck

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| path | Health check GET request path. See default values below | string | false | 

Default path values:
- `/actuator/heatlh` for JAVA apps
- `/ping` for NodeJS and Web applications

## ApplicationSecret

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | name of the secret | string | true |
| value | value of the environment variable to bind | string | true |

## ApplicationStatus

To be defined

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
