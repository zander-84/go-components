package CValidate

import (
	"encoding/json"
	"errors"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/go-playground/validator/v10"
	"log"
	"reflect"
)

// 兼容gin
type Validate struct {
	obj   *validator.Validate
	trans ut.Translator
	uni   *ut.UniversalTranslator
	conf  Conf
}

//
func NewValidate(opts ...func(interface{})) *Validate {
	_validator := &Validate{
		obj: validator.New(),
	}

	for _, opt := range opts {
		opt(_validator)
	}

	_validator.build()
	return _validator
}

func BuildValidate(opts ...func(interface{})) interface{} {
	return NewValidate(opts...)
}
func SetConfig(conf Conf) func(interface{}) {
	return func(i interface{}) {
		g := i.(*Validate)
		g.conf = conf
		g.conf.SetDefault()
	}
}
func (v *Validate) build() {
	locales := map[string]bool{
		"zh": true,
		"en": true,
	}
	if _, ok := locales[v.conf.Locale]; !ok {
		log.Fatalln("locale err")
	}
	v.obj = validator.New()
	v.obj.SetTagName(v.conf.ValidateTag) //binding
	v.obj.RegisterTagNameFunc(func(fld reflect.StructField) string {
		return fld.Tag.Get(v.conf.CommentTag)
	})

	v.uni = ut.New(en.New(), zh.New())
	v.SetTranslator(v.conf.Locale)
}

// ValidateStruct receives any kind of type, but only performed struct or pointer to struct type.
func (v *Validate) ValidateStruct(obj interface{}) error {
	if kindOfData(obj) == reflect.Struct {
		if err := v.obj.Struct(obj); err != nil {
			errMap := make(map[string]interface{}, 0)
			source := structTagMapJson(obj, v.conf.CommentTag, v.conf.JsonTag)
			for _, err := range err.(validator.ValidationErrors) {
				if key, ok := source[err.Field()]; ok {
					errMap[key] = err.Translate(v.trans)
				} else {
					errMap[err.Field()] = err.Translate(v.trans)
				}
			}

			if len(errMap) > 0 {
				mjson, _ := json.Marshal(errMap)
				return errors.New(string(mjson))
			}

			return err
		}
	}
	return nil
}

// Engine returns the underlying validator engine which powers the default
// Validator instance. This is useful if you want to register custom validations
// or struct level validations. See validator GoDoc for more info -
// https://godoc.org/gopkg.in/go-playground/validator.v8
func (v *Validate) Engine() interface{} {
	return v.obj
}

func (v *Validate) lazyinit() {

}

func (v *Validate) SetTranslator(locale string) {
	switch locale {
	case "en":
		v.trans, _ = v.uni.GetTranslator(locale) //zh en
		_ = en_translations.RegisterDefaultTranslations(v.obj, v.trans)
	default:
		v.trans, _ = v.uni.GetTranslator("zh") //zh en
		_ = zh_translations.RegisterDefaultTranslations(v.obj, v.trans)
	}
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func structTagMapJson(data interface{}, comment string, json string) map[string]string {
	res := make(map[string]string)
	var t reflect.Type
	if reflect.TypeOf(data).Kind() == reflect.Ptr {
		t = reflect.TypeOf(data).Elem()
	} else {
		t = reflect.TypeOf(data)
	}

	for i := 0; i < t.NumField(); i++ {
		res[t.Field(i).Tag.Get(comment)] = t.Field(i).Tag.Get(json) //comment:json
	}
	return res
}
