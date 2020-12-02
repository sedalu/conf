package json

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/ardanlabs/conf"
)

type JSONSourcer struct {
	m map[string]string
}

func (s *JSONSourcer) Source(fld conf.Field) (string, bool) {
	if fld.Options.ShortFlagChar != 0 {
		flagKey := fld.Options.ShortFlagChar
		k := strings.ToLower(string(flagKey))
		if val, found := s.m[k]; found {
			return val, found
		}
	}

	k := strings.ToLower(strings.Join(fld.FlagKey, `_`))
	val, found := s.m[k]
	return val, found
}

// NewSource returns a conf.Sourcer and, potentially, an error if a
// read error occurs or the Reader contains an invalid JSON document.
func NewSource(r io.Reader) (conf.Sourcer, error) {
	if r == nil {
		return &JSONSourcer{m: nil}, nil
	}

	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	tmpMap := make(map[string]interface{})
	err = json.Unmarshal(src, &tmpMap)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for key, value := range tmpMap {
		switch v := value.(type) {
		case float64:
			m[key] = strings.TrimRight(fmt.Sprintf("%f", v), "0.")
		case bool:
			m[key] = fmt.Sprintf("%t", v)
		case string:
			m[key] = value.(string)
		}
	}

	return &JSONSourcer{m: m}, nil
}
