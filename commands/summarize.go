package commands

import (
	"fmt"
	"os"

	"github.com/KambojRajan/ship/core/llm"
	"github.com/KambojRajan/ship/core/trace"
)

const summarizeModel = "claude-3-5-haiku-20241022"

// Summarize computes the current diff, sends it to the Claude API, and prints
// a suggested commit message in Conventional Commits format.
//
// cached=false: diff working tree vs index.
// cached=true:  diff index vs HEAD (staged changes).
func Summarize(path string, cached bool) error {
	// 1. Compute diff as a string
	doneStep := trace.Step("DiffString")
	diffStr, err := DiffString(path, cached)
	doneStep(err)
	if err != nil {
		return fmt.Errorf("computing diff: %w", err)
	}

	// 2. Strip ANSI colour codes
	doneStep = trace.Step("StripANSI")
	cleanDiff := llm.StripANSI(diffStr)
	doneStep(nil)

	if cleanDiff == "" {
		fmt.Fprintln(os.Stdout, "No changes to summarize.")
		return nil
	}

	trace.Meta("diff_bytes", fmt.Sprintf("%d", len(cleanDiff)))

	// 3. Call Claude API
	prompt := buildSummarizePrompt(cleanDiff)
	doneStep = trace.Step("llm.Complete")
	suggestion, err := llm.Complete(summarizeModel, 256, prompt)
	doneStep(err)
	if err != nil {
		return fmt.Errorf("calling Claude API: %w", err)
	}

	// 4. Print result
	fmt.Fprintln(os.Stdout, suggestion)
	return nil
}

func buildSummarizePrompt(diff string) string {
	return `You are an expert software engineer. Analyze the following code diff and suggest a single commit message in Conventional Commits format (e.g. "feat: add user login", "fix: handle nil pointer in parser").

Rules:
- Output ONLY the commit message — no explanation, no preamble, no quotes.
- Keep it under 72 characters.
- Use one of: feat, fix, docs, style, refactor, perf, test, build, ci, chore.
- Be specific about what changed.

Diff:
` + diff
}
