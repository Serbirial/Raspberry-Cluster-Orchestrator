package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type ProcStat struct {
	PID        string  `json:"pid"`
	Command    string  `json:"command"`
	RSSMB      float64 `json:"rss_mb"`
	CPUTime    float64 `json:"cpu_time"`
	CPUPercent float64 `json:"cpu_percent"` // cpu usage percentage over interval

}

// get total CPU time (user+nice+system+idle+... all fields from /proc/stat first line)
func getTotalCPUTime() (float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}
	lines := strings.Split(string(data), "\n")
	if len(lines) == 0 {
		return 0, fmt.Errorf("/proc/stat is empty")
	}

	fields := strings.Fields(lines[0])
	if len(fields) < 8 || fields[0] != "cpu" {
		return 0, fmt.Errorf("unexpected /proc/stat format")
	}

	var total uint64 = 0
	for _, f := range fields[1:] {
		v, err := strconv.ParseUint(f, 10, 64)
		if err != nil {
			return 0, err
		}
		total += v
	}

	return float64(total), nil
}

// read process CPU time from /proc/[pid]/stat (utime + stime)
func getProcCPUTime(pid string) (float64, error) {
	statPath := filepath.Join("/proc", pid, "stat")
	data, err := os.ReadFile(statPath)
	if err != nil {
		return 0, err
	}
	fields := strings.Fields(string(data))
	if len(fields) < 17 {
		return 0, fmt.Errorf("unexpected stat format for pid %s", pid)
	}
	utime, err := strconv.ParseFloat(fields[13], 64)
	if err != nil {
		return 0, err
	}
	stime, err := strconv.ParseFloat(fields[14], 64)
	if err != nil {
		return 0, err
	}
	return utime + stime, nil
}

func getProcStats(match string) ([]ProcStat, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	type procTimes struct {
		totalCPU1 float64
		totalCPU2 float64
		procCPU1  map[string]float64
		procCPU2  map[string]float64
	}

	pt := procTimes{
		procCPU1: make(map[string]float64),
		procCPU2: make(map[string]float64),
	}

	// Sample total CPU time and proc CPU time first time
	pt.totalCPU1, err = getTotalCPUTime()
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() || !isNumeric(entry.Name()) {
			continue
		}
		pid := entry.Name()

		cmdlinePath := filepath.Join("/proc", pid, "cmdline")
		cmdlineBytes, err := os.ReadFile(cmdlinePath)
		if err != nil || len(cmdlineBytes) == 0 {
			continue
		}
		command := strings.ReplaceAll(string(cmdlineBytes), "\x00", " ")
		if !strings.Contains(command, match) {
			continue
		}

		cpuTime, err := getProcCPUTime(pid)
		if err != nil {
			continue
		}
		pt.procCPU1[pid] = cpuTime
	}

	time.Sleep(100 * time.Millisecond)

	// Sample again after interval
	pt.totalCPU2, err = getTotalCPUTime()
	if err != nil {
		return nil, err
	}

	for pid := range pt.procCPU1 {
		cpuTime2, err := getProcCPUTime(pid)
		if err != nil {
			pt.procCPU2[pid] = 0
		} else {
			pt.procCPU2[pid] = cpuTime2
		}
	}

	var results []ProcStat

	pageSize := float64(os.Getpagesize())

	for pid := range pt.procCPU1 {
		cmdlinePath := filepath.Join("/proc", pid, "cmdline")
		cmdlineBytes, err := os.ReadFile(cmdlinePath)
		if err != nil || len(cmdlineBytes) == 0 {
			continue
		}
		command := strings.ReplaceAll(string(cmdlineBytes), "\x00", " ")

		statPath := filepath.Join("/proc", pid, "stat")
		statFile, err := os.Open(statPath)
		if err != nil {
			continue
		}
		data, err := io.ReadAll(statFile)
		statFile.Close()
		if err != nil {
			continue
		}

		fields := strings.Fields(string(data))
		if len(fields) < 24 {
			continue
		}

		// RSS is number of pages
		rssPages, err := strconv.ParseInt(fields[23], 10, 64)
		if err != nil {
			continue
		}
		rssMB := float64(rssPages) * pageSize / (1024 * 1024)

		// Total CPU time in seconds (jiffies / 100)
		totalTime := pt.procCPU2[pid] / 100.0

		// Calculate CPU usage %
		totalCPUDelta := pt.totalCPU2 - pt.totalCPU1
		procCPUDelta := pt.procCPU2[pid] - pt.procCPU1[pid]
		cpuPercent := 0.0
		if totalCPUDelta > 0 {
			cpuPercent = (procCPUDelta / totalCPUDelta) * 100.0
		}

		results = append(results, ProcStat{
			PID:        pid,
			Command:    command,
			RSSMB:      rssMB,
			CPUTime:    totalTime,
			CPUPercent: cpuPercent,
		})
	}

	return results, nil
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func exportProc(name string) ([]byte, error) {
	stats, err := getProcStats(name)
	if err != nil {
		return nil, err
	}

	jsonOutput, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		fmt.Println("JSON Marshal Error:", err)
		return nil, err
	}

	return jsonOutput, nil
}
