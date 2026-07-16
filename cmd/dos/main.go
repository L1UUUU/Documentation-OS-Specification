// Command dos exposes the Documentation OS engine through a thin scriptable CLI.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"documentation-os/engine"
)

// main executes the CLI and returns its deterministic exit status to the shell.
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run parses one CLI invocation and delegates all repository behavior to the engine.
func run(args []string, stdout, stderr io.Writer) int {
	root, jsonOutput, commandArgs, err := parseGlobalOptions(args)
	if err != nil {
		writeCLIError(stderr, jsonOutput, err)
		return 2
	}
	if len(commandArgs) == 0 || commandArgs[0] == "help" || commandArgs[0] == "--help" || commandArgs[0] == "-h" {
		_, _ = io.WriteString(stdout, cliUsage)
		return 0
	}
	if commandArgs[0] == "version" {
		return writeVersion(root, jsonOutput, stdout, stderr)
	}
	instance, err := engine.New(root)
	if err != nil {
		writeCLIError(stderr, jsonOutput, err)
		return 1
	}
	switch commandArgs[0] {
	case "init":
		if err := instance.Initialize(); err != nil {
			writeCLIError(stderr, jsonOutput, err)
			return 1
		}
		return writeOutput(stdout, jsonOutput, map[string]any{"initialized": true, "root": instance.Root}, "Initialized Documentation OS repository at "+instance.Root)
	case "inspect":
		result, err := instance.Inspect()
		return finishResult(result, err, stdout, stderr, jsonOutput, 1)
	case "validate":
		result, err := instance.Validate()
		if err != nil {
			writeCLIError(stderr, jsonOutput, err)
			return 1
		}
		if !result.Passed() {
			_ = writeOutput(stdout, jsonOutput, result, result.String())
			return 1
		}
		return writeOutput(stdout, jsonOutput, result, result.String())
	case "health":
		result, err := instance.Health()
		return finishResult(result, err, stdout, stderr, jsonOutput, 1)
	case "sync":
		result, err := instance.Synchronize()
		return finishResult(result, err, stdout, stderr, jsonOutput, 1)
	case "generate":
		return runGenerate(instance, commandArgs[1:], stdout, stderr, jsonOutput)
	case "complete":
		return runComplete(instance, commandArgs[1:], stdout, stderr, jsonOutput)
	case "migrate":
		return runMigrate(instance, commandArgs[1:], stdout, stderr, jsonOutput)
	default:
		writeCLIError(stderr, jsonOutput, fmt.Errorf("unknown command %q", commandArgs[0]))
		return 2
	}
}

// parseGlobalOptions extracts options that are shared by every command.
func parseGlobalOptions(args []string) (string, bool, []string, error) {
	root := "."
	jsonOutput := false
	commandArgs := make([]string, 0, len(args))
	for index := 0; index < len(args); index++ {
		arg := args[index]
		switch {
		case arg == "--json":
			jsonOutput = true
		case arg == "--root":
			if index+1 >= len(args) {
				return root, jsonOutput, nil, errors.New("--root requires a path")
			}
			root = args[index+1]
			index++
		case strings.HasPrefix(arg, "--root="):
			root = strings.TrimPrefix(arg, "--root=")
		default:
			commandArgs = append(commandArgs, arg)
		}
	}
	return root, jsonOutput, commandArgs, nil
}

// runGenerate dispatches the Generate command category.
func runGenerate(instance *engine.Engine, args []string, stdout, stderr io.Writer, jsonOutput bool) int {
	if len(args) < 2 {
		writeCLIError(stderr, jsonOutput, errors.New("generate requires `work <slug>` or `id <category> <draft-path>"))
		return 2
	}
	switch args[0] {
	case "work":
		result, err := instance.GenerateWork(args[1])
		return finishResult(result, err, stdout, stderr, jsonOutput, 1)
	case "id":
		if len(args) < 3 {
			writeCLIError(stderr, jsonOutput, errors.New("generate id requires a category and draft path"))
			return 2
		}
		result, err := instance.AllocateKnowledgeIdentifier(args[1], args[2])
		return finishResult(result, err, stdout, stderr, jsonOutput, 1)
	default:
		writeCLIError(stderr, jsonOutput, fmt.Errorf("unknown generate category %q", args[0]))
		return 2
	}
}

// runComplete parses caller-supplied outcome and delegates Complete.
func runComplete(instance *engine.Engine, args []string, stdout, stderr io.Writer, jsonOutput bool) int {
	if len(args) == 0 {
		writeCLIError(stderr, jsonOutput, errors.New("complete requires a Work slug and --outcome"))
		return 2
	}
	slug := args[0]
	outcome := ""
	for index := 1; index < len(args); index++ {
		if args[index] == "--outcome" && index+1 < len(args) {
			outcome = args[index+1]
			index++
		} else if strings.HasPrefix(args[index], "--outcome=") {
			outcome = strings.TrimPrefix(args[index], "--outcome=")
		}
	}
	if outcome == "" {
		writeCLIError(stderr, jsonOutput, errors.New("complete requires --outcome=succeeded|cancelled|superseded|failed"))
		return 2
	}
	result, err := instance.Complete(slug, outcome)
	if err != nil {
		_ = writeOutput(stdout, jsonOutput, result, err.Error())
		return 1
	}
	return writeOutput(stdout, jsonOutput, result, fmt.Sprintf("Completed Work %s (%s)", slug, outcome))
}

// runMigrate dispatches the supported migration target.
func runMigrate(instance *engine.Engine, args []string, stdout, stderr io.Writer, jsonOutput bool) int {
	target := engine.SpecificationVersion
	for index := 0; index < len(args); index++ {
		if args[index] == "--to" && index+1 < len(args) {
			target = args[index+1]
			index++
		} else if strings.HasPrefix(args[index], "--to=") {
			target = strings.TrimPrefix(args[index], "--to=")
		}
	}
	result, err := instance.Migrate(target)
	if err != nil {
		_ = writeOutput(stdout, jsonOutput, result, err.Error())
		return 1
	}
	return writeOutput(stdout, jsonOutput, result, "Migration completed")
}

// writeVersion emits independent version dimensions without requiring a valid repository.
func writeVersion(root string, jsonOutput bool, stdout, stderr io.Writer) int {
	result := map[string]string{
		"specification_version": engine.SpecificationVersion,
		"engine_version":        engine.EngineVersion,
		"cli_version":           engine.EngineVersion,
		"repository_profile":    engine.ProfileName,
	}
	if instance, err := engine.New(root); err == nil {
		info := instance.Version()
		result["repository_version"] = info.RepositoryVersion
	}
	if err := writeOutput(stdout, jsonOutput, result, fmt.Sprintf("Documentation OS %s; engine %s", engine.SpecificationVersion, engine.EngineVersion)); err != 0 {
		writeCLIError(stderr, jsonOutput, errors.New("write version output"))
		return 1
	}
	return 0
}

// finishResult writes a successful result or an operation error.
func finishResult(result any, err error, stdout, stderr io.Writer, jsonOutput bool, errorCode int) int {
	if err != nil {
		writeCLIError(stderr, jsonOutput, err)
		return errorCode
	}
	return writeOutput(stdout, jsonOutput, result, humanResult(result))
}

// humanResult renders useful stable summaries for the non-JSON CLI surface.
func humanResult(result any) string {
	switch value := result.(type) {
	case engine.InspectReport:
		return fmt.Sprintf("%s: %d Knowledge document(s), %d active Work(s), %d completed Work(s)", value.RepositoryProfile, value.KnowledgeDocuments, value.ActiveWorks, value.CompletedWorks)
	case engine.HealthReport:
		return fmt.Sprintf("Health: %s (%d/100)", value.Level, value.Score)
	case engine.IndexResult:
		return fmt.Sprintf("INDEX regenerated: %s", value.Path)
	case engine.SyncResult:
		return fmt.Sprintf("Synchronization complete; no Knowledge edits required; INDEX: %s", value.Index.Path)
	case engine.WorkResult:
		return fmt.Sprintf("Generated Work %s at %s", value.Slug, value.Path)
	case engine.AllocationResult:
		return fmt.Sprintf("Allocated %s: %s -> %s", value.Identifier, value.OldPath, value.NewPath)
	case engine.MigrationReport:
		return fmt.Sprintf("Migration %s -> %s completed", value.Before, value.After)
	default:
		return "success"
	}
}

// writeOutput emits JSON for automation or concise human-readable output.
func writeOutput(stdout io.Writer, jsonOutput bool, value any, human string) int {
	if jsonOutput {
		encoder := json.NewEncoder(stdout)
		encoder.SetIndent("", "  ")
		if err := encoder.Encode(value); err != nil {
			return 1
		}
		return 0
	}
	_, err := fmt.Fprintln(stdout, human)
	if err != nil {
		return 1
	}
	return 0
}

// writeCLIError emits a deterministic CLI error representation.
func writeCLIError(stderr io.Writer, jsonOutput bool, err error) {
	if jsonOutput {
		_ = json.NewEncoder(stderr).Encode(map[string]string{"error": err.Error()})
		return
	}
	_, _ = fmt.Fprintln(stderr, "error:", err)
}

const cliUsage = `Documentation OS CLI

Usage:
  dos [--root PATH] [--json] <command>

Commands:
  init
  version
  inspect
  validate
  health
  sync
  generate work <slug>
  generate id <category> <draft-path>
  complete <slug> --outcome <succeeded|cancelled|superseded|failed>
  migrate [--to 1.0]
`
