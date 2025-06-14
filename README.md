# 3xui-metrics-collector

Service for collecting metrics from 3xui control panel and exporting them in Prometheus format.

## Description

This service collects metrics from the 3xui control panel and exports them in Prometheus format. Metrics include:
- Number of clients
- Total traffic
- Service status
- Other metrics available through the 3xui API

## Requirements

- Go 1.21 or higher (for building from source)
- Access to 3xui control panel
- Prometheus (for collecting metrics)

## Installation Methods

### Method 1: Using Docker (Recommended)

Run the service using Docker:

```bash
docker run -d \
  --name 3xui-metrics-collector \
  -p 2112:2112 \
  -e X_UI_HOST=your-xui-host \
  -e X_UI_PORT=54321 \
  -e X_UI_USERNAME=admin \
  -e X_UI_PASSWORD=your-password \
  ra1zee/3xui-metrics-collector:latest
```

With docker-compose:

```yaml
version: '3.8'
services:
  3xui-metrics-collector:
    image: ra1zee/3xui-metrics-collector:latest
    ports:
      - "2112:2112"
    environment:
      - X_UI_HOST=your-xui-host
      - X_UI_PORT=54321
      - X_UI_USERNAME=admin
      - X_UI_PASSWORD=your-password
    restart: unless-stopped
```

### Method 2: Download Pre-built Binary

Download the latest release for your platform:

```bash
# For Linux amd64
curl -L -o 3xui-metrics-collector.tar.gz https://github.com/Ra1ze505/3xui-metrics-collector/releases/latest/download/3xui-metrics-collector-linux-amd64.tar.gz
tar -xzf 3xui-metrics-collector.tar.gz

# For Linux arm64
curl -L -o 3xui-metrics-collector.tar.gz https://github.com/Ra1ze505/3xui-metrics-collector/releases/latest/download/3xui-metrics-collector-linux-arm64.tar.gz
tar -xzf 3xui-metrics-collector.tar.gz

# For macOS amd64
curl -L -o 3xui-metrics-collector.tar.gz https://github.com/Ra1ze505/3xui-metrics-collector/releases/latest/download/3xui-metrics-collector-darwin-amd64.tar.gz
tar -xzf 3xui-metrics-collector.tar.gz

# For macOS arm64 (Apple Silicon)
curl -L -o 3xui-metrics-collector.tar.gz https://github.com/Ra1ze505/3xui-metrics-collector/releases/latest/download/3xui-metrics-collector-darwin-arm64.tar.gz
tar -xzf 3xui-metrics-collector.tar.gz

# For Windows amd64
curl -L -o 3xui-metrics-collector.zip https://github.com/Ra1ze505/3xui-metrics-collector/releases/latest/download/3xui-metrics-collector-windows-amd64.zip
unzip 3xui-metrics-collector.zip
```

### Method 3: Build from Source

1. Clone the repository:
```bash
git clone https://github.com/Ra1ze505/3xui-metrics-collector.git
cd 3xui-metrics-collector
```

2. Build the project:
```bash
go build -o 3xui-metrics-collector
```

## Configuration

The service uses the following environment variables:

- `X_UI_HOST` - 3xui control panel host
- `X_UI_PORT` - 3xui control panel port  
- `X_UI_BASEPATH` - API base path (empty by default)
- `X_UI_USERNAME` - username for API access
- `X_UI_PASSWORD` - password for API access

## Setting up systemd service

After downloading the binary (Method 2), you can set it up as a systemd service:

1. Create the necessary directories and move files:
```bash
# Download and extract the binary (example for Linux amd64)
curl -L -o 3xui-metrics-collector.tar.gz https://github.com/Ra1ze505/3xui-metrics-collector/releases/latest/download/3xui-metrics-collector-linux-amd64.tar.gz
tar -xzf 3xui-metrics-collector.tar.gz

# Create directories
sudo mkdir -p /opt/3xui-metrics-collector
sudo mkdir -p /etc/3xui-metrics-collector

# Move binary
sudo mv 3xui-metrics-collector /opt/3xui-metrics-collector/
```

2. Create configuration file `/etc/3xui-metrics-collector/config.env`:
```bash
sudo tee /etc/3xui-metrics-collector/config.env > /dev/null <<EOF
X_UI_HOST=your_xui_host
X_UI_PORT=your_xui_port
X_UI_BASEPATH=your_base_path
X_UI_USERNAME=your_username
X_UI_PASSWORD=your_password
```

3. Create systemd service file `/etc/systemd/system/3xui-metrics-collector.service`:
```bash
sudo tee /etc/systemd/system/3xui-metrics-collector.service > /dev/null <<EOF
[Unit]
Description=3xui Metrics Collector
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/3xui-metrics-collector
ExecStart=/opt/3xui-metrics-collector/3xui-metrics-collector
EnvironmentFile=/etc/3xui-metrics-collector/config.env
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

4. Set proper permissions:
```bash
sudo chown -R root:root /opt/3xui-metrics-collector
sudo chown -R root:root /etc/3xui-metrics-collector
sudo chmod 755 /opt/3xui-metrics-collector/3xui-metrics-collector
sudo chmod 600 /etc/3xui-metrics-collector/config.env
```

5. Start the service:
```bash
sudo systemctl daemon-reload
sudo systemctl enable 3xui-metrics-collector
sudo systemctl start 3xui-metrics-collector
```

## Service Management

- Check status:
```bash
sudo systemctl status 3xui-metrics-collector
```

- Stop service:
```bash
sudo systemctl stop 3xui-metrics-collector
```

- Restart service:
```bash
sudo systemctl restart 3xui-metrics-collector
```

- View logs:
```bash
sudo journalctl -u 3xui-metrics-collector -f
```

## Metrics

The service exports metrics in Prometheus format on port 2112. Available metrics:

- `xui_clients_total` - total number of clients
- `xui_traffic_total_bytes` - total traffic in bytes
- `xui_service_status` - service status (1 - active, 0 - inactive)

Access metrics at: `http://localhost:2112/metrics`

## Prometheus Configuration

1. Install Prometheus if not already installed:
```bash
# For Ubuntu/Debian
sudo apt-get update
sudo apt-get install prometheus

# For CentOS/RHEL
sudo yum install prometheus
```

2. Add metrics collection configuration to `/etc/prometheus/prometheus.yml`:
```yaml
scrape_configs:
  - job_name: '3xui'
    static_configs:
      - targets: ['localhost:2112']
    scrape_interval: 15s
    scrape_timeout: 10s
```

3. Restart Prometheus to apply changes:
```bash
sudo systemctl restart prometheus
```

4. Verify metrics are being collected:
   - Open Prometheus web interface (usually available at http://localhost:9090)
   - Go to Status -> Targets
   - Make sure the `3xui` target is in UP state

5. Example queries for Grafana:
```promql
# Total number of clients
xui_clients_total

# Total traffic in gigabytes
xui_traffic_total_bytes / 1024 / 1024 / 1024

# Service status
xui_service_status
```

## Docker Image Tags

Available Docker image tags:
- `latest` - latest stable release
- `v1.0.0` - specific version tags
- `main` - latest development build

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## License

MIT 