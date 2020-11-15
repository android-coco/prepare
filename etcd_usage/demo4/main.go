package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	//建立一个客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		// handle error!
		fmt.Println("连接失败：", err)
		return
	}
	defer cli.Close()

	//用于读写etcd的键值对
	kv := clientv3.NewKV(cli)
	if getResp, err := kv.Get(context.TODO(), "/cron/jobs/",clientv3.WithPrefix()); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println( getResp.Kvs,getResp.Count)
	}

}
