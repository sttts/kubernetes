#!/usr/bin/env bash

set -o nounset
set -o errexit
set -o pipefail

# This helper git rm's all non-root OWNERS files. Our fork only uses
# the root OWNERS file updated to include only members of the
# openshift org, and all other OWNERS files will be flagged by a
# downstream CI check.

ROOT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")/.."; pwd -P)"
find . -mindepth 2 -type f -name 'OWNERS' | grep -ve '^./vendor' | grep -ve '^./openshift-hack' | xargs -I {} git rm {}
git rm "${ROOT_PATH}/OWNERS_ALIASES"
