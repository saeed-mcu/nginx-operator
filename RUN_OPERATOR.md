## Run the Operator
There are three ways to run the operator:
* As a Go program outside a cluster
* As a Deployment inside a Kubernetes cluster
* Managed by the Operator Lifecycle Manager (OLM) in bundle format

1. Run locally outside the cluster
```bash
make install run
```
OR
```bash
# Apply CRD into Cluster
kustomize build config/crd/ | kubectl apply -f -
# run controller
go run ./cmd/main.go
```
Then apply your CR:

```bash
kubectl create -f config/samples/operator_v1alpha1_nginxoperator.yaml
```
