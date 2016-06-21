package validations

import (
	"github.com/asaskevich/govalidator"
	"github.com/jinzhu/gorm"
)

var skipValidations = "validations:skip_validations"

func validate(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		if result, ok := scope.DB().Get(skipValidations); !(ok && result.(bool)) {
			if !scope.HasError() {
				scope.CallMethod("Validate")
				_, err := govalidator.ValidateStruct(scope.IndirectValue().Interface())
				if err != nil {
					scope.DB().AddError(err)
				}
			}
		}
	}
}

// RegisterCallbacks register callback into GORM DB
func RegisterCallbacks(db *gorm.DB) {
	callback := db.Callback()
	callback.Create().Before("gorm:before_create").Register("validations:validate", validate)
	callback.Update().Before("gorm:before_update").Register("validations:validate", validate)
}
