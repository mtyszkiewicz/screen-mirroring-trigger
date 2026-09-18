# screen-mirroring-trigger

Runs a command when it detects AirPlay screen-mirroring traffic on macOS. I use it to switch my receiver to the TV profile.

It watches `awdl0` for packets destined for port `7000`. The first packet triggers the command; subsequent packets do nothing until there have been 30 seconds of silence. Nothing runs when traffic stops.

This is a traffic-based guess, not an AirPlay session API. Unrelated traffic can trigger it, and a long enough pause can cause it to fire again.

## Install

Requires macOS, Go, and sudo access for packet capture. Uses macOS’s bundled libpcap.

Edit `dev.mtyszkiewicz.screen-mirroring-trigger.plist` and replace my command with yours:

```xml
<string>-command</string>
<string>/usr/bin/curl -X PUT http://10.205.0.5:8001/profile?name=tv</string>
```

Then run:

```sh
./install.sh
```

The installer builds the program, copies it to `/usr/local/bin`, and starts a system LaunchDaemon. It runs as root and starts automatically at boot.

To change the command or silence timeout, edit the plist in this directory and rerun `./install.sh`.

## Run manually

Without installing the daemon:

```sh
go build -o screen-mirroring-trigger .
sudo ./screen-mirroring-trigger -command "/usr/bin/echo Mirroring detected"
```

| Flag | Description |
| --- | --- |
| `-command` | Command to run when traffic starts. Required. |
| `-silence` | Seconds without packets before another trigger is allowed. Default: `30`. |

Commands are split on whitespace, not interpreted by a shell. Quotes inside the command, pipes, redirects, and variable expansion are not supported. For anything more involved, use an executable script and pass its absolute path.

The command runs synchronously, so keep it short.

## Logs

```sh
tail -f /tmp/screen-mirroring-trigger.log /tmp/screen-mirroring-trigger.err
```

## Uninstall

```sh
sudo ./uninstall.sh
```

Stops the daemon and removes the installed binary, plist, and logs.