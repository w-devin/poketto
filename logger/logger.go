package logger

import (
	"fmt"
	"github.com/w-devin/logrus"
	"io"
	"os"
	"regexp"
	"runtime"
	"strings"
)

var Logger = logrus.New()

var frameIgnored = regexp.MustCompile(`(?)(github.com/sirupsen/logrus)|(logger.go)`)

func init() {
	Init("info", false)
}

func CallerPrettifier(frame *runtime.Frame) (function string, file string) {
	file = "???"
	line := 0

	pc := make([]uintptr, 64)
	n := runtime.Callers(3, pc)
	if n != 0 {
		pc = pc[:n]
		frames := runtime.CallersFrames(pc)

		for {
			frame, more := frames.Next()
			if !frameIgnored.MatchString(frame.File) {
				file = frame.File
				line = frame.Line
				break
			}
			if !more {
				break
			}
		}
	}

	slices := strings.Split(file, "/")
	file = slices[len(slices)-1]
	return fmt.Sprintf(" [%s:%d]", file, line), ""
}

func Init(level string, disableColors bool, writers ...io.Writer) {
	Logger.SetFormatter(&logrus.TextFormatter{
		DisableQuote:     true,
		DisableSorting:   true,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		DisableColors:    disableColors,
		CallerPrettyfier: CallerPrettifier,
	})

	writers = append(writers, os.Stdout)
	Logger.SetOutput(io.MultiWriter(writers...))

	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.DebugLevel
	}
	Logger.SetLevel(lvl)
	Logger.SetReportCaller(true)
}

func Print(args ...interface{}) {
	Logger.Info(args...)
}

func Printf(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	Logger.Infof(format, args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Debugf(format string, args ...interface{}) {
	Logger.Debugf(format, args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Warnf(format string, args ...interface{}) {
	Logger.Warnf(format, args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Logger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Fatalf(format string, args ...interface{}) {
	Logger.Fatalf(format, args...)
}
