package JsonDB

import (
	"JsonDB/services"
	"encoding/json"
	"errors"
)

// ---------------- 高层服务 ----------------

type DBManager struct {
	Ctx *services.DBContext
}

// NewDBManager 创建新实例
func NewDBManager(dbName, collectionName string) *DBManager {
	ctx := &services.DBContext{
		CurrentDB:         dbName,
		CurrentCollection: collectionName,
	}
	return &DBManager{Ctx: ctx}
}

// ---------------- Doc操作封装 ----------------

func (m *DBManager) Find(filter map[string]interface{}, opts *services.FindOptions) (services.DocumentList, error) {
	return m.Ctx.Find(filter, opts)
}

func (m *DBManager) FindOne(filter map[string]interface{}) (services.Document, error) {
	return m.Ctx.FindOne(filter)
}

func (m *DBManager) Insert(doc services.Document) (services.Document, error) {
	return m.Ctx.InsertOne(doc)
}

func (m *DBManager) InsertMany(docs []services.Document) ([]services.Document, error) {
	return m.Ctx.InsertMany(docs)
}

func (m *DBManager) Update(filter map[string]interface{}, update services.Document) (services.Document, error) {
	return m.Ctx.UpdateOne(filter, update)
}

func (m *DBManager) UpdateMany(filter map[string]interface{}, update services.Document) ([]services.Document, error) {
	return m.Ctx.UpdateMany(filter, update)
}

func (m *DBManager) Delete(filter map[string]interface{}) (int, error) {
	return m.Ctx.Delete(filter)
}

// ---------------- Collection操作封装 ----------------

func (m *DBManager) SwitchCollection(name string) {
	m.Ctx.CurrentCollection = name
}

func (m *DBManager) CreateCollection(name string) error {
	return m.Ctx.CollectionCreate(name)
}

func (m *DBManager) DeleteCollection(name string) error {
	return m.Ctx.CollectionDelete(name)
}

func (m *DBManager) ListCollections() ([]string, error) {
	return m.Ctx.CollectionList(m.Ctx.CurrentDB)
}

func (m *DBManager) RenameCollection(oldName, newName string) error {
	return m.Ctx.CollectionRename(oldName, newName)
}

// ---------------- Database操作封装 ----------------

func (m *DBManager) SwitchDB(name string) {
	m.Ctx.CurrentDB = name
}

func (m *DBManager) CreateDB(name string) error {
	return m.Ctx.DBCreate(name)
}

func (m *DBManager) DeleteDB(name string) error {
	return m.Ctx.DBDelete(name)
}

func (m *DBManager) RenameDB(oldName, newName string) error {
	return m.Ctx.DBRename(oldName, newName)
}

func (m *DBManager) ListDBs() ([]string, error) {
	return m.Ctx.DBList()
}

// ---------------- Field操作封装 ----------------

func (m *DBManager) SetUniqueField(field string) error {
	return m.Ctx.SetUniqueField(m.Ctx.CurrentCollection, field)
}

func (m *DBManager) UnSetUniqueField(field string) error {
	return m.Ctx.UnSetUniqueField(m.Ctx.CurrentCollection, field)
}

func (m *DBManager) SetUniqueFields(fields []string) error {
	return m.Ctx.SetUniqueFields(m.Ctx.CurrentCollection, fields)
}

func (m *DBManager) UnSetUniqueFields(fields []string) error {
	return m.Ctx.UnSetUniqueFields(m.Ctx.CurrentCollection, fields)
}

func (m *DBManager) CreateIndex(field string) error {
	return m.Ctx.CreateIndex(m.Ctx.CurrentCollection, field)
}

func (m *DBManager) DropIndex(field string) error {
	return m.Ctx.DropIndex(m.Ctx.CurrentCollection, field)
}

func (m *DBManager) CreateIndexes(fields []string) error {
	return m.Ctx.CreateIndexes(m.Ctx.CurrentCollection, fields)
}

func (m *DBManager) DropIndexes(fields []string) error {
	return m.Ctx.DropIndexes(m.Ctx.CurrentCollection, fields)
}

// ParseJSON 将字符串解析为 map[string]interface{}
func ParseJSON(input string) (map[string]interface{}, error) {
	if input == "" {
		return nil, errors.New("输入内容为空，无法解析")
	}

	var result map[string]interface{}
	err := json.Unmarshal([]byte(input), &result)
	if err != nil {
		return nil, errors.New("JSON 解析失败: " + err.Error())
	}

	return result, nil
}
