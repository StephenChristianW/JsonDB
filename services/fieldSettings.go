package services

import (
	"github.com/StephenChristianW/JsonDB/fileIO"
	Config "github.com/StephenChristianW/JsonDB/fileIO/configFileIO"
)

// FieldService 字段约束 & 索引服务接口
// 用于对集合的字段设置唯一约束或普通索引
type FieldService interface {
	// ==================== 唯一字段 uniqueField ====================

	SetUniqueField(collectionName string, field string) error
	UnSetUniqueField(collectionName string, field string) error
	SetUniqueFields(collectionName string, fields []string) error
	UnSetUniqueFields(collectionName string, fields []string) error
	// ==================== 普通索引 index ====================

	CreateIndex(collectionName string, index string) error
	DropIndex(collectionName string, index string) error
	CreateIndexes(collectionName string, indexes []string) error
	DropIndexes(collectionName string, indexes []string) error
}

const fieldSettingsPath = "JsonDB/services/fieldSettings.go"

// writeSettingsError 错误统一处理函数
// - funcName: 出错的函数名
// - err: 实际捕获的错误
// - msg: 可选的补充说明
// 如果 err 不为 nil，则写入错误日志并返回原始错误
func writeSettingsError(funcName string, err error, msg string) error {
	if err == nil {
		return nil
	}
	var errInfo string
	if msg != "" {
		errInfo = err.Error() + " | msg: " + msg
	} else {
		errInfo = err.Error()
	}
	fileIO.WriteErrorInfo(errInfo, funcName, fieldSettingsPath)
	return err
}

// ==================== 唯一字段 uniqueField ====================

// SetUniqueField 为集合设置 单个唯一字段索引
// - collectionName: 集合名
// - field: 需要设置唯一约束的字段
func (db *DBContext) SetUniqueField(collectionName string, field string) error {
	err := Config.SetUniqueField(db.CurrentDB, collectionName, field)
	return writeSettingsError("SetUniqueField", err, "")
}

// UnSetUniqueField 取消集合的 单个唯一字段索引
// - collectionName: 集合名
// - field: 需要取消唯一约束的字段
func (db *DBContext) UnSetUniqueField(collectionName string, field string) error {
	err := Config.UnSetUniqueField(db.CurrentDB, collectionName, field)
	return writeSettingsError("UnSetUniqueField", err, "")
}

// SetUniqueFields 为集合设置 多个唯一字段索引
// - collectionName: 集合名
// - fields: 需要设置唯一约束的字段列表
func (db *DBContext) SetUniqueFields(collectionName string, fields []string) error {
	err := Config.SetUniqueFields(db.CurrentDB, collectionName, fields)
	return writeSettingsError("SetUniqueFields", err, "")
}

// UnSetUniqueFields 取消集合的 多个唯一字段索引
// - collectionName: 集合名
// - fields: 需要取消唯一约束的字段列表
func (db *DBContext) UnSetUniqueFields(collectionName string, fields []string) error {
	err := Config.UnSetUniqueFields(db.CurrentDB, collectionName, fields)
	return writeSettingsError("UnSetUniqueFields", err, "")
}

// ==================== 普通索引 index ====================

// CreateIndex 为集合创建 单个普通索引
// - collectionName: 集合名
// - index: 索引字段名
func (db *DBContext) CreateIndex(collectionName string, index string) error {
	err := Config.CreateIndex(db.CurrentDB, collectionName, index)
	return writeSettingsError("CreateIndexConfig", err, "")
}

// DropIndex 删除集合的 单个普通索引
// - collectionName: 集合名
// - index: 要删除的索引字段名
func (db *DBContext) DropIndex(collectionName string, index string) error {
	err := Config.DropIndex(db.CurrentDB, collectionName, index)
	return writeSettingsError("DropIndexConfig", err, "")
}

// CreateIndexes 为集合批量创建 普通索引
// - collectionName: 集合名
// - indexes: 需要创建索引的字段列表
func (db *DBContext) CreateIndexes(collectionName string, indexes []string) error {
	err := Config.CreateIndexes(db.CurrentDB, collectionName, indexes)
	return writeSettingsError("CreateIndexes", err, "")
}

// DropIndexes 批量删除集合的 普通索引
// - collectionName: 集合名
// - indexes: 需要删除索引的字段列表
func (db *DBContext) DropIndexes(collectionName string, indexes []string) error {
	err := Config.DropIndexes(db.CurrentDB, collectionName, indexes)
	return writeSettingsError("DropIndexes", err, "")
}
