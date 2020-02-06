# gramarr
## A [Radarr](https://github.com/Radarr/Radarr) and [Sonarr](https://github.com/Sonarr/Sonarr) Telegram Bot featuring user authentication/level access.

![Grammar](https://extraimage.com/images/2020/02/06/gramarr.jpg)

## Features

### Sonarr

- Search for TV Shows by name.
- Pick which seasons you want to download.
- Choose which quality and language profile you want to download.

### Radarr

- Search for Movies by name.
- Choose which quality profile you want to download.

## Requirements

- A running instance of Radarr
- A running instance of Sonarr V3 (preview)

### If running from source

- [Go](https://golang.org/)

### If running from docker

- [Docker](https://docker.io)
- [Docker Compose](https://docs.docker.com/compose/)

## Configuration

- Copy the `config.json.template` file to `config.json` and set-up your configuration;

#### If running from downloaded binaries
- Put the copied `config.json` alongside with the binary downloaded from [releases](https://github.com/alcmoraes/gramarr/releases);

## Running it

### From Docker

```bash
$ docker-compose up -d
```

### From source

```bash
$ go get github.com/alcmoraes/gramarr
$ cd $GOPATH/src/github.com/alcmoraes/gramarr
$ go get
$ go run .
```

### From release

Just [download](https://github.com/alcmoraes/gramarr/releases/latest) the respective binary for your System.

*Obs: Don't forget to put the `config.json` in the same folder as the binary file.*

## TODO

- **Package oriented**: Reorganize the project.
