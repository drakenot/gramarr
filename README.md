#gramarr
## A [Radarr](https://github.com/Radarr/Radarr) and [Sonarr](https://github.com/Sonarr/Sonarr) Telegram Bot featuring user authentication/level access.

## Features

### Sonarr

- Search for TV Shows by name.
- Pick which seasons you want to download.

### Radarr

- Search for Movies by name.
- Choose which quality profile you want to download.

## Requirements

- A running instance of Radarr
- A running instance of Sonarr V3 (preview)

### If running from source

- [Go](https://golang.org/)

## Configuration

- Copy the `config.json.template` file to `config.json` and set-up your configuration;

#### If running from downloaded binaries
- Put the copied `config.json` alongside with the binary downloaded from [releases](https://github.com/drakenot/gramarr/releases);

## Running it


### From source

```bash
$ go get github.com/drakenot/gramarr
$ cd $GOPATH/src/github.com/drakenot/gramarr
$ go get
$ go run .
```

### From release

Just [download](https://github.com/drakenot/gramarr/releases/latest) the respective binary for your System.

*Obs: Don't forget to put the `config.json` in the same folder as the binary file.*

