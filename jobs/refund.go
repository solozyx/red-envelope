package jobs

import (
	"fmt"
	"time"

	"github.com/go-redsync/redsync"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"

	"github.com/solozyx/red-envelope/comm"
	"github.com/solozyx/red-envelope/core/envelopes"
	"github.com/solozyx/red-envelope/infra"
)

// 过期红包退款 定时任务
type RefundExpiredJobStarter struct {
	infra.BaseStarter
	ticker *time.Ticker
	mutex  *redsync.Mutex
}

func (r *RefundExpiredJobStarter) Init(ctx infra.StarterContext) {
	// 创建定时器
	d := ctx.Props().GetDurationDefault("jobs.refund.interval", 1*time.Minute)
	r.ticker = time.NewTicker(d)

	// redis
	maxIdle := ctx.Props().GetIntDefault("redis.maxIdle", 2)
	maxActive := ctx.Props().GetIntDefault("redis.maxActive", 5)
	idleTimeout := ctx.Props().GetDurationDefault("redis.idleTimeout", 20*time.Second)
	addr := ctx.Props().GetDefault("redis.addr", "127.0.0.1:6379")

	pools := make([]redsync.Pool, 0)
	pool := &redis.Pool{
		Dial: func() (conn redis.Conn, e error) {
			return redis.Dial("tcp", addr)
		},
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: idleTimeout,
	}
	pools = append(pools, pool)
	rsync := redsync.New(pools)

	ip := comm.GetIP()

	r.mutex = rsync.NewMutex("lock:RefundExpired",
		redsync.SetExpiry(50*time.Second),
		redsync.SetTries(3),
		redsync.SetGenValueFunc(func() (s string, e error) {
			now := time.Now()
			logrus.Infof("节点%s正在执行过期红包的退款业务", ip)
			return fmt.Sprintf("%d:%s", now.Unix(), ip), nil
		}))
}

func (r *RefundExpiredJobStarter) Start(ctx infra.StarterContext) {
	// Go协程异步执行红包过期退款 定时任务
	go func() {
		// 迭代 r.Ticker.C channel 获取到值则触发定时任务
		for {
			c := <-r.ticker.C
			err := r.mutex.Lock()
			if err == nil {
				logrus.Debug("过期红包退款开始...", c)
				// 红包过期退款业务
				domain := new(envelopes.ExpiredEnvelopeDomain)
				domain.Expired()
			} else {
				logrus.Info("已经有节点在运行该任务,err=", err.Error())
			}
			r.mutex.Unlock()
		}
	}()
}

func (r *RefundExpiredJobStarter) Stop(ctx infra.StarterContext) {
	r.ticker.Stop()
}
