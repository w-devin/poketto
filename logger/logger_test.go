package logger

import "testing"

func TestLogger(t *testing.T) {
	Logger.SetFormatter(GetDefaultTextFormatter(false, "hello"))

	Logger.Infof("test")
}
