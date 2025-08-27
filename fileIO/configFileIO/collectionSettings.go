package configFileIO

// ==================== collection 设置对外接口 ====================

// SetUniqueField 设置 UniqueField
func SetUniqueField(dbName, collectionName, uniqueField string) error {
	return updateFieldMap(dbName, collectionName, []string{uniqueField}, "UniqueField", true)
}

// UnSetUniqueField 取消 UniqueField
func UnSetUniqueField(dbName, collectionName, uniqueField string) error {
	return updateFieldMap(dbName, collectionName, []string{uniqueField}, "UniqueField", false)
}

// CreateIndexConfig 设置 Index
func CreateIndexConfig(dbName, collectionName, index string) error {
	return updateFieldMap(dbName, collectionName, []string{index}, "Index", true)
}

// DropIndexConfig 取消 Index
func DropIndexConfig(dbName, collectionName, index string) error {
	return updateFieldMap(dbName, collectionName, []string{index}, "Index", false)
}

// ==================== collection 设置批量操作接口 ====================

// SetUniqueFields 批量设置集合的唯一字段
func SetUniqueFields(dbName, collectionName string, fields []string) error {
	return updateFieldMap(dbName, collectionName, fields, "UniqueField", true)
}

// UnSetUniqueFields 批量取消集合的唯一字段
func UnSetUniqueFields(dbName, collectionName string, fields []string) error {
	return updateFieldMap(dbName, collectionName, fields, "UniqueField", false)
}

// CreateIndexes 批量创建索引
func CreateIndexes(dbName, collectionName string, indexes []string) error {
	return updateFieldMap(dbName, collectionName, indexes, "Index", true)
}

// DropIndexes 批量删除索引
func DropIndexes(dbName, collectionName string, indexes []string) error {
	return updateFieldMap(dbName, collectionName, indexes, "Index", false)
}
