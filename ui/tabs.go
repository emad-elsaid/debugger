package ui

var (
	SelectedTabBgColor = MixColor(ACCENT_COLOR_100, WHITE, 20)
	SelectedTabFgColor = ACCENT_COLOR_500
)

type TabChild struct {
	Name  string
	Panel W
}

type Tabs struct {
	List      List
	Active    string
	Clickable map[string]*Clickable
}

func (l *Tabs) Layout(tabs ...*TabChild) W {
	if l.Clickable == nil {
		l.Clickable = map[string]*Clickable{}
	}

	var panel = EmptyWidget
	for i, t := range tabs {
		active := l.Active == t.Name || l.Active == "" && i == 0
		if active {
			panel = t.Panel
		}
	}

	ele := func(c C, i int) D {
		t := tabs[i]
		var b *Clickable
		var ok bool
		if b, ok = l.Clickable[t.Name]; !ok {
			b = &Clickable{}
			l.Clickable[t.Name] = b
		} else if b.Clicked() {
			l.Active = t.Name
		}

		active := l.Active == t.Name || l.Active == "" && i == 0
		return Columns(
			Rigid(TabButton(b, active, t.Name)),
			Rigid(VR(1)),
		)(c)
	}

	return Rows(
		Rigid(HR(1)),
		Rigid(func(c C) D { return l.List.Layout(c, len(tabs), ele) }),
		Rigid(HR(1)),
		Flexed(1, panel),
	)
}

func TabButton(cl *Clickable, active bool, l string) W {
	var b W
	if cl.Hovered() && !active {
		b = Background(CardColor, Inset1(Label(l)))
	} else if active {
		b = Background(SelectedTabBgColor, Wrap(Label(l), Inset1, TextColor(SelectedTabFgColor)))
	} else {
		b = Inset1(Label(l))
	}

	return func(c C) D {
		return cl.Layout(c, b)
	}
}
