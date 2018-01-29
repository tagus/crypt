# Crypt
A convenient credential store that securely saved account info.
This `CLI` provides simple commands to add, modify and list all stored credentials.

Crypt stores credentials in an encrypted `encoded` file.

The strength of a password is limited by its convention and length.
Easy to remember usually means easier to crack. Committing a long password is hard
and storing the password in plaintext in your pc is a major point of failure.

Crypt will use a 256 MD5 hash as the key. The hash will be computed from
a key string from the users.

## Commands

1. Add
2. Edit
3. List