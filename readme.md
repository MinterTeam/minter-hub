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
cargo install --path orchestrator
cargo install --path register_delegate_keys
```

## Run
1. Install and sync Minter Node
```bash
minter node --testnet
```

2. Install and sync Ethereum node
```bash
geth --ropsten
```

3. Sync Minter Hub Node
```bash
# Download genesis
mkdir -p ~/.mhub/config/
curl https://raw.githubusercontent.com/MinterTeam/minter-hub/master/testnet-genesis.json > ~/.mhub/config/genesis.json

# Start and sync Minter Hub node
mhub start \
	--p2p.persistent_peers="bb75bf42dd14f55bb6528e7588d8e63cd2db2a44@46.101.215.17:36656"
```

4. Generate Hub account
```bash
mhub keys add --keyring-backend test validator1
```
	- ***WARNING: save generated key***
	- Request some test HUB to your generated address

5. Create Hub validator
```bash
mhub tendermint show-validator # show validator's public key
mhub tx staking create-validator \
	--from=validator1 \
	--keyring-backend=test \
	--amount=10hub \
	--pubkey=<VALIDATOR PUBLIC KEY>  \
	--commission-max-change-rate="0.1" \
	--commission-max-rate="1" \
	--commission-rate="0.1" \
	--min-self-delegation="1" \
	--chain-id=mhub-test
```
	- ***WARNING: save tendermint validator's key***

6. Generate Minter & Ethereum keys
```bash
mhub-keys-generator
```
	- ***WARNING: save generated keys***
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

	- **Start Hub ↔ Ethereum oracle.** Ethereum Contract for testnet: 0xfe9E069E52986ac50614A51590eBe183cc87Fc30
```bash
RUST_LOG=info orchestrator \
	--cosmos-phrase=<COSMOS MNEMONIC> \
	--ethereum-key=<ETHEREUM PRIVATE KEY> \
	--cosmos-grpc="http://127.0.0.1:1090" \
	--cosmos-legacy-rpc="http://127.0.0.1:1317" \
	--ethereum-rpc="http://127.0.0.1:8545/" \
	--fees=hub \
	--contract-address=<ADDRESS OF ETHEREUM CONTRACT> 
```

	- **Start Hub ↔ Minter oracle.** Minter Multisig for testnet: Mxffffffffffffffffffffffffffffffffffffffff, Start Minter Block for testnet: 2561976

```bash
mhub-minter-connector \
	--minter-multisig=<ADDRESS OF MINTER MULTISIG> \
	--minter-chain=testnet \
	--minter-mnemonic=<MINTER MNEMONIC> \
	--minter-node-url="127.0.0.1:8843/v2/" \
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

