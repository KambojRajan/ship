package utils

import (
	"fmt"
	"strings"
)

// DiffOp is the type of a single diff operation.
type DiffOp int

const (
	DiffEqual  DiffOp = iota
	DiffInsert        // line present in new but not in old
	DiffDelete        // line present in old but not in new
)

// DiffLine is a single entry in a diff result.
type DiffLine struct {
	Op   DiffOp
	Text string
}

// SplitLines splits s into lines, dropping the final empty element that
// results from a trailing newline.
func SplitLines(s string) []string {
	if s == "" {
		return nil
	}
	lines := strings.Split(s, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// DiffLines computes a line-level diff between oldLines and newLines using
// Myers' O(ND) algorithm.
func DiffLines(oldLines, newLines []string) []DiffLine {
	n, m := len(oldLines), len(newLines)

	if n == 0 && m == 0 {
		return nil
	}
	if n == 0 {
		out := make([]DiffLine, m)
		for i, l := range newLines {
			out[i] = DiffLine{DiffInsert, l}
		}
		return out
	}
	if m == 0 {
		out := make([]DiffLine, n)
		for i, l := range oldLines {
			out[i] = DiffLine{DiffDelete, l}
		}
		return out
	}

	return myersDiff(oldLines, newLines)
}

// myersDiff runs the Myers shortest-edit-script algorithm and backtracks the
// trace to produce an ordered slice of DiffLines.
func myersDiff(a, b []string) []DiffLine {
	n, m := len(a), len(b)
	max := n + m
	offset := max

	// v[k+offset] = furthest x reached along diagonal k.
	v := make([]int, 2*max+2)

	// traces[d] is a snapshot of v taken *before* computing d-paths.
	// This lets us backtrack: traces[d] holds the (d-1)-path frontier.
	traces := make([][]int, 0, max+1)

	finalD := 0
	found := false

	for d := 0; d <= max && !found; d++ {
		snap := make([]int, len(v))
		copy(snap, v)
		traces = append(traces, snap)

		for k := -d; k <= d; k += 2 {
			var x int
			if k == -d || (k != d && v[k-1+offset] < v[k+1+offset]) {
				// Move down from diagonal k+1 (insert from b).
				x = v[k+1+offset]
			} else {
				// Move right from diagonal k-1 (delete from a).
				x = v[k-1+offset] + 1
			}
			y := x - k
			// Slide diagonally as far as possible (equal lines).
			for x < n && y < m && a[x] == b[y] {
				x++
				y++
			}
			v[k+offset] = x
			if x >= n && y >= m {
				found = true
				finalD = d
				break
			}
		}
	}

	// --- Backtrack ---
	// We reconstruct the edit script from (n,m) back to (0,0).
	result := make([]DiffLine, 0, n+m)
	x, y := n, m

	for d := finalD; d > 0; d-- {
		snap := traces[d] // v after (d-1) steps = what guided step d
		k := x - y

		// Determine whether step d was an insert (come from k+1) or delete (from k-1).
		var prevK int
		if k == -d || (k != d && snap[k-1+offset] < snap[k+1+offset]) {
			prevK = k + 1 // insert: x unchanged, y advanced
		} else {
			prevK = k - 1 // delete: x advanced, y unchanged
		}

		prevX := snap[prevK+offset]
		prevY := prevX - prevK

		// Walk back along the snake (equal lines) before this edit.
		for x > prevX && y > prevY {
			x--
			y--
			result = append([]DiffLine{{DiffEqual, a[x]}}, result...)
		}

		// Record the non-diagonal (edit) move.
		if x == prevX {
			// Insert: y came from prevY → prevY+1 → ... current y.
			y--
			result = append([]DiffLine{{DiffInsert, b[y]}}, result...)
		} else {
			// Delete.
			x--
			result = append([]DiffLine{{DiffDelete, a[x]}}, result...)
		}
	}

	// Any remaining position means equal lines at the very start.
	for x > 0 && y > 0 {
		x--
		y--
		result = append([]DiffLine{{DiffEqual, a[x]}}, result...)
	}

	return result
}

// FormatUnifiedDiff formats diff lines into a coloured unified diff string.
// Returns "" when there are no actual changes.
func FormatUnifiedDiff(oldPath, newPath string, diff []DiffLine, contextLines int) string {
	if len(diff) == 0 {
		return ""
	}

	hasChanges := false
	for _, d := range diff {
		if d.Op != DiffEqual {
			hasChanges = true
			break
		}
	}
	if !hasChanges {
		return ""
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("diff --ship a/%s b/%s\n", oldPath, newPath))
	sb.WriteString(fmt.Sprintf("--- a/%s\n", oldPath))
	sb.WriteString(fmt.Sprintf("+++ b/%s\n", newPath))

	n := len(diff)
	i := 0
	for i < n {
		// Advance past unchanged lines to find the next edit.
		for i < n && diff[i].Op == DiffEqual {
			i++
		}
		if i >= n {
			break
		}

		// Extend hunk start backward for context.
		hunkStart := i - contextLines
		if hunkStart < 0 {
			hunkStart = 0
		}

		// Walk forward, merging adjacent edit regions separated by ≤ 2*contextLines
		// equal lines into a single hunk.
		j := i
		for j < n {
			if diff[j].Op != DiffEqual {
				j++
				continue
			}
			// Count consecutive equal lines starting at j.
			eqEnd := j
			for eqEnd < n && diff[eqEnd].Op == DiffEqual {
				eqEnd++
			}
			// If the gap is large enough, end this hunk here.
			if eqEnd >= n || eqEnd-j > 2*contextLines {
				trailing := eqEnd - j
				if trailing > contextLines {
					trailing = contextLines
				}
				j += trailing
				break
			}
			// Small gap: include it and continue to the next edit region.
			j = eqEnd
		}
		hunkEnd := j
		if hunkEnd > n {
			hunkEnd = n
		}

		// Compute 1-based starting line numbers for old and new files.
		oldStart, newStart := 1, 1
		for k := 0; k < hunkStart; k++ {
			if diff[k].Op != DiffInsert {
				oldStart++
			}
			if diff[k].Op != DiffDelete {
				newStart++
			}
		}

		oldCount, newCount := 0, 0
		for k := hunkStart; k < hunkEnd; k++ {
			if diff[k].Op != DiffInsert {
				oldCount++
			}
			if diff[k].Op != DiffDelete {
				newCount++
			}
		}

		sb.WriteString(fmt.Sprintf("@@ -%d,%d +%d,%d @@\n", oldStart, oldCount, newStart, newCount))

		for k := hunkStart; k < hunkEnd; k++ {
			switch diff[k].Op {
			case DiffEqual:
				sb.WriteString(" " + diff[k].Text + "\n")
			case DiffDelete:
				sb.WriteString(ShipRed + "-" + diff[k].Text + Reset + "\n")
			case DiffInsert:
				sb.WriteString(ShipGreen + "+" + diff[k].Text + Reset + "\n")
			}
		}

		i = hunkEnd
	}

	return sb.String()
}
