/*
 * This file was last modified at 2024-04-20 00:12 by Victor N. Skurikhin.
 * time.go
 * $Id$
 */

package model

import (
	"github.com/mailru/easyjson/jwriter"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) MarshalEasyJSON(w *jwriter.Writer) {
	w.RawString(t.Time.Local().Round(time.Second).Format(time.RFC3339))
}

func (t *Time) MarshalJSON() ([]byte, error) {
	return []byte(t.Time.Local().Round(time.Second).Format(time.RFC3339)), nil
}
