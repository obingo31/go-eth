#!/usr/bin/env bash
set -euo pipefail

# Helper to run geth with Prometheus scraping enabled and metrics pushed to InfluxDB.
# Customize behaviour through the environment variables documented below or by
# passing extra geth flags as CLI arguments when invoking this script.

: "${GETH_BIN:=geth}"
: "${GETH_DATADIR:=$PWD/.ethereum}"
: "${GETH_DEFAULT_ARGS:=--dev --http --http.addr 0.0.0.0 --http.vhosts * --http.api eth,net,web3,debug --ws --ws.addr 0.0.0.0 --ws.origins * --ws.api eth,net,web3,debug}"
: "${METRICS_ADDR:=0.0.0.0}"
: "${METRICS_PORT:=6060}"
: "${INFLUX_ENDPOINT:=http://localhost:8086}"
: "${INFLUX_DB:=geth}"
: "${INFLUX_USER:=}"
: "${INFLUX_PASSWORD:=}"

if [[ -z "$INFLUX_USER" || -z "$INFLUX_PASSWORD" ]]; then
  echo "INFLUX_USER and INFLUX_PASSWORD must be set before running this script" >&2
  exit 1
fi

read -r -a DEFAULT_ARGS <<< "$GETH_DEFAULT_ARGS"

exec "$GETH_BIN" \
  --datadir "$GETH_DATADIR" \
  "${DEFAULT_ARGS[@]}" \
  --metrics \
  --metrics.addr "$METRICS_ADDR" \
  --metrics.port "$METRICS_PORT" \
  --metrics.expensive \
  --metrics.influxdb \
  --metrics.influxdb.endpoint "$INFLUX_ENDPOINT" \
  --metrics.influxdb.database "$INFLUX_DB" \
  --metrics.influxdb.username "$INFLUX_USER" \
  --metrics.influxdb.password "$INFLUX_PASSWORD" \
  "$@"
