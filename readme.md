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
cd ~/minter-hub/orchestrator
cargo install --path orchestrator
cargo install --path register_delegate_keys
```

## Run

1. Install and sync Minter Node
2. Install and sync Ethereum node
3. Install and sync Minter Hub node:

```bash
mkdir -p ~/.mhub/config/
curl https://raw.githubusercontent.com/MinterTeam/minter-hub/master/testnet-genesis.json > ~/.mhub/config/genesis.json
```

```bash
mhub start --p2p.persistent_peers="6f43047669b7a499b4c824624a748f5943739b6f@46.101.215.17:36656"
```

```bash
mhub keys add --keyring-backend test validator1
```

```bash
mhub tendermint show-validator
mhub tx staking create-validator --from=validator1 --keyring-backend test --amount=10hub --pubkey=cosmosvalconspub1zcjduepqfhcsrg04lmyyyzu9g2t72stduvt6j89lrsjxxww906g8zl69999qp05uhf  --commission-max-change-rate="0.1" --commission-max-rate="1" --commission-rate="0.1" --min-self-delegation="1"
```
