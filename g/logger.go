package g

import (
	"fmt"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/coreos/go-log/log"
	"os"
	"time"
)

var (
	Logger *log.Logger
	Ipaddr string
	SnmpServerDict []SnmpServers
)


func GetHostname(ipaddr string) (hostname string) {
	for _, info := range SnmpServerDict {
		if ipaddr == info.IPAddr {
			return info.HostName
		}
	}
	hostname, _ = os.Hostname()
	return
}

func InitLog() *log.Logger {
	LogMaxAge := Config().LogMaxAge
	LogRotateAge := Config().LogRotateAge
	logfile := Config().LogFile
	writer, err := rotatelogs.New(
		fmt.Sprintf("%s.%s", logfile, "%Y%m%d_%H%M%S.log"),
		rotatelogs.WithLinkName(logfile),
		rotatelogs.WithMaxAge(time.Second * time.Duration(LogMaxAge)),
		rotatelogs.WithRotationTime(time.Second * time.Duration(LogRotateAge)),
	)

	if err != nil {
		panic(fmt.Errorf("error opening file: %v", err))
	}

	Logger := log.NewSimple(
		log.WriterSink(writer,
			"[%s] [%s] [%d] [%s:%d] >>> [%s] msg=%s\n",
			[]string{"full_time", "priority", "pid", "filename", "lineno",
				"executable", "message"}))

	return Logger
}
