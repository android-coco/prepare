package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/mvcc/mvccpb"
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

	go func() {
		for {
			kv.Put(context.TODO(), "/cron/jobs/job7", "i am job7")

			kv.Delete(context.TODO(), "/cron/jobs/job7")

			time.Sleep(1 * time.Second)
		}
	}()

	if getResp, err := kv.Get(context.TODO(), "/cron/jobs/job7");err !=nil {
		fmt.Println(err)
		return
	}else{
		if len(getResp.Kvs) != 0 {
			fmt.Println( "当前的值:",string(getResp.Kvs[0].Value))
		}

		watchStartRevision := getResp.Header.Revision + 1

		watcher := clientv3.NewWatcher(cli)

		fmt.Println("从该版本向后监听:",watchStartRevision)


		//5s后取消监听
		ctx,cancelFunc := context.WithCancel(context.TODO())

		time.AfterFunc(5*time.Second, func() {
			cancelFunc()
		})

		watchRespChan := watcher.Watch(ctx, "/cron/jobs/job7", clientv3.WithRev(watchStartRevision))

		for watchResp := range watchRespChan {
			for _,event := range watchResp.Events {
				switch event.Type {
				case mvccpb.PUT:
					fmt.Println("修改为:",string(event.Kv.Value),"Revision:",event.Kv.CreateRevision,event.Kv.ModRevision)
				case mvccpb.DELETE:
					fmt.Println("删除了","Revision:",event.Kv.ModRevision)
				}
			}
		}
	}

}
