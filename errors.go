package kvql

import (
	"fmt"
	"strings"
)

var (
	_ error       = (*SyntaxError)(nil)
	_ error       = (*ExecuteError)(nil)
	_ QueryBinder = (*SyntaxError)(nil)
	_ QueryBinder = (*ExecuteError)(nil)

	DefaultErrorPadding = 7
)

type QueryBinder interface {
	BindQuery(query string)
	SetPadding(pad int)
}

type SyntaxError struct {
	Query   string
	Message string
	Pos     int
	Padding int
}

func NewSyntaxError(pos int, msg string, args ...any) error {
	return &SyntaxError{
		Message: fmt.Sprintf(msg, args...),
		Pos:     pos,
		Padding: DefaultErrorPadding,
	}
}

func (e *SyntaxError) BindQuery(query string) {
	e.Query = query
}

func (e *SyntaxError) SetPadding(pad int) {
	e.Padding = pad
}

func (e *SyntaxError) Error() string {
	if e.Query == "" {
		return e.simpleError()
	}
	return e.queryError()
}

func (e *SyntaxError) queryError() string {
	ret := outputQueryAndErrPos(e.Query, e.Pos, e.Padding)
	pad := generatePads(e.Padding)
	ret += fmt.Sprintf("%sSyntax Error: %s", pad, e.Message)
	return ret
}

func (e *SyntaxError) simpleError() string {
	return fmt.Sprintf("Syntax Error: %s at %d", e.Message, e.Pos)
}

type ExecuteError struct {
	Query   string
	Message string
	Pos     int
	Padding int
}

func NewExecuteError(pos int, msg string, args ...any) error {
	return &ExecuteError{
		Pos:     pos,
		Message: fmt.Sprintf(msg, args...),
		Padding: DefaultErrorPadding,
	}
}

func (e *ExecuteError) BindQuery(query string) {
	e.Query = query
}

func (e *ExecuteError) SetPadding(pad int) {
	e.Padding = pad
}

func (e *ExecuteError) Error() string {
	if e.Query == "" {
		return e.simpleError()
	}
	return e.queryError()
}

func (e *ExecuteError) simpleError() string {
	return fmt.Sprintf("Execute Error: %s at %d", e.Message, e.Pos)
}

func (e *ExecuteError) queryError() string {
	ret := outputQueryAndErrPos(e.Query, e.Pos, e.Padding)
	pad := generatePads(e.Padding)
	ret += fmt.Sprintf("%sExecute Error: %s", pad, e.Message)
	return ret
}

func generatePads(pad int) string {
	ret := ""
	for i := 0; i < pad; i++ {
		ret += " "
	}
	return ret
}

func outputQueryAndErrPos(query string, pos int, adjust int) string {
	tquery := strings.TrimSpace(query)
	qlen := len(tquery)
	if pos == -1 {
		pos = qlen
	}
	trimLeft := false
	trimRight := false
	if qlen > 70 {
		if pos <= 35 {
			tquery = tquery[0:70]
			trimRight = true
		} else {
			trimLeft = true
			trim := pos - 35
			restLen := qlen - trim
			if restLen > 70 {
				restLen = 70
				trimRight = true
			}
			tquery = tquery[trim : trim+restLen]
			pos -= trim
		}
	}
	ret := ""
	errPos := pos + adjust
	if trimLeft {
		ret = "... "
		errPos += 4
	}
	ret += tquery
	if trimRight {
		ret += " ..."
	}
	ret += "\n"
	for i := 0; i < errPos; i++ {
		ret += " "
	}
	ret += "^--\n"
	return ret
}
