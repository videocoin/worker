package sysinfo

import (
	"encoding/json"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
)

type SystemInfo struct {
	AppVersion string
}

func (si *SystemInfo) GetInfo() (map[string]interface{}, []byte, error) {
	info := map[string]interface{}{}

	vm, err := mem.VirtualMemory()
	if err == nil {
		info["memory"] = map[string]uint64{
			"total": vm.Total,
			"free":  vm.Free,
			"used":  vm.Used,
		}
	}

	cpuInfo, err := cpu.Info()
	if err == nil {
		info["cpu"] = map[string]interface{}{
			"cores": cpuInfo[0].Cores,
			"freq":  cpuInfo[0].Mhz,
		}
	}

	info["cpu_usage"] = GetCPUUsage()

	loadInfo, err := load.Avg()
	if err == nil {
		info["load"] = map[string]float64{
			"load1":  loadInfo.Load1,
			"load5":  loadInfo.Load5,
			"load15": loadInfo.Load15,
		}
	}

	info["app_version"] = si.AppVersion

	b, err := json.Marshal(info)
	if err != nil {
		return nil, nil, err
	}

	return info, b, nil
}

func GetCPUUsage() float64 {
	cpuUsage, err := cpu.Percent(time.Second, false)
	if err == nil {
		return cpuUsage[0]
	}

	return -1
}
