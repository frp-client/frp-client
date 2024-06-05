package utils

import (
	"os/exec"
	"runtime"
	"strings"
)

func OpenUrl(openUrl string) {
	var osName = strings.ToLower(runtime.GOOS)
	switch osName {
	case "windows":
		cmd := exec.Command("cmd", "/c", "start", openUrl)
		_ = cmd.Run()
	case "darwin":
		cmd := exec.Command("open", openUrl)
		_ = cmd.Run()
	}
}
