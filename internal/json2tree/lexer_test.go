package json2tree

import (
	assertions "github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestSpaceNotZero(t *testing.T) {
	assert := assertions.New(t)
	assert.NotEqual(0, ' ')
}

func TestStringToken(t *testing.T) {
	assert := assertions.New(t)

	origin := `"\n\t\"abz_ABZ-123\'''\b\f\/\\\c"`
	expect := "\n\t\"abz_ABZ-123'''\b\f/\\\\c"
	lex := lexer{}
	lex.Reset(strings.NewReader(origin))
	assert.True(lex.HasNext())
	typ, err := lex.Next()
	assert.Nil(err)
	assert.Equal(TokenString, typ)
	assert.Equal(expect, lex.String())
}

func TestIntegerToken(t *testing.T) {
	assert := assertions.New(t)

	origin := `-123`
	expect := int64(-123)
	lex := lexer{}
	lex.Reset(strings.NewReader(origin))
	assert.True(lex.HasNext())
	typ, err := lex.Next()
	assert.Nil(err)
	assert.Equal(TokenInt, typ)
	assert.Equal(expect, lex.Int())
}

func TestFloatToken(t *testing.T) {
	assert := assertions.New(t)

	origin := `-12.3`
	expect := float64(-12.3)
	lex := lexer{}
	lex.Reset(strings.NewReader(origin))
	assert.True(lex.HasNext())
	typ, err := lex.Next()
	assert.Nil(err)
	assert.Equal(TokenFloat, typ)
	assert.Equal(expect, lex.Float())
}

func TestBoolToken(t *testing.T) {
	assert := assertions.New(t)

	func() {
		origins := []string{`true`, `True`, `TRUE`, `tRUE`}
		for _, origin := range origins {
			lex := lexer{}
			lex.Reset(strings.NewReader(origin))
			assert.True(lex.HasNext())
			typ, err := lex.Next()
			assert.Nil(err)
			assert.Equal(TokenBool, typ)
			assert.True(lex.Bool())
		}
	}()

	func() {
		origins := []string{`false`, `False`, `FALSE`, `fALSE`}
		for _, origin := range origins {
			lex := lexer{}
			lex.Reset(strings.NewReader(origin))
			assert.True(lex.HasNext())
			typ, err := lex.Next()
			assert.Nil(err)
			assert.Equal(TokenBool, typ)
			assert.False(lex.Bool())
		}
	}()
}

func TestNullToken(t *testing.T) {
	assert := assertions.New(t)

	origins := []string{`null`, `Null`, `NULL`, `nULL`}
	for _, origin := range origins {
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.Nil(err)
		if err != nil {
			t.Log(err)
		}
		assert.Equal(TokenNull, typ)
	}
}

func TestSimpleToken(t *testing.T) {
	assert := assertions.New(t)
	origin := "{}:,[]"
	typs := []TokenType{
		TokenLeftBrace,
		TokenRightBrace,
		TokenColon,
		TokenComma,
		TokenLeftSquare,
		TokenRightSquare,
	}
	lex := lexer{}
	lex.Reset(strings.NewReader(origin))
	for _, expect := range typs {
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.Nil(err)
		assert.Equal(expect, typ)
	}
}

type Token struct {
	typ TokenType
	obj interface{}
}

func runAndCompare(t *testing.T, origin string, tokens []Token) {
	assert := assertions.New(t)

	lex := new(lexer)
	lex.Reset(strings.NewReader(origin))
	for i, token := range tokens {
		typ := token.typ
		obj := token.obj
		assert.True(lex.HasNext(), i)
		tokenType, err := lex.Next()
		assert.Nil(err, i)
		if err != nil {
			t.Log(err)
		}
		switch typ {
		case TokenString:
			assert.Equal(TokenString, tokenType, i)
			str := lex.String()
			assert.Equal(obj.(string), str, i)
		case TokenInt:
			assert.Equal(TokenInt, tokenType)
			switch val := obj.(type) {
			case int64:
				assert.Equal(val, lex.Int(), i)
			case int:
				assert.Equal(int64(val), lex.Int(), i)
			}
		case TokenFloat:
			assert.Equal(TokenFloat, tokenType, i)
			assert.Equal(obj.(float64), lex.Float(), i)
		case TokenBool:
			assert.Equal(TokenBool, tokenType, i)
			assert.Equal(obj.(bool), lex.Bool(), i)
		default:
			assert.Equal(typ, tokenType, i)
		}
	}
}

func TestSkipSpace(t *testing.T) {
	func() {
		origin := `
	{
	"str "  :  123
}  `
		tokens := []Token{
			{TokenLeftBrace, nil},
			{TokenString, "str "},
			{TokenColon, nil},
			{TokenInt, 123},
			{TokenRightBrace, nil},
		}
		runAndCompare(t, origin, tokens)
	}()

	func() {
		origin := "\t{\n\n\"str\\t\"\t:\t-32\t}"
		tokens := []Token{
			{TokenLeftBrace, nil},
			{TokenString, "str\t"},
			{TokenColon, nil},
			{TokenInt, -32},
			{TokenRightBrace, nil},
		}
		runAndCompare(t, origin, tokens)
	}()
}

func TestSkipComment(t *testing.T) {
	func() {
		origin := `{/*"str // this is a fake comment"/*"this is a wrong string"*/"str"://123
456/*
*/}//`
		tokens := []Token{
			{TokenLeftBrace, nil},
			{TokenString, "str"},
			{TokenColon, nil},
			{TokenInt, 456},
			{TokenRightBrace, nil},
		}
		runAndCompare(t, origin, tokens)
	}()
	func() {
		origin := `
// useless
	{
	"str // this is a fake comment"/*"this is a wrong string"*/  :  //123
321
}// `
		tokens := []Token{
			{TokenLeftBrace, nil},
			{TokenString, "str // this is a fake comment"},
			{TokenColon, nil},
			{TokenInt, 321},
			{TokenRightBrace, nil},
		}
		runAndCompare(t, origin, tokens)
	}()
}

func TestUnexpectCharacterError(t *testing.T) {
	assert := assertions.New(t)

	func() {
		origin := "123..321"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		unexpect, ok := err.(*UnexpectCharacterError)
		assert.True(ok)
		assert.Equal(1, unexpect.Line)
		assert.Equal(5, unexpect.Row)
		assert.Equal('.', unexpect.Char)
	}()

	func() {
		origin := "ture"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		unexpect, ok := err.(*UnexpectCharacterError)
		assert.True(ok)
		assert.Equal(1, unexpect.Line)
		assert.Equal(2, unexpect.Row)
		assert.Equal('u', unexpect.Char)
	}()

	func() {
		origin := "flase"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		unexpect, ok := err.(*UnexpectCharacterError)
		assert.True(ok)
		assert.Equal(1, unexpect.Line)
		assert.Equal(2, unexpect.Row)
		assert.Equal('l', unexpect.Char)
	}()

	func() {
		origin := "123x321"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.Nil(err)
		assert.Equal(TokenInt, typ)

		typ, err = lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		unexpect, ok := err.(*UnexpectCharacterError)
		assert.True(ok)
		assert.Equal(1, unexpect.Line)
		assert.Equal(4, unexpect.Row)
		assert.Equal('x', unexpect.Char)
	}()

	func() {
		origin := "123O321"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.Nil(err)
		assert.Equal(TokenInt, typ)

		typ, err = lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		unexpect, ok := err.(*UnexpectCharacterError)
		assert.True(ok)
		assert.Equal(1, unexpect.Line)
		assert.Equal(4, unexpect.Row)
		assert.Equal('O', unexpect.Char)
	}()
}

func TestUnexpectIOError(t *testing.T) {
	assert := assertions.New(t)

	func() {
		origin := ""
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.False(lex.HasNext())
		err := lex.IoError()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(io.EOF, err)
	}()

	func() {
		origin := "   // \n/**/"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.False(lex.HasNext())
		err := lex.IoError()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(io.EOF, err)
	}()

	func() {
		origin := "\"str"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		ioe, ok := err.(*IOError)
		assert.True(ok)
		assert.Equal(io.EOF, ioe.Inner)
	}()

	func() {
		origin := "\"\\"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		ioe, ok := err.(*IOError)
		assert.True(ok)
		assert.Equal(io.EOF, ioe.Inner)
	}()

	func() {
		origin := "\""
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		ioe, ok := err.(*IOError)
		assert.True(ok)
		assert.Equal(io.EOF, ioe.Inner)
	}()

	func() {
		origin := "fal"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		ioe, ok := err.(*IOError)
		assert.True(ok)
		assert.Equal(io.EOF, ioe.Inner)
	}()

	func() {
		origin := "tr"
		lex := lexer{}
		lex.Reset(strings.NewReader(origin))
		assert.True(lex.HasNext())
		typ, err := lex.Next()
		assert.NotNil(err)
		t.Log(err)
		assert.Equal(TokenInvalid, typ)
		ioe, ok := err.(*IOError)
		assert.True(ok)
		assert.Equal(io.EOF, ioe.Inner)
	}()
}

func TestTokenTypeString(t *testing.T) {
	assert := assertions.New(t)

	type void struct{}
	member := void{}
	set := make(map[string]void)
	for typ := TokenInvalid; typ < tokenTail; typ++ {
		str := typ.String()
		_, ok := set[str]
		assert.False(ok)
		set[str] = member
	}
}
