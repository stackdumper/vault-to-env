# vault-to-env

`vault-to-env` reads secrets from vault and saves them to shell file as environment variables.

## Usage

```bash
Usage of vault-to-env:
  -e value
    	environment variable to fetch
  -o string
    	output file path
```

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

# read a secret and output it into stdout
# if you use KV v2, append /data/ before secret path
$ go run main.go -e FOO=/secret/data/hello#foo -e EXCITED=/secret/data/hello#excited
export FOO="world"
export EXCITED="yes"

# read a secret and output it to file
$ go run main.go -e FOO=/secret/data/hello#foo -e EXCITED=/secret/data/hello#excited -o /tmp/hello.sh
$ cat /tmp/hello.sh
export FOO="world"
export EXCITED="yes"
```
