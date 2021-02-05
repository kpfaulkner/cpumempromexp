package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/process"
	log "github.com/sirupsen/logrus"
	"net/http"
	"github.com/shirou/gopsutil/v3/mem"
)

type CPUMemPromExp struct {

	memTotalUsedMetric     *prometheus.Desc
	memProcessUsedMetric     *prometheus.Desc
	memTotalFreeMetric     *prometheus.Desc

	cpuProcessPercentageUsedMetric *prometheus.Desc
	memProcessPercentageUsedMetric     *prometheus.Desc
	cpuTotalPercentageTotalMetric *prometheus.Desc
	memTotalPercentageUsedMetric     *prometheus.Desc
}

func NewCPUMemPromExp() *CPUMemPromExp {
	e := CPUMemPromExp{}

	namespace := "minecraft"

	e.cpuProcessPercentageUsedMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cpu_process_used_percentage"),
		"CPU percentage used for minecraft client",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	e.memTotalUsedMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mem_total_used"),
		"Memory used in total",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	e.memProcessUsedMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mem_process_used"),
		"Memory used for minecraft client",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	e.cpuTotalPercentageTotalMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "cpu_total_used_percentage"),
		"CPU percentage used in total",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	e.memTotalFreeMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mem_total_free"),
		"Memory free",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	e.memProcessPercentageUsedMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mem_process_used_percentage"),
		"Memory percentage used for minecraft client",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	e.memTotalPercentageUsedMetric = prometheus.NewDesc(
		prometheus.BuildFQName(namespace, "", "mem_total_used_percentage"),
		"Memory percentage used in total ",
		[]string{"metric1string"}, prometheus.Labels{"label1": "label1val"},
	)

	return &e
}

func (e *CPUMemPromExp) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.cpuProcessPercentageUsedMetric
	ch <- e.memTotalUsedMetric
	ch <- e.cpuTotalPercentageTotalMetric
	ch <- e.memTotalFreeMetric
	ch <- e.memTotalFreeMetric
	ch <- e.memProcessPercentageUsedMetric
	ch <- e.memTotalPercentageUsedMetric
}

func (e *CPUMemPromExp) Collect(ch chan<- prometheus.Metric) {

	v, err := mem.VirtualMemory()
	if err != nil {
		log.Errorf("Unable to get memory stats : %s", err.Error())
		return
	}

	ch <- prometheus.MustNewConstMetric(e.memTotalFreeMetric, prometheus.GaugeValue, float64(v.Free), "")
	ch <- prometheus.MustNewConstMetric(e.memTotalUsedMetric, prometheus.GaugeValue, float64(v.Used), "")
	ch <- prometheus.MustNewConstMetric(e.memTotalPercentageUsedMetric, prometheus.GaugeValue, v.UsedPercent, "")

	c,err := cpu.Times(false)
	if err != nil {
		log.Errorf("Unable to get cpu stats : %s", err.Error())
		return
	}

	if len(c) <= 0 {
		log.Errorf("Unable to get cpu stats, although no stats error : %s", err.Error())
		return
	}

	t := c[0].Total()
	i := c[0].Idle
	fmt.Printf("%f %f\n", t,i)
	s := c[0].System
	u := c[0].User
	totalUsed := s+u
	totalUsedPercentage := (totalUsed / t) * 100
	ch <- prometheus.MustNewConstMetric(e.cpuTotalPercentageTotalMetric, prometheus.GaugeValue, totalUsedPercentage, "")

	processes,_ := process.Processes()
	for _,p := range processes {
		n,_ := p.Name()
		if n == "notepad.exe" {
			fmt.Printf("Name %s\n", n)
			c,_ := p.CPUPercent()
			fmt.Printf("cpu percentage %f\n", c)

			m,_ := p.MemoryInfo()
			fmt.Printf("mem %d\n", m.VMS)
		}

	}

}

func main() {

	e := NewCPUMemPromExp( )

	prometheus.MustRegister(e)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":8081", nil))
}
