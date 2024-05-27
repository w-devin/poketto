package logger

import (
	"fmt"
	"github.com/orandin/sentrus"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Logger = logrus.New()

var frameIgnored = regexp.MustCompile(`(?)(github.com/sirupsen/logrus)|(logger.go)`)

func init() {
	Init("info", "")
}

func Init(level string, fileDir string) {
	//Logger.SetFormatter(&logrus.JSONFormatter{})
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
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
		},
	})
	if fileDir == "" {
		Logger.SetOutput(os.Stdout)
	} else {
		logFile := filepath.Join(fileDir, "access.scanner")
		f, err := os.OpenFile(logFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0700)
		if err != nil {
			panic(err)
		}
		Logger.SetOutput(io.MultiWriter(os.Stdout, f))
	}
	lvl, err := logrus.ParseLevel(level)
	if err != nil {
		lvl = logrus.DebugLevel
	}
	Logger.SetLevel(lvl)
	Logger.SetReportCaller(true)

	Logger.AddHook(sentrus.NewHook([]logrus.Level{logrus.WarnLevel, logrus.ErrorLevel, logrus.FatalLevel}))
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
