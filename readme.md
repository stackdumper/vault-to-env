# vault-to-env
VTE reads secrets from Hashicorp Vault and outputs them as environment variables.

<br />

## Usage

#### Read
Read secrets and output them as env variables

```
Usage:
  vte read [flags]

Flags:
  -h, --help                 help for read
      --lease-duration int   adjust secret lease duration
      --save-leases          save secret leases
      --vars strings         list of vars to read

Global Flags:
      --address string      Vault address (default "http://vault.admin.e4f.cc")
      --auth-data strings   Vault auth data
      --auth-path string    Vault auth path
```

<br />

## Examples

```bash
# create a secret (kv v2)
$ vault kv put secret/hello foo=world excited=yes
Key              Value
---              -----
created_time     2019-05-14T14:51:38.856822Z
deletion_time    n/a
destroyed        false
version          1


# read secrets and output them into stdout
# if you use KV v2, append /data/ before secret path
$ vte read
    --auth-path /auth/userpass/login/tester
    --auth-data password=tester
    --vars FOO=/secret/data/hello#foo
    --vars EXCITED=/secret/data/hello#excited
export FOO="world"
export EXCITED="yes"
```

<br />

## License
[MIT](./license)
