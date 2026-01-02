package llm

import (
	"sort"
	"strings"
)

type shellAdapter struct{}

func (shellAdapter) BuildLaunchCommand(options LaunchOptions) (string, error) {
	var builder strings.Builder

	if len(options.Env) > 0 {
		exports := make([]string, 0, len(options.Env))
		for key, value := range options.Env {
			exports = append(exports, key+"="+value)
		}
		sort.Strings(exports)
		builder.WriteString("export ")
		builder.WriteString(strings.Join(exports, " "))
		builder.WriteString(" && ")
	}

	builder.WriteString(buildRuntimeCommand(options.Runtime, options.Prompt))

	return builder.String(), nil
}

func buildRuntimeCommand(settings RuntimeSettings, prompt string) string {
	cmd := settings.Command
	if cmd == "" {
		cmd = "claude"
	}

	args := settings.Args
	if args == nil {
		args = []string{"--dangerously-skip-permissions"}
	}

	commandParts := []string{cmd}
	if len(args) > 0 {
		commandParts = append(commandParts, args...)
	}

	command := strings.Join(commandParts, " ")

	effectivePrompt := prompt
	if effectivePrompt == "" {
		effectivePrompt = settings.InitialPrompt
	}

	if effectivePrompt == "" {
		return command
	}

	return command + " " + quoteForShell(effectivePrompt)
}

func quoteForShell(s string) string {
	escaped := strings.ReplaceAll(s, `\\`, `\\\\`)
	escaped = strings.ReplaceAll(escaped, `"`, `\\"`)
	return `"` + escaped + `"`
}
