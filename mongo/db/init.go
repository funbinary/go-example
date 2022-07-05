package db

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	HEART_BEAT_INTERVAL = 10
	CONNECT_TIMEOUT     = 3
	MAX_CONNIDLE_TIME   = 60
)

var dbmr = &DbManager{}
var dbmronce sync.Once

type DbManager struct {
	Mongo  *mongo.Client
	DbName string
}

func (dbm *DbManager) InitDb(addr []string, maxpollsize uint64, username, password, source string) error {
	var retryWrites bool = false
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	clientOptions := options.Client().SetHosts(addr).
		SetMaxPoolSize(maxpollsize).
		SetMinPoolSize(1).
		SetServerSelectionTimeout(time.Second * 60).
		SetHeartbeatInterval(HEART_BEAT_INTERVAL * time.Second).
		SetConnectTimeout(CONNECT_TIMEOUT * time.Second).
		SetMaxConnIdleTime(MAX_CONNIDLE_TIME * time.Second).
		SetRetryWrites(retryWrites)

	//设置用户名和密码
	dbm.DbName = source
	if len(username) > 0 && len(password) > 0 && len(source) > 0 {
		clientOptions.SetAuth(options.Credential{Username: username, Password: password, AuthSource: source, PasswordSet: true})
	}

	json_p, err := json.Marshal(clientOptions)
	if err != nil {
		logrus.Infof("json_p %v", string(json_p))
		return err
	}

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}

	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logrus.Errorf(err.Error())
		return err
	}

	logrus.WithFields(logrus.Fields{
		"addrs": addr,
	}).Infof("Connected to MongoDB!")
	dbm.Mongo = client

	return nil
}

func (dbm *DbManager) CloseDb(ctx context.Context) {
	if dbm.Mongo != nil {
		dbm.Mongo.Disconnect(ctx)
		dbm.Mongo = nil
	}
}
