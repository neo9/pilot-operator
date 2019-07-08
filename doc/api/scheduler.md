## Scheduler

Scheduler describes common cronJobs requests

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| metadata | Standard object’s metadata | [metav1.ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#objectmeta-v1-meta) | false |
| spec | Specification of the desired behavior of the Scheduler| [SchedulerSpec](#schedulerspec) | true |
| status | Most recent observed status of the Scheduler. Read-only | *[SchedulerStatus](#schedulerstatus) | false |

## SchedulerSpec

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| requests | CronJobs requests | *[][SchedulerRequest](#schedulerrequest) | false |

## SchedulerRequest

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Unique request name | string | true |
| method | Request method. One of ['get', 'post', 'put', 'patch', 'delete']. Default `post` | string | false |
| scheduler | CronJob schedule. Example: `0 * * * *` | string | true |
| service | Service target | string | true |
| port | Target http port. Default `80` | int32 | false |
| path | Target http path | string | true |

## SchedulerStatus

To be defined

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
