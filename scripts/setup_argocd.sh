#!/bin/bash

NAMESPACE="argocd"
ARGOCD_VERSION="stable"

kubectl create namespace $NAMESPACE

kubectl apply -n $NAMESPACE -f https://raw.githubusercontent.com/argoproj/argo-cd/$ARGOCD_VERSION/manifests/install.yaml

echo "Waiting for ArgoCD server to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n $NAMESPACE

ARGOCD_PASSWORD=$(kubectl -n $NAMESPACE get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)

curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd
rm argocd-linux-amd64

kubectl port-forward svc/argocd-server -n $NAMESPACE 8080:443 &

sleep 5

# Login to ArgoCD
argocd login localhost:8080 --username admin --password $ARGOCD_PASSWORD --insecure

echo "ArgoCD setup complete!"
echo "ArgoCD UI is available at: https://localhost:8080"
echo "Username: admin"
echo "Password: $ARGOCD_PASSWORD"