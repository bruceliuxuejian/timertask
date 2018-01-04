package task

import (
	"fmt"
	"sync/atomic"
	"time"
)

type TaskPayload interface {
	Before()
	Doing()
	Finish()
}

type Task struct {
	id          int64
	executeTime int64
	Payload     TaskPayload
	recycleTime int64
	recycleNum  int64
}

type TaskManager struct {
	taskID                int64
	executingTimestamp    int64
	lastExecutedTimestamp int64
	taskMap               map[int64]*Task
	taskQueue             map[int64][]int64
	taskChan              chan *Task
	taskCancel            chan int64
	ticker                *time.Ticker
	cancelMap             map[int64]struct{}
}

func NewTask(timeInterval time.Duration) *TaskManager {
	tmranager := new(TaskManager)
	if timeInterval <= 0 {
		timeInterval = 1
	}
	tmranager.loding(timeInterval)
	return tmranager
}

func (manager *TaskManager) loding(timeInterval time.Duration) {
	manager.taskMap = make(map[int64]*Task)
	manager.taskChan = make(chan *Task, 100)
	manager.taskQueue = make(map[int64][]int64)
	manager.cancelMap = make(map[int64]struct{})
	manager.taskCancel = make(chan int64, 100)

	manager.ticker = time.NewTicker(timeInterval * time.Second)
	manager.lastExecutedTimestamp = getTimestamp()
	go manager.taskRoutine()
}
func (manager *TaskManager) RegistTask(executeTime int64, payload TaskPayload) int64 {
	newId := atomic.AddInt64(&manager.taskID, 1)
	task := &Task{id: newId, executeTime: executeTime, Payload: payload}
	manager.taskChan <- task
	return newId
}
func (manager *TaskManager) RegistRecycleTask(recycleTime int64, payload TaskPayload) int64 {
	newId := atomic.AddInt64(&manager.taskID, 1)
	task := &Task{id: newId, recycleTime: recycleTime, Payload: payload, executeTime: getTimestamp() + recycleTime, recycleNum: 1}
	fmt.Println("executeTime", task.executeTime)
	manager.taskChan <- task
	return newId
}
func (manager *TaskManager) Cancel(taskId int64) {
	if taskId <= 0 {
		return
	}
	manager.taskCancel <- taskId
}

func (manager *TaskManager) taskRoutine() {
	for {
		select {
		case <-manager.ticker.C:
			manager.runTasks()
		case task, ok := <-manager.taskChan:
			if ok {
				manager.registTask(task)
			}
		case taskId, ok := <-manager.taskCancel:
			if ok {
				manager.cancelTask(taskId)
			}
		}
	}
}

func (manager *TaskManager) runTasks() {
	nowTime := getTimestamp()
	manager.executingTimestamp = nowTime
	for sec := manager.lastExecutedTimestamp + 1; sec <= nowTime; sec++ {
		tasks, ok := manager.taskQueue[sec]
		if !ok {
			continue
		}

		for _, taskID := range tasks {

			// Check whether the task is cancelled
			if _, ok := manager.cancelMap[taskID]; ok {
				delete(manager.cancelMap, taskID)
				continue
			}
			task, ok := manager.taskMap[taskID]
			if !ok {
				continue
			}
			manager.beforeTask(task)
			manager.doTask(task)
			manager.finishTask(task)
			if task.recycleTime <= 0 {
				delete(manager.taskMap, task.id)
			} else {
				time := task.executeTime + task.recycleNum*task.recycleTime
				task.recycleNum++
				_, ok := manager.taskQueue[time]
				if !ok {
					manager.taskQueue[time] = append(make([]int64, 0, 100), task.id)
				} else {
					manager.taskQueue[time] = append(manager.taskQueue[time], task.id)
				}
			}

		}

		delete(manager.taskQueue, sec)
	}
	manager.lastExecutedTimestamp = nowTime
}
func (manager *TaskManager) registTask(task *Task) {
	nowTime := getTimestamp()
	//The overdue task is carried out immediately
	if task.executeTime-1 <= nowTime { // Keep a second buffer
		fmt.Println("obsolete task", task.id)
		manager.doTask(task)
		return
	}

	manager.taskMap[task.id] = task
	tasks, ok := manager.taskQueue[task.executeTime]
	if !ok {
		tasks = make([]int64, 0, 100)
		manager.taskQueue[task.executeTime] = tasks
	}
	manager.taskQueue[task.executeTime] = append(tasks, task.id)
}
func (manager *TaskManager) cancelTask(taskId int64) {
	manager.cancelMap[taskId] = struct{}{}
	delete(manager.taskMap, taskId)
}
func (manager *TaskManager) doTask(task *Task) {
	task.Payload.Doing()
}
func (manager *TaskManager) beforeTask(task *Task) {
	task.Payload.Before()
}
func (manager *TaskManager) finishTask(task *Task) {
	task.Payload.Finish()
}
func getTimestamp() int64 {
	return time.Now().Unix()
}
