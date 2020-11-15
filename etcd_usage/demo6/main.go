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

	//申请一个lease租约
	lease := clientv3.NewLease(cli)

	//申请一个10s的一个租约
	if leaseGrantResp, err := lease.Grant(context.TODO(), 10); err != nil {
		fmt.Println(err)
		return
	} else {
		// 租约ID
		leaseId := leaseGrantResp.ID

		//ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)
		if keepRespChan, err := lease.KeepAlive(context.TODO(), leaseId);err!=nil{
			fmt.Println(err)
			return
		}else{
			go func() {
				for  {

					select {
					case keepResp := <- keepRespChan:
						if keepRespChan == nil {
							fmt.Println("租约已经失效了")
							goto END
						}else{//每秒会续租一次,所以就会有一次应答
							fmt.Println("收到自动续租应答:",keepResp.ID)
						}
					}

				}
				END:
			}()
		}

		// 获得kv对象
		kv := clientv3.NewKV(cli)
		if putResp, err := kv.Put(context.TODO(), "/cron/lock/job1", "1", clientv3.WithLease(leaseId)); err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println("写入成功：", putResp.Header.Revision)
		}
	}

	//定时查看是否过期
	for {
		kv := clientv3.NewKV(cli)
		if getResp, err := kv.Get(context.TODO(), "/cron/lock/job1"); err != nil {
			fmt.Println(err)
			return
		} else {
			if getResp.Count == 0 {
				fmt.Println("已经过期了")
				break
			}
			fmt.Println("还没有过期：",getResp.Kvs)
		}
		time.Sleep(2 * time.Second)
	}


}
