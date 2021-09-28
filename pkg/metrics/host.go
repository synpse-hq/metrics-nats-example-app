package metrics

import (
	"context"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/go-osstat/memory"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

func init() {
	prometheus.MustRegister(memTotal)
	prometheus.MustRegister(memUsed)
	prometheus.MustRegister(memCached)
	prometheus.MustRegister(memFree)
	prometheus.MustRegister(cpuUser)
	prometheus.MustRegister(cpuNice)
	prometheus.MustRegister(cpuSystem)
	prometheus.MustRegister(cpuCount)
	prometheus.MustRegister(cpuIdle)
	prometheus.MustRegister(cpuIrq)
	prometheus.MustRegister(cpuIowait)
	prometheus.MustRegister(cpuSteal)
}

var (
	memTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_mem_total",
			Help: "Agent total memory bytes",
		},
	)

	memUsed = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_mem_used",
			Help: "Agent used memory bytes",
		},
	)

	memCached = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_mem_cached",
			Help: "Agent total memory cached bytes",
		},
	)

	memFree = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_mem_free",
			Help: "Agent free memory",
		},
	)

	cpuUser = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_user_percentage",
			Help: "Agent user cpu usage percentage",
		},
	)

	cpuNice = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_nice_percentage",
			Help: "Agent nice cpu usage percentage",
		},
	)

	cpuSystem = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_system_percentage",
			Help: "Agent system cpu usage percentage",
		},
	)

	cpuCount = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_count",
			Help: "Agent system cpu count",
		},
	)

	cpuIdle = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_idle_percentage",
			Help: "Agent idle cpu usage percentage",
		},
	)

	cpuIowait = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_iowait_percentage",
			Help: "Agent cpu io wait percentage",
		},
	)

	cpuIrq = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_irq_percentage",
			Help: "Agent cpu interrupt requests percentage",
		},
	)

	cpuSteal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_steal_percentage",
			Help: "Agent cpu steal time percentage",
		},
	)

	cpuGuest = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "agent_cpu_guest_percentage",
			Help: "Agent cpu guest wait time percentage",
		},
	)
)

func (m *Metrics) updateMemory(ctx context.Context) {
	mem, err := memory.Get()
	if err != nil {
		m.log.Warn("failed to get host memory stats", zap.Error(err))
	}

	memTotal.Set(float64(mem.Total))
	memUsed.Set(float64(mem.Used))
	memFree.Set(float64(mem.Free))
	memCached.Set(float64(mem.Cached))
}

func (m *Metrics) updateCPU(ctx context.Context) {
	cpu, err := cpu.Get()
	if err != nil {
		m.log.Warn("failed to get host cpu stats", zap.Error(err))
	}

	firstRun := m.cpuCache == nil

	// on first run we dont have ref point so set cache and return early
	if firstRun {
		m.updateCPUCache(cpu)
		return
	}

	totalDiff := float64(cpu.Total - m.cpuCache.Total)

	// calculate cpu usage % from last run
	// TODO: math looks right but results does not... Needs checking
	// https://github.com/mackerelio/go-osstat#note-for-counter-values
	// https://github.com/mackerelio/mackerel-agent/blob/master/metrics/linux/cpuusage.go
	// https://github.com/torvalds/linux/blob/4ec9f7a18/kernel/sched/cputime.c#L151-L158

	cpuUser.Set(float64((cpu.User-cpu.Guest)-(m.cpuCache.User-m.cpuCache.Guest)) / totalDiff * 100)
	cpuNice.Set(float64(cpu.Nice-m.cpuCache.Nice) / totalDiff * 100)
	cpuSystem.Set(float64(cpu.System-m.cpuCache.System) / totalDiff * 100)
	cpuIdle.Set(float64(cpu.Idle-m.cpuCache.Idle) / totalDiff * 100)

	if cpu.StatCount >= 5 {
		cpuIowait.Set(float64(cpu.Iowait-m.cpuCache.Iowait) / totalDiff * 100)
	}
	if cpu.StatCount >= 6 {
		cpuIowait.Set(float64(cpu.Iowait-m.cpuCache.Iowait) / totalDiff * 100)
	}
	if cpu.StatCount >= 7 {
		cpuIrq.Set(float64(cpu.Softirq-m.cpuCache.Softirq) / totalDiff * 100)
	}
	if cpu.StatCount >= 8 {
		cpuSteal.Set(float64(cpu.Steal-m.cpuCache.Steal) / totalDiff * 100)
	}
	if cpu.StatCount >= 9 {
		cpuGuest.Set(float64(cpu.Guest-m.cpuCache.Guest) / totalDiff * 100)
	}
	m.updateCPUCache(cpu)
}

func (m *Metrics) updateCPUCache(cpu *cpu.Stats) {
	m.cpuCache = &cpuCache{}
	m.cpuCache.Idle = cpu.Idle
	m.cpuCache.Iowait = cpu.Iowait
	m.cpuCache.User = cpu.User
	m.cpuCache.Guest = cpu.Guest
	m.cpuCache.Nice = cpu.Nice
	m.cpuCache.System = cpu.System
	m.cpuCache.Total = cpu.Total
	m.cpuCache.Count = cpu.CPUCount
	m.cpuCache.Softirq = cpu.Softirq
	m.cpuCache.Steal = cpu.Steal
}
