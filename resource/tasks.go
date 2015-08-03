package resource

import (
	"container/heap"
	"reflect"
	"runtime"
	"sync"
	"time"

	"code.google.com/p/go-uuid/uuid"
	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
)

// https://talks.golang.org/2010/io/balance.go
// https://talks.golang.org/2012/waza.slide#53
// http://habrahabr.ru/post/198150/

const (
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
	name    string
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
	return t.name
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
type worker struct {
	dispatcher     *Dispatcher
	index          int
	localWaitGroup *sync.WaitGroup
	newTask        chan *Task
	quit           chan bool // канал для завершения исполнителя

	workerID string
	task     *Task
	status   int
	created  time.Time

	logger *logrus.Entry
}

// GetID возвращает уникальный идентифкатор исполнителя
func (w *worker) GetID() string {
	return w.workerID
}

// GetStatus возвращает статус исполнителя
func (w *worker) GetStatus() int {
	return w.status
}

// GetTask возвращает текущее задание
func (w *worker) GetTask() *Task {
	return w.task
}

// Kill завершает работу исполнителя
func (w *worker) Kill() {
	w.quit <- true
}

// GetCreated возвращает дату создания исполнителя
func (w *worker) GetCreated() time.Time {
	return w.created
}

// work выполняет задачу
func (w *worker) work(done chan<- *worker) {
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
						w.logger.WithFields(logrus.Fields{
							"task":  w.task.GetName(),
							"args":  w.task.args,
							"error": err,
						}).Warn("Failed")

						w.task.status = taskStatusFail
					} else {
						w.logger.WithFields(logrus.Fields{
							"task": w.task.GetName(),
							"args": w.task.args,
						}).Debug("Success")
						w.task.status = taskStatusSuccess
					}

					w.localWaitGroup.Done()
					w.dispatcher.waitGroup.Done()

					w.task = nil
				}()

				if repeat, duration := w.task.fn(w.task.args...); repeat {
					t := w.task
					t.status = taskStatusRepeatWait
					w.logger.WithFields(logrus.Fields{
						"task": w.task.GetName(),
						"args": w.task.args,
					}).Debug("Repeat")

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
type pool []*worker

func (p pool) Len() int {
	return len(p)
}

func (p pool) Less(i, j int) bool {
	return p[i].status < p[j].status
}

func (p *pool) Swap(i, j int) {
	a := *p

	if i >= 0 && i < len(a) && j >= 0 && j < len(a) {
		a[i], a[j] = a[j], a[i]
		a[i].index = i
		a[j].index = j
	}
}

func (p *pool) Push(x interface{}) {
	n := len(*p)
	worker := x.(*worker)
	worker.index = n
	*p = append(*p, worker)
}

func (p *pool) Pop() interface{} {
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
	workers         pool         // пул исполнителей
	done            chan *worker // канал уведомления о завершении выполнения заданий
	allowProcessing chan bool    // канал для блокировки выполнения новых задач для случая, когда все исполнители заняты
	quit            chan bool    // канал для завершения диспетчера
	workersBusy     int          // количество занятых исполнителей
	tasksWait       []*Task      // задачи, ожидающие назначения исполнителя

	waitGroup   *sync.WaitGroup
	mutex       sync.RWMutex
	application *shadow.Application
	config      *Config
	logger      *logrus.Entry
}

func (d *Dispatcher) GetName() string {
	return "tasks"
}

func (d *Dispatcher) GetConfigVariables() []ConfigVariable {
	return []ConfigVariable{
		ConfigVariable{
			Key:   "tasks.workers",
			Value: 2,
			Usage: "Default workers count",
		},
	}
}

func (d *Dispatcher) Init(a *shadow.Application) error {
	d.application = a

	resourceConfig, err := a.GetResource("config")
	if err != nil {
		return err
	}
	d.config = resourceConfig.(*Config)

	d.newTasks = make(chan *Task)
	d.queue = make(chan *Task)
	d.workers = make(pool, 0)
	d.done = make(chan *worker)
	d.allowProcessing = make(chan bool)
	d.quit = make(chan bool)
	d.waitGroup = new(sync.WaitGroup)
	d.workersBusy = 0
	d.tasksWait = make([]*Task, 0)

	return nil
}

func (d *Dispatcher) Run() error {
	resourceLogger, err := d.application.GetResource("logger")
	if err != nil {
		return err
	}
	d.logger = resourceLogger.(*Logger).Get(d.GetName())

	// отслеживание квоты на занятость исполнителей
	go func() {
		for {
			d.queue <- <-d.newTasks

			d.mutex.Lock()
			d.tasksWait = append(d.tasksWait[1:])
			d.mutex.Unlock()

			<-d.allowProcessing
		}
	}()

	// инициализируем исполнителей
	heap.Init(&d.workers)

	var i int64
	for i = 0; i < d.config.GetInt64("tasks.workers"); i++ {
		d.AddWorker()
	}

	go d.process()
	return nil
}

// AddWorker добавляет еще одного исполнителя в пулл
func (d *Dispatcher) AddWorker() {
	id := uuid.New()

	w := &worker{
		dispatcher:     d,
		localWaitGroup: new(sync.WaitGroup),
		newTask:        make(chan *Task),
		quit:           make(chan bool),
		workerID:       id,
		status:         workerStatusWait,
		created:        time.Now(),
		logger:         d.logger.WithField("worker", id),
	}

	heap.Push(&d.workers, w)
	go w.work(d.done)
}

// AddTask добавляет задание в очередь на выполнение и возвращает саму задачу
func (d *Dispatcher) AddNamedTask(name string, fn func(...interface{}) (bool, time.Duration), args ...interface{}) {
	t := &Task{
		taskID:  uuid.New(),
		name:    name,
		fn:      fn,
		args:    args,
		status:  taskStatusWait,
		created: time.Now(),
	}

	d.sendTask(t)
}

func (d *Dispatcher) AddTask(fn func(...interface{}) (bool, time.Duration), args ...interface{}) {
	d.AddNamedTask(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), fn, args...)
}

func (d *Dispatcher) sendTask(t *Task) {
	go func() {
		d.mutex.Lock()
		d.tasksWait = append(d.tasksWait, t)
		d.mutex.Unlock()

		d.newTasks <- t
	}()
}

// Kill завершает работы диспетчера
func (d *Dispatcher) Kill() {
	d.quit <- true
}

// GetStats возвращает статистику
func (d *Dispatcher) GetStats() map[string]interface{} {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

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
	worker := heap.Pop(&d.workers).(*worker)
	worker.newTask <- t
	heap.Push(&d.workers, worker)

	// проверяем есть ли еще свободные исполнители для задач
	if d.workersBusy++; d.workersBusy < d.workers.Len() {
		d.allowProcessing <- true
	}
}

func (d *Dispatcher) completed(w *worker) {
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
