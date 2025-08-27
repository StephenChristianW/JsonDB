package main

import (
	"JsonDB"
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
)

var (
	ColorRed    = color.New(color.FgRed)
	ColorGreen  = color.New(color.FgGreen)
	ColorYellow = color.New(color.FgYellow)
	ColorBlue   = color.New(color.FgBlue)
	ColorCyan   = color.New(color.FgCyan)
)

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	} else {
		_, _ = fmt.Print("\033[H\033[2J")
	}
}

func main() {
	manager := JsonDB.NewDBManager("", "")
	reader := bufio.NewReader(os.Stdin)

	for {
		clearScreen()
		printStatus(manager)
		_, _ = ColorCyan.Print("1. 数据库操作\n2. 集合操作\n3. 文档操作\n0. 退出\n请选择: ")
		choice := readChoice(reader)

		switch choice {
		case 0:
			_, _ = ColorYellow.Println("退出 JsonDB")
			return
		case 1:
			dbMenu(manager, reader)
		case 2:
			collectionMenu(manager, reader)
		case 3:
			documentMenu(manager, reader)
		default:
			_, _ = ColorRed.Println("无效选项，请重新选择")
			pause(reader)
		}
	}
}

// -------------------- 状态栏 --------------------
func printStatus(manager *JsonDB.DBManager) {
	_, _ = ColorBlue.Println("================== JsonDB 菜单 ==================")
	_, _ = ColorYellow.Printf("当前数据库: %s\n", manager.Ctx.CurrentDB)
	_, _ = ColorYellow.Printf("当前集合: %s\n\n", manager.Ctx.CurrentCollection)
}

// -------------------- 数据库菜单 --------------------
func dbMenu(manager *JsonDB.DBManager, reader *bufio.Reader) {
	for {
		clearScreen()
		printStatus(manager)
		_, _ = ColorCyan.Println("---- 数据库操作 ----")
		_, _ = ColorCyan.Println("1. 列出数据库\n2. 创建数据库\n3. 删除数据库\n4. 切换数据库\n0. 返回主菜单")
		_, err := ColorCyan.Print("请选择: ")
		if err != nil {
			return
		}
		choice := readChoice(reader)

		switch choice {
		case 0:
			return
		case 1:
			dbs, _ := manager.ListDBs()
			printTable("数据库列表", dbs)
			pause(reader)
		case 2:
			fmt.Print("请输入数据库名: ")
			name := readLine(reader)
			if name == "" {
				_, _ = ColorRed.Println("❌ 数据库名不能为空")
			} else if err := manager.CreateDB(name); err != nil {
				_, _ = ColorRed.Println("错误:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 数据库创建成功:", name)
			}
			pause(reader)
		case 3:
			fmt.Print("请输入数据库名: ")
			name := readLine(reader)
			if name == "" {
				_, _ = ColorRed.Println("❌ 数据库名不能为空")
			} else if err := manager.DeleteDB(name); err != nil {
				_, _ = ColorRed.Println("错误:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 数据库已删除:", name)
			}
			pause(reader)
		case 4:
			fmt.Print("请输入数据库名: ")
			name := readLine(reader)
			if name == "" {
				_, _ = ColorRed.Println("❌ 数据库名不能为空，切换失败")
			} else {
				manager.SwitchDB(name)
				_, _ = ColorGreen.Println("✅ 已切换到数据库:", name)
			}
			pause(reader)
		default:
			_, _ = ColorRed.Println("无效选项，请重新选择")
			pause(reader)
		}
	}
}

// -------------------- 集合菜单 --------------------
func collectionMenu(manager *JsonDB.DBManager, reader *bufio.Reader) {
	for {
		clearScreen()
		printStatus(manager)
		_, _ = ColorCyan.Println("---- 集合操作 ----")
		_, _ = ColorCyan.Println("1. 列出集合\n2. 创建集合\n3. 删除集合\n4. 切换集合\n5. 创建索引\n6. 删除索引\n0. 返回主菜单")
		_, _ = ColorCyan.Print("请选择: ")
		choice := readChoice(reader)

		switch choice {
		case 0:
			return
		case 1:
			cols, _ := manager.ListCollections()
			printTable("集合列表", cols)
			pause(reader)
		case 2:
			fmt.Print("请输入集合名: ")
			name := readLine(reader)
			if name == "" {
				_, _ = ColorRed.Println("❌ 集合名不能为空")
			} else if err := manager.CreateCollection(name); err != nil {
				_, _ = ColorRed.Println("错误:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 集合创建成功:", manager.Ctx.CurrentDB+"."+name)
			}
			pause(reader)
		case 3:
			fmt.Print("请输入集合名: ")
			name := readLine(reader)
			if name == "" {
				_, _ = ColorRed.Println("❌ 集合名不能为空")
			} else if err := manager.DeleteCollection(name); err != nil {
				_, _ = ColorRed.Println("错误:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 集合已删除:", name)
			}
			pause(reader)
		case 4:
			fmt.Print("请输入集合名: ")
			name := readLine(reader)
			if name == "" {
				_, _ = ColorRed.Println("❌ 集合名不能为空，切换失败")
			} else {
				manager.SwitchCollection(name)
				_, _ = ColorGreen.Println("✅ 已切换到集合:", name)
			}
			pause(reader)
		case 5: // 创建索引
			fmt.Print("请输入要创建索引的字段名（多个用逗号分隔）: ")
			input := readLine(reader)
			fields := strings.Split(input, ",")
			for i := range fields {
				fields[i] = strings.TrimSpace(fields[i])
			}
			if err := manager.CreateIndexes(fields); err != nil {
				_, _ = ColorRed.Println("❌ 创建索引失败:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 索引创建成功:", fields)
			}
			pause(reader)
		case 6: // 删除索引
			fmt.Print("请输入要删除索引的字段名（多个用逗号分隔）: ")
			input := readLine(reader)
			fields := strings.Split(input, ",")
			for i := range fields {
				fields[i] = strings.TrimSpace(fields[i])
			}
			if err := manager.DropIndexes(fields); err != nil {
				_, _ = ColorRed.Println("❌ 删除索引失败:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 索引删除成功:", fields)
			}
			pause(reader)
		default:
			_, _ = ColorRed.Println("无效选项，请重新选择")
			pause(reader)
		}
	}
}

// -------------------- 文档菜单 --------------------
func documentMenu(manager *JsonDB.DBManager, reader *bufio.Reader) {
	for {
		clearScreen()
		printStatus(manager)
		_, _ = ColorCyan.Println("---- 文档操作 ----")
		_, _ = ColorCyan.Println("1. 插入文档\n2. 查询文档\n3. 删除文档\n0. 返回主菜单")
		_, _ = ColorCyan.Print("请选择: ")
		choice := readChoice(reader)

		switch choice {
		case 0:
			return
		case 1:
			fmt.Print("请输入文档 (JSON 格式): ")
			docStr := readLine(reader)
			doc, err := JsonDB.ParseJSON(docStr)
			if err != nil {
				_, _ = ColorRed.Println("❌ JSON 解析错误:", err.Error())
				pause(reader)
				continue
			}
			insertedDoc, err := manager.Insert(doc)
			if err != nil {
				_, _ = ColorRed.Println("❌ 插入失败:", err.Error())
			} else {
				_, _ = ColorGreen.Println("✅ 插入成功:")
				jsonBytes, _ := json.MarshalIndent(insertedDoc, "", "  ")
				fmt.Println(string(jsonBytes))
			}
			pause(reader)
		case 2:
			fmt.Print("请输入查询条件 (JSON 格式): ")
			filterStr := readLine(reader)
			filter, err := JsonDB.ParseJSON(filterStr)
			if err != nil {
				_, _ = ColorRed.Println("❌ JSON 解析错误:", err.Error())
				pause(reader)
				continue
			}
			docs, err := manager.Find(filter, nil)
			if err != nil {
				_, _ = ColorRed.Println("❌ 查询失败:", err.Error())
			} else {
				_, _ = ColorBlue.Println("\n查询结果:")
				if len(docs) == 0 {
					fmt.Println("  （空）")
				} else {
					jsonBytes, _ := json.MarshalIndent(docs, "", "  ")
					fmt.Println(string(jsonBytes))
				}
			}
			pause(reader)
		case 3:
			fmt.Print("请输入删除条件 (JSON 格式): ")
			filterStr := readLine(reader)
			filter, err := JsonDB.ParseJSON(filterStr)
			if err != nil {
				_, _ = ColorRed.Println("❌ JSON 解析错误:", err.Error())
				pause(reader)
				continue
			}
			count, err := manager.Delete(filter)
			if err != nil {
				_, _ = ColorRed.Println("❌ 删除失败:", err.Error())
			} else {
				_, _ = ColorGreen.Printf("✅ 已删除 %d 条文档\n", count)
			}
			pause(reader)
		default:
			_, _ = ColorRed.Println("无效选项，请重新选择")
			pause(reader)
		}
	}
}

// -------------------- 工具函数 --------------------
func readChoice(reader *bufio.Reader) int {
	line := readLine(reader)
	val, err := strconv.Atoi(line)
	if err != nil {
		_, _ = ColorRed.Println("无效输入，请输入数字")
		return -1
	}
	return val
}

func readLine(reader *bufio.Reader) string {
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	switch line {
	case `/exit`, `/quit`:
		_, _ = ColorYellow.Println("全局退出 JsonDB")
		os.Exit(0)
	case "/help":
		printHelp()
		fmt.Print("按 Enter 返回...")
		_, _ = reader.ReadString('\n')
		return readLine(reader)
	}
	return line
}

func pause(reader *bufio.Reader) {
	fmt.Print("按 Enter 键继续...")
	_, _ = reader.ReadString('\n')
}

func printTable(title string, items []string) {
	_, _ = ColorBlue.Println("==== " + title + " ====")
	if len(items) == 0 {
		fmt.Println("（空）")
		return
	}

	numWidth := len(fmt.Sprintf("%d", len(items)))
	nameWidth := len("名称")
	for _, item := range items {
		if len(item) > nameWidth {
			nameWidth = len(item)
		}
	}

	fmt.Printf("| %-*s | %-*s |\n", numWidth, "#", nameWidth, "名称")
	fmt.Printf("|%s|%s|\n", strings.Repeat("-", numWidth+2), strings.Repeat("-", nameWidth+2))

	for i, item := range items {
		fmt.Printf("| %*d | %-*s |\n", numWidth, i+1, nameWidth, item)
	}

	fmt.Printf("%s\n", strings.Repeat("-", numWidth+nameWidth+7))
}

func printHelp() {
	_, _ = ColorCyan.Println("===== JsonDB 全局帮助 =====")
	_, _ = ColorGreen.Println("可用命令:")
	fmt.Println("  /exit | /quit      退出程序（任意输入处可使用）")
	fmt.Println("  /help                显示帮助文档")
	fmt.Println()
	fmt.Println("菜单操作说明:")
	fmt.Println("  1. 数据库操作: 列出/创建/删除/切换数据库")
	fmt.Println("  2. 集合操作: 列出/创建/删除/切换集合 + 索引操作")
	fmt.Println("  3. 文档操作: 插入/查询/删除文档")
	fmt.Println()
	fmt.Println("提示: 在文档和集合菜单中，也可以直接输入 \\exit 或 \\quit 退出程序")
	_, _ = ColorCyan.Println("=============================")
}
