package main

import (
	"context"
	"go/ast"
)

type IServiceCase3 interface {
	Test(ctx context.Context)
	Test0(context.Context, string)
	Test1(ctx context.Context,
		field1 string) error
	Test2(ctx context.Context,
		field2 string,
	)
	error
	Test3(ctx context.Context,
		field3 string,
		decl ast.GenDecl,
	) error
	Test4(ctx context.Context, field4 string)
	error
	Test5(ctx context.Context, field5, field5a, field5b string) error
	Test6(context.Context, string, int, any) error
}
