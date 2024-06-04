# crypt

The strength of a password is limited by its convention and length.
Easy to remember usually means easier to crack. Committing a long password is hard
and storing the password in plaintext in your pc is a major point of failure.

crypt is a convenient credential store that securely saved account info.
The `crypt` cli provides simple commands to add, modify and list all stored credentials.

Crypt stores credentials in a sqlite db where sensitive data (password, security questions etc.)
is encrypted using 32-bit AES encryption.

The cli will look for the sqlite db file in the following locations in order.
1. `--crypt-db`, `-c` flag
2. `CRYPT_DB` env variable
3. `./.crypt.db`
4. `~/.crypt.db`

Multiple crypts can be stores in a single crypt db however the cli allows you to access one at a 
time by specifying the name with the `--crypt-name` flag.

## commands
1. `add` adds a credential to the resolved `crypt` db
2. `archive` archives the given service from crypt
3. `info` shows metadata about the current crypt
4. `list` exports the `cryptfile` to plain json (mostly for debugging, be careful)
5. `show` show information about a service while copying the pwd to your clipboard

The `--help` flag can be used in conjunction with any command to get more details
about usage and optional flags.
