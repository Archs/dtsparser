package dts

type Kind int

const (
	TopLevel Kind = iota
	Module
	Class
	Interface
	Enum
	Obj
)

type Identifier struct {
	Name     string
	Modifier []string
}

/*
KNOWN type: string -> string
            number -> float64
            object -> js.M
            funcion   -> func

            class  -> struct
            module -> struct
            interface -> struct
            toplevel  -> convert the same name package
                         toplevel object can have no Indentifier
// var/func overrides -> comment
// other types just output as comment
*/
type Variable struct {
	Identifier
	// only convert the KNOWN type to go
	Type []string // multipy types, return type
}

// is this usefull?
type Assignment struct {
	Identifier
	Value string // just a string
}

type Function struct {
	Identifier
	Args       []Variable
	ReturnType []string
	// Type is the function return type
}

type Object struct {
	Identifier
	// module/class/interface/enum/js object/ or top level
	Kind Kind

	// for class/interface/object
	Extents    map[string]*Object `json:"omitempty"`
	Implements map[string]*Object `json:"omitempty"`
	// var difinitions
	Vars        map[string]*Variable   `json:"omitempty"`
	Assignments map[string]*Assignment `json:"omitempty"`
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

func newObject(kind Kind, name string) *Object {
	o := &Object{
		Kind: kind,
	}
	o.Identifier.Name = name
	return o
}

type DTS struct {
	// TopLevel object
	Object
	// current parsing ojbect
	current *Object
}

func (d DTS) Init() {
	d.Object.Kind = TopLevel
}

func (d DTS) NewModule(text string) {
	println("module", text)
	d.current = newObject(Module, text)
	if d.Modules == nil {
		d.Modules = make(map[string]*Object)
	}
	d.Modules[d.current.Name] = d.current
}

func (d DTS) NewClass(text string) {
	println("class", text)
}

func (d DTS) NewInterface(text string) {
	println("interface", text)

}

func (d DTS) NewEnum(text string) {
	println("enum", text)
}

func (d DTS) NewVariable(text string) {
	println("variable", text)
}

func (d DTS) NewFunction(text string) {
	println("function", text)
}
