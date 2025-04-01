package config

import (
	"database/sql"
	"gorm.io/gorm"
	"gorm.io/gorm/migrator"
	"strings"
)

//type ITableInfo interface {
//	GetTableColumns(schemaName string, tableName string) (result []*model.Column, err error)
//	GetTableIndex(schemaName string, tableName string) (indexes []gorm.Index, err error)
//}

func (m *DatabricksMigrator) CurrentSchema(stmt *gorm.Statement, table string) (string, string) {
	if tables := strings.Split(table, `.`); len(tables) == 2 {
		return tables[0], tables[1]
	}
	m.DB = m.DB.Table(table)
	return m.Schema, table
}

// ColumnTypes column types return columnTypes,error
func (m *DatabricksMigrator) ColumnTypes(value interface{}) ([]gorm.ColumnType, error) {
	columnTypes := make([]gorm.ColumnType, 0)
	err := m.Migrator.RunWithValue(value, func(stmt *gorm.Statement) error {
		var (
			currentDatabase, table = m.CurrentSchema(stmt, stmt.Table)
			columnTypeSQL          = "SELECT column_name, column_default, is_nullable = 'YES', data_type, full_data_type, character_maximum_length, numeric_precision, numeric_scale "
			//"SELECT column_name, column_default, is_nullable = 'YES', data_type, character_maximum_length, column_type, column_key, extra, column_comment, numeric_precision, numeric_scale "
			rows, err = m.GormDB.Session(&gorm.Session{}).Table(table).Limit(1).Rows()
		)

		if err != nil {
			return err
		}

		rawColumnTypes, err := rows.ColumnTypes()

		if err != nil {
			return err
		}

		if err := rows.Close(); err != nil {
			return err
		}

		columnTypeSQL += "FROM information_schema.columns WHERE table_schema = ? AND table_name = ? ORDER BY ORDINAL_POSITION"

		columns, rowErr := m.GormDB.Table(table).Raw(columnTypeSQL, currentDatabase, table).Rows()
		if rowErr != nil {
			return rowErr
		}

		defer columns.Close()

		for columns.Next() {
			var (
				column            migrator.ColumnType
				datetimePrecision sql.NullInt64
				extraValue        sql.NullString
				columnKey         sql.NullString
				values            = []interface{}{
					&column.NameValue, &column.DefaultValueValue, &column.NullableValue, &column.DataTypeValue, &column.ColumnTypeValue, &column.LengthValue, &column.DecimalSizeValue, &column.ScaleValue,
					// &column.NameValue, &column.DefaultValueValue, &column.NullableValue, &column.DataTypeValue, &column.LengthValue, &column.ColumnTypeValue, &columnKey, &extraValue, &column.CommentValue, &column.DecimalSizeValue, &column.ScaleValue,
				}
			)

			if scanErr := columns.Scan(values...); scanErr != nil {
				return scanErr
			}

			column.PrimaryKeyValue = sql.NullBool{Bool: false, Valid: true}
			column.UniqueValue = sql.NullBool{Bool: false, Valid: true}
			switch columnKey.String {
			case "PRI":
				column.PrimaryKeyValue = sql.NullBool{Bool: true, Valid: true}
			case "UNI":
				column.UniqueValue = sql.NullBool{Bool: true, Valid: true}
			}

			if strings.Contains(extraValue.String, "auto_increment") {
				column.AutoIncrementValue = sql.NullBool{Bool: true, Valid: true}
			}

			// only trim paired single-quotes
			s := column.DefaultValueValue.String
			for (len(s) >= 3 && s[0] == '\'' && s[len(s)-1] == '\'' && s[len(s)-2] != '\\') ||
				(len(s) == 2 && s == "''") {
				s = s[1 : len(s)-1]
			}
			column.DefaultValueValue.String = s

			if datetimePrecision.Valid {
				column.DecimalSizeValue = datetimePrecision
			}

			for _, c := range rawColumnTypes {
				if c.Name() == column.NameValue.String {
					column.SQLColumnType = c
					break
				}
			}

			columnTypes = append(columnTypes, column)
		}

		return nil
	})

	return columnTypes, err
}
