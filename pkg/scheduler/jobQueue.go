package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"time"
)

type pendingJobs []*scrapeJob

func (pjs pendingJobs) Len() int {
	return len(pjs)
}

func (pjs pendingJobs) Swap(i, j int) {
	pjs[i], pjs[j] = pjs[j], pjs[i]
}

func (pjs pendingJobs) Less(i, j int) bool {
	return pjs[i].executeTime.Before(pjs[j].executeTime)
}

func (pjs pendingJobs) Peek() *scrapeJob {
	return pjs[0]
}

func (pjs *pendingJobs) Pop() any {
	i := pjs.Len() - 1
	out := (*pjs)[i]
	*pjs = (*pjs)[:i]
	return out
}

func (pjs *pendingJobs) Push(x any) {
	*pjs = append(*pjs, x.(*scrapeJob))
}

type jobQueue struct {
	pending    *pendingJobs
	addJobChan <-chan *scrapeJob
}

func NewJobQueue(addJobChan <-chan *scrapeJob, jobs pendingJobs) *jobQueue {
	heap.Init(&jobs)
	return &jobQueue{
		pending:    &jobs,
		addJobChan: addJobChan,
	}
}

func (q *jobQueue) Start(ctx context.Context) <-chan *scrapeJob {
	ready := make(chan *scrapeJob, 10)
	go func() {
		defer close(ready)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancelled")
				return
			case j, ok := <-q.addJobChan:
				if !ok {
					fmt.Println("add job channel cancelled")
					return
				}
				heap.Push(q.pending, j)
			default:
				for q.pending.Len() != 0 && q.pending.Peek().isDue() {
					j := heap.Pop(q.pending).(*scrapeJob)
					fmt.Printf("Ready to execute: %v at %v\n", j.scrapeConfig.ID, j.executeTime.String())
					ready <- j
				}
				time.Sleep(1 * time.Second)
			}

		}
	}()
	return ready
}
