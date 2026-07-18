// Command external-consumer-smoke verifies the published Engine from a fresh
// module that cannot inherit this repository's development replace directive.
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const engineModule = "github.com/L1UUUU/Documentation-OS-Specification/engine"

func main() {
	version := flag.String("version", "v0.1.0-rc.3", "published Engine module version to verify")
	flag.Parse()
	if strings.TrimSpace(*version) == "" {
		fatalf("version must not be empty")
	}

	consumerRoot, err := os.MkdirTemp("", "dos-external-consumer-*")
	if err != nil {
		fatalf("create external consumer: %v", err)
	}
	defer os.RemoveAll(consumerRoot)

	goMod := fmt.Sprintf("module example.com/dos-external-consumer\n\ngo 1.22\n\nrequire %s %s\n", engineModule, *version)
	mainSource := fmt.Sprintf("package main\n\nimport (\n\t\"fmt\"\n\tengine %q\n)\n\nfunc main() { fmt.Print(engine.EngineVersion) }\n", engineModule)
	writeFile(filepath.Join(consumerRoot, "go.mod"), goMod)
	writeFile(filepath.Join(consumerRoot, "main.go"), mainSource)

	runGo(consumerRoot, "mod", "tidy")
	resolvedReplace := strings.TrimSpace(runGo(consumerRoot, "list", "-m", "-f", "{{if .Replace}}{{.Replace.Dir}}{{end}}", engineModule))
	if resolvedReplace != "" {
		fatalf("external consumer unexpectedly resolved a replace directive to %q", resolvedReplace)
	}
	updatedGoMod, err := os.ReadFile(filepath.Join(consumerRoot, "go.mod"))
	if err != nil {
		fatalf("read external go.mod: %v", err)
	}
	if containsReplaceDirective(string(updatedGoMod)) {
		fatalf("external consumer go.mod contains a replace directive")
	}
	runGo(consumerRoot, "test", "./...")
	fmt.Printf("external consumer resolved %s@%s without replace\n", engineModule, *version)
}

func containsReplaceDirective(goMod string) bool {
	for _, line := range strings.Split(goMod, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "replace ") || strings.TrimSpace(line) == "replace (" {
			return true
		}
	}
	return false
}

func writeFile(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		fatalf("write %s: %v", filepath.Base(path), err)
	}
}

func runGo(directory string, args ...string) string {
	command := exec.Command("go", args...)
	command.Dir = directory
	command.Env = append(os.Environ(), "GOWORK=off")
	output, err := command.CombinedOutput()
	if err != nil {
		fatalf("go %s failed: %v\n%s", strings.Join(args, " "), err, output)
	}
	return string(output)
}

func fatalf(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
