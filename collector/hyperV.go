package collector

import (
	"log"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/prometheus/client_golang/prometheus"
)

// HyperVCollector is a Prometheus collector for hyper-v
type HyperVCollector struct {
	// OperatingSystem：获取宿主机系统版本、物理内存使用情况、虚拟内存使用情况
	PhysicalMemoryFreeBytes *prometheus.Desc
	PagingFreeBytes         *prometheus.Desc
	VirtualMemoryFreeBytes  *prometheus.Desc
	ProcessesLimit          *prometheus.Desc
	ProcessMemoryLimitBytes *prometheus.Desc
	Processes               *prometheus.Desc
	Users                   *prometheus.Desc
	PagingLimitBytes        *prometheus.Desc
	VirtualMemoryBytes      *prometheus.Desc
	VisibleMemoryBytes      *prometheus.Desc
	Time                    *prometheus.Desc
	Timezone                *prometheus.Desc

	//

}

// NewHyperVCollector ...
func NewHyperVCollector() (Collector, error) {
	return &HyperVCollector{
		PagingLimitBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "paging_limit_bytes"),
			"OperatingSystem.SizeStoredInPagingFiles",
			nil,
			nil,
		),
		PagingFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "paging_free_bytes"),
			"OperatingSystem.FreeSpaceInPagingFiles",
			nil,
			nil,
		),
		PhysicalMemoryFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "physical_memory_free_bytes"),
			"OperatingSystem.FreePhysicalMemory",
			nil,
			nil,
		),
		Time: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "time"),
			"OperatingSystem.LocalDateTime",
			nil,
			nil,
		),
		Timezone: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "timezone"),
			"OperatingSystem.LocalDateTime",
			[]string{"timezone"},
			nil,
		),
		Processes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "processes"),
			"OperatingSystem.NumberOfProcesses",
			nil,
			nil,
		),
		ProcessesLimit: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "processes_limit"),
			"OperatingSystem.MaxNumberOfProcesses",
			nil,
			nil,
		),
		ProcessMemoryLimitBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "process_memory_limix_bytes"),
			"OperatingSystem.MaxProcessMemorySize",
			nil,
			nil,
		),
		Users: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "users"),
			"OperatingSystem.NumberOfUsers",
			nil,
			nil,
		),
		VirtualMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "virtual_memory_bytes"),
			"OperatingSystem.TotalVirtualMemorySize",
			nil,
			nil,
		),
		VisibleMemoryBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "visible_memory_bytes"),
			"OperatingSystem.TotalVisibleMemorySize",
			nil,
			nil,
		),
		VirtualMemoryFreeBytes: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "os", "virtual_memory_free_bytes"),
			"OperatingSystem.FreeVirtualMemory",
			nil,
			nil,
		),
	}, nil
}

// Collect sends the metric values for each metric
// to the provided prometheus Metric channel.
func (c *HyperVCollector) Collect(ch chan<- prometheus.Metric) error {
	if desc, err := c.collectOs(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV os metrics:", desc, err)
		return err
	}
	return nil
}

type Win32_OperatingSystem struct {
	FreePhysicalMemory      uint64
	FreeSpaceInPagingFiles  uint64
	FreeVirtualMemory       uint64
	MaxNumberOfProcesses    uint32
	MaxProcessMemorySize    uint64
	NumberOfProcesses       uint32
	NumberOfUsers           uint32
	SizeStoredInPagingFiles uint64
	TotalVirtualMemorySize  uint64
	TotalVisibleMemorySize  uint64
	LocalDateTime           time.Time
}

func (c *HyperVCollector) collectOs(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_OperatingSystem
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	ch <- prometheus.MustNewConstMetric(
		c.PhysicalMemoryFreeBytes,
		prometheus.GaugeValue,
		float64(dst[0].FreePhysicalMemory*1024), // KiB -> bytes
	)

	time := dst[0].LocalDateTime

	ch <- prometheus.MustNewConstMetric(
		c.Time,
		prometheus.GaugeValue,
		float64(time.Unix()),
	)

	timezoneName, _ := time.Zone()

	ch <- prometheus.MustNewConstMetric(
		c.Timezone,
		prometheus.GaugeValue,
		1.0,
		timezoneName,
	)

	ch <- prometheus.MustNewConstMetric(
		c.PagingFreeBytes,
		prometheus.GaugeValue,
		float64(dst[0].FreeSpaceInPagingFiles*1024), // KiB -> bytes
	)

	ch <- prometheus.MustNewConstMetric(
		c.VirtualMemoryFreeBytes,
		prometheus.GaugeValue,
		float64(dst[0].FreeVirtualMemory*1024), // KiB -> bytes
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessesLimit,
		prometheus.GaugeValue,
		float64(dst[0].MaxNumberOfProcesses),
	)

	ch <- prometheus.MustNewConstMetric(
		c.ProcessMemoryLimitBytes,
		prometheus.GaugeValue,
		float64(dst[0].MaxProcessMemorySize*1024), // KiB -> bytes
	)

	ch <- prometheus.MustNewConstMetric(
		c.Processes,
		prometheus.GaugeValue,
		float64(dst[0].NumberOfProcesses),
	)

	ch <- prometheus.MustNewConstMetric(
		c.Users,
		prometheus.GaugeValue,
		float64(dst[0].NumberOfUsers),
	)

	ch <- prometheus.MustNewConstMetric(
		c.PagingLimitBytes,
		prometheus.GaugeValue,
		float64(dst[0].SizeStoredInPagingFiles*1024), // KiB -> bytes
	)

	ch <- prometheus.MustNewConstMetric(
		c.VirtualMemoryBytes,
		prometheus.GaugeValue,
		float64(dst[0].TotalVirtualMemorySize*1024), // KiB -> bytes
	)

	ch <- prometheus.MustNewConstMetric(
		c.VisibleMemoryBytes,
		prometheus.GaugeValue,
		float64(dst[0].TotalVisibleMemorySize*1024), // KiB -> bytes
	)

	return nil, nil
}
