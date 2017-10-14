static:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/openvpn_exporter openvpn_exporter.go

docker: clean static docker-image

docker-image:
	docker build -t openvpn_exporter  .

clean:
	rm bin/openvpn_exporter
