package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Tasker interface {
	Exec()
}

type Worker struct {
	ID     int
	taskCh chan Tasker
	t      int64
	wg     *sync.WaitGroup
}

func NewWorker(ID int, wg *sync.WaitGroup) *Worker {
	return &Worker{
		ID:     ID,
		taskCh: make(chan Tasker),
		wg:     wg,
	}
}

// Run 通道中取任务执行
func (w *Worker) Run() {
	go func() {
		defer func() {
			fmt.Printf("ID:%d worker 关闭\n", w.ID)
			w.wg.Done()
		}()
		for task := range w.taskCh {
			fmt.Printf("ID:%d worker 执行\n", w.ID)
			task.Exec()
		}
	}()
}

type Pool struct {
	workers  []*Worker
	queue    *list.List
	mx       sync.Mutex
	wg       sync.WaitGroup
	duration int
}

// NewPool 初始化协程池，duration为等待*秒后工作协程关闭，默认5
func NewPool(num int, duration int) *Pool {
	p := &Pool{
		queue:    list.New(),
		duration: duration,
	}
	if duration <= 0 {
		p.duration = 5
	}
	for i := 0; i < num; i++ {
		p.wg.Add(1)
		w := NewWorker(i, &p.wg)
		w.Run()
		p.workers = append(p.workers, w)
	}
	p.distribute()
	return p
}

// distribute 分配任务给worker
func (p *Pool) distribute() {
	for _, w := range p.workers {
		go func(w *Worker) {
			defer func() {
				//fmt.Printf("ID:%d worker 任务分配协程等待%ds关闭\n", w.ID, p.duration)
				close(w.taskCh)
			}()
			for {
				task := p.shiftTask()
				if task == nil {
					//等待*秒后退出
					if w.t > 0 && time.Now().Unix()-w.t >= int64(p.duration) {
						return
					} else if w.t == 0 {
						w.t = time.Now().Unix()
					}
				} else {
					w.t = 0
					w.taskCh <- task
				}
			}
		}(w)
	}
}

// shiftTask 队列中取任务
func (p *Pool) shiftTask() Tasker {
	p.mx.Lock()
	defer p.mx.Unlock()
	if p.queue.Len() == 0 {
		return nil
	}
	task := p.queue.Front()
	if task == nil {
		return nil
	}
	p.queue.Remove(task)
	return task.Value.(Tasker)
}

// Submit 提交任务
func (p *Pool) Submit(task Tasker) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.queue.PushBack(task)
}

// Wait 等待所有工作协程完成关闭
func (p *Pool) Wait() {
	p.wg.Wait()
}
