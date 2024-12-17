package goarmfunctions

import (
	"context"
	"log/slog"
	"strings"

	"github.com/matt-FFFFFF/goarmfunctions/armlexer"
	"github.com/matt-FFFFFF/goarmfunctions/armparser"
	"github.com/matt-FFFFFF/goarmfunctions/logger"
)

func LexAndParse(ctx context.Context, s string, evalCtx armparser.EvalContext, lgr logger.Logger) (any, error) {
	if lgr == nil {
		lgr = logger.LoggerFromContext(ctx)
	}
	ctx = context.WithValue(ctx, logger.LoggerContextKey, lgr)
	lgr.Debug("LexAndParse", slog.String("input", s), slog.Any("evalCtx", evalCtx))
	defer lgr.Debug("LexAndParse done")
	parser := armparser.New()
	if lgr.Enabled(ctx, slog.LevelDebug) {
		slog.DebugContext(ctx, "Lexing", slog.String("input", s))
		reader := strings.NewReader(s)
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
	lgr.Debug("Parsing", slog.String("input", s))
	f, err := parser.ParseString("test", s)
	if err != nil {
		lgr.Error("Parser error", slog.String("error", err.Error()))
		return nil, err
	}
	return f.Evaluate(ctx, evalCtx)
}
