package json2tree

import (
	"io"
	"strconv"
	"strings"
)

/*
TODO-List:
	For JSON lexer:
		1. hex integer leading with 0x
		2. '_' separator for integer, e.g. 123_456_789
		3. enable to store 0xFFFFFFFF_FFFFFFFF as uint64
		4. in parseString(), "\uhhhh" escape to single character
*/

//type IntegerOverflowError struct {
//	line int
//	row  int
//}
//
//func (err *IntegerOverflowError) Error() string {
//	return fmt.Sprintf("integer overflow at Line %d, Row %d", err.line, err.row)
//}

type TokenType int

const (
	TokenInvalid     TokenType = iota
	TokenLeftBrace             // '{'
	TokenRightBrace            // '}'
	TokenLeftSquare            // '['
	TokenRightSquare           // ']'
	TokenComma                 // ','
	TokenColon                 // ':'
	TokenInt
	TokenFloat
	TokenBool
	TokenNull
	TokenString

	// for test
	tokenTail
)

func (typ TokenType) String() string {
	switch typ {
	default:
		return "TokenInvalid"
	case TokenLeftBrace:
		return "TokenLeftBrace"
	case TokenRightBrace:
		return "TokenRightBrace"
	case TokenLeftSquare:
		return "TokenLeftSquare"
	case TokenRightSquare:
		return "TokenRightSquare"
	case TokenComma:
		return "TokenComma"
	case TokenColon:
		return "TokenColon"
	case TokenInt:
		return "TokenInt"
	case TokenFloat:
		return "TokenFloat"
	case TokenBool:
		return "TokenBool"
	case TokenNull:
		return "TokenNull"
	case TokenString:
		return "TokenString"
	}
}

//var charToIntForHex = buildCharToIntForHexMap()
//
//func buildCharToIntForHexMap() map[rune]int {
//	ret := make(map[rune]int)
//	for i, c := 0, '0'; i <= 9; i, c = i+1, c+1 {
//		ret[c] = i
//	}
//	for i, lo, up := 10, 'a', 'A'; i <= 15; i, lo, up = i+1, lo+1, up+1 {
//		ret[lo] = i
//		ret[up] = i
//	}
//	return ret
//}

type lexer struct {
	// results
	currentString strings.Builder
	currentInt    int64
	currentFloat  float64
	currentBool   bool

	// status
	currentToken TokenType
	reader       io.RuneReader
	currentChar  rune
	line         int
	row          int
	ioError      error // errors happened during io except io.EOF
	startLine    int
	startRow     int
}

func (lex *lexer) Reset(reader io.RuneReader) {
	lex.line = 1
	lex.row = 0
	lex.currentToken = TokenInvalid
	lex.reader = reader
	lex.ioError = nil

	if lex.getChar() == 0 {
		return
	}
	if shouldSkip(lex.currentChar) {
		lex.skip()
	}
}

// <<==== accessors begin ====>

func (lex *lexer) StartAt() (line int, row int) {
	return lex.startLine, lex.startRow
}

func (lex *lexer) Int() int64 {
	return lex.currentInt
}

func (lex *lexer) Float() float64 {
	return lex.currentFloat
}

func (lex *lexer) Bool() bool {
	return lex.currentBool
}

func (lex *lexer) String() string {
	return lex.currentString.String()
}

func (lex *lexer) IoError() error {
	return lex.ioError
}

func (lex *lexer) HasNext() bool {
	return lex.ioError == nil && lex.currentChar != 0
}

func (lex *lexer) CurrentType() TokenType {
	return lex.currentToken
}

/*
Next

Errors:
	IOError: unexpect EOF
	UnexpectCharacterError: lex error
*/
func (lex *lexer) Next() (TokenType, error) {
	if DEBUG {
		if !lex.HasNext() {
			panic("call Next() when HasNext() returns false")
		}
	}
	lex.currentToken = TokenInvalid
	err := lex.moveNext()
	if err != nil {
		lex.currentChar = 0
		lex.currentToken = TokenInvalid
		return lex.currentToken, err
	}
	return lex.currentToken, nil
}

// <----- accessors end ----->

// <<===== utils begin =====>

func (lex *lexer) recordPosition() {
	lex.startLine = lex.line
	lex.startRow = lex.row
}

func (lex *lexer) unexpectError() error {
	return &UnexpectCharacterError{
		Line: lex.line,
		Row:  lex.row,
		Char: lex.currentChar,
	}
}

func (lex *lexer) eofError() error {
	return &IOError{
		Inner: lex.ioError,
		Line:  lex.line,
		Row:   lex.row,
	}
}

func (lex *lexer) getChar() rune {
	if lex.currentChar == '\n' {
		lex.line++
		lex.row = 1
	} else {
		lex.row++
	}
	var err error
	lex.currentChar, _, err = lex.reader.ReadRune()
	if err != nil {
		lex.currentChar = 0
		lex.ioError = err
	}
	return lex.currentChar
}

func shouldSkip(ch rune) bool {
	switch ch {
	case ' ', '\t', '\n', '\r', '/':
		return true
	default:
		return false
	}
}

func isSpace(ch rune) bool {
	switch ch {
	case ' ', '\t', '\n', '\r':
		return true
	default:
		return false
	}
}

// <----- utils end ----->

// <<==== process functions begin ====>

func (lex *lexer) skip() {
	for {
		if isSpace(lex.currentChar) {
			lex.skipSpace()
			continue
		} else if lex.currentChar == '/' {
			lex.skipComment()
			continue
		}
		return
	}
}

func (lex *lexer) skipSpace() {
	for isSpace(lex.currentChar) {
		if lex.getChar() == 0 {
			return
		}
	}
}

func (lex *lexer) skipComment() {
	if DEBUG {
		if lex.currentChar != '/' {
			panic("comment not start with '/'")
		}
	}
	if lex.getChar() == 0 {
		return
	}
	if lex.currentChar == '*' {
		// multi-Line comment
		lastChar := rune(0)
		for {
			if lex.getChar() == 0 {
				return
			}
			if lastChar == '*' {
				if lex.currentChar == '/' {
					lex.getChar()
					return
				}
			}
			lastChar = lex.currentChar
		}
	} else {
		// single-Line comment
		for {
			if lex.getChar() == 0 {
				return
			}
			if lex.currentChar == '\n' {
				lex.getChar()
				return
			}
		}
	}
}

func (lex *lexer) parseNumber() (err error) {
	if DEBUG {
		if lex.currentChar != '-' && (lex.currentChar < '0' || lex.currentChar > '9') {
			panic("number not start with digit or '-'")
		}
	}
	buffer := strings.Builder{}
	isFloat := false
	for {
		buffer.WriteRune(lex.currentChar)
		if lex.getChar() == 0 {
			break
		}
		if lex.currentChar >= '0' && lex.currentChar <= '9' {
			continue
		}
		if lex.currentChar == '.' {
			if isFloat {
				return lex.unexpectError()
			}
			isFloat = true
			continue
		}
		break
	}
	if isFloat {
		lex.currentFloat, err = strconv.ParseFloat(buffer.String(), 64)
		if err != nil {
			return err
		}
		lex.currentToken = TokenFloat
		lex.getChar()
		return nil
	}
	lex.currentInt, err = strconv.ParseInt(buffer.String(), 10, 64)
	if err != nil {
		return err
	}
	lex.currentToken = TokenInt
	// lex.getChar() currentChar is already not a part of number
	return nil
}

func (lex *lexer) parseString() error {
	if DEBUG {
		if lex.currentChar != '"' {
			panic("string not start with '\"'")
		}
	}
	if lex.getChar() == 0 {
		return lex.eofError()
	}
	lex.currentString.Reset()
	for lex.currentChar != '"' {
		if lex.currentChar == '\\' {
			if lex.getChar() == 0 {
				return lex.eofError()
			}
			switch lex.currentChar {
			case 'b':
				lex.currentString.WriteRune('\b')
			case 'f':
				lex.currentString.WriteRune('\f')
			case '/':
				lex.currentString.WriteRune('/')
			case 'n':
				lex.currentString.WriteRune('\n')
			case '\r':
				lex.currentString.WriteRune('\r')
			case 't':
				lex.currentString.WriteRune('\t')
			case '"':
				lex.currentString.WriteRune('"')
			case '\\':
				lex.currentString.WriteRune('\\')
			case '\'':
				lex.currentString.WriteRune('\'')
			default:
				lex.currentString.WriteRune('\\')
				lex.currentString.WriteRune(lex.currentChar)
			}
		} else {
			lex.currentString.WriteRune(lex.currentChar)
		}
		if lex.getChar() == 0 {
			return lex.eofError()
		}
	}
	lex.currentToken = TokenString
	lex.getChar()
	return nil
}

func (lex *lexer) parseTrue() error {
	if DEBUG {
		if lex.currentChar != 'T' && lex.currentChar != 't' {
			panic("true not start with 't' or 'T'")
		}
	}
	for _, c := range []rune{'r', 'u', 'e'} {
		if lex.getChar() == 0 {
			return lex.eofError()
		}
		if lex.currentChar != c && lex.currentChar != (c-'a'+'A') {
			return lex.unexpectError()
		}
	}
	lex.currentToken = TokenBool
	lex.currentBool = true
	lex.getChar()
	return nil
}

func (lex *lexer) parseFalse() (err error) {
	if DEBUG {
		if lex.currentChar != 'f' && lex.currentChar != 'F' {
			panic("false not start with 'f' or 'F'")
		}
	}
	for _, c := range []rune{'a', 'l', 's', 'e'} {
		if lex.getChar() == 0 {
			return lex.eofError()
		}
		if lex.currentChar != c && lex.currentChar != (c-'a'+'A') {
			return lex.unexpectError()
		}
	}
	lex.currentToken = TokenBool
	lex.currentBool = false
	lex.getChar()
	return nil
}

func (lex *lexer) parseNull() (err error) {
	if DEBUG {
		if lex.currentChar != 'N' && lex.currentChar != 'n' {
			panic("null not start with 'n' or 'N'")
		}
	}
	for _, c := range []rune{'u', 'l', 'l'} {
		if lex.getChar() == 0 {
			return lex.eofError()
		}
		if lex.currentChar != c && lex.currentChar != (c-'a'+'A') {
			return lex.unexpectError()
		}
	}
	lex.currentToken = TokenNull
	lex.getChar()
	return nil
}

func (lex *lexer) moveNext() (err error) {
	err = lex.moveNext2()
	if err != nil {
		return err
	}
	if shouldSkip(lex.currentChar) {
		lex.skip()
	}
	return nil
}

func (lex *lexer) moveNext2() (err error) {
	for {
		lex.recordPosition()
		switch lex.currentChar {
		case 0:
			return lex.eofError()
		default:
			return lex.unexpectError()
		case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return lex.parseNumber()
		//case ' ', '\t', '\n':
		//	lex.skipSpace()
		//case '/':
		//	lex.skipComment()
		case 't', 'T':
			return lex.parseTrue()
		case 'f', 'F':
			return lex.parseFalse()
		case 'n', 'N':
			return lex.parseNull()
		case '"':
			return lex.parseString()
		case '{':
			lex.currentToken = TokenLeftBrace
			lex.getChar()
			return
		case '}':
			lex.currentToken = TokenRightBrace
			lex.getChar()
			return
		case '[':
			lex.currentToken = TokenLeftSquare
			lex.getChar()
			return
		case ']':
			lex.currentToken = TokenRightSquare
			lex.getChar()
			return
		case ',':
			lex.currentToken = TokenComma
			lex.getChar()
			return
		case ':':
			lex.currentToken = TokenColon
			lex.getChar()
			return
		}
	}
}

// <----- process functions end ----->
