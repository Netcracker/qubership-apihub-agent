# qubership-apihub-agent

qubership-apihub-agent is a golang application (microservice) providing the following features:

- Designed for k8s deployment, Requires ClusterView permissions
- Using k8s API able to execute REST calls to any service
- Discovers API specificiations (OpenAPI, GQL) exposed by services
- Integrates with Qubership-APIHUB application

Not production ready, under development.

## Arch diagramm

todo

## Installation

Qubership APIHUB Agent is designed to work in k8s cluster only.
Deployment to k8s is possible via Helm Chart. Please refer to [Installation Notes](./helm-templates/README.md)


## Build

Just run build.cmd(sh) file from this repository
