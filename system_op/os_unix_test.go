package system_op

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestRunCommand(t *testing.T) {
	ctx := context.Background()
	p, err := RunCommand(ctx, "cd ..\\..\\gui\\tools\\xray\\ && xray_windows.exe ws --listen 127.0.0.1:8090 --html-output results\\xray.html", func(s string) {
		fmt.Print(s)
	})
	assert.NoErrorf(t, err, "failed to start process, %v", err)

	time.Sleep(20 * time.Second)

	err = KillProcess(p)
	assert.NoErrorf(t, err, "failed to kill process, %v, %v", p.Pid, err)

	time.Sleep(20 * time.Second)
}
