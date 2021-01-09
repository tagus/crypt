# Crypt
A convenient credential store that securely saved account info.
The `crypt` cli provides simple commands to add, modify and list all stored credentials.

Crypt stores credentials in an aes encrypted file `cryptfile`. The `crypt` binary will
resolve the default crypfile in the following locations in order.
1. `--cryptfile`, `-c` flag
2. `CRYPTFILE` env variable
3. `~/.cryptfile`

The strength of a password is limited by its convention and length.
Easy to remember usually means easier to crack. Committing a long password is hard
and storing the password in plaintext in your pc is a major point of failure.

Crypt will use a SHA256 hash as the key to unlock the `cryptfile`.
The hash will be computed from a password that the user provides.

## Commands
1. `add` adds a credential to the resolved `cryptfile`
2. `delete` deletes the given service from crypt
3. `edit` edits individual fields for the given service
4. `export` exports the `cryptfile` to plain json (mostly for debugging, be careful)
5. `show` show information about a service while copying the pwd to your clipboard
6. `info` metadata about the resolved `cryptfile`
7. `ls` list all stored services in the resolved `cryptfile`
8. `new` creates a new `cryptfile`
9. `pwd` copies the pwd of the specified service in your clipboard

## Structure
Each service is stored as a `credential` with the given structure.

```
credential {
  service: string
  email: string
  username: string
  password: string
  description: string
}
```