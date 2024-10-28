package worker

import (
	"database/sql"
	"time"

	"github.com/xiusin/pinecms/src/common/river"
	"github.com/xiusin/pinecms/src/common/river/args"
)

func init() {
	river.RegisterWorker(new(CleanRecycleBinWorker))

	// 注册回收站清理任务
	river.RegisterCrontab(time.Second*10, args.CronJobArgs{Name: "回收站清理"}, river.QueueCleanRecycleBin)
}

func Start(db *sql.DB) {
	river.InitRiverJob(db)
}
