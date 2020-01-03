package sysinfo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

type SystemInfo struct {
	AppVersion string
	Logger     *logrus.Entry
}

var IP_CHECKERS = []string{
	"https://api.ipify.org/",
	"https://api.my-ip.io/ip",
	"http://ipv4bot.whatismyipaddress.com/",
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
			"cores": len(cpuInfo),
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

	hostInfo, err := host.Info()
	if err == nil {
		info["host"] = map[string]string{
			"hostname":         hostInfo.Hostname,
			"os":               hostInfo.OS,
			"platform":         hostInfo.Platform,
			"platform_version": hostInfo.PlatformVersion,
			"platform_family":  hostInfo.PlatformFamily,
			"kernel_version":   hostInfo.KernelVersion,
			"virt_role":        hostInfo.VirtualizationRole,
			"virt_system":      hostInfo.VirtualizationSystem,
		}
	}
	info["ip"] = si.GetIP()

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

func (si *SystemInfo) GetIP() string {
	for _, url := range IP_CHECKERS {
		resp, err := http.Get(url)
		if err != nil {
			si.Logger.Error(err)
			continue
		}
		defer resp.Body.Close()
		html, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			si.Logger.Error(err)
			continue
		}
		if IsIPv4(string(html)) {
			return string(html)
		}
	}

	return ""
}

func IsIPv4(host string) bool {
	parts := strings.Split(host, ".")

	if len(parts) < 4 {
		return false
	}

	for _, x := range parts {
		if i, err := strconv.Atoi(x); err == nil {
			if i < 0 || i > 255 {
				return false
			}
		} else {
			return false
		}

	}
	return true
}
