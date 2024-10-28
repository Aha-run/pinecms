package worker

import (
	"database/sql"
	"time"

	"github.com/xiusin/pinecms/src/common/job"
	"github.com/xiusin/pinecms/src/common/job/args"
)

func init() {
	job.RegisterWorker(new(CleanRecycleBinWorker))

	// 注册回收站清理任务
	job.RegisterCrontab(time.Second*10, args.CronJobArgs{Name: "回收站清理"}, job.QueueCleanRecycleBin)
}

func Start(db *sql.DB) {
	job.InitRiverJob(db)
}
