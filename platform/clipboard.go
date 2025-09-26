package platform

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

var (
	ErrNoClipboardTool = errors.New("no clipboard tool available")
	ErrClipboardFailed = errors.New("clipboard operation failed")
)

type ClipboardManager interface {
	Copy(text string) error
}

type clipboardTool struct {
	name string
	args []string
}

type LinuxClipboardManager struct {
	tools   []clipboardTool
	timeout time.Duration
}

func NewLinuxClipboardManager() *LinuxClipboardManager {
	return &LinuxClipboardManager{
		tools: []clipboardTool{
			{"wl-copy", nil},
			{"xclip", []string{"-selection", "clipboard"}},
			{"xsel", []string{"--clipboard", "--input"}},
		},
		timeout: 5 * time.Second,
	}
}

func (m *LinuxClipboardManager) Copy(text string) error {
	if text == "" {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var lastErr error
	for _, tool := range m.tools {
		if err := m.copyWithTool(ctx, tool, text); err == nil {
			return nil
		} else {
			lastErr = err
		}
	}

	if lastErr != nil {
		return fmt.Errorf("%w: %v", ErrClipboardFailed, lastErr)
	}
	return ErrNoClipboardTool
}

func (m *LinuxClipboardManager) copyWithTool(ctx context.Context, tool clipboardTool, text string) error {
	cmd := exec.CommandContext(ctx, tool.name, tool.args...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start %s: %w", tool.name, err)
	}

	if _, err := stdin.Write([]byte(text)); err != nil {
		stdin.Close()
		return fmt.Errorf("failed to write to %s: %w", tool.name, err)
	}
	stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("%s failed: %w", tool.name, err)
	}

	return nil
}

var defaultClipboard ClipboardManager = NewLinuxClipboardManager()

func CopyToClipboard(text string) error {
	return defaultClipboard.Copy(text)
}
