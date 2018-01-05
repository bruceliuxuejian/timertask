package timertask

import (
	"fmt"
)

type Task struct {
	Num int
}

func (u *Task) Doing() {
	fmt.Println(u.Num)
}

func main() {

	myTask := NewTaskManager(1)
	u := &Task{Num: 1}
	u1 := &Task{Num: 1000}
	u2 := &Task{Num: 10000}
	/**
	RecycleTask
	*/
	myTask.RegistRecycleTask(2, u)
	myTask.RegistRecycleTask(3, u1)
	myTask.RegistRecycleTask(4, u2)
	/**
	OnceTask
	*/
	myTask.RegistTask(1515058675, &Task{Num: 0})

}
