package workerpool

type Job struct {
	Payload func()
}

func (j *Job) Execute() {
	// Execute given function
	j.Payload()
}

func NewJob(p func()) *Job {
	return &Job{Payload: p}
}
