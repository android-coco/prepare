package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	output []byte
	err    error
}

func main() {

	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result
		res        *result
	)
	resultChan = make(chan *result)
	ctx, cancelFunc = context.WithCancel(context.TODO())
	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2;echo hello")
		output, err = cmd.CombinedOutput()
		resultChan <- &result{
			output: output,
			err:    err,
		}
	}()

	time.Sleep(1 * time.Second)
	cancelFunc()

	res = <-resultChan
	if res.err != nil {
		fmt.Println(res.err)
		return
	}
	fmt.Println(string(res.output))
}
