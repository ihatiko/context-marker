package main

import (
	"context"
	"fmt"
)

type IService interface {
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
	) error
	Test4(ctx context.Context, field4 string)
	error
	Test5(ctx context.Context, field5, field5a, field5b string) error
	Test6(context.Context, string, int, any, fmt.Formatter) error
	Test7(string, int, any, context.Context) error
	Test8(ctx context.Context) error
}

type IService2 interface {
	Test(ctx context.Context)
	Test2(context.Context, string)
	Test3(context.Context, string)
}

type Test struct {
	Field1 int
	Field2 int
	Field3 int
}

func (t Test) SomeFunc() {

}

type IService3 interface {
	Test(ctx context.Context)
	Test2(context.Context, string)
	Test3(context.Context, string)
}

type Test2 struct {
	Field1 int
	Field2 int
	Field3 int
}

func (t Test2) SomeFunc() {

}
