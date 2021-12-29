package xml

import (
	"bytes"
	"encoding/json"

	xj "github.com/basgys/goxml2json"
	"github.com/tomwright/dasel"
)

// convert XML to JSON, then unmarshal JSON
func decodeXML(data []byte, v interface{}) error {
	buf, err := xj.Convert(bytes.NewReader(data))
	if err != nil {
		return err
	}
	return json.Unmarshal(buf.Bytes(), &v)
}

// Get queries the unmarshaled XML data for the given path.
func Get(xml []byte, daselQuery string) (string, error) {
	var data interface{}
	err := decodeXML(xml, &data)
	rootNode := dasel.New(data)
	val, err := rootNode.Query(daselQuery)
	if err != nil {
		return "", err
	}
	return val.String(), nil
}
