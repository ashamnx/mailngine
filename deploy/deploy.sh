#!/usr/bin/env bash
set -euo pipefail

# Mailngine Deployment Script — Zero-Downtime with Automatic Rollback
# Usage: ./deploy/deploy.sh [full|api|dashboard|website|nginx|rollback]

SERVER="root@178.128.208.168"
REMOTE_DIR="/opt/mailngine"
RELEASES_DIR="$REMOTE_DIR/releases"
SHARED_DIR="$REMOTE_DIR/shared"
CURRENT_LINK="$REMOTE_DIR/current"
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
RELEASE_NAME="$(date +%Y%m%d-%H%M%S)"
HEALTH_URL="http://localhost:8090/health"
HEALTH_RETRIES=15
HEALTH_INTERVAL=2
KEEP_RELEASES=5

cd "$PROJECT_DIR"

component="${1:-full}"

log()  { echo "==> $*"; }
warn() { echo "WARNING: $*" >&2; }
die()  { echo "FATAL: $*" >&2; exit 1; }

# --- Remote helpers ---

remote() {
    ssh "$SERVER" "$@"
}

# Ensure the release directory structure exists (idempotent).
ensure_structure() {
    remote "mkdir -p $RELEASES_DIR $SHARED_DIR"
    # Migrate legacy flat .env to shared/ on first run
    remote "[ -f $REMOTE_DIR/.env ] && mv $REMOTE_DIR/.env $SHARED_DIR/.env || true"
}

# Returns the release path that /opt/mailngine/current points to.
get_current_release() {
    remote "readlink -f $CURRENT_LINK 2>/dev/null || echo ''"
}

# Atomically swap the current symlink to a new release.
swap_symlink() {
    local target="$1"
    log "Swapping symlink → $target"
    remote "ln -sfn $target $CURRENT_LINK"
}

# Health check: curl /health up to HEALTH_RETRIES times.
health_check() {
    local service_name="$1"
    log "Health check ($service_name)..."
    for i in $(seq 1 "$HEALTH_RETRIES"); do
        if remote "curl -sf $HEALTH_URL" >/dev/null 2>&1; then
            log "Health check passed ($service_name)"
            return 0
        fi
        sleep "$HEALTH_INTERVAL"
    done
    warn "Health check FAILED ($service_name) after $((HEALTH_RETRIES * HEALTH_INTERVAL))s"
    return 1
}

# Rollback: swap symlink to previous release and restart services.
rollback() {
    local previous="$1"
    if [ -z "$previous" ] || [ ! -d "$previous" ] 2>/dev/null; then
        die "No previous release to roll back to"
    fi
    warn "Rolling back to $previous"
    swap_symlink "$previous"
    remote "systemctl restart mailngine-api" || true
    remote "systemctl restart mailngine-worker" || true
    # Wait for the rolled-back version to become healthy
    if health_check "rollback"; then
        warn "Rollback complete. Services running from $previous"
    else
        die "Rollback health check also failed — manual intervention required"
    fi
    exit 1
}

# Remove old releases, keeping the most recent $KEEP_RELEASES.
cleanup_releases() {
    log "Cleaning up old releases (keeping last $KEEP_RELEASES)..."
    remote "ls -1dt $RELEASES_DIR/*/ 2>/dev/null | tail -n +$((KEEP_RELEASES + 1)) | xargs rm -rf 2>/dev/null || true"
}

# --- Build steps ---

build_api() {
    log "Building Go binaries (linux/amd64)..."
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/linux/server ./cmd/server
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/linux/worker ./cmd/worker
    log "Go binaries: OK"
}

build_dashboard() {
    log "Building Angular dashboard..."
    cd web && npx ng build --configuration=production 2>&1 | tail -3
    cd "$PROJECT_DIR"
    log "Dashboard: OK"
}

build_website() {
    log "Building Qwik website (static)..."
    cd website && npm run build.static 2>&1 | tail -5
    cd "$PROJECT_DIR"
    log "Website: OK"
}

# --- Deploy steps ---

deploy_api() {
    ensure_structure

    local release_path="$RELEASES_DIR/$RELEASE_NAME"
    local previous
    previous=$(get_current_release)

    log "Creating release $RELEASE_NAME..."
    remote "mkdir -p $release_path"

    # Upload binaries to the new release directory (services still running old version)
    log "Uploading binaries..."
    scp bin/linux/server bin/linux/worker "$SERVER:$release_path/"
    remote "chmod +x $release_path/server $release_path/worker"

    # Deploy systemd services (env is managed on the server, never overwritten)
    log "Deploying systemd services..."
    if ! remote "test -s $SHARED_DIR/.env"; then
        die "Production .env not found at $SHARED_DIR/.env — create it on the server first"
    fi
    scp deploy/production/mailngine-api.service deploy/production/mailngine-worker.service "$SERVER:/etc/systemd/system/"
    remote "systemctl daemon-reload"

    # Atomic symlink swap
    swap_symlink "$release_path"

    # Restart API first, verify health before touching worker
    log "Restarting API..."
    remote "systemctl restart mailngine-api"
    if ! health_check "api"; then
        warn "API failed to start with new release"
        rollback "$previous"
    fi

    # Restart worker (asynq drains in-flight tasks via ShutdownTimeout)
    log "Restarting worker (draining in-flight tasks)..."
    remote "systemctl restart mailngine-worker"
    remote "systemctl enable mailngine-api mailngine-worker"

    cleanup_releases
    log "API deployed: $RELEASE_NAME"
}

deploy_dashboard() {
    log "Deploying Angular dashboard..."
    remote "mkdir -p $REMOTE_DIR/dashboard"
    rsync -az --delete web/dist/web/browser/ "$SERVER:$REMOTE_DIR/dashboard/"
    log "Dashboard deployed"
}

deploy_website() {
    log "Deploying Qwik website..."
    remote "mkdir -p $REMOTE_DIR/website"
    rsync -az --delete website/dist/ "$SERVER:$REMOTE_DIR/website/"
    log "Website deployed"
}

deploy_nginx() {
    log "Deploying nginx configs..."
    scp deploy/production/security-headers.conf "$SERVER:/etc/nginx/conf.d/security-headers.conf"
    scp deploy/production/security.txt "$SERVER:$REMOTE_DIR/security.txt"
    scp deploy/production/mailngine-api.nginx "$SERVER:/etc/nginx/sites-available/mailngine-api"
    scp deploy/production/mailngine-website.nginx "$SERVER:/etc/nginx/sites-available/mailngine-website"
    remote "ln -sf /etc/nginx/sites-available/mailngine-api /etc/nginx/sites-enabled/ && ln -sf /etc/nginx/sites-available/mailngine-website /etc/nginx/sites-enabled/ && nginx -t && systemctl reload nginx"
    log "Nginx deployed and reloaded"
}

# --- Manual rollback command ---

do_rollback() {
    local current
    current=$(get_current_release)
    log "Current release: $current"

    # Find the previous release (second newest)
    local previous
    previous=$(remote "ls -1dt $RELEASES_DIR/*/ 2>/dev/null | sed -n '2p' | tr -d '/'")
    if [ -z "$previous" ]; then
        die "No previous release found to roll back to"
    fi

    log "Rolling back from $current to $previous"
    swap_symlink "$previous"
    remote "systemctl restart mailngine-api"
    health_check "api" || die "Rollback failed — API unhealthy"
    remote "systemctl restart mailngine-worker"
    log "Rollback complete: now running $previous"
}

# --- Entrypoint ---

case "$component" in
    full)
        log "=== Mailngine Deploy: full ==="
        build_api
        build_dashboard
        build_website
        deploy_api
        deploy_dashboard
        deploy_website
        deploy_nginx
        ;;
    api)
        log "=== Mailngine Deploy: api ==="
        build_api
        deploy_api
        ;;
    dashboard)
        log "=== Mailngine Deploy: dashboard ==="
        build_dashboard
        deploy_dashboard
        remote "systemctl reload nginx"
        ;;
    website)
        log "=== Mailngine Deploy: website ==="
        build_website
        deploy_website
        remote "systemctl reload nginx"
        ;;
    nginx)
        log "=== Mailngine Deploy: nginx ==="
        deploy_nginx
        ;;
    rollback)
        log "=== Mailngine Rollback ==="
        do_rollback
        ;;
    *)
        echo "Usage: $0 [full|api|dashboard|website|nginx|rollback]"
        exit 1
        ;;
esac

log "=== Deploy complete ==="
