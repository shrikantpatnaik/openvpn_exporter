# OpenVPN Exporter for Prometheus

A simple exporter that reads the OpenVPN status file and exposes the data as Prometheus Metrics

## Usage

### Local

```
go build openvpn_exporter.go
./openvpn_exporter --openvpn.status_path=/path/to/openvpn.status
```

### Docker

To use with docker you must mount your status file to /etc/openvpn_exporter/server.status

Example:
```
docker run -p 9176:9176 -v /path/to/your/server.status:/etc/openvpn_exporter/server.status shrikantpatnaik/openvpn_exporter
```
