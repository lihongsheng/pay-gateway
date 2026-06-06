package cron

import (
	"github.com/lihongsheng/pay-gateway/global"
	"go.uber.org/zap"
)

type log struct{}

// // Info logs routine messages about cron's operation.
//	Info(msg string, keysAndValues ...interface{})
//	// Error logs an error condition.
//	Error(err error, msg string, keysAndValues ...interface{})

func (l *log) log(keysAndValues ...interface{}) []zap.Field {
	if len(keysAndValues) == 0 {
		return []zap.Field{zap.String("source", "cron")}
	}
	keyLen := len(keysAndValues)
	if keyLen == 0 || keyLen%2 != 0 {
		//	global.GVA_LOG.Warn(fmt.Sprintf("Keyvalues must appear in pairs: %d", keyLen))
	}
	common := make([]zap.Field, 0, (keyLen/2)+1)
	for i := 0; i < keyLen; i += 2 {
		if msgKey, ok := keysAndValues[i].(string); ok {
			common = append(common, zap.Any(msgKey, keysAndValues[i+1]))
		}
	}
	if keyLen > 0 && keyLen%2 != 0 {
		common = append(common, zap.Any("any", keysAndValues[keyLen-1]))
	}
	return append(common, zap.String("source", "cron"))
}

func (l *log) Info(msg string, keysAndValues ...interface{}) {
	global.GVA_LOG.Info(msg, l.log(keysAndValues)...)
}

func (l *log) Error(err error, msg string, keysAndValues ...interface{}) {
	fields := append(l.log(keysAndValues), zap.Error(err))
	global.GVA_LOG.Error(msg, fields...)
}
