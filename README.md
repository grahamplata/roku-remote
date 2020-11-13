# Roku Remote CLI

A go cli for interacting with Roku devices on your home network.

## Setup

Build it yourself

```shell
git clone git@github.com:grahamplata/roku-remote.git
cd roku-remote
go build -o roku-remote -v .
```

## Usage

```shell
Using SSDP (Simple Service Discovery Protocol) access your Roku's RESTful API

Usage:
  roku-remote [command]

Available Commands:
  apps        List the applications on your Roku.
  find        Find Roku Remotes on your local network.
  help        Help about any command
  live        Status of the Roku media player.
  send        Send an action to your Roku Device.

Flags:
      --config string   config file (default is $HOME/.roku-remote.yaml)
  -h, --help            help for roku-remote
      --host string     host ip of the roku

Use "roku-remote [command] --help" for more information about a command.
```

### find

```shell
roku-remote find

Use the arrow keys to navigate: ↓ ↑ → ←
? Select a default Roku from your network:
  ▸ http://192.168.10.95:8060/
    http://192.168.10.122:8060/
```

### help

```shell
roku-remote find

Using SSDP (Simple Service Discovery Protocol) access your Roku's RESTful API

Usage:
  roku-remote [command]

Available Commands:
  apps        List the applications on your Roku.
  find        Find Roku Remotes on your local network.
  help        Help about any command
  live        Status of the Roku media player.
  send        Send an action to your Roku Device.

Flags:
      --config string   config file (default is $HOME/.roku-remote.yaml)
  -h, --help            help for roku-remote
      --host string     host ip of the roku

Use "roku-remote [command] --help" for more information about a command.
```

### live

```shell
roku-remote live

Playing: Plex - Stream for Free
Watched: 23m10.633
```

### send

```shell
roku-remote send select

Sent select action to Roku
```

#### Available Actions

```shell
Navigation:  left, right, up, down, select, home, search
Keyboard:    backspace, enter
Remote:      fwd, rev, play, replay, tuner, poweroff, channeldown, channelup, volumedown, volumeup,  info, mute, replay
Inputs:      HDMI1, HDMI2, HDMI3, HDMI4
```

### apps

```shell
roku-remote send apps

The following apps are installed on this Roku
---
ID: 31012 Name: FandangoNOW Movies & TV Version: 1.9.34
ID: 12 Name: Netflix Version: 4.2.81121084
ID: 837 Name: YouTube Version: 1.0.80000314
ID: 13 Name: Prime Video Version: 11.4.2020092816
ID: 50539 Name: Twitch Version: 1.0.16
ID: 13535 Name: Plex - Stream for Free Version: 6.6.16
ID: 151908 Name: The Roku Channel Version: 2.9.42
ID: 34376 Name: ESPN Version: 3.2.2020092400
ID: 552944 Name: Roku Tips & Tricks Version: 1.2.55
```

## Configuration

```shell
# .roku-remote.yaml
roku:
  host: http://192.168.1.1:8060/
```

## Tree

```shell
.
├── LICENSE
├── README.md
├── cmd
│   ├── apps.go
│   ├── find.go
│   ├── helpers.go
│   ├── live.go
│   ├── root.go
│   └── send.go
├── go.mod
├── go.sum
├── main.go
└── roku
    ├── config.go
    ├── instructions.go
    └── roku.go

8 directories, 22 files

```

## Notes

- [Roku documentation](https://developer.roku.com/docs/developer-program/debugging/external-control-api.md)
- Rokus use External Control Protocol (ECP)
  - Enables a Roku device to be controlled over a local area network by providing a number of external control services.
  - The Roku devices offering these external control services are discoverable using SSDP (Simple Service Discovery Protocol).
  - ECP is a simple RESTful API that can be accessed by programs in virtually any programming environment.
- On a Mac
  - [To avoid the apple network warning, you need to build the executable file once and codesign it.](https://apple.stackexchange.com/a/393721)
  - `go build -o roku-remote main.go && codesign -s - roku-remote # build the executable file and codesign`
