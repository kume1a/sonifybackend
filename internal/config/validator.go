package config

import (
	"github.com/asaskevich/govalidator"
	"github.com/google/uuid"
)

func ConfigureGoValidator() {
	govalidator.SetFieldsRequiredByDefault(true)

	// after
	govalidator.CustomTypeTagMap.Set("uuidSliceNotEmpty", func(i interface{}, o interface{}) bool {
		slice, ok := i.([]uuid.UUID)
		if !ok {
			return false
		}
		return len(slice) > 0
	})
}
