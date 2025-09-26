package platform

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

const (
	ClipboardExecTimeout    = 4 * time.Second
	ClipboardAutoClearDelay = 11 * time.Second
)

var (
	ErrNoClipboardTool = errors.New("no clipboard tool available")
	ErrClipboardFailed = errors.New("clipboard operation failed")
)

type ClipboardManager interface {
	Copy(text string) error
	Read() (string, error)
}

type clipboardTool struct {
	name     string
	copyArgs []string
	readArgs []string
}

type LinuxClipboardManager struct {
	tools   []clipboardTool
	timeout time.Duration
}

func NewLinuxClipboardManager() *LinuxClipboardManager {
	return &LinuxClipboardManager{
		tools: []clipboardTool{
			{"wl-copy", nil, []string{"wl-paste"}},
			{"xclip", []string{"-selection", "clipboard"}, []string{"-selection", "clipboard", "-o"}},
			{"xsel", []string{"--clipboard", "--input"}, []string{"--clipboard", "--output"}},
		},
		timeout: ClipboardExecTimeout,
	}
}

func (m *LinuxClipboardManager) Copy(text string) error {
	if text == "" {
		return m.clearClipboard()
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var lastErr error
	for _, tool := range m.tools {
		if err := m.copyWithTool(ctx, tool, text); err == nil {
			m.startBackgroundClear(text, ClipboardAutoClearDelay)
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

func (m *LinuxClipboardManager) clearClipboard() error {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var lastErr error
	for _, tool := range m.tools {
		if err := m.copyWithTool(ctx, tool, ""); err == nil {
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
	cmd := exec.CommandContext(ctx, tool.name, tool.copyArgs...)

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

func (m *LinuxClipboardManager) Read() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	var lastErr error
	for _, tool := range m.tools {
		if content, err := m.readWithTool(ctx, tool); err == nil {
			return content, nil
		} else {
			lastErr = err
		}
	}

	if lastErr != nil {
		return "", fmt.Errorf("%w: %v", ErrClipboardFailed, lastErr)
	}
	return "", ErrNoClipboardTool
}

func (m *LinuxClipboardManager) readWithTool(ctx context.Context, tool clipboardTool) (string, error) {
	var cmd *exec.Cmd
	if tool.name == "wl-copy" {
		cmd = exec.CommandContext(ctx, "wl-paste")
	} else {
		cmd = exec.CommandContext(ctx, tool.name, tool.readArgs...)
	}

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("%s failed: %w", tool.name, err)
	}

	return string(output), nil
}

func (m *LinuxClipboardManager) startBackgroundClear(originalText string, delay time.Duration) {
	cmd := exec.Command("./auto_clear.sh", fmt.Sprintf("%d", int(delay.Seconds())), originalText)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Start()
}

var defaultClipboard ClipboardManager = NewLinuxClipboardManager()

func CopyToClipboard(text string) error {
	return defaultClipboard.Copy(text)
}

func ClearClipboard() error {
	return defaultClipboard.Copy("")
}
