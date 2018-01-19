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

	// Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition：获取被分配的物理页面、远程物理页面
	PhysicalPagesAllocated *prometheus.Desc
	PreferredNUMANodeIndex *prometheus.Desc
	RemotePhysicalPages    *prometheus.Desc
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
		PhysicalPagesAllocated: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "page", "physical_pages_allocated"),
			"The number of physical pages allocated",
			nil,
			nil,
		),
		PreferredNUMANodeIndex: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "page", "preferred_numa_node_index"),
			"The preferred NUMA node index associated with this partition",
			nil,
			nil,
		),
		RemotePhysicalPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "page", "remote_physical_pages"),
			"The number of physical pages not allocated from the preferred NUMA node",
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

	if desc, err := c.collectVmPages(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV pages metrics:", desc, err)
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

	for _, health := range dst {
		label := health.Name

		ch <- prometheus.MustNewConstMetric(
			c.HealthCritical,
			prometheus.GaugeValue,
			float64(health.HealthCritical),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.HealthOk,
			prometheus.GaugeValue,
			float64(health.HealthOk),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition ..,
type Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition struct {
	Name                   string
	PhysicalPagesAllocated uint64
	PreferredNUMANodeIndex uint64
	RemotePhysicalPages    uint64
}

func (c *HyperVCollector) collectVmPages(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_VidPerfProvider_HyperVVMVidPartition
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, page := range dst {
		label := page.Name

		ch <- prometheus.MustNewConstMetric(
			c.PhysicalPagesAllocated,
			prometheus.GaugeValue,
			float64(page.PhysicalPagesAllocated),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PreferredNUMANodeIndex,
			prometheus.GaugeValue,
			float64(page.PreferredNUMANodeIndex),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.RemotePhysicalPages,
			prometheus.GaugeValue,
			float64(page.RemotePhysicalPages),
			label,
		)

	}

	return nil, nil
}
