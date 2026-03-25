package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	tracecmd "github.com/KambojRajan/ship/commands/trace"
	coretrace "github.com/KambojRajan/ship/core/trace"
	"github.com/spf13/cobra"
)

var (
	traceFormat     string
	traceOtel       bool
	traceOtelOutput string
)

var TraceCmd = &cobra.Command{
	Use:   "trace <command> [args...]",
	Short: "Trace internal execution of a ship command",
	Long: `Trace records and prints every internal pipeline step for the given command.

SUPPORTED COMMANDS
  commit <message>          trace the commit pipeline
  add    <path> [paths...]  trace the staging pipeline
  status [path]             trace the status pipeline

OUTPUT FORMATS
  --format text   (default) human-readable colour output + summary footer
  --format json   NDJSON – one JSON object per step; pipe to jq for filtering

OPENTELEMETRY
  --otel                    export spans to an OTel-compatible JSON file
  --otel-output <file>      destination file (default: stderr)

EXAMPLES
  ship trace commit "feat: add auth"
  ship trace add src/ README.md
  ship trace status
  ship trace commit "fix: bug" --format json | jq 'select(.status=="error")'
  ship trace commit "feat: observability" --otel --otel-output spans.json`,

	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		operation := args[0]
		rest := args[1:]

		var sink coretrace.Sink
		var pretty *coretrace.PrettySink
		var otelEnd func()

		switch {
		case traceOtel:
			outW, cleanup, err := openOTelOutput(traceOtelOutput)
			if err != nil {
				return err
			}
			if cleanup != nil {
				defer cleanup()
			}

			tp, shutdown, err := coretrace.NewOTelProvider(outW)
			if err != nil {
				return fmt.Errorf("otel init: %w", err)
			}
			defer func() { _ = shutdown(context.Background()) }()

			var otelSink *coretrace.OTelSink
			otelSink, otelEnd = coretrace.NewOTelSink(context.Background(), tp, operation)
			defer func() {
				if otelEnd != nil {
					otelEnd()
				}
			}()
			sink = otelSink

		case traceFormat == "json":
			sink = coretrace.NewJSONSink(os.Stdout)

		default:
			pretty = coretrace.NewPrettySink(os.Stdout)
			sink = pretty
		}

		if pretty != nil {
			ops := tracecmd.SupportedOps()
			fmt.Fprintf(os.Stdout,
				"\n  \033[1m\033[36mtrace ship %s\033[0m   \033[90m[%s]\033[0m\n\n",
				operation, strings.Join(ops, " · "),
			)
		}

		restore := coretrace.SetSink(sink)
		defer restore()

		var runErr error
		runErr = tracecmd.Dispatch(operation, &tracecmd.ExecContext{
			RepoBasePath: repoBasePath,
			Args:         rest,
		})

		if pretty != nil {
			pretty.PrintSummary()
		}
		if traceOtel && traceOtelOutput != "" && traceOtelOutput != "-" {
			fmt.Fprintf(os.Stdout, "  Spans written to %s\n\n", traceOtelOutput)
		}

		return runErr
	},
}

func openOTelOutput(path string) (*os.File, func(), error) {
	if path == "" || path == "-" {
		return os.Stderr, nil, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("cannot open otel output %q: %w", path, err)
	}
	return f, func() { _ = f.Close() }, nil
}

func init() {
	TraceCmd.Flags().StringVar(&traceFormat, "format", "text", `Output format: "text" or "json"`)
	TraceCmd.Flags().BoolVar(&traceOtel, "otel", false, "Export pipeline steps as OpenTelemetry spans")
	TraceCmd.Flags().StringVar(&traceOtelOutput, "otel-output", "-", `File to write OTEL spans to ("-" = stderr)`)
	rootCmd.AddCommand(TraceCmd)
}
