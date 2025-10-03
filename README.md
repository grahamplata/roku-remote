# Roku Remote CLI

A go cli for interacting with Roku devices on your home network.

## Features
- Discover Roku devices on your network via SSDP
- Interactive device control with keyboard shortcuts
- Launch and manage applications
- View device information and media player status
- Multiple device support with easy switching

## Installation

### Homebrew
```shell
# Using Homebrew
brew tap grahamplata/tap
brew install roku-remote
# Or if you already have the tap
brew install grahamplata/tap/roku-remote
```

### Binary Releases
Download pre-built binaries from the [releases page](https://github.com/grahamplata/roku-remote/releases).

### From Source
```shell
git clone https://github.com/grahamplata/roku-remote.git
cd roku-remote
go build -o roku-remote .
```

## Roku Setup

> Note: to use third-party apps to control your Roku device you must now [enable it](https://support.roku.com/en-gb/article/217288467#section-2)

1. Use the directional pad on your Roku TV remote to scroll down and select Settings
2. Navigate to System > Advanced system settings
3. Select Control by mobile apps
4. Select from the following settings based on the level of access you need:
  - `Limited` - Control by mobile apps is restricted to text input, app launches and accessing your activity within the app itself. The app can only control devices that are part of your Wi-Fi network
  - `Enabled` - Your Roku device can always be controlled by mobile apps, but it will only respond to commands from devices that are connected to the same local network
  - `Permissive` - Any device within and outside your network could potentially send all commands to your Roku device

## Usage

```shell
# Find devices (first time setup)
roku-remote find

# Control device interactively
roku-remote control

# Launch Netflix
roku-remote apps launch netflix

# Check what's currently running
roku-remote apps active
```

### Help

```shell
Using SSDP (Simple Service Discovery Protocol) access your Roku's RESTful API

Usage:
  roku [command]

app
  active      Show the currently active application on your Roku.
  add         Add applications to your Roku.
  launch      Launch applications on your Roku.
  list        List the applications on your Roku.

device
  control     Control a Roku device via keyboard
  describe    Describes the currently selected Roku
  find        Find Roku Remotes on your local network.
  live        Status of the Roku media player.
  send        Send an action to your Roku Device.
  switch      Switch the default Roku device.

Additional Commands:
  help        Help about any command
  completion  Generate the autocompletion script for the specified shell

Flags:
      --config string   config file (default is $HOME/.roku-remote.yaml)
  -h, --help            help for roku
      --host string     host ip of the roku

Use "roku [command] --help" for more information about a command.
```

### find

```shell
roku-remote find

Use the arrow keys to navigate: ↓ ↑ → ←
? Select a default Roku from your network:
  ▸ http://192.168.10.95
    http://192.168.10.122
```

## Configuration
The CLI stores device information in `~/.roku-remote.yaml`. You can manually edit this file or use the `find` and `switch` commands to manage devices.

## Notes

- [Roku documentation](https://developer.roku.com/docs/developer-program/debugging/external-control-api.md)
- Rokus use External Control Protocol (ECP)
  - Enables a Roku device to be controlled over a local area network by providing a number of external control services.
  - The Roku devices offering these external control services are discoverable using SSDP (Simple Service Discovery Protocol).
  - ECP is a simple RESTful API that can be accessed by programs in virtually any programming environment.
- On a Mac
  - [To avoid the apple network warning, you need to build the executable file once and codesign it.](https://apple.stackexchange.com/a/393721)
  - `go build -o roku-remote main.go && codesign -s - roku-remote # build the executable file and codesign`

## Troubleshooting

### "No Roku device configured" error
Run `roku-remote find` to discover and set a default device.

### "Limited mode" error
Your Roku is in restricted mode. Press the Home button 5 times quickly on your physical remote, or use `roku-remote apps active` to check the current app.

### Device not found
Ensure your computer and Roku are on the same Wi-Fi network.

### Interactive Control
Use `roku-remote device control` for keyboard-based control:

- `q`: Quit
- `p`: Power
- `+/-`: Volume up/down
- `m`: Mute
- Arrow keys: Navigate
- `Enter`: Select
- `Space`: Play/Pause
- `b`: Back
- `h`: Home
