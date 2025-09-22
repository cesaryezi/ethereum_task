package task2

import "testing"

func TestPointerUse(t *testing.T) {
	num := 12
	pointerUse(&num)
	t.Log(num)
}

func TestSliceMul(t *testing.T) {
	num := []int{1, 2, 3, 4, 5}
	sliceMul(&num)
	t.Log(num)
}

func TestGoroutine(t *testing.T) {
	goroutine()
}

func TestTaskScheduler(t *testing.T) {
	tasks := []func(){
		func() {
			t.Log("task1")
		},
		func() {
			t.Log("task2")
		},
		func() {
			t.Log("task3")
		},
	}
	taskScheduler(tasks)
}

func TestObjectOriented(t *testing.T) {
	objectOriented()
}

func TestPrintInfo(t *testing.T) {
	e := Employee{
		EmployeeID: 1,
		Person: Person{
			Age:  12,
			Name: "张三",
		},
	}

	e.PrintInfo()
}

func TestChannel(t *testing.T) {
	channel()
}

func TestBufferedChannel(t *testing.T) {
	bufferedChannel()
}

func TestMutex(t *testing.T) {
	res := mutex()
	t.Log(res)
}

func TestAtomic(t *testing.T) {
	res := atomic()
	t.Log(res)
}
