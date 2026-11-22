package cli

import (
	"fmt"
	"io"
	"strings"
	"sync/atomic"
	"testing"
	"time"
)

func TestBatchParallelRunsTargetsConcurrentlyAndReportsInOrder(t *testing.T) {
	var active int32
	var maximum int32
	run := func(output io.Writer, args []string) error {
		current := atomic.AddInt32(&active, 1)
		for {
			previous := atomic.LoadInt32(&maximum)
			if current <= previous || atomic.CompareAndSwapInt32(&maximum, previous, current) {
				break
			}
		}
		time.Sleep(30 * time.Millisecond)
		atomic.AddInt32(&active, -1)
		_, err := fmt.Fprintln(output, args[0])
		return err
	}
	var output strings.Builder
	failures := runBatchTargets(&output, "FORGE", []string{"one", "two", "three"}, func(writer io.Writer, args []string) error {
		return run(writer, args)
	}, true, false)
	if len(failures) != 0 {
		t.Fatalf("unexpected batch failures: %v", failures)
	}
	if atomic.LoadInt32(&maximum) < 2 {
		t.Fatal("parallel targets did not overlap")
	}
	text := output.String()
	last := -1
	for _, target := range []string{"one", "two", "three"} {
		index := strings.Index(text, target+"\n")
		if index <= last {
			t.Fatalf("parallel output was not reported in input order: %q", text)
		}
		last = index
	}
}

func TestBatchFailFastStopsSequentialTargets(t *testing.T) {
	runs := 0
	run := func(_ io.Writer, args []string) error {
		runs++
		if args[0] == "first" {
			return fmt.Errorf("failed")
		}
		return nil
	}
	failures := runBatchTargets(&strings.Builder{}, "VERIFY",
		[]string{"first", "second"}, func(output io.Writer, args []string) error {
			return run(output, args)
		}, false, true)
	if len(failures) != 1 || runs != 1 {
		t.Fatalf("fail-fast ran %d targets and returned %d errors", runs, len(failures))
	}
}
