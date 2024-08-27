# cyberpower_exporter
[![Build](https://github.com/kmulvey/cyberpower_exporter/actions/workflows/build.yml/badge.svg)](https://github.com/kmulvey/cyberpower_exporter/actions/workflows/build.yml) [![Release](https://github.com/kmulvey/cyberpower_exporter/actions/workflows/release.yml/badge.svg)](https://github.com/kmulvey/cyberpower_exporter/actions/workflows/release.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/kmulvey/cpwatch)](https://goreportcard.com/report/github.com/kmulvey/cpwatch) [![Go Reference](https://pkg.go.dev/badge/github.com/kmulvey/cpwatch.svg)](https://pkg.go.dev/github.com/kmulvey/cpwatch)

Monitor and store CyberPower UPS statistics.

## Installation and Usage
Several linux package formats are available in the releases. Becasue pwrstat needs to be run as root, this tool needs to be run as root as well.

### Manual linux install:
- `sudo cp cyberpower_exporter /usr/bin/` (this path can be changed if you like, just be sure to change the path in the service file as well)
- `sudo cp cyberpower_exporter.service /etc/systemd/system/`
- `sudo systemctl daemon-reload`
- `sudo systemctl enable cyberpower_exporter`
- `sudo systemctl restart cyberpower_exporter`
- Import grafana-config.json to your grafana instance
- enjoy!

![Screenshot](https://github.com/kmulvey/cyberpower_exporter/blob/main/screenshot.jpg?raw=true)
