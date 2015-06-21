package tasks

import (
	"container/heap"
	"reflect"
	"runtime"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
)

// https://talks.golang.org/2010/io/balance.go
// https://talks.golang.org/2012/waza.slide#53
// http://habrahabr.ru/post/198150/

const (
	// количество исполнителей, запускаемых по умолчанию
	defaultWorkers = 1

	// статусы исполнителя
	workerStatusWait = iota
	workerStatusBusy
)

// статусы задачи
const (
	taskStatusWait = iota
	taskStatusProcess
	taskStatusSuccess
	taskStatusFail
	taskStatusRepeatWait
)

/*
 * Task
 */
type Task struct {
	taskID  string
	fn      func(...interface{}) (bool, time.Duration)
	args    []interface{}
	status  int
	created time.Time
}

// GetID возвращает уникальный идентификатор задачи
func (j *Task) GetID() string {
	return j.taskID
}

// GetName возвращает имя выполняемой задачи
func (t *Task) GetName() string {
	funcName := runtime.FuncForPC(reflect.ValueOf(t.fn).Pointer()).Name()

	//parts := strings.Split(funcName, "(")

	return funcName
}

// GetStatus возвращает статус задачи
func (t *Task) GetStatus() int {
	return t.status
}

// GetCreated возвращает дату создания задания
func (t *Task) GetCreated() time.Time {
	return t.created
}

/*
 * Worker
 */
type Worker struct {
	dispatcher     *Dispatcher
	index          int
	localWaitGroup *sync.WaitGroup
	newTask        chan *Task
	quit           chan bool // канал для завершения исполнителя

	workerID string
	task     *Task
	status   int
	created  time.Time
}

// GetID возвращает уникальный идентифкатор исполнителя
func (w *Worker) GetID() string {
	return w.workerID
}

// GetStatus возвращает статус исполнителя
func (w *Worker) GetStatus() int {
	return w.status
}

// GetTask возвращает текущее задание
func (w *Worker) GetTask() *Task {
	return w.task
}

// Kill завершает работу исполнителя
func (w *Worker) Kill() {
	w.quit <- true
}

// GetCreated возвращает дату создания исполнителя
func (w *Worker) GetCreated() time.Time {
	return w.created
}

// work выполняет задачу
func (w *Worker) work(done chan<- *Worker) {
	for {
		select {
		// пришло новое задание на выполнение
		case w.task = <-w.newTask:
			w.status = workerStatusBusy
			w.task.status = taskStatusProcess

			func() {
				w.dispatcher.waitGroup.Add(1)
				w.localWaitGroup.Add(1)

				defer func() {
					if err := recover(); err != nil {
						w.task.status = taskStatusFail
					} else {
						w.task.status = taskStatusSuccess
					}

					w.localWaitGroup.Done()
					w.dispatcher.waitGroup.Done()

					w.task = nil
				}()

				if repeat, duration := w.task.fn(w.task.args...); repeat {
					t := w.task
					t.status = taskStatusRepeatWait

					time.AfterFunc(duration, func() {
						w.dispatcher.sendTask(t)
					})
				}
			}()

			done <- w

		// пришел сигнал на завершение исполнителя
		case <-w.quit:
			// ждем завершения текущего задания, если таковое есть и выходим
			w.localWaitGroup.Wait()
			return
		}
	}
}

/*
 * Pool
 */
type Pool []*Worker

func (p Pool) Len() int {
	return len(p)
}

func (p Pool) Less(i, j int) bool {
	return p[i].status < p[j].status
}

func (p *Pool) Swap(i, j int) {
	a := *p

	if i >= 0 && i < len(a) && j >= 0 && j < len(a) {
		a[i], a[j] = a[j], a[i]
		a[i].index = i
		a[j].index = j
	}
}

func (p *Pool) Push(x interface{}) {
	n := len(*p)
	worker := x.(*Worker)
	worker.index = n
	*p = append(*p, worker)
}

func (p *Pool) Pop() interface{} {
	a := *p
	n := len(a)

	item := a[n-1]
	item.index = -1

	*p = a[0 : n-1]

	return item
}

/*
 * Dispatcher
 */
type Dispatcher struct {
	newTasks        chan *Task   // очередь новых заданий
	queue           chan *Task   // очередь выполняемых заданий
	workers         Pool         // пул исполнителей
	done            chan *Worker // канал уведомления о завершении выполнения заданий
	allowProcessing chan bool    // канал для блокировки выполнения новых задач для случая, когда все исполнители заняты
	quit            chan bool    // канал для завершения диспетчера
	workersBusy     int          // количество занятых исполнителей
	tasksWait       []*Task      // задачи, ожидающие назначения исполнителя

	waitGroup *sync.WaitGroup
}

func NewDispatcher() *Dispatcher {
	d := &Dispatcher{
		newTasks:        make(chan *Task),
		queue:           make(chan *Task),
		workers:         make(Pool, 0),
		done:            make(chan *Worker),
		allowProcessing: make(chan bool),
		quit:            make(chan bool),
		waitGroup:       new(sync.WaitGroup),
		workersBusy:     0,
		tasksWait:       make([]*Task, 0),
	}

	// отслеживание квоты на занятость исполнителей
	go func() {
		for {
			d.queue <- <-d.newTasks
			d.tasksWait = append(d.tasksWait[1:])

			<-d.allowProcessing
		}
	}()

	// инициализируем исполнителей
	heap.Init(&d.workers)
	for i := 0; i < defaultWorkers; i++ {
		d.AddWorker()
	}

	go d.process()
	return d
}

// AddWorker добавляет еще одного исполнителя в пулл
func (d *Dispatcher) AddWorker() {
	w := &Worker{
		dispatcher:     d,
		localWaitGroup: new(sync.WaitGroup),
		newTask:        make(chan *Task),
		quit:           make(chan bool),
		workerID:       uuid.New(),
		status:         workerStatusWait,
		created:        time.Now(),
	}

	heap.Push(&d.workers, w)
	go w.work(d.done)
}

// AddTask добавляет задание в очередь на выполнение и возвращает саму задачу
func (d *Dispatcher) AddTask(fn func(...interface{}) (bool, time.Duration), args ...interface{}) {
	t := &Task{
		taskID:  uuid.New(),
		fn:      fn,
		args:    args,
		status:  taskStatusWait,
		created: time.Now(),
	}

	d.sendTask(t)
}

func (d *Dispatcher) sendTask(t *Task) {
	go func() {
		d.tasksWait = append(d.tasksWait, t)
		d.newTasks <- t
	}()
}

// Kill завершает работы диспетчера
func (d *Dispatcher) Kill() {
	d.quit <- true
}

// GetStats возвращает статистику
func (d *Dispatcher) GetStats() map[string]interface{} {
	workers := make([]map[string]interface{}, 0, d.workers.Len())
	for i := range d.workers {
		worker := d.workers[i]
		stat := map[string]interface{}{
			"id":      worker.GetID(),
			"status":  worker.GetStatus(),
			"created": worker.GetCreated(),
		}

		if worker.task != nil {
			stat["task"] = map[string]interface{}{
				"id":      worker.task.GetID(),
				"name":    worker.task.GetName(),
				"status":  worker.task.GetStatus(),
				"created": worker.task.GetCreated(),
			}
		} else {
			stat["task"] = nil
		}

		workers = append(workers, stat)
	}

	tasksWait := make([]map[string]interface{}, 0, len(d.tasksWait))
	for i := range d.tasksWait {
		task := d.tasksWait[i]
		stat := map[string]interface{}{
			"id":      task.GetID(),
			"name":    task.GetName(),
			"status":  task.GetStatus(),
			"created": task.GetCreated(),
		}

		tasksWait = append(tasksWait, stat)
	}

	stats := map[string]interface{}{
		"workers_count": d.workers.Len(),
		"workers_busy":  d.workersBusy,
		"workers_wait":  d.workers.Len() - d.workersBusy,
		"workers":       workers,
		"tasks_wait":    tasksWait,
	}

	return stats
}

// dispatch отправляет задание свободному исполнителю
func (d *Dispatcher) dispatch(t *Task) {
	worker := heap.Pop(&d.workers).(*Worker)
	worker.newTask <- t
	heap.Push(&d.workers, worker)

	// проверяем есть ли еще свободные исполнители для задач
	if d.workersBusy++; d.workersBusy < d.workers.Len() {
		d.allowProcessing <- true
	}
}

func (d *Dispatcher) completed(w *Worker) {
	heap.Remove(&d.workers, w.index)
	w.status = workerStatusWait
	heap.Push(&d.workers, w)

	// проверяем не освободился ли какой-нибудь исполнитель
	if d.workersBusy--; d.workersBusy == d.workers.Len()-1 {
		d.allowProcessing <- true
	}
}

func (d *Dispatcher) process() {
	for {
		select {
		// пришел новый таск на выполнение от flow контроллера
		case task := <-d.queue:
			d.dispatch(task)

		// пришло уведомление, что рабочий закончил выполнение задачи
		case worker := <-d.done:
			d.completed(worker)

		// завершение работы диспетчера
		case <-d.quit:
			// ждем завершения всех заданий и всех исполнителей
			d.waitGroup.Wait()
			return
		}
	}
}
