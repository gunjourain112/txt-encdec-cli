package platform

import (
	"os"
	"path/filepath"
	"strings"
	"unicode"
)

type SystemStateDetector interface {
	IsCapsLockOn() bool
	IsKoreanInput(r rune) bool
	IsLatinInput(r rune) bool
}

type LinuxSystemDetector struct {
	capsLockPath string
}

func NewLinuxSystemDetector() *LinuxSystemDetector {
	detector := &LinuxSystemDetector{}
	detector.findCapsLockPath()
	return detector
}

func (d *LinuxSystemDetector) findCapsLockPath() {
	files, err := filepath.Glob("/sys/class/leds/input*::capslock/brightness")
	if err == nil && len(files) > 0 {
		d.capsLockPath = files[0]
	}
}

func (d *LinuxSystemDetector) IsCapsLockOn() bool {
	if d.capsLockPath == "" {
		return false
	}

	content, err := os.ReadFile(d.capsLockPath)
	if err != nil {
		return false
	}

	return strings.TrimSpace(string(content)) == "1"
}

func (d *LinuxSystemDetector) IsKoreanInput(r rune) bool {
	if r >= 0xAC00 && r <= 0xD7A3 {
		return true
	}

	if r >= 0x3131 && r <= 0x318E {
		return true
	}

	if r >= 0x3200 && r <= 0x32FF {
		return true
	}

	return false
}

func (d *LinuxSystemDetector) HasKoreanInput(runes []rune) bool {
	for _, r := range runes {
		if d.IsKoreanInput(r) {
			return true
		}
	}
	return false
}

func (d *LinuxSystemDetector) IsLatinInput(r rune) bool {
	return unicode.In(r, unicode.Latin)
}

var defaultDetector SystemStateDetector = NewLinuxSystemDetector()

func IsCapsOnLinux() bool {
	return defaultDetector.IsCapsLockOn()
}
