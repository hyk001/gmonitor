package procs

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/simplejia/clog"
)

func StartProc(rootPath string, cmd string, env string) (process *os.Process, err error) {
	if process, err = GetProc(cmd); err != nil || process != nil {
		return
	}

	dirname := ""
	pos := strings.Index(cmd, " ")
	if pos != -1 {
		dirname = filepath.Dir(cmd[:pos])
	} else {
		dirname = filepath.Dir(cmd)
	}

	env = strings.Trim(env, ";")
	if env == "" {
		env = "true"
	}
	cmdStr := fmt.Sprintf("cd %s; %s; nohup %s >> gmonitor.log 2>&1 &", rootPath, env, cmd)

	clog.Info("exec.Command() start............ cmd: %s", cmdStr)

	err = exec.Command("sh", "-c", cmdStr).Run()
	if err != nil {
		return
	}
	process, err = GetProc(cmd)
	if err != nil {
		return
	}
	if process != nil {
		return
	}
	content, err := ioutil.ReadFile(filepath.Join(dirname, "gmonitor.log"))
	if err != nil {
		return
	}
	err = errors.New(string(content))
	return
}

func GetProc(cmd string) (process *os.Process, err error) {
	output, err := exec.Command("ps", "-e", "-opid", "-oppid", "-ocommand").CombinedOutput()
	if err != nil {
		err = fmt.Errorf("err: %v, output: %s", err, output)
		return
	}

	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")
	pid := ""
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		_pid, _ppid, _cmd := fields[0], fields[1], strings.Join(fields[2:], " ")

		if !strings.Contains(_cmd, cmd) {
			continue
		}
		if pid == "" {
			if _ppid == "1" {
				pid = _pid
			} else {
				pid = _ppid
			}
		} else {
			if _ppid != pid {
				err = fmt.Errorf("GetProc() %s multi process exist, ppid:%v pid:%v", cmd, _ppid, pid)
				return
			}
		}
	}

	if pid == "" {
		return
	}

	i, err := strconv.Atoi(pid)
	if err != nil {
		return
	}
	process, err = os.FindProcess(i)
	return
}

func StopProc(process *os.Process) (err error) {
	if process == nil {
		return
	}
	if err = process.Kill(); err != nil {
		return
	}

	process.Release()
	return
}

func GStopProc(process *os.Process) (err error) {
	if process == nil {
		return
	}
	// SIGHUP: 1
	if err = process.Signal(syscall.Signal(1)); err != nil {
		return
	}
	process.Release()
	return
}

func CheckProc(process *os.Process) (ok bool) {
	if process == nil {
		return
	}
	err := process.Signal(syscall.Signal(0))
	if err == nil {
		ok = true
	}
	return
}
