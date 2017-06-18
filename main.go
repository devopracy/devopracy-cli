// this is the main package for the devopracy cli

// go:generate go run ./scripts/generate-plugins.go

package main

import (
  "os"
  "math/rand"
  "time"
)

func main() {
  // Call realMain instead of doing the work here so we can use
  // `defer` statements within the function and have them work properly.
  // (defers aren't called with os.Exit)
  os.Exit(realMain())
}

// realMain is executed from main and returns the exit status to exit.
func realMain() int {
  // Call the real real Main
  return wrappedMain()
}

// wrappedMain is called only when we're wrapped panicwrap and
// returns the exit status to exit.
func wrappedMain() int {
  return 0
}

func init() {
  // Seed the random number generator
  rand.Seed(time.Now().UTC().UnixNano())
}
