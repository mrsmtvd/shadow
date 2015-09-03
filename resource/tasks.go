package resource

import (
	"container/heap"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/kihamo/shadow"
	"github.com/nu7hatch/gouuid"
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
type task struct {
	id       string
	name     string
	status   int
	created  time.Time
	duration time.Duration
	fn       func(...interface{}) (bool, time.Duration)
	args     []interface{}
}

/*
 * Worker
 */
type worker struct {
	id          string
	status      int
	index       int
	created     time.Time
	waitGroup   *sync.WaitGroup
	newTask     chan *task
	quit        chan bool // канал для завершения исполнителя
	executeTask *task
	logger      *logrus.Entry
}

// kill worker shutdown
func (w *worker) kill() {
	w.quit <- true
}

// work выполняет задачу
func (w *worker) work(done chan<- *worker, repeat chan<- *task) {
	for {
		select {
		// пришло новое задание на выполнение
		case w.executeTask = <-w.newTask:
			w.executeTask.status = taskStatusProcess

			func() {
				w.waitGroup.Add(1)

				defer func() {
					if err := recover(); err != nil {
						w.logger.WithFields(logrus.Fields{
							"task":  w.executeTask.name,
							"args":  w.executeTask.args,
							"error": err,
						}).Warn("Failed")

						w.executeTask.status = taskStatusFail
					} else {
						w.logger.WithFields(logrus.Fields{
							"task": w.executeTask.name,
							"args": w.executeTask.args,
						}).Debug("Success")
						w.executeTask.status = taskStatusSuccess
					}

					w.executeTask = nil
					w.waitGroup.Done()

					done <- w
				}()

				var repeated bool
				if repeated, w.executeTask.duration = w.executeTask.fn(w.executeTask.args...); repeated {
					repeat <- w.executeTask
				}
			}()

		// пришел сигнал на завершение исполнителя
		case <-w.quit:
			// ждем завершения текущего задания, если таковое есть и выходим
			w.waitGroup.Wait()
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
	workers pool // пул исполнителей

	workersBusy int     // количество занятых исполнителей
	tasksWait   []*task // задачи, ожидающие назначения исполнителя

	waitGroup   *sync.WaitGroup
	mutex       sync.RWMutex
	application *shadow.Application
	config      *Config
	logger      *logrus.Entry

	waitQueue    chan *task // очередь новых заданий
	executeQueue chan *task // очередь выполняемых заданий
	repeatQueue  chan *task // канал уведомления о повторном выполнении заданий

	done chan *worker // канал уведомления о завершении выполнения заданий

	quit            chan bool // канал для завершения диспетчера
	allowProcessing chan bool // канал для блокировки выполнения новых задач для случая, когда все исполнители заняты
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

	d.workers = make(pool, 0)
	d.waitGroup = new(sync.WaitGroup)
	d.workersBusy = 0
	d.tasksWait = make([]*task, 0)
	d.waitQueue = make(chan *task)
	d.executeQueue = make(chan *task)
	d.repeatQueue = make(chan *task)
	d.done = make(chan *worker)
	d.quit = make(chan bool)
	d.allowProcessing = make(chan bool)

	return nil
}

func (d *Dispatcher) Run(wg *sync.WaitGroup) error {
	resourceLogger, err := d.application.GetResource("logger")
	if err != nil {
		return err
	}
	d.logger = resourceLogger.(*Logger).Get(d.GetName())

	// отслеживание квоты на занятость исполнителей
	go func() {
		for {
			d.executeQueue <- <-d.waitQueue

			d.mutex.Lock()
			d.tasksWait = append(d.tasksWait[1:])
			d.mutex.Unlock()

			<-d.allowProcessing
		}
	}()

	// инициализируем исполнителей
	heap.Init(&d.workers)

	maxWorkers := int(d.config.GetInt64("tasks.workers"))
	for len(d.workers) < maxWorkers {
		d.AddWorker()
	}

	go func() {
		defer wg.Done()

		for {
			select {
			// пришел новый таск на выполнение от flow контроллера
			case task := <-d.executeQueue:
				d.dispatch(task)

			// пришло уведомление, что рабочий закончил выполнение задачи
			case worker := <-d.done:
				d.completed(worker)

			// пришло уведомление, что необходимо повторить задачу
			case task := <-d.repeatQueue:
				task.status = taskStatusRepeatWait

				d.logger.WithFields(logrus.Fields{
					"task": task.name,
					"args": task.args,
				}).Debug("Repeat")

				d.sendTask(task)

			// завершение работы диспетчера
			case <-d.quit:
				// ждем завершения всех заданий и всех исполнителей
				d.waitGroup.Wait()
				return
			}
		}
	}()

	return nil
}

// AddWorker добавляет еще одного исполнителя в пулл
func (d *Dispatcher) AddWorker() {
	id, _ := uuid.NewV4()

	w := &worker{
		id:        id.String(),
		status:    workerStatusWait,
		created:   time.Now(),
		waitGroup: d.waitGroup,
		newTask:   make(chan *task, 1),
		quit:      make(chan bool),
		logger:    d.logger.WithField("worker", id),
	}

	go w.work(d.done, d.repeatQueue)
	heap.Push(&d.workers, w)
}

func (d *Dispatcher) AddNamedTask(name string, fn func(...interface{}) (bool, time.Duration), args ...interface{}) {
	id, _ := uuid.NewV4()

	t := &task{
		id:      id.String(),
		name:    name,
		status:  taskStatusWait,
		created: time.Now(),
		fn:      fn,
		args:    args,
	}

	d.sendTask(t)
}

func (d *Dispatcher) AddTask(fn func(...interface{}) (bool, time.Duration), args ...interface{}) {
	d.AddNamedTask(runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name(), fn, args...)
}

func (d *Dispatcher) sendTask(t *task) {
	time.AfterFunc(t.duration, func() {
		go func() {
			d.mutex.Lock()
			d.tasksWait = append(d.tasksWait, t)
			d.mutex.Unlock()

			d.waitQueue <- t
		}()
	})
}

// Kill dispatcher shutdown
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
			"id":      worker.id,
			"status":  worker.status,
			"created": worker.created,
		}

		if worker.executeTask != nil {
			stat["task"] = map[string]interface{}{
				"id":      worker.executeTask.id,
				"name":    worker.executeTask.name,
				"status":  worker.executeTask.status,
				"created": worker.executeTask.created,
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
			"id":      task.id,
			"name":    task.name,
			"status":  task.status,
			"created": task.created,
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
func (d *Dispatcher) dispatch(t *task) {
	worker := heap.Pop(&d.workers).(*worker)
	worker.newTask <- t
	worker.status = workerStatusBusy
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
