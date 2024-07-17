#!/bin/bash

# Define the common parts of the image name
REGISTRY="localhost:5001/karmada"
TAG="v1.11.0-alpha.0-147-gd3adcf68e-dirty"
NAMESPACE="karmada-system"

# List of deployments
DEPLOYMENTS=(
    # "karmada-aggregated-apiserver"
    # "karmada-apiserver"
    "karmada-controller-manager"
    # "karmada-descheduler"
    # "karmada-kube-controller-manager"
    # "karmada-metrics-adapter"
    # "karmada-scheduler"
    # "karmada-search"
    # "karmada-webhook"
)

# Update image for each deployment
for DEPLOYMENT in "${DEPLOYMENTS[@]}"; do
    NEW_IMAGE="${REGISTRY}/${DEPLOYMENT}:${TAG}"
    echo "Updating deployment ${DEPLOYMENT} to use image ${NEW_IMAGE}..."
    
    kubectl -n ${NAMESPACE} set image deployment/${DEPLOYMENT} "*=${NEW_IMAGE}" --record
    
    if [ $? -eq 0 ]; then
        echo "Successfully updated ${DEPLOYMENT}"
    else
        echo "Failed to update ${DEPLOYMENT}"
    fi
done

# Update image for scheduler estimator members
ESTIMATOR_DEPLOYMENTS=(
    "karmada-scheduler-estimator-member1"
    "karmada-scheduler-estimator-member2"
    "karmada-scheduler-estimator-member3"
)

NEW_IMAGE="${REGISTRY}/karmada-scheduler-estimator:${TAG}"

for DEPLOYMENT in "${ESTIMATOR_DEPLOYMENTS[@]}"; do
    echo "Updating deployment ${DEPLOYMENT} to use image ${NEW_IMAGE}..."
    
    kubectl -n ${NAMESPACE} set image deployment/${DEPLOYMENT} "*=${NEW_IMAGE}" --record
    
    if [ $? -eq 0 ]; then
        echo "Successfully updated ${DEPLOYMENT}"
    else
        echo "Failed to update ${DEPLOYMENT}"
    fi
done

echo "All deployments have been updated."
