package validate

import (
	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v4"
)

type Validator struct {
	v *validator.Validate
}

func NewValidator() (echo.Validator, error) {
	v := validator.New()

	zh := zh_Hans_CN.New()
	uni := ut.New(zh, zh)

	trans, _ := uni.GetTranslator("zh")

	if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	validator := &Validator{v: v}

	return validator, nil
}

func (v *Validator) Validate(i any) error {
	return v.v.Struct(i)
}
