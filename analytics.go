package main

import (
	"fmt"
	"os"
)

const ANALYTICS_LIMIT = 1000

type Analytics struct {
	OpenFiles   []int
	Threads     []uint
	Vsize       []uint64
	ResidentMem []uint64
	SharedMem   []uint64
	TextMem     []uint64
	DataMem     []uint64
	PID         int
}

func NewAnalytics() Analytics {
	return Analytics{
		OpenFiles:   make([]int, 0, ANALYTICS_LIMIT),
		Threads:     make([]uint, 0, ANALYTICS_LIMIT),
		Vsize:       make([]uint64, 0, ANALYTICS_LIMIT),
		ResidentMem: make([]uint64, 0, ANALYTICS_LIMIT),
		SharedMem:   make([]uint64, 0, ANALYTICS_LIMIT),
		TextMem:     make([]uint64, 0, ANALYTICS_LIMIT),
		DataMem:     make([]uint64, 0, ANALYTICS_LIMIT),
	}
}

func (a *Analytics) Collect(pid int) {
	if pid != a.PID {
		a.PID = pid
		a.Clear()
	}
	a.openFiles(pid)
	a.stat(pid)
	a.statm(pid)
}

func (a *Analytics) Clear() {
	a.OpenFiles = make([]int, 0, ANALYTICS_LIMIT)
	a.Threads = make([]uint, 0, ANALYTICS_LIMIT)
	a.Vsize = make([]uint64, 0, ANALYTICS_LIMIT)
	a.ResidentMem = make([]uint64, 0, ANALYTICS_LIMIT)
	a.SharedMem = make([]uint64, 0, ANALYTICS_LIMIT)
	a.TextMem = make([]uint64, 0, ANALYTICS_LIMIT)
	a.DataMem = make([]uint64, 0, ANALYTICS_LIMIT)
}

func (a *Analytics) openFiles(pid int) {
	f, err := os.Open(fmt.Sprintf("/proc/%d/fd", pid))
	if err != nil {
		return
	}
	defer f.Close()

	dirs, _ := f.ReadDir(-1)
	limitedAppend(&a.OpenFiles, len(dirs))
}

// https://man.archlinux.org/man/proc.5.en
func (a *Analytics) stat(pid int) {
	s, err := os.ReadFile(fmt.Sprintf("/proc/%d/stat", pid))
	if err != nil {
		return
	}

	var (
		comm                                                     string
		state                                                    rune
		ppid, pgrp, session, tty_nr, tpgid                       int
		flags                                                    uint
		minflt, cminflt, majflt, cmajflt, utime, stime           uint64
		cutime, cstime, priority, nice, num_threads, itrealvalue uint
		starttime, vsize, rss                                    uint64
		rsslim, startcode, endcode, startstack, kstkesp          uint64
		kstkeip, signal, blocked, sigignore, sigcatch            uint64
		wchan, nswap, cnswap                                     uint64
		exit_signal, processor                                   int
		rt_priority, policy                                      uint
		delayacct_blkio_ticks, guest_time, cguest_time           uint64
		start_data, end_data, start_brk, arg_start               uint64
		arg_end, env_start, env_end                              uint64
		exit_code                                                int
	)

	fmt.Sscanf(
		string(s),
		""+
			"%d %s %c %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d %d %d %d "+
			"%d %d",
		&pid, &comm, &state, &ppid, &pgrp,
		&session, &tty_nr, &tpgid, &flags, &minflt,
		&cminflt, &majflt, &cmajflt, &utime, &stime,
		&cutime, &cstime, &priority, &nice, &num_threads,
		&itrealvalue, &starttime, &vsize, &rss, &rsslim,
		&startcode, &endcode, &startstack, &kstkesp, &kstkeip,
		&signal, &blocked, &sigignore, &sigcatch, &wchan,
		&nswap, &cnswap, &exit_signal, &processor, &rt_priority,
		&policy, &delayacct_blkio_ticks, &guest_time, &cguest_time, &start_data,
		&end_data, &start_brk, &arg_start, &arg_end, &env_start,
		&env_end, &exit_code,
	)

	limitedAppend(&a.Threads, num_threads)
	limitedAppend(&a.Vsize, vsize)
}

func (a *Analytics) statm(pid int) {
	s, err := os.ReadFile(fmt.Sprintf("/proc/%d/statm", pid))
	if err != nil {
		return
	}

	var size, resident, shared, text, lib, data, dt uint64
	fmt.Sscanf(
		string(s),
		"%d %d %d %d %d %d %d",
		&size, &resident, &shared, &text, &lib, &data, &dt,
	)

	limitedAppend(&a.ResidentMem, resident)
	limitedAppend(&a.SharedMem, shared)
	limitedAppend(&a.TextMem, text)
	limitedAppend(&a.DataMem, data)
}
