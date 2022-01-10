package json2tree

import "fmt"

// <<<==== lexer error begin ====>>>

type UnexpectCharacterError struct {
	Line int
	Row  int
	Char rune
}

func (err *UnexpectCharacterError) Error() string {
	return fmt.Sprintf("unexpect character: '%c' at line %d, row %d", err.Char, err.Line, err.Row)
}

type IOError struct {
	Inner error
	Line  int
	Row   int
}

func (err *IOError) Error() string {
	return fmt.Sprintf("unexpect IO error at line %d, row %d, may be caused by: %s", err.Line, err.Row, err.Inner.Error())
}

// <<----- lexer error end ----->>

// <<<==== parser error begin ====>>>

type UnexpectTokenError struct {
	Line    int
	Row     int
	Token   TokenType
	content string
}

func (e *UnexpectTokenError) Error() string {
	return fmt.Sprintf("unexpect token (%s) at line %d, row %d",
		e.content, e.Line, e.Row)
}

type ParseError struct {
	Line  int
	Row   int
	Inner error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parsing error at line %d, row %d, may caused by %s",
		e.Line, e.Row, e.Inner.Error())
}

type DuplicateKeyError struct {
	Line int
	Row  int
	Key  string
}

func (e *DuplicateKeyError) Error() string {
	return fmt.Sprintf("duplicate key (%s) at line %d, row %d",
		e.Key, e.Line, e.Row)
}

// <<----- parser error end ----->>
