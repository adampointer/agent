package proc

import (
	"context"

	"github.com/adampointer/agent/internal/app/agent/sensors"
	"github.com/pkg/errors"
	"github.com/prometheus/procfs"
)

const (
	procName = "algod"
	label    = "procfs"
)

// PidFinder is a function which takes in a process name and returns its PID
type PidFinder func(processName string) (int, error)

// DefaultPidFinder returns a PidFinder which uses an initialised procfs.FS to find the PID of the given process
func DefaultPidFinder(fs procfs.FS) PidFinder {
	return func(processName string) (int, error) {
		procs, err := fs.AllProcs()
		if err != nil {
			return 0, err
		}
		for _, p := range procs {
			stat, err := p.Stat()
			if err != nil {
				return 0, err
			}
			if stat.Comm == processName {
				return stat.PID, nil
			}
		}
		return 0, errors.Errorf("pid not found for process '%s'", processName)
	}
}

// Sensor implements sensors.Sensor and samples from the selected stats from the proc filesystem
type Sensor struct {
	procFSPath string
	fs         procfs.FS
	pidFinder  PidFinder
}

// NewSensorFromPath creates a new Sensor from a path to the proc filesystem. Usually this is /proc but can be
// overriden to a fixtured copy for testing
// Bu default, DefaultPidFinder is used to determine the PID of the algod service
func NewSensorFromPath(path string) (*Sensor, error) {
	fs, err := procfs.NewFS(path)
	if err != nil {
		return nil, errors.Wrap(err, "initialise proc fs")
	}
	sensor := &Sensor{
		procFSPath: path,
		fs:         fs,
		pidFinder:  DefaultPidFinder(fs),
	}
	return sensor, nil
}

// WithPidFinder replaces DefaultPidFinder
func (s *Sensor) WithPidFinder(f PidFinder) *Sensor {
	s.pidFinder = f
	return s
}

// DoScan performs the scan and returns a sensors.SnapshotEvent or an error
func (s *Sensor) DoScan(_ context.Context) (*sensors.SnapshotEvent, error) {
	payload := make(map[string]interface{})

	pid, err := s.pidFinder(procName)
	if err != nil {
		return nil, err
	}
	payload["algod_pid"] = pid

	proc, err := s.fs.Proc(pid)
	if err != nil {
		return nil, err
	}

	stat, err := proc.Stat()
	if err != nil {
		return nil, err
	}
	payload["utime"] = stat.UTime
	payload["stime"] = stat.STime
	payload["starttime"] = stat.Starttime

	netstat, err := proc.Netstat()
	if err != nil {
		return nil, err
	}
	payload["InOctets"] = netstat.InOctets
	payload["OutOctets"] = netstat.OutOctets

	snmp, err := proc.Snmp()
	if err != nil {
		return nil, err
	}
	payload["InErrs"] = snmp.InErrs
	payload["OutRsts"] = snmp.OutRsts

	sstat, err := s.fs.Stat()
	if err != nil {
		return nil, err
	}
	payload["uptime"] = sstat.BootTime

	cpu, err := s.fs.PSIStatsForResource("cpu")
	if err != nil {
		return nil, err
	}
	payload["pressure_cpu"] = cpu.Some

	mem, err := s.fs.PSIStatsForResource("memory")
	if err != nil {
		return nil, err
	}
	payload["pressure_mem"] = mem.Some

	return sensors.NewSnapshotEvent(label, payload), nil
}
