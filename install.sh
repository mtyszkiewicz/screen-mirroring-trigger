#!/bin/bash
set -e

echo "Building screen-mirroring-trigger..."
go build -o screen-mirroring-trigger screen-mirroring-trigger.go

echo "Installing binary to /usr/local/bin..."
sudo cp screen-mirroring-trigger /usr/local/bin/
sudo chown root:wheel /usr/local/bin/screen-mirroring-trigger
sudo chmod 755 /usr/local/bin/screen-mirroring-trigger

echo "Installing LaunchDaemon..."
sudo cp dev.mtyszkiewicz.screen-mirroring-trigger.plist /Library/LaunchDaemons/
sudo chown root:wheel /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist
sudo chmod 644 /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist

# Reload the service if already loaded so the new binary/plist takes effect.
# bootout is best-effort: it's a no-op (and non-fatal) on a fresh install where
# nothing is loaded yet. Without this, bootstrap fails with I/O error (code 5)
# when reinstalling over an already-running daemon.
echo "Reloading service..."
sudo launchctl bootout system /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist 2>/dev/null || true
sudo launchctl bootstrap system /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist

echo ""
echo "✓ screen-mirroring-trigger installed successfully!"
echo ""
echo "Default settings:"
echo "  Silence: 30s"
echo "  Command: curl -X PUT http://10.205.0.5:8001/profile?name=tv"
echo ""
echo "To customize, edit /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist"
echo "Then reload: sudo launchctl bootout system /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist"
echo "             sudo launchctl bootstrap system /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist"
echo ""
echo "Check logs: tail -f /tmp/screen-mirroring-trigger.log"
