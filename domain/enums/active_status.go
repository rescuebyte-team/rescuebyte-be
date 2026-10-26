package enums

import (
	"database/sql/driver"
	"errors"
	"math"
)

type ActiveStatusIdx int16
type ActiveStatusKey string

const (
	ActiveStatusInactiveIdx ActiveStatusIdx = 0
	ActiveStatusActiveIdx   ActiveStatusIdx = 1

	ActiveStatusInactiveKey ActiveStatusKey = "inactive"
	ActiveStatusActiveKey   ActiveStatusKey = "active"
)

type ActiveStatusValue struct {
	Idx        int16           `json:"idx"`
	Key        string          `json:"key"`
	LongLabel  LocalizedString `json:"long_label"`
	ShortLabel LocalizedString `json:"short_label"`
}

var (
	ActiveStatusInactiveValue = ActiveStatusValue{
		Idx:        int16(ActiveStatusInactiveIdx),
		Key:        string(ActiveStatusInactiveKey),
		LongLabel:  LocalizedString{"id": "Tidak Aktif", "en": "Inactive"},
		ShortLabel: LocalizedString{"id": "Tidak Aktif", "en": "Inactive"},
	}
	ActiveStatusActiveValue = ActiveStatusValue{
		Idx:        int16(ActiveStatusActiveIdx),
		Key:        string(ActiveStatusActiveKey),
		LongLabel:  LocalizedString{"id": "Aktif", "en": "Active"},
		ShortLabel: LocalizedString{"id": "Aktif", "en": "Active"},
	}
)

var ActiveStatusMapIdx = map[ActiveStatusIdx]ActiveStatusValue{
	ActiveStatusInactiveIdx: ActiveStatusInactiveValue,
	ActiveStatusActiveIdx:   ActiveStatusActiveValue,
}

var ActiveStatusMapKey = map[ActiveStatusKey]ActiveStatusValue{
	ActiveStatusInactiveKey: ActiveStatusInactiveValue,
	ActiveStatusActiveKey:   ActiveStatusActiveValue,
}

var ActiveStatusMapReverse = map[ActiveStatusKey]ActiveStatusIdx{
	ActiveStatusInactiveKey: ActiveStatusInactiveIdx,
	ActiveStatusActiveKey:   ActiveStatusActiveIdx,
}

func (d ActiveStatusIdx) String() string {
	if v, ok := ActiveStatusMapIdx[d]; ok {
		return v.Key
	}
	return string(ActiveStatusInactiveKey)
}

func (d ActiveStatusIdx) IsValidIdx() bool {
	switch d {
	case ActiveStatusInactiveIdx, ActiveStatusActiveIdx:
		return true
	default:
		return false
	}
}

func (d ActiveStatusIdx) Toggle() ActiveStatusIdx {
	if d == ActiveStatusActiveIdx {
		return ActiveStatusInactiveIdx
	}
	return ActiveStatusActiveIdx
}

func (d *ActiveStatusIdx) Scan(value any) error {
	if value == nil {
		*d = ActiveStatusInactiveIdx
		return nil
	}

	switch v := value.(type) {
	case bool:
		if v {
			*d = ActiveStatusActiveIdx
		} else {
			*d = ActiveStatusInactiveIdx
		}
		return nil
	case int8:
		*d = ActiveStatusIdx(v)
		return nil
	case int16:
		*d = ActiveStatusIdx(v)
		return nil
	case int32:
		if v > math.MaxInt16 || v < math.MinInt16 {
			return errors.New("integer overflow")
		}
		*d = ActiveStatusIdx(v)
		return nil
	case int64:
		if v > math.MaxInt16 || v < math.MinInt16 {
			return errors.New("integer overflow")
		}
		*d = ActiveStatusIdx(v)
		return nil
	}

	return errors.New("invalid active status value")
}

func (d ActiveStatusIdx) Value() (driver.Value, error) {
	return int16(d), nil
}

func init() {
	RegisterEnum("active_status", "active_statuses", ActiveStatusMapIdx)
}
