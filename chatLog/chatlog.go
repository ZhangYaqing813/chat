package chatlog

import (
	_ "fmt"
	log "github.com/sirupsen/logrus"
	"os"
)

var (
	Std = log.New()
)

func Init() {
	Std.SetReportCaller(true)
	Std.SetOutput(os.Stdout)
	Std.SetFormatter(&log.TextFormatter{})
}
