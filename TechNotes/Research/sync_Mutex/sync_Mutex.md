

# Gemini4 Hydrogen Question No: 1 ***sync.Mutex*** - 2026-11-12 22:55

Using the 'go' programming language (a.k.a. Golang), please create sample code showing the best procedures for using a sync.Mutex to prevent race conditions. This example code should create a simple 'struct' containing a method with a receiver pointer. This method needs to lock the mutex before performing an operation and release the lock afterwards.



# Gemini4 Hydrogen Response No: 1 ***sync.Mutex*** - 2026-11-12 22:55



Here is a complete, runnable ***Go*** example demonstrating the idiomatic best practices for using `sync.Mutex` to prevent data races.

```go
package main

import (
	"fmt"
	"sync"
)

// SafeCounter is safe to use concurrently across multiple goroutines.
type SafeCounter struct {
	mu    sync.Mutex // Protects the count field below. Zero value is ready to use.
	count int
}

// Increment safely increases the counter by 1.
// A pointer receiver (*SafeCounter) is required so the struct (and its mutex) is not copied.
func (c *SafeCounter) Increment() {
	c.mu.Lock()
	defer c.mu.Unlock() // Guarantees the lock is released when the method returns.

	c.count++
}

// Value safely reads and returns the current count.
// Reads must also acquire the lock to prevent a data race with concurrent writes.
func (c *SafeCounter) Value() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.count
}

func main() {
	// No explicit constructor is needed because the zero value of sync.Mutex is unlocked.
	var counter SafeCounter
	var wg sync.WaitGroup

	const workers = 1000
	wg.Add(workers)

	// Launch 1000 concurrent goroutines that all modify the same SafeCounter.
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			counter.Increment()
		}()
	}

	wg.Wait()
	fmt.Printf("Final count: %d\n", counter.Value())
}
```

### Best Practices Illustrated

1. **Always use a pointer receiver (`*SafeCounter`):**
   Go passes value receivers by copy. Copying a `sync.Mutex` duplicates its internal lock state, completely breaking synchronization and potentially causing deadlocks (the `go vet` tool will even flag this as an error). Once a struct containing a mutex is created, it must never be copied.
2. **Use `defer c.mu.Unlock()` immediately after `c.mu.Lock()`:**
   Placing `defer` directly below `Lock()` ensures the mutex is always released—even if the method grows more complex, returns early from an `if` branch, or panics.
3. **Keep the mutex unexported (`mu` instead of `Mu`):**
   Making the mutex field private prevents external packages from locking or unlocking it directly, keeping all synchronization logic encapsulated inside your struct's methods.
4. **Place the mutex directly above the fields it protects:**
   By Go convention, declaring `mu sync.Mutex` directly above `count int` (with a brief comment) makes it immediately obvious to other developers which fields the mutex guards.
5. **Protect reads as well as writes:**
   Even a simple read like `Value()` needs to acquire the mutex. In Go, a concurrent read and write to the same memory without synchronization is a data race. *(Note: If your struct experiences heavy concurrent reads and infrequent writes, you can swap `sync.Mutex` for `sync.RWMutex` and use `RLock()`/`RUnlock()` in read-only methods.)*
6. **Rely on the zero value:**
   `sync.Mutex` does not require initialization. Declaring `var counter SafeCounter` creates a valid, unlocked mutex ready for immediate use.

### Verify with the Race Detector
You can verify that this code has no race conditions by running it with Go's built-in race detector:

```bash
go run -race main.go
```