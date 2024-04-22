/*
 * This file was last modified at 2024-04-21 15:31 by Victor N. Skurikhin.
 * float.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"strconv"
)

type Float float64

func (f *Float) MarshalEasyJSON(w *jwriter.Writer) {
	w.RawString(strconv.FormatFloat(float64(*f), 'f', -1, 64))
}
func (f *Float) UnmarshalEasyJSON(l *jlexer.Lexer) {
	fl := l.Float64()
	*f = Float(fl)
}
