/*
 * This file was last modified at 2024-05-07 16:32 by Victor N. Skurikhin.
 * error_easyjson.go
 * $Id$
 */

// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package model

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonE34310f8DecodeGithubComVskurikhinGophermartInternalModel(in *jlexer.Lexer, out *JSONError) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "error":
			out.Error = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonE34310f8EncodeGithubComVskurikhinGophermartInternalModel(out *jwriter.Writer, in JSONError) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"error\":"
		out.RawString(prefix[1:])
		out.String(string(in.Error))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v JSONError) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonE34310f8EncodeGithubComVskurikhinGophermartInternalModel(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v JSONError) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonE34310f8EncodeGithubComVskurikhinGophermartInternalModel(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *JSONError) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonE34310f8DecodeGithubComVskurikhinGophermartInternalModel(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *JSONError) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonE34310f8DecodeGithubComVskurikhinGophermartInternalModel(l, v)
}
