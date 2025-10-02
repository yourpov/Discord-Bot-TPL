#!/usr/bin/env bash
# Discord Bot Management for Ubuntu/Debian/CentOS with systemctl support (this script run in /root/discord-bot)
#  Usage: bash./manage.sh [command] (options: build, restart, status, logs, screen, full)
#    Much Love - YourPov (@uhhhwhatever)

# Server Requirements:
#  - Ram: 512MB+ recommended (bot typically uses 50-200MB)
#  - Disk: 100MB+ free (bot binary ~11MB + dependencies ~50MB)
#  - OS: 
set -eEuo pipefail

command -v go >/dev/null 2>&1 || die "go is  not installed"
command -v systemctl >/dev/null 2>&1 || die "systemctl is not available"
command -v screen >/dev/null 2>&1 || die "screen is not installed"

BOT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
cd "$BOT_DIR"

SERVICE="discord-bot"
BIN_NAME="template"
SCREEN_NAME="botbuild"
BUILD_LOG="$BOT_DIR/build.log"
RUN_LOG="$BOT_DIR/run.log"

die() { echo "❌ ERROR: $*" >&2; exit 1; }

log_info() { echo "ℹ️  $*"; }
log_success() { echo "✅ $*"; }
log_warning() { echo "⚠️  $*"; }
log_error() { echo "❌ $*"; }
log_activity() { echo "🔄 $*"; }
log_bot() { echo "🤖 $*"; }
log_build() { echo "🔨 $*"; }
log_screen() { echo "🖥️  $*"; }

write_status() {
  local status="$1"
  local timestamp=$(date '+%Y-%m-%d %H:%M:%S')
  echo "[$timestamp] $status" >> "$BOT_DIR/status.log"
}

get_bot_status() {
  if systemctl is-active --quiet "$SERVICE" 2>/dev/null; then
    echo "running"
  elif systemctl is-failed --quiet "$SERVICE" 2>/dev/null; then
    echo "failed"
  else
    echo "stopped"
  fi
}

get_screen_status() {
  if screen -list 2>/dev/null | grep -q "[.]$SCREEN_NAME"; then
    echo "active"
  else
    echo "inactive"
  fi
}

build_bot() {
  clear
  log_build "Starting Discord Bot Build Process"
  write_status "Build started"
  echo
  
  if [ ! -f "go.mod" ]; then
    log_error "go.mod not found. Are you in a Go project directory?"
    write_status "Build failed - no go.mod"
    exit 1
  fi
  
  log_info "Found go.mod - Go project detected"
  
  if [ -f "$BIN_NAME" ]; then
    log_activity "Removing old binary: $BIN_NAME"
    rm -f "$BIN_NAME"
  fi
  
  log_build "Compiling Go source code..."
  if ! go build -v -o "$BIN_NAME" 2>&1 | tee "$BUILD_LOG"; then
    log_error "Build failed! Check $BUILD_LOG for details"
    write_status "Build failed - compilation error"
    echo
    echo "Last few lines of build log:"
    tail -n 5 "$BUILD_LOG" 2>/dev/null || echo "No build log available"
    exit 1
  fi
  
  chmod +x "$BIN_NAME"
  log_success "Bot compiled successfully: $BIN_NAME"
  log_info "Binary is executable and ready to run"
  write_status "Build completed successfully"
}

restart_service() {
  clear
  log_activity "Restarting Discord Bot service: $SERVICE"
  write_status "Service restart initiated"
  echo
  
  log_info "Stopping existing service..."
  if systemctl stop "$SERVICE" 2>/dev/null; then
    log_success "Service stopped successfully"
  else
    log_warning "Service was not running or failed to stop"
  fi
  
  sleep 1
  
  log_activity "Starting service..."
  if systemctl restart "$SERVICE"; then
    log_success "Bot service restarted successfully"
    write_status "Bot restarted and running"
  else
    log_error "Failed to restart bot service"
    write_status "Service restart failed"
    exit 1
  fi
  
  sleep 2
  check_and_report_status
}

start_service() {
  clear
  log_activity "Starting Discord Bot service: $SERVICE"
  write_status "Service start initiated"
  echo
  
  if systemctl start "$SERVICE"; then
    log_success "Bot service started successfully"
    write_status "Bot started and initializing"
  else
    log_error "Failed to start bot service"
    write_status "Service start failed"
    exit 1
  fi
  
  log_info "Waiting for bot to initialize..."
  sleep 3
  
  check_and_report_status
}

check_and_report_status() {
  echo
  log_info "Checking bot status..."
  
  local bot_status=$(get_bot_status)
  local screen_status=$(get_screen_status)
  
  case $bot_status in
    "running")
      log_success "Bot connected to Discord"
      log_bot "Ready to process commands"
      write_status "Bot connected and ready"
      log_info "Discord connection established"
      log_info "Command handlers loaded"
      log_info "Bot monitoring channels"
      ;;
    "failed")
      log_error "Bot went offline (service failed)"
      write_status "Bot offline - service failed"
      echo "Recent error logs:"
      journalctl -u "$SERVICE" -n 3 --no-pager --since "5 minutes ago" 2>/dev/null | grep -i "error\|fail\|exception" | head -3 || echo "  No recent errors found"
      ;;
    "stopped")
      log_warning "Bot is offline (service stopped)"
      write_status "Bot offline - service stopped"
      ;;
  esac
}

screen_attach_or_create() {
  clear
  log_screen "Managing screen session: $SCREEN_NAME"
  write_status "Screen session management started"
  echo
  
  if screen -list 2>/dev/null | grep -q "[.]$SCREEN_NAME"; then
    log_info "Existing screen session found"
    log_success "Bot screen session is active"
    log_screen "Attaching to running session..."
    write_status "Attached to existing screen session"
    echo
    echo "Screen Commands:"
    echo "   Detach: Press Ctrl+A then D"
    echo "   Kill:   screen -S $SCREEN_NAME -X quit"
    echo
    exec screen -r "$SCREEN_NAME"
  else
    log_info "No existing screen session found"
    
    if [ ! -x "$BIN_NAME" ]; then
      log_warning "Bot binary not found or not executable - building first..."
      build_bot
    fi
    
    log_activity "Creating new screen session..."
    log_screen "Launching bot in background..."
    
    if screen -S "$SCREEN_NAME" -d -m bash -lc "cd '$BOT_DIR'; ./'$BIN_NAME' >> '$RUN_LOG' 2>&1"; then
      log_success "Screen session created successfully"
      log_bot "Bot is now running in background"
      log_info "Bot process started and logging to $RUN_LOG"
      write_status "Screen session created - bot running"
      
      sleep 2
      
      log_success "Bot initialization complete"
      log_bot "Discord bot is now online"
      
      echo
      echo "Screen Management:"
      echo "   Attach: screen -r $SCREEN_NAME"
      echo "   Detach: Ctrl+A then D"
      echo "   Kill:   screen -S $SCREEN_NAME -X quit"
      echo "   Logs:   tail -f $RUN_LOG"
      
    else
      log_error "Screen session failed to create"
      write_status "Screen session creation failed"
      exit 1
    fi
  fi
}

show_status() {
  clear
  echo "Discord Bot Status Dashboard"
  echo "══════════════════════════════════════════════════════════════════════"
  echo
  
  local bot_status=$(get_bot_status)
  local screen_status=$(get_screen_status)
  
  echo "Current Status:"
  case $bot_status in
    "running")
      log_success "Bot is online and running"
      log_bot "Connected to Discord"
      log_info "Ready to process commands"
      ;;
    "failed")
      log_error "Bot is offline (service failed)"
      log_warning "Discord connection lost"
      ;;
    "stopped")
      log_warning "Bot is offline (service stopped)"
      log_info "Not connected to Discord"
      ;;
  esac
  
  echo
  echo "Screen Session:"
  if [ "$screen_status" = "active" ]; then
    log_success "Screen session '$SCREEN_NAME' is running"
    log_info "Bot process is in background"
  else
    log_warning "No screen session found"
    log_info "Bot not running in screen"
  fi
  
  echo
  echo "Recent Activity:"
  echo "─────────────────────────────────────────────────────────────────"
  if [ -f "$BOT_DIR/status.log" ]; then
    tail -n 8 "$BOT_DIR/status.log" | while IFS= read -r line; do
      echo "  $line"
    done
  else
    echo "  No activity log found"
  fi
  
  echo
  echo "Bot Runtime Output:"
  echo "─────────────────────────────────────────────────────────────────"
  if [ -f "$RUN_LOG" ] && [ -s "$RUN_LOG" ]; then
    echo "  Recent bot messages:"
    tail -n 5 "$RUN_LOG" | while IFS= read -r line; do
      [ -n "$line" ] && echo "    $line"
    done
  else
    echo "No runtime log available yet"
  fi
  
  echo
  echo "System Diagnostics:"
  echo "─────────────────────────────────────────────────────────────────"
  local recent_errors=$(journalctl -u "$SERVICE" -n 5 --no-pager --since "10 minutes ago" 2>/dev/null | grep -i "error\|fail\|exception" | head -3)
  if [ -n "$recent_errors" ]; then
    echo "  Recent system errors:"
    echo "$recent_errors" | while IFS= read -r line; do
      echo "    $line"
    done
  else
    log_success "No recent system errors detected"
  fi
  
  echo
  echo "Information:"
  echo "─────────────────────────────────────────────────────────────────"
  echo "  Bot binary: $BIN_NAME $([ -f "$BIN_NAME" ] && echo "(exists)" || echo "(missing)")"
  echo "  Service: $SERVICE"
  echo "  Screen: $SCREEN_NAME"
  echo "  Logs: $RUN_LOG"
}

show_logs_follow() {
  clear
  echo "=== Live Logs (Ctrl+C to exit) ==="
  echo
  journalctl -u "$SERVICE" -f
}

full_flow() {
  clear
  echo "=== Full Flow ==="
  echo
  [ -f "$BUILD_LOG" ] && tail -n 40 "$BUILD_LOG" || echo "No build.log yet"
  echo
  [ -f "$RUN_LOG" ] && tail -n 40 "$RUN_LOG" || echo "No run.log yet"

  echo
  echo "Build..."
  build_bot

  echo
  echo "Restart service..."
  systemctl stop "$SERVICE" || true
  start_service

  echo
  echo "Recent logs:"
  journalctl -u "$SERVICE" -n 20 --no-pager || true

  echo
  echo "Screen session..."
  screen_attach_or_create
}

case "${1:-help}" in
  restart) restart_service ;;
  logs)    show_logs_follow ;;
  status)  show_status ;;
  screen)  screen_attach_or_create ;;
  build)   build_bot ;;
  full|"") full_flow ;;
  *)
  echo "=== Bot Manager ==="

    echo
    cat <<EOF
Usage: $0 [OPTION]

Commands:
  restart → Restart the bot service
  logs    → View live bot logs (Ctrl+C to exit)
  status  → Check bot status and recent activity
  screen  → Create/attach to bot screen session
  build   → Build the bot from source
  full    → Complete deployment (build → restart → screen)

Examples:
  $0 status    # Quick status check
  $0 full      # Full deployment workflow
  $0 screen    # Manage bot in screen session
EOF
    exit 1
    ;;
esac
