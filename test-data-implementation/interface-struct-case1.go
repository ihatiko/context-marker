package test_data_implementation

const (
	a = iota
)

var test = map[string]string{}

const name = "name"

func TESTCASEFUNC() {

}

type InterfaceStructCase1 interface {
	DoSome(string2 string)
	DoSome2()
}

type InterfaceStructCase1Realization struct {
}

func (r InterfaceStructCase1Realization) DoSome(string2 string) {

}
func (r *InterfaceStructCase1Realization) DoSome2() {

}
func TESTCASEFUNC2() {

}
