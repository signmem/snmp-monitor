package g

import (
	"fmt"
	"github.com/coreos/go-log/log"
	"github.com/lestrrat/go-file-rotatelogs"
	"time"
)

var (
	Logger *log.Logger
)

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
		log.Panicf("can not open file. err:%s\n", err)
	}

	Logger := log.NewSimple(
		log.WriterSink(writer,
			"[%s] [%s] [%d] [%s:%d] >>> [%s] msg=%s\n",
			[]string{"full_time", "priority", "pid", "filename", "lineno",
				"executable", "message"}))

	return Logger
}
