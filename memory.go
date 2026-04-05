package main

import (
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/go-delve/delve/pkg/proc"
)

type MemoryPanel struct {
	*gtk.Box
	addrEntry   *gtk.Entry
	textView    *gtk.TextView
	textBuffer  *gtk.TextBuffer
	debugger    *Debugger
	bytesPerRow int
}

func NewMemoryPanel(debugger *Debugger) *MemoryPanel {
	builder := gtk.NewBuilderFromFile("ui/memory.ui")

	panel := &MemoryPanel{
		Box:         builder.GetObject("memory_panel").Cast().(*gtk.Box),
		debugger:    debugger,
		bytesPerRow: 16,
		addrEntry:   builder.GetObject("address_entry").Cast().(*gtk.Entry),
		textView:    builder.GetObject("memory_text_view").Cast().(*gtk.TextView),
	}

	panel.textBuffer = panel.textView.Buffer()

	// Connect button and entry activation
	goButton := builder.GetObject("go_button").Cast().(*gtk.Button)
	goButton.ConnectClicked(panel.onGoClicked)
	panel.addrEntry.ConnectActivate(panel.onGoClicked)

	return panel
}

func (p *MemoryPanel) onGoClicked() {
	addrText := p.addrEntry.Text()
	addrText = strings.TrimSpace(addrText)

	if addrText == "" {
		p.displayError("Please enter a memory address")
		return
	}

	// Parse address (support hex with or without 0x prefix)
	addrText = strings.TrimPrefix(addrText, "0x")
	addrText = strings.TrimPrefix(addrText, "0X")

	addr, err := strconv.ParseUint(addrText, 16, 64)
	if err != nil {
		p.displayError(fmt.Sprintf("Invalid address: %v", err))
		return
	}

	// Read 256 bytes of memory
	const bytesToRead = 256
	data := make([]byte, bytesToRead)

	p.debugger.WithTarget(func(target *proc.Target) {
		if ok, _ := target.Valid(); !ok {
			p.displayError("No process attached")
			return
		}

		_, err = target.Memory().ReadMemory(data, addr)
		if err != nil {
			p.displayError(fmt.Sprintf("Failed to read memory: %v", err))
			return
		}

		p.displayMemory(addr, data)
	})
}

func (p *MemoryPanel) displayMemory(startAddr uint64, data []byte) {
	var builder strings.Builder

	for i := 0; i < len(data); i += p.bytesPerRow {
		addr := startAddr + uint64(i)
		builder.WriteString(fmt.Sprintf("0x%016X  ", addr))

		// Hex bytes
		rowData := data[i:]
		if len(rowData) > p.bytesPerRow {
			rowData = rowData[:p.bytesPerRow]
		}

		hexDump := hex.EncodeToString(rowData)
		for j := 0; j < len(hexDump); j += 2 {
			builder.WriteString(hexDump[j : j+2])
			if j == 14 {
				builder.WriteString("  ")
			} else {
				builder.WriteString(" ")
			}
		}

		// Pad if incomplete row
		for j := len(rowData); j < p.bytesPerRow; j++ {
			builder.WriteString("   ")
		}

		builder.WriteString(" |")

		// ASCII representation
		for _, b := range rowData {
			if b >= 32 && b <= 126 {
				builder.WriteByte(b)
			} else {
				builder.WriteString(".")
			}
		}
		builder.WriteString("|")

		builder.WriteString("\n")
	}

	p.textBuffer.SetText(builder.String())
}

func (p *MemoryPanel) displayError(message string) {
	p.textBuffer.SetText(fmt.Sprintf("Error: %s\n", message))
}

func (p *MemoryPanel) Update() {
	// Memory panel updates on demand via Go button
}
