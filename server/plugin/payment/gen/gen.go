// go run gen.go  -h 47.110.88.165:2624 -u root -p J4l71RPO14 -d cabinets -t sc_device_domain,sc_tenant_device,sc_tenant_device_attr,sc_tenant_device_config,sc_tenant_device_group,sc_tenant_device_info,sc_tenant_device_task,sc_tenant_sub_device
//
//nolint:all
package main

// 日志[链路追踪],  日志报警,  异常
// 框架基础包
//    时间转换的
//    json的可以直接用原生的或者 github.com/tidwall/gjson
import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var (
	// tableName
	host      string
	user      string
	pwd       string
	work      string
	database  string
	tableName string
)

// application
// company
// notify_record
// payment_account
// payment_method
// payment_order
// payment_order_product
// payment_request_log
// refund_order
// refund_order_product
func init() {
	flag.StringVar(&host, "h", "47.110.88.165:2624", "数据库地址 IP:端口")
	flag.StringVar(&user, "u", "root", "用户名")
	flag.StringVar(&pwd, "p", "J4l71RPO14", " 密码")
	flag.StringVar(&work, "w", "repo", " 模块名")
	flag.StringVar(&database, "d", "payment", "数据库名")
	flag.StringVar(&tableName, "t", "application,merchant,notify_record,router_account_statistics,statistics,payment_expire_record, event_process_record, trade_statistics,payment_account,payment_method,payment_order,payment_order_product,payment_request_log,refund_order,refund_order_product", "表名字，多个表名用逗号分隔，为空则生成全部表")
}

func main() {
	flag.Parse()

	// 连接数据库
	MySQLDSN := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&timeout=30s", user, pwd, host, database)
	db, err := gorm.Open(mysql.Open(MySQLDSN))
	if err != nil {
		panic(fmt.Errorf("cannot establish db connection: %w", err))
	}

	// 生成实例
	g := gen.NewGenerator(gen.Config{
		// 相对执行`go run`时的路径, 会自动创建目录
		OutPath:      filepath.Join(".") + "/repo/dao",
		ModelPkgPath: filepath.Join(".") + "/repo/model",
		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
		// 表字段可为 null 值时, 对应结体字段使用指针类型
		//FieldNullable: true, // generate pointer when field is nullable
		// 表字段默认值与模型结构体字段零值不一致的字段, 在插入数据时需要赋值该字段值为零值的, 结构体字段须是指针类型才能成功, 即`FieldCoverable:true`配置下生成的结构体字段.
		// 因为在插入时遇到字段为零值的会被GORM赋予默认值. 如字段`age`表默认值为10, 即使你显式设置为0最后也会被GORM设为10提交.
		// 如果该字段没有上面提到的插入时赋零值的特殊需要, 则字段为非指针类型使用起来会比较方便.
		FieldCoverable: false, // generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values
		// 模型结构体字段的数字类型的符号表示是否与表字段的一致, `false`指示都用有符号类型
		FieldSignable: false, // detect integer field's unsigned type, adjust generated data type
		// 生成 gorm 标签的字段索引属性
		FieldWithIndexTag: false, // generate with gorm index tag
		// 生成 gorm 标签的字段类型属性
		FieldWithTypeTag: true, // generate with gorm column type tag
	})

	// 设置目标 db
	g.UseDB(db)

	// 自定义字段的数据类型
	// 统一数字类型为int64,兼容protobuf
	dataMap := map[string]func(detailType gorm.ColumnType) (dataType string){
		"tinyint":   func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(detailType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(detailType gorm.ColumnType) (dataType string) { return "int64" },
	}

	// 要先于`ApplyBasic`执行
	g.WithDataTypeMap(dataMap)

	// 模型自定义选项组
	fieldOpts := []gen.ModelOpt{}

	var allModel []interface{}

	// 根据tableName参数决定生成哪些表
	if tableName == "" {
		// 如果tableName为空，生成全部表
		fmt.Println("生成全部表的model和dao...")
		allModel = g.GenerateAllTable(fieldOpts...)
	} else {
		// 如果tableName不为空，按逗号分隔并生成指定的表
		tableNames := parseTableNames(tableName)
		fmt.Printf("生成指定表的model和dao: %v\n", tableNames)

		// 验证表是否存在
		existingTables, err := getExistingTables(db)
		if err != nil {
			panic(fmt.Errorf("获取数据库表列表失败: %w", err))
		}

		// 生成指定的表
		for _, table := range tableNames {
			table = strings.TrimSpace(table)
			if table == "" {
				continue
			}

			// 检查表是否存在
			if !contains(existingTables, table) {
				fmt.Printf("警告: 表 '%s' 不存在，跳过生成\n", table)
				continue
			}

			fmt.Printf("正在生成表 '%s' 的model和dao...\n", table)
			model := g.GenerateModel(table, fieldOpts...)
			allModel = append(allModel, model)
		}

		if len(allModel) == 0 {
			fmt.Println("没有找到有效的表，程序退出")
			return
		}
	}

	// 创建模型的方法,生成文件在 query 目录
	g.ApplyBasic(allModel...)
	g.Execute()

	fmt.Println("代码生成完成！")
}

// parseTableNames 解析表名字符串，按逗号分隔
func parseTableNames(tableNameStr string) []string {
	if tableNameStr == "" {
		return []string{}
	}

	tables := strings.Split(tableNameStr, ",")
	var result []string

	for _, table := range tables {
		table = strings.TrimSpace(table)
		if table != "" {
			result = append(result, table)
		}
	}

	return result
}

// getExistingTables 获取数据库中所有存在的表名
func getExistingTables(db *gorm.DB) ([]string, error) {
	var tables []string

	// 查询所有表名
	rows, err := db.Raw("SHOW TABLES").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		if err := rows.Scan(&tableName); err != nil {
			return nil, err
		}
		tables = append(tables, tableName)
	}

	return tables, nil
}

// contains 检查切片中是否包含指定元素
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func removeAllContents(dir string) error {
	// 获取目录下所有文件和子目录
	files, err := filepath.Glob(filepath.Join(dir, "*"))
	if err != nil {
		return err
	}

	// 遍历所有文件和子目录
	for _, file := range files {
		// 删除文件或递归删除子目录
		err := os.RemoveAll(file)
		if err != nil {
			return err
		}
	}

	// 删除目录本身
	err = os.RemoveAll(dir)
	if err != nil {
		return err
	}

	return nil
}
