package shared

func Contains[T comparable](elems []T, v T) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func MapList[T interface{}, DTO interface{}](entities []T, mapper func(*T) DTO) *[]DTO {
	dtos := make([]DTO, 0, len(entities))

	for _, v := range entities {
		dtos = append(dtos, mapper(&v))
	}

	return &dtos
}
