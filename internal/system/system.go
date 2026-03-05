package system

import (
	"fmt"
	"os/exec"
)

type Manager interface {
	Reboot() error
	PowerOff() error
}

type RealManager struct {
	HostRoot string // Default should be /host if running in k8s
}

func NewManager(hostRoot string) *RealManager {
	return &RealManager{HostRoot: hostRoot}
}

func (m *RealManager) Reboot() error {
	return m.executeHostCommand("reboot")
}

func (m *RealManager) PowerOff() error {
	return m.executeHostCommand("poweroff")
}

func (m *RealManager) executeHostCommand(command string) error {
	// If HostRoot is set, we use chroot to execute the command on the host
	if m.HostRoot != "" && m.HostRoot != "/" {
		cmd := exec.Command("chroot", m.HostRoot, command)
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to execute %s on host: %v (output: %s)", command, err, string(output))
		}
		return nil
	}

	// Otherwise, execute directly (useful for local dev or if not in k8s)
	cmd := exec.Command(command)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to execute %s: %v (output: %s)", command, err, string(output))
	}
	return nil
}
