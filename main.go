package main

import (
	"fmt"
	"github.com/samber/lo"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
	"reflect"
	"sort"
	"strings"
)

const template = "context.Context"
const templateNamed = "ctx context.Context"
const appendContextTemplate = "context"
const importContextTemplateWithComma = "import \"context\" \n\n"
const appendContextTemplateWithComma = "\"context\""

const comma byte = 44 // "," https://unicode-table.com/en/002C/
const space byte = 32 // \s https://unicode-table.com/en/0020/
func main() {
	startFolder := "test-data-implementation"
	files, err := os.ReadDir(startFolder)
	dirPath := startFolder
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		if file.IsDir() {
			dirPath = path.Join(dirPath, file.Name())
		}
	}

	processDir(startFolder)
}

func processDir(folder string) {
	files, err := os.ReadDir(folder)
	fmt.Println("Process folder:", folder)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		filePath := path.Join(folder, file.Name())
		if file.IsDir() {
			processDir(filePath)
			continue
		}
		markInterfaceRealization(folder, filePath)
		processFile(folder, filePath)
	}
}

type TokenInfo struct {
	Start   int
	End     int
	Imports []string
}

// TODO в 2 фазы 1) вычитываем все файлы и собираем реализации по пакету (текущий уровень вложенности по путям)
var InterfaceRealization = map[string][]string{}

func markInterfaceRealization(fileFolder, file string) {
	fSet := token.NewFileSet()
	node, err := parser.ParseFile(fSet, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		var specs ast.Spec
		if ok {
			specs = genDecl.Specs[0]
		} else {
			genDecl, _ := decl.(*ast.FuncDecl)
			if genDecl.Recv != nil {
				_, ok := genDecl.Recv.List[0].Type.(*ast.Ident)
				params := lo.Map[*ast.Field, string](genDecl.Type.Params.List, func(item *ast.Field, index int) string {
					return item.Type.(*ast.Ident).Name
				})
				sort.Strings(params)
				fnParams := strings.Join(params, ",")
				if ok {
					fmt.Println(fmt.Sprintf("%s(%s)", genDecl.Name, fnParams))
				} else {
					/*					starExpr, _ := genDecl.Recv.List[0].Type.(*ast.StarExpr)
										astIdent := starExpr.X.(*ast.Ident)*/
					fmt.Println(fmt.Sprintf("%s(%s)", genDecl.Name, fnParams))
				}
			}
			continue
		}

		if genDecl.Tok == token.IMPORT {
			continue
		}
		v := reflect.TypeOf(specs).Elem()
		switch v.Name() {
		case "ValueSpec":
			typeSpec, _ := genDecl.Specs[0].(*ast.ValueSpec)
			typeSpec = typeSpec
			break
		case "TypeSpec":
			typeSpec, _ := genDecl.Specs[0].(*ast.TypeSpec)
			typeSpec = typeSpec
			break
		default:
			fmt.Println("unsupported type", v.Name())
		}
	}
}

// TODO stringBuilder
func processFile(fileFolder, file string) {
	tokenInfo := TokenInfo{}
	needImport := true
	haveImports := false
	cursor := 0
	firstPosition := 0
	fileSet := token.NewFileSet()
	node, err := parser.ParseFile(fileSet, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	fileStream, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var replacingBuffer []byte

	for declIndex, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			replacingBuffer = append(replacingBuffer, fileStream[cursor:(decl.End())]...)
			cursor = int(decl.End())
			continue
		}

		if genDecl.Tok == token.IMPORT {
			haveImports = true
			tokenInfo.Start = int(genDecl.Pos())
			tokenInfo.End = int(genDecl.End())
			for _, spec := range genDecl.Specs {
				astSpec := spec.(*ast.ImportSpec)
				tokenInfo.Imports = append(tokenInfo.Imports, astSpec.Path.Value)
				if astSpec.Path.Value == fmt.Sprintf("\"%s\"", appendContextTemplate) {
					needImport = false
				}
			}
			continue
		}
		if declIndex == 0 {
			firstPosition = int(genDecl.Pos())
		}
		v := reflect.TypeOf(genDecl.Specs[0]).Elem()
		switch v.Name() {
		case "ValueSpec":
			continue
		case "TypeSpec":
			break
		default:
			fmt.Println("unsupported type", v.Name())
		}
		typeSpec := genDecl.Specs[0].(*ast.TypeSpec)

		interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		for i, t := range interfaceType.Methods.List {
			if t.Names == nil {
				continue
			}
			fName := t.Names[0]
			astInterfaceFunc := fName.Obj.Decl.(*ast.Field)
			astInterfaceFuncParams := astInterfaceFunc.Type.(*ast.FuncType).Params.List
			replacingBuffer = append(replacingBuffer, fileStream[cursor:(fName.End())]...)
			if !FindContext(astInterfaceFuncParams) {
				if needImport {
					if haveImports {
						tokenInfo.Imports = append(tokenInfo.Imports, appendContextTemplateWithComma)
						sort.Strings(tokenInfo.Imports)
						imports := strings.Join(tokenInfo.Imports, "\n\t")
						formattedImport := fmt.Sprintf("import (\n\t%s\n)", imports)
						var replacingBufferWithImport []byte
						replacingBufferWithImport = append(replacingBufferWithImport, replacingBuffer[:tokenInfo.Start-1]...)
						replacingBufferWithImport = append(replacingBufferWithImport, []byte(formattedImport)...)
						replacingBufferWithImport = append(replacingBufferWithImport, replacingBuffer[tokenInfo.End:]...)
						replacingBuffer = replacingBufferWithImport
					} else {
						var replacingBufferWithImport []byte
						replacingBufferWithImport = append(replacingBufferWithImport, fileStream[:firstPosition-1]...)
						lastPosition := len(replacingBufferWithImport)
						replacingBufferWithImport = append(replacingBufferWithImport, []byte(importContextTemplateWithComma)...)
						replacingBufferWithImport = append(replacingBufferWithImport, replacingBuffer[lastPosition:]...)
						replacingBuffer = replacingBufferWithImport
					}
					needImport = false
				}
				if IsNamedTemplate(astInterfaceFuncParams) {
					replacingBuffer = append(replacingBuffer, []byte(templateNamed)...)
					if len(astInterfaceFuncParams) > 0 {
						replacingBuffer = append(replacingBuffer, comma, space)
					}
				} else {
					replacingBuffer = append(replacingBuffer, []byte(template)...)
					if len(astInterfaceFuncParams) > 0 {
						replacingBuffer = append(replacingBuffer, comma, space)
					}
				}
			}
			replacingBuffer = append(replacingBuffer, fileStream[fName.End():(t.End())]...)
			cursor = int(t.End())
			if i == len(interfaceType.Methods.List)-1 {
				replacingBuffer = append(replacingBuffer, fileStream[cursor:(interfaceType.End())]...)
				cursor = int(interfaceType.End())
			}
		}
	}
	_, err = os.ReadDir(path.Join("gen", fileFolder))
	if err != nil {
		err = os.MkdirAll(path.Join("gen", fileFolder), os.ModePerm)
		if err != nil {
			panic(err)
		}
	}

	err = os.WriteFile(path.Join("gen", file), replacingBuffer, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func IsNamedTemplate(astInterfaceFuncParams []*ast.Field) bool {
	if len(astInterfaceFuncParams) > 0 {
		prmF := astInterfaceFuncParams[0]
		_, ok := prmF.Type.(*ast.Ident)
		if !ok {
			innerTypeSel, ok := prmF.Type.(*ast.SelectorExpr)
			if !ok {
				return true
			}
			innerType, ok := innerTypeSel.X.(*ast.Ident)
			if !ok {
				return true
			}
			fmt.Println(innerType)
		}
		if prmF.Names == nil {
			return false
		}
	}
	return true
}

func FindContext(astInterfaceFuncParams []*ast.Field) bool {
	for _, prmF := range astInterfaceFuncParams {
		_, ok := prmF.Type.(*ast.Ident)
		if !ok {
			innerTypeSel, ok := prmF.Type.(*ast.SelectorExpr)
			if !ok {
				continue
			}
			innerType, ok := innerTypeSel.X.(*ast.Ident)
			if !ok {
				continue
			}

			if fmt.Sprintf("%s.%s", innerType.Name, innerTypeSel.Sel.Name) == template {
				return true
			}
		}
	}
	return false
}
