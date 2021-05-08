# fixtures-cli

---

This is a cli that generates initial data for testing features.

## Install

Make sure you are in the app folder where you can find `Makefile`. To install the cli
use the following command:

```shell
make install_fixtures
```

Run the following command to make sure that it is installed right:

```shell
fixtures
```

You should see the following response.

```shell
Fixtures can generate testing data a user with a few campaigns alongside with a few templates. Also about 
hundreds of subscribers in a few segments

Usage:
  fixtures [command]

Available Commands:
  clear       clear deletes all the stuff created with full-init
  help        Help about any command
  init        init creates testing data

Flags:
  -h, --help              help for fixtures
  -p, --password string   Password for the user with fixtures
  -s, --secret string     Secret for api key for the user with fixtures
  -u, --username string   Username for the user with fixtures

Use "fixtures [command] --help" for more information about a command.
```

## Commands

### init

This command creates new user with api key, along with a segment with a hundred subscribers
and five campaigns with one same template.

```shell
fixtures init -u="username" -p="password" -s="secret"
```

The `--username` or short `-u` flag is required for the username, the `--password`, `-p` and
`--secret`, `-s` flags are optional by default random uuid is generated for both.

## clear

This command deletes all data for specified user. Only the data for users that are created
with the cli can be cleared.

```shell
fixtures init -u="username"
```

For This command only the `--username`, `-u` flag is required. Then you will be asked for 
the user's password.

```shell
password for "username": "password"
```