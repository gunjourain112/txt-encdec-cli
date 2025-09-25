package platform

import (
	"os"
	"path/filepath"
	"strings"
)

func IsCapsOnLinux() bool {
	files, err := filepath.Glob("/sys/class/leds/input*::capslock/brightness")
	if err != nil || len(files) == 0 {
		return false
	}
	content, err := os.ReadFile(files[0])
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(content)) == "1"
}
