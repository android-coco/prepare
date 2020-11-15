package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	var (
		leaseId clientv3.LeaseID
	)
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


	//lease实现锁自动过期
	// op操作
	//txn事务: if else then

	//1,上锁(创建租约,自动续租,拿着租约去抢占一个key)

	//申请一个lease租约
	lease := clientv3.NewLease(cli)

	//申请一个5s的一个租约
	if leaseGrantResp, err := lease.Grant(context.TODO(), 5); err != nil {
		fmt.Println(err)
		return
	} else {
		// 租约ID
		leaseId = leaseGrantResp.ID

		ctx, cancelFunc := context.WithCancel(context.TODO())

		defer cancelFunc()

		defer lease.Revoke(context.TODO(),leaseId)

		if keepRespChan, err := lease.KeepAlive(ctx, leaseId);err!=nil{
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
	}




	//用于读写etcd的键值对
	kv := clientv3.NewKV(cli)

	//事务
	txn := kv.Txn(context.TODO())


	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/jobs/job9"),"=",0)).
		Then(clientv3.OpPut("/cron/jobs/job9", "i am job9",clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/jobs/job9"))//否则抢锁失败

	if txnResp, err := txn.Commit();err != nil {
		fmt.Println(err)
		return
	}else{
		if !txnResp.Succeeded {
			fmt.Println("锁被占用:",string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
			return
		}
		//2,处理业务
		fmt.Println("处理业务")
		time.Sleep(5 * time.Second)
	}

	//3,释放锁(取消自动续租,释放租约)
	// defer 会把租约释放掉,关联的KV就被删除了


}
