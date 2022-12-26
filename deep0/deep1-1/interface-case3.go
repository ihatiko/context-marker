package main

import "go/ast"

type IServiceCase3 interface {
	Test()
	Test0(string)
	Test1(
		field1 string) error
	Test2(
		field2 string,
	)
	error
	Test3(
		field3 string,
		decl ast.GenDecl,
	) error
	Test4(field4 string)
	error
	Test5(field5, field5a, field5b string) error
	Test6(string, int, any) error
}
