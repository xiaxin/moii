package web

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	zh "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
	"sync"
)

var (
	validate *validator.Validate
	trans    ut.Translator
	uni      *ut.UniversalTranslator
)

func init() {
	zh := zh.New()
	uni := ut.New(zh, zh)
	validate = validator.New()
	// TODO 这部分应放到中间件中
	trans, _ = uni.GetTranslator("zh")
	zh_translations.RegisterDefaultTranslations(validate, trans)
}

type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var _ binding.StructValidator = &DefaultValidator{}

// ValidateStruct 如果接收到的类型是一个结构体或指向结构体的指针，则执行验证。
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {
		v.lazyinit()

		// 如果传递不合规则的值，则返回InvalidValidationError，否则返回nil。
		// 如果返回err != nil，可通过err.(validator.ValidationErrors)来访问错误数组。
		if err := validate.Struct(obj); err != nil {
			return err
		}
	}
	return nil
}
// Engine 返回支持`StructValidator`实现的底层验证引擎
func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validate
	})
}

func kindOfData(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func Validate(ctx *gin.Context, obj interface{}) map[string]string {
	resp := make(map[string]string, 0)

	if err := ctx.ShouldBind(obj); nil != err {
		tranFn, bool := obj.(ValidateFieldTranslates)
		m := make(map[string]string, 0)

		if bool {
			m = tranFn.FieldTranslate()
		}

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			transtr := e.Translate(trans)
			field := e.Field()

			//判断错误字段是否在命名中，如果在，则替换错误信息中的字段
			if rp, ok := m[e.Field()]; ok {
				resp[field] = strings.Replace(transtr, e.Field(), rp, 1)
			} else {
				resp[field] = transtr
			}
		}

		return resp
	}
	return nil
}


type ValidateFieldTranslates interface{
	FieldTranslate() map[string]string
}
