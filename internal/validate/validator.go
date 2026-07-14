package validate

import (
	"errors"
	"strings"

	"github.com/go-playground/locales/zh_Hans_CN"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/labstack/echo/v5"
)

type Validator struct {
	v     *validator.Validate
	trans ut.Translator
}

func NewValidator() (echo.Validator, error) {
	v := validator.New()

	zh := zh_Hans_CN.New()
	uni := ut.New(zh, zh)

	trans, _ := uni.GetTranslator("zh")

	if err := zhTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, err
	}

	validator := &Validator{
		v:     v,
		trans: trans, // 保存翻译器实例
	}

	return validator, nil
}

func (v *Validator) Validate(i any) error {
	err := v.v.Struct(i)
	if err == nil {
		return nil
	}

	if validationErrs, ok := errors.AsType[validator.ValidationErrors](err); ok {
		var errMsgs []string
		for _, e := range validationErrs {
			errMsgs = append(errMsgs, e.Translate(v.trans))
		}
		// 将所有的翻译后的错误信息拼接成一个字符串返回
		return errors.New(strings.Join(errMsgs, "; "))
	}

	// 如果是其他类型的错误（如传入了非结构体类型），直接返回
	return err
}
