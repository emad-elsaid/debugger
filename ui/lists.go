package ui

import "gioui.org/layout"

func NewVerticalList() List {
	return List{
		Axis: layout.Vertical,
	}
}

func ZebraList(l *List, len int, ele layout.ListElement) W {
	return func(c C) D {
		return l.Layout(c, len, func(c C, i int) D {
			w := func(c C) D {
				return ele(c, i)
			}

			if i%2 == 0 {
				return Background(SILVER_100, w)(c)
			}

			return w(c)
		})
	}
}

type ClickableList struct {
	List
	Btns map[int]*Clickable
}

func NewClickableList() ClickableList {
	return ClickableList{
		List: NewVerticalList(),
		Btns: map[int]*Clickable{},
	}
}

func (l *ClickableList) Layout(len int, ele layout.ListElement, onClick func(i int)) W {
	for i, cl := range l.Btns {
		if cl.Clicked() {
			onClick(i)
		}
	}

	return ZebraList(&l.List, len, func(c C, i int) D {
		if _, ok := l.Btns[i]; !ok {
			l.Btns[i] = &Clickable{}
		}

		w := func(c C) D { return ele(c, i) }
		return l.Btns[i].Layout(c, w)
	})
}
