package tree2json

import "strings"

type jsonWriter interface {
	Enter()
	Exit()
	StartComment()
	EndComment()

	CommentAndNewLine(comment string)

	/*
		NewLine create a new line.
	*/
	NewLine()

	/*
		EndLine create a new line.When calling EndLine just after another calling
		of EndLine or NewLine, no operation will be done.
	*/
	EndLine()

	WriteSpace()
	WriteString(str string)
	WriteRune(char rune)
}

type emptyType struct{}

var emptyValue = emptyType{}

type stringBuilderWriter struct {
	builder        strings.Builder
	level          int
	commentLevel   map[int]emptyType
	endLineNeedNew bool
}

func (s *stringBuilderWriter) Enter() {
	s.endLineNeedNew = true
	s.level++
}

func (s *stringBuilderWriter) Exit() {
	s.endLineNeedNew = true
	if DEBUG {
		if s.level == 0 {
			panic("Exit() without Enter()")
		}
	}
	s.level--
}

func (s *stringBuilderWriter) StartComment() {
	s.endLineNeedNew = true
	s.builder.WriteString("// ")
	s.commentLevel[s.level] = emptyValue
}

func (s *stringBuilderWriter) EndComment() {
	delete(s.commentLevel, s.level)
	s.EndLine()
}

func (s *stringBuilderWriter) CommentAndNewLine(comment string) {
	s.builder.WriteString("// ")
	s.builder.WriteString(comment)
	s.NewLine()
}

func (s *stringBuilderWriter) EndLine() {
	if !s.endLineNeedNew {
		return
	}
	s.NewLine()
}

func (s *stringBuilderWriter) NewLine() {
	s.endLineNeedNew = false
	s.builder.WriteRune('\n')
	for i := 0; i < s.level; i++ {
		s.builder.WriteRune('\t')
		_, ok := s.commentLevel[i+1]
		if ok {
			s.builder.WriteString("// ")
		}
	}
}

func (s *stringBuilderWriter) WriteSpace() {
	s.endLineNeedNew = true
	s.builder.WriteRune(' ')
}

func (s *stringBuilderWriter) WriteString(str string) {
	s.endLineNeedNew = true
	s.builder.WriteString(str)
}

func (s *stringBuilderWriter) WriteRune(char rune) {
	s.endLineNeedNew = true
	s.builder.WriteRune(char)
}
