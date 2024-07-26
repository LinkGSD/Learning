package main

import (
	_ "embed"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type ZRanking struct {
	Redis          *redis.Client
	Key            string        // redis zset key
	Expiration     time.Duration // 数据保存过期时间
	StartTimestamp int64         // 排行活动开始时间
	EndTimestamp   int64         // 排行活动结束时间戳
	TimePadWidth   int           // 排行榜活动结束时间与用户排序值更新时间的差值补0宽度
}

// New 创建ZRanking实例
func New(rds *redis.Client, key string, startTs, endTs int64, expiration time.Duration) (*ZRanking, error) {
	deltaTs := endTs - startTs
	if deltaTs <= 0 {
		return nil, fmt.Errorf("invalid deltaTs:%v", deltaTs)
	}
	timePadWidth := len(fmt.Sprint(deltaTs))
	return &ZRanking{
		Redis:          rds,
		Key:            key,
		Expiration:     expiration,
		StartTimestamp: startTs,
		EndTimestamp:   endTs,
		TimePadWidth:   timePadWidth,
	}, nil
}

// val 转为 score:
// score = float64(val.deltaTs)
//func (r *ZRanking) val2score(ctx context.Context, val int64) (float64, error) {
//	nowts := time.Now().Unix()
//	deltaTs := r.EndTimestamp - nowts
//	if deltaTs < 0 {
//		return 0, errors.New("val2score error: deltaTs < 0")
//	}
//	scoreFormat := fmt.Sprintf("%%v.%%0%dd", r.TimePadWidth)
//	scoreStr := fmt.Sprintf(scoreFormat, val, deltaTs)
//	score, err := strconv.ParseFloat(scoreStr, 64)
//	if err != nil {
//		err = fmt.Errorf("%v,%s", err, "ZRanking val2score ParseFloat error")
//		return 0, err
//	}
//	return score, nil
//}

// 从 score 中获取 val
//func (r *ZRanking) score2val(ctx context.Context, score float64) (int64, error) {
//	scoreStr := fmt.Sprint(score)
//	ss := strings.Split(scoreStr, ".")
//	valStr := ss[0]
//	val, err := strconv.ParseInt(valStr, 10, 64)
//	if err != nil {
//		err = fmt.Errorf("%v,%s", err, "ZRanking score2val ParseInt error")
//		return 0, err
//	}
//	return val, nil
//}

//go:embed lua/add.lua
var addScript string

func main() {
	client := redis.NewClient(&redis.Options{Addr: "192.168.182.132:6379"})
	zRanking, err := New(client, "test", time.Now().Add(-3*time.Hour).Unix(), time.Now().Add((30*24-3)*time.Hour).Unix(), time.Hour)
	if err != nil {
		return
	}
	script := redis.NewScript(addScript)
	i, err := script.Eval(zRanking.Redis, []string{"test"}, 4, 1).Result()
	fmt.Println(i)
	//zs, err := zRanking.Redis.ZRevRangeWithScores("test", 0, 9).Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(zs)
	//t, err := zRanking.Redis.ZRevRank("test", "4").Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(t)
	f, err := zRanking.Redis.ZScore("test", "4").Result()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(f)
	//result, err := zRanking.Redis.ZCard("test").Result()
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(result)
}
