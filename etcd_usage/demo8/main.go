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

	putOp := clientv3.OpPut("/cron/jobs/job8", "i am job8")

	if opResp, err := kv.Do(context.TODO(), putOp); err != nil {
		fmt.Println(err)
		return
	}else {
		fmt.Println("写入的Revision:",opResp.Put().Header.Revision)
	}


	getOp := clientv3.OpGet("/cron/jobs/job8")

	if getResp, err := kv.Do(context.TODO(), getOp); err != nil {
		fmt.Println(err)
		return
	}else {
		fmt.Println("数据Revision:",getResp.Get().Kvs[0].ModRevision)
		fmt.Println("数据Value:",string(getResp.Get().Kvs[0].Value))
	}
}
