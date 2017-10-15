package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	openvpnStatus "github.com/shrikantpatnaik/go-openvpn-status"
)

type openVPNExporter struct {
	statusPath string
}

var (
	openvpnUpDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "", "up"),
		"Whether scraping OpenVPN's metrics was successful.",
		nil, nil)
	openvpnLastUpdatedDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "", "last_updated"),
		"Whether scraping OpenVPN's metrics was successful.",
		nil, nil)
	openvpnConnectedClientsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "", "connected_clients"),
		"Number Of Connected Clients", nil, nil)
	openvpnGlobalStatsDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "global_stats", "max_bcast_mcast_queue_len"),
		"Global Stats", nil, nil)
	openvpnClientConnectedSinceDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "connected_since"),
		"Client Connected Since",
		[]string{"name", "real_address"}, nil)
	openvpnClientBytesReceivedDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "bytes_received"),
		"Client Bytes Received",
		[]string{"name"}, nil)
	openvpnClientBytesSentDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "client", "bytes_sent"),
		"Client Bytes Sent",
		[]string{"name"}, nil)
	openvpnRoutingLastRegDesc = prometheus.NewDesc(
		prometheus.BuildFQName("openvpn", "routing", "last_ref"),
		"Routing last reference time",
		[]string{"name", "virtual_address", "real_address"}, nil)
)

func (e *openVPNExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- openvpnUpDesc
}

func (e *openVPNExporter) Collect(ch chan<- prometheus.Metric) {
	status, err := openvpnStatus.ParseFile(e.statusPath)
	up := 0.0
	if status.IsUp {
		up = 1
	}
	ch <- prometheus.MustNewConstMetric(
		openvpnUpDesc,
		prometheus.GaugeValue,
		up)
	if err == nil {
		ch <- prometheus.MustNewConstMetric(
			openvpnConnectedClientsDesc,
			prometheus.GaugeValue,
			float64(len(status.ClientList)))
		ch <- prometheus.MustNewConstMetric(
			openvpnGlobalStatsDesc,
			prometheus.GaugeValue,
			float64(status.GlobalStats.MaxBcastMcastQueueLen))
		ch <- prometheus.MustNewConstMetric(
			openvpnLastUpdatedDesc,
			prometheus.GaugeValue,
			float64(status.UpdatedAt.Unix()))
		for _, client := range status.ClientList {
			nameSlice := []string{client.CommonName}
			nameAndAddressSlice := append(nameSlice, client.RealAddress)
			ch <- prometheus.MustNewConstMetric(
				openvpnClientConnectedSinceDesc,
				prometheus.GaugeValue,
				float64(client.ConnectedSince.Unix()),
				nameAndAddressSlice...)
			bytesReceived, _ := strconv.ParseFloat(client.BytesReceived, 64)
			ch <- prometheus.MustNewConstMetric(
				openvpnClientBytesReceivedDesc,
				prometheus.GaugeValue,
				bytesReceived,
				nameSlice...)
			bytesSent, _ := strconv.ParseFloat(client.BytesSent, 64)
			ch <- prometheus.MustNewConstMetric(
				openvpnClientBytesSentDesc,
				prometheus.GaugeValue,
				bytesSent,
				nameSlice...)
		}
		for _, route := range status.RoutingTable {
			labelSlice := []string{route.CommonName, route.VirtualAddress, route.RealAddress}
			ch <- prometheus.MustNewConstMetric(
				openvpnRoutingLastRegDesc,
				prometheus.GaugeValue,
				float64(route.LastRef.Unix()),
				labelSlice...)
		}
	}
}

func newOpenVPNExporter(statusPath string) *openVPNExporter {
	return &openVPNExporter{
		statusPath: statusPath,
	}
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}

func main() {
	var (
		listenAddress     = flag.String("web.listen-address", ":9176", "Address to listen on for web interface and telemetry.")
		metricsPath       = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
		openvpnStatusPath = flag.String("openvpn.status_path", "examples/server.status", "Paths at which OpenVPN places its status files.")
	)
	flag.Parse()
	exporter := newOpenVPNExporter(*openvpnStatusPath)
	log.Printf("Starting OpenVPN Exporter\n")
	prometheus.MustRegister(exporter)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>OpenVPN Exporter</title></head>
			<body>
			<h1>OpenVPN Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			</body>
			</html>`))
	})
	log.Printf("Listening on %s\n", *listenAddress)
	http.ListenAndServe(*listenAddress, logRequest(http.DefaultServeMux))
}
