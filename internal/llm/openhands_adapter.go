package llm

// OpenHandsAdapter configures the OpenHands CLI for agent sessions while
// preserving Gas Town metadata via environment variables.
type OpenHandsAdapter struct {
	shell shellAdapter
}

// NewOpenHandsAdapter constructs a new adapter instance.
func NewOpenHandsAdapter() Adapter {
	return OpenHandsAdapter{}
}

func (a OpenHandsAdapter) BuildLaunchCommand(opts LaunchOptions) (string, error) {
	env := cloneEnv(opts.Env)
	if opts.Role != "" {
		env["OPENHANDS_AGENT_ROLE"] = opts.Role
	}
	if opts.Actor != "" {
		env["OPENHANDS_AGENT_NAME"] = opts.Actor
	}
	if opts.RigName != "" {
		env["OPENHANDS_RIG_NAME"] = opts.RigName
	}
	if opts.Runtime.Model != "" {
		env["OPENHANDS_MODEL"] = opts.Runtime.Model
	}
	env["OPENHANDS_PROVIDER"] = "openhands"

	runtime := opts.Runtime
	if runtime.Command == "" {
		runtime.Command = "openhands"
	}
	if len(runtime.Args) == 0 {
		runtime.Args = []string{"--exp"}
	} else if !containsArg(runtime.Args, "--exp") {
		runtime.Args = append(runtime.Args, "--exp")
	}

	next := opts
	next.Env = env
	next.Runtime = runtime

	return a.shell.BuildLaunchCommand(next)
}

func cloneEnv(src map[string]string) map[string]string {
	if len(src) == 0 {
		return map[string]string{}
	}
	clone := make(map[string]string, len(src))
	for k, v := range src {
		clone[k] = v
	}
	return clone
}

func containsArg(args []string, needle string) bool {
	for _, arg := range args {
		if arg == needle {
			return true
		}
	}
	return false
}

func init() {
	RegisterAdapter("openhands", NewOpenHandsAdapter())
}
