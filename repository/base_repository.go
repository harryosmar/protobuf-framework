package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-sql-driver/mysql"
	appError "github.com/harryosmar/protobuf-go/error"
	"gorm.io/gorm/schema"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type (
	PkType interface {
		~string | ~int64 | ~uint64 | ~int32 | ~uint32 | int16 | ~uint16 | int8 | ~uint8 | ~int | ~uint
	}

	ServiceRepository[T schema.Tabler, P PkType] interface {
		DB(ctx context.Context) *gorm.DB
		Create(ctx context.Context, row *T) (*T, error)
		CreatePerBatch(ctx context.Context, rows []*T) ([]*T, int64, error)
		UpdateColumns(ctx context.Context, row *T, updatedColumns []string) (int64, error)
		Update(ctx context.Context, row *T) (int64, error)
		Upsert(ctx context.Context, row *T, onConflictUpdatedColumns []string) (int64, error)
		UpdateWhere(ctx context.Context, wheres []Where, values map[string]interface{}) (int64, error)
		GetById(ctx context.Context, id P) (*T, error)
		Delete(ctx context.Context, id P) error
		GetFirst(ctx context.Context, wheres []Where) (*T, error)
		GetAll(ctx context.Context, orders []OrderBy, wheres []Where) ([]T, error)
		GetPerPage(ctx context.Context, page int32, pageSize int32, orders []OrderBy, wheres []Where) ([]T, *Paginator, error)
	}

	Paginator struct {
		Page    int32 `json:"page"`
		PerPage int32 `json:"per_page"`
		Total   int64 `json:"total"`
	}

	Where struct {
		Name             string      `json:"name"`
		IsLike           bool        `json:"is_like"`             // use "%keyword%": WHERE name LIKE '%ware%'
		IsFullTextSearch bool        `json:"is_full_text_search"` // use "*keyword*" : WHERE MATCH(name) AGAINST ('*ware*' IN BOOLEAN MODE) : To fully optimize this, create index "FULLTEXT KEY `idx_fulltext_columName` (`columName`)", read also about stopwords https://dev.mysql.com/doc/refman/8.4/en/fulltext-stopwords.html
		Value            interface{} `json:"value"`
	}

	OrderBy struct {
		Field     string `json:"field"`
		Direction string `json:"direction"` // asc, desc
	}
)

func (o OrderBy) String() string {
	if o.Field != "" && (o.Direction == "asc" || o.Direction == "desc") {
		return fmt.Sprintf("%s %s", o.Field, o.Direction)
	}

	return ""
}

func (c *Where) String() string {
	whereSql := fmt.Sprintf("%s = ?", c.Name)
	if c.IsFullTextSearch {
		whereSql = fmt.Sprintf("MATCH(%s) AGAINST (? IN BOOLEAN MODE)", c.Name)
	} else if c.IsLike {
		whereSql = fmt.Sprintf("%s LIKE ?", c.Name)
	}

	return whereSql
}

type BaseGorm[T schema.Tabler, P PkType] struct {
	db *gorm.DB
}

func NewBaseGorm[T schema.Tabler, P PkType](db *gorm.DB) *BaseGorm[T, P] {
	return &BaseGorm[T, P]{db: db}
}

func (o *BaseGorm[T, P]) Detail(ctx context.Context, id P) (*T, error) {
	var (
		db  = o.db.WithContext(ctx)
		row T
		err error
	)

	if err = db.First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &row, nil
}

func (o *BaseGorm[T, P]) GetById(ctx context.Context, id P) (*T, error) {
	var (
		e  T
		db = o.db.WithContext(ctx).Model(e.TableName())
	)

	if err := db.WithContext(ctx).First(&e, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Not found is not an error at repository level
		}
		return nil, err
	}

	return &e, nil
}
func (o *BaseGorm[T, P]) Delete(ctx context.Context, id P) error {
	var (
		e  T
		db = o.db.WithContext(ctx)
	)

	result := db.WithContext(ctx).Delete(&e, id)
	if result.Error != nil {
		return result.Error
	}
	// Return success even if no rows affected - idempotent delete
	return nil
}

func (o *BaseGorm[T, P]) GetFirst(ctx context.Context, wheres []Where) (*T, error) {
	var (
		e   T
		db  = o.db.WithContext(ctx).Model(e.TableName())
		err error
	)

	for _, v := range wheres {
		if v.IsLike {
			v.Value = fmt.Sprintf("%%%s%%", v.Value)
		}
		db.Where(v.String(), v.Value)
	}

	if err = db.First(&e).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &e, nil
}

func (o *BaseGorm[T, P]) GetAll(ctx context.Context, orders []OrderBy, wheres []Where) ([]T, error) {
	var (
		e    T
		db   = o.db.WithContext(ctx).Table(e.TableName())
		rows []T
		err  error
	)

	for _, v := range wheres {
		if v.IsLike {
			v.Value = fmt.Sprintf("%%%s%%", v.Value)
		}
		db.Where(v.String(), v.Value)
	}

	for _, order := range orders {
		orderByStr := order.String()
		if orderByStr != "" {
			db.Order(orderByStr)
		}
	}

	if err = db.Find(&rows).Error; err != nil {
		return rows, err
	}

	return rows, nil
}

func (o *BaseGorm[T, P]) GetPerPage(ctx context.Context, page int32, pageSize int32, orders []OrderBy, wheres []Where) ([]T, *Paginator, error) {
	var (
		e         T
		db        = o.db.WithContext(ctx).Table(e.TableName())
		rows      []T
		count     int64
		err       error
		paginator = &Paginator{
			Page:    page,
			PerPage: pageSize,
			Total:   0,
		}
	)

	for _, v := range wheres {
		if v.IsLike {
			v.Value = fmt.Sprintf("%%%s%%", v.Value)
		}
		db.Where(v.String(), v.Value)
	}

	for _, order := range orders {
		orderByStr := order.String()
		if orderByStr != "" {
			db.Order(orderByStr)
		}
	}

	if err = db.Count(&count).Error; err != nil {
		return rows, nil, err
	}

	paginator.Total = count
	if count == 0 {
		return rows, paginator, nil
	}

	if err = db.Offset(int((page - 1) * pageSize)).Limit(int(pageSize)).Find(&rows).Error; err != nil {
		return rows, paginator, err
	}

	return rows, paginator, nil
}

func (o *BaseGorm[T, P]) Create(ctx context.Context, row *T) (*T, error) {
	var (
		e   T
		db  = o.db.WithContext(ctx).Table(e.TableName())
		err error
	)

	// cannot handle upsert condition, will get err Duplicate entry
	if err = db.Create(row).Error; err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, appError.ErrRecordAlreadyExists
		}
		return nil, err
	}

	return row, nil
}

func (o *BaseGorm[T, P]) DB(ctx context.Context) *gorm.DB {
	return o.db.WithContext(ctx)
}

func (o *BaseGorm[T, P]) CreatePerBatch(ctx context.Context, rows []*T) ([]*T, int64, error) {
	var (
		rowsAffected int64
	)

	if len(rows) == 0 {
		return rows, rowsAffected, nil
	}

	var (
		e   T
		db  = o.db.WithContext(ctx).Table(e.TableName())
		err error
	)

	result := db.Create(rows)
	err = result.Error
	rowsAffected = result.RowsAffected

	return rows, rowsAffected, err
}

func (o *BaseGorm[T, P]) Update(ctx context.Context, row *T) (int64, error) {
	var (
		e  T
		db = o.db.WithContext(ctx).Table(e.TableName())
	)

	result := db.Model(row).Updates(row)
	return result.RowsAffected, result.Error
}

func (o *BaseGorm[T, P]) UpdateColumns(ctx context.Context, row *T, updatedColumns []string) (int64, error) {
	var (
		e   T
		db  = o.db.WithContext(ctx).Table(e.TableName())
		err error
	)

	if len(updatedColumns) > 0 {
		db = db.Select(updatedColumns)
	}

	// Use the model to get the correct table and add WHERE clause for the primary key
	result := db.Model(row).Updates(row)
	err = result.Error

	return result.RowsAffected, err
}

func (o *BaseGorm[T, P]) UpdateWhere(ctx context.Context, wheres []Where, values map[string]interface{}) (int64, error) {
	var (
		e   T
		db  = o.db.WithContext(ctx).Table(e.TableName())
		err error
	)

	// Build where clauses
	for _, v := range wheres {
		if v.IsLike {
			v.Value = fmt.Sprintf("%%%s%%", v.Value)
		}
		db.Where(v.String(), v.Value)
	}

	// Execute update
	result := db.Updates(values)
	err = result.Error

	return result.RowsAffected, err
}

func (o *BaseGorm[T, P]) Upsert(ctx context.Context, row *T, onConflictUpdatedColumns []string) (int64, error) {
	var (
		e  T
		db = o.db.WithContext(ctx).Table(e.TableName())
	)

	result := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{},
		DoUpdates: clause.AssignmentColumns(onConflictUpdatedColumns),
	}).Create(&row)

	return result.RowsAffected, result.Error
}
