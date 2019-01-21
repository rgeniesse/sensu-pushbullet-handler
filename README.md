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

_NOTE: It is recommended to use [`state_change_only` filter][3] with this handler. If not used, you will get a notification after each check run during a non-0 exit status._

Example Sensu Go definition:

```json
{
    "api_version": "core/v2",
    "type": "Handler",
    "metadata": {
        "namespace": "default",
        "name": "pushbullet"
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
[3]: https://docs.sensu.io/sensu-go/5.1/reference/filters/#handling-state-change-only
