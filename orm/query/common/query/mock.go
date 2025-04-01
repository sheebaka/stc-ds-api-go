package query

import (
	"github.com/stc-ds-databricks-go/orm/query/common"
	"gorm.io/gorm"
)

type SfPlaceHolder struct {
	SfAccount common.Filter
}

func Use(*gorm.DB) (p SfPlaceHolder) {
	return
}
