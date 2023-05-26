package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Status string

const (
	Eating   Status = "eating"
	Thinking Status = "thinking"
	Finished Status = "finished"
  Unseated Status = "unseated"
  NoStatus Status = ""
)

func GetStatusEmoji(s Status) (Status, error) {
	switch s {
	case Eating:
		{
      return "🍴", nil
		}
	case Thinking:
		{
			return "🤔", nil
		}
  case Finished:
    {
      return "❤️", nil
    }
  case Unseated:
    {
      return "😴", nil
    }
	}
  return NoStatus, fmt.Errorf("Invalid status")
}

type Update struct {
	philosopherId int
	apetite       int
	status        Status
}

func NewUpdate(id int, apetite int, status Status) *Update {
	return &Update{
		philosopherId: id,
		apetite:       apetite,
		status:        status,
	}
}

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
func (p Philosopher) Dine(wg *sync.WaitGroup, uc chan Update) {
	// signal to stop waiting on this goroutine when function is finished executing
	defer wg.Done()
	rand.Seed(int64(time.Now().Nanosecond()))

	// while the philosopher is hungry
	for p.apetite > 0 {

		// fmt.Printf("Philosopher %d is thinking\n", p.id)
		uc <- *NewUpdate(p.id, p.apetite, Thinking)

		// even philosophers reach for left fork then right
		if p.id%2 == 0 {
			p.leftFork.Lock()
			p.rightFork.Lock()
		} else {
			// odd philosophers reach for right fork then left
			p.rightFork.Lock()
			p.leftFork.Lock()
		}

		// fmt.Printf("Philosopher %d is eating\n", p.id)
		uc <- *NewUpdate(p.id, p.apetite, Eating)
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)

		// release hold of forks (order doesn't matter here)
		p.leftFork.Unlock()
		p.rightFork.Unlock()
		// decrement apetite
		p.apetite--
	}
	uc <- *NewUpdate(p.id, p.apetite, Finished)
}

func DineManager(philosophers []Philosopher, numPhilosophers int, updatesChannel chan Update) {
	defer close(updatesChannel)
	var wg sync.WaitGroup
	wg.Add(numPhilosophers)
	for i := 0; i < numPhilosophers; i++ {
		go philosophers[i].Dine(&wg, updatesChannel)
	}
	wg.Wait()
}

func RenderProgress(apetite int, startingApetite int) string {
  totalCols := 20
  prcnt := 1 - (float32(apetite) / float32(startingApetite))
  filledCols := int(float32(totalCols) * prcnt)
  result := fmt.Sprintf("[%s%s] %d%%", strings.Repeat("#", filledCols), strings.Repeat(" ", totalCols-filledCols), int(prcnt * float32(100)))
  return result
}

func RenderState(state map[int]Update, startingApetite int) {
	// clear screen
	fmt.Printf("\033[2J")
	// for _, p := range state {
  for i := 0; i < len(state); i++ {
    p := state[i]
		tabs := "\t"
		if p.status == Eating {
			tabs += "\t"
		}
    emoji, err := GetStatusEmoji(p.status)
    if err != nil {
      panic(err)
    }
		fmt.Printf("Philosopher:\t%d\t%s%s\t%s\tApetite: %d\t%s\n", p.philosopherId, p.status, tabs, emoji, p.apetite, RenderProgress(p.apetite, startingApetite))
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
  // set up philosopher update map
	philosopherState := map[int]Update{}
  for i, p := range philosophers {
    philosopherState[i] = *NewUpdate(p.id, p.apetite, Unseated)
  }

	// Should be plenty of space in the buffer for philosopher updates to flow in a queue
	updatesChannel := make(chan Update, numPhilosophers)
	// Start the dining manager
	go DineManager(philosophers, numPhilosophers, updatesChannel)
	for update := range updatesChannel {
		// fmt.Printf("Philosopher %d is %s. Apetite = %d\n", update.philosopherId, update.status, update.apetite)
		philosopherState[update.philosopherId] = update
		// display state
		RenderState(philosopherState, startingApetite)
	}
}
