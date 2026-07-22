//=============================================================================
//===
//=== Copyright (C) 2026-present Andrea Carboni
//===
//=== This source code is licensed under the Elastic License 2.0 (ELv2) available at:
//=== https://github.com/algotiqa/docs/blob/main/LICENSE.md
//=== By using this file, you agree to the terms and conditions of that license.
//=============================================================================

package core

import (
	"os"
	"strconv"
	"strings"
)

//=============================================================================

type MemoryInfo struct {
	UsedMB  int
	FreeMB  int
	TotalMB int
}

//=============================================================================

func GetMemoryInfo() *MemoryInfo {
	info := &MemoryInfo{}

	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return info
	}

	totalKB, totalOK := readMemFieldKB(data, "MemTotal")
	availKB, availOK := readMemFieldKB(data, "MemAvailable")

	if totalOK {
		info.TotalMB = int(totalKB / 1024)
	}

	if availOK {
		info.FreeMB = int(availKB / 1024)
		info.UsedMB = info.TotalMB - info.FreeMB
	}

	return info
}

//=============================================================================
//===
//=== Private methods
//===
//=============================================================================

func readMemFieldKB(data []byte, field string) (uint64, bool) {
	prefix := field + ":"
	for _, line := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(line, prefix) {
			f := strings.Fields(line)
			if len(f) >= 2 {
				if kb, err := strconv.ParseUint(f[1], 10, 64); err == nil {
					return kb, true
				}
			}
			return 0, false
		}
	}
	return 0, false
}

//=============================================================================
