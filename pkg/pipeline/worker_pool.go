package pipeline

import (
	"context"
	"sync"
	"sync/atomic"
)

// WorkerPool 工作协程池
type WorkerPool struct {
	size    int32         // 池大小
	tasks   chan func()   // 任务队列
	workers []*Worker     // 工作协程
	active  int32         // 活跃数
	stopped int32         // 停止标志
	wg      sync.WaitGroup // 等待组
}

// Worker 工作协程
type Worker struct {
	id     int32
	pool   *WorkerPool
	ctx    context.Context
	cancel context.CancelFunc
}

// NewWorkerPool 创建新的工作协程池
func NewWorkerPool(size int) *WorkerPool {
	if size <= 0 {
		size = 1
	}
	
	pool := &WorkerPool{
		size:  int32(size),
		tasks: make(chan func(), 1000), // 默认缓冲区大小
	}
	
	// 启动工作协程
	for i := int32(0); i < pool.size; i++ {
		worker := &Worker{
			id:   i,
			pool: pool,
		}
		worker.ctx, worker.cancel = context.WithCancel(context.Background())
		pool.workers = append(pool.workers, worker)
		
		pool.wg.Add(1)
		go worker.run()
	}
	
	return pool
}

// Submit 提交任务
func (p *WorkerPool) Submit(task func()) error {
	if atomic.LoadInt32(&p.stopped) == 1 {
		return nil // 池已停止，不接受新任务
	}
	
	select {
	case p.tasks <- task:
		return nil
	default:
		// 任务队列已满，直接在调用协程中执行
		go task()
		return nil
	}
}

// Stop 停止工作池
func (p *WorkerPool) Stop() {
	if !atomic.CompareAndSwapInt32(&p.stopped, 0, 1) {
		return // 已经停止
	}
	
	// 关闭任务通道
	close(p.tasks)
	
	// 取消所有工作协程
	for _, worker := range p.workers {
		worker.cancel()
	}
	
	// 等待所有工作协程完成
	p.wg.Wait()
}

// ActiveCount 获取活跃协程数
func (p *WorkerPool) ActiveCount() int32 {
	return atomic.LoadInt32(&p.active)
}

// Size 获取池大小
func (p *WorkerPool) Size() int32 {
	return atomic.LoadInt32(&p.size)
}

// run 工作协程运行函数
func (w *Worker) run() {
	defer w.pool.wg.Done()
	
	atomic.AddInt32(&w.pool.active, 1)
	defer atomic.AddInt32(&w.pool.active, -1)
	
	for {
		select {
		case <-w.ctx.Done():
			return
		case task, ok := <-w.pool.tasks:
			if !ok {
				return // 任务通道已关闭
			}
			
			// 执行任务
			task()
		}
	}
}
