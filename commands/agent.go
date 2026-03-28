package commands

import (
	"fmt"
	"os"

	"github.com/KambojRajan/ship/core/llm"
	"github.com/KambojRajan/ship/core/trace"
)

const agentModel = "claude-3-5-haiku-20241022"

// Agent stages all working-tree changes, asks Claude to generate a commit
// message based on the staged diff and the developer's description, then
// commits automatically.
func Agent(description string, path string) error {
	// 1. Stage all changes (ship add .)
	doneStep := trace.Step("Add(.)")
	if err := Add("."); err != nil {
		doneStep(err)
		return fmt.Errorf("staging changes: %w", err)
	}
	doneStep(nil)

	// 2. Get staged diff as a string (index vs HEAD)
	doneStep = trace.Step("DiffString(cached)")
	diffStr, err := DiffString(path, true)
	doneStep(err)
	if err != nil {
		return fmt.Errorf("computing staged diff: %w", err)
	}

	// 3. Strip ANSI colour codes
	doneStep = trace.Step("StripANSI")
	cleanDiff := llm.StripANSI(diffStr)
	doneStep(nil)

	if cleanDiff == "" {
		fmt.Fprintln(os.Stdout, "Nothing staged to commit.")
		return nil
	}

	trace.Meta("diff_bytes", fmt.Sprintf("%d", len(cleanDiff)))
	trace.Meta("description_len", fmt.Sprintf("%d", len(description)))

	// 4. Ask Claude to generate a commit message
	prompt := buildAgentPrompt(description, cleanDiff)
	doneStep = trace.Step("llm.Complete")
	commitMessage, err := llm.Complete(agentModel, 256, prompt)
	doneStep(err)
	if err != nil {
		return fmt.Errorf("calling Claude API: %w", err)
	}

	// 5. Commit with the generated message
	fmt.Fprintf(os.Stdout, "Committing with message:\n  %s\n", commitMessage)
	doneStep = trace.Step("Commit")
	if err := Commit(commitMessage, path); err != nil {
		doneStep(err)
		return fmt.Errorf("committing: %w", err)
	}
	doneStep(nil)

	fmt.Fprintln(os.Stdout, "Done. All changes staged and committed.")
	return nil
}

func buildAgentPrompt(description, diff string) string {
	return `You are an expert software engineer. A developer has described their change and provided the staged diff. Generate a single commit message in Conventional Commits format.

Rules:
- Output ONLY the commit message — no explanation, no preamble, no quotes.
- Keep it under 72 characters.
- Use one of: feat, fix, docs, style, refactor, perf, test, build, ci, chore.
- Prefer the developer's intent (description) over the raw diff when they differ.
- Be specific.

Developer's description:
` + description + `

Staged diff:
` + diff
}
