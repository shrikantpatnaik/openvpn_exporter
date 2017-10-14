FROM scratch

ADD bin/openvpn_exporter /
ADD examples/server.status /etc/openvpn_exporter/

EXPOSE 9176

CMD ["/openvpn_exporter", "--openvpn.status_path=/etc/openvpn_exporter/server.status"]
