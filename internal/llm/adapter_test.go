package llm

import (
	"strings"
	"testing"
)

type stubAdapter struct {
	response string
	last     LaunchOptions
}

func (s *stubAdapter) BuildLaunchCommand(options LaunchOptions) (string, error) {
	s.last = options
	return s.response, nil
}

func TestBuildLaunchCommandUsesRegisteredAdapter(t *testing.T) {
	adapter := &stubAdapter{response: "adapter-output"}
	RegisterAdapter("test-provider", adapter)
	t.Cleanup(func() {
		RegisterAdapter("test-provider", nil)
	})

	opts := LaunchOptions{
		Provider: "test-provider",
		Runtime:  RuntimeSettings{Command: "ignored"},
	}

	cmd, err := BuildLaunchCommand(opts)
	if err != nil {
		t.Fatalf("BuildLaunchCommand returned error: %v", err)
	}
	if cmd != "adapter-output" {
		t.Fatalf("unexpected command: %s", cmd)
	}

	if adapter.last.Provider != "test-provider" {
		t.Fatalf("adapter did not receive provider hint")
	}
}

func TestBuildLaunchCommandFallsBackToShellAdapter(t *testing.T) {
	opts := LaunchOptions{
		Env: map[string]string{
			"BD_ACTOR": "gastown/witness",
			"GT_ROLE":  "witness",
		},
		Runtime: RuntimeSettings{
			Command: "custom-cli",
			Args:    []string{"--flag"},
		},
	}

	cmd, err := BuildLaunchCommand(opts)
	if err != nil {
		t.Fatalf("BuildLaunchCommand returned error: %v", err)
	}

	expected := "export BD_ACTOR=gastown/witness GT_ROLE=witness && custom-cli --flag"
	if cmd != expected {
		t.Fatalf("command mismatch\nwant: %s\ngot:  %s", expected, cmd)
	}
}

func TestBuildLaunchCommandDerivesAdapterFromCommandPath(t *testing.T) {
	adapter := &stubAdapter{response: "path-adapter"}
	RegisterAdapter("custom", adapter)
	t.Cleanup(func() {
		RegisterAdapter("custom", nil)
	})

	opts := LaunchOptions{
		Runtime: RuntimeSettings{Command: "/usr/local/bin/custom"},
	}

	cmd, err := BuildLaunchCommand(opts)
	if err != nil {
		t.Fatalf("BuildLaunchCommand returned error: %v", err)
	}
	if cmd != "path-adapter" {
		t.Fatalf("expected adapter output, got %s", cmd)
	}
}

func TestOpenHandsAdapterAddsMetadata(t *testing.T) {
	adapter := NewOpenHandsAdapter()
	options := LaunchOptions{
		Env: map[string]string{
			"GT_ROLE":  "polecat",
			"BD_ACTOR": "gastown/polecats/toast",
		},
		Runtime: RuntimeSettings{},
		Role:    "polecat",
		RigName: "gastown",
		Actor:   "gastown/polecats/toast",
	}
	cmd, err := adapter.BuildLaunchCommand(options)
	if err != nil {
		t.Fatalf("BuildLaunchCommand returned error: %v", err)
	}
	if !strings.Contains(cmd, "OPENHANDS_AGENT_ROLE=polecat") {
		t.Fatalf("expected role metadata in command: %s", cmd)
	}
	if !strings.Contains(cmd, "OPENHANDS_AGENT_NAME=gastown/polecats/toast") {
		t.Fatalf("expected actor metadata in command: %s", cmd)
	}
	if !strings.Contains(cmd, "openhands --exp") {
		t.Fatalf("expected openhands invocation with --exp, got %s", cmd)
	}
}
