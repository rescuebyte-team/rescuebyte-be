package enums

var DefaultLanguage = "en"

type LocalizedString map[string]string

func (l LocalizedString) String() string {
	val, ok := l[DefaultLanguage]
	if !ok || val == "" {
		return l["id"]
	}
	return val
}
