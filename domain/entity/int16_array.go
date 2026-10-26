package entity

import (
	"database/sql/driver"
	"fmt"
	"strconv"
	"strings"
)

// Int16Array maps a PostgreSQL SMALLINT[] column to a Go []int16.
type Int16Array []int16

func (a *Int16Array) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*a = Int16Array{}
		return nil
	case string:
		return a.parse(v)
	case []byte:
		return a.parse(string(v))
	default:
		return fmt.Errorf("unsupported Scan source for Int16Array: %T", src)
	}
}

func (a *Int16Array) parse(s string) error {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")
	if s == "" {
		*a = Int16Array{}
		return nil
	}

	parts := strings.Split(s, ",")
	res := make(Int16Array, len(parts))
	for i, p := range parts {
		n, err := strconv.ParseInt(strings.TrimSpace(p), 10, 16)
		if err != nil {
			return fmt.Errorf("invalid Int16Array element %q: %w", p, err)
		}
		res[i] = int16(n)
	}
	*a = res
	return nil
}

func (a Int16Array) Value() (driver.Value, error) {
	if a == nil {
		return "{}", nil
	}

	var sb strings.Builder
	sb.WriteByte('{')
	for i, n := range a {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.FormatInt(int64(n), 10))
	}
	sb.WriteByte('}')
	return sb.String(), nil
}
