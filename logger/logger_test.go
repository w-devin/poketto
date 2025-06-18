package logger

import "testing"

func TestLogger(t *testing.T) {
	Logger.SetFormatter(GetDefaultTextFormatter(false, "", false))

	Logger.Infof("test")
}
