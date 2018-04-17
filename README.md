# OpenVPN Exporter for Prometheus

A simple exporter that reads the OpenVPN status file and exposes the data as Prometheus Metrics

## Usage

### Local

```
go build openvpn_exporter.go
./openvpn_exporter --openvpn.status_path=/path/to/openvpn.status
```

### Docker

Build the image:
```
docker build --force-rm=true -t openvpn_exporter .
```

The final image is around 8MB. A temporary image has been downloaded to make the final one. Once built, this temporary image become orphan, you can delete it:
```
docker rmi -f $(docker images | grep "<none>" | awk "{print \$3}")
```

To use with docker you must mount your status file to /etc/openvpn_exporter/server.status.
```
docker run -it -p 9176:9176 -v /path/to/openvpn_server.status:/etc/openvpn_exporter/server.status openvpn_exporter
```

Metrics should be available on your host IP: http://<host_ip>:9176/metrics. E.g: http://10.39.9.94:9176/metrics

## TODO
Figure out a good way to see if the server is up, I currently just assume its down if the last update happened more than 10 minutes ago
