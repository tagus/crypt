# Crypt
A convenient credential store that securely saved account info.
This `CLI` provides simple commands to add, modify and list all stored credentials.

Crypt stores credentials in an encrypted `encoded` file, `~/.cryptfile`

The strength of a password is limited by its convention and length.
Easy to remember usually means easier to crack. Committing a long password is hard
and storing the password in plaintext in your pc is a major point of failure.

Crypt will use a SHA256 hash as the key to unlock the `.cryptfile`.
The hash will be computed from a password provided that the user provides.

## Commands
1. `ls` prints the name of all stored credentials
1. `add` enables the user to add a credential to the `.cryptfile`
3. `show` prints the fields of the specified *service*

## Structure
Each credential is stored is a `Credential` struct that has the following structure.

```go
type Credential struct {
  Service     string
  Email       string
  Username    string
  Password    string
  Description string
}
```