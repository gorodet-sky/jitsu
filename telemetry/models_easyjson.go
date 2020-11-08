// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package telemetry

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

func easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry(in *jlexer.Lexer, out *Usage) {
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
		case "server_start":
			out.ServerStart = int(in.Int())
		case "server_stop":
			out.ServerStop = int(in.Int())
		case "events":
			out.Events = uint64(in.Uint64())
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
func easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry(out *jwriter.Writer, in Usage) {
	out.RawByte('{')
	first := true
	_ = first
	if in.ServerStart != 0 {
		const prefix string = ",\"server_start\":"
		first = false
		out.RawString(prefix[1:])
		out.Int(int(in.ServerStart))
	}
	if in.ServerStop != 0 {
		const prefix string = ",\"server_stop\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int(int(in.ServerStop))
	}
	if in.Events != 0 {
		const prefix string = ",\"events\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Uint64(uint64(in.Events))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Usage) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Usage) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Usage) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Usage) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry(l, v)
}
func easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry1(in *jlexer.Lexer, out *Request) {
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
		case "timestamp":
			out.Timestamp = string(in.String())
		case "instance_info":
			if in.IsNull() {
				in.Skip()
				out.InstanceInfo = nil
			} else {
				if out.InstanceInfo == nil {
					out.InstanceInfo = new(InstanceInfo)
				}
				(*out.InstanceInfo).UnmarshalEasyJSON(in)
			}
		case "metric_type":
			out.MetricType = string(in.String())
		case "usage":
			if in.IsNull() {
				in.Skip()
				out.Usage = nil
			} else {
				if out.Usage == nil {
					out.Usage = new(Usage)
				}
				(*out.Usage).UnmarshalEasyJSON(in)
			}
		case "errors":
			if in.IsNull() {
				in.Skip()
				out.Errors = nil
			} else {
				if out.Errors == nil {
					out.Errors = new(Errors)
				}
				(*out.Errors).UnmarshalEasyJSON(in)
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
func easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry1(out *jwriter.Writer, in Request) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Timestamp != "" {
		const prefix string = ",\"timestamp\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Timestamp))
	}
	if in.InstanceInfo != nil {
		const prefix string = ",\"instance_info\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.InstanceInfo).MarshalEasyJSON(out)
	}
	if in.MetricType != "" {
		const prefix string = ",\"metric_type\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.MetricType))
	}
	if in.Usage != nil {
		const prefix string = ",\"usage\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Usage).MarshalEasyJSON(out)
	}
	if in.Errors != nil {
		const prefix string = ",\"errors\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		(*in.Errors).MarshalEasyJSON(out)
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Request) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Request) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Request) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Request) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry1(l, v)
}
func easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry2(in *jlexer.Lexer, out *InstanceInfo) {
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
		case "id":
			out.Id = string(in.String())
		case "commit":
			out.Commit = string(in.String())
		case "tag":
			out.Tag = string(in.String())
		case "built_at":
			out.BuiltAt = string(in.String())
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
func easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry2(out *jwriter.Writer, in InstanceInfo) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Id != "" {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.String(string(in.Id))
	}
	if in.Commit != "" {
		const prefix string = ",\"commit\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Commit))
	}
	if in.Tag != "" {
		const prefix string = ",\"tag\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.Tag))
	}
	if in.BuiltAt != "" {
		const prefix string = ",\"built_at\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.String(string(in.BuiltAt))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v InstanceInfo) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry2(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v InstanceInfo) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry2(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *InstanceInfo) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry2(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *InstanceInfo) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry2(l, v)
}
func easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry3(in *jlexer.Lexer, out *Errors) {
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
		case "id":
			out.Id = int64(in.Int64())
		case "quantity":
			out.Quantity = int64(in.Int64())
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
func easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry3(out *jwriter.Writer, in Errors) {
	out.RawByte('{')
	first := true
	_ = first
	if in.Id != 0 {
		const prefix string = ",\"id\":"
		first = false
		out.RawString(prefix[1:])
		out.Int64(int64(in.Id))
	}
	if in.Quantity != 0 {
		const prefix string = ",\"quantity\":"
		if first {
			first = false
			out.RawString(prefix[1:])
		} else {
			out.RawString(prefix)
		}
		out.Int64(int64(in.Quantity))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v Errors) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry3(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v Errors) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonD2b7633eEncodeGithubComJitsucomEventnativeTelemetry3(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *Errors) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry3(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *Errors) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonD2b7633eDecodeGithubComJitsucomEventnativeTelemetry3(l, v)
}
