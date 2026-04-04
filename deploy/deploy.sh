#!/usr/bin/env bash
set -euo pipefail

# Hello Mail Deployment Script
# Usage: ./deploy/deploy.sh [full|api|dashboard|website]

SERVER="root@178.128.208.168"
REMOTE_DIR="/opt/mailngine"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"

cd "$PROJECT_DIR"

component="${1:-full}"

echo "=== Hello Mail Deploy: $component ==="

build_api() {
    echo "Building Go binaries (linux/amd64)..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/linux/server ./cmd/server
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/linux/worker ./cmd/worker
    echo "Go binaries: OK"
}

build_dashboard() {
    echo "Building Angular dashboard..."
    cd web && npx ng build --configuration=production 2>&1 | tail -3
    cd "$PROJECT_DIR"
    echo "Dashboard: OK"
}

build_website() {
    echo "Building Qwik website (static)..."
    cd website && npm run build.static 2>&1 | tail -5
    cd "$PROJECT_DIR"
    echo "Website: OK"
}

deploy_api() {
    echo "Deploying API binaries..."
    ssh "$SERVER" "mkdir -p $REMOTE_DIR/bin"
    scp bin/linux/server bin/linux/worker "$SERVER:$REMOTE_DIR/bin/"
    ssh "$SERVER" "chmod +x $REMOTE_DIR/bin/server $REMOTE_DIR/bin/worker"

    echo "Deploying env and systemd services..."
    scp deploy/production/.env.production "$SERVER:$REMOTE_DIR/.env"
    scp deploy/production/mailngine-api.service deploy/production/mailngine-worker.service "$SERVER:/etc/systemd/system/"

    ssh "$SERVER" "systemctl daemon-reload && systemctl restart mailngine-api && systemctl restart mailngine-worker && systemctl enable mailngine-api mailngine-worker"
    echo "API deployed and restarted"
}

deploy_dashboard() {
    echo "Deploying Angular dashboard..."
    ssh "$SERVER" "mkdir -p $REMOTE_DIR/dashboard"
    rsync -az --delete web/dist/web/browser/ "$SERVER:$REMOTE_DIR/dashboard/"
    echo "Dashboard deployed"
}

deploy_website() {
    echo "Deploying Qwik website..."
    ssh "$SERVER" "mkdir -p $REMOTE_DIR/website"
    rsync -az --delete website/dist/ "$SERVER:$REMOTE_DIR/website/"
    echo "Website deployed"
}

deploy_nginx() {
    echo "Deploying nginx configs..."
    scp deploy/production/mailngine-api.nginx "$SERVER:/etc/nginx/sites-available/mailngine-api"
    scp deploy/production/mailngine-website.nginx "$SERVER:/etc/nginx/sites-available/mailngine-website"
    ssh "$SERVER" "ln -sf /etc/nginx/sites-available/mailngine-api /etc/nginx/sites-enabled/ && ln -sf /etc/nginx/sites-available/mailngine-website /etc/nginx/sites-enabled/ && nginx -t && systemctl reload nginx"
    echo "Nginx deployed and reloaded"
}

case "$component" in
    full)
        build_api
        build_dashboard
        build_website
        deploy_api
        deploy_dashboard
        deploy_website
        deploy_nginx
        ;;
    api)
        build_api
        deploy_api
        ;;
    dashboard)
        build_dashboard
        deploy_dashboard
        ssh "$SERVER" "systemctl reload nginx"
        ;;
    website)
        build_website
        deploy_website
        ssh "$SERVER" "systemctl reload nginx"
        ;;
    nginx)
        deploy_nginx
        ;;
    *)
        echo "Usage: $0 [full|api|dashboard|website|nginx]"
        exit 1
        ;;
esac

echo "=== Deploy complete ==="
