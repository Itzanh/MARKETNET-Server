package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"os"
	"runtime"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

type HardwareId struct {
	CPUs     []HardwareIdCPU `json:"cpus"`
	RAMTotal uint64          `json:"RAMTotal"`
	Hostname string          `json:"hostname"`
	OS       string          `json:"os"`
}

type HardwareIdCPU struct {
	CPU        int32   `json:"cpu"`
	VendorID   string  `json:"vendorId"`
	Family     string  `json:"family"`
	PhysicalID string  `json:"physicalId"`
	Cores      int32   `json:"cores"`
	ModelName  string  `json:"modelName"`
	Mhz        float64 `json:"mhz"`
}

func hardwareId() string {
	h := HardwareId{}
	h.CPUs = make([]HardwareIdCPU, 0)
	cpus, _ := cpu.Info()
	for i := 0; i < len(cpus); i++ {
		cpu := HardwareIdCPU{}
		cpu.CPU = cpus[i].CPU
		cpu.VendorID = cpus[i].VendorID
		cpu.Family = cpus[i].Family
		cpu.PhysicalID = cpus[i].PhysicalID
		cpu.Cores = cpus[i].Cores
		cpu.ModelName = cpus[i].ModelName
		cpu.Mhz = cpus[i].Mhz
		h.CPUs = append(h.CPUs, cpu)
	}
	v, _ := mem.VirtualMemory()
	h.RAMTotal = v.Total
	h.Hostname, _ = os.Hostname()
	h.OS = runtime.GOOS

	hwid, _ := json.Marshal(h)

	sha512 := sha512.New()
	sha512.Write(hwid)

	return base64.StdEncoding.EncodeToString(sha512.Sum(nil))
}
