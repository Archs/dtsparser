package dts

import (
	"encoding/json"
	"path/filepath"
	"regexp"
	"strings"
)

type Kind string

const (
	TopLevel  Kind = "TopLevel"
	Module         = "Module"
	Class          = "Class"
	Interface      = "Interface"
	Enum           = "Enum"
	Obj            = "Obj"
)

var (
	// modifiers = "declare|export|private|static|import|var|function"
	modifiers = regexp.MustCompile("declare|export|private|static|import|var|function")
	space     = regexp.MustCompile("[ \r\n\t]+")
	or        = regexp.MustCompile("\\|")
	and       = regexp.MustCompile(",")
)

type Identifier struct {
	Name     string
	Modifier []string `json:",omitempty"`
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
	Extents    []string `json:",omitempty"` // ids
	Implements []string `json:",omitempty"` // ids
	// var difinitions
	Vars        map[string]*Variable   `json:",omitempty"`
	Assignments map[string]*Assignment `json:",omitempty"`
	// using slice here, incase of function override
	Funcs []*Function `json:",omitempty"`
	// helpers
	// constructor for class
	Constructor *Function `json:",omitempty"`

	// for module/class/interface
	parent     *Object
	Classes    map[string]*Object `json:",omitempty"`
	Interfaces map[string]*Object `json:",omitempty"`
	Modules    map[string]*Object `json:",omitempty"`
	Enums      map[string]*Object `json:",omitempty"`
	Objs       map[string]*Object `json:",omitempty"`
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

func sepBy(text string, exp *regexp.Regexp) []string {
	text = strings.TrimSpace(text)
	ret := []string{}
	ss := exp.Split(text, -1)
	for _, s := range ss {
		s = strings.TrimSpace(s)
		if s != "" {
			ret = append(ret, s)
		}
	}
	return ret
}

func (d *DTS) NewBlock(modifiers string, kind Kind) {
	// create the object
	o := &Object{
		Kind: kind,
	}
	modifiers = strings.TrimSpace(modifiers)
	o.Modifier = sepBy(modifiers, space)
	// set parent
	if d.current != nil {
		o.parent = d.current
	} else {
		// toplevel
		o.parent = &d.Object
	}
	d.current = o
}

func (d *DTS) SetBlockID(name string) {
	o := d.current
	parent := o.parent
	o.Name = name
	// add to mapping
	switch o.Kind {
	case Module:
		if parent.Modules == nil {
			parent.Modules = make(map[string]*Object)
		}
		parent.Modules[o.Name] = o
	case Class:
		if parent.Classes == nil {
			parent.Classes = make(map[string]*Object)
		}
		parent.Classes[o.Name] = o
	case Interface:
		if parent.Interfaces == nil {
			parent.Interfaces = make(map[string]*Object)
		}
		parent.Interfaces[o.Name] = o
	case Enum:
		if parent.Enums == nil {
			parent.Enums = make(map[string]*Object)
		}
		parent.Enums[o.Name] = o
	case Obj:
		if parent.Objs == nil {
			parent.Objs = make(map[string]*Object)
		}
		parent.Objs[o.Name] = o
	}
}

func (d *DTS) Extends(text string) {
	d.current.Extents = sepBy(text, and)
}

func (d *DTS) Implements(text string) {
	d.current.Implements = sepBy(text, and)
}

func (d *DTS) EndBlock(msg string) {
	d.current = d.current.parent
}

func (d *DTS) NewVariable(text string) {
	d.v = new(Variable)
	text = strings.TrimSpace(text)
	d.v.Modifier = sepBy(text, space)
}

func (d *DTS) VSetIdentifier(text string) {
	if strings.HasSuffix(text, "?") {
		d.v.IsOptional = true
		text = text[:len(text)-1]
	}
	d.v.Name = text
}

func (d *DTS) VSetType(text string) {
	d.v.Type = sepBy(text, or)
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
	d.f.Modifier = sepBy(text, space)
}

func (d *DTS) FSetIdentifier(text string) {
	d.f.Name = text
}

func (d *DTS) FSetType(text string) {
	d.f.ReturnType = sepBy(text, or)
}

func (d *DTS) NewArg(text string) {
	d.v = new(Variable)
	text = strings.TrimSpace(text)
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
