package log

import (
	"fmt"
	log "github.com/cihub/seelog"
)

func logs() {
	logger, err := log.LoggerFromConfigAsFile("seelog.xml")
	if err != nil {
		fmt.Println("log 初始化失败")
	}

	log.ReplaceLogger(logger)

}
