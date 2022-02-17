package main

import (
	"fmt"
	"math/rand"
)

type Job struct {
	ID      int
	RandNum int
}

type Result struct {
	job *Job
	sum int
}

func main() {
	jobChan := make(chan *Job, 128)
	resultChan := make(chan *Result, 128)
	createPool(64, jobChan, resultChan)
	go func(resultChan chan *Result) {
		for result := range resultChan {
			fmt.Printf("job id:%v randum:%v result:%v\n",
				result.job.ID, result.job.RandNum, result.sum)
		}
	}(resultChan)
	var id int
	for {
		id++
		r_num := rand.Int()
		job := &Job{
			ID:      id,
			RandNum: r_num,
		}
		jobChan <- job
	}
}
func createPool(num int, jobChan chan *Job, resultChan chan *Result) {
	for i := 0; i < num; i++ {
		go func(jobChan chan *Job, resultChan chan *Result) {
			for job := range jobChan {
				r_num := job.RandNum
				var sum int
				for r_num != 0 {
					tmp := r_num % 10
					sum += tmp
					r_num /= 10
				}
				r := &Result{
					job: job,
					sum: sum,
				}
				resultChan <- r
			}
		}(jobChan, resultChan)
	}
}
