package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	batchparser "secret-bitpack/internal/batch/parser"
)

type batchJob struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	File      string `json:"file"`
	Status    string `json:"status"`
	StartedAt int64  `json:"started_at"`
	EndedAt   int64  `json:"ended_at,omitempty"`
	Error     string `json:"error,omitempty"`
}

func bitpackBatch(output io.Writer, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: batch <file.bpbatch>")
	}
	file, err := os.Open(args[0])
	if err != nil {
		return fmt.Errorf("open batch file: %w", err)
	}
	plan, parseErr := batchparser.Parse(file)
	closeErr := file.Close()
	if parseErr != nil {
		return fmt.Errorf("parse batch file: %w", parseErr)
	}
	if closeErr != nil {
		return fmt.Errorf("close batch file: %w", closeErr)
	}

	started := time.Now()
	job := batchJob{
		ID:   started.UTC().Format("20060102T150405.000000000Z"),
		Name: plan.Name, File: args[0], Status: "running",
		StartedAt: started.Unix(),
	}
	jobPath := dataPath("batch", job.ID+".json")
	if err := writeJSON(jobPath, job, 0o600); err != nil {
		return fmt.Errorf("record batch job: %w", err)
	}
	fmt.Fprintf(output, "[BATCH] %s\n", plan.Name)

	stages := []struct {
		name    string
		targets []string
		run     func(io.Writer, []string) error
	}{
		{name: "FORGE", targets: plan.Forge, run: bitpackForge},
		{name: "VERIFY", targets: plan.Verify, run: bitpackVerify},
		{name: "PUBLISH", targets: plan.Publish, run: bitpackPublish},
		{name: "ASSEMBLE", targets: plan.Assemble, run: bitpackAssemble},
		{name: "IGNITE", targets: plan.Ignite, run: bitpackIgnite},
	}
	var stageErrors []error
	for _, stage := range stages {
		if len(stage.targets) == 0 {
			continue
		}
		fmt.Fprintf(output, "[%s]\n", stage.name)
		stageErrors = append(stageErrors,
			runBatchTargets(output, stage.name, stage.targets, stage.run, plan.Parallel, plan.FailFast)...)
		if len(stageErrors) > 0 && plan.FailFast {
			break
		}
	}

	job.EndedAt = time.Now().Unix()
	job.Status = "succeeded"
	if len(stageErrors) > 0 {
		job.Status = "failed"
		job.Error = stageErrors[0].Error()
	}
	if err := writeJSON(jobPath, job, 0o600); err != nil {
		return fmt.Errorf("save batch result: %w", err)
	}
	if len(stageErrors) > 0 {
		return fmt.Errorf("batch %q failed: %w", plan.Name, stageErrors[0])
	}
	fmt.Fprintln(output, "[SUCCESS] Batch complete")
	return nil
}

func runBatchTargets(
	output io.Writer,
	stageName string,
	targets []string,
	run func(io.Writer, []string) error,
	parallel, failFast bool,
) []error {
	if !parallel || len(targets) < 2 {
		var failures []error
		for _, target := range targets {
			fmt.Fprintf(output, "  %s\n", target)
			if err := run(output, []string{target}); err != nil {
				failures = append(failures, fmt.Errorf("%s %s: %w", stageName, target, err))
				if failFast {
					break
				}
			}
		}
		return failures
	}

	buffers := make([]bytes.Buffer, len(targets))
	failures := make([]error, len(targets))
	var wait sync.WaitGroup
	for index, target := range targets {
		fmt.Fprintf(output, "  %s\n", target)
		wait.Add(1)
		go func(index int, target string) {
			defer wait.Done()
			if err := run(&buffers[index], []string{target}); err != nil {
				failures[index] = fmt.Errorf("%s %s: %w", stageName, target, err)
			}
		}(index, target)
	}
	wait.Wait()
	for index := range buffers {
		if _, err := output.Write(buffers[index].Bytes()); err != nil {
			failures[index] = fmt.Errorf("write %s output for %s: %w", stageName, targets[index], err)
		}
	}
	var result []error
	for _, err := range failures {
		if err != nil {
			result = append(result, err)
		}
	}
	return result
}

func bitpackBatchStatus(output io.Writer) error {
	root := dataPath("batch")
	entries, err := os.ReadDir(root)
	if os.IsNotExist(err) {
		fmt.Fprintln(output, "No batch jobs.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("list batch jobs: %w", err)
	}
	count := 0
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, entry.Name()))
		if err != nil {
			return fmt.Errorf("read batch job %s: %w", entry.Name(), err)
		}
		var job batchJob
		if err := json.Unmarshal(data, &job); err != nil {
			return fmt.Errorf("decode batch job %s: %w", entry.Name(), err)
		}
		fmt.Fprintf(output, "%s\t%s\t%s\n", job.ID, job.Name, job.Status)
		count++
	}
	if count == 0 {
		fmt.Fprintln(output, "No batch jobs.")
	}
	return nil
}
