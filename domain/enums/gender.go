package enums

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"math"
	"strconv"
	"strings"
)

type GenderIdx int16
type GenderKey string

type GenderValue struct {
	Idx        int16           `json:"idx"`
	Key        string          `json:"key"`
	LongLabel  LocalizedString `json:"long_label"`
	ShortLabel LocalizedString `json:"short_label"`
}

const (
	GenderUnknownIdx GenderIdx = 0
	GenderMaleIdx    GenderIdx = 1
	GenderFemaleIdx  GenderIdx = 2

	GenderUnknownKey GenderKey = "unknown"
	GenderMaleKey    GenderKey = "male"
	GenderFemaleKey  GenderKey = "female"
)

var (
	GenderUnknownValue = GenderValue{
		Idx:        int16(GenderUnknownIdx),
		Key:        string(GenderUnknownKey),
		LongLabel:  LocalizedString{"id": "Tidak Diketahui", "en": "Unknown"},
		ShortLabel: LocalizedString{"id": "Tidak Diketahui", "en": "Unknown"},
	}
	GenderMaleValue = GenderValue{
		Idx:        int16(GenderMaleIdx),
		Key:        string(GenderMaleKey),
		LongLabel:  LocalizedString{"id": "Laki-laki", "en": "Male"},
		ShortLabel: LocalizedString{"id": "Laki-laki", "en": "Male"},
	}
	GenderFemaleValue = GenderValue{
		Idx:        int16(GenderFemaleIdx),
		Key:        string(GenderFemaleKey),
		LongLabel:  LocalizedString{"id": "Perempuan", "en": "Female"},
		ShortLabel: LocalizedString{"id": "Perempuan", "en": "Female"},
	}
)

var (
	GenderMapIdx = map[GenderIdx]GenderValue{
		GenderUnknownIdx: GenderUnknownValue,
		GenderMaleIdx:    GenderMaleValue,
		GenderFemaleIdx:  GenderFemaleValue,
	}

	GenderMapKey = map[GenderKey]GenderValue{
		GenderUnknownKey: GenderUnknownValue,
		GenderMaleKey:    GenderMaleValue,
		GenderFemaleKey:  GenderFemaleValue,
	}

	GenderMapReverse = map[GenderKey]GenderIdx{
		GenderUnknownKey: GenderUnknownIdx,
		GenderMaleKey:    GenderMaleIdx,
		GenderFemaleKey:  GenderFemaleIdx,
	}
)

func (d GenderIdx) String() string {
	if v, ok := GenderMapIdx[d]; ok {
		return v.Key
	}
	return string(GenderUnknownKey)
}

func (d *GenderIdx) Scan(value any) error {
	if value == nil {
		*d = GenderUnknownIdx
		return nil
	}

	switch v := value.(type) {
	case int8:
		*d = GenderIdx(v)
		return nil
	case int16:
		*d = GenderIdx(v)
		return nil
	case int32:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return errors.New("invalid gender value: int32 out of int16 range")
		}
		*d = GenderIdx(v)
		return nil
	case int64:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return errors.New("invalid gender value: int64 out of int16 range")
		}
		*d = GenderIdx(v)
		return nil
	case string:
		if idx, ok := GenderMapReverse[GenderKey(strings.ToLower(v))]; ok {
			*d = idx
			return nil
		}
	case []byte:
		strVal := strings.ToLower(string(v))
		if idx, ok := GenderMapReverse[GenderKey(strVal)]; ok {
			*d = idx
			return nil
		}
	}

	return errors.New("invalid gender value")
}

func (d GenderIdx) Value() (driver.Value, error) {
	return int64(d), nil
}

type NullGenderIdx struct {
	GenderIdx GenderIdx
	Valid     bool
}

func (d *NullGenderIdx) Scan(value any) error {
	if value == nil {
		d.GenderIdx, d.Valid = GenderUnknownIdx, false
		return nil
	}

	d.Valid = true
	return d.GenderIdx.Scan(value)
}

func (d NullGenderIdx) Value() (driver.Value, error) {
	if !d.Valid {
		return driver.Value(nil), nil
	}

	return d.GenderIdx.Value()
}

func (d GenderIdx) IsValidIdx() bool {
	switch d {
	case GenderMaleIdx, GenderFemaleIdx:
		return true
	default:
		return false
	}
}

func (d *GenderIdx) UnmarshalJSON(b []byte) error {
	str := strings.TrimSpace(string(b))
	if str == "null" || str == `""` {
		*d = GenderUnknownIdx
		return nil
	}

	if val, err := strconv.ParseInt(str, 10, 16); err == nil {
		*d = GenderIdx(val)
		return nil
	}

	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		if idx, ok := GenderMapReverse[GenderKey(strings.ToLower(s))]; ok {
			*d = idx
			return nil
		}
	}

	return errors.New("invalid gender value")
}

func (d GenderIdx) MarshalJSON() ([]byte, error) {
	return json.Marshal(int16(d))
}

func init() {
	RegisterEnum("gender", "genders", GenderMapIdx)
}
