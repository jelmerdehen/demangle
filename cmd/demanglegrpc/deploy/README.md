# demanglegrpc deploy

Ships the gRPC service. Stage 6.5 in the v5.1 plan — not deployed
until a concrete non-Go non-skynet caller materialises.

## Files in this directory

| File | Purpose |
|---|---|
| `demanglegrpc.service` | Systemd unit for generic installs (`User=demanglegrpc`) |
| `skynet-demangle.service` | Systemd unit for the lux/skynet deployment (`User=demangle`) |
| `build.sh` | Build the binary and verify it is under 15 MB |
| `README.md` | This file |

## Building

```sh
bash cmd/demanglegrpc/deploy/build.sh
```

Produces `./demanglegrpc` in the repository root. Fails if the binary
exceeds 15 MB.

## Endpoints

| Service       | Port                | Proto  | Purpose                          |
|---------------|---------------------|--------|----------------------------------|
| gRPC          | 10.13.38.4:50061    | gRPC   | Demangle / Detect / Schemes / DemangleStream / context upload |
| health        | 10.13.38.4:50062    | HTTP   | `/healthz`, `/readyz`, `/metrics` |

## One-shot install on lux

```
# As root on 10.13.38.4.

# 1. Build + copy the binary.
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /usr/local/bin/skynet-demangle ./cmd/demanglegrpc

# 2. Create the service user + data dir.
useradd -r -s /usr/sbin/nologin demangle
install -d -o demangle -g demangle -m 0750 /var/lib/skynet-demangle

# 3. Install the systemd unit.
install -m 0644 cmd/demanglegrpc/deploy/skynet-demangle.service /etc/systemd/system/
systemctl daemon-reload

# 4. Start + enable.
systemctl enable --now skynet-demangle.service
systemctl status skynet-demangle.service

# 5. Smoke test.
curl -s http://10.13.38.4:50062/healthz
grpcurl -plaintext 10.13.38.4:50061 demangle.v1.Demangle/Schemes
```

## Monitoring

Prometheus scrape config (add to the existing skynet Prometheus):

```yaml
scrape_configs:
  - job_name: demangle
    static_configs:
      - targets: ['10.13.38.4:50062']
```

Metrics exposed:

| Metric                               | Type    | Description |
|--------------------------------------|---------|-------------|
| `demangle_uptime_seconds`            | gauge   | Seconds since process start |
| `demangle_requests_total`            | counter | Successful + failed demangle calls |
| `demangle_errors_total`              | counter | Failed demangle calls |
| `demangle_bytes_in_total`            | counter | Input bytes processed |
| `demangle_bytes_out_total`           | counter | Output bytes produced |
| `demangle_requests_by_scheme_total{scheme=…}` | counter | Per-scheme call counter |
| `demangle_goroutines`                | gauge   | `runtime.NumGoroutine()` |
| `demangle_heap_bytes`                | gauge   | Current heap allocation |

## Integration with skynet deploy-skynet skill

Add this entry to `deploy-skynet` under the lux tier:

```
when changed: /data/p/demangle/cmd/demanglegrpc/**
action:
  rsync /data/p/demangle/cmd/demanglegrpc/deploy/skynet-demangle.service \
        root@lux:/etc/systemd/system/
  ssh lux 'cd /data/p/demangle && GOOS=linux GOARCH=amd64 go build \
           -ldflags="-s -w" -o /usr/local/bin/skynet-demangle ./cmd/demanglegrpc'
  ssh lux 'systemctl daemon-reload && systemctl restart skynet-demangle'
  ssh lux 'curl -s http://10.13.38.4:50062/healthz'
```

## Rollback

```
systemctl stop skynet-demangle
cp /usr/local/bin/skynet-demangle.prev /usr/local/bin/skynet-demangle
systemctl start skynet-demangle
```

## Why the deploy was deferred

Per v5.1 decision 4: day-one consumers (skynet-scan, skynet-ghidra,
skynet-flows, behavox via skynet-GraphQL) all reach the library via
direct import. `cmd/demanglegrpc` exists in-repo + in-tests so the
moment a concrete non-Go non-skynet caller appears the deploy work
lands as a follow-up stage (6.5) without a green-field effort.
