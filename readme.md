# vault-to-env

VTE reads secrets from Hashicorp Vault and outputs them as environment variables into stdout or file.

<br />

## Usage

#### Read
Read secrets and output them as env variables

```
Usage:
  vte read [flags]

Flags:
  -h, --help           help for read
      --out string
      --vars strings

Global Flags:
      --address string      Vault address (default "http://localhost:8200")
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
$ go run main.go
    --auth-path /auth/userpass/login/tester
    --auth-data password=tester
      read
        --vars FOO=/secret/data/hello#foo
        --vars EXCITED=/secret/data/hello#excited
export FOO="world"
export EXCITED="yes"


# read secrets and output them to file
$ ➜ go run main.go
    --auth-path /auth/userpass/login/tester
    --auth-data password=tester
      read
        --vars FOO=/secret/data/hello#foo
        --vars EXCITED=/secret/data/hello#excited
        --out /tmp/hello.sh

# get generated file contents
$ cat /tmp/hello.sh
export FOO="world"
export EXCITED="yes"
```

<br />

## LICENSE
[MIT](./license)
