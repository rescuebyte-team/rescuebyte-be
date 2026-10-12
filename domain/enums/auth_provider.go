package enums

type AuthProvider string

const (
	ProviderLocal  AuthProvider = "LOCAL"
	ProviderGoogle AuthProvider = "GOOGLE"
)

func (p AuthProvider) String() string {
	return string(p)
}
