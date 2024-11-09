package armparser

// Equals returns true if the first argument is equal to the second argument.
func Equals(f *FunctionCall, ctx EvalContext) (any, error) {
	if len(f.Args) != 2 {
		return nil, NewArgumentError("equals", 2, len(f.Args))
	}
	arg1, err := f.Args[0].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	arg2, err := f.Args[1].Evaluate(ctx)
	if err != nil {
		return nil, err
	}
	return arg1 == arg2, nil
}
