package pagination

var entityLimits = map[string]uint64{
	"product": 10,
	"store":   10,
}

func Normalize(limit, page, offset *uint64, entity string) {
	if *limit < 1 {
		if l, ok := entityLimits[entity]; ok {
			*limit = l
		} else {
			*limit = 10
		}
	}

	if *page < 1 {
		*page = 1
	}

	*offset = (*page - 1) * *limit
}
