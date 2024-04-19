/*
 * This file was last modified at 2024-04-19 16:13 by Victor N. Skurikhin.
 * float.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson/jwriter"
	"strconv"
)

type Float float64

func (f *Float) MarshalEasyJSON(w *jwriter.Writer) {
	w.RawString(strconv.FormatFloat(float64(*f), 'f', -1, 64))
}
