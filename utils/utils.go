package utils

import (
	"regexp"
	"strings"
	"syscall"
)

func InfoHashFromMagnet(magnet string) string {
	r, _ := regexp.Compile(`urn:btih:([a-fA-F0-9]{40})`)
	matches := r.FindStringSubmatch(magnet)
	if len(matches) > 1 {
		return strings.ToLower(matches[1])
	}
	return ""
}

func IsProcessRunning(pid int) bool {
	h, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(h)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(h, &exitCode)
	if err != nil {
		return false
	}

	// STILL_ACTIVE == 259
	return exitCode == 259
}
