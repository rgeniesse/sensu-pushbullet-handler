# Sensu Go Pushbullet Plugin
TravisCI: [![TravisCI Build Status](https://travis-ci.com/rgeniesse/sensu-pushbullet-handler.svg?branch=master)](https://travis-ci.com/rgeniesse/sensu-pushbullet-handler)

A Sensu handler plugin for sending notifications via [Pushbullet][2]

## Installation

Download the latest version of the sensu-pushbullet-handler from [releases][1],
or create an executable script from this source.

From the local path of the sensu-pushbullet-handler repository:

```
go build -o /usr/local/bin/sensu-pushbullet-handler main.go
```

## Configuration

Example Sensu Go definition:

```json
{
    "api_version": "core/v2",
    "type": "CHANGEME",
    "metadata": {
        "namespace": "default",
        "name": "CHANGEME"
    },
    "spec": {
        "type": "pipe",
        "command": "sensu-pushbullet-handler",
        "timeout": 10,
        "env_vars": [
            "PUSHBULLET_API_KEY=your_api_key_here"
        ],
        "filters": [
            "is_incident",
            "not_silenced",
            "state_change_only"
        ],
        "runtime_assets": [ "sensu-pushbullet-handler" ]
    }
}
```

## Usage Examples

Help:

```
The Sensu Go handler for Pushover

Usage:
  sensu-pushbullet-handler [flags]

Flags:
  -a, --api.token string   Pushbullet API app token, use default from PUSHBULLET_APP_TOKEN env var
  -h, --help               help for sensu-pushbullet-handler
```

## Contributing

See https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md

[1]: https://github.com/CHANGEME/sensu-CHANGEME/releases
[2]: https://www.pushbullet.com/
