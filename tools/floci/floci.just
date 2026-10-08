# Managed mode on Floci (local AWS). Imported by the root justfile.
#
# `just identityx managed` deploys the same image the hosted container runs
# as two Lambda functions, wired the way AWS will wire them:
#
#   identityx-api     <- HTTP API (API Gateway v2), RUNTIME_ROLE=api
#   identityx-worker  <- SQS queue identityx-jobs (DLQ identityx-jobs-dlq)
#                     <- EventBridge Scheduler, one schedule per periodic job
#                        (cadence from `identityx schedules`)
#
# Both read their config from Secrets Manager (`just inject identityx`) and
# send email through SES (listed at http://localhost:4566/_aws/ses).

floci_endpoint := "http://localhost:4566"
floci_region := "us-east-1"
floci_account := "000000000000"

# Later files win over earlier ones. Targets Floci.
#
# Seed a service's secret (trieoh/<svc>) from .env + api/<svc>/.env
inject SVC:
    #!/usr/bin/env bash
    set -euo pipefail
    COMPOSE_PROFILES=managed docker compose up --detach --wait floci
    SECRETS_PROVIDER=aws-secretsmanager \
    SECRETS_ID="trieoh/{{SVC}}" \
    AWS_ENDPOINT_URL="{{floci_endpoint}}" \
    AWS_REGION="{{floci_region}}" \
    AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test \
      go run lib/cmd/secrets inject .env "api/{{SVC}}/.env"

# Deploy IdentityX to Floci as api + worker functions (idempotent).
[private]
identityx-floci-deploy:
    #!/usr/bin/env bash
    set -euo pipefail
    export COMPOSE_PROFILES=managed
    aws() {
      docker compose exec -T floci aws --endpoint-url "{{floci_endpoint}}" \
        --region "{{floci_region}}" --output text "$@"
    }
    set -a; source api/identityx/.env; set +a

    image="{{floci_account}}.dkr.ecr.{{floci_region}}.amazonaws.com/trieoh/identityx:dev"
    exec_role="arn:aws:iam::{{floci_account}}:role/identityx-local"

    docker compose up --detach --wait floci
    docker compose build identityx
    just inject identityx

    # Job queue, dead-lettering after 5 failed receives. The visibility
    # timeout outlasts the worker's function timeout.
    aws sqs create-queue --queue-name identityx-jobs-dlq >/dev/null
    dlq_url=$(aws sqs get-queue-url --queue-name identityx-jobs-dlq --query QueueUrl)
    dlq_arn=$(aws sqs get-queue-attributes --queue-url "$dlq_url" --attribute-names QueueArn --query Attributes.QueueArn)
    printf '{"VisibilityTimeout":"180","RedrivePolicy":"{\\"deadLetterTargetArn\\":\\"%s\\",\\"maxReceiveCount\\":\\"5\\"}"}' "$dlq_arn" \
      | aws sqs create-queue --queue-name identityx-jobs --attributes file:///dev/stdin >/dev/null
    queue_url=$(aws sqs get-queue-url --queue-name identityx-jobs --query QueueUrl)
    queue_arn=$(aws sqs get-queue-attributes --queue-url "$queue_url" --attribute-names QueueArn --query Attributes.QueueArn)

    aws ses verify-email-identity --email-address "$EMAIL_FROM" >/dev/null

    # One image, two roles. Only the bootstrap and the managed overrides are
    # function environment; everything else comes from the secret. The
    # worker alone maps 5xx to a failed invocation, so a failed scheduled
    # job is retried by the scheduler (queue batches report per-message
    # failures in a 200 body instead); on the api role a 5xx must stay an
    # HTTP response.
    for role in api worker; do
      fn="identityx-$role"
      extra=""
      [ "$role" = worker ] && extra=',"AWS_LWA_ERROR_STATUS_CODES":"500-599"'
      env_json=$(printf '{"Variables":{"RUNTIME_MODE":"managed","RUNTIME_ROLE":"%s","SECRETS_PROVIDER":"aws-secretsmanager","SECRETS_ID":"trieoh/identityx","AWS_ENDPOINT_URL":"http://floci:4566","JOBS_QUEUE_URL":"%s","EMAIL_TRANSPORT":"ses"%s}}' "$role" "$queue_url" "$extra")
      if aws lambda get-function --function-name "$fn" >/dev/null 2>&1; then
        aws lambda update-function-code --function-name "$fn" --image-uri "$image" >/dev/null
        aws lambda wait function-updated-v2 --function-name "$fn"
        echo "$env_json" | aws lambda update-function-configuration --function-name "$fn" \
          --timeout 60 --memory-size 512 --environment file:///dev/stdin >/dev/null
      else
        echo "$env_json" | aws lambda create-function --function-name "$fn" \
          --package-type Image --code "ImageUri=$image" --role "$exec_role" \
          --timeout 60 --memory-size 512 --environment file:///dev/stdin >/dev/null
      fi
      aws lambda wait function-active-v2 --function-name "$fn"
    done
    api_arn=$(aws lambda get-function-configuration --function-name identityx-api --query FunctionArn)
    worker_arn=$(aws lambda get-function-configuration --function-name identityx-worker --query FunctionArn)

    # Queue -> worker, reporting per-message failures.
    mapping=$(aws lambda list-event-source-mappings --function-name identityx-worker \
      --event-source-arn "$queue_arn" --query 'EventSourceMappings[0].UUID')
    if [ "$mapping" = "None" ] || [ -z "$mapping" ]; then
      aws lambda create-event-source-mapping --function-name identityx-worker \
        --event-source-arn "$queue_arn" --batch-size 10 \
        --function-response-types ReportBatchItemFailures >/dev/null
    fi

    # Scheduler -> worker, one schedule per periodic job. The list is read
    # on fd 3: `docker compose exec` would otherwise drain the loop's stdin.
    while read -r kind seconds <&3; do
      name="identityx-${kind//./-}"
      minutes=$(( seconds / 60 )); [ "$minutes" -lt 1 ] && minutes=1
      unit=minutes; [ "$minutes" -eq 1 ] && unit=minute
      aws scheduler delete-schedule --name "$name" >/dev/null 2>&1 || true
      printf '{"Arn":"%s","RoleArn":"%s","Input":"{\\"kind\\":\\"%s\\"}"}' "$worker_arn" "$exec_role" "$kind" \
        | aws scheduler create-schedule --name "$name" \
            --schedule-expression "rate($minutes $unit)" \
            --flexible-time-window Mode=OFF --target file:///dev/stdin >/dev/null
      echo "schedule $name: rate($minutes $unit)"
    done 3< <(go run ./api/identityx/cmd/identityx schedules)

    # HTTP API -> api function.
    api_id=$(aws apigatewayv2 get-apis --query "Items[?Name=='identityx-managed'].ApiId | [0]")
    if [ "$api_id" = "None" ] || [ -z "$api_id" ]; then
      api_id=$(aws apigatewayv2 create-api --name identityx-managed --protocol-type HTTP --query ApiId)
    fi
    integration_id=$(aws apigatewayv2 get-integrations --api-id "$api_id" \
      --query "Items[?IntegrationUri=='$api_arn'].IntegrationId | [0]")
    if [ "$integration_id" = "None" ] || [ -z "$integration_id" ]; then
      integration_id=$(aws apigatewayv2 create-integration --api-id "$api_id" \
        --integration-type AWS_PROXY --integration-uri "$api_arn" \
        --payload-format-version 2.0 --query IntegrationId)
    fi
    routes=$(aws apigatewayv2 get-routes --api-id "$api_id" --query 'Items[].RouteKey')
    for route in 'ANY /' 'ANY /{proxy+}'; do
      if ! grep -qF "$route" <<<"$routes"; then
        aws apigatewayv2 create-route --api-id "$api_id" --route-key "$route" \
          --target "integrations/$integration_id" >/dev/null
      fi
    done
    stages=$(aws apigatewayv2 get-stages --api-id "$api_id" --query 'Items[].StageName')
    if ! grep -qF '$default' <<<"$stages"; then
      aws apigatewayv2 create-stage --api-id "$api_id" --stage-name '$default' --auto-deploy >/dev/null
    fi
    aws lambda add-permission --function-name identityx-api \
      --statement-id identityx-managed-apigw --action lambda:InvokeFunction \
      --principal apigateway.amazonaws.com \
      --source-arn "arn:aws:execute-api:{{floci_region}}:{{floci_account}}:$api_id/*/*/*" >/dev/null 2>&1 || true

    echo "IdentityX managed endpoint: http://$api_id.execute-api.localhost.floci.io:4566"
    echo "Emails sent through SES:    {{floci_endpoint}}/_aws/ses"
