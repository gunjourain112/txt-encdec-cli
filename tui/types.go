package tui

import (
	"errors"
	"fmt"
)

type AppState int

const (
	StateSelectMode AppState = iota
	StateEnterSecret
	StateEnterText
	StateShowResult
	StateShowError
	StateWaitingToClear
)

func (s AppState) String() string {
	switch s {
	case StateSelectMode:
		return "SelectMode"
	case StateEnterSecret:
		return "EnterSecret"
	case StateEnterText:
		return "EnterText"
	case StateShowResult:
		return "ShowResult"
	case StateShowError:
		return "ShowError"
	default:
		return fmt.Sprintf("Unknown(%d)", int(s))
	}
}

type OperationMode int

const (
	ModeEncrypt OperationMode = iota
	ModeDecrypt
)

func (m OperationMode) String() string {
	switch m {
	case ModeEncrypt:
		return "Encrypt"
	case ModeDecrypt:
		return "Decrypt"
	default:
		return fmt.Sprintf("Unknown(%d)", int(m))
	}
}

type InputState struct {
	CapsLockOn   bool
	KoreanActive bool
}

func (s InputState) HasIndicators() bool {
	return s.CapsLockOn || s.KoreanActive
}

type TerminalSize struct {
	Width  int
	Height int
}

func (t TerminalSize) IsValid() bool {
	return t.Width > 0 && t.Height > 0
}

type AppConfig struct {
	MinInputWidth    int
	MaxInputWidth    int
	DefaultWidth     int
	DefaultHeight    int
	InputCharLimit   int
	MinTerminalWidth int
}

func DefaultConfig() AppConfig {
	return AppConfig{
		MinInputWidth:    50,
		MaxInputWidth:    100,
		DefaultWidth:     80,
		DefaultHeight:    24,
		InputCharLimit:   1024,
		MinTerminalWidth: 66,
	}
}

type AppError struct {
	Op  string 
	Err error  
}

func (e *AppError) Error() string {
	if e.Op == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrInvalidState     = errors.New("invalid application state")
	ErrInvalidOperation = errors.New("invalid operation")
	ErrEmptyInput       = errors.New("input cannot be empty")
)
