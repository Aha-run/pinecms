package args

type CronJobArgs struct {
	Name string `json:"name" river:"unique"`
}

func (CronJobArgs) Kind() string { return "cron" }
