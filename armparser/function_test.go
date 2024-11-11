package armparser

import (
	"context"
	"log/slog"
	"strings"
	"testing"

	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
	"github.com/matt-FFFFFF/goarmfunctions/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCases []struct {
	desc     string
	in       string
	ctx      EvalContext
	expected any
	err      error
}

func runFunctionTest(ctx context.Context, t *testing.T, tcs testCases) {
	lgr := logger.NewDebugLogger()
	for _, tC := range tcs {
		t.Run(tC.desc, func(t *testing.T) {
			parser := New()
			if lgr.Enabled(ctx, slog.LevelDebug) {
				slog.DebugContext(ctx, "Lexing", slog.String("input", tC.in))
				reader := strings.NewReader(tC.in)
				lexer, err := parser.Lexer().Lex("debug", reader)
				if err != nil {
					lgr.Error("Lexer error", slog.String("error", err.Error()))
				}
				require.NoError(t, err)
				symbols := armlexer.TokenType2Str(parser.Lexer().Symbols())
				for tok, err := lexer.Next(); err == nil && !tok.EOF(); tok, err = lexer.Next() {
					lgr.Debug("Lexer token", slog.String("type", symbols[tok.Type]), slog.String("value", tok.Value))
				}
			}
			f, err := parser.ParseString("test", tC.in)
			require.NoError(t, err)
			result, err := f.Evaluate(ctx, tC.ctx)
			require.Equalf(t, tC.err, err, "unexpected evaluate error: %v", err)
			if err != nil {
				return
			}
			assert.Equalf(t, tC.expected, result, "unexpected result: %v", result)
		})
	}
}
