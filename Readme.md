TODO

Добавить возможно анализа код стайла для интерфейса на основе текущей функции
по дефолту вставлять ctx context.Context

type IInterface interface {
    Test() // ast будет смотреть на Test1 чтобы определить тип шаблона
    Test1(ctx context.Context)
}


type IInterface interface {
    Test() // ast будет смотреть на Test1 чтобы определить тип шаблона
    Test1(context.Context)
}