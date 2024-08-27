package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"net/http"
	"reflect"
	"strings"
)

type LoginForm struct {
	Id     uint64 `form:"id" json:"id" binding:"required"`
	Name   string `form:"name" json:"name" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required"`
}

/*
curl -X POST http://localhost:8081/login -d "id=123" -d "name=example_name"
*/
func login(c *gin.Context) {
	var loginForm LoginForm
	if err := c.ShouldBind(&loginForm); err != nil {
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "翻译失败"})
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": removeTopStruct(errs.Translate(trans))})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":     loginForm.Id,
		"name":   loginForm.Name,
		"passwd": loginForm.Passwd,
	})
}

var trans ut.Translator

func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}
	for field, err := range fileds {
		rsp[field[strings.Index(field, ".")+1:]] = err
	}
	return rsp
}

func InitTrans(locale string) (err error) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册一个获取json的tag的自定义方法
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})

		zhTr := zh.New()
		enTr := en.New()
		// 第一个是备用
		uni := ut.New(enTr, zhTr)
		trans, ok = uni.GetTranslator(locale)
		if !ok {
			return fmt.Errorf("translator not found(%s)", locale)
		}
		switch locale {
		case "zh":
			zh_translations.RegisterDefaultTranslations(v, trans)
		case "en":

			en_translations.RegisterDefaultTranslations(v, trans)
		default:
			en_translations.RegisterDefaultTranslations(v, trans)
		}
		return nil
	}
	return nil
}

func main() {
	err := InitTrans("zh")
	if err != nil {
		fmt.Println(err)
		return
	}
	r := gin.Default()
	r.POST("/login", login)
	r.Run(":8081") // 监听并在 0.0.0.0:8080 上启动服务
}
