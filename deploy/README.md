# deploy

This folder space is dedicated to k8s, helm, ArgoCD configurations files.
Consider it the infrastructure files to be created and replicated from one command.

## North Star
The applications and db configurations should come from open-source tools only. The ideology is to be able to setup same environment anywhere.

- we always have two cluster per environment (infra, app)
	- infra cluster should hold applications that are not app related.
		Usually those apps access will have something like `{app-name}.infra.{base-url}` and should have authentication enabled at all time.
		Example: ArgoCD, BackStage, Image Builder, Image Registry, pkgs registry, tracer, logs db. Grafana etc
	- app cluster: should hold all app related to create a fully functional app. Example: microservices, databases, pubsub, cache, istio, Trafeik

## Tools

- Kubectl
- k3s
- [k3d](http://k3d.io)
- cloud-provider-kind

