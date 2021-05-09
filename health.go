package main

import (
	"math"
	"runtime"
	"time"
)

var start = time.Now()

const MB float64 = 1.0 * 1024 * 1024

// swagger:model imaginary_HealthStats
type HealthStats struct {
	//Number of seconds of uptime
	//
	//example: 1293
	Uptime int64 `json:"uptime"`
	//Current Allocated Memory
	//
	//example: 1.42
	AllocatedMemory float64 `json:"allocatedMemory"`
	//Total Allocated Memory
	//
	//example: 1.42
	TotalAllocatedMemory float64 `json:"totalAllocatedMemory"`
	//Number of GoRoutines
	//
	//example: 4
	Goroutines int `json:"goroutines"`
	//Completed Garbage Collector Cycles
	//
	//example: 0
	GCCycles uint32 `json:"completedGCCycles"`
	//Number of CPUs
	//
	//example: 8
	NumberOfCPUs int `json:"cpus"`
	//Max Heap Usage
	//
	//example: 63.56
	HeapSys float64 `json:"maxHeapUsage"`
	//Current heap in use
	//
	//example: 1.42
	HeapAllocated float64 `json:"heapInUse"`
	//Total number of objects in use
	//
	//example: 7745
	ObjectsInUse uint64 `json:"objectsInUse"`
	//Operating System Memory Used
	//
	//example: 68.58
	OSMemoryObtained float64 `json:"OSMemoryObtained"`
}

func GetHealthStats() *HealthStats {
	mem := &runtime.MemStats{}
	runtime.ReadMemStats(mem)

	return &HealthStats{
		Uptime:               GetUptime(),
		AllocatedMemory:      toMegaBytes(mem.Alloc),
		TotalAllocatedMemory: toMegaBytes(mem.TotalAlloc),
		Goroutines:           runtime.NumGoroutine(),
		NumberOfCPUs:         runtime.NumCPU(),
		GCCycles:             mem.NumGC,
		HeapSys:              toMegaBytes(mem.HeapSys),
		HeapAllocated:        toMegaBytes(mem.HeapAlloc),
		ObjectsInUse:         mem.Mallocs - mem.Frees,
		OSMemoryObtained:     toMegaBytes(mem.Sys),
	}
}

func GetUptime() int64 {
	return time.Now().Unix() - start.Unix()
}

func toMegaBytes(bytes uint64) float64 {
	return toFixed(float64(bytes)/MB, 2)
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
