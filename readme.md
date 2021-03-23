# Minter Hub

## Install

```bash
apt-get install git build-essential wget curl
wget https://golang.org/dl/go1.16.2.linux-amd64.tar.gz
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.16.2.linux-amd64.tar.gz
echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.profile
```

```bash
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
```

```bash
source ~/.profile
```

```bash
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
cargo build --all --release
```
