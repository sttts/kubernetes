#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

time make WHAT=test/e2e/e2e.test

kubeconfig="${WORKDIR}/.kcp2/data/admin.kubeconfig"

## each of the single-cluster contexts passes normal e2e
for context in "admin" "user" "other"; do
  _output/local/bin/linux/amd64/e2e.test --e2e-verify-service-account=false --ginkgo.focus "Watchers" \
    --kubeconfig "${kubeconfig}" --context "${context}" \
    --provider local
done

# the multi-cluster list/watcher works with all the single-cluster contexts
_output/local/bin/linux/amd64/e2e.test --e2e-verify-service-account=false --ginkgo.focus "MultiClusterWorkflow" \
  --kubeconfig "${kubeconfig}" --context "admin" \
  --provider kcp \
  --kcp-multi-cluster-kubeconfig "${kubeconfig}" --kcp-multi-cluster-context "cross-cluster" \
  --kcp-secondary-kubeconfig "${kubeconfig}" --kcp-secondary-context "user" \
  --kcp-tertiary-kubeconfig "${kubeconfig}" --kcp-tertiary-context "other" \
  --kcp-clusterless-kubeconfig "${kubeconfig}" --kcp-clusterless-context "admin"
