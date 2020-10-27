# Roku Remote CLI

A go cli for interacting with Roku devices on your home network.

## Setup

```shell
some command
```

## Usage

```shell
Using External Control Protocol your Roku device is exposed over a local area network.

Using SSDP (Simple Service Discovery Protocol) and asimple RESTful API

Usage:
  roku-remote [command]

Available Commands:
  find        Find Roku Remotes on your local network.
  help        Help about any command
  live        A glimpse into what is currently playing on the Roku.
  send        Send an action to your Roku Device.

Flags:
      --config string   config file (default is $HOME/.roku-remote.yaml)
  -h, --help            help for roku-remote
      --host string     host ip of the roku

Use "roku-remote [command] --help" for more information about a command.
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
├── cmd
│   ├── root.go
│   ├── find.go
│   ├── live.go
│   ├── send.go
│   └── helpers.go
├── roku
│   ├── config.go
│   ├── instructions.go
│   └── roku.go
├── LICENSE
├── README.md
├── go.mod
├── go.sum
└── main.go

3 directories, 15 files

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
