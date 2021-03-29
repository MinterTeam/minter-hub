# Minter Hub

## Install

```bash
apt-get update
apt-get install -y git build-essential wget curl libssl-dev pkg-config
wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.profile
```

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

```bash
source ~/.profile
```

```bash
cd ~
git clone https://github.com/MinterTeam/minter-hub.git
```

```bash
cd ~/minter-hub/chain
make install
```

```bash
cd ~/minter-hub/minter-connector
make install
```

```bash
cd ~/minter-hub/oracle
make install
```

```bash
cd ~/minter-hub/keys-generator
make install
```

```bash
cd ~/minter-hub/orchestrator
cargo install --path orchestrator
cargo install --path register_delegate_keys
```

## Run

1. Install and sync Minter Node
2. Install and sync Ethereum node
3. Install and sync Minter Hub node
4. Start services

## Download Minter Hub genesis
```bash
mkdir -p ~/.mhub/config/
curl https://raw.githubusercontent.com/MinterTeam/minter-hub/master/testnet-genesis.json > ~/.mhub/config/genesis.json
```

## Start and sync Minter Hub node
```bash
mhub start --p2p.persistent_peers="6f43047669b7a499b4c824624a748f5943739b6f@46.101.215.17:36656"
```

## Generate Hub account
```bash
mhub keys add --keyring-backend test validator1
```

## Create Hub validator
```bash
mhub tendermint show-validator
mhub tx staking create-validator --from=validator1 --keyring-backend test --amount=10hub --pubkey=cosmosvalconspub.......  --commission-max-change-rate="0.1" --commission-max-rate="1" --commission-rate="0.1" --min-self-delegation="0"
```

## Generate Minter & Ethereum keys
```bash
mhub-keys-generator
```

## Register Ethereum keys
```bash
register-peggy-delegate-keys --cosmos-phrase="COSMOS MNEMONIC" --validator-phrase="COSMOS MNEMONIC" --ethereum-key=ETHEREUM PRIVATE KEY --cosmos-rpc="http://127.0.0.1:1317" --fees=hub
```

## Start services

### Start Hub <-> Ethereum oracle
```bash
orchestrator --cosmos-phrase="COSMOS MNEMONIC" --ethereum-key=ETHEREUM PRIVATE KEY --cosmos-grpc="http://127.0.0.1:1090" --cosmos-legacy-rpc="http://127.0.0.1:1317" --ethereum-rpc="http://127.0.0.1:8545/" --fees=hub --contract-address=ADDRESS OF ETHEREUM CONTRACT 
```

### Start Hub <-> Minter oracle
```bash
mhub-minter-connector --minter-multisig=Mxffffffffffffffffffffffffffffffffffffffff --minter-chain=testnet --minter-mnemonic="MINTER MNEMONIC" --minter-node-url="https://node-api.taconet.minter.network/v2/" --cosmos-mnemonic="COSMOS MNEMONIC" --cosmos-node-url="127.0.0.1:9090" --tm-node-url="127.0.0.1:26657" --minter-start-block=1 --minter-start-event-nonce=1 --minter-start-batch-nonce=1 --minter-start-valset-nonce=1
```

### Start price oracle
```bash
mhub-oracle --minter-node-url="https://node-api.taconet.minter.network/v2/" --cosmos-mnemonic="COSMOS MNEMONIC" --cosmos-node-url="127.0.0.1:9090" --tm-node-url="127.0.0.1:26657" 
```

