package logger

import (
	xlog "xorm.io/xorm/log"

	"github.com/rs/zerolog"
)

var customXormLogger *CustomXormLogger

func init() {

	customXormLogger = new(CustomXormLogger)
	customXormLogger.Logger = GetAppLogger()
	customXormLogger.showSQL = true
}

func GetCustomXormLogger() *CustomXormLogger {
	return customXormLogger
}

type CustomXormLogger struct {
	Logger  *zerolog.Logger
	showSQL bool
}

func (c *CustomXormLogger) Debug(v ...interface{}) {
	c.Logger.Debug().Msgf("%v", v...)
}
func (c *CustomXormLogger) Debugf(format string, v ...interface{}) {
	c.Logger.Debug().Msgf(format, v...)
}
func (c *CustomXormLogger) Error(v ...interface{}) {
	c.Logger.Error().Msgf("%v", v...)
}
func (c *CustomXormLogger) Errorf(format string, v ...interface{}) {
	c.Logger.Error().Msgf(format, v...)
}
func (c *CustomXormLogger) Info(v ...interface{}) {
	c.Logger.Info().Msgf("%v", v...)
}
func (c *CustomXormLogger) Infof(format string, v ...interface{}) {
	c.Logger.Info().Msgf(format, v...)
}
func (c *CustomXormLogger) Warn(v ...interface{}) {
	c.Logger.Warn().Msgf("%v", v...)
}
func (c *CustomXormLogger) Warnf(format string, v ...interface{}) {
	c.Logger.Warn().Msgf(format, v...)
}

func (c *CustomXormLogger) Level() xlog.LogLevel {
	level := c.Logger.GetLevel()
	switch level {
	case zerolog.ErrorLevel:
		return xlog.LOG_ERR
	case zerolog.DebugLevel:
		return xlog.LOG_DEBUG
	case zerolog.FatalLevel:
		return xlog.LOG_ERR
	case zerolog.InfoLevel:
		return xlog.LOG_INFO
	case zerolog.NoLevel:
		return xlog.LOG_OFF
	case zerolog.WarnLevel:
		return xlog.LOG_WARNING
	default:
		return xlog.LOG_UNKNOWN
	}
}
func (c *CustomXormLogger) SetLevel(l xlog.LogLevel) {
	c.Logger.Debug().Msgf("SetLevel called %v", l)
}

func (c *CustomXormLogger) ShowSQL(show ...bool) {
	c.showSQL = show[0]
}
func (c *CustomXormLogger) IsShowSQL() bool {
	return c.showSQL
}
