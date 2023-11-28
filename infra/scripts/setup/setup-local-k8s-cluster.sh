
#!/bin/bash

# Create a new Minikube Cluster
minikube start --kubernetes-version=v1.28 --cpus=6 --memory=12000 --disk-size=50g
killall kubectl

# Enable Ingress Addon
echo 'Enabling Ingress addon...'
minikube addons enable ingress
echo

# Deploying K8s Dashboard
echo 'Deploying K8s Dashboard...'
kubectl apply -f https://raw.githubusercontent.com/kubernetes/dashboard/v2.7.0/aio/deploy/recommended.yaml
echo

# Install CRDs
echo 'Installing CRDs...'

kubectl apply -f https://raw.githubusercontent.com/cloudnative-pg/cloudnative-pg/release-1.20/releases/cnpg-1.20.0.yaml

helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo add signoz https://charts.signoz.io

helm repo update

echo
echo "Local K8s cluster up and running..."
