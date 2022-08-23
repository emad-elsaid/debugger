package main

import (
	"encoding/json"
	"errors"
	"os"
	"path"

	"gioui.org/app"
	"github.com/go-delve/delve/service/api"
	"github.com/google/uuid"
)

const CONFIG_DIR = "debugger"
const SESSIONS_FILE = "sessions.json"
const MAX_SESSIONS = 10

func ConfigPath() (string, error) {
	datadir, err := app.DataDir()
	if err != nil {
		return "", err
	}

	configDir := path.Join(datadir, CONFIG_DIR)
	if _, err := os.Stat(configDir); errors.Is(err, os.ErrNotExist) {
		os.Mkdir(configDir, 0777)
	}

	return configDir, nil
}

func SessionsFile() (string, error) {
	config, err := ConfigPath()
	if err != nil {
		return "", err
	}

	return path.Join(config, SESSIONS_FILE), nil
}

type Session struct {
	ID             string
	Path           string
	BinName        string
	Args           string
	RunImmediately bool
	Test           bool
	Breakpoints    []SessionBreakpoint
	Watches        []string
}

type SessionBreakpoint struct {
	Name    string
	File    string
	Line    int
	Enabled bool
}

func ListSessions() []Session {
	sessionsFile, err := SessionsFile()
	if err != nil {
		return []Session{}
	}

	data, err := os.ReadFile(sessionsFile)
	if err != nil {
		return []Session{}
	}

	var sessions []Session
	if err := json.Unmarshal(data, &sessions); err != nil {
		return []Session{}
	}

	return sessions
}

func SaveSession(session Session) error {
	newSessions := []Session{session}

	// add all saved sessions to the array execluding the current session
	savedSessions := ListSessions()
	for _, s := range savedSessions {
		if s.ID != session.ID {
			newSessions = append(newSessions, s)
		}
	}

	if len(newSessions) > MAX_SESSIONS {
		newSessions = newSessions[:MAX_SESSIONS]
	}

	SaveSessions(newSessions)

	return nil
}

func SaveSessions(sessions []Session) error {
	sessionsFile, err := SessionsFile()
	if err != nil {
		return err
	}

	data, err := json.Marshal(sessions)
	if err != nil {
		return err
	}

	if err := os.WriteFile(sessionsFile, data, 0777); err != nil {
		return err
	}

	return nil
}

func ToSession(d *Debugger) Session {
	if d.SessionID == "" {
		u, _ := uuid.NewUUID()
		d.SessionID = u.String()
	}

	bps := []SessionBreakpoint{}
	for _, bp := range d.Breakpoints() {
		bps = append(bps, SessionBreakpoint{
			Name:    bp.Name,
			File:    bp.File,
			Line:    bp.Line,
			Enabled: !bp.Disabled,
		})
	}

	return Session{
		ID:             d.SessionID,
		Path:           d.Path,
		BinName:        d.BinName,
		Args:           d.Args,
		RunImmediately: d.RunImmediately,
		Test:           d.Test,
		Breakpoints:    bps,
		Watches:        d.WatchesExpr,
	}
}

func FromSession(s Session) (*Debugger, error) {
	d, err := NewDebugger(s.Path, s.BinName, s.Args, false, s.Test)
	if err != nil {
		return nil, err
	}

	d.SessionID = s.ID
	for _, bp := range s.Breakpoints {
		_, err = d.CreateBreakpoint(&api.Breakpoint{
			Name:     bp.Name,
			File:     bp.File,
			Line:     bp.Line,
			Disabled: !bp.Enabled,
		})
	}

	d.WatchesExpr = s.Watches

	if s.RunImmediately {
		d.Continue()
	}
	return d, nil
}

func DeleteSession(id string) {
	sessions := ListSessions()
	newSessions := make([]Session, 0, len(sessions))
	for _, s := range sessions {
		if s.ID != id {
			newSessions = append(newSessions, s)
		}
	}

	SaveSessions(newSessions)
}
