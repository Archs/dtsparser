package dts

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
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

var (
	// modifiers = "declare|export|private|static|import|var|function"
	modifiers = regexp.MustCompile("declare|export|private|static|import|var|function")
	space     = regexp.MustCompile("[ \r\n\t]+")
)

type Identifier struct {
	Name     string
	Modifier []string
	// text holding the ide
	Text string `json:"-"`
}

func (i *Identifier) isPrivate() bool {
	for _, v := range i.Modifier {
		if v == "private" {
			return true
		}
	}
	return false
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
	// use for func type
	IsOptional bool
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
	Args       []*Variable
	ReturnType []string
	// Type is the function return type
}

type Object struct {
	Identifier
	// module/class/interface/enum/js object/ or top level
	Kind Kind

	// for class/interface/object
	Extents    map[string]*Object
	Implements map[string]*Object
	// var difinitions
	Vars        map[string]*Variable
	Assignments map[string]*Assignment
	// using slice here, incase of function override
	Funcs []*Function
	// helpers
	// constructor for class
	Constructor *Function

	// for module/class/interface
	parent     *Object
	Classes    map[string]*Object
	Interfaces map[string]*Object
	Modules    map[string]*Object
	Enums      map[string]*Object
}

type DTS struct {
	// TopLevel object
	Object
	// current parsing ojbect
	current *Object
	// current variable
	v *Variable
	// current function
	f *Function
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
	d.newObject(Module, text)
}

func (d *DTS) NewClass(text string) {
	d.newObject(Class, text)
}

func (d *DTS) NewInterface(text string) {
	d.newObject(Interface, text)
}

func (d *DTS) NewEnum(text string) {
	d.newObject(Enum, text)
}

func (d *DTS) EndBlock(msg string) {
	d.current = d.current.parent
}

func (d *DTS) NewVariable(text string) {
	d.v = new(Variable)
	d.v.Modifier = space.Split(text, -1)
}

func (d *DTS) VSetIdentifier(text string) {
	if strings.HasSuffix(text, "?") {
		d.v.IsOptional = true
		text = text[:len(text)-1]
	}
	d.v.Name = text
}

func (d *DTS) VSetType(text string) {
	d.v.Type = strings.Split(text, "|")
}

func (d *DTS) EndVariable(text string) {
	d.v.Text = text
	if d.current.Vars == nil {
		d.current.Vars = make(map[string]*Variable)
	}
	d.current.Vars[d.v.Name] = d.v
	d.v = nil
}

func (d *DTS) NewFunction(text string) {
	d.f = new(Function)
	d.f.Modifier = space.Split(text, -1)
}

func (d *DTS) FSetIdentifier(text string) {
	d.f.Name = text
}

func (d *DTS) FSetType(text string) {
	d.f.ReturnType = strings.Split(text, "|")
}

func (d *DTS) NewArg(text string) {
	d.v = new(Variable)
	if strings.HasSuffix(text, "?") {
		d.v.IsOptional = true
		text = text[:len(text)-1]
	}
	d.v.Name = text
}

func (d *DTS) EndArg(text string) {
	d.v.Text = text
	if d.f.Args == nil {
		d.f.Args = make([]*Variable, 0)
	}
	d.f.Args = append(d.f.Args, d.v)
	d.v = nil
}

func (d *DTS) EndFunction(text string) {
	d.f.Text = text
	if d.current.Funcs == nil {
		d.current.Funcs = make([]*Function, 0)
	}
	d.current.Funcs = append(d.current.Funcs, d.f)
	if d.f.Name == "constructor" {
		d.current.Constructor = d.f
	}
	d.f = nil
}

func (d *DTS) Show(text string) {
	println(text)
}

func (d *DTS) Json() (string, error) {
	dat, err := json.Marshal(d.Object)
	if err != nil {
		return "", err
	}
	return string(dat), nil
}
