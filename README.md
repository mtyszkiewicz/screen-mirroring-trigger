# screen-mirroring-trigger

Executes a command when AirPlay screen mirroring starts on macOS.

## How it works

Monitors the `awdl0` (Apple Wireless Direct Link) interface for AirPlay connections to port 7000, which is used for screen mirroring. When a connection is detected, it executes a configurable command with a cooldown period to prevent duplicate triggers.

## Requirements

- macOS with Go installed
- Root/sudo access (required for packet capture)
- libpcap (pre-installed on macOS)

## Setup

**1. Edit the plist before installing**

Open `dev.mtyszkiewicz.screen-mirroring-trigger.plist` and change the command:

```xml
<string>-command</string>
<string>curl -X PUT http://10.205.0.5:8001/profile?name=tv</string>
```

This example switches my Onkyo amplituner to TV mode. **Replace with your own command.**

Optionally adjust the cooldown (seconds between triggers):
```xml
<string>-cooldown</string>
<string>30</string>
```

**2. Install**

```bash
chmod +x install.sh uninstall.sh
./install.sh
```

This will:
1. Build the binary
2. Install it to `/usr/local/bin/screen-mirroring-trigger`
3. Install and load a LaunchDaemon to run at startup

## Manual Usage

You can also run it manually:

```bash
sudo screen-mirroring-trigger -cooldown 60 -command "echo 'Screen mirroring started!'"
```

## Logs

View logs:
```bash
tail -f /tmp/screen-mirroring-trigger.log
```

## Uninstall

```bash
./uninstall.sh
```