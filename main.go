package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path"
)

const template = "context.Context"
const templateNamed = "ctx context.Context"
const appendContextTemplate = "context"
const importContextTemplateWithComma = "import \"context\" \n\n"
const appendContextTemplateWithComma = "\"context\""

const comma byte = 44 // "," https://unicode-table.com/en/002C/
const space byte = 32 // \s https://unicode-table.com/en/0020/
func main() {
	startFolder := "deep0"
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
	//processFile("deep0/interface-case1.go")
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
		processFile(folder, filePath)
	}
}

type TokenInfo struct {
	Start int
	End   int
}

func processFile(fileFolder, file string) {
	tokenInfo := TokenInfo{}
	needImport := true
	haveImports := false
	cursor := 0
	firstPosition := 0
	firstPosition = firstPosition
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
	if err != nil {
		log.Fatal(err)
	}
	fileStream, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	var replacingBuffer []byte

	for declIndex, decl := range node.Decls {
		d1, ok := decl.(*ast.GenDecl)
		if !ok {
			replacingBuffer = append(replacingBuffer, fileStream[cursor:(decl.End())]...)
			cursor = int(decl.End())
			continue
		}

		if d1.Tok == token.IMPORT {
			haveImports = true
			tokenInfo.Start = int(d1.Pos())
			tokenInfo.End = int(d1.End())
			for _, spec := range d1.Specs {
				astSpec := spec.(*ast.ImportSpec)
				if astSpec.Path.Value == fmt.Sprintf("\"%s\"", appendContextTemplate) {
					needImport = false
				}
			}
			continue
		}
		if declIndex == 0 {
			firstPosition = int(d1.Pos())
		}
		d2 := d1.Specs[0].(*ast.TypeSpec)

		d3, ok := d2.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}
		for i, t := range d3.Methods.List {
			if t.Names == nil {
				continue
			}
			fName := t.Names[0]
			astInterfaceFunc := fName.Obj.Decl.(*ast.Field)
			astInterfaceFuncParams := astInterfaceFunc.Type.(*ast.FuncType).Params.List
			replacingBuffer = append(replacingBuffer, fileStream[cursor:(fName.End())]...)
			cursor = int(fName.End())
			if !FindContext(astInterfaceFuncParams) {
				if needImport {
					if haveImports {

					} else {
						var replacingBufferWithImport []byte
						replacingBufferWithImport = append(replacingBufferWithImport, fileStream[:firstPosition-1]...)
						lastPosition := len(replacingBufferWithImport)
						replacingBufferWithImport = append(replacingBufferWithImport, []byte(importContextTemplateWithComma)...)
						cursor += len(importContextTemplateWithComma)
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
					cursor = int(fName.End())
					if i == len(d3.Methods.List)-1 {
						replacingBuffer = append(replacingBuffer, fileStream[cursor:(d3.End())]...)
					}
					continue
				}
				replacingBuffer = append(replacingBuffer, []byte(template)...)
				if len(astInterfaceFuncParams) > 0 {
					replacingBuffer = append(replacingBuffer, comma, space)
				}
				cursor = int(fName.End())
			}
			if i == len(d3.Methods.List)-1 {
				replacingBuffer = append(replacingBuffer, fileStream[cursor:(d3.End())]...)
				cursor = int(d3.End())
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
