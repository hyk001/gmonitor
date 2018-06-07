// 进程监控服务
package main

import (
	"fmt"
	"time"

	"github.com/simplejia/clog"
	_ "github.com/simplejia/cmonitor/clog"
	"github.com/simplejia/cmonitor/comm"
	"github.com/simplejia/cmonitor/conf"
	"github.com/simplejia/cmonitor/svr"
	"github.com/simplejia/utils"
)

func request(command string, service string) {
	url := fmt.Sprintf("http://%s:%d", utils.LocalIp, conf.C.Port)
	params := map[string]string{
		"command": command,
		"service": service,
	}
	gpp := &utils.GPP{
		Uri:     url,
		Timeout: time.Second * 8,
		Params:  params,
	}
	body, err := utils.Get(gpp)
	if err != nil {
		fmt.Printf("Error: [cmonitor maybe down!] %v, %s\n", err, body)
		return
	}

	fmt.Println(string(body))
	return
}

func main() {
	switch {
	case conf.Start != "":
		request(comm.START, conf.Start)
	case conf.Stop != "":
		request(comm.STOP, conf.Stop)
	case conf.Restart != "":
		request(comm.RESTART, conf.Restart)
	case conf.GraceRestart != "":
		request(comm.GRESTART, conf.GraceRestart)
	case conf.Status != "":
		request(comm.STATUS, conf.Status)
	default:
		clog.Info("main() StartSvr")
		svr.StartSvr()
	}
}
