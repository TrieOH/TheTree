set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

build service:
    docker compose -f api/{{service}}/compose.yml build

build *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml build; \
    done

run service:
    docker compose -f api/{{service}}/compose.yml up --build

run *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml up --build; \
    done

rund service:
    docker compose -f api/{{service}}/compose.yml up -d --build

rund *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml up -d --build; \
    done

down service:
    docker compose -f api/{{service}}/compose.yml down

down *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml down; \
    done

down-v service:
    docker compose -f api/{{service}}/compose.yml down -v

down-v *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml down -v; \
    done

logs service:
    docker compose -f api/{{service}}/compose.yml logs -f

logs *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml logs -f; \
    done

restart service:
    docker compose -f api/{{service}}/compose.yml restart

restart *services:
    for s in {{services}}; do \
        docker compose -f api/"$s"/compose.yml restart; \
    done

all *args:
    for f in api/*/compose.yml; do \
        docker compose -f "$f" {{args}}; \
    done
