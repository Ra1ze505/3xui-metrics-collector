# 3x-ui Metrics Collector

A Prometheus metrics collector for 3x-ui panel that collects client traffic statistics.

## Features

- Collects upload and download statistics for all clients
- Exposes metrics in Prometheus format
- Configurable collection interval
- Secure HTTPS communication with 3x-ui panel

## Configuration

The collector is configured using environment variables:

- `X_UI_HOST` - Hostname or IP address of the 3x-ui panel
- `X_UI_PORT` - Port number of the 3x-ui panel
- `X_UI_BASEPATH` - Base path of the 3x-ui panel (e.g., "/")
- `X_UI_USERNAME` - Username for 3x-ui panel authentication
- `X_UI_PASSWORD` - Password for 3x-ui panel authentication

## Metrics

The collector exposes the following metrics:

- `xui_client_up_bytes` - Total uploaded bytes for each client
- `xui_client_down_bytes` - Total downloaded bytes for each client

Each metric includes the following labels:
- `email` - Client email
- `client_id` - Client ID
- `inbound_id` - Inbound ID
- `inbound_remark` - Inbound remark

## Usage

1. Set the required environment variables
2. Run the collector:
   ```bash
   ./3xui-metrics-collector
   ```
3. Access metrics at `http://localhost:2112/metrics`

## Prometheus Configuration

Add the following to your Prometheus configuration:

```yaml
scrape_configs:
  - job_name: '3xui'
    static_configs:
      - targets: ['localhost:2112']
``` 