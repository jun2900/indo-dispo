// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dal

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"github.com/jun2900/indo-dispo/model"
)

func newSupplier(db *gorm.DB, opts ...gen.DOOption) supplier {
	_supplier := supplier{}

	_supplier.supplierDo.UseDB(db, opts...)
	_supplier.supplierDo.UseModel(&model.Supplier{})

	tableName := _supplier.supplierDo.TableName()
	_supplier.ALL = field.NewAsterisk(tableName)
	_supplier.SupplierID = field.NewInt32(tableName, "supplier_id")
	_supplier.SupplierName = field.NewString(tableName, "supplier_name")
	_supplier.SupplierEmail = field.NewString(tableName, "supplier_email")
	_supplier.SupplierTelephone = field.NewString(tableName, "supplier_telephone")
	_supplier.SupplierWeb = field.NewString(tableName, "supplier_web")
	_supplier.SupplierNpwp = field.NewString(tableName, "supplier_npwp")
	_supplier.SupplierAddress = field.NewString(tableName, "supplier_address")
	_supplier.SupplierType = field.NewString(tableName, "supplier_type")
	_supplier.SupplierWhatsapp = field.NewString(tableName, "supplier_whatsapp")
	_supplier.SupplierDescription = field.NewString(tableName, "supplier_description")
	_supplier.SupplierCity = field.NewString(tableName, "supplier_city")
	_supplier.SupplierState = field.NewString(tableName, "supplier_state")
	_supplier.SupplierZipCode = field.NewString(tableName, "supplier_zip_code")

	_supplier.fillFieldMap()

	return _supplier
}

type supplier struct {
	supplierDo

	ALL                 field.Asterisk
	SupplierID          field.Int32
	SupplierName        field.String
	SupplierEmail       field.String
	SupplierTelephone   field.String
	SupplierWeb         field.String
	SupplierNpwp        field.String
	SupplierAddress     field.String
	SupplierType        field.String
	SupplierWhatsapp    field.String
	SupplierDescription field.String
	SupplierCity        field.String
	SupplierState       field.String
	SupplierZipCode     field.String

	fieldMap map[string]field.Expr
}

func (s supplier) Table(newTableName string) *supplier {
	s.supplierDo.UseTable(newTableName)
	return s.updateTableName(newTableName)
}

func (s supplier) As(alias string) *supplier {
	s.supplierDo.DO = *(s.supplierDo.As(alias).(*gen.DO))
	return s.updateTableName(alias)
}

func (s *supplier) updateTableName(table string) *supplier {
	s.ALL = field.NewAsterisk(table)
	s.SupplierID = field.NewInt32(table, "supplier_id")
	s.SupplierName = field.NewString(table, "supplier_name")
	s.SupplierEmail = field.NewString(table, "supplier_email")
	s.SupplierTelephone = field.NewString(table, "supplier_telephone")
	s.SupplierWeb = field.NewString(table, "supplier_web")
	s.SupplierNpwp = field.NewString(table, "supplier_npwp")
	s.SupplierAddress = field.NewString(table, "supplier_address")
	s.SupplierType = field.NewString(table, "supplier_type")
	s.SupplierWhatsapp = field.NewString(table, "supplier_whatsapp")
	s.SupplierDescription = field.NewString(table, "supplier_description")
	s.SupplierCity = field.NewString(table, "supplier_city")
	s.SupplierState = field.NewString(table, "supplier_state")
	s.SupplierZipCode = field.NewString(table, "supplier_zip_code")

	s.fillFieldMap()

	return s
}

func (s *supplier) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := s.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (s *supplier) fillFieldMap() {
	s.fieldMap = make(map[string]field.Expr, 13)
	s.fieldMap["supplier_id"] = s.SupplierID
	s.fieldMap["supplier_name"] = s.SupplierName
	s.fieldMap["supplier_email"] = s.SupplierEmail
	s.fieldMap["supplier_telephone"] = s.SupplierTelephone
	s.fieldMap["supplier_web"] = s.SupplierWeb
	s.fieldMap["supplier_npwp"] = s.SupplierNpwp
	s.fieldMap["supplier_address"] = s.SupplierAddress
	s.fieldMap["supplier_type"] = s.SupplierType
	s.fieldMap["supplier_whatsapp"] = s.SupplierWhatsapp
	s.fieldMap["supplier_description"] = s.SupplierDescription
	s.fieldMap["supplier_city"] = s.SupplierCity
	s.fieldMap["supplier_state"] = s.SupplierState
	s.fieldMap["supplier_zip_code"] = s.SupplierZipCode
}

func (s supplier) clone(db *gorm.DB) supplier {
	s.supplierDo.ReplaceConnPool(db.Statement.ConnPool)
	return s
}

func (s supplier) replaceDB(db *gorm.DB) supplier {
	s.supplierDo.ReplaceDB(db)
	return s
}

type supplierDo struct{ gen.DO }

func (s supplierDo) Debug() *supplierDo {
	return s.withDO(s.DO.Debug())
}

func (s supplierDo) WithContext(ctx context.Context) *supplierDo {
	return s.withDO(s.DO.WithContext(ctx))
}

func (s supplierDo) ReadDB() *supplierDo {
	return s.Clauses(dbresolver.Read)
}

func (s supplierDo) WriteDB() *supplierDo {
	return s.Clauses(dbresolver.Write)
}

func (s supplierDo) Session(config *gorm.Session) *supplierDo {
	return s.withDO(s.DO.Session(config))
}

func (s supplierDo) Clauses(conds ...clause.Expression) *supplierDo {
	return s.withDO(s.DO.Clauses(conds...))
}

func (s supplierDo) Returning(value interface{}, columns ...string) *supplierDo {
	return s.withDO(s.DO.Returning(value, columns...))
}

func (s supplierDo) Not(conds ...gen.Condition) *supplierDo {
	return s.withDO(s.DO.Not(conds...))
}

func (s supplierDo) Or(conds ...gen.Condition) *supplierDo {
	return s.withDO(s.DO.Or(conds...))
}

func (s supplierDo) Select(conds ...field.Expr) *supplierDo {
	return s.withDO(s.DO.Select(conds...))
}

func (s supplierDo) Where(conds ...gen.Condition) *supplierDo {
	return s.withDO(s.DO.Where(conds...))
}

func (s supplierDo) Order(conds ...field.Expr) *supplierDo {
	return s.withDO(s.DO.Order(conds...))
}

func (s supplierDo) Distinct(cols ...field.Expr) *supplierDo {
	return s.withDO(s.DO.Distinct(cols...))
}

func (s supplierDo) Omit(cols ...field.Expr) *supplierDo {
	return s.withDO(s.DO.Omit(cols...))
}

func (s supplierDo) Join(table schema.Tabler, on ...field.Expr) *supplierDo {
	return s.withDO(s.DO.Join(table, on...))
}

func (s supplierDo) LeftJoin(table schema.Tabler, on ...field.Expr) *supplierDo {
	return s.withDO(s.DO.LeftJoin(table, on...))
}

func (s supplierDo) RightJoin(table schema.Tabler, on ...field.Expr) *supplierDo {
	return s.withDO(s.DO.RightJoin(table, on...))
}

func (s supplierDo) Group(cols ...field.Expr) *supplierDo {
	return s.withDO(s.DO.Group(cols...))
}

func (s supplierDo) Having(conds ...gen.Condition) *supplierDo {
	return s.withDO(s.DO.Having(conds...))
}

func (s supplierDo) Limit(limit int) *supplierDo {
	return s.withDO(s.DO.Limit(limit))
}

func (s supplierDo) Offset(offset int) *supplierDo {
	return s.withDO(s.DO.Offset(offset))
}

func (s supplierDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *supplierDo {
	return s.withDO(s.DO.Scopes(funcs...))
}

func (s supplierDo) Unscoped() *supplierDo {
	return s.withDO(s.DO.Unscoped())
}

func (s supplierDo) Create(values ...*model.Supplier) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Create(values)
}

func (s supplierDo) CreateInBatches(values []*model.Supplier, batchSize int) error {
	return s.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (s supplierDo) Save(values ...*model.Supplier) error {
	if len(values) == 0 {
		return nil
	}
	return s.DO.Save(values)
}

func (s supplierDo) First() (*model.Supplier, error) {
	if result, err := s.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.Supplier), nil
	}
}

func (s supplierDo) Take() (*model.Supplier, error) {
	if result, err := s.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.Supplier), nil
	}
}

func (s supplierDo) Last() (*model.Supplier, error) {
	if result, err := s.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.Supplier), nil
	}
}

func (s supplierDo) Find() ([]*model.Supplier, error) {
	result, err := s.DO.Find()
	return result.([]*model.Supplier), err
}

func (s supplierDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.Supplier, err error) {
	buf := make([]*model.Supplier, 0, batchSize)
	err = s.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (s supplierDo) FindInBatches(result *[]*model.Supplier, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return s.DO.FindInBatches(result, batchSize, fc)
}

func (s supplierDo) Attrs(attrs ...field.AssignExpr) *supplierDo {
	return s.withDO(s.DO.Attrs(attrs...))
}

func (s supplierDo) Assign(attrs ...field.AssignExpr) *supplierDo {
	return s.withDO(s.DO.Assign(attrs...))
}

func (s supplierDo) Joins(fields ...field.RelationField) *supplierDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Joins(_f))
	}
	return &s
}

func (s supplierDo) Preload(fields ...field.RelationField) *supplierDo {
	for _, _f := range fields {
		s = *s.withDO(s.DO.Preload(_f))
	}
	return &s
}

func (s supplierDo) FirstOrInit() (*model.Supplier, error) {
	if result, err := s.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.Supplier), nil
	}
}

func (s supplierDo) FirstOrCreate() (*model.Supplier, error) {
	if result, err := s.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.Supplier), nil
	}
}

func (s supplierDo) FindByPage(offset int, limit int) (result []*model.Supplier, count int64, err error) {
	result, err = s.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = s.Offset(-1).Limit(-1).Count()
	return
}

func (s supplierDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = s.Count()
	if err != nil {
		return
	}

	err = s.Offset(offset).Limit(limit).Scan(result)
	return
}

func (s supplierDo) Scan(result interface{}) (err error) {
	return s.DO.Scan(result)
}

func (s supplierDo) Delete(models ...*model.Supplier) (result gen.ResultInfo, err error) {
	return s.DO.Delete(models)
}

func (s *supplierDo) withDO(do gen.Dao) *supplierDo {
	s.DO = *do.(*gen.DO)
	return s
}