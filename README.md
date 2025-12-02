`````markdown
````markdown
# go-eth

## DemoToken (ERC-20)

- Source: `contracts/DemoToken.sol`
- Solidity ^0.8.20, owner-mintable supply (constructor mints to deployer, `mint` allows more tokens for testing).

### Quick Deploy (Remix or Hardhat)

1. Compile the contract with the injected compiler (0.8.20+).
2. Deploy `DemoToken` providing the initial supply in wei (e.g., `1000000000000000000000` for 1000 DEMO with 18 decimals).
3. Optionally call `mint(<address>, <amount>)` from the deployer/owner to top up any account for CLI demos.
4. The deployer receives the constructor supply; use `transfer`, `mint`, `approve`, or `transferFrom` for testing flows.

### Local Deployment Helper

```bash
# Requires Ganache (or any RPC) running on http://127.0.0.1:8545
RPC_URL=http://127.0.0.1:8545 \
PRIVATE_KEY=<ganache-privkey> \
INITIAL_SUPPLY=0 \
MINT_TO=<recipient-address> \
MINT_AMOUNT=1000000000000000000000 \
node scripts/deploy_demo_token.js
```

The script compiles, deploys, optionally mints, and prints the token address plus mint tx hash for use with Go CLIs.

The design keeps the implementation lean so it is easy to read alongside the Go examples in `cmd/`.

### Token Metadata & Balance CLI

Generate / refresh the Go binding (stored in `token/erc20.go`) and query any ERC-20 with the new helper:

```bash
npx solcjs --abi contracts/DemoToken.sol -o build --base-path . --include-path node_modules
cp build/contracts_DemoToken_sol_DemoToken.abi erc20_sol_ERC20.abi
abigen --abi=erc20_sol_ERC20.abi --pkg=token --out=token/erc20.go

go run ./cmd/token-balance \
  --rpc=https://mainnet.infura.io/v3/<PROJECT_ID> \
  --contract=0xa74476443119A942dE498590Fe1f2454d7D4aC0d \
  --account=0x0536806df512d6cdde913cf95c9886f65b1d3462
```

The CLI connects to the supplied RPC, instantiates the local binding, prints name/symbol/decimals, and reports both the raw wei balance and a human-friendly value that accounts for token decimals.

### InfluxDB helper

Some monitoring scripts expect an InfluxDB HTTP API at `http://localhost:8086`. Start a disposable instance with Docker before issuing queries:

```
docker run --rm -p 8086:8086 influxdb:1.8
```

Or rely on the bundled compose setup, which persists data in the `influxdb-data` volume:

```
docker compose up -d influxdb
docker compose logs -f influxdb
```

Wait for the container/compose service to finish booting, then run commands such as `curl -XPOST "http://localhost:8086/query" --data-urlencode "q=CREATE USER ..."`.

### Metrics pipeline quickstart

1. Start the local services:
   ```bash
   docker compose up -d influxdb
   docker run -d --name goeth-prom -p 9090:9090 \
     -v "$PWD/monitoring/prometheus:/etc/prometheus" \
     prom/prometheus:latest
   ```
2. Provision the Influx database and user (skip if already created):
   ```bash
   curl -u admin:adminpass -XPOST "http://localhost:8086/query" \
     --data-urlencode "q=CREATE DATABASE geth"
   curl -u admin:adminpass -XPOST "http://localhost:8086/query" \
     --data-urlencode "q=CREATE USER username WITH PASSWORD 'password'"
   curl -u admin:adminpass -XPOST "http://localhost:8086/query" \
     --data-urlencode "q=GRANT ALL ON geth TO username"
   ```
3. Export the credentials and start geth with metrics already wired up:
   ```bash
   export INFLUX_USER=username
   export INFLUX_PASSWORD=password
   scripts/run_geth_metrics.sh --syncmode snap --mainnet
   ```
   The helper enables `/debug/metrics/prometheus` on port `6060` for Prometheus while publishing the same metrics to InfluxDB using the provided credentials.
4. Open `http://localhost:9090`, confirm `Status → Targets` lists the `go-ethereum` job as *UP*, then graph queries like `geth_p2p_peers`, `sum(rate(geth_txpool_pending{}[1m]))`, or `histogram_quantile(0.95, sum(rate(geth_rpc_durations_histogram_bucket[5m])) by (le, method))`.
5. Use Influx to spot-check raw data:
   ```bash
   curl -u username:password -G http://localhost:8086/query \
     --data-urlencode "db=geth" \
     --data-urlencode "q=SHOW MEASUREMENTS"
   ```

Stop everything with `docker compose down` and `docker rm -f goeth-prom` when finished.

````

`````
