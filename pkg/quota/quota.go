package quota

import (
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis"
	"github.com/saeede-bellefille/quota-manager/pkg/config"
)

const (
	defaultRequstPerMinute = 10
	defaultSizePerMonth    = 50 * 1024 * 1024 // 50 MB
)

type Quota struct {
	redis     *redis.Client
	userQuota map[uint]config.UserQuota
}

func New(redisAddress string, userQuota map[uint]config.UserQuota) *Quota {
	return &Quota{
		redis: redis.NewClient(&redis.Options{
			Addr: redisAddress,
		}),
		userQuota: userQuota,
	}
}

func (q *Quota) Check(id uint, userId uint, size int64) error {
	idKey := fmt.Sprintf("id:%d", id)
	if result := q.redis.Incr(idKey); result.Err() == nil {
		q.redis.Expire(idKey, 10*time.Minute)
		if result.Val() > 1 {
			log.Print("unique err")
			return uniqueError
		}
	} else {
		log.Printf("ERROR: %v", result.Err())
		return result.Err()
	}

	userQuota, ok := q.userQuota[userId]
	if !ok {
		userQuota.RequstPerMinute = defaultRequstPerMinute
		userQuota.SizePerMonth = defaultSizePerMonth
	}
	t := time.Now()
	countKey := fmt.Sprintf("count:%d:%d", userId, t.Unix()/60)
	if result := q.redis.Incr(countKey); result.Err() == nil {
		q.redis.Expire(countKey, time.Minute)
		if result.Val() > userQuota.RequstPerMinute {
			q.redis.Decr(idKey)
			log.Print("count err")
			return countError
		}
	} else {
		log.Printf("ERROR: %v", result.Err())
		return result.Err()
	}

	sizeKey := fmt.Sprintf("size:%d:%d-%d", userId, t.Year(), t.Month())
	if result := q.redis.IncrBy(sizeKey, size); result.Err() == nil {
		q.redis.Expire(sizeKey, 31*24*time.Hour)
		if result.Val() > userQuota.SizePerMonth {
			log.Print("size err")
			q.redis.Decr(idKey)
			q.redis.Decr(countKey)
			return sizeError
		}
	} else {
		log.Printf("ERROR: %v", result.Err())
		return result.Err()
	}

	return nil
}
