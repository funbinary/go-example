package db

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bin-work/go-example/mongo/model"
)

func InitDb(addr []string, maxpollsize uint64, username, password, source string) error {
	return dbmr.InitDb(addr, maxpollsize, username, password, source)
}

func Close(ctx context.Context) {
	dbmr.CloseDb(ctx)
}

func FindRoom(ctx context.Context, roomId string) bool {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	collection := dbmr.Mongo.Database(dbmr.DbName).Collection("record.cache")
	filter := bson.D{
		{"room", roomId},
	}
	video := &model.RecordCache{}
	collection.FindOne(ctx, filter).Decode(&video)

	return !video.Id.IsZero()
}

func GetCache(ctx context.Context, roomId string) (cache *model.RecordCache, err error) {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	collection := dbmr.Mongo.Database(dbmr.DbName).Collection("record.cache")
	filter := bson.D{
		{"room", roomId},
	}
	cache = &model.RecordCache{}
	err = collection.FindOne(ctx, filter).Decode(&cache)

	return cache, err
}

func CreateRecordCache(ctx context.Context, cache *model.RecordCache) error {
	cache.Id = primitive.NewObjectID()

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := dbmr.Mongo.Database(dbmr.DbName).Collection("record.cache").InsertOne(ctx, cache)
	if err != nil {
		return err
	}
	return nil
}

func AppendRecordForRecordCache(ctx context.Context, roomId string, record *model.Record) error {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	collection := dbmr.Mongo.Database(dbmr.DbName).Collection("record.cache")
	filter := bson.M{"room": roomId}
	update := bson.D{
		{"$push", bson.D{{"records", record}}},
	}
	_, err := collection.UpdateOne(ctx, filter, update)
	return err
}

func SetEnd(ctx context.Context, roomId, display string, startTime uint64) error {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	collection := dbmr.Mongo.Database(dbmr.DbName).Collection("record.cache")
	filter := bson.M{"room": roomId}
	update := bson.D{
		{"$set", bson.D{
			{"records.$[record].end", true},
		},
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: append(make([]interface{}, 0), bson.D{
			{"record.display", display},
			{"record.startTime", startTime},
		},
		),
	})
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}

func SetUrl(ctx context.Context, roomId, display string, startTime uint64, url string) error {

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	collection := dbmr.Mongo.Database(dbmr.DbName).Collection("record.cache")
	filter := bson.M{"room": roomId}
	update := bson.D{
		{"$set", bson.D{
			{"records.$[record].url", url},
		},
		},
	}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: append(make([]interface{}, 0), bson.D{
			{"record.display", display},
			{"record.startTime", startTime},
		},
		),
	})
	_, err := collection.UpdateOne(ctx, filter, update, opts)
	return err
}
