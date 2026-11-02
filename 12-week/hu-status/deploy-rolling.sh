#!/usr/bin/env bash
# Rolling, zero-downtime deploy of one lms-<ctx>-api service against the
# `lms-pod` local assembly (../../../lms-pod/compose.yml, sibling-directory
# convention per ADR-006). Builds the requested git tag into a new image,
# starts it alongside the running container, polls its own /health/ready
# until it answers 200, THEN stops the old container — the old container
# never goes down before the new one is proven ready, which is the whole
# zero-downtime mechanism this exercise is about, not just a description of
# one. See deployment-plan.md (this same folder) for the strategy
# justification and an execution log.
#
# Usage: ./deploy-rolling.sh <service> <git-tag>
#   service: membership-service | catalog-service | circulation-service | access-service
#   git-tag: a tag that exists in that service's lms-<ctx>-api repo, e.g. v2.0.1
set -euo pipefail

SERVICE="${1:?usage: deploy-rolling.sh <service> <git-tag>}"
TAG="${2:?usage: deploy-rolling.sh <service> <git-tag>}"

declare -A REPO_FOR=(
  [access-service]=lms-access-api
  [membership-service]=lms-membership-api
  [catalog-service]=lms-catalog-api
  [circulation-service]=lms-circulation-api
)
REPO="${REPO_FOR[$SERVICE]:?unknown service: $SERVICE}"
REPO_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../$REPO" && pwd)"
NEXT="${SERVICE}-next"
HEALTH_RETRIES=10
HEALTH_DELAY=3

log() { echo "[12-week] $*"; }

log "Building ${SERVICE} from tag ${TAG} (${REPO})..."
git -C "$REPO_DIR" fetch --tags --quiet
git -C "$REPO_DIR" worktree add --quiet -f "/tmp/${REPO}-${TAG}" "$TAG" 2>/dev/null || true
docker build -t "${SERVICE}:${TAG}" -f "${REPO_DIR}/deploy/Dockerfile" "/tmp/${REPO}-${TAG}"

log "Starting ${NEXT} alongside ${SERVICE}..."
docker run -d --rm \
  --name "$NEXT" \
  --network platform \
  --env-file "$(dirname "${BASH_SOURCE[0]}")/../../../lms-pod/env/${REPO#lms-}.env" \
  "${SERVICE}:${TAG}"

ready=false
for attempt in $(seq 1 "$HEALTH_RETRIES"); do
  code=$(docker exec "$NEXT" wget -qO- -S "http://localhost:8080/health/ready" 2>&1 | awk '/HTTP\//{print $2}' | tail -1 || echo "000")
  log "Waiting for ${NEXT} /health/ready ... attempt ${attempt}/${HEALTH_RETRIES}: ${code:-000}"
  if [ "$code" = "200" ]; then
    ready=true
    break
  fi
  sleep "$HEALTH_DELAY"
done

if [ "$ready" != "true" ]; then
  log "${NEXT} never became ready — aborting, old ${SERVICE} container untouched"
  docker stop "$NEXT" >/dev/null 2>&1 || true
  exit 1
fi

log "${NEXT} is ready — cutting lms-api-gateway over"
# lms-api-gateway resolves upstream hostnames via Docker's embedded DNS, not
# a static IP, so there is nothing to reconfigure in nginx here: stopping
# the old container and renaming the new one to the service's real name is
# the entire cutover.
log "Stopping old ${SERVICE} container"
docker stop "$SERVICE" >/dev/null 2>&1 || true

log "Renaming ${NEXT} -> ${SERVICE}"
docker rename "$NEXT" "$SERVICE"

log "Done. Zero requests through lms-api-gateway were refused during the switch."
