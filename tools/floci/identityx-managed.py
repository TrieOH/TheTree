#!/usr/bin/env python3
"""Build and deploy IdentityX as a local Floci Lambda + HTTP API."""

import json
import os
import subprocess
import sys
from pathlib import Path


ROOT = Path(__file__).resolve().parents[2]
IMAGE = "000000000000.dkr.ecr.us-east-1.amazonaws.com/trieoh/identityx:dev"
FUNCTION = "identityx"
API_NAME = "identityx-managed"
REGION = "us-east-1"
ENDPOINT = "http://localhost:4566"


def compose(*args: str) -> None:
    env = os.environ | {"COMPOSE_PROFILES": "managed"}
    subprocess.run(
        ["docker", "compose", *args],
        cwd=ROOT,
        env=env,
        check=True,
        text=True,
    )


def aws(*args: str, input_text: str | None = None) -> dict:
    """Run the AWS CLI bundled in Floci's compat image."""
    command = [
        "docker", "compose", "exec", "-T", "floci", "aws",
        "--endpoint-url", ENDPOINT,
        "--region", REGION,
        "--output", "json",
        *args,
    ]
    result = subprocess.run(
        command,
        cwd=ROOT,
        env=os.environ | {"COMPOSE_PROFILES": "managed"},
        check=True,
        text=True,
        input=input_text,
        stdout=subprocess.PIPE,
    )
    return json.loads(result.stdout or "{}")


def dotenv(path: Path) -> dict[str, str]:
    values: dict[str, str] = {}
    for raw_line in path.read_text().splitlines():
        line = raw_line.strip()
        if not line or line.startswith("#") or "=" not in line:
            continue
        key, value = line.split("=", 1)
        value = value.strip()
        if len(value) >= 2 and value[0] == value[-1] and value[0] in "\"'":
            value = value[1:-1]
        values[key.strip()] = value
    return values


def function_environment() -> dict[str, str]:
    values = dotenv(ROOT / "api/identityx/.env")
    # Postgres image initialization variables belong to the DB container.
    for key in ("POSTGRES_DB", "POSTGRES_USER", "POSTGRES_PASSWORD"):
        values.pop(key, None)
    return values


def main() -> None:
    compose("up", "--detach", "--wait", "identityx-db", "floci")
    compose("build", "identityx")

    env_json = json.dumps({"Variables": function_environment()})
    functions = aws("lambda", "list-functions").get("Functions", [])
    if not any(item.get("FunctionName") == FUNCTION for item in functions):
        aws(
            "lambda", "create-function",
            "--function-name", FUNCTION,
            "--package-type", "Image",
            "--code", f"ImageUri={IMAGE}",
            "--role", "arn:aws:iam::000000000000:role/identityx-local",
            "--timeout", "30",
            "--memory-size", "512",
            "--environment", "file:///dev/stdin",
            input_text=env_json,
        )
    else:
        aws("lambda", "update-function-code", "--function-name", FUNCTION,
            "--image-uri", IMAGE)
        aws("lambda", "wait", "function-updated-v2", "--function-name", FUNCTION)
        aws("lambda", "update-function-configuration", "--function-name", FUNCTION,
            "--timeout", "30", "--memory-size", "512",
            "--environment", "file:///dev/stdin", input_text=env_json)
    aws("lambda", "wait", "function-active-v2", "--function-name", FUNCTION)
    function = aws("lambda", "get-function-configuration", "--function-name", FUNCTION)
    function_arn = function["FunctionArn"]

    apis = aws("apigatewayv2", "get-apis").get("Items", [])
    api = next((item for item in apis if item.get("Name") == API_NAME), None)
    if api is None:
        api = aws("apigatewayv2", "create-api", "--name", API_NAME,
                  "--protocol-type", "HTTP")
    api_id = api["ApiId"]

    integrations = aws("apigatewayv2", "get-integrations", "--api-id", api_id).get("Items", [])
    integration = next((item for item in integrations
                        if item.get("IntegrationUri") == function_arn), None)
    if integration is None:
        integration = aws(
            "apigatewayv2", "create-integration",
            "--api-id", api_id,
            "--integration-type", "AWS_PROXY",
            "--integration-uri", function_arn,
            "--payload-format-version", "2.0",
        )
    integration_id = integration["IntegrationId"]
    routes = aws("apigatewayv2", "get-routes", "--api-id", api_id).get("Items", [])
    route_keys = {item.get("RouteKey") for item in routes}
    for route_key in ("ANY /", "ANY /{proxy+}"):
        if route_key not in route_keys:
            aws("apigatewayv2", "create-route", "--api-id", api_id,
                "--route-key", route_key,
                "--target", f"integrations/{integration_id}")
    stages = aws("apigatewayv2", "get-stages", "--api-id", api_id).get("Items", [])
    if not any(item.get("StageName") == "$default" for item in stages):
        aws("apigatewayv2", "create-stage", "--api-id", api_id,
            "--stage-name", "$default", "--auto-deploy")

    try:
        aws("lambda", "add-permission",
            "--function-name", FUNCTION,
            "--statement-id", "identityx-managed-apigw",
            "--action", "lambda:InvokeFunction",
            "--principal", "apigateway.amazonaws.com",
            "--source-arn", f"arn:aws:execute-api:{REGION}:000000000000:{api_id}/*/*/*")
    except subprocess.CalledProcessError:
        # The permission already exists on repeated deploys.
        pass

    print(f"IdentityX Managed endpoint: http://{api_id}.execute-api.localhost.floci.io:4566")


if __name__ == "__main__":
    try:
        main()
    except (subprocess.CalledProcessError, KeyError, json.JSONDecodeError) as error:
        print(f"IdentityX managed deployment failed: {error}", file=sys.stderr)
        sys.exit(1)
