package service

import (
	"strings"

	"github.com/rescuebyte/rescuebyte-be/domain/enums"
)

func (s *enumService) GetEnums(types string) (map[string]any, error) {
	requested := make(map[string]bool)
	if types != "" {
		for _, t := range strings.Split(types, ",") {
			requested[strings.TrimSpace(t)] = true
		}
	}

	includeAll := len(requested) == 0
	result := make(map[string]any)

	for _, desc := range enums.GetEnumRegistry() {
		if includeAll || requested[desc.QueryKey] {
			result[desc.ResponseKey] = desc.Values
		}
	}

	return result, nil
}
