package pnet

import (
	"os"
	"os/exec"
)

// CurrentCommand returns the execution file of the current program
func CurrentCommand() (*exec.Cmd, error) {
	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	execPath, err := os.Executable()
	if err != nil {
		return nil, err
	}
	execPath, err = exec.LookPath(execPath)
	if err != nil {
		return nil, err
	}
	cmd := exec.Command(execPath)
	cmd.Dir = workDir
	cmd.Args = os.Args
	err = InjectionNetFiles(cmd)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}
