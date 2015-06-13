package dts

type Variable struct {
	Identifier string
	Modifier   []string
	Type       []string // multipy types, return type
}

type Function struct {
	Variable
	Args []Variable
	// Type is the function return type
}

type Object struct {
	// module/class/interface/enum/js object
	Type string

	// for class/interface/object
	Identifier string
	Modifier   []string           `json:"omitempty"`
	Extens     map[string]*Object `json:"omitempty"`
	Implements map[string]*Object `json:"omitempty"`
	// var difinitions
	Vars map[string]*Variable `json:"omitempty"`
	// using slice here, incase of function override
	Funcs []*Function `json:"omitempty"`
	// helpers
	// constructor for class
	Constructor *Function

	// for module
	Classes    map[string]*Object `json:"omitempty"`
	Interfaces map[string]*Object `json:"omitempty"`
	Modules    map[string]*Object `json:"omitempty"`
}

type Class struct {
	Identifier string
	Modifier   []string
	Extens     map[string]*Class
	Implements map[string]*Interface
	Vars       map[string]*Variable
	Funcs      []*Function // using slice here, incase of function override
	// helpers
	Constructor *Function
}

type Interface struct {
	Class
}

type Module struct {
	Identifier string
	Modifier   []string
	Classes    map[string]*Class
	SubModules map[string]*Module
	Vars       map[string]*Variable
	Funcs      []*Function // using slice here, incase of function override
}

type DTS struct {
	// type register
	Modules    map[string]*Module
	Classes    map[string]*Class
	Interfaces map[string]*Interface
	// for parsing
	currentModule    *Module
	currentClass     *Class
	currentInterface *Interface
}

func (d *DTS) Init() {
	d.Modules = make(map[string]*Module)
	d.Classes = make(map[string]*Class)
	d.Interfaces = make(map[string]*Interface)
}

func (d *DTS) NewModule(text string) {
	println("module", text)
}

func (d *DTS) NewClass(text string) {
	println("class", text)
}

func (d *DTS) NewInterface(text string) {
	println("interface", text)

}

func (d *DTS) NewEnum(text string) {
	println("enum", text)
}

func (d *DTS) NewVariable(text string) {
	println("variable", text)
}

func (d *DTS) NewFunction(text string) {
	println("function", text)
}
