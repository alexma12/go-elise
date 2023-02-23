package scheduler

import (
	"container/heap"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type pendingJobs []scrapeJob

func (pjs pendingJobs) Len() int {
	return len(pjs)
}

func (pjs pendingJobs) Swap(i, j int) {
	pjs[i], pjs[j] = pjs[j], pjs[i]
}

func (pjs pendingJobs) Less(i, j int) bool {
	return pjs[i].executeTime.Before(pjs[j].executeTime)
}

func (pjs *pendingJobs) Peek() any {
	return (*pjs)[0]
}

func (pjs *pendingJobs) Pop() any {
	i := pjs.Len() - 1
	out := (*pjs)[i]
	*pjs = (*pjs)[:i]
	return out
}

func (pjs *pendingJobs) Push(x any) {
	*pjs = append(*pjs, x.(scrapeJob))
}

type jobQueue struct {
	pending    *pendingJobs
	addJobChan <-chan uuid.UUID
}

func NewJobQueue(addJobChan <-chan uuid.UUID, jobs pendingJobs) *jobQueue {
	heap.Init(&jobs)
	return &jobQueue{
		pending:    &jobs,
		addJobChan: addJobChan,
	}
}

func (q *jobQueue) Start(ctx context.Context) <-chan uuid.UUID {
	ready := make(chan uuid.UUID)
	go func() {
		defer close(ready)
		for {
			select {
			case <-ctx.Done():
				fmt.Println("cancelled")
				return
			case j := <-q.addJobChan:
				heap.Push(q.pending, j)
			default:
				for q.pending.Len() != 0 && q.pending.Peek().(scrapeJob).isDue() {
					j := heap.Pop(q.pending).(scrapeJob)
					fmt.Printf("TODO: execute: %v \n", j.id)
				}
				time.Sleep(1 * time.Second)
			}

		}
	}()
	return ready
}
