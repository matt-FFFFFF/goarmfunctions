package goarmfunctions

import (
	"context"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
	"github.com/matt-FFFFFF/goarmfunctions/armparser"
	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

// Evaluate parses and evaluates an ARM template expression with the given EvalContext and FuncRegistry.
func Evaluate(ctx context.Context, expr string, evalCtx armparser.EvalContext, registry *armparser.FuncRegistry, lgr logger.Logger) (any, error) {
	if lgr == nil {
		lgr = logger.LoggerFromContext(ctx)
	}
	ctx = context.WithValue(ctx, logger.LoggerContextKey, lgr)
	lgr.Debug("Evaluate", slog.String("input", expr))
	defer lgr.Debug("Evaluate done")
	parser := armparser.New()
	if lgr.Enabled(ctx, slog.LevelDebug) {
		slog.DebugContext(ctx, "Lexing", slog.String("input", expr))
		reader := strings.NewReader(expr)
		lexer, err := parser.Lexer().Lex("debug", reader)
		if err != nil {
			lgr.Error("Lexer error", slog.String("error", err.Error()))
			return nil, err
		}
		symbols := armlexer.TokenType2Str(parser.Lexer().Symbols())
		for tok, err := lexer.Next(); err == nil && !tok.EOF(); tok, err = lexer.Next() {
			lgr.Debug("Lexer token", slog.String("type", symbols[tok.Type]), slog.String("value", tok.Value))
		}
	}
	lgr.Debug("Parsing", slog.String("input", expr))
	f, err := parser.ParseString("test", expr)
	if err != nil {
		lgr.Error("Parser error", slog.String("error", err.Error()))
		return nil, err
	}
	if registry == nil {
		registry = armparser.DefaultRegistry()
	}
	return f.Evaluate(ctx, evalCtx, registry)
}

// LexAndParse is a convenience wrapper that parses and evaluates an ARM template expression
// using a map-based evaluation context and the default function registry.
func LexAndParse(ctx context.Context, s string, evalCtx map[string]any, lgr logger.Logger) (any, error) {
	return Evaluate(ctx, s, armparser.FromMap(evalCtx), nil, lgr)
}
