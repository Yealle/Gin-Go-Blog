package app

import (
	"gin-blog/pkg/logging"

	"github.com/astaxie/beego/validation"
)

func MarkErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info("err.key: %s, err.message: %s", err.Key, err.Message)
	}

	return
}
