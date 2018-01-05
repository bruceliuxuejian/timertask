package timertask

import (
	"fmt"
)

type TaskTest struct {
	Num int
}

func (u *TaskTest) Doing() {
	fmt.Println(u.Num)
}

func main() {

	myTask := NewTaskManager(1)
	u := &TaskTest{Num: 1}
	u1 := &TaskTest{Num: 1000}
	u2 := &TaskTest{Num: 10000}
	/**
	RecycleTask
	*/
	myTask.RegistRecycleTask(2, u)
	myTask.RegistRecycleTask(3, u1)
	myTask.RegistRecycleTask(4, u2)
	/**
	OnceTask
	*/
	myTask.RegistTask(1515058675, &TaskTest{Num: 0})

}
