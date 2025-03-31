package orm

import (
	"github.com/stc-ds-databricks-go/orm/model"
	"github.com/stc-ds-databricks-go/orm/query"
	"gorm.io/gorm"
)

func FilterByDotNumber(db *gorm.DB, dotNumber string) (out model.SfAccount, err error) {
	//out, err = Use(a.GormDb).SfAccount.FilterByDotNumber(dotNumber)
	account := query.Use(db).SfAccount
	out, err = account.FilterWithColumn(account.DOTNumberC.ColumnName().String(), dotNumber)
	return
}
