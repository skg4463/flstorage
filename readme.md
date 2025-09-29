# flstorage
**flstorage** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).


## Environment
Ignite CLI version:             v29.3.1-dev

Ignite CLI source hash:         845a1a8886b8a098ed56372bab45ddee5caea526

Ignite CLI config version:      v1

Cosmos SDK version:             v0.53.3

Your go version:                go version go1.24.7 linux/amd64

## Get started

```
ignite chain serve

계정 및 키 확인
flstoraged keys list

업로드 확인된 alice 키 삽입 [fls-client - uploader]
go run main.go ../model_round_1_client_1.bin "1-[aliceAddr]-flstorage" [aliceAddr]

결과의 originalHash를 통해 query 
flstoraged query storage show-stored-file [originalHash]

txhash를 통해 블록 생성 확인 
flstoraged query tx [txHash]

block height를 통한 블록 확인
flstoraged query block --type=height 30248

downloader를 통해 originalHash를 인자로 복원 [fls-client - downloader]
go run main.go [originalHash] ../restored_round_weight.bin

checksum 비교를 통해 복원 확인
sha256sum model_round_1_client_1.bin restored_round_weight.bin
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

Port info: 
    tendermint node(RPC): 26657,
    blockchain API: 1317,
    token faucet 4500,
    gPRC: 9090

### Configure

Your blockchain in development can be configured with `config.yml`. To learn more, see the [Ignite CLI docs](https://docs.ignite.com).

### Web Frontend

Additionally, Ignite CLI offers a frontend scaffolding feature (based on Vue) to help you quickly build a web frontend for your blockchain:

Use: `ignite scaffold vue`
This command can be run within your scaffolded blockchain project.


For more information see the [monorepo for Ignite front-end development](https://github.com/ignite/web).

## Release
To release a new version of your blockchain, create and push a new tag with `v` prefix. A new draft release with the configured targets will be created.

```
git tag v0.1
git push origin v0.1
```

After a draft release is created, make your final changes from the release page and publish it.

### Install
To install the latest version of your blockchain node's binary, execute the following command on your machine:

```
curl https://get.ignite.com/username/flstorage@latest! | sudo bash
```
`username/flstorage` should match the `username` and `repo_name` of the Github repository to which the source code was pushed. Learn more about [the install process](https://github.com/ignite/installer).

## Learn more

- [Ignite CLI](https://ignite.com/cli)
- [Tutorials](https://docs.ignite.com/guide)
- [Ignite CLI docs](https://docs.ignite.com)
- [Cosmos SDK docs](https://docs.cosmos.network)
- [Developer Chat](https://discord.com/invite/ignitecli)
