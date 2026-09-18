# screen-mirroring-trigger

Executes a command when AirPlay screen mirroring starts on macOS.

## How it works

Monitors the `awdl0` (Apple Wireless Direct Link) interface for AirPlay connections to port 7000, which is used for screen mirroring. When a stream **starts**, it executes a configurable command **once**; further packets while the stream is active are ignored so manual receiver changes (e.g. volume) are not overwritten. After `silence` seconds with no packets the stream is considered ended, and the next packet triggers the command again.

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

Optionally adjust the silence threshold (seconds without packets before a stream is considered ended):
```xml
<string>-silence</string>
<string>30</string>
```
AirPlay mirroring sends packets continuously while active, so silence reliably means the stream stopped. 30s tolerates brief WiFi blips without re-triggering.

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
sudo screen-mirroring-trigger -silence 60 -command "echo 'Screen mirroring started!'"
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