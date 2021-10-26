#!/bin/bash

# Copyright The Shipwright Contributors
# 
# SPDX-License-Identifier: Apache-2.0

set -euo pipefail

echo "Building container image"

echo "Adding io.shipwright.vcs-ref label with value: ${GITHUB_SHA}"

# Using defaults, this pushes to:
# - quay.io/shipwright/git:latest
# - quay.io/shipwright/shipwright-build-controller:latest
KO_DOCKER_REPO="${IMAGE_HOST}/${IMAGE_NAMESPACE}" GOFLAGS="${GO_FLAGS}" ko resolve \
  --base-import-paths \
  --tags "${TAG}" \
  --image-label "io.shipwright.vcs-ref=${GITHUB_SHA}" \
  --platform=all -R -f deploy/ > release.yaml

KO_DOCKER_REPO="${IMAGE_HOST}/${IMAGE_NAMESPACE}" GOFLAGS="${GO_FLAGS} -tags=pprof_enabled" ko resolve \
  --base-import-paths \
  --tags "${TAG}-debug" \
  --image-label "io.shipwright.vcs-ref=${GITHUB_SHA}" \
  --platform=all -R -f deploy/ > release-debug.yaml

# Bundle the sample cluster build strategies, remove namespace strategies first
find samples/buildstrategy -type f -print0 | xargs -0 grep -l "kind: BuildStrategy" | xargs rm -f
KO_DOCKER_REPO="${IMAGE_HOST}/${IMAGE_NAMESPACE}" ko resolve -R -f samples/buildstrategy/ > sample-strategies.yaml
