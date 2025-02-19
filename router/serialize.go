package router

import (
	"encoding/json"
	"io"
)

var _ Serializer = JSONSerializer{}

type Serializer interface {
	Marshal(w io.Writer, v any) error
	Unmarshal(data []byte, v any) error
	ContentType() string
}

func defaultSerializers() map[string]Serializer {
	return map[string]Serializer{
		contentTypeJson: JSONSerializer{},
	}
}

// -----------------

type JSONSerializer struct{}

func (JSONSerializer) Marshal(w io.Writer, v any) error {
	return json.NewEncoder(w).Encode(v)
}

func (JSONSerializer) Unmarshal(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func (JSONSerializer) ContentType() string {
	return contentTypeJson
}
