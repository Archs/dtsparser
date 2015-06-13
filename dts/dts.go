package dts

import (
	"path/filepath"
)

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
                         toplevel object has the base name of the .d.ts file
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

	// for module/class/interface
	parent     *Object
	Classes    map[string]*Object `json:"omitempty"`
	Interfaces map[string]*Object `json:"omitempty"`
	Modules    map[string]*Object `json:"omitempty"`
	Enums      map[string]*Object `json:"omitempty"`
}

type DTS struct {
	// TopLevel object
	Object
	// current parsing ojbect
	current *Object
}

// fpath, the file path of the .d.ts file
func (d *DTS) Init(fpath string) {
	base := filepath.Base(fpath)
	// set the toplevel module name
	d.Identifier.Name = base
	d.Object.Kind = TopLevel
	d.current = &d.Object
}

func (d *DTS) newObject(kind Kind, name string) *Object {
	// create the object
	o := &Object{
		Kind: kind,
	}
	o.Identifier.Name = name
	// set parent
	if d.current != nil {
		o.parent = d.current
	} else {
		// toplevel
		o.parent = &d.Object
	}
	// add to mapping
	switch kind {
	case Module:
		if d.Modules == nil {
			d.Modules = make(map[string]*Object)
		}
		d.Modules[o.Name] = o
	case Class:
		if d.Classes == nil {
			d.Classes = make(map[string]*Object)
		}
		d.Classes[o.Name] = o
	case Interface:
		if d.Interfaces == nil {
			d.Interfaces = make(map[string]*Object)
		}
		d.Interfaces[o.Name] = o
	case Enum:
		if d.Enums == nil {
			d.Enums = make(map[string]*Object)
		}
		d.Enums[o.Name] = o
	}
	d.current = o
	return o
}

func (d *DTS) NewModule(text string) {
	println("\tmodule", text)
	d.newObject(Module, text)
}

func (d *DTS) NewClass(text string) {
	println("\tclass", text)
	d.newObject(Class, text)
}

func (d *DTS) NewInterface(text string) {
	println("\tinterface", text)
	d.newObject(Interface, text)
}

func (d *DTS) NewEnum(text string) {
	println("\tenum", text)
	d.newObject(Enum, text)
}

func (d *DTS) EndBlock(msg string) {
	println("\tend block:", msg)
	d.current = d.current.parent
}

func (d *DTS) NewVariable(text string) {
	println("variable", text)
}

func (d *DTS) NewFunction(text string) {
	println("function", text)
}
