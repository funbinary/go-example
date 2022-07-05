package main

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/bin-work/go-example/mongo/model"

	"github.com/bin-work/go-example/mongo/db"
)

func main() {
	db.InitDb([]string{"192.168.3.250:27017"},
		3,
		"root",
		"123456",
		"nuvedb")
	ctx := context.Background()

	// 增加记录
	//Cache1 := &model.RecordCache{
	//	Room:    "62b91672afdb86de4868273f",
	//	Records: []model.Record{},
	//}
	//if err := db.CreateRecordCache(ctx, Cache1); err != nil {
	//	panic(err)
	//}

	roomId := "62b91672afdb86de4868273f"
	// 增加record

	record := &model.Record{
		Display:   "789",
		StartTime: time.Now().Unix(),
		End:       false,
	}
	if err := db.AppendRecordForRecordCache(ctx, roomId, record); err != nil {
		panic(err)
	}

	// 修改record
	if err := db.SetEnd(ctx, roomId, "789", 1656389483); err != nil {
		panic(err)
	}

	if err := db.SetUrl(ctx, roomId, "789", 1656389483, "record/62ba6c15afdb86de486827a2/1a3597061716/2022-06-28/10-49-04.mp4"); err != nil {
		panic(err)
	}

	fmt.Println(db.FindRoom(ctx, roomId))
	cache, _ := db.GetCache(ctx, roomId)

	for _, v := range cache.Records {

		if v.Url != "" {
			ts := strings.TrimPrefix(v.Url, "record/"+"62ba6c15afdb86de486827a2/"+"1a3597061716/")
			ts = strings.TrimSuffix(ts, ".mp4")
			fmt.Println(ts)
			t, _ := time.Parse("2006-01-02/15-04-05", ts)
			now, _ := time.Parse("2006-01-02 15-04-05", "2022-06-28 10-49-06")
			fmt.Println(math.Abs(t.Sub(now).Seconds()))
		}

	}

}
