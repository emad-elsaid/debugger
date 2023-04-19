package main

import . "github.com/emad-elsaid/debugger/ui"

type MemoryPanel struct {
	list List
}

func NewMemoryPanel() MemoryPanel {
	return MemoryPanel{
		list: NewVerticalList(),
	}
}

func (m *MemoryPanel) Layout(a *Analytics) W {
	charts := []struct {
		title     string
		analytics *[]uint64
	}{
		{"Virtual Mem", &a.Vsize},
		{"Resident Mem", &a.ResidentMem},
		{"Shared Mem", &a.SharedMem},
		{"Text Mem", &a.TextMem},
		{"Data Mem", &a.DataMem},
	}

	ele := func(c C, i int) D {
		a := *charts[i].analytics
		var last uint64
		if len(a) > 0 {
			last = a[len(a)-1]
		}

		return Inset1(
			Rows(
				Rigid(
					Columns(
						Flexed(1, Label(*&charts[i].title)),
						Rigid(Label(ByteCountToDecimal(last))),
					),
				),
				Rigid(Chart(a, 100)),
			),
		)(c)
	}

	return Grid(&m.list, len(charts), 300, ele)
}
