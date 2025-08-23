package logger

import (
	"fmt"
	"os"
	"runtime"

	"github.com/rs/zerolog"
)

var appLogger *ZeroLogger

type ZeroLogger struct {
	logger zerolog.Logger
}

func init() {
	appLogger = newZeroLogger()
}

func newZeroLogger() *ZeroLogger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	output.TimeFormat = "2006-01-02 15:04:05"
	logger := zerolog.New(output).With().Timestamp().Logger()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return &ZeroLogger{logger: logger}
}

func GetAppLogger() *zerolog.Logger {
	return &appLogger.logger
}

// Info logs an informational message.
func Info(message string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	appLogger.logger.Info().Msg(fmt.Sprintf("[%v] %v", funcName, message))
}

// Error logs an error message.
func Error(message string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	appLogger.logger.Error().Msg(fmt.Sprintf("[%v] %v", funcName, message))
}

// Debug logs a debug message.
func Debug(message string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	appLogger.logger.Debug().Msg(fmt.Sprintf("[%v] %v", funcName, message))
}

// Panic logs a panic message and panics.
func Panic(message string) {
	pc, _, _, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	appLogger.logger.Panic().Msg(fmt.Sprintf("[%v] %v", funcName, message))
}
