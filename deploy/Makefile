K8S_NAMESPACE := "startup-monitor"

.PHONY: helm-repo
helm-repo:
	helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts
	helm install my-otel-demo open-telemetry/opentelemetry-demo

.PHONY: port-forward
port-forward:
	kubectl port-forward svc/my-otel-demo-frontendproxy 8080:8080


.PHONY: ksm
ksm:
	helm install ksm promethues-community/kube-state-metrics -n "$K8S_NAMESPACE"

.PHONY: nodeexporter
nodeexporter:
	helm install nodeexporter prometheus-community/promethues-node-exporter -n "startup-monitor"

.PHONY: setup-monitor
setup-monitor:
	echo "Installing Helm Charts"
setup-monitor: ksm nodeexporter
