package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
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

	//过滤条件
	delOption := options.Delete()
	delResult, err := collection.DeleteMany(context.TODO(), bson.D{{"time_point.start_time", bson.D{{"$lt", time.Now().Unix()}}}}, delOption)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("删除行数:", delResult.DeletedCount)

}
