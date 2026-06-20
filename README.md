# 3x-ui Metrics Exporter (3.3.1)

Prometheus exporter for [3x-ui](https://github.com/MHSanaei/3x-ui) **3.3.1** master panels. It polls panel APIs in the background, caches snapshots, and exposes metrics for Prometheus and Grafana.

## Features

- **API token authentication** (Bearer) — no cookie login
- **Multi-panel** support with `panel` label
- Metrics for **master + remote nodes**, **system**, **inbounds**, **clients**, and optional **outbounds**
- Per-client traffic from inbound `clientStats` (correct counters for `rate()`)
- Docker Compose stack: exporter + Prometheus + Grafana with pre-built dashboards

## Requirements

- 3x-ui **3.3.1** master panel
- API token: **Settings → Security → API Token** (full-admin credential)

## Configuration

Copy `config.example.yaml` to `config.yaml` and edit:

```yaml
listen_addr: ":2112"
poll_interval: 30s
request_timeout: 10s

panels:
  - name: "master-eu"
    base_url: "https://panel.example.com:54321/secret"
    api_token: "your-api-token"
    insecure_skip_verify: false
    collect_outbounds: true
```

| Field | Description |
|-------|-------------|
| `name` | Panel label (`panel` in metrics) |
| `base_url` | Full URL including optional base path |
| `api_token` | Bearer token from panel settings |
| `insecure_skip_verify` | Skip TLS verification (self-signed certs) |
| `collect_outbounds` | Scrape `/panel/api/xray/getOutboundsTraffic` |

### Environment overrides

- `CONFIG_PATH` — config file path (default `config.yaml`)
- `XUI_PANEL_<NAME>_TOKEN` — override `api_token` for panel `<name>` (e.g. `XUI_PANEL_MASTER_EU_TOKEN`)

## Run locally

```bash
go run ./cmd/exporter -config config.yaml
curl http://localhost:2112/metrics
curl http://localhost:2112/healthz
```

## Docker Compose (Prometheus + Grafana)

```bash
cp config.example.yaml config.yaml
# edit config.yaml with your panel URL and token

cd deploy
docker compose up -d
```

| Service | URL |
|---------|-----|
| Exporter | http://localhost:2112/metrics |
| Prometheus | http://localhost:9090 |
| Grafana | http://localhost:3000 (admin / admin) |

Grafana dashboards (folder **3x-ui**):

1. **Fleet Overview** — nodes table, online clients, traffic by node
2. **Node Detail** — master CPU/mem/disk/load/network, xray, uptime
3. **Users** — top clients, online, limits, expiry, groups
4. **Inbounds** — traffic and clients by protocol/inbound

## Metrics

All metrics use prefix `xui_`.

### Exporter

| Metric | Type | Labels |
|--------|------|--------|
| `xui_up` | gauge | `panel` |
| `xui_scrape_duration_seconds` | gauge | `panel` |
| `xui_scrape_errors_total` | counter | `panel`, `endpoint` |
| `xui_panel_info` | gauge | `panel`, `version` |

### Nodes (master + remote)

| Metric | Type | Labels |
|--------|------|--------|
| `xui_node_up` | gauge | `panel`, `node`, `role` |
| `xui_node_cpu_percent` | gauge | `panel`, `node`, `role` |
| `xui_node_mem_percent` | gauge | `panel`, `node`, `role` |
| `xui_node_online_clients` | gauge | `panel`, `node`, `role` |
| `xui_node_client_count` | gauge | `panel`, `node`, `role` |
| `xui_node_inbound_count` | gauge | `panel`, `node`, `role` |
| `xui_node_xray_up` | gauge | `panel`, `node`, `role` |

### System (master)

| Metric | Type | Labels |
|--------|------|--------|
| `xui_server_cpu_percent` | gauge | `panel` |
| `xui_server_mem_bytes` | gauge | `panel`, `state` |
| `xui_server_net_traffic_bytes_total` | counter | `panel`, `direction` |
| `xui_server_load` | gauge | `panel`, `period` |
| `xui_xray_up` | gauge | `panel` |

### Inbounds & clients

| Metric | Type | Labels |
|--------|------|--------|
| `xui_inbound_up_bytes_total` | counter | `panel`, `node`, `inbound_id`, ... |
| `xui_inbound_down_bytes_total` | counter | same |
| `xui_client_up_bytes_total` | counter | `panel`, `node`, `email`, `inbound_id`, ... |
| `xui_client_down_bytes_total` | counter | same |
| `xui_client_online` | gauge | `panel`, `email` |
| `xui_client_traffic_limit_bytes` | gauge | `panel`, `email` |
| `xui_client_group_info` | gauge | `panel`, `email`, `group` |

### Outbounds (optional)

| Metric | Type | Labels |
|--------|------|--------|
| `xui_outbound_up_bytes_total` | counter | `panel`, `tag` |
| `xui_outbound_down_bytes_total` | counter | `panel`, `tag` |

## Build

```bash
go build -o 3xui-metrics-collector ./cmd/exporter
```

## License

MIT
