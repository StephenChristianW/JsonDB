package configFileIO

import (
	"sync"
)

const (
	collectionSettingsError = "JsonDB/fileIO/configFileIO/collectionSettings.go"
	collectionIOError       = "JsonDB/fileIO/configFileIO/collectionIO.go"
	dbIOError               = "JsonDB/fileIO/configFileIO/dbIO.go"
)

var configMu sync.RWMutex // 配置文件读写锁，保证并发安全

// ==================== 配置结构体 ====================

// configuration 根结构体，保存所有数据库配置
type configuration struct {
	Databases map[string]dbConfig `json:"databases"`
}

// dbConfig 数据库配置，包含数据库名、创建时间、更新时间以及集合信息
type dbConfig struct {
	DBName      string                      `json:"db_name"`
	CreateAt    string                      `json:"create_at"`
	UpdateAt    string                      `json:"update_at"`
	Collections map[string]collectionConfig `json:"collections"`
}

// collectionConfig 集合配置，包含创建/更新时间、自定义设置、文档数量
type collectionConfig struct {
	CreateAt  string             `json:"create_at"`
	UpdateAt  string             `json:"update_at"`
	Settings  collectionSettings `json:"settings"`
	DocsCount int                `json:"fields_count"`
}

// collectionSettings 集合的自定义约束，包括唯一字段和索引
type collectionSettings struct {
	UniqueField map[string]struct{} `json:"unique_field"`
	Index       map[string]struct{} `json:"index"`
}
