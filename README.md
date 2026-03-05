# Lullaby - Orange Pi 5 Manager Bot

A simple, production-ready Telegram bot to manage your Orange Pi 5 (reboot and power off) via Telegram.

## Features
- 🔄 Reboot host
- 🔌 Power off host
- 🛂 Security: Only responds to authorized user IDs
- 🏗️ K8s Ready: Designed to run in a k3s cluster with host privileges

## Tech Stack
- **Go**: Core application logic
- **Telebot**: Telegram Bot API wrapper
- **Docker**: Multi-stage ARM64 build
- **Kubernetes**: Deployment manifests for k3s

## Setup

1.  **Get a Telegram Bot Token**: Create a bot via [@BotFather](https://t.me/BotFather).
2.  **Get your User ID**: Use [@userinfobot](https://t.me/userinfobot) to find your Telegram ID.
3.  **Configure Secrets**:
    Update `deploy/k8s/secret.yaml` with your token and user ID.
4.  **Build & Deploy**:
    ```bash
    # Build docker image
    make docker-build
    
    # Apply manifests
    kubectl apply -f deploy/k8s/secret.yaml
    kubectl apply -f deploy/k8s/deployment.yaml
    ```

## Security Notice
This bot requires `privileged: true` and mounts the host root `/` to be able to execute `reboot` and `poweroff` on the host machine. Ensure you only allow trusted IDs in `AUTHORIZED_USER_IDS`.
