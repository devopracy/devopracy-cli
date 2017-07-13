// this is the main package for the devopracy cli

// go:generate go run ./scripts/generate-plugins.go

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/mitchellh/panicwrap"
	"github.com/mitchellh/prefixedio"
	"github.com/munjeli/devopracy-cli/version"
)

// this pattern of three mains is something from Hashicorp. You'll see by
// the third call what makes it so great.
func main() {
	// Call realMain instead of doing the work here so we can use
	// `defer` statements within the function and have them work properly.
	// (defers aren't called with os.Exit)
	os.Exit(realMain())
}

// realMain is executed from main and returns the exit status to exit.
func realMain() int {
	var wrapConfig panicwrap.WrapConfig

	if !panicwrap.Wrapped(&wrapConfig) {

		// this is generating a UUID for the application
		// run. I'm not sure how to use it yet but I'll start
		// off with it in place because I like the idea of
		// tracking an operation with a id. But it always
		// returns a nil error based on rand.read which is
		// ignored. Hmmm
		UUID, _ := uuid.GenerateUUID()
		os.Setenv("DEVO_RUN_UUID", UUID)

		// Determine where logs should go in general (requested by the user)
		// This is going to be useful down the road when we start to optimize
		// the builds... setting up our logs all in the same place with some
		// logic will make it a lot more obvious where to look and what to
		// monitor. I've never seen this error actually using their stuff.
		logWriter, err := logOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't setup log output: %s", err)
			return 1
		}
		if logWriter == nil {
			logWriter = ioutil.Discard
		}

		// We always send logs to a temporary file that we use in case
		// there is a panic. Otherwise, we delete it.
		logTempFile, err := ioutil.TempFile("", "devo-log")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't setup logging tempfile: %s", err)
			return 1
		}
		defer os.Remove(logTempFile.Name())
		defer logTempFile.Close()

		// Tell the logger to log to this file
		os.Setenv(EnvLog, "")
		os.Setenv(EnvLogFile, "")

		// Setup the prefixed readers that send data properly to
		// stdout/stderr. The style of the naming is interesting here.
		// I'm going to start with it same but I'm not sure it's
		// low cognitive load... outR? If there's consistency in this
		// style we'll keep it.
		doneCh := make(chan struct{})
		outR, outW := io.Pipe()
		go copyOutput(outR, doneCh)

		// Create the configuration for panicwrap and wrap our executable
		wrapConfig.Handler = panicHandler(logTempFile)
		wrapConfig.Writer = io.MultiWriter(logTempFile, logWriter)
		wrapConfig.Stdout = outW
		exitStatus, err := panicwrap.Wrap(&wrapConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Couldn't start Devo: %s", err)
			return 1
		}

		// If >= 0, we're the parent, so just exit
		if exitStatus >= 0 {
			// Close the stdout writer so that our copy process can finish
			outW.Close()

			// Wait for the output copying to finish
			<-doneCh

			return exitStatus
		}

		// We're the child, so just close the tempfile we made in order to
		// save file handles since the tempfile is only used by the parent.
		logTempFile.Close()
	}

	// Call the real real Main
	return wrappedMain()
}

// Panicwrap is one of my favorite features of Hashicorp's clis -
// when the application crashes with a bug it will give the user a
// error message and the url to submit an issue. Remember: delagation
// is also a kind of automation! How can we get users to do what we
// want?
// wrappedMain is called only when we're wrapped panicwrap and
// returns the exit status to exit.
func wrappedMain() int {
	// If there is no explicit number of Go threads to use, then set it.
	// I've never used this but it's a nice feature. Possibly useful when
	// we get into optimization.
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	log.SetOutput(os.Stderr)

	// I don't know that I'll support running the cli on Windows or Mac, but
	// for now I'll stay in the game. Definitely not supporting builds of anything
	// not open source for the devo cloud. We'll optimize this later when I Determine
	// how much work it is to run my own code on the proprietary arches.
	log.Printf("[INFO] Devo version: %s", version.FormattedVersion())
	log.Printf("Devo Target OS/Arch: %s %s", runtime.GOOS, runtime.GOARCH)
	log.Printf("Built with Go Version: %s", runtime.Version())
	return 0
}

// copyOutput uses output prefixes to determine whether data on stdout
// should go to stdout or stderr. This is due to panicwrap using stderr
// as the log and error channel.
func copyOutput(r io.Reader, doneCh chan<- struct{}) {
	defer close(doneCh)

	pr, err := prefixedio.NewReader(r)
	if err != nil {
		panic(err)
	}

	stderrR, err := pr.Prefix(ErrorPrefix)
	if err != nil {
		panic(err)
	}
	stdoutR, err := pr.Prefix(OutputPrefix)
	if err != nil {
		panic(err)
	}
	defaultR, err := pr.Prefix("")
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		io.Copy(os.Stderr, stderrR)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, stdoutR)
	}()
	go func() {
		defer wg.Done()
		io.Copy(os.Stdout, defaultR)
	}()

	wg.Wait()
}

func init() {
	// Seed the random number generator
	rand.Seed(time.Now().UTC().UnixNano())
}
