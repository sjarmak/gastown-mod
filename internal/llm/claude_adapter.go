package llm

// ClaudeAdapter is the default adapter that mirrors the previous shell-based
// Claude launcher semantics.
type ClaudeAdapter struct {
	shellAdapter
}

// NewClaudeAdapter constructs a ClaudeAdapter instance.
func NewClaudeAdapter() Adapter {
	return ClaudeAdapter{}
}

func init() {
	RegisterAdapter("claude", NewClaudeAdapter())
}
