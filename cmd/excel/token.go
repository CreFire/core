package excel

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
)

type StructInfo struct {
	Name   string
	Fields []FieldInfo
}

type FieldInfo struct {
	Name string
	Type string
}

func Token() {
	// 指定要解析的 Go 文件
	filePath := "./cmd/excel/kvPairs.go"

	// 打开 Go 文件
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个 Go 语法分析器的文件集
	fset := token.NewFileSet()

	// 使用 Go 语法分析器解析文件
	node, err := parser.ParseFile(fset, filePath, file, parser.AllErrors)
	if err != nil {
		fmt.Println("Error parsing file:", err)
		return
	}
	structInfos := make([]StructInfo, 0)

	// 遍历文件中的所有声明并查找结构体
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					// 提取结构体的名称
					structName := typeSpec.Name.Name
					// 创建用于存储字段信息的切片
					fields := make([]FieldInfo, 0)
					if structType, isStruct := typeSpec.Type.(*ast.StructType); isStruct {
						// 提取结构体的名称
						if len(structType.Fields.List) > 0 {
							for _, field := range structType.Fields.List {
								fmt.Println(field.Tag)
								// 提取字段的名称和类型
								fieldName := field.Names[0].Name
								fieldType := types.ExprString(field.Type)

								// 创建字段信息对象并添加到字段切片
								fieldInfo := FieldInfo{
									Name: fieldName,
									Type: fieldType,
								}
								fields = append(fields, fieldInfo)
							}
						}
						fmt.Println("Found struct:", structName)
					}

					// 创建结构体信息对象并添加到结构体切片
					structInfo := StructInfo{
						Name:   structName,
						Fields: fields,
					}
					structInfos = append(structInfos, structInfo)
				}
			}
		}
	}
	// 打印提取的结构体信息
	for _, structInfo := range structInfos {
		fmt.Printf("Struct Name: %s\n", structInfo.Name)
		for _, field := range structInfo.Fields {
			fmt.Sprintf("%v", field)
		}
	}
}
