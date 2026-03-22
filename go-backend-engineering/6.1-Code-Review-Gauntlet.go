// ============================================================
// CHALLENGE 6.1 — Code Review Gauntlet
// Rules:
// - Do NOT run the code. Read only.
// - For each program: list ALL bugs, explain why, write the fix.
// - Time yourself: target 2 hours total (24 min per program)
// ============================================================

// ============================================================
// PROGRAM 1 — The Goroutine Leak
// Topic: goroutines, channels
// Question: What is wrong with this program?
//           Will it ever exit cleanly? Why or why not?
// ============================================================

package main

import (
	"fmt"
	"time"
)

func produce(ch chan int) {
	for i := 0; i < 5; i++ {
		ch <- i
	}
}

func main() {
	ch := make(chan int)

	go produce(ch)

	for v := range ch {
		fmt.Println(v)
		time.Sleep(100 * time.Millisecond)
	}
}

// ============================================================
// PROGRAM 2 — The Race Condition
// Topic: goroutines, shared memory, sync.Mutex
// Question: This program sometimes prints wrong results.
//           Why? How do you fix it without changing the logic?
// ============================================================

package main

import (
	"fmt"
	"sync"
)

func main() {
	counter := 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter++
		}()
	}

	wg.Wait()
	fmt.Println("counter:", counter) // should print 1000, but doesn't always
}

// ============================================================
// PROGRAM 3 — The Nil Interface Trap
// Topic: interfaces, nil, type system
// Question: This program panics at runtime even though
//           the nil check passes. Why?
// ============================================================

package main

import "fmt"

type Logger interface {
	Log(msg string)
}

type FileLogger struct {
	filename string
}

func (f *FileLogger) Log(msg string) {
	fmt.Println("logging to file:", msg)
}

func getLogger(useFile bool) Logger {
	var fl *FileLogger // nil pointer

	if useFile {
		fl = &FileLogger{filename: "app.log"}
	}

	return fl // always returns fl, even when nil
}

func main() {
	logger := getLogger(false)

	if logger != nil {
		logger.Log("hello") // panics here — why?
	}
}

// ============================================================
// PROGRAM 4 — The Broken Interface Design
// Topic: interfaces, API design
// Question: This code compiles and works, but has a serious
//           design problem that will cause pain as the codebase
//           grows. What is wrong and how would you redesign it?
// ============================================================

package main

import "fmt"

type Animal interface {
	Breathe()
	Eat()
	Sleep()
	Swim()
	Fly()
	Run()
	MakeSound()
}

type Dog struct{}

func (d Dog) Breathe()   { fmt.Println("dog breathes") }
func (d Dog) Eat()       { fmt.Println("dog eats") }
func (d Dog) Sleep()     { fmt.Println("dog sleeps") }
func (d Dog) Swim()      { fmt.Println("dog swims") }
func (d Dog) Fly()       { fmt.Println("dog cannot fly") } // ← problem
func (d Dog) Run()       { fmt.Println("dog runs") }
func (d Dog) MakeSound() { fmt.Println("woof") }

type Bird struct{}

func (b Bird) Breathe()   { fmt.Println("bird breathes") }
func (b Bird) Eat()       { fmt.Println("bird eats") }
func (b Bird) Sleep()     { fmt.Println("bird sleeps") }
func (b Bird) Swim()      { fmt.Println("bird cannot swim") } // ← problem
func (b Bird) Fly()       { fmt.Println("bird flies") }
func (b Bird) Run()       { fmt.Println("bird runs") }
func (b Bird) MakeSound() { fmt.Println("tweet") }

func describeAnimal(a Animal) {
	a.Breathe()
	a.MakeSound()
}

func main() {
	describeAnimal(Dog{})
	describeAnimal(Bird{})
}

// ============================================================
// PROGRAM 5 — The Performance Anti-Pattern
// Topic: strings, memory, allocations
// Question: This function is correct but extremely slow
//           at scale. What is wrong and how do you fix it?
//           What Go tool would you use to confirm the problem?
// ============================================================

package main

import (
	"fmt"
	"strings"
)

// buildReport takes a list of log lines and builds a summary string
func buildReport(lines []string) string {
	report := ""

	for _, line := range lines {
		if strings.Contains(line, "ERROR") {
			report = report + "[ERROR] " + line + "\n"
		} else if strings.Contains(line, "WARN") {
			report = report + "[WARN]  " + line + "\n"
		} else {
			report = report + "[INFO]  " + line + "\n"
		}
	}

	return report
}

func main() {
	lines := []string{
		"user logged in",
		"ERROR: disk full",
		"WARN: high memory",
		"payment processed",
	}
	fmt.Println(buildReport(lines))
}