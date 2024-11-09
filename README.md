# goarmfunctions

This go module provides an implementation of some ARM template functions in Go.
It is a work in progress and not all functions are implemented yet, and may never will be.

Some functions are implemented in a way that they are not 100% compatible with the ARM template functions.
For example, the `format()` function does not support format specifiers.

Finally, at the moment all the output is a string. This may change in the future as the evolves.s

## Comparison Functions

- `equals()`

## Deployment Functions

- `parameters()`

## String Functions

- `format()`
- `if()`
- `parameters()`
- `replace()`

## Example

```go
parser := armparser.New()
f, err := parser.ParseString("example", "[if(equals('a', 'b'), 'a is equal to b', 'a is not equal to b')]")
if err != nil {
  panic(err)
}
result, err := f.Evaluate(nil)
if err != nil {
  panic(err)
}
fmt.Println(result)
// Output: a is not equal to b
```
