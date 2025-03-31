//go:build !windows

package system_op

import (
	"context"
	"fmt"
	"github.com/w-devin/poketto/logger"
	"os"
	"os/exec"
)

func RunCommand(ctx context.Context, command string, onNewOutput func(string)) (*os.Process, error) {
	logger.Debugf("execute command: %s", command)

	cmd := exec.CommandContext(ctx, "bash", "-c", command)

	// 命令的错误输出和标准输出都连接到同一个管道
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return nil, err
	}

	if err = cmd.Start(); err != nil {
		return nil, err
	}

	go func() {
		// 从管道中实时获取输出并打印到终端
		for {
			tmp := make([]byte, 10240)
			n, err := stdout.Read(tmp)
			if onNewOutput != nil {
				onNewOutput(string(tmp[:n]))
			}
			if err != nil {
				break
			}
		}

		if err = cmd.Wait(); err != nil {
			if onNewOutput != nil {
				onNewOutput(fmt.Sprintf("error when execute [%s], %v", command, err))
			}
			return
		}
	}()

	return cmd.Process, nil
}

func KillProcess(process *os.Process) error {
	if err := process.Kill(); err != nil {
		return err
	}

	return nil
}
