package worker

import (
	"context"
	"fmt"

	"github.com/riverqueue/river"
	"github.com/xiusin/pinecms/src/common/river/args"
)

type CleanRecycleBinWorker struct {
	river.WorkerDefaults[args.CronJobArgs]
}

func (w *CleanRecycleBinWorker) Work(_ context.Context, job *river.Job[args.CronJobArgs]) error {
	fmt.Println("定时清理回收站 ~")

	return nil
}
