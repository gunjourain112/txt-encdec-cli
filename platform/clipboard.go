package platform

import (
	"fmt"
	"os/exec"
)

func CopyToClipboard(text string) error {
	tools := []struct {
		name string
		args []string
	}{
		{"wl-copy", nil},
		{"xclip", []string{"-selection", "clipboard"}},
		{"xsel", []string{"--clipboard", "--input"}},
	}

	for _, tool := range tools {
		cmd := exec.Command(tool.name, tool.args...)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			continue
		}
		if err := cmd.Start(); err != nil {
			continue
		}
		if _, err := stdin.Write([]byte(text)); err != nil {
			stdin.Close()
			continue
		}
		stdin.Close()
		if err := cmd.Wait(); err == nil {
			return nil
		}
	}
	return fmt.Errorf("no clipboard tool available")
}
