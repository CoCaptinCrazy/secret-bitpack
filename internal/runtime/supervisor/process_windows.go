package supervisor

import (
	"fmt"
	"os/exec"
	"strconv"
	"syscall"
)

const processQueryLimitedInformation = 0x1000
const processStillActive = 259
const processNotFound = syscall.Errno(87)

func configureProcess(command *exec.Cmd) {
	command.SysProcAttr = &syscall.SysProcAttr{CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP}
}

func processRunning(pid int) (bool, error) {
	handle, err := syscall.OpenProcess(processQueryLimitedInformation, false, uint32(pid))
	if err != nil {
		if err == processNotFound {
			return false, nil
		}
		return false, err
	}
	defer syscall.CloseHandle(handle)
	var exitCode uint32
	if err := syscall.GetExitCodeProcess(handle, &exitCode); err != nil {
		return false, err
	}
	return exitCode == processStillActive, nil
}

func terminateProcess(pid int) error {
	command := exec.Command("taskkill.exe", "/PID", strconv.Itoa(pid), "/T", "/F")
	output, err := command.CombinedOutput()
	if err != nil {
		running, checkErr := processRunning(pid)
		if checkErr == nil && !running {
			return nil
		}
		return fmt.Errorf("%w: %s", err, output)
	}
	return nil
}
