package xml

import (
	"bytes"
	"encoding/json"

	xj "github.com/basgys/goxml2json"
	gojsonq "github.com/thedevsaddam/gojsonq/v2"
)

type xmlDecoder struct {
}

func (i *xmlDecoder) Decode(data []byte, v interface{}) error {
	buf, err := xj.Convert(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), &v)
}

type Query struct {
	jq *gojsonq.JSONQ
}

func New() *Query {
	return &Query{gojsonq.New(gojsonq.SetDecoder(&xmlDecoder{}))}
}

// Get returns the value at path.
func (q Query) Get(xml []byte, path string) (string, error) {
	val, err := q.jq.FromString(string(xml)).From(path).GetR()
	if err != nil {
		return "", err
	}
	return val.String()
}
