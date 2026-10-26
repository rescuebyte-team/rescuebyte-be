package enums

import (
	"database/sql/driver"
	"errors"
	"math"
)

type UserRoleIdx int16
type UserRoleKey string
type UserRoleValue struct {
	Idx        int16           `json:"idx"`
	Key        string          `json:"key"`
	LongLabel  LocalizedString `json:"long_label"`
	ShortLabel LocalizedString `json:"short_label"`
}

const (
	RoleUnknownIdx UserRoleIdx = 0
	RoleUserIdx    UserRoleIdx = 1
	RoleAdminIdx   UserRoleIdx = 2
	RoleSellerIdx  UserRoleIdx = 3

	RoleUnknownKey UserRoleKey = "unknown"
	RoleUserKey    UserRoleKey = "user"
	RoleAdminKey   UserRoleKey = "admin"
	RoleSellerKey  UserRoleKey = "seller"
)

var (
	RoleUnknownValue = UserRoleValue{
		Idx:        int16(RoleUnknownIdx),
		Key:        string(RoleUnknownKey),
		LongLabel:  LocalizedString{"id": "", "en": ""},
		ShortLabel: LocalizedString{"id": "", "en": ""},
	}
	RoleUserValue = UserRoleValue{
		Idx:        int16(RoleUserIdx),
		Key:        string(RoleUserKey),
		LongLabel:  LocalizedString{"id": "Pengguna", "en": "User"},
		ShortLabel: LocalizedString{"id": "User", "en": "User"},
	}
	RoleAdminValue = UserRoleValue{
		Idx:        int16(RoleAdminIdx),
		Key:        string(RoleAdminKey),
		LongLabel:  LocalizedString{"id": "Administrator", "en": "Administrator"},
		ShortLabel: LocalizedString{"id": "Admin", "en": "Admin"},
	}
	RoleSellerValue = UserRoleValue{
		Idx:        int16(RoleSellerIdx),
		Key:        string(RoleSellerKey),
		LongLabel:  LocalizedString{"id": "Penjual", "en": "Seller"},
		ShortLabel: LocalizedString{"id": "Seller", "en": "Seller"},
	}
)

var (
	UserRoleMapIdx = map[UserRoleIdx]UserRoleValue{
		RoleUnknownIdx: RoleUnknownValue,
		RoleUserIdx:    RoleUserValue,
		RoleAdminIdx:   RoleAdminValue,
		RoleSellerIdx:  RoleSellerValue,
	}

	UserRoleMapKey = map[UserRoleKey]UserRoleValue{
		RoleUnknownKey: RoleUnknownValue,
		RoleUserKey:    RoleUserValue,
		RoleAdminKey:   RoleAdminValue,
		RoleSellerKey:  RoleSellerValue,
	}

	UserRoleMapReverse = map[UserRoleKey]UserRoleIdx{
		RoleUnknownKey: RoleUnknownIdx,
		RoleUserKey:    RoleUserIdx,
		RoleAdminKey:   RoleAdminIdx,
		RoleSellerKey:  RoleSellerIdx,
	}
)

func (d UserRoleIdx) String() string {
	if role, ok := UserRoleMapIdx[d]; ok {
		return role.Key
	}
	return string(RoleUnknownKey)
}

func (d *UserRoleIdx) Scan(value any) error {
	if value == nil {
		*d = RoleUnknownIdx
		return nil
	}

	switch v := value.(type) {
	case int8:
		*d = UserRoleIdx(v)
		return nil
	case int16:
		*d = UserRoleIdx(v)
		return nil
	case int32:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return errors.New("invalid user role value: int32 out of int16 range")
		}
		*d = UserRoleIdx(v)
		return nil
	case int64:
		if v < math.MinInt16 || v > math.MaxInt16 {
			return errors.New("invalid user role value: int64 out of int16 range")
		}
		*d = UserRoleIdx(v)
		return nil
	case string:
		if idx, ok := UserRoleMapReverse[UserRoleKey(v)]; ok {
			*d = idx
			return nil
		}
	case []byte:
		strVal := string(v)
		if idx, ok := UserRoleMapReverse[UserRoleKey(strVal)]; ok {
			*d = idx
			return nil
		}
	}

	return errors.New("invalid user role value")
}

func (d UserRoleIdx) Value() (driver.Value, error) {
	return int64(d), nil
}

type NullUserRoleIdx struct {
	UserRoleIdx UserRoleIdx
	Valid       bool
}

func (d *NullUserRoleIdx) Scan(value any) error {
	if value == nil {
		d.UserRoleIdx, d.Valid = RoleUnknownIdx, false
		return nil
	}

	d.Valid = true
	return d.UserRoleIdx.Scan(value)
}

func (d NullUserRoleIdx) Value() (driver.Value, error) {
	if !d.Valid {
		return driver.Value(nil), nil
	}

	return d.UserRoleIdx.Value()
}

func (d UserRoleIdx) IsValidIdx() bool {
	switch d {
	case RoleUnknownIdx:
		return false
	case RoleUserIdx, RoleAdminIdx, RoleSellerIdx:
		return true
	default:
		return false
	}
}

func init() {
	RegisterEnum("user_role", "user_roles", UserRoleMapIdx)
}
