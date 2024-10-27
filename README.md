# crypt

The strength of a password is limited by its convention and length.
Easy to remember usually means easier to crack. Committing a long password is hard
and storing the password in plaintext in your pc is a major point of failure.

crypt is a convenient credential store that securely saved account info.
The `crypt` cli provides simple commands to add, modify and list all stored credentials.

Crypt stores credentials in a sqlite db where sensitive data (password, security questions etc.)
is encrypted using 32-bit AES encryption. This means that the metadata about service accounts is
still stored in plaintext.

The cli will look for the sqlite db file in the following locations in order.
1. `--crypt-db`, `-c` flag
2. `CRYPT_DB` env variable
3. `./.crypt.db`
4. `~/.crypt.db`

Multiple crypts can be stores in a single crypt db however the cli allows you to access one at a 
time by specifying the name with the `--crypt-name` flag.

## usage

The `--help` flag can be used in conjunction with any command to get more details
about usage and optional flags.

```sh
crypt is CLI application to securely store your credentials
so that you don't have to worry about remembering all of your
internet accounts.

Crypt stores your credentials in a sqlite db that will encrypt private details
(e.g. password, secret keys, security questions) with a master password so that 
they cannot be read even if the sqlite file is examined.

The db file can be specified using the following methods listed here in decreasing priority.

        1. db flag
        2. CRYPT_DB env variable
        3. ./.crypt.db
        4. ~/.crypt.db

Usage:
  crypt [command]

Available Commands:
  add         add a service to crypt
  archive     deletes the given service from the crypt db
  completion  Generate the autocompletion script for the specified shell
  edit        edit fields for the given service
  help        Help about any command
  info        information about your crypt file
  list        list stored credentials
  show        show information about a service

Flags:
  -c, --crypt-db string     the crypt db location
  -n, --crypt-name string   the crypt name (default "main")
  -h, --help                help for crypt
      --init                whether to initialize the crypt db file if it doesn't exist
  -v, --verbose             print out any debug logs
      --version             version for crypt

Use "crypt [command] --help" for more information about a command.
```
