package main

import "context"

type IServiceCase4StyleCase1 interface {
	Test1(ctx context.Context)
	Test2(ctx context.Context)
	Test3(ctx context.Context)
	Test0(context.Context, string)
}

type IServiceCase4StyleCase2 interface {
	Test1(ctx context.Context)
	Test2(ctx context.Context)
	Test3(ctx context.Context)
	Test0(ctx context.Context, t1 string)
}
