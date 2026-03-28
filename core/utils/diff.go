package utils

import (
	"fmt"
	"strings"
)

type DiffOp int

const (
	DiffEqual  DiffOp = iota
	DiffInsert
	DiffDelete
)

type DiffLine struct {
	Op   DiffOp
	Text string
}

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

func myersDiff(a, b []string) []DiffLine {
	n, m := len(a), len(b)
	max := n + m
	offset := max

	v := make([]int, 2*max+2)
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
				x = v[k+1+offset]
			} else {
				x = v[k-1+offset] + 1
			}
			y := x - k
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

	result := make([]DiffLine, 0, n+m)
	x, y := n, m

	for d := finalD; d > 0; d-- {
		snap := traces[d]
		k := x - y

		var prevK int
		if k == -d || (k != d && snap[k-1+offset] < snap[k+1+offset]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX := snap[prevK+offset]
		prevY := prevX - prevK

		for x > prevX && y > prevY {
			x--
			y--
			result = append([]DiffLine{{DiffEqual, a[x]}}, result...)
		}

		if x == prevX {
			y--
			result = append([]DiffLine{{DiffInsert, b[y]}}, result...)
		} else {
			x--
			result = append([]DiffLine{{DiffDelete, a[x]}}, result...)
		}
	}

	for x > 0 && y > 0 {
		x--
		y--
		result = append([]DiffLine{{DiffEqual, a[x]}}, result...)
	}

	return result
}

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
		for i < n && diff[i].Op == DiffEqual {
			i++
		}
		if i >= n {
			break
		}

		hunkStart := i - contextLines
		if hunkStart < 0 {
			hunkStart = 0
		}

		j := i
		for j < n {
			if diff[j].Op != DiffEqual {
				j++
				continue
			}
			eqEnd := j
			for eqEnd < n && diff[eqEnd].Op == DiffEqual {
				eqEnd++
			}
			if eqEnd >= n || eqEnd-j > 2*contextLines {
				trailing := eqEnd - j
				if trailing > contextLines {
					trailing = contextLines
				}
				j += trailing
				break
			}
			j = eqEnd
		}
		hunkEnd := j
		if hunkEnd > n {
			hunkEnd = n
		}

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
