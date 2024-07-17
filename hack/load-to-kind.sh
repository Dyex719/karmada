#!/bin/bash

# Define the target images
TARGETS=(
    # "karmada-aggregated-apiserver"
    "karmada-controller-manager"
    # "karmada-scheduler"
    # "karmada-descheduler"
    # "karmada-webhook"
    # "karmada-agent"
    "karmada-scheduler-estimator"
    # "karmada-interpreter-webhook-example"
    # "karmada-search"
    # "karmada-operator"
    # "karmada-metrics-adapter"
)

# Define the registry and tag
REGISTRY="localhost:5001/karmada"
TAG="v1.11.0-alpha.0-147-gd3adcf68e-dirty"

# Define the kind cluster name
KIND_CLUSTER_NAME="karmada-host"

# Load each image into the kind cluster
for TARGET in "${TARGETS[@]}"; do
    IMAGE="${REGISTRY}/${TARGET}:${TAG}"
    echo "Loading image ${IMAGE} into kind cluster ${KIND_CLUSTER_NAME}..."
    kind load docker-image "${IMAGE}" --name "${KIND_CLUSTER_NAME}"
done

echo "All images have been loaded into the kind cluster."

