package rds_lock

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"math/rand"
	"time"
)

type RedisLock struct {
	rdsCli      *redis.Client
	key, value  string
	expire      time.Duration                   //过期时间
	delayFunc   func(tryTime int) time.Duration //重试间隔
	genValuFunc func(valLen int) string         //value生成
}

const (
	maxTries = 32
	valueLen = 32
)

func defaultDelayFunc(tries int) time.Duration {
	return (time.Duration(rand.Intn(1000)) * time.Millisecond) + (time.Duration(tries) * time.Millisecond * 10)
}

func defaultGenValuFunc(valLen int) string {
	b := make([]byte, valLen)
	_, _ = rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)[:valLen]
}

func NewRedisLock(rdsCli *redis.Client, key string, expireSecond int) *RedisLock {
	return &RedisLock{
		rdsCli:    rdsCli,
		key:       key,
		value:     defaultGenValuFunc(valueLen),
		expire:    time.Duration(expireSecond) * time.Second,
		delayFunc: defaultDelayFunc,
	}
}

func (r *RedisLock) Lock() error {
	ctx := context.Background()

	for i := 0; i < maxTries; i++ {
		time.Sleep(r.delayFunc(i))

		// 尝试上锁
		ok, err := r.rdsCli.SetNX(ctx, r.key, r.value, r.expire).Result()
		if ok && err == nil {
			zap.S().Infof("第%d次上锁成功", i)
			return nil
		}

		// 上锁失败,等待后重试
		zap.S().Infof("第%d次上锁失败", i)
	}

	return errors.New(fmt.Sprintf("Lock failed, tried %d times", maxTries))
}

var deleteScript = redis.NewScript(`
	if redis.call("GET", KEYS[1]) == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
`)

func (r *RedisLock) Unlock() error {
	keys := []string{r.key}
	values := []interface{}{r.value}

	// 使用lua脚本删除,保证原子性
	_, err := deleteScript.Run(context.Background(), r.rdsCli, keys, values).Result()
	return err
}
