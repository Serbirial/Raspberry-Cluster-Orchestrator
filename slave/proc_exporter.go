package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type ProcStat struct {
	PID     string  `json:"pid"`
	Command string  `json:"command"`
	RSSMB   float64 `json:"rss_mb"`
	CPUTime float64 `json:"cpu_time"`
}

func getProcStats(match string) ([]ProcStat, error) {
	entries, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}

	var results []ProcStat

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

		utime, _ := strconv.ParseFloat(fields[13], 64)
		stime, _ := strconv.ParseFloat(fields[14], 64)
		totalTime := (utime + stime) / 100.0 // jiffies to seconds

		rssPages, _ := strconv.ParseInt(fields[23], 10, 64)
		rssMB := float64(rssPages*int64(os.Getpagesize())) / (1024 * 1024)

		results = append(results, ProcStat{
			PID:     pid,
			Command: command,
			RSSMB:   rssMB,
			CPUTime: totalTime,
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
