package main

import (
	"fmt"
	"task"
)

type User struct {
	Num int
}

func (u *User) Before() {
	u.Num++
}
func (u *User) Doing() {
	fmt.Println(u.Num)
}
func (u *User) Finish() {

}

func main() {

	myTask := task.NewTask(1)
	u := &User{Num: 1}
	u1 := &User{Num: 1000}
	u2 := &User{Num: 10000}

	myTask.RegistRecycleTask(2, u)
	myTask.RegistRecycleTask(3, u1)
	myTask.RegistRecycleTask(4, u2)

	for {

	}
}
