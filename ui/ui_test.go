package main

import (
	"os"
	"testing"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// Initialize GTK for testing
	gtk.Init()
	adw.Init()
	os.Exit(m.Run())
}

func TestWindowUIFile(t *testing.T) {
	t.Run("loads without error", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("window.ui")
		require.NotNil(t, builder, "builder should not be nil")
	})

	t.Run("has main window", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("window.ui")
		obj := builder.GetObject("main_window")
		require.NotNil(t, obj, "main_window object should exist")

		win := obj.Cast().(*adw.ApplicationWindow)
		require.NotNil(t, win, "should cast to AdwApplicationWindow")
		require.Equal(t, "Go Debugger", win.Title(), "window title should be 'Go Debugger'")
	})

	t.Run("has header bar", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("window.ui")
		obj := builder.GetObject("header_bar")
		require.NotNil(t, obj, "header_bar object should exist")

		headerBar := obj.Cast().(*adw.HeaderBar)
		require.NotNil(t, headerBar, "should cast to AdwHeaderBar")
	})

	t.Run("has toolbar buttons", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("window.ui")

		buttons := []struct {
			id      string
			icon    string
			tooltip string
		}{
			{"continue_btn", "media-playback-start-symbolic", "Continue (F5)"},
			{"pause_btn", "media-playback-pause-symbolic", "Pause"},
			{"restart_btn", "view-refresh-symbolic", "Restart (Ctrl+Shift+F5)"},
			{"step_btn", "go-down-symbolic", "Step Into (F11)"},
			{"next_btn", "go-next-symbolic", "Step Over (F10)"},
			{"step_out_btn", "go-up-symbolic", "Step Out (Shift+F11)"},
		}

		for _, btn := range buttons {
			t.Run(btn.id, func(t *testing.T) {
				obj := builder.GetObject(btn.id)
				require.NotNil(t, obj, "%s should exist", btn.id)

				button := obj.Cast().(*gtk.Button)
				require.NotNil(t, button, "%s should cast to GtkButton", btn.id)
				require.Equal(t, btn.icon, button.IconName(), "icon should match")
				require.Equal(t, btn.tooltip, button.TooltipText(), "tooltip should match")
			})
		}
	})

	t.Run("has paned layout", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("window.ui")

		mainPane := builder.GetObject("main_pane")
		require.NotNil(t, mainPane, "main_pane should exist")
		pane := mainPane.Cast().(*gtk.Paned)
		require.NotNil(t, pane, "should cast to GtkPaned")
		require.Equal(t, gtk.OrientationHorizontal, pane.Orientation(), "should be horizontal")
		require.Equal(t, 900, pane.Position(), "position should be 900")

		bottomPane := builder.GetObject("bottom_pane")
		require.NotNil(t, bottomPane, "bottom_pane should exist")
		bPane := bottomPane.Cast().(*gtk.Paned)
		require.NotNil(t, bPane, "should cast to GtkPaned")
		require.Equal(t, gtk.OrientationVertical, bPane.Orientation(), "should be vertical")
		require.Equal(t, 600, bPane.Position(), "position should be 600")
	})

	t.Run("has tab infrastructure", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("window.ui")

		tabView := builder.GetObject("tab_view")
		require.NotNil(t, tabView, "tab_view should exist")
		tv := tabView.Cast().(*adw.TabView)
		require.NotNil(t, tv, "should cast to AdwTabView")

		tabBar := builder.GetObject("tab_bar")
		require.NotNil(t, tabBar, "tab_bar should exist")
		tb := tabBar.Cast().(*adw.TabBar)
		require.NotNil(t, tb, "should cast to AdwTabBar")
	})
}

func TestSidebarUIFile(t *testing.T) {
	t.Run("loads without error", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("sidebar.ui")
		require.NotNil(t, builder, "builder should not be nil")
	})

	t.Run("has sidebar box", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("sidebar.ui")
		obj := builder.GetObject("sidebar_box")
		require.NotNil(t, obj, "sidebar_box should exist")

		box := obj.Cast().(*gtk.Box)
		require.NotNil(t, box, "should cast to GtkBox")
		require.Equal(t, gtk.OrientationVertical, box.Orientation(), "should be vertical")
		require.Equal(t, 12, box.Spacing(), "spacing should be 12")
	})

	t.Run("has watches section", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("sidebar.ui")

		group := builder.GetObject("watches_group")
		require.NotNil(t, group, "watches_group should exist")
		wg := group.Cast().(*adw.PreferencesGroup)
		require.NotNil(t, wg, "should cast to AdwPreferencesGroup")
		require.Equal(t, "Watches", wg.Title(), "title should be 'Watches'")

		list := builder.GetObject("watches_list")
		require.NotNil(t, list, "watches_list should exist")
		wl := list.Cast().(*gtk.ListBox)
		require.NotNil(t, wl, "should cast to GtkListBox")
		require.Equal(t, gtk.SelectionNone, wl.SelectionMode(), "selection mode should be none")
	})

	t.Run("has breakpoints section", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("sidebar.ui")

		group := builder.GetObject("breakpoints_group")
		require.NotNil(t, group, "breakpoints_group should exist")
		bg := group.Cast().(*adw.PreferencesGroup)
		require.NotNil(t, bg, "should cast to AdwPreferencesGroup")
		require.Equal(t, "Breakpoints", bg.Title(), "title should be 'Breakpoints'")

		list := builder.GetObject("breakpoints_list")
		require.NotNil(t, list, "breakpoints_list should exist")
		bl := list.Cast().(*gtk.ListBox)
		require.NotNil(t, bl, "should cast to GtkListBox")
		require.Equal(t, gtk.SelectionNone, bl.SelectionMode(), "selection mode should be none")
	})
}

func TestConsoleUIFile(t *testing.T) {
	t.Run("loads without error", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("console.ui")
		require.NotNil(t, builder, "builder should not be nil")
	})

	t.Run("has console panel", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("console.ui")
		obj := builder.GetObject("console_panel")
		require.NotNil(t, obj, "console_panel should exist")

		panel := obj.Cast().(*gtk.ScrolledWindow)
		require.NotNil(t, panel, "should cast to GtkScrolledWindow")
		require.True(t, panel.VExpand(), "should expand vertically")
		require.True(t, panel.HExpand(), "should expand horizontally")
	})

	t.Run("has console text view", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("console.ui")
		obj := builder.GetObject("console_text_view")
		require.NotNil(t, obj, "console_text_view should exist")

		textView := obj.Cast().(*gtk.TextView)
		require.NotNil(t, textView, "should cast to GtkTextView")
		require.False(t, textView.Editable(), "should not be editable")
		require.Equal(t, gtk.WrapWordChar, textView.WrapMode(), "wrap mode should be word-char")
		require.True(t, textView.Monospace(), "should be monospace")
		require.Equal(t, 8, textView.MarginStart(), "margin-start should be 8")
		require.Equal(t, 8, textView.MarginEnd(), "margin-end should be 8")
		require.Equal(t, 8, textView.MarginTop(), "margin-top should be 8")
		require.Equal(t, 8, textView.MarginBottom(), "margin-bottom should be 8")
	})
}

func TestTypesUIFile(t *testing.T) {
	t.Run("loads without error", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("types.ui")
		require.NotNil(t, builder, "builder should not be nil")
	})

	t.Run("has types panel", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("types.ui")
		obj := builder.GetObject("types_panel")
		require.NotNil(t, obj, "types_panel should exist")

		panel := obj.Cast().(*gtk.Box)
		require.NotNil(t, panel, "should cast to GtkBox")
		require.Equal(t, gtk.OrientationVertical, panel.Orientation(), "should be vertical")
		require.Equal(t, 6, panel.Spacing(), "spacing should be 6")
		require.Equal(t, 6, panel.MarginStart(), "margin-start should be 6")
		require.Equal(t, 6, panel.MarginEnd(), "margin-end should be 6")
		require.Equal(t, 6, panel.MarginTop(), "margin-top should be 6")
		require.Equal(t, 6, panel.MarginBottom(), "margin-bottom should be 6")
	})

	t.Run("has search entry", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("types.ui")
		obj := builder.GetObject("types_search_entry")
		require.NotNil(t, obj, "types_search_entry should exist")

		entry := obj.Cast().(*gtk.SearchEntry)
		require.NotNil(t, entry, "should cast to GtkSearchEntry")
		require.Equal(t, "Search types...", entry.PlaceholderText(), "placeholder should match")
	})

	t.Run("has list box", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("types.ui")
		obj := builder.GetObject("types_list")
		require.NotNil(t, obj, "types_list should exist")

		listBox := obj.Cast().(*gtk.ListBox)
		require.NotNil(t, listBox, "should cast to GtkListBox")
		require.Equal(t, gtk.SelectionSingle, listBox.SelectionMode(), "selection mode should be single")
	})

	t.Run("has scrolled window", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("types.ui")
		obj := builder.GetObject("types_scrolled")
		require.NotNil(t, obj, "types_scrolled should exist")

		scrolled := obj.Cast().(*gtk.ScrolledWindow)
		require.NotNil(t, scrolled, "should cast to GtkScrolledWindow")
		require.True(t, scrolled.VExpand(), "should expand vertically")
	})
}

func TestStackUIFile(t *testing.T) {
	t.Run("loads without error", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("stack.ui")
		require.NotNil(t, builder, "builder should not be nil")
	})

	t.Run("has stack panel", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("stack.ui")
		obj := builder.GetObject("stack_panel")
		require.NotNil(t, obj, "stack_panel should exist")

		panel := obj.Cast().(*gtk.ScrolledWindow)
		require.NotNil(t, panel, "should cast to GtkScrolledWindow")
		require.True(t, panel.VExpand(), "should expand vertically")
		require.True(t, panel.HExpand(), "should expand horizontally")
	})

	t.Run("has stack list", func(t *testing.T) {
		builder := gtk.NewBuilderFromFile("stack.ui")
		obj := builder.GetObject("stack_list")
		require.NotNil(t, obj, "stack_list should exist")

		listBox := obj.Cast().(*gtk.ListBox)
		require.NotNil(t, listBox, "should cast to GtkListBox")
		require.Equal(t, gtk.SelectionSingle, listBox.SelectionMode(), "selection mode should be single")
	})
}

func TestAllUIFilesExist(t *testing.T) {
	files := []string{
		"window.ui",
		"sidebar.ui",
		"console.ui",
		"types.ui",
		"stack.ui",
		"memory.ui",
		"sources.ui",
		"functions.ui",
		"openfiles.ui",
		"disassembly.ui",
		"packages.ui",
		"breakpoints.ui",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			_, err := os.Stat(file)
			require.NoError(t, err, "%s should exist", file)
		})
	}
}

func TestAllUIFilesLoadWithoutError(t *testing.T) {
	files := []string{
		"window.ui",
		"sidebar.ui",
		"console.ui",
		"types.ui",
		"stack.ui",
		"memory.ui",
		"sources.ui",
		"functions.ui",
		"openfiles.ui",
		"disassembly.ui",
		"packages.ui",
		"breakpoints.ui",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			builder := gtk.NewBuilderFromFile(file)
			require.NotNil(t, builder, "%s should load without error", file)
		})
	}
}
