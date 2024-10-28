package job

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/robfig/cron/v3"
)

const (
	QueueDefault          = river.QueueDefault
	QueueProductShareLock = "product_share_lock"
	QueueProductCrawlData = "product_crawl_data"
	QueueDouPlusValid     = "dou_plus_valid"
	QueueCleanRecycleBin  = "clean_recycle_bin"
)

var riverClient *river.Client[*sql.Tx]
var workers = river.NewWorkers()
var riverOnceLocker sync.Once
var periodicJobs []*river.PeriodicJob

// RegisterWorker 注册river任务worker
func RegisterWorker[T river.JobArgs, W river.Worker[T]](work W) {
	fmt.Println("注册工作者", reflect.TypeOf(work).Elem().Name(), river.AddWorkerSafely(workers, work))
}

// Enqueue 任务入队
func Enqueue(ctx context.Context, queue string, args river.JobArgs, delays ...time.Duration) error {
	if len(queue) == 0 {
		queue = QueueDefault
	}

	opt := &river.InsertOpts{Queue: queue}

	if len(delays) > 0 { // 延迟时间
		opt.ScheduledAt.Add(delays[0])
	}

	_, err := GetRiverClient().Insert(ctx, args, opt)
	return err
}

// RegisterCrontab 注册周期性任务 t:任务类型 time.Duration|cron.Schedule，args:任务参数
func RegisterCrontab[T time.Duration | string, A river.JobArgs](t T, args A, queueName ...string) {
	var schedule river.PeriodicSchedule
	var err error

	if len(queueName) == 0 {
		queueName = append(queueName, QueueDefault)
	}

	switch t := any(t).(type) {
	case time.Duration:
		schedule = river.PeriodicInterval(t)
	case string:
		schedule, err = cron.ParseStandard(t)
	default:
		panic(errors.New("t must be time.Duration or cron.Schedule"))
	}

	if err != nil {
		panic(err)
	}

	periodicJobs = append(periodicJobs, river.NewPeriodicJob(
		schedule,
		func() (river.JobArgs, *river.InsertOpts) {
			return args, &river.InsertOpts{Queue: queueName[0]}
		},
		&river.PeriodicJobOpts{RunOnStart: true},
	))
}

func InitRiverJob(db *sql.DB) {
	riverOnceLocker.Do(func() {
		ctx := context.Background()
		var err error
		riverClient, err = river.NewClient(
			riverdatabasesql.New(db),
			&river.Config{
				Workers:      workers,      // 工作者
				PeriodicJobs: periodicJobs, // 周期性任务
				Queues: map[string]river.QueueConfig{
					// TODO: 配置队列
					QueueDefault:          {MaxWorkers: 10},
					QueueProductShareLock: {MaxWorkers: 10},
					QueueCleanRecycleBin:  {MaxWorkers: 1},
					QueueDouPlusValid:     {MaxWorkers: 10},
				},
			},
		)
		if err != nil {
			panic(err)
		}

		// 启动任务调度
		fmt.Println("启动river任务调度...", riverClient.Start(ctx))
	})
}

func GetRiverClient() *river.Client[*sql.Tx] {
	if riverClient == nil {
		panic(errors.New("riverClient is nil"))
	}

	return riverClient
}
