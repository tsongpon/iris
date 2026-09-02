#!/usr/bin/env bash
# One-time setup for the Cloud Scheduler job that calls the iris Cloud Run
# service (POST /v1/notify) every weekday at 08:00 Asia/Bangkok.
#
# Prerequisites:
#   - The iris Cloud Run service must already be deployed (run
#     scripts/setup-gcloud.sh and the GitHub Actions deploy first).
#   - gcloud authenticated with an account that can manage IAM, service
#     accounts and Cloud Scheduler jobs in the project.
#
# Usage:
#   ./scripts/setup-scheduler.sh PROJECT_ID
#
# The script is idempotent: re-running updates the job if it already exists.
set -euo pipefail

PROJECT_ID="${1:?usage: setup-scheduler.sh PROJECT_ID}"
REGION="asia-southeast1"
SERVICE="iris"
SA_NAME="iris-scheduler"
SCHEDULE="0 8 * * 1-5"                 # weekdays 08:00
TIMEZONE="Asia/Bangkok"
JOB_NAME="iris-notify-daily"

SA_EMAIL="${SA_NAME}@${PROJECT_ID}.iam.gserviceaccount.com"
PROJECT_NUMBER=$(gcloud projects describe "$PROJECT_ID" --format='value(projectNumber)')
SCHEDULER_AGENT="service-${PROJECT_NUMBER}@gcp-sa-cloudscheduler.iam.gserviceaccount.com"

gcloud config set project "$PROJECT_ID"

echo "==> Enabling Cloud Scheduler API"
gcloud services enable cloudscheduler.googleapis.com

echo "==> Checking Cloud Run service '${SERVICE}' exists"
if ! gcloud run services describe "$SERVICE" --region "$REGION" >/dev/null 2>&1; then
  echo "ERROR: Cloud Run service '${SERVICE}' not found in ${REGION}." >&2
  echo "       Deploy the service first, then re-run this script." >&2
  exit 1
fi
SERVICE_URL=$(gcloud run services describe "$SERVICE" --region "$REGION" --format='value(status.url)')
echo "    ${SERVICE_URL}"

echo "==> Creating scheduler service account (idempotent)"
gcloud iam service-accounts create "$SA_NAME" \
  --display-name="Cloud Scheduler for iris" 2>/dev/null || \
  echo "   service account already exists"

echo "==> Granting roles/run.invoker on the service"
gcloud run services add-iam-policy-binding "$SERVICE" \
  --region "$REGION" \
  --member="serviceAccount:${SA_EMAIL}" \
  --role="roles/run.invoker" --quiet >/dev/null

echo "==> Allowing Cloud Scheduler service agent to mint OIDC tokens"
gcloud iam service-accounts add-iam-policy-binding "$SA_EMAIL" \
  --member="serviceAccount:${SCHEDULER_AGENT}" \
  --role="roles/iam.serviceAccountTokenCreator" --quiet >/dev/null

echo "==> Creating or updating scheduler job '${JOB_NAME}'"
SCHEDULER_ARGS=(
  --schedule="${SCHEDULE}"
  --time-zone="${TIMEZONE}"
  --http-method=POST
  --uri="${SERVICE_URL}/v1/notify"
  --oidc-service-account-email="${SA_EMAIL}"
  --oidc-token-audience="${SERVICE_URL}"
  --max-retry-attempts=3
  --min-backoff=10s
  --max-backoff=60s
)
if gcloud scheduler jobs describe "$JOB_NAME" --location "$REGION" >/dev/null 2>&1; then
  gcloud scheduler jobs update http "$JOB_NAME" --location "$REGION" "${SCHEDULER_ARGS[@]}"
else
  gcloud scheduler jobs create http "$JOB_NAME" --location "$REGION" "${SCHEDULER_ARGS[@]}"
fi

cat <<EOF

============================================================
Scheduler job '${JOB_NAME}' is configured.

  Schedule:  ${SCHEDULE} (${TIMEZONE})
  Target:    POST ${SERVICE_URL}/v1/notify
  Auth:      OIDC as ${SA_EMAIL}

Useful commands:

  # Trigger immediately to test (check Teams for the message)
  gcloud scheduler jobs run ${JOB_NAME} --location=${REGION}

  # Inspect recent executions
  gcloud logging read 'resource.type="cloud_scheduler_job" AND
    resource.labels.job_id="${JOB_NAME}" AND severity>=WARNING' \\
    --project=${PROJECT_ID} --limit=20

  # Pause / resume
  gcloud scheduler jobs pause ${JOB_NAME} --location=${REGION}
  gcloud scheduler jobs resume ${JOB_NAME} --location=${REGION}
============================================================
EOF