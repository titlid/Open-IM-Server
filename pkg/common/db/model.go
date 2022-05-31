package db

import (
	"Open_IM/pkg/common/config"
	//"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"fmt"
	go_redis "github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo/options"
	//	"context"
	//	"fmt"
	"github.com/garyburd/redigo/redis"
	"gopkg.in/mgo.v2"
	"time"

	"context"
	//"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	//	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB DataBases

type DataBases struct {
	MysqlDB     mysqlDB
	mgoSession  *mgo.Session
	redisPool   *redis.Pool
	mongoClient *mongo.Client
	rdb         *go_redis.Client
}

func key(dbAddress, dbName string) string {
	return dbAddress + "_" + dbName
}

func init() {
	//log.NewPrivateLog(constant.LogFileName)
	//var mgoSession *mgo.Session
	var mongoClient *mongo.Client
	var err1 error
	//mysql init
	initMysqlDB()
	// mongo init
	// "mongodb://sysop:moon@localhost/records"
	uri := "mongodb://sample.host:27017/?maxPoolSize=20&w=majority"
	if config.Config.Mongo.DBUri != "" {
		// example: mongodb://$user:$password@mongo1.mongo:27017,mongo2.mongo:27017,mongo3.mongo:27017/$DBDatabase/?replicaSet=rs0&readPreference=secondary&authSource=admin&maxPoolSize=$DBMaxPoolSize
		uri = config.Config.Mongo.DBUri
	} else {
		if config.Config.Mongo.DBPassword != "" && config.Config.Mongo.DBUserName != "" {
			uri = fmt.Sprintf("mongodb://%s:%s@%s/%s?maxPoolSize=%d", config.Config.Mongo.DBUserName, config.Config.Mongo.DBPassword, config.Config.Mongo.DBAddress[0],
				config.Config.Mongo.DBDatabase, config.Config.Mongo.DBMaxPoolSize)
		} else {
			uri = fmt.Sprintf("mongodb://%s/%s/?maxPoolSize=%d",
				config.Config.Mongo.DBAddress[0], config.Config.Mongo.DBDatabase,
				config.Config.Mongo.DBMaxPoolSize)
		}
	}
	mongoClient, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		fmt.Println(" mongo.Connect  failed, try ", utils.GetSelfFuncName(), err.Error(), uri)
		time.Sleep(time.Duration(30) * time.Second)
		mongoClient, err1 = mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
		if err1 != nil {
			fmt.Println(" mongo.Connect retry failed, panic", err.Error(), uri)
			panic(err1.Error())
		}
	}
	fmt.Println("0", utils.GetSelfFuncName(), "mongo driver client init success: ", uri)

	DB.mongoClient = mongoClient

	//mgoDailInfo := &mgo.DialInfo{
	//	Addrs:     config.Config.Mongo.DBAddress,
	//	Direct:    config.Config.Mongo.DBDirect,
	//	Timeout:   time.Second * time.Duration(config.Config.Mongo.DBTimeout),
	//	Database:  config.Config.Mongo.DBDatabase,
	//	Source:    config.Config.Mongo.DBSource,
	//	Username:  config.Config.Mongo.DBUserName,
	//	Password:  config.Config.Mongo.DBPassword,
	//	PoolLimit: config.Config.Mongo.DBMaxPoolSize,
	//}
	//mgoSession, err = mgo.DialWithInfo(mgoDailInfo)
	//
	//if err != nil {
	//
	//	mgoSession, err1 = mgo.DialWithInfo(mgoDailInfo)
	//	if err1 != nil {
	//		log.NewError(" mongo.Connect  failed, panic", err.Error())
	//		panic(err1.Error())
	//	}
	//}

	//DB.mgoSession = mgoSession
	//DB.mgoSession.SetMode(mgo.Monotonic, true)
	//c := DB.mgoSession.DB(config.Config.Mongo.DBDatabase).C(cChat)
	//err = c.EnsureIndexKey("uid")
	//if err != nil {
	//	panic(err.Error())
	//}
	//

	// redis pool init
	DB.redisPool = &redis.Pool{
		MaxIdle:     config.Config.Redis.DBMaxIdle,
		MaxActive:   config.Config.Redis.DBMaxActive,
		IdleTimeout: time.Duration(config.Config.Redis.DBIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial(
				"tcp",
				config.Config.Redis.DBAddress,
				redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
				redis.DialDatabase(0),
				redis.DialPassword(config.Config.Redis.DBPassWord),
			)
		},
	}
	DB.rdb = go_redis.NewClient(&go_redis.Options{
		Addr:     config.Config.Redis.DBAddress,
		Password: config.Config.Redis.DBPassWord, // no password set
		DB:       0,                              // use default DB
		PoolSize: 100,                            // 连接池大小
	})
	//DB.rdb = go_redis.NewClusterClient(&go_redis.ClusterOptions{
	//	Addrs:    []string{config.Config.Redis.DBAddress},
	//	PoolSize: 100,
	//	Password: config.Config.Redis.DBPassWord,
	//})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err = DB.rdb.Ping(ctx).Result()
	if err != nil {
		panic(err.Error())
	}
}
