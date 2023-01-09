package main

import (
	"context"
	"fmt"
)

type IService interface {
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
	) error
	Test4(field4 string)
	error
	Test5(field5, field5a, field5b string) error
	Test6(string, int, any, fmt.Formatter) error
	Test7(string, int, any, context.Context) error
	Test8(ctx context.Context) error
}

type IService2 interface {
	Test()
	Test2(string)
	Test3(string)
}

type Test struct {
	Field1 int
	Field2 int
	Field3 int
}

func (t Test) SomeFunc() {

}

type IService3 interface {
	Test()
	Test2(string)
	Test3(string)
}

type Test3 interface {
	SomeFunc()
}

type Test2 struct {
	Field1 int
	Field2 int
	Field3 int
}

func (t Test2) SomeFunc() {

}
