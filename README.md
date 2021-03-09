# YAUTVending

YAUT Vending Machine

This is a sample for how to create a native token on Cardano and write a Golang program to sell it automatically.

## Disclaimer

USE AT YOUR OWN RISK.

## Requirements

* `cardano-node`
* `posrgres`
* `cardano-db-sync`
* `cardano-explorer-api`

You can set all of them up and running by using docker-compose, see https://github.com/input-output-hk/cardano-rest

**You have to wait till `cardano-node` fully synchroize with `mainnet` and `cardano-db-sync` fully synchronize
with `cardano-node`**

You need also generate your own signing key and verification for both Issuer and Distributor.

```shell
mkdir -p keys
cardano-cli address key-gen \
    --verification-key-file keys/issuer.vkey \
    --signing-key-file keys/issuer.skey
cardano-cli address key-gen \
    --verification-key-file keys/dist.vkey \
    --signing-key-file keys/dist.skey
```

## Configure

Configurations are hard-coded, you have to modify `main.go`

## Build

Golang can product static binary with no dynamic libraries linked, you can build `yautvending`, copy it
to `cardano-node` docker container, and run it inside the container.

Clone this repository, modify `main.go` and run

```shell
CGO_ENABLED=0 go build -mod vendor
```

You will have an executable named `yautvending`, copy that file into `cardano-node` container.

```shell
# create a workspace inside the container
docker run cardano-node mkdir /workspace

docker cp yautvending cardano-node:/workspace/
```


## Usage

Copy you key files into container.

```shell
docker cp keys cardano-node:/workspace/
```

Execute `yautvending` inside container

```shell
docker exec -ti cardano-node bash

cd /workspace
./yautvending -submit
```

## Explain

The `yautvending` program do these things.

1. Calculate distributor address (by invoking cardano-cli)
2. Calculate issuer verification key hash (by invoking cardano-cli)
3. Build policy script and generate policy id (by invoking cardano-cli)
4. Collect UTXOs for distributor address (by invoking cardano-cli)
5. Query transaction details (by invoking cardano-explorer-api)
6. Build a transaction that send back YAUT token along with minimum ADAs required and send all remaining ADAs to gringotts address (by invoking cardano-cli)
7. Calculate fee (by invoking cardano-cli)
8. Rebuild transaction with proper fee (by invoking cardano-cli)
9. Sign transaction (by invoking cardano-cli)
10. Submit transaction (by invoking cardano-cli)

## Credits

Guo Y.K., MIT License
