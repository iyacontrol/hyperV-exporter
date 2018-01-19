package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/iyacontrol/HyperV-exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"golang.org/x/sys/windows/svc"
)

// WmiCollector implements the prometheus.Collector interface.
type WmiCollector struct {
	collector collector.Collector
}

const (
	serviceName = "hyperV_exporter"
)

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "exporter", "collector_duration_seconds"),
		"hyperV_exporter: Duration of a collection.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(collector.Namespace, "exporter", "collector_success"),
		"hyperV_exporter: Whether the collector was successful.",
		[]string{"collector"},
		nil,
	)
)

// Describe sends all the descriptors of the collectors included to
// the provided channel.
func (coll WmiCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

// Collect sends the collected metrics from each of the collectors to
// prometheus. Collect could be called several times concurrently
// and thus its run is protected by a single mutex.
func (coll WmiCollector) Collect(ch chan<- prometheus.Metric) {
	go func(c collector.Collector) {
		execute(c, ch)
	}(coll.collector)
}

func execute(c collector.Collector, ch chan<- prometheus.Metric) {
	begin := time.Now()
	err := c.Collect(ch)
	duration := time.Since(begin)
	var success float64

	if err != nil {
		log.Errorf("ERROR: %s collector failed after %fs: %s", "hyperV", duration.Seconds(), err)
		success = 0
	} else {
		log.Debugf("OK: %s collector succeeded after %fs.", "hyperV", duration.Seconds())
		success = 1
	}
	ch <- prometheus.MustNewConstMetric(
		scrapeDurationDesc,
		prometheus.GaugeValue,
		duration.Seconds(),
		"hyperV",
	)
	ch <- prometheus.MustNewConstMetric(
		scrapeSuccessDesc,
		prometheus.GaugeValue,
		success,
		"hyperV",
	)
}

func loadCollector() (collector.Collector, error) {
	return collector.NewHyperVCollector()
}

func init() {
	prometheus.MustRegister(version.NewCollector("hyperV_exporter"))
}

func initWbem() {
	// This initialization prevents a memory leak on WMF 5+. See
	// https://github.com/martinlindhe/wmi_exporter/issues/77 and linked issues
	// for details.
	log.Debugf("Initializing SWbemServices")
	s, err := wmi.InitializeSWbemServices(wmi.DefaultClient)
	if err != nil {
		log.Fatal(err)
	}
	wmi.DefaultClient.SWbemServicesClient = s
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	var (
		showVersion   = flag.Bool("version", false, "Print version information.")
		listenAddress = flag.String("telemetry.addr", ":9182", "host:port for WMI exporter.")
		metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics.")
	)
	flag.Usage = usage
	flag.Parse()

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("hyperV_exporter"))
		os.Exit(0)
	}

	initWbem()

	isInteractive, err := svc.IsAnInteractiveSession()
	if err != nil {
		log.Fatal(err)
	}

	stopCh := make(chan bool)
	if !isInteractive {
		go svc.Run(serviceName, &wmiExporterService{stopCh: stopCh})
	}

	collector, err := loadCollector()
	if err != nil {
		log.Fatalf("Couldn't load collector: %s", err)
	}

	hyperVCollector := WmiCollector{collector: collector}
	prometheus.MustRegister(hyperVCollector)

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/health", healthCheck)

	// landingPage contains the HTML served at '/'.
	// TODO: Make this nicer and more informative.
	var landingPage = []byte(`<html>
	<head><title>Hyper-V exporter</title></head>
	<body>
	<h1>Hyper-V exporter</h1>
	<p><a href='` + *metricsPath + `'>Metrics</a></p>
	</body>
	</html>
	`)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	log.Infoln("Starting HyperV exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	go func() {
		log.Infoln("Starting server on", *listenAddress)
		log.Fatalf("cannot start HyperV exporter: %s", http.ListenAndServe(*listenAddress, nil))
	}()

	for {
		if <-stopCh {
			log.Info("Shutting down HyperV exporter")
			break
		}
	}

}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, `{"status":"ok"}`)
}

type wmiExporterService struct {
	stopCh chan<- bool
}

func (s *wmiExporterService) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				s.stopCh <- true
				break loop
			default:
				log.Error(fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
