package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type Philosopher struct {
	leftFork  *sync.Mutex
	rightFork *sync.Mutex
	id        int
	apetite   int
}

func NewPhilosopher(id int, apetite int, leftFork *sync.Mutex, rightFork *sync.Mutex) *Philosopher {
	return &Philosopher{
		id:        id,
		apetite:   apetite,
		leftFork:  leftFork,
		rightFork: rightFork,
	}
}

func getLeftFork(id int, philosophersCount int) int {
	result := id - 1

	if result < 0 {
		result = philosophersCount - 1
	}

	return result
}

/*
The selected philosopher dines (should be ran as go routine)
*/
func (p Philosopher) Dine(wg *sync.WaitGroup) {
	// signal to stop waiting on this goroutine when function is finished executing
	defer wg.Done()
  rand.Seed(int64(time.Now().Nanosecond()))

	// while the philosopher is hungry
	for p.apetite > 0 {

		fmt.Printf("Philosopher %d is thinking\n", p.id)

		// even philosophers reach for left fork then right
		if p.id % 2 == 0 {
			p.leftFork.Lock()
			p.rightFork.Lock()
		} else {
			// odd philosophers reach for right fork then left
			p.rightFork.Lock()
			p.leftFork.Lock()
		}

		fmt.Printf("Philosopher %d is eating\n", p.id)
		time.Sleep(time.Duration(rand.Intn(80)) * time.Millisecond)

		// release hold of forks (order doesn't matter here)
		p.leftFork.Unlock()
		p.rightFork.Unlock()
		// decrement apetite
		p.apetite--
	}
}

func main() {
	args := os.Args

	// verify that user provided required arguments
	if len(args) != 3 {
		fmt.Println("Missing arguments")
		fmt.Printf("Usage: %s {number of philosophers} {philosopher apetite}\n", args[0])
		os.Exit(1)
	}

	numPhilosophers, err := strconv.Atoi(args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	startingApetite, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		os.Exit(1)
	}

	// create forks. Should be one to every philosopher
	var forks []sync.Mutex

	for i := 0; i < numPhilosophers; i++ {
		forks = append(forks, sync.Mutex{})
	}

	// create philosophers
	var philosophers []Philosopher

	for i := 0; i < numPhilosophers; i++ {
		philosophers = append(philosophers, *NewPhilosopher(i, startingApetite, &forks[getLeftFork(i, numPhilosophers)], &forks[i]))
	}

	var wg sync.WaitGroup
	wg.Add(numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		go philosophers[i].Dine(&wg)
	}
	wg.Wait()
}
