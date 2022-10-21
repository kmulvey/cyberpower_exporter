# cpwatch
[![cpwatch](https://github.com/kmulvey/cpwatch/actions/workflows/release_build.yml/badge.svg)](https://github.com/kmulvey/cpwatch/actions/workflows/release_build.yml) [![codecov](https://codecov.io/gh/kmulvey/cpwatch/branch/main/graph/badge.svg?token=wp6NcwDC5k)](https://codecov.io/gh/kmulvey/cpwatch) [![Go Report Card](https://goreportcard.com/badge/github.com/kmulvey/cpwatch)](https://goreportcard.com/report/github.com/kmulvey/cpwatch) [![Go Reference](https://pkg.go.dev/badge/github.com/kmulvey/cpwatch.svg)](https://pkg.go.dev/github.com/kmulvey/cpwatch)

Monitor and store CyberPower UPS statistics.

## Usage
Becasue pwrstat needs to be run as root, this tool needs to be run as root as well.
```
suod cpwatch -h
```

## HTTP Endpoints
| Route        | Description   
|--------------|----------------------------------|
| /latest      | return the latest statistics     |
