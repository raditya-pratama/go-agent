package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var goRoutine sync.WaitGroup

type Element struct {
	data string
	next *Element
}

type LinkedList struct {
	first *Element
	last  *Element
}

func (l *LinkedList) IsEmpty() bool {
	return l.first == nil && l.last == nil
}

func (l *LinkedList) Add(data string) {
	newEl := &Element{data: data}
	if l.IsEmpty() {
		l.first = newEl
	} else {
		l.last.next = newEl
	}
	l.last = newEl
}

func (l *LinkedList) Process() *Element {
	if l.IsEmpty() {
		return nil
	}
	processedEl := l.first
	l.first = l.first.next
	processedEl.next = nil

	if l.first == nil { // became empty
		l.last = nil
	}

	return processedEl
}

func (l *LinkedList) Print() {
	el := l.first
	for el != nil {
		fmt.Printf("%s; ", el.data)
		el = el.next
	}
	fmt.Println("")
}

func main() {
	rand.Seed(time.Now().Unix())

	// Create a buffered channel to manage the employee vs project load.
	projects := make(chan string, 10)

	// Launch 5 goroutines to handle the projects.
	goRoutine.Add(5)
	for i := 1; i <= 5; i++ {
		go employee(projects, i)
	}

	for j := 1; j <= 10; j++ {
		projects <- fmt.Sprintf("Project :%d", j)
	}

	// Close the channel so the goroutines will quit
	close(projects)
	fmt.Println("projects closed")
	goRoutine.Wait()

	fmt.Println("linkedlist sample start")
	list := &LinkedList{}
	list.Add("data1")
	list.Print()
	list.Add("data2")
	list.Print()
	list.Add("data3")
	list.Print()
	list.Add("data4")
	list.Print()

	for {
		if list := list.Process(); list == nil {
			break
		}
		list.Print()
	}
}

func employee(projects chan string, employee int) {
	defer goRoutine.Done()
	for {
		// Wait for project to be assigned.
		project, result := <-projects

		if result == false {
			// This means the channel is empty and closed.
			fmt.Printf("Employee : %d : Exit\n", employee)
			return
		}

		fmt.Printf("Employee : %d : Started   %s\n", employee, project)

		// Randomly wait to simulate work time.
		sleep := rand.Int63n(50)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		// Display time to wait
		fmt.Println("\nTime to sleep", sleep, "ms\n")

		// Display project completed by employee.
		fmt.Printf("Employee : %d : Completed %s\n", employee, project)
	}

}
