package llm

import (
	"fmt"
)

// Adapter builds runtime-specific startup commands for agent sessions.
type Adapter interface {
	BuildLaunchCommand(options LaunchOptions) (string, error)
}

// LaunchOptions describes how an agent session should be launched.
type LaunchOptions struct {
	// Env contains inline environment exports that must be applied before the
	// runtime command executes (e.g., GT_ROLE, GT_RIG).
	Env map[string]string

	// Prompt optionally overrides the runtime's initial prompt.
	Prompt string

	// Runtime contains the resolved runtime configuration for the rig.
	Runtime RuntimeSettings

	// RigPath identifies the rig directory that owns the launch request.
	RigPath string

	// Role identifies the GT_ROLE for this launch.
	Role string

	// RigName is the GT_RIG value when available.
	RigName string

	// Actor is the BD_ACTOR identity for the agent.
	Actor string

	// Provider optionally specifies the adapter name to use. When empty, the
	// runtime command is used to derive an adapter key.
	Provider string
}

// RuntimeSettings is a copy of config.RuntimeConfig that avoids an import cycle.
type RuntimeSettings struct {
	Command       string
	Args          []string
	InitialPrompt string
	Model         string
}

// BuildLaunchCommand resolves the appropriate adapter and builds the final
// startup command string with inline exports.
func BuildLaunchCommand(options LaunchOptions) (string, error) {
	adapter := getAdapter(options)
	if adapter == nil {
		return "", fmt.Errorf("llm: no adapter registered")
	}
	return adapter.BuildLaunchCommand(options)
}
