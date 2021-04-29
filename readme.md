# Minter Hub

## Build & Install

1. Install dependencies
```bash
apt-get update && \
	apt-get install -y git build-essential wget curl libssl-dev pkg-config
```

2. Install Golang
```bash
wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz && \
	rm -rf /usr/local/go && \
	tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin:~/go/bin' >> ~/.profile
```

3. Install Rust
```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.profile
```

4. Clone Minter Hub repository
```bash
cd ~ && git clone https://github.com/MinterTeam/minter-hub.git
```

5. Compile & install 
```bash
# Minter Hub node
cd ~/minter-hub/chain
make install

# Hub ↔ Minter oracle
cd ~/minter-hub/minter-connector
make install

# Prices oracle
cd ~/minter-hub/oracle
make install

# Keys generator
cd ~/minter-hub/keys-generator
make install

# Hub ↔ Ethereum oracle
cd ~/minter-hub/orchestrator
cargo install --locked --path orchestrator
cargo install --locked --path register_delegate_keys
```

## Run
1. Install and sync Minter Node 
```bash
minter node
```

2. Install and sync Ethereum node
```bash
geth --rpc --rpcaddr "127.0.0.1" --rpcport "8545"
```

3. Sync Minter Hub Node
```bash
# Download genesis
mkdir -p ~/.mhub/config/
curl https://raw.githubusercontent.com/MinterTeam/minter-hub/master/mainnet-genesis.json > ~/.mhub/config/genesis.json

# Start and sync Minter Hub node
mhub start \
	--p2p.persistent_peers="b740ff04fadabce115b4bcb296cab9812694e4d5@104.236.213.173:26656"
```

- **IMPORTANT**: After syncing you must edit `~/.mhub/config/app.toml`: enable API in respective section.

4. Generate Hub account
```bash
mhub keys add validator1
```

- **WARNING: save generated key**
- Request some test HUB to your generated address

5. Create Hub validator
```bash
mhub tendermint show-validator # show validator's public key
mhub tx staking create-validator \
	--from=validator1 \
	--amount=1000000000000000000hub \
	--pubkey=<VALIDATOR PUBLIC KEY>  \
	--commission-max-change-rate="0.1" \
	--commission-max-rate="1" \
	--commission-rate="0.1" \
	--min-self-delegation="1" \
	--chain-id=mhub-mainnet-1
```

- **WARNING: save tendermint validator's key**
- An important point: the validator is turned off if it does not commit data for a long time. You can turn in on again by sending an unjail transaction. Docs: `mhub tx slashing unjail --help`

6. Generate Minter & Ethereum keys
```bash
mhub-keys-generator
```
- **WARNING: save generated keys**
- Request some test ETH to your generated address

7. Register Ethereum keys
```bash
register-peggy-delegate-keys \
	--cosmos-phrase=<COSMOS MNEMONIC> \
	--validator-phrase=<COSMOS MNEMONIC> \
	--ethereum-key=<ETHEREUM PRIVATE KEY> \
	--cosmos-rpc="http://127.0.0.1:1317" \
	--fees=hub
```

8. Start services. *You can set them up as services or run in different terminal screens.*

- **Start Hub ↔ Ethereum oracle.** 
```
Ethereum Contract for testnet: 0x8D0B99eAE97a247B2C3E9B3dA6774B8359a37537

Ethereum Contract for mainnet: 0xc735478ef7562ecc37662fc7c5e521eb835f9dab
```
```bash
RUST_LOG=info orchestrator \
	--cosmos-phrase=<COSMOS MNEMONIC> \
	--ethereum-key=<ETHEREUM PRIVATE KEY> \
	--cosmos-grpc="http://127.0.0.1:9090" \
	--cosmos-legacy-rpc="http://127.0.0.1:1317" \
	--ethereum-rpc="http://127.0.0.1:8545/" \
	--fees=hub \
	--contract-address=<ADDRESS OF ETHEREUM CONTRACT> 
```

- **Start Hub ↔ Minter oracle.** 
```
Minter Multisig for testnet: Mx703880f64588b3247f8a583f1bef5a6ac5aeac59
Start Minter Block for testnet: 3393443

Minter Multisig for mainnet: Mx68f4839d7f32831b9234f9575f3b95e1afe21a56
Start Minter Block for mainnet: 3442652
```
```bash
mhub-minter-connector \
	--minter-multisig=<ADDRESS OF MINTER MULTISIG> \
	--minter-chain=<testnet|mainnet> \
	--minter-mnemonic=<MINTER MNEMONIC> \
	--minter-node-url="tcp://127.0.0.1:8843/v2/" \
	--cosmos-mnemonic=<COSMOS MNEMONIC> \
	--cosmos-node-url="127.0.0.1:9090" \
	--tm-node-url="127.0.0.1:26657" \
	--minter-start-block=<MINTER START BLOCK> \
	--minter-start-event-nonce=1 \
	--minter-start-batch-nonce=1 \
	--minter-start-valset-nonce=1
```
	
- **Start price oracle**
```bash
mhub-oracle \
	--minter-node-url="127.0.0.1:8843/v2/" \
	--cosmos-mnemonic=<COSMOS MNEMONIC> \
	--cosmos-node-url="127.0.0.1:9090" \
	--tm-node-url="127.0.0.1:26657" 
```

