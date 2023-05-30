#!/bin/bash

command -v minikube || echo "Please install minikube"
command -v kubectl || echo "Please install kubectl"
command -v helm || echo "Please install helm"

minikube start --profile startupbuilder
minikube start --profile startuptools

kubectl create nampespace startupbuilder

# adding prometheus repo
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helo repo update

# Create namespace for infra toolings
kubectl create ns startupmonitor

# Install prometheus using helm
helm upgrade --install prom prometheus-community/kube-prometheus-stack -n startupmonitor --values Charts/prometheus/values.yaml
