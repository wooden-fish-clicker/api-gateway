package app

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func Valid(c *gin.Context, params interface{}, isQuery bool) (int, map[string]string) {
	var (
		errors = make(map[string]string)
		err    error
	)

	if isQuery {
		err = c.BindQuery(params)
	} else {
		err = c.Bind(params) //取body
	}

	if err != nil {
		return http.StatusInternalServerError, map[string]string{
			"驗證例外": err.Error(),
		}
	}

	valid := validation.Validation{}
	check, err := valid.Valid(params)
	if err != nil {
		return http.StatusInternalServerError, map[string]string{
			"驗證例外": err.Error(),
		}
	}
	if !check {
		for _, e := range valid.Errors {
			errors[e.Field] = e.Message
		}
		return http.StatusBadRequest, errors
	}

	return http.StatusOK, nil
}
