package json2tree

import (
	"fmt"
	"github.com/SnowPhoenix0105/cfgm/internal/tree"
	"io"
	"strconv"
	"strings"
)

type parser struct {
	lex        lexer
	innerError error
	walker     tree.Walker
}

func (env *parser) Reset(reader io.RuneReader) error {
	env.lex.Reset(reader)
	env.innerError = nil

	if env.getToken() == TokenInvalid {
		return env.lexerError()
	}
	return nil
}

func (env *parser) assertType(name string, targets ...TokenType) {
	if len(targets) == 0 {
		return
	}
	typ := env.lex.CurrentType()
	targetStrings := make([]string, 0)
	for _, target := range targets {
		if typ == target {
			return
		}
		targetStrings = append(targetStrings, target.String())
	}
	panic(name + " not start with " +
		strings.Join(targetStrings[:len(targetStrings)-1], ", ") +
		" or " + targetStrings[len(targetStrings)-1] + " but " + typ.String())
}

func (env *parser) tokenType() TokenType {
	return env.lex.CurrentType()
}

func (env *parser) getToken() TokenType {
	if !env.lex.HasNext() {
		env.innerError = env.lex.IoError()
		return TokenInvalid
	}
	typ, err := env.lex.Next()
	if err != nil {
		env.innerError = err
		return TokenInvalid
	}
	return typ
}

func (env *parser) wrapError(err error) error {
	line, row := env.lex.StartAt()
	return &ParseError{
		Line:  line,
		Row:   row,
		Inner: err,
	}
}

func (env *parser) lexerError() error {
	if env.innerError != nil {
		return env.wrapError(env.innerError)
	}
	if env.lex.IoError() != nil {
		return env.wrapError(env.lex.IoError())
	}
	return env.unexpectError()
}

func (env *parser) unexpectError() error {
	var content string
	switch env.tokenType() {
	case TokenString:
		content = "\"" + env.lex.String() + "\""
	case TokenBool:
		if env.lex.Bool() {
			content = "true"
		} else {
			content = "false"
		}
	case TokenInt:
		content = strconv.FormatInt(env.lex.Int(), 10)
	case TokenFloat:
		content = strconv.FormatFloat(env.lex.Float(), 'f', 10, 64)
	case TokenNull:
		content = "null"
	case TokenLeftSquare:
		content = "'['"
	case TokenRightSquare:
		content = "']'"
	case TokenLeftBrace:
		content = "'{'"
	case TokenRightBrace:
		content = "'}'"
	case TokenComma:
		content = "','"
	case TokenColon:
		content = "':'"
	default:
		content = "<invalid>"
	}
	line, row := env.lex.StartAt()
	return &UnexpectTokenError{
		Line:    line,
		Row:     row,
		Token:   env.tokenType(),
		content: content,
	}
}

func (env *parser) inFirstSetForNode() bool {
	switch env.tokenType() {
	case TokenInt,
		TokenString,
		TokenBool,
		TokenFloat,
		TokenNull,
		TokenLeftBrace,
		TokenLeftSquare:
		return true
	default:
		return false
	}
}

func (env *parser) parseNode() error {
	if DEBUG {
		env.assertType("Node",
			TokenInt,
			TokenString,
			TokenBool,
			TokenFloat,
			TokenNull,
			TokenLeftBrace,
			TokenLeftSquare)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing Node")
		}
	}
	switch env.tokenType() {

	case TokenInt:
		integer := env.lex.Int()
		env.getToken()
		env.walker.SetInt(integer)
		env.walker.SetNullFor(tree.NodeKeyInt, false)
		env.walker.SetFloat(float64(integer))
		env.walker.SetNullFor(tree.NodeKeyFloat, false)
		return nil

	case TokenFloat:
		f := env.lex.Float()
		env.getToken()
		env.walker.SetFloat(f)
		env.walker.SetNullFor(tree.NodeKeyFloat, false)
		return nil

	case TokenBool:
		b := env.lex.Bool()
		env.getToken()
		env.walker.SetBool(b)
		env.walker.SetNullFor(tree.NodeKeyBool, false)
		return nil

	case TokenString:
		str := env.lex.String()
		env.getToken()
		env.walker.SetString(str)
		env.walker.SetNullFor(tree.NodeKeyString, false)

		return nil
	case TokenNull:
		env.getToken()
		setNull(env.walker)
		return nil
	case TokenLeftBrace:
		return env.parseObject()

	case TokenLeftSquare:
		return env.parseList()

	default:
		return env.unexpectError()
	}
}

func (env *parser) parseObject() error {
	// '{'
	if DEBUG {
		env.assertType("Object", TokenLeftBrace)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing Object")
		}
	}
	if env.getToken() == TokenInvalid {
		return env.lexerError()
	}

	// KvPairs
	if env.tokenType() != TokenRightBrace && env.tokenType() != TokenString {
		return env.unexpectError()
	}
	err := env.parseKvPairs()
	if err != nil {
		return err
	}

	// '}'
	if env.tokenType() != TokenRightBrace {
		return env.unexpectError()
	}
	env.getToken()
	return nil
}

func (env *parser) parseKvPairs() error {
	if DEBUG {
		env.assertType("KvPairs", TokenString, TokenRightBrace)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing KvPairs")
		}
	}
	keySet := make(map[string]emptyType)
	for env.tokenType() == TokenString {
		key, err := env.parseKvPair()
		if err != nil {
			return nil
		}
		_, ok := keySet[key]
		if ok {
			line, row := env.lex.StartAt()
			return &DuplicateKeyError{
				Line: line,
				Row:  row,
				Key:  key,
			}
		}
	}
	return nil
}

func (env *parser) parseKvPair() (string, error) {
	// string
	if DEBUG {
		env.assertType("KvPair", TokenString)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing KvPair")
		}
	}
	key := env.lex.String()
	if env.getToken() == TokenInvalid {
		return "", env.lexerError()
	}

	// ':'
	if env.tokenType() != TokenColon {
		return "", env.unexpectError()
	}
	if env.getToken() == TokenInvalid {
		return "", env.lexerError()
	}

	// Node
	if !env.inFirstSetForNode() {
		return "", env.unexpectError()
	}
	env.walker.EnterObj(key)
	err := env.parseNode()
	env.walker.Exit()
	if err != nil {
		return "", err
	}

	if env.tokenType() == TokenComma {
		// ','
		env.getToken()
	}
	return key, nil
}

func (env *parser) parseList() error {
	// '['
	if DEBUG {
		env.assertType("List", TokenLeftSquare)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing List")
		}
	}
	if env.getToken() == TokenInvalid {
		return env.lexerError()
	}

	// Elements
	if !env.inFirstSetForNode() && env.tokenType() != TokenRightSquare {
		return env.unexpectError()
	}
	err := env.parseElements()
	if err != nil {
		return err
	}

	// ']'
	if env.tokenType() != TokenRightSquare {
		return env.unexpectError()
	}
	env.getToken()
	return nil
}

func (env *parser) parseElements() error {
	if DEBUG {
		env.assertType("Elements",
			TokenInt,
			TokenString,
			TokenBool,
			TokenFloat,
			TokenNull,
			TokenLeftBrace,
			TokenLeftSquare,
			TokenRightSquare)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing Elements")
		}
	}
	index := 0
	for env.inFirstSetForNode() {
		env.walker.EnterList(index)
		err := env.parseElement()
		env.walker.Exit()
		index++
		if err != nil {
			return err
		}
	}
	return nil
}

func (env *parser) parseElement() error {
	// Node
	if DEBUG {
		env.assertType("Element",
			TokenInt,
			TokenString,
			TokenBool,
			TokenFloat,
			TokenNull,
			TokenLeftBrace,
			TokenLeftSquare)
		if DEBUG_ENABLE_PARSER_LOG {
			fmt.Println("parsing Elements")
		}
	}
	err := env.parseNode()
	if err != nil {
		return err
	}

	if env.tokenType() == TokenComma {
		// ','
		env.getToken()
	}
	return nil
}
