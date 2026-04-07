package main

import (
	"fmt"
	"html"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syscall"

	"github.com/diamondburned/gotk4-adwaita/pkg/adw"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	. "github.com/emad-elsaid/debugger/structs"
	"github.com/go-delve/delve/pkg/proc"
)

type FileInfo struct {
	FD   int
	Name string
	Path string
	Size int64
	Type string
	Mode string
}

type SortBy int

const (
	SortByName SortBy = iota
	SortByPath
	SortBySize
	SortByType
	SortByMode
)

type OpenFilesPanel struct {
	*gtk.Box
	debugger   *Debugger
	listBox    *gtk.ListBox
	searchBox  *gtk.SearchEntry
	sortCombo  *gtk.DropDown
	files      []FileInfo
	sortBy     SortBy
	filterText string
}

func NewOpenFilesPanel(debugger *Debugger) *OpenFilesPanel {
	builder := gtk.NewBuilderFromFile("ui/openfiles.ui")

	panel := &OpenFilesPanel{
		Box:       builder.GetObject("openfiles_panel").Cast().(*gtk.Box),
		debugger:  debugger,
		listBox:   builder.GetObject("openfiles_list").Cast().(*gtk.ListBox),
		searchBox: builder.GetObject("search_entry").Cast().(*gtk.SearchEntry),
		sortCombo: builder.GetObject("sort_combo").Cast().(*gtk.DropDown),
		files:     []FileInfo{},
		sortBy:    SortByName,
	}

	panel.setupHandlers()
	panel.subscribeToEvents()

	return panel
}

func (p *OpenFilesPanel) setupHandlers() {
	p.searchBox.ConnectSearchChanged(func() {
		p.filterText = strings.ToLower(p.searchBox.Text())
		p.refresh()
	})

	p.sortCombo.ConnectActivate(func() {
		p.sortBy = SortBy(p.sortCombo.Selected())
		p.sortAndRefresh()
	})
}

func (p *OpenFilesPanel) subscribeToEvents() {
	p.debugger.Signals.Connect(EventStateChanged, func(event DebuggerEvent, data EventData) {
		if data.State != nil && !data.State.Running {
			p.updateFileList()
		}
	})
}

func (p *OpenFilesPanel) updateFileList() {
	if p.debugger.LastState == nil || p.debugger.LastState.Exited {
		return
	}

	var pid int
	p.debugger.WithTarget(func(target *proc.Target) {
		pid = target.Pid()
	})

	if pid == 0 {
		return
	}

	files := []FileInfo{}
	fdPath := fmt.Sprintf("/proc/%d/fd", pid)

	entries, err := os.ReadDir(fdPath)
	if err != nil {
		return
	}

	for _, entry := range entries {
		fdNum := 0
		fmt.Sscanf(entry.Name(), "%d", &fdNum)

		linkPath := filepath.Join(fdPath, entry.Name())
		target, err := os.Readlink(linkPath)
		if err != nil {
			continue
		}

		info := FileInfo{
			FD:   fdNum,
			Name: filepath.Base(target),
			Path: target,
		}

		var stat syscall.Stat_t
		if err := syscall.Stat(target, &stat); err == nil {
			info.Size = stat.Size
			info.Type = getFileType(stat.Mode)
			info.Mode = getFileMode(target)
		} else {
			info.Type = detectTypeFromPath(target)
			info.Mode = "unknown"
		}

		files = append(files, info)
	}

	p.files = files
	p.sortAndRefresh()
}

func (p *OpenFilesPanel) sortAndRefresh() {
	sort.Slice(p.files, func(i, j int) bool {
		switch p.sortBy {
		case SortByName:
			return p.files[i].Name < p.files[j].Name
		case SortByPath:
			return p.files[i].Path < p.files[j].Path
		case SortBySize:
			return p.files[i].Size < p.files[j].Size
		case SortByType:
			return p.files[i].Type < p.files[j].Type
		case SortByMode:
			return p.files[i].Mode < p.files[j].Mode
		default:
			return p.files[i].Name < p.files[j].Name
		}
	})
	p.refresh()
}

func (p *OpenFilesPanel) refresh() {
	for child := p.listBox.FirstChild(); child != nil; child = p.listBox.FirstChild() {
		p.listBox.Remove(child)
	}

	for _, file := range p.files {
		if p.matchesFilter(file) {
			p.addFileRow(file)
		}
	}
}

func (p *OpenFilesPanel) matchesFilter(file FileInfo) bool {
	if p.filterText == "" {
		return true
	}

	searchIn := strings.ToLower(file.Name + " " + file.Path)
	parts := strings.Fields(p.filterText)

	for _, part := range parts {
		if !strings.Contains(searchIn, part) {
			return false
		}
	}
	return true
}

func (p *OpenFilesPanel) addFileRow(file FileInfo) {
	row := adw.NewActionRow()
	row.SetTitle(html.EscapeString(file.Name))

	subtitle := fmt.Sprintf("%s • %s • %s • %d bytes",
		file.Type, file.Mode, file.Path, file.Size)
	row.SetSubtitle(html.EscapeString(subtitle))

	openDirBtn := gtk.NewButtonFromIconName("folder-open-symbolic")
	openDirBtn.SetVAlign(gtk.AlignCenter)
	openDirBtn.AddCSSClass("flat")
	openDirBtn.SetTooltipText("Open Directory")
	openDirBtn.ConnectClicked(func() {
		p.openDirectory(file.Path)
	})

	row.AddSuffix(openDirBtn)
	p.listBox.Append(row)
}

func (p *OpenFilesPanel) openDirectory(path string) {
	dir := filepath.Dir(path)
	if info, err := os.Stat(dir); err == nil && info.IsDir() {
		file := gio.NewFileForPath(dir)
		launcher := gtk.NewFileLauncher(file)
		launcher.OpenContainingFolder(nil, nil, nil)
	}
}

func getFileType(mode uint32) string {
	switch mode & syscall.S_IFMT {
	case syscall.S_IFREG:
		return "file"
	case syscall.S_IFDIR:
		return "directory"
	case syscall.S_IFLNK:
		return "symlink"
	case syscall.S_IFIFO:
		return "pipe"
	case syscall.S_IFSOCK:
		return "socket"
	case syscall.S_IFCHR:
		return "char device"
	case syscall.S_IFBLK:
		return "block device"
	default:
		return "unknown"
	}
}

func detectTypeFromPath(path string) string {
	if strings.HasPrefix(path, "pipe:") {
		return "pipe"
	}
	if strings.HasPrefix(path, "socket:") {
		return "socket"
	}
	if strings.HasPrefix(path, "anon_inode:") {
		return "special"
	}
	if strings.HasPrefix(path, "/dev/") {
		return "device"
	}
	return "unknown"
}

func getFileMode(path string) string {
	info, err := os.Stat(path)
	if err != nil {
		return "unknown"
	}

	mode := info.Mode()
	readable := mode&0444 != 0
	writable := mode&0222 != 0

	if readable && writable {
		return "read-write"
	} else if readable {
		return "read-only"
	} else if writable {
		return "write-only"
	}
	return "no-access"
}
