package pgerror

type Code string

const (
	UniqueViolation Code = "23505"
	ForeignKey      Code = "23503"
)

func (e Code) String() string {
	return string(e)
}
