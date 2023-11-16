package sql

import (
	"capstonea03/be/src/libs/validator"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ModelI interface {
	TableName() string
}

type Service[T ModelI] struct {
	db *DB
}

func NewService[T ModelI](db *DB) *Service[T] {
	if db == nil {
		logger.Panic("service.DB must exist")
	}

	model := new(T)
	if err := db.AutoMigrate(model); err != nil {
		logger.Panic(err)
	}

	return &Service[T]{db: db}
}

func Transaction(db *DB, txs ...func(tx *DB) *DB) error {
	if err := db.Transaction(func(tx *DB) error {
		for i := range txs {
			if err := txs[i](tx).Error; err != nil {
				logger.Error(err)
				return err
			}
		}
		return nil
	}); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (s *Service[T]) DB() *DB {
	return s.db
}

func (s *Service[T]) Exists(existsOptions *ExistsOptions) (*bool, error) {
	docStruct := new(T)

	existsQuery := s.db.Model(docStruct)

	if existsOptions.Where != nil {
		for _, where := range *existsOptions.Where {
			existsQuery = existsQuery.Where(where.Query, where.Args...)
		}
	}
	if existsOptions.IsUnscoped {
		existsQuery = existsQuery.Unscoped()
	}

	exists := new(bool)
	if err := existsQuery.Select("count(*) > 0").Find(exists).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return exists, nil
}

func (s *Service[T]) Count(countOptions *CountOptions) (*int64, error) {
	docStruct := new(T)

	countQuery := s.db.Model(docStruct)

	if countOptions.Where != nil {
		for _, where := range *countOptions.Where {
			countQuery = countQuery.Where(where.Query, where.Args...)
		}
	}
	if countOptions.IsUnscoped {
		countQuery = countQuery.Unscoped()
	}

	count := new(int64)
	if err := countQuery.Count(count).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return count, nil
}

func (s *Service[T]) FindOne(findOptions *FindOneOptions) (*T, error) {
	docStruct := new(T)

	selectQuery := s.db.Model(docStruct)

	if findOptions.IncludeTables != nil {
		for _, table := range *findOptions.IncludeTables {
			selectQuery = selectQuery.Preload(table.Query, table.Args...)
		}
	}
	if findOptions.Where != nil {
		for _, where := range *findOptions.Where {
			selectQuery = selectQuery.Where(where.Query, where.Args...)
		}
	}
	if findOptions.Order != nil {
		for _, order := range *findOptions.Order {
			selectQuery = selectQuery.Order(order)
		}
	}
	if findOptions.IsUnscoped {
		selectQuery = selectQuery.Unscoped()
	}

	if err := selectQuery.Take(docStruct).Error; err != nil {
		if !IsErrRecordNotFound(err) {
			logger.Error(err)
		}
		return nil, err
	}

	return docStruct, nil
}

func (s *Service[T]) FindAll(findOptions *FindAllOptions) (*[]*T, *Pagination, error) {
	docStruct := &[]*T{}

	selectQuery := s.db.Model(docStruct)

	if findOptions != nil {
		if findOptions.IncludeTables != nil {
			for _, table := range *findOptions.IncludeTables {
				selectQuery = selectQuery.Preload(table.Query, table.Args...)
			}
		}
		if findOptions.Distinct != nil {
			distincts := make([]interface{}, 0, len(*findOptions.Distinct))
			for _, where := range *findOptions.Distinct {
				distincts = append(distincts, where)
			}
			selectQuery.Distinct(distincts...)
		}
		if findOptions.Where != nil {
			for _, where := range *findOptions.Where {
				if where.IncludeInCount {
					selectQuery = selectQuery.Where(where.Where.Query, where.Where.Args...)
				}
			}
		}
		if findOptions.IsUnscoped {
			selectQuery = selectQuery.Unscoped()
		}
	}

	total := new(int64)
	if err := selectQuery.Count(total).Error; err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	if findOptions != nil {
		if findOptions.Where != nil {
			for _, where := range *findOptions.Where {
				if !where.IncludeInCount {
					selectQuery = selectQuery.Where(where.Where.Query, where.Where.Args...)
				}
			}
		}
		if findOptions.Order != nil {
			for _, order := range *findOptions.Order {
				selectQuery = selectQuery.Order(order)
			}
		}
		if findOptions.Offset != nil && *findOptions.Offset > 0 {
			selectQuery = selectQuery.Offset(*findOptions.Offset)
		}
	}
	if findOptions != nil && findOptions.Limit != nil {
		if *findOptions.Limit > FindAllMaximumLimit {
			*findOptions.Limit = FindAllMaximumLimit
		} else if *findOptions.Limit == 0 {
			*findOptions.Limit = FindAllDefaultLimit
		}
	} else {
		*findOptions.Limit = FindAllDefaultLimit
	}
	if *findOptions.Limit != -1 {
		selectQuery = selectQuery.Limit(*findOptions.Limit)
	}

	if err := selectQuery.Find(docStruct).Error; err != nil {
		logger.Error(err)
		return nil, nil, err
	}

	return docStruct, &Pagination{
		Limit: *findOptions.Limit,
		Count: len(*docStruct),
		Total: int(*total),
	}, nil
}

func (s *Service[T]) Create(data *T, createOptions ...*CreateOptions) (*T, error) {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return nil, err
	}

	if err := s.CreateTx(s.db, data, createOptions...).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) BulkCreate(data *[]*T, createOptions ...*CreateOptions) (*[]*T, error) {
	for _, doc := range *data {
		if err := validator.Struct(doc); err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	if err := s.BulkCreateTx(s.db, data, createOptions...).Error; err != nil {
		logger.Error(err)
		return nil, err
	}

	return data, nil
}

func (s *Service[T]) Update(data *T, updateOptions ...*UpdateOptions) (*T, error) {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return nil, err
	}

	tx := s.UpdateTx(s.db, data, updateOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, nil
}

func (s *Service[T]) BulkUpdate(data *[]*T, updateOptions ...*UpdateOptions) (*[]*T, error) {
	for _, doc := range *data {
		if err := validator.Struct(doc); err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	tx := s.BulkUpdateTx(s.db, data, updateOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return nil, err
	}
	if tx.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return data, nil
}

func (s *Service[T]) Replace(data *T, replaceOptions ...*ReplaceOptions) error {
	if err := validator.Struct(data); err != nil {
		logger.Error(err)
		return err
	}

	tx := s.ReplaceTx(s.db, data, replaceOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *Service[T]) Destroy(destroyOptions ...*DestroyOptions) error {
	tx := s.DestroyTx(s.db, destroyOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *Service[T]) BulkDestroy(data *[]*T, destroyOptions ...*DestroyOptions) error {
	tx := s.BulkDestroyTx(s.db, data, destroyOptions...)
	if err := tx.Error; err != nil {
		logger.Error(err)
		return err
	}
	if tx.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

func (s *Service[T]) CreateTx(tx *DB, data *T, createOptions ...*CreateOptions) *DB {
	docStruct := new(T)

	insertQuery := tx.Model(docStruct)

	if len(createOptions) > 0 {
		if createOptions[0].IsUpsert {
			insertQuery = insertQuery.Clauses(clause.OnConflict{UpdateAll: true})
		}
	}

	return insertQuery.Create(data)
}

func (s *Service[T]) BulkCreateTx(tx *DB, data *[]*T, createOptions ...*CreateOptions) *DB {
	docStruct := new(T)

	insertQuery := tx.Model(docStruct)

	if len(createOptions) > 0 {
		if createOptions[0].IsUpsert {
			insertQuery = insertQuery.Clauses(clause.OnConflict{UpdateAll: true})
		}
	}

	return insertQuery.Create(data)
}

func (s *Service[T]) UpdateTx(tx *DB, data *T, updateOptions ...*UpdateOptions) *DB {
	docStruct := new(T)

	updateQuery := tx.Model(docStruct)

	if len(updateOptions) > 0 && updateOptions[0].Where != nil {
		for _, where := range *updateOptions[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
		if updateOptions[0].IsUnscoped {
			updateQuery = updateQuery.Unscoped()
		}
	}

	return updateQuery.Updates(data)
}

func (s *Service[T]) BulkUpdateTx(tx *DB, data *[]*T, updateOptions ...*UpdateOptions) *DB {
	docStruct := new(T)

	updateQuery := tx.Model(docStruct)

	if len(updateOptions) > 0 && updateOptions[0].Where != nil {
		for _, where := range *updateOptions[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
		if updateOptions[0].IsUnscoped {
			updateQuery = updateQuery.Unscoped()
		}
	}

	return updateQuery.Updates(data)
}

func (s *Service[T]) ReplaceTx(tx *DB, data *T, replaceOptions ...*ReplaceOptions) *DB {
	docStruct := new(T)

	updateQuery := tx.Model(docStruct)

	if len(replaceOptions) > 0 && replaceOptions[0].Where != nil {
		for _, where := range *replaceOptions[0].Where {
			updateQuery = updateQuery.Where(where.Query, where.Args...)
		}
		if replaceOptions[0].IsUnscoped {
			updateQuery = updateQuery.Unscoped()
		}
	}

	return updateQuery.Updates(data)
}

func (s *Service[T]) DestroyTx(tx *DB, destroyOptions ...*DestroyOptions) *DB {
	docStruct := new(T)

	deleteQuery := tx.Model(docStruct)

	if len(destroyOptions) > 0 && destroyOptions[0].Where != nil {
		for _, where := range *destroyOptions[0].Where {
			deleteQuery = deleteQuery.Where(where.Query, where.Args...)
		}
		if destroyOptions[0].IsUnscoped {
			deleteQuery = deleteQuery.Unscoped()
		}
	}

	return deleteQuery.Delete(docStruct)
}

func (s *Service[T]) BulkDestroyTx(tx *DB, data *[]*T, destroyOptions ...*DestroyOptions) *DB {
	deleteQuery := tx

	if len(destroyOptions) > 0 && destroyOptions[0].Where != nil {
		for _, where := range *destroyOptions[0].Where {
			deleteQuery = deleteQuery.Where(where.Query, where.Args...)
		}
		if destroyOptions[0].IsUnscoped {
			deleteQuery = deleteQuery.Unscoped()
		}
	}

	return deleteQuery.Delete(data)
}
