#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

# hack/update-vendor.sh doesn't play nicely with the circular dependency
# introduced by the openshift shared libs (api, client-go, library-go,
# apiserver-library-go) and will end up setting staging repos to versions
# instead of paths. Provided an internal kube version via KUBE_VERSION
# (e.g. KUBE_VERSION=v0.19.0) this script will reset the staging repos back to
# using paths.

KUBE_VERSION="${KUBE_VERSION:-}"
if [[ -z "${KUBE_VERSION}" ]]; then
  echo >&2 "KUBE_VERSION is required"
  exit 1
fi

grep -e '^\sk8s.io.*=>' go.mod |\
  grep "${KUBE_VERSION}" |\
  sed -e 's/=>.*//' -e 's/^[[:space:]]*//' -e 's/[[:space:]]*$//' |\
  xargs -I {} sed -i -e 's+'{}' =>.*+'{}' => ./staging/src/'{}'+' go.mod
