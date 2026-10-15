# upsApi

REST API and Prometheus exporter for UPS monitoring via [NUT](https://networkupstools.org/) (`upsc`).

Built in Go with [Gin](https://github.com/gin-gonic/gin). Designed to run on a Raspberry Pi (or any Linux host) connected to a UPS over USB.

## Features

- JSON endpoint to read full UPS status (`/api/read`)
- Prometheus metrics endpoint (`/metrics`) with automatic refresh every 10s
- Health check endpoints (`/health/liveness`, `/health/readiness`)
- Multi-arch Docker image (amd64, arm64) published to GHCR

## Prometheus Metrics

All metrics are prefixed with `upsapi_`:

| Metric | Description |
|---|---|
| `upsapi_battery_charge` | Battery charge (%) |
| `upsapi_battery_charge_low` | Low charge threshold (%) |
| `upsapi_battery_runtime` | Battery runtime (seconds) |
| `upsapi_input_transfer_high` | Input transfer high voltage |
| `upsapi_input_transfer_low` | Input transfer low voltage |
| `upsapi_output_voltage` | Output voltage |
| `upsapi_output_frequency_nominal` | Output nominal frequency |
| `upsapi_ups_load` | UPS load (%) |
| `upsapi_ups_real_power` | Real power (W) |
| `upsapi_ups_power_nominal` | Nominal power (VA/W) |
| `upsapi_ups_info` | UPS info (labels: status, manufacturer, model, serial, firmware, type) |
| `upsapi_ups_beeper_enabled` | Beeper status (1=enabled, 0=disabled) |

## Configuration

All settings are configured via environment variables (or a `.env` file):

| Variable | Default | Description |
|---|---|---|
| `ENV` | `dev` | Environment (`dev` enables Gin debug mode) |
| `LOG_LEVEL` | `0` | Zerolog level (0=debug, 1=info, 2=warn, 3=error, 4=fatal) |
| `LISTEN_PORT` | `9695` | HTTP listen port |
| `UPS_USER` | `eaton` | NUT UPS name (the `<ups>` part of `upsc <ups>@<host>`) |
| `UPS_HOST` | `localhost` | NUT server host |

## Prerequisites

The host (or container) must have **NUT client tools** installed (`upsc` command) and a NUT server (`upsd`) running and reachable.

If the UPS is connected via USB directly to the host, the NUT driver (e.g. `usbhid-ups`) must be configured on the host or inside the container.

## Run locally

```bash
# Install NUT client
sudo apt install nut-client

# Check your UPS is visible
upsc eaton@localhost

# Run
go run .
```

## Docker

### Build

```bash
docker build -t upsapi .
```

### Run with USB passthrough

The UPS communicates over USB. To give the container access to the USB device, you need to pass it through from the host.

**1. Find your UPS USB device on the host:**

```bash
lsusb
# Example output:
# Bus 001 Device 004: ID 0463:ffff EATON 5E
```

The device file is typically `/dev/bus/usb/<bus>/<device>`, e.g. `/dev/bus/usb/001/004`.

**2. Run the container with `--device`:**

```bash
docker run -d \
  --name upsapi \
  --device /dev/bus/usb/001/004 \
  -p 9695:9695 \
  -e UPS_USER=eaton \
  -e UPS_HOST=localhost \
  -e ENV=production \
  -e LOG_LEVEL=1 \
  ghcr.io/mandelv/upsapi:latest
```

> **Note:** The device number can change after a reboot. For a stable mapping, use a udev rule (see below).

**3. Alternative: udev rule for a stable device path:**

Create `/etc/udev/rules.d/99-ups.rules`:

```
SUBSYSTEM=="usb", ATTR{idVendor}=="0463", ATTR{idProduct}=="ffff", SYMLINK+="ups0", MODE="0666"
```

Then reload:

```bash
sudo udevadm control --reload-rules && sudo udevadm trigger
```

And run the container with the symlink:

```bash
docker run -d \
  --name upsapi \
  --device /dev/ups0 \
  -p 9695:9695 \
  -e UPS_USER=eaton \
  -e UPS_HOST=localhost \
  ghcr.io/mandelv/upsapi:latest
```

**4. Alternative: privileged mode (not recommended):**

```bash
docker run -d \
  --name upsapi \
  --privileged \
  -v /dev/bus/usb:/dev/bus/usb \
  -p 9695:9695 \
  -e UPS_USER=eaton \
  -e UPS_HOST=localhost \
  ghcr.io/mandelv/upsapi:latest
```

### NUT inside vs. outside the container

The Docker image includes `nut` (via `apk`). Two setups are possible:

#### NUT inside the container

Mount the USB device (`--device`) and configure `/etc/nut/` inside the container (ups.conf, upsd.conf, upsd.users, upsdrv). Use `UPS_HOST=localhost`.

#### NUT on the host (recommended)

The container only needs the `upsc` client, no USB passthrough required. Two ways to reach the host's `upsd`:

**Option A: `network_mode: host` (simplest)**

The container shares the host network stack, so `UPS_HOST=localhost` works directly:

```yaml
services:
  upsapi:
    image: ghcr.io/mandelv/upsapi:latest
    network_mode: host
    devices:
      - "/dev/ups0"
    environment:
      UPS_USER: eaton
      UPS_HOST: localhost
      ENV: production
      LOG_LEVEL: "1"
    restart: unless-stopped
```

**Option B: Docker bridge network**

Point `UPS_HOST` to the host's IP on the Docker bridge (default `172.17.0.1`):

```yaml
environment:
  UPS_HOST: 172.17.0.1
```

You can find the actual IP with `ip addr show docker0`.

For this to work, `upsd` on the host must listen on that interface. Edit `/etc/nut/upsd.conf`:

```
LISTEN 127.0.0.1 3493
LISTEN 172.17.0.1 3493
```

Then restart the NUT server:

```bash
sudo systemctl restart nut-server
```

> **Note:** `host.docker.internal` does not work natively on Linux (it is a Docker Desktop feature for macOS/Windows).

### Docker Compose

```yaml
services:
  upsapi:
    image: ghcr.io/mandelv/upsapi:latest
    ports:
      - "9695:9695"
    devices:
      - "/dev/ups0"
    environment:
      UPS_USER: eaton
      UPS_HOST: localhost
      ENV: production
      LOG_LEVEL: "1"
    restart: unless-stopped
```

## API

| Endpoint | Description |
|---|---|
| `GET /api/read` | Full UPS status as JSON (parsed from `upsc`) |
| `GET /metrics` | Prometheus metrics |
| `GET /health/liveness` | Liveness probe (returns `200 OK`) |
| `GET /health/readiness` | Readiness probe (returns `200 OK`) |

## CI/CD

A GitHub Actions workflow (`.github/workflows/build.yml`) builds and pushes a multi-arch Docker image to GHCR on every tag push (`v*`).

## License

MIT
