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

	// Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition：获取虚拟TLB页面、保存页面数据
	AddressSpaces                 *prometheus.Desc
	AttachedDevices               *prometheus.Desc
	DepositedPages                *prometheus.Desc
	DeviceDMAErrors               *prometheus.Desc
	DeviceInterruptErrors         *prometheus.Desc
	DeviceInterruptMappings       *prometheus.Desc
	DeviceInterruptThrottleEvents *prometheus.Desc
	GPAPages                      *prometheus.Desc
	GPASpaceModificationsPersec   *prometheus.Desc
	IOTLBFlushCost                *prometheus.Desc
	IOTLBFlushesPersec            *prometheus.Desc
	RecommendedVirtualTLBSize     *prometheus.Desc
	SkippedTimerTicks             *prometheus.Desc
	Value1Gdevicepages            *prometheus.Desc
	Value1GGPApages               *prometheus.Desc
	Value2Mdevicepages            *prometheus.Desc
	Value2MGPApages               *prometheus.Desc
	Value4Kdevicepages            *prometheus.Desc
	Value4KGPApages               *prometheus.Desc
	VirtualTLBFlushEntiresPersec  *prometheus.Desc
	VirtualTLBPages               *prometheus.Desc

	// Win32_PerfRawData_HvStats_HyperVHypervisor：获取逻辑处理器数量、虚拟处理器数量
	LogicalProcessors *prometheus.Desc
	VirtualProcessors *prometheus.Desc

	// Win32_PerfRawData_HvStats_HyperVHypervisorVirtualProcessor：获取宾客CPU使用率、管理程序CPU使用率、CPU空闲率（需要通过RunTime计算）
	PercentGuestRunTime      *prometheus.Desc
	PercentHypervisorRunTime *prometheus.Desc
	PercentRemoteRunTime     *prometheus.Desc
	PercentTotalRunTime      *prometheus.Desc

	// Win32_PerfRawData_NvspSwitchStats_HyperVVirtualSwitch：获取虚拟交换机信息
	BroadcastPacketsReceivedPersec         *prometheus.Desc
	BroadcastPacketsSentPersec             *prometheus.Desc
	BytesPersec                            *prometheus.Desc
	BytesReceivedPersec                    *prometheus.Desc
	BytesSentPersec                        *prometheus.Desc
	DirectedPacketsReceivedPersec          *prometheus.Desc
	DirectedPacketsSentPersec              *prometheus.Desc
	DroppedPacketsIncomingPersec           *prometheus.Desc
	DroppedPacketsOutgoingPersec           *prometheus.Desc
	ExtensionsDroppedPacketsIncomingPersec *prometheus.Desc
	ExtensionsDroppedPacketsOutgoingPersec *prometheus.Desc
	LearnedMacAddresses                    *prometheus.Desc
	LearnedMacAddressesPersec              *prometheus.Desc
	MulticastPacketsReceivedPersec         *prometheus.Desc
	MulticastPacketsSentPersec             *prometheus.Desc
	NumberofSendChannelMovesPersec         *prometheus.Desc
	NumberofVMQMovesPersec                 *prometheus.Desc
	PacketsFlooded                         *prometheus.Desc
	PacketsFloodedPersec                   *prometheus.Desc
	PacketsPersec                          *prometheus.Desc
	PacketsReceivedPersec                  *prometheus.Desc
	PacketsSentPersec                      *prometheus.Desc
	PurgedMacAddresses                     *prometheus.Desc
	PurgedMacAddressesPersec               *prometheus.Desc

	// Win32_PerfRawData_EthernetPerfProvider_HyperVLegacyNetworkAdapter：获取Legacy虚拟网络适配器信息
	AdapterBytesDropped         *prometheus.Desc
	AdapterBytesReceivedPersec  *prometheus.Desc
	AdapterBytesSentPersec      *prometheus.Desc
	AdapterFramesDropped        *prometheus.Desc
	AdapterFramesReceivedPersec *prometheus.Desc
	AdapterFramesSentPersec     *prometheus.Desc
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

		//

		PhysicalPagesAllocated: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "vid", "physical_pages_allocated"),
			"The number of physical pages allocated",
			nil,
			nil,
		),
		PreferredNUMANodeIndex: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "vid", "preferred_numa_node_index"),
			"The preferred NUMA node index associated with this partition",
			nil,
			nil,
		),
		RemotePhysicalPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "vid", "remote_physical_pages"),
			"The number of physical pages not allocated from the preferred NUMA node",
			nil,
			nil,
		),

		//

		AddressSpaces: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "address_spaces"),
			"The number of address spaces in the virtual TLB of the partition",
			nil,
			nil,
		),
		AttachedDevices: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "attached_devices"),
			"The number of devices attached to the partition",
			nil,
			nil,
		),
		DepositedPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "deposited_pages"),
			"The number of pages deposited into the partition",
			nil,
			nil,
		),
		DeviceDMAErrors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_dma_errors"),
			"An indicator of illegal DMA requests generated by all devices assigned to the partition",
			nil,
			nil,
		),
		DeviceInterruptErrors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_interrupt_errors"),
			"An indicator of illegal interrupt requests generated by all devices assigned to the partition",
			nil,
			nil,
		),
		DeviceInterruptMappings: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_interrupt_mappings"),
			"The number of device interrupt mappings used by the partition",
			nil,
			nil,
		),
		DeviceInterruptThrottleEvents: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "device_interrupt_throttle_events"),
			"The number of times an interrupt from a device assigned to the partition was temporarily throttled because the device was generating too many interrupts",
			nil,
			nil,
		),
		GPAPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "preferred_numa_node_index"),
			"The number of pages present in the GPA space of the partition (zero for root partition)",
			nil,
			nil,
		),
		GPASpaceModificationsPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "gpa_space_modifications_persec"),
			"The rate of modifications to the GPA space of the partition",
			nil,
			nil,
		),
		IOTLBFlushCost: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "io_tlb_flush_cost"),
			"The average time (in nanoseconds) spent processing an I/O TLB flush",
			nil,
			nil,
		),
		IOTLBFlushesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "io_tlb_flush_persec"),
			"The rate of flushes of I/O TLBs of the partition",
			nil,
			nil,
		),
		RecommendedVirtualTLBSize: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "recommended_virtual_tlb_size"),
			"The recommended number of pages to be deposited for the virtual TLB",
			nil,
			nil,
		),
		SkippedTimerTicks: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "physical_pages_allocated"),
			"The number of timer interrupts skipped for the partition",
			nil,
			nil,
		),
		Value1Gdevicepages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "1G_device_pages"),
			"The number of 1G pages present in the device space of the partition",
			nil,
			nil,
		),
		Value1GGPApages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "1G_gpa_pages"),
			"The number of 1G pages present in the GPA space of the partition",
			nil,
			nil,
		),
		Value2Mdevicepages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "2M_device_pages"),
			"The number of 2M pages present in the device space of the partition",
			nil,
			nil,
		),
		Value2MGPApages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "2M_gpa_pages"),
			"The number of 2M pages present in the GPA space of the partition",
			nil,
			nil,
		),
		Value4Kdevicepages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "4K_device_pages"),
			"The number of 4K pages present in the device space of the partition",
			nil,
			nil,
		),
		Value4KGPApages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "4K_gpa_pages"),
			"The number of 4K pages present in the GPA space of the partition",
			nil,
			nil,
		),
		VirtualTLBFlushEntiresPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "virtual_tlb_flush_entires_persec"),
			"The rate of flushes of the entire virtual TLB",
			nil,
			nil,
		),
		VirtualTLBPages: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "hv", "virtual_tlb_pages"),
			"The number of pages used by the virtual TLB of the partition",
			nil,
			nil,
		),

		//

		VirtualProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "processor", "virtual_processors"),
			"The number of virtual processors present in the system",
			nil,
			nil,
		),
		LogicalProcessors: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "processor", "logical_processors"),
			"The number of logical processors present in the system",
			nil,
			nil,
		),

		//

		PercentGuestRunTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rate", "guest_run_time"),
			"The percentage of time spent by the virtual processor in guest code",
			nil,
			nil,
		),
		PercentHypervisorRunTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rate", "hypervisor_run_time"),
			"The percentage of time spent by the virtual processor in hypervisor code",
			nil,
			nil,
		),
		PercentRemoteRunTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rate", "remote_run_time"),
			"The percentage of time spent by the virtual processor running on a remote node",
			nil,
			nil,
		),
		PercentTotalRunTime: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "rate", "total_run_time"),
			"The percentage of time spent by the virtual processor in guest and hypervisor code",
			nil,
			nil,
		),

		//
		BroadcastPacketsReceivedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "broadcast_packets_received_persec"),
			"This counter represents the total number of broadcast packets received per second by the virtual switch",
			nil,
			nil,
		),
		BroadcastPacketsSentPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "broadcast_packets_sent_persec"),
			"This counter represents the total number of broadcast packets sent per second by the virtual switch",
			nil,
			nil,
		),
		BytesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "bytes_persec"),
			"This counter represents the total number of bytes per second traversing the virtual switch",
			nil,
			nil,
		),
		BytesReceivedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "bytes_received_persec"),
			"This counter represents the total number of bytes received per second by the virtual switch",
			nil,
			nil,
		),
		BytesSentPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "bytes_sent_persec"),
			"This counter represents the total number of bytes sent per second by the virtual switch",
			nil,
			nil,
		),
		DirectedPacketsReceivedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "directed_packets_received_persec"),
			"This counter represents the total number of directed packets received per second by the virtual switch",
			nil,
			nil,
		),
		DirectedPacketsSentPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "directed_packets_send_persec"),
			"This counter represents the total number of directed packets sent per second by the virtual switch",
			nil,
			nil,
		),
		DroppedPacketsIncomingPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "dropped_packets_incoming_persec"),
			"This counter represents the total number of packet dropped per second by the virtual switch in the incoming direction",
			nil,
			nil,
		),
		DroppedPacketsOutgoingPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "dropped_packets_outcoming_persec"),
			"This counter represents the total number of packet dropped per second by the virtual switch in the outgoing direction",
			nil,
			nil,
		),
		ExtensionsDroppedPacketsIncomingPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "extensions_dropped_packets_incoming_persec"),
			"This counter represents the total number of packet dropped per second by the virtual switch extensions in the incoming direction",
			nil,
			nil,
		),
		ExtensionsDroppedPacketsOutgoingPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "extensions_dropped_packets_outcoming_persec"),
			"This counter represents the total number of packet dropped per second by the virtual switch extensions in the outgoing direction",
			nil,
			nil,
		),
		LearnedMacAddresses: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "learned_mac_addresses"),
			"This counter represents the total number of learned MAC addresses of the virtual switch",
			nil,
			nil,
		),
		LearnedMacAddressesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "learned_mac_addresses_persec"),
			"This counter represents the total number MAC addresses learned per second by the virtual switch",
			nil,
			nil,
		),
		MulticastPacketsReceivedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "multicast_packets_received_persec"),
			"This counter represents the total number of multicast packets received per second by the virtual switch",
			nil,
			nil,
		),
		MulticastPacketsSentPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "multicast_packets_sent_persec"),
			"This counter represents the total number of multicast packets sent per second by the virtual switch",
			nil,
			nil,
		),
		NumberofSendChannelMovesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "number_of_send_channel_moves_persec"),
			"This counter represents the total number of send channel moves per second on this virtual switch",
			nil,
			nil,
		),
		NumberofVMQMovesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "number_of_vmq_moves_persec"),
			"This counter represents the total number of VMQ moves per second on this virtual switch",
			nil,
			nil,
		),
		PacketsFlooded: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "packets_flooded"),
			"This counter represents the total number of packets flooded by the virtual switch",
			nil,
			nil,
		),
		PacketsFloodedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "packets_flooded_persec"),
			"This counter represents the total number of packets flooded per second by the virtual switch",
			nil,
			nil,
		),
		PacketsPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "packets_persec"),
			"This counter represents the total number of packets per second traversing the virtual switch",
			nil,
			nil,
		),
		PacketsReceivedPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "packets_received_persec"),
			"This counter represents the total number of packets received per second by the virtual switch",
			nil,
			nil,
		),
		PacketsSentPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "packets_sent_persec"),
			"This counter represents the total number of packets send per second by the virtual switch",
			nil,
			nil,
		),
		PurgedMacAddresses: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "purged_mac_addresses"),
			"This counter represents the total number of purged MAC addresses of the virtual switch",
			nil,
			nil,
		),
		PurgedMacAddressesPersec: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "switch", "purged_mac_addresses_persec"),
			"This counter represents the total number MAC addresses purged per second by the virtual switch",
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

	if desc, err := c.collectVmVid(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV pages metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmHv(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV hv status metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmProcessor(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV processor metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmRate(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV rate metrics:", desc, err)
		return err
	}

	if desc, err := c.collectVmSwitch(ch); err != nil {
		log.Println("[ERROR] failed collecting hyperV switch metrics:", desc, err)
		return err
	}
	return nil
}

// Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary vm health status
type Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary struct {
	HealthCritical uint32
	HealthOk       uint32
}

func (c *HyperVCollector) collectVmHealth(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_VmmsVirtualMachineStats_HyperVVirtualMachineHealthSummary
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, health := range dst {
		ch <- prometheus.MustNewConstMetric(
			c.HealthCritical,
			prometheus.GaugeValue,
			float64(health.HealthCritical),
		)

		ch <- prometheus.MustNewConstMetric(
			c.HealthOk,
			prometheus.GaugeValue,
			float64(health.HealthOk),
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

func (c *HyperVCollector) collectVmVid(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
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

// Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition ...
type Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition struct {
	Name                          string
	AddressSpaces                 uint64
	AttachedDevices               uint64
	DepositedPages                uint64
	DeviceDMAErrors               uint64
	DeviceInterruptErrors         uint64
	DeviceInterruptMappings       uint64
	DeviceInterruptThrottleEvents uint64
	GPAPages                      uint64
	GPASpaceModificationsPersec   uint64
	IOTLBFlushCost                uint64
	IOTLBFlushesPersec            uint64
	RecommendedVirtualTLBSize     uint64
	SkippedTimerTicks             uint64
	Value1Gdevicepages            uint64
	Value1GGPApages               uint64
	Value2Mdevicepages            uint64
	Value2MGPApages               uint64
	Value4Kdevicepages            uint64
	Value4KGPApages               uint64
	VirtualTLBFlushEntiresPersec  uint64
	VirtualTLBPages               uint64
}

func (c *HyperVCollector) collectVmHv(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_HvStats_HyperVHypervisorRootPartition
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, obj := range dst {
		label := obj.Name

		ch <- prometheus.MustNewConstMetric(
			c.AddressSpaces,
			prometheus.GaugeValue,
			float64(obj.AddressSpaces),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.AttachedDevices,
			prometheus.GaugeValue,
			float64(obj.AttachedDevices),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DepositedPages,
			prometheus.GaugeValue,
			float64(obj.DepositedPages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DeviceDMAErrors,
			prometheus.GaugeValue,
			float64(obj.DeviceDMAErrors),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DeviceInterruptErrors,
			prometheus.GaugeValue,
			float64(obj.DeviceInterruptErrors),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DeviceInterruptThrottleEvents,
			prometheus.GaugeValue,
			float64(obj.DeviceInterruptThrottleEvents),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.GPAPages,
			prometheus.GaugeValue,
			float64(obj.GPAPages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.GPASpaceModificationsPersec,
			prometheus.GaugeValue,
			float64(obj.GPASpaceModificationsPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IOTLBFlushCost,
			prometheus.GaugeValue,
			float64(obj.IOTLBFlushCost),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.IOTLBFlushesPersec,
			prometheus.GaugeValue,
			float64(obj.IOTLBFlushesPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.RecommendedVirtualTLBSize,
			prometheus.GaugeValue,
			float64(obj.RecommendedVirtualTLBSize),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.SkippedTimerTicks,
			prometheus.GaugeValue,
			float64(obj.SkippedTimerTicks),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Value1Gdevicepages,
			prometheus.GaugeValue,
			float64(obj.Value1Gdevicepages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Value1GGPApages,
			prometheus.GaugeValue,
			float64(obj.Value1GGPApages),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.Value2Mdevicepages,
			prometheus.GaugeValue,
			float64(obj.Value2Mdevicepages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Value2MGPApages,
			prometheus.GaugeValue,
			float64(obj.Value2MGPApages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Value4Kdevicepages,
			prometheus.GaugeValue,
			float64(obj.Value4Kdevicepages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.Value4KGPApages,
			prometheus.GaugeValue,
			float64(obj.Value4KGPApages),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.VirtualTLBFlushEntiresPersec,
			prometheus.GaugeValue,
			float64(obj.VirtualTLBFlushEntiresPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.VirtualTLBPages,
			prometheus.GaugeValue,
			float64(obj.VirtualTLBPages),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_HvStats_HyperVHypervisor ...
type Win32_PerfRawData_HvStats_HyperVHypervisor struct {
	LogicalProcessors uint64
	VirtualProcessors uint64
}

func (c *HyperVCollector) collectVmProcessor(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_HvStats_HyperVHypervisor
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, obj := range dst {

		ch <- prometheus.MustNewConstMetric(
			c.LogicalProcessors,
			prometheus.GaugeValue,
			float64(obj.LogicalProcessors),
		)

		ch <- prometheus.MustNewConstMetric(
			c.VirtualProcessors,
			prometheus.GaugeValue,
			float64(obj.VirtualProcessors),
		)

	}

	return nil, nil
}

// Win32_PerfRawData_HvStats_HyperVHypervisorRootVirtualProcessor ...
type Win32_PerfRawData_HvStats_HyperVHypervisorRootVirtualProcessor struct {
	Name                     string
	PercentGuestRunTime      uint64
	PercentHypervisorRunTime uint64
	PercentRemoteRunTime     uint64
	PercentTotalRunTime      uint64
}

func (c *HyperVCollector) collectVmRate(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_HvStats_HyperVHypervisorRootVirtualProcessor
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, obj := range dst {
		label := obj.Name

		ch <- prometheus.MustNewConstMetric(
			c.PercentGuestRunTime,
			prometheus.GaugeValue,
			float64(obj.PercentGuestRunTime),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PercentHypervisorRunTime,
			prometheus.GaugeValue,
			float64(obj.PercentHypervisorRunTime),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PercentRemoteRunTime,
			prometheus.GaugeValue,
			float64(obj.PercentRemoteRunTime),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PercentTotalRunTime,
			prometheus.GaugeValue,
			float64(obj.PercentTotalRunTime),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_NvspSwitchStats_HyperVVirtualSwitch ...
type Win32_PerfRawData_NvspSwitchStats_HyperVVirtualSwitch struct {
	Name                                   string
	BroadcastPacketsReceivedPersec         uint64
	BroadcastPacketsSentPersec             uint64
	BytesPersec                            uint64
	BytesReceivedPersec                    uint64
	BytesSentPersec                        uint64
	DirectedPacketsReceivedPersec          uint64
	DirectedPacketsSentPersec              uint64
	DroppedPacketsIncomingPersec           uint64
	DroppedPacketsOutgoingPersec           uint64
	ExtensionsDroppedPacketsIncomingPersec uint64
	ExtensionsDroppedPacketsOutgoingPersec uint64
	LearnedMacAddresses                    uint64
	LearnedMacAddressesPersec              uint64
	MulticastPacketsReceivedPersec         uint64
	MulticastPacketsSentPersec             uint64
	NumberofSendChannelMovesPersec         uint64
	NumberofVMQMovesPersec                 uint64
	PacketsFlooded                         uint64
	PacketsFloodedPersec                   uint64
	PacketsPersec                          uint64
	PacketsReceivedPersec                  uint64
	PacketsSentPersec                      uint64
	PurgedMacAddresses                     uint64
	PurgedMacAddressesPersec               uint64
}

func (c *HyperVCollector) collectVmSwitch(ch chan<- prometheus.Metric) (*prometheus.Desc, error) {
	var dst []Win32_PerfRawData_NvspSwitchStats_HyperVVirtualSwitch
	if err := wmi.Query(wmi.CreateQuery(&dst, ""), &dst); err != nil {
		return nil, err
	}

	for _, obj := range dst {
		label := obj.Name

		ch <- prometheus.MustNewConstMetric(
			c.BroadcastPacketsReceivedPersec,
			prometheus.GaugeValue,
			float64(obj.BroadcastPacketsReceivedPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.BroadcastPacketsSentPersec,
			prometheus.GaugeValue,
			float64(obj.BroadcastPacketsSentPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.BytesPersec,
			prometheus.GaugeValue,
			float64(obj.BytesPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.BytesReceivedPersec,
			prometheus.GaugeValue,
			float64(obj.BytesReceivedPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.BytesSentPersec,
			prometheus.GaugeValue,
			float64(obj.BytesSentPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DirectedPacketsReceivedPersec,
			prometheus.GaugeValue,
			float64(obj.DirectedPacketsReceivedPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DirectedPacketsSentPersec,
			prometheus.GaugeValue,
			float64(obj.DirectedPacketsSentPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.DroppedPacketsIncomingPersec,
			prometheus.GaugeValue,
			float64(obj.DroppedPacketsIncomingPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.DroppedPacketsOutgoingPersec,
			prometheus.GaugeValue,
			float64(obj.DroppedPacketsOutgoingPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ExtensionsDroppedPacketsIncomingPersec,
			prometheus.GaugeValue,
			float64(obj.ExtensionsDroppedPacketsIncomingPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.ExtensionsDroppedPacketsOutgoingPersec,
			prometheus.GaugeValue,
			float64(obj.ExtensionsDroppedPacketsOutgoingPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.LearnedMacAddresses,
			prometheus.GaugeValue,
			float64(obj.LearnedMacAddresses),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.LearnedMacAddressesPersec,
			prometheus.GaugeValue,
			float64(obj.LearnedMacAddressesPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.MulticastPacketsReceivedPersec,
			prometheus.GaugeValue,
			float64(obj.MulticastPacketsReceivedPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.MulticastPacketsSentPersec,
			prometheus.GaugeValue,
			float64(obj.MulticastPacketsSentPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.NumberofSendChannelMovesPersec,
			prometheus.GaugeValue,
			float64(obj.NumberofSendChannelMovesPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.NumberofVMQMovesPersec,
			prometheus.GaugeValue,
			float64(obj.NumberofVMQMovesPersec),
			label,
		)

		// ...
		ch <- prometheus.MustNewConstMetric(
			c.PacketsFlooded,
			prometheus.GaugeValue,
			float64(obj.PacketsFlooded),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PacketsFloodedPersec,
			prometheus.GaugeValue,
			float64(obj.PacketsFloodedPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PacketsPersec,
			prometheus.GaugeValue,
			float64(obj.PacketsPersec),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PacketsReceivedPersec,
			prometheus.GaugeValue,
			float64(obj.PacketsReceivedPersec),
			label,
		)
		ch <- prometheus.MustNewConstMetric(
			c.PurgedMacAddresses,
			prometheus.GaugeValue,
			float64(obj.PurgedMacAddresses),
			label,
		)

		ch <- prometheus.MustNewConstMetric(
			c.PurgedMacAddressesPersec,
			prometheus.GaugeValue,
			float64(obj.PurgedMacAddressesPersec),
			label,
		)

	}

	return nil, nil
}

// Win32_PerfRawData_EthernetPerfProvider_HyperVLegacyNetworkAdapter ...
type Win32_PerfRawData_EthernetPerfProvider_HyperVLegacyNetworkAdapter struct {
	Name                        string
	AdapterBytesDropped         uint64
	AdapterBytesReceivedPersec  uint64
	AdapterBytesSentPersec      uint64
	AdapterFramesDropped        uint64
	AdapterFramesReceivedPersec uint64
	AdapterFramesSentPersec     uint64
}
