package initialize

import (
	"context"
	"fmt"
	"log"
	"time"

	"gmserver/global"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func initIndexes(ctx context.Context) error {
	// 表名:索引列表 列: "表名": [][]string{{"index1", "index2"}}

	return nil
}

func InitMongo() error {
	uri := fmt.Sprintf("mongodb://%s:%s", global.GVA_CONFIG.Mongo.Hosts[0].Host, global.GVA_CONFIG.Mongo.Hosts[0].Port)
	dbName := global.GVA_CONFIG.Mongo.Database

	// 配置连接池
	clientOptions := options.Client().ApplyURI(uri).
		SetMaxPoolSize(global.GVA_CONFIG.Mongo.MaxPoolSize).
		SetMinPoolSize(global.GVA_CONFIG.Mongo.MinPoolSize)

	// 只有在配置了用户名和密码时才设置认证
	if global.GVA_CONFIG.Mongo.Username != "" && global.GVA_CONFIG.Mongo.Password != "" {
		clientOptions.SetAuth(options.Credential{
			Username: global.GVA_CONFIG.Mongo.Username,
			Password: global.GVA_CONFIG.Mongo.Password,
		})
	}

	// 创建客户端
	MongoClient, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = MongoClient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("mongo ping错误")
	}

	mongodb := MongoClient.Database(dbName)

	// 赋值给全局变量
	global.GVA_MONGO = MongoClient
	global.GVA_MONGO_DATABASE = mongodb
	global.GVA_MONGO_PLAYER_INFO = mongodb.Collection("player_info")
	global.GVA_MONGO_SYS_USER = mongodb.Collection("sys_user")
	global.GVA_MONGO_SYS_USER_AUTHORITY = mongodb.Collection("sys_user_authority")
	global.GVA_MONGO_SYS_OPERATION_RECORDS = mongodb.Collection("sys_operation_records")
	global.GVA_MONGO_SYS_IGNORE_APIS = mongodb.Collection("sys_ignore_apis")
	global.GVA_MONGO_SYS_BASE_MENUS = mongodb.Collection("sys_base_menus")
	global.GVA_MONGO_SYS_AUTHORITY_MENUS = mongodb.Collection("sys_authority_menus")
	global.GVA_MONGO_SYS_AUTHORITIES = mongodb.Collection("sys_authorities")
	global.GVA_MONGO_SYS_APIS = mongodb.Collection("sys_apis")

	initIndexes(ctx)
	return nil
}
