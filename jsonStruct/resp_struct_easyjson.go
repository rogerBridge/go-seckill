//out.Data: false// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jsonStruct

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

func easyjsonAc765e9cDecodeGoRedisJsonStruct(in *jlexer.Lexer, out *CommonResponse) {
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
		case "code":
			out.Code = int(in.Int())
		case "msg":
			out.Msg = string(in.String())
		case "data":
			if m, ok := out.Data.(easyjson.Unmarshaler); ok {
				m.UnmarshalEasyJSON(in)
			} else if m, ok := out.Data.(json.Unmarshaler); ok {
				_ = m.UnmarshalJSON(in.Raw())
			} else {
				out.Data = in.Interface()
			}
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
func easyjsonAc765e9cEncodeGoRedisJsonStruct(out *jwriter.Writer, in CommonResponse) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"code\":"
		out.RawString(prefix[1:])
		out.Int(int(in.Code))
	}
	{
		const prefix string = ",\"msg\":"
		out.RawString(prefix)
		out.String(string(in.Msg))
	}
	{
		const prefix string = ",\"data\":"
		out.RawString(prefix)
		if m, ok := in.Data.(easyjson.Marshaler); ok {
			m.MarshalEasyJSON(out)
		} else if m, ok := in.Data.(json.Marshaler); ok {
			out.Raw(m.MarshalJSON())
		} else {
			out.Raw(json.Marshal(in.Data))
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v CommonResponse) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonAc765e9cEncodeGoRedisJsonStruct(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v CommonResponse) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonAc765e9cEncodeGoRedisJsonStruct(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *CommonResponse) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonAc765e9cDecodeGoRedisJsonStruct(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *CommonResponse) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonAc765e9cDecodeGoRedisJsonStruct(l, v)
}