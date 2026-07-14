set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

up:
    docker compose up

down:
    docker compose down

identityx:
    docker compose up --build identityx

univents:
    docker compose up --build univents

payssage:
    docker compose up --build payssage

informd:
    docker compose up --build informd

generate +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      (cd api/$svc && tygo generate)
    done
