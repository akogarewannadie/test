package main

import (
	"fmt"
	"time"
)

type Ttype struct {
	id         int
	cT         string // время создания
	fT         string // время выполнения
	taskRESULT []byte
}

func main() {
	superChan := make(chan Ttype, 10)
	doneTasks := make(chan Ttype)
	undoneTasks := make(chan Ttype)

	taskCreator := func() {
		for {
			ft := time.Now().Format(time.RFC3339)
			if time.Now().Nanosecond()%2 > 0 { // условие появления ошибочных заданий
				ft = "Some error occurred"
			}
			superChan <- Ttype{cT: ft, id: int(time.Now().Unix())}
			time.Sleep(time.Second) // пауза между созданием заданий
		}
	}

	taskWorker := func() {
		for t := range superChan {
			tt, _ := time.Parse(time.RFC3339, t.cT)
			if tt.After(time.Now().Add(-20 * time.Second)) {
				t.taskRESULT = []byte("task has been succeeded")
			} else {
				t.taskRESULT = []byte("something went wrong")
			}
			t.fT = time.Now().Format(time.RFC3339Nano)
			time.Sleep(time.Millisecond * 150) // имитация времени выполнения задания
			if string(t.taskRESULT) == "task has been succeeded" {
				doneTasks <- t
			} else {
				undoneTasks <- t
			}
		}
	}

	printDoneTasks := func() {
		for t := range doneTasks {
			fmt.Printf("Done task: ID %d, Time %s\n", t.id, t.fT)
		}
	}

	printUndoneTasks := func() {
		for t := range undoneTasks {
			fmt.Printf("Undone task: ID %d, Time %s, Error: %s\n", t.id, t.cT, t.taskRESULT)
		}
	}

	go taskCreator()
	go taskWorker()
	go printDoneTasks()
	go printUndoneTasks()

	time.Sleep(time.Minute)
}
