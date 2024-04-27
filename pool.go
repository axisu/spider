package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type Task func()

type Worker struct {
	ID     int
	taskCh chan Task
	t      int64
}

func NewWorker(ID int) *Worker {
	return &Worker{
		ID:     ID,
		taskCh: make(chan Task),
	}
}

// Start 通道中取任务执行
func (w *Worker) Start() {
	go func() {
		defer fmt.Printf("ID:%d worker 关闭\n", w.ID)
		for task := range w.taskCh {
			task()
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

// NewPool 初始化协程池，duration为等待*秒后工作协程关闭，duration为0时如果没有新的任务工作协程会立刻关闭
func NewPool(num int, duration int) *Pool {
	p := &Pool{
		queue:    list.New(),
		duration: duration,
	}
	for i := 0; i < num; i++ {
		w := NewWorker(i)
		w.Start()
		p.workers = append(p.workers, w)
	}
	p.distribute()
	return p
}

// distribute 分配任务给worker
func (p *Pool) distribute() {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(w *Worker) {
			defer func() {
				close(w.taskCh)
				p.wg.Done()
			}()
			for {
				task := p.popTask()
				if task == nil {
					if p.duration > 0 {
						//等待*秒后退出
						if w.t > 0 && time.Now().Unix()-w.t >= int64(p.duration) {
							fmt.Printf("ID:%d worker 任务分配协程等待%ds关闭\n", w.ID, p.duration)
							return
						} else if w.t == 0 {
							w.t = time.Now().Unix()
						}
					} else {
						//直接退出
						return
					}
				} else {
					w.t = 0
					w.taskCh <- task
				}
			}
		}(w)
	}
}

// popTask 队列中取任务
func (p *Pool) popTask() Task {
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
	return task.Value.(Task)
}

// Submit 提交任务
func (p *Pool) Submit(task Task) {
	p.mx.Lock()
	defer p.mx.Unlock()
	p.queue.PushBack(task)
}

// Wait 等待任务执行完成
func (p *Pool) Wait() {
	p.wg.Wait()
}
