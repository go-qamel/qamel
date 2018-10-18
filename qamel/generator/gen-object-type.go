package generator

import "fmt"

type goTypeConverter struct {
	inC    string
	inCpp  string
	inCgo  string
	cgo2Go func(name string) string
	go2C   func(name string) string
}

var mapGoType = map[string]goTypeConverter{
	"int": goTypeConverter{
		inC:   "int",
		inCpp: "int",
		inCgo: "C.int",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("int(int32(%s))", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.int(int32(%s))", name)
		},
	},

	"int32": goTypeConverter{
		inC:   "int",
		inCpp: "int",
		inCgo: "C.int",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("int32(%s)", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.int(%s)", name)
		},
	},

	"int64": goTypeConverter{
		inC:   "long",
		inCpp: "long",
		inCgo: "C.long",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("int64(%s)", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.long(%s)", name)
		},
	},

	"float32": goTypeConverter{
		inC:   "float",
		inCpp: "float",
		inCgo: "C.float",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("float32(%s)", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.float(%s)", name)
		},
	},

	"float64": goTypeConverter{
		inC:   "double",
		inCpp: "double",
		inCgo: "C.double",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("float64(%s)", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.double(%s)", name)
		},
	},

	"bool": goTypeConverter{
		inC:   "bool",
		inCpp: "bool",
		inCgo: "C.bool",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("bool(%s)", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.bool(%s)", name)
		},
	},

	"string": goTypeConverter{
		inC:   "char*",
		inCpp: "QString",
		inCgo: "*C.char",
		cgo2Go: func(name string) string {
			return fmt.Sprintf("C.GoString(%s)", name)
		},
		go2C: func(name string) string {
			return fmt.Sprintf("C.CString(%s)", name)
		},
	},
}
