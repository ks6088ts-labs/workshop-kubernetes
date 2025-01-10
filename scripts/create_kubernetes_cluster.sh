#!/bin/sh

# Quickstart: Deploy an Azure Kubernetes Service (AKS) cluster using Azure CLI
# https://learn.microsoft.com/azure/aks/learn/quick-kubernetes-deploy-cli

# Variables
LOCATION="japaneast"
RANDOM_SUFFIX=$(openssl rand -hex 4)
RESOURCE_GROUP_NAME="rg-workshop-kubernetes-$RANDOM_SUFFIX"
AKS_CLUSTER_NAME="workshop-kubernetes-$RANDOM_SUFFIX"

# Create resource group
az group create \
  --name "$RESOURCE_GROUP_NAME" \
  --location "$LOCATION" \
  --verbose

# Create AKS cluster
az aks create \
  --resource-group "$RESOURCE_GROUP_NAME" \
  --name "$AKS_CLUSTER_NAME" \
  --enable-managed-identity \
  --node-count 1 \
  --generate-ssh-keys \
  --verbose

# Get credentials
az aks get-credentials \
  --resource-group "$RESOURCE_GROUP_NAME" \
  --name "$AKS_CLUSTER_NAME" \
  --verbose

# Verify
kubectl get nodes
