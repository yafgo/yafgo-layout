package validators

import (
	"fmt"

	"github.com/gin-gonic/gin/binding"
	localeEn "github.com/go-playground/locales/en"
	localeZh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transEn "github.com/go-playground/validator/v10/translations/en"
	transZh "github.com/go-playground/validator/v10/translations/zh"
)

var trans ut.Translator

// initTrans
//
//	locale 通常取决于 http 请求头的 'Accept-Language'
func initTrans(local string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		zhT := localeZh.New() //chinese
		enT := localeEn.New() //english
		uni := ut.New(enT, zhT, enT)

		var found bool
		trans, found = uni.GetTranslator(local)
		if !found {
			return fmt.Errorf("uni.GetTranslator(%s) failed", local)
		}
		//register translate
		switch local {
		case "zh":
			err = transZh.RegisterDefaultTranslations(v, trans)
		case "en":
			fallthrough
		default:
			err = transEn.RegisterDefaultTranslations(v, trans)
		}
		return
	}
	return
}

func TranslateErrors(err error) (msgs map[string]string) {
	if trans == nil {
		initTrans("zh")
	}
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		// 非validator.ValidationErrors类型错误直接返回
		return nil
	}
	// validator.ValidationErrors类型错误则进行翻译
	return errs.Translate(trans)
}
