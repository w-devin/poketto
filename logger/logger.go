package logger

import (
	"fmt"
	"github.com/w-devin/logrus"
	"github.com/w-devin/poketto/array"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var Logger = logrus.New()

var frameIgnored = regexp.MustCompile(`(?)(github.com/w-devin/logrus)|(logger.go)`)

func init() {
	Init("info", "", false)
}

func GetDefaultTextFormatter(disableColors bool, projectName string) *logrus.TextFormatter {
	return &logrus.TextFormatter{
		DisableQuote:    true,
		DisableSorting:  true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		DisableColors:   disableColors,
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			function, file = GetCallerPrettifier(projectName)(frame)
			return fmt.Sprintf("[%s]", function), file
		},
	}
}

func GetDefaultJsonFormatter(projectName string) *logrus.JSONFormatter {
	return &logrus.JSONFormatter{
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: GetCallerPrettifier(projectName),
	}
}

func GetCallerPrettifier(projectName string) func(frame *runtime.Frame) (function string, file string) {
	return func(frame *runtime.Frame) (function string, file string) {
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

		if projectName != "" {
			projectNameIndex := array.IndexOfFold(projectName, slices)
			if projectNameIndex == -1 {
				Fatalf("projectName Error")
				os.Exit(-1)
			}

			file = strings.Join(slices[projectNameIndex:], string(filepath.Separator))
		} else {
			file = slices[len(slices)-1]
		}

		return fmt.Sprintf("%s:%d", file, line), ""
	}
}

func Init(level, projectName string, disableColors bool, writers ...io.Writer) {
	Logger.SetFormatter(GetDefaultTextFormatter(disableColors, projectName))

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
