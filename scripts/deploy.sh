#!/usr/bin/env bash
set -euox pipefail
SCRIPT_DIR=$(dirname "$0")
echo "${SCRIPT_DIR}"

project=$(gcloud secrets versions access latest --secret="project-id")
if [[ -z "${project}" ]]; then
  echo -n "need project"
  exit 1
fi
echo "${project}"

gcloud run deploy go-publisher \
  --image gcr.io/"${project}"/go-publisher:latest \
  --platform managed \
  --project "${project}" \
  --region asia-northeast1 \
  --allow-unauthenticated \
  --set-env-vars PUB_PROJECT="${project}"

# MEMO: use all user access
#  --allow-unauthenticated \
