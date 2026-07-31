set shell := ["bash", "-cu"]

default:
    just --list

ps:
    docker ps

up:
    docker compose up

down:
    docker compose down

identityx +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx ;;
      lint) golangci-lint run ./api/identityx/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

univents +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build payssage -d
            docker compose up --build univents ;;
      lint) golangci-lint run ./api/univents/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

payssage +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build payssage ;;
      lint) golangci-lint run ./api/payssage/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

informd +CMD="":
    #!/usr/bin/env bash
    case "{{CMD}}" in
      "")   docker compose up --build identityx -d
            docker compose up --build informd ;;
      lint) golangci-lint run ./api/informd/... ;;
      *)    echo "unknown command: {{CMD}}" && exit 1 ;;
    esac

test:
    cd api/identityx && just test
    cd api/informd && just test
    cd api/payssage && just test
    cd api/univents && just test

generate +SERVICES="identityx informd payssage univents":
    #!/usr/bin/env bash
    for svc in {{SERVICES}}; do
      (cd api/$svc && tygo generate)
    done

lint +NAME="":
    #!/usr/bin/env bash
    case "{{NAME}}" in
      identityx) golangci-lint run ./api/identityx/... ;;
      informd)   golangci-lint run ./api/informd/...   ;;
      payssage)  golangci-lint run ./api/payssage/...  ;;
      univents)  golangci-lint run ./api/univents/...  ;;
      lib)       golangci-lint run ./lib/go/...        ;;
      sdk)       golangci-lint run ./sdk/go/IdentityX/... ./sdk/go/Payssage/... ;;       
      *)         golangci-lint run ./api/identityx/... ./api/informd/... ./api/payssage/... ./api/univents/... ./lib/go/... ./sdk/go/IdentityX/... ./sdk/go/Payssage/... ;;       
    esac
