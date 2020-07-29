package eastmoney

import _log "github.com/xiaxin/moii/log"

//  TODO 日志的使用方式
var (
	log = _log.Clone().Named("service.eastmoney")
)

func SetLog(l *_log.Logger) {
	log = l
}
