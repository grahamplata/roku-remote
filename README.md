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
  roku [command]

app
  active      Show the currently active application on your Roku.
  add         Add applications to your Roku.
  launch      Launch applications on your Roku.
  list        List the applications on your Roku.

device
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
  ▸ http://192.168.10.95:8060/
    http://192.168.10.122:8060/
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
