package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type TimePoint struct {
	StartTime int64 `bson:"start_time"`
	EndTime   int64 `bson:"end_time"`
}

type LogRecord struct {
	JobName   string    `bson:"job_name"`   //任务名
	Command   string    `bson:"command"`    //shell命令
	Err       string    `bson:"err"`        //脚本错误
	Content   string    `bson:"content"`    //脚本输出
	TimePoint TimePoint `bson:"time_point"` //执行时间
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		fmt.Println(err)
		return
	}
	database := client.Database("cron")

	collection := database.Collection("log")

	record := &LogRecord{
		JobName:   "job10",
		Command:   "echo hello",
		Err:       "",
		Content:   "hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(), EndTime: time.Now().Unix() + 10},
	}
	logArr := []interface{}{
		record, record, record}
	result, err := collection.InsertMany(context.TODO(), logArr)
	if err != nil {
		fmt.Println(err)
		return
	}

	// _id:默认生成一个全局唯一ID,ObjectID: 12字节的二进制
	//推特很早的时候开源的,tweet的ID
	//snowflake: 毫秒/微妙的当前时间 + 机器ID + 当前毫秒/微妙内的自增ID(每当毫秒变化了,会重置为0,继续自增)
	for _, docID := range result.InsertedIDs {
		fmt.Println("自增ID:", docID.(primitive.ObjectID).Hex())
	}
}
