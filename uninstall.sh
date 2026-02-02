#!/bin/bash
set -e

echo "Stopping service..."
sudo launchctl bootout system /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist 2>/dev/null || true

echo "Removing LaunchDaemon..."
sudo rm -f /Library/LaunchDaemons/dev.mtyszkiewicz.screen-mirroring-trigger.plist

echo "Removing binary..."
sudo rm -f /usr/local/bin/screen-mirroring-trigger

echo "Removing logs..."
rm -f /tmp/screen-mirroring-trigger.log /tmp/screen-mirroring-trigger.err

echo "✓ screen-mirroring-trigger has been uninstalled."
