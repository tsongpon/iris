#!/usr/bin/env bash
# One-time GCP setup for the iris GitHub Actions -> Cloud Run pipeline.
#
# Prerequisites:
#   gcloud installed, authenticated with an account that is Project Owner
#   (or has privileges to enable APIs, create SAs, pools and IAM bindings).
#
# Usage:
#   ./scripts/setup-gcloud.sh PROJECT_ID
#
# After it completes, fill in the GitHub secrets/vars printed at the end.
set -euo pipefail

PROJECT_ID="${1:?usage: setup-gcloud.sh PROJECT_ID}"
REGION="asia-southeast1"
AR_HOST="${REGION}-docker.pkg.dev"
AR_REPO="iris"
GITHUB_OWNER="tsongpon"
GITHUB_REPO="iris"
CI_SA="github-deploy-iris"
RUNTIME_SA="iris-calendar-reader"        # existing SA the service runs as
POOL="github-pool"
PROVIDER="iris"

gcloud config set project "$PROJECT_ID"

echo "==> Enabling APIs"
gcloud services enable iamcredentials.googleapis.com \
  artifactregistry.googleapis.com \
  run.googleapis.com \
  cloudbuild.googleapis.com

echo "==> Ensuring Artifact Registry repo '${AR_REPO}'"
gcloud artifacts repositories describe "$AR_REPO" --region "$REGION" >/dev/null 2>&1 || \
  gcloud artifacts repositories create "$AR_REPO" \
    --repository-format=docker --location="$REGION"

echo "==> Creating CI service account (idempotent)"
gcloud iam service-accounts create "$CI_SA" \
  --display-name="GitHub Actions CI for iris" 2>/dev/null || \
  echo "   service account already exists"

CI_SA_EMAIL="${CI_SA}@${PROJECT_ID}.iam.gserviceaccount.com"
RUNTIME_SA_EMAIL="${RUNTIME_SA}@${PROJECT_ID}.iam.gserviceaccount.com"

echo "==> Granting CI roles"
# Push images
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${CI_SA_EMAIL}" \
  --role="roles/artifactregistry.writer" --quiet >/dev/null
# Create/update Cloud Run services
gcloud projects add-iam-policy-binding "$PROJECT_ID" \
  --member="serviceAccount:${CI_SA_EMAIL}" \
  --role="roles/run.admin" --quiet >/dev/null
# Act as the runtime SA when deploying
gcloud iam service-accounts add-iam-policy-binding "$RUNTIME_SA_EMAIL" \
  --member="serviceAccount:${CI_SA_EMAIL}" \
  --role="roles/iam.serviceAccountUser" --quiet >/dev/null

echo "==> Creating Workload Identity pool and provider"
gcloud iam workload-identity-pools describe "$POOL" \
  --location=global >/dev/null 2>&1 || \
  gcloud iam workload-identity-pools create "$POOL" \
    --location=global --display-name="GitHub Actions pool"

gcloud iam workload-identity-pools providers describe "$PROVIDER" \
  --location=global --workload-identity-pool="$POOL" >/dev/null 2>&1 || \
  gcloud iam workload-identity-pools providers create-oidc "$PROVIDER" \
    --location=global \
    --workload-identity-pool="$POOL" \
    --display-name="iris GitHub OIDC" \
    --issuer-uri="https://token.actions.githubusercontent.com" \
    --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
    --attribute-condition="assertion.repository=='${GITHUB_OWNER}/${GITHUB_REPO}'"

echo "==> Binding pool to CI service account"
gcloud iam service-accounts add-iam-policy-binding "$CI_SA_EMAIL" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/projects/$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')/locations/global/workloadIdentityPools/${POOL}/attribute.repository/${GITHUB_OWNER}/${GITHUB_REPO}" \
  --quiet >/dev/null

PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')

cat <<EOF

============================================================
Setup complete. Now configure the GitHub repository:

Secrets (Settings -> Secrets and variables -> Actions):
  WIF_PROVIDER            projects/${PROJECT_NUMBER}/locations/global/workloadIdentityPools/${POOL}/providers/${PROVIDER}
  WIF_SERVICE_ACCOUNT     ${CI_SA_EMAIL}
  RUNTIME_SERVICE_ACCOUNT ${RUNTIME_SA_EMAIL}
  GOOGLE_CALENDAR_ID      <your calendar id>
  MSTEAM_WEBHOOK_URL      <your Teams webhook url>

Variables (same page, Variables tab):
  GCP_PROJECT_ID    ${PROJECT_ID}
  GCP_REGION        ${REGION}
  AR_HOST           ${AR_HOST}
  AR_REPO           ${AR_REPO}
  CLOUDRUN_SERVICE  iris
  CALENDAR_TIMEZONE Asia/Bangkok
  LOG_LEVEL         info

Reminder: the runtime SA ${RUNTIME_SA_EMAIL} must have access to the
Google Calendar (share the calendar with that email if it hasn't been
granted yet).
============================================================
EOF