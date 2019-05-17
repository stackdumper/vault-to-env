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
  --lease-duration int   adjust secret lease duration
  --save-leases          save secret leases
  --vars strings         list of vars to read

Global Flags:
  --address string      Vault address
  --auth-token string   Vault auth token
  --auth-data strings   Vault auth data
  --auth-path string    Vault auth pat
```

#### Renew
Renew secrets leases

```
Usage:
  vte renew [flags]

Flags:
  --duration int         lease renew duration (default 3600)
  --leases stringArray   list leases to renew

Global Flags:
  --address string      Vault address
  --auth-data strings   Vault auth data
  --auth-path string    Vault auth path
  --auth-token string   Vault auth token
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
$ vte read \
    --auth-path /auth/userpass/login/tester \
    --auth-data password=tester \
    --vars FOO=/secret/data/hello#data.foo \
    --vars EXCITED=/secret/data/hello#data.excited

export FOO="world"
export EXCITED="yes"


# read MongoDB credentials, adjust and save leases
$ vte read \
    --auth-path /auth/userpass/login/tester \
    --auth-data password=tester \
    --vars MONGO_USER=/database/creds/admin#username \
    --vars MONGO_PASS=/database/creds/admin#password \
    --save-leases \
    --lease-duration 3600

export MONGO_USER="v-userpass-tester-admin-blCs4nhV8rJjoiWErvCn-1558078240"
export MONGO_USER_LEASE_ID="database/creds/admin/yh9NzzVbUReDyDm80kDdqgGw"
export MONGO_PASS="A1a-NbVnfS5zvxssQpVW"
export MONGO_PASS_LEASE_ID="database/creds/admin/wIK2rggEjlqVfZ3ybGCQ8Vvz"



# renew MongoDB credentials leases
$ go run main.go renew \
    --auth-path /auth/userpass/login/tester \
    --auth-data password=tester \
    --leases $MONGO_USER_LEASE_ID \
    --leases $MONGO_PASS_LEASE_ID \
    --duration 9000
```

<br />

## License
[MIT](./license)
