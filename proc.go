package main

import (
	"fmt"
	"reflect"

	. "github.com/emad-elsaid/debugger/ui"
	"github.com/go-delve/delve/pkg/proc"
)

const MaxVarUINest = 4

var ProcLoadConfig = proc.LoadConfig{
	FollowPointers:     true,
	MaxStringLen:       10000,
	MaxArrayValues:     100,
	MaxStructFields:    100,
	MaxMapBuckets:      100,
	MaxVariableRecurse: 2,
}

type VarWidget interface {
	Layout(c C) D
}

func NewVarWidget(v *proc.Variable) VarWidget {
	if v == nil {
		return &UnknownVarWidget{Variable: v}
	}

	switch v.Kind {
	case reflect.String:
		return &StringVarWidget{Variable: v}
	case reflect.Struct:
		return &StructVarWidget{
			Variable:  v,
			open:      false,
			clickable: new(Clickable),
		}
	case reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Uintptr:
		return &NumberVarWidget{Variable: v}
	case reflect.Slice, reflect.Array:
		return &SliceVarWidget{
			Variable:  v,
			open:      false,
			clickable: new(Clickable),
		}
	case reflect.Bool:
		return &BoolVarWidget{Variable: v}
	case reflect.Func:
		return &FuncVarWidget{Variable: v}
	case reflect.Ptr, reflect.UnsafePointer:
		return &PtrVarWidget{Variable: v}
	case reflect.Map:
		return &MapVarWidget{
			Variable:  v,
			open:      false,
			clickable: new(Clickable),
		}
	case reflect.Complex64, reflect.Complex128:
		return &ComplexVarWidget{Variable: v}
	case reflect.Chan:
		return &ChanVarWidget{Variable: v}
	case reflect.Interface:
		return &InterfaceVarWidget{Variable: v}
	default:
		return &UnknownVarWidget{Variable: v}
	}
}

// String widget
type StringVarWidget struct {
	*proc.Variable
}

func (w *StringVarWidget) Layout(c C) D {
	var val string = "nil"

	if w.Value != nil {
		val = w.Value.ExactString()
	}
	return Label(fmt.Sprintf("%s %s = %s", w.Name, w.TypeString(), val))(c)
}

// Struct widget
type StructVarWidget struct {
	*proc.Variable
	open      bool
	clickable *Clickable
	children  []VarWidget
}

func (w *StructVarWidget) Layout(c C) D {
	if w.clickable.Clicked() {
		w.open = !w.open
	}

	if w.open {
		if w.children == nil {
			w.children = make([]VarWidget, 0, len(w.Children))

			for _, f := range w.Children {
				w.children = append(w.children, NewVarWidget(&f))
			}
		}

		fields := []FlexChild{
			Rigid(
				LayoutToWidget(
					w.clickable.Layout,
					Label(fmt.Sprintf("%s %s = {", w.Name, w.TypeString())),
				),
			),
		}

		padding := Margin(0, 0, 0, SpaceUnit*3)
		for _, c := range w.children {
			fields = append(fields, Rigid(padding(c.Layout)))
		}

		fields = append(fields, Rigid(Label("}")))

		return Rows(fields...)(c)
	} else {
		return LayoutToWidget(
			w.clickable.Layout,
			Label(fmt.Sprintf("%s %s = {...}", w.Name, w.TypeString())),
		)(c)
	}
}

// Map widget
type MapVarWidget struct {
	*proc.Variable
	open      bool
	clickable *Clickable
	keys      []VarWidget
	values    []VarWidget
}

func (w *MapVarWidget) Layout(c C) D {
	if w.clickable.Clicked() {
		w.open = !w.open
	}

	if w.open {
		if w.keys == nil || w.values == nil {
			w.keys = make([]VarWidget, 0, len(w.Children)/2)
			w.values = make([]VarWidget, 0, len(w.Children)/2)

			for i := range w.Children {
				if i%2 == 0 {
					w.keys = append(w.keys, NewVarWidget(&w.Children[i]))
				} else {
					w.values = append(w.values, NewVarWidget(&w.Children[i]))
				}
			}
		}

		fields := []FlexChild{
			Rigid(
				LayoutToWidget(
					w.clickable.Layout, Label(fmt.Sprintf("%s %s = {", w.Name, w.TypeString())),
				),
			),
		}

		for k := range w.keys {
			fields = append(fields,
				Rigid(
					Columns(
						ColSpacer3,
						Rigid(w.keys[k].Layout),
						Rigid(Label(" : ")),
						Rigid(w.values[k].Layout),
					),
				),
			)
		}

		fields = append(fields, Rigid(Label("}")))

		return Rows(fields...)(c)
	} else {
		return LayoutToWidget(
			w.clickable.Layout,
			Label(fmt.Sprintf("%s %s = {...(%d)...}", w.Name, w.TypeString(), w.Len/2)),
		)(c)
	}
}

// Numbers widget like int, float, uint...
type NumberVarWidget struct {
	*proc.Variable
}

func (v *NumberVarWidget) Layout(c C) D {
	var val string
	if v.Value == nil {
		val = "nil"
	} else {
		val = v.Value.ExactString()
	}

	return Label(fmt.Sprintf("%s %s = %s", v.Name, v.TypeString(), val))(c)
}

// Slice of values
type SliceVarWidget struct {
	*proc.Variable
	open      bool
	clickable *Clickable
	children  []VarWidget
}

func (w *SliceVarWidget) Layout(c C) D {
	if w.clickable.Clicked() {
		w.open = !w.open
	}

	if w.open {
		if w.children == nil {
			w.children = make([]VarWidget, 0, len(w.Children))

			for _, f := range w.Children {
				w.children = append(w.children, NewVarWidget(&f))
			}
		}

		fields := []FlexChild{
			Rigid(LayoutToWidget(w.clickable.Layout, Label(fmt.Sprintf("%s %s = {", w.Name, w.TypeString())))),
		}

		padding := Margin(0, 0, 0, SpaceUnit*3)
		for _, c := range w.children {
			fields = append(fields, Rigid(padding(c.Layout)))
		}

		fields = append(fields, Rigid(Label("}")))

		return Rows(fields...)(c)
	} else {
		return LayoutToWidget(w.clickable.Layout, Label(fmt.Sprintf("%s %s = {...(%d)...}", w.Name, w.TypeString(), w.Len)))(c)
	}
}

// Boolean widget
type BoolVarWidget struct {
	*proc.Variable
}

func (v *BoolVarWidget) Layout(c C) D {
	var val string
	if v.Value != nil {
		val = v.Value.ExactString()
	}
	return Label(fmt.Sprintf("%s %s = %s", v.Name, v.TypeString(), val))(c)
}

// Func widget
type FuncVarWidget struct {
	*proc.Variable
}

func (v *FuncVarWidget) Layout(c C) D {
	var val string
	if v.Value != nil {
		val = v.Value.ExactString()
	}
	return Label(fmt.Sprintf("%s %s = %s", v.Name, v.TypeString(), val))(c)
}

// Pointer widget
type PtrVarWidget struct {
	*proc.Variable
	child VarWidget
}

func (v *PtrVarWidget) Layout(c C) D {

	var val W = EmptyWidget

	if len(v.Children) > 0 {
		if v.child == nil {
			v.child = NewVarWidget(&v.Children[0])
		}
		val = v.child.Layout
	}

	return Columns(
		Rigid(Label(fmt.Sprintf("%s %s ->", v.Name, v.TypeString()))),
		Flexed(1, val),
	)(c)
}

// Complex numbers widget
type ComplexVarWidget struct {
	*proc.Variable
}

func (v *ComplexVarWidget) Layout(c C) D {
	var val string
	if v.Value == nil {
		val = "nil"
	} else {
		val = v.Value.ExactString()
	}

	return Label(fmt.Sprintf("%s %s = %s", v.Name, v.TypeString(), val))(c)
}

// Chan widget
type ChanVarWidget struct {
	*proc.Variable
}

func (v *ChanVarWidget) Layout(c C) D {
	return Label(fmt.Sprintf("%s %s", v.Name, v.TypeString()))(c)
}

// Interface widget
type InterfaceVarWidget struct {
	*proc.Variable
	child VarWidget
}

func (v *InterfaceVarWidget) Layout(c C) D {
	var val W = EmptyWidget

	if len(v.Children) > 0 {
		if v.child == nil {
			v.child = NewVarWidget(&v.Children[0])
		}
		val = v.child.Layout
	}

	return Columns(
		Rigid(Label(fmt.Sprintf("%s %s ->", v.Name, v.TypeString()))),
		Flexed(1, val),
	)(c)
}

// Rest of variables widgets
type UnknownVarWidget struct {
	*proc.Variable
}

func (v *UnknownVarWidget) Layout(c C) D {
	if v.Variable == nil {
		return Label("nil")(c)
	}

	return Label(fmt.Sprintf("Can't render '%s' type: %s", v.Name, v.Kind))(c)
}
