package core

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var Trans ut.Translator

func InitializeTrans() (err error) {
	//配置gin以支持中文
	// Accept-Language
	// 修改gin框架validator引擎属性
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := fld.Tag.Get("json")
			return name
		})
		zhT := zh.New()
		uni := ut.New(zhT, zhT)
		Trans, _ = uni.GetTranslator("zh")
		err = zhTranslations.RegisterDefaultTranslations(v, Trans)
		return
	}
	return
}

//GetFirstValidateError 验证错误
func GetFirstValidateError(err error) (errStr string) {
	if err != nil {
		if errors, ok := err.(validator.ValidationErrors); ok {
			trans := errors.Translate(Trans)
			for _, err1 := range trans {
				errStr = err1
				return
			}
		}
		return err.Error()
	}
	return
}

func removeTopStruct(fields validator.ValidationErrorsTranslations) validator.ValidationErrorsTranslations {
	r := make(validator.ValidationErrorsTranslations)
	for f, v := range fields {
		r[f[strings.Index(f, ".")+1:]] = v
	}
	return r
}
