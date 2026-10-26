package contracts

type EnumService interface {
	GetEnums(types string) (map[string]any, error)
}
