package collector

import (
	"log"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
)

// HyperVCollector is a Prometheus collector for hyper-v
type HyperVCollector struct {
	// Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary：获取虚拟机健康状态
	HealthCritical *prometheus.Desc
	HealthOk       *prometheus.Desc
}

// NewHyperVCollector ...
func NewHyperVCollector() (Collector, error) {
	return &HyperVCollector{
		HealthCritical: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "health", "health_critical"),
			"This counter represents the number of virtual machines with critical health",
			nil,
			nil,
		),
		HealthOk: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "health", "health_ok"),
			"This counter represents the number of virtual machines with ok health",
			nil,
			nil,
		),
	}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *HyperVCollector) Collect(ch chan<- prometheus.Metric) error {
	if desc, err := c.collectVmHealth(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV health status metrics:", desc, err)
		return err
	}
	return nil
}

// Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary vm health status
type Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary struct {
	Name           string
	HealthCritical uint32
	HealthOk       uint32
}

func (c *HyperVCollector) collectVmHealth(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	label := dst[0].Name

	ch <- prometheus.MustNewConstMetric(
		c.HealthCritical,
		prometheus.GaugeValue,
		float64(dst[0].HealthCritical),
		label,
	)

	ch <- prometheus.MustNewConstMetric(
		c.HealthOk,
		prometheus.GaugeValue,
		float64(dst[0].HealthOk),
		label,
	)

	return nil, nil
}
