/*
 * This file was last modified at 2024-04-20 17:09 by Victor N. Skurikhin.
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

// UnmarshalJSON supports json.Unmarshaler interface
//func (f Float) UnmarshalJSON(data []byte) error {
//	if data[0] == 34 {
//		err := json.Unmarshal(data[1:len(data)-1], &f)
//		if err != nil {
//			return errors.New("Float: UnmarshalJSON: " + err.Error())
//		}
//	} else {
//		err := json.Unmarshal(data, &f)
//		if err != nil {
//			return errors.New("Float: UnmarshalJSON: " + err.Error())
//		}
//	}
//	fmt.Fprintf(os.Stderr, "f: %f\n", f)
//	return nil
//}
