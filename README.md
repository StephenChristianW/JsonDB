# JsonDB

JsonDB 是一个基于 JSON 文件的轻量级数据库管理系统，提供简单易用的命令行界面（CLI），支持数据库、集合和文档操作，并提供索引与唯一字段设置功能。适合小型项目、测试和快速原型开发。

---

## 特性

- 基于 JSON 文件存储数据
- 支持多数据库、多集合
- 文档 CRUD（插入、查询、删除）
- 支持唯一字段与索引管理
- 跨平台兼容 CLI（Windows CMD / Linux / macOS）
- 全局命令支持 `/exit`, `/quit`, `/help`

---

## 安装

1. 克隆仓库：

   ```bash
   git clone https://github.com/StephenChristianW/JsonDB.git
   cd jsondb
   ```

2. 构建或运行：

   ```bash
   cd main
   go build -o jsondb main.go
   ./jsondb
   ```


   或直接运行 

   ```bash
   cd main
   go run main.go
   ```
1. GO Get：

   ```bash
   go get github.com/StephenChristianW/JsonDB@629337f

   ```

2. 调用代码
   ```go
   package main
   
   import (
   "fmt"
   JsonDB "github.com/StephenChristianW/JsonDB"
   )
   
   func main() {
   // ----------------- 初始化 -----------------
   ctx, err := JsonDB.NewDBContext()
   if err != nil {
   fmt.Println("初始化 DBContext 失败:", err)
   return
   }
   manager := JsonDB.NewDBManager(ctx)
   fmt.Println("DBManager 初始化成功")
   
       // ----------------- 创建数据库 -----------------
       dbName := "testDB"
       err = manager.CreateDB(dbName)
       if err != nil {
           fmt.Println("创建数据库失败:", err)
       } else {
           fmt.Println("数据库创建成功:", dbName)
       }
   
       // ----------------- 创建集合 -----------------
       collectionName := "users"
       err = manager.CreateCollection(dbName, collectionName)
       if err != nil {
           fmt.Println("创建集合失败:", err)
       } else {
           fmt.Println("集合创建成功:", collectionName)
       }
   
       // ----------------- 插入文档 -----------------
       doc := map[string]interface{}{
           "username": "yuanlao",
           "age":      28,
           "level":    1,
       }
       err = manager.InsertOne(dbName, collectionName, doc)
       if err != nil {
           fmt.Println("插入文档失败:", err)
       } else {
           fmt.Println("文档插入成功:", doc)
       }
   
       // ----------------- 查询文档 -----------------
       filter := map[string]interface{}{"username": "yuanlao"}
       results, err := manager.Find(dbName, collectionName, filter)
       if err != nil {
           fmt.Println("查询文档失败:", err)
       } else {
           fmt.Println("查询结果:")
           for i, r := range results {
               fmt.Printf("文档 %d: %+v\n", i+1, r)
           }
       }
   
       // ----------------- 更新文档 -----------------
       update := map[string]interface{}{"level": 2}
       err = manager.UpdateOne(dbName, collectionName, filter, update)
       if err != nil {
           fmt.Println("更新文档失败:", err)
       } else {
           fmt.Println("文档更新成功")
       }
   
       // ----------------- 查询更新后的文档 -----------------
       results, _ = manager.Find(dbName, collectionName, filter)
       fmt.Println("更新后的查询结果:")
       for i, r := range results {
           fmt.Printf("文档 %d: %+v\n", i+1, r)
       }
   
       // ----------------- 删除文档 -----------------
       err = manager.DeleteOne(dbName, collectionName, filter)
       if err != nil {
           fmt.Println("删除文档失败:", err)
       } else {
           fmt.Println("文档删除成功")
       }
   
       // ----------------- 查询删除后的文档 -----------------
       results, _ = manager.Find(dbName, collectionName, filter)
       fmt.Println("删除后的查询结果:", results)
   }

   ```

## 命令行界面 (CLI)

启动后，将看到主菜单：

```tex
================== JsonDB 菜单 ==================
当前数据库: None
当前集合: None

1. 数据库操作
2. 集合操作
3. 文档操作
0. 退出
请选择:
```

### 全局命令

- `/exit` 或 `/quit`：退出程序（在任意输入处可使用）
- `/help`：显示帮助文档

## 数据库操作

- 列出数据库
- 创建数据库
- 删除数据库
- 切换数据库

### 示例

```tex
请输入数据库名: testdb
✅ 数据库创建成功: testdb

请输入数据库名: testdb
✅ 已切换到数据库: testdb
```

## 集合操作

- 列出集合
- 创建集合
- 删除集合
- 切换集合

### 示例

```tex
请输入集合名: users
✅ 集合创建成功: testdb.users

请输入集合名: users
✅ 已切换到集合: users
```

## 文档操作

- 插入 JSON 文档
- 查询文档
- 删除文档

### 示例

插入文档：

```json
{
  "name": "Alice",
  "age": 30,
  "email": "alice@example.com"
}
```

查询文档：

```json
{
  "age": 30
}
```

删除文档：

```json
{
  "name": "Alice"
}
```

## 索引与唯一字段操作

在当前集合中，你可以设置字段为唯一或创建索引：

```go
manager.SetUniqueField("username")         // 设置单个字段唯一
manager.UnSetUniqueField("username")       // 取消字段唯一
manager.SetUniqueFields([]string{"id", "email"}) // 设置多个字段唯一
manager.CreateIndex("age")                  // 创建单个字段索引
manager.CreateIndexes([]string{"age","score"}) // 创建多个索引
manager.DropIndex("age")                    // 删除单个索引
manager.DropIndexes([]string{"age","score"}) // 删除多个索引
```

- **唯一字段**：保证字段在集合中不重复
- **索引字段**：加快查询速度

------

## 存储结构

```bash
JsonDataBase/
├── 数据库名/
│   ├── 集合1.json
│   ├── 集合2.json
│   └── ...
└── ...

```

每个集合对应一个 JSON 文件存储所有文档，文件中可包含索引和唯一字段信息。

------

## 注意事项

- JSON 文档必须符合标准格式，否则会解析失败
- 全局命令 `/exit`、`/quit` 可随时退出程序
- 全局命令 `/help` 显示此帮助文档
- 支持跨平台命令行兼容，Windows CMD 可直接使用

------

## 贡献

欢迎提交 Issue 和 Pull Request，帮助 JsonDB 更完善。



# 非商业使用许可 / Non-Commercial Use License

版权所有 © 2025 StephenChristianW  
联系方式: yuanlao1016@gmail.com

---

## 许可说明 / License Terms

### 1. 非商业用途 / Non-Commercial Use
- 个人或组织可 **免费** 使用、复制、修改本软件及其文档，仅限 **非商业目的**（例如学习、研究、个人项目）。
- 非商业用途不得产生直接或间接的利润。

### 2. 商业用途 / Commercial Use
- 商业使用本软件（包括但不限于销售、提供付费服务、企业内部盈利性使用）必须 **事先获得版权所有者的书面授权**。
- 商业授权需支付相应的许可费用（可通过上述邮箱联系作者洽谈）。

### 3. 保留版权 / Copyright
- 使用、复制或修改本软件时，必须保留本版权声明及本许可文件。

### 4. 免责声明 / Disclaimer
- 本软件按“原样”提供，不附带任何明示或暗示的保证，包括但不限于适销性、特定用途适用性及非侵权保证。
- 作者不对因使用本软件产生的任何直接或间接损失承担责任，无论合同、侵权或其他法律形式。

### 5. 法律适用 / Governing Law
- 本许可受中华人民共和国法律管辖。
- 任何未经授权的商业使用可能会承担法律责任。

---

## 联系方式 / Contact
如需商业授权或有其他许可相关问题，请通过邮箱联系作者：  
**yuanlao1016@gmail.com**

---

# Non-Commercial Use License

Copyright © 2025 StephenChristianW  
Contact: yuanlao1016@gmail.com

---

## 1. Non-Commercial Use
- Individuals or organizations may use, copy, and modify this software and its documentation **for non-commercial purposes only**, free of charge.
- Non-commercial purposes must not generate any direct or indirect profit.

## 2. Commercial Use
- Commercial use of this software (including but not limited to selling, providing paid services, or internal profit-making use) requires **prior written authorization** from the copyright holder.
- Commercial authorization requires a licensing fee (please contact the author via the above email).

## 3. Copyright
- All copies or substantial portions of this software must retain this copyright notice and this license file.

## 4. Disclaimer
- This software is provided "as is", without any express or implied warranty, including but not limited to warranties of merchantability, fitness for a particular purpose, and non-infringement.
- The author is not liable for any direct or indirect damages arising from the use of this software, under contract, tort, or any other legal theory.

## 5. Governing Law
- This license is governed by the laws of the People's Republic of China.
- Any unauthorized commercial use may result in legal liability.

## Contact
For commercial licensing or any license-related questions, please contact the author via:  
**yuanlao1016@gmail.com**


