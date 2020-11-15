package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {

	var (
		expr *cronexpr.Expression
		err error
		now time.Time
		nextTime time.Time
	)
	if expr, err = cronexpr.Parse("*/5 * * * * * *");err !=nil {
		fmt.Println(err)
		return
	}

	now = time.Now()
	nextTime = expr.Next(now)

	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了:",	nextTime)
	})


	time.Sleep(5 * time.Second)
}
