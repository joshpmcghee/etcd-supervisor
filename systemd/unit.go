package systemd

import (
	"html/template"
	"os"

	"github.com/ansel1/merry"
	"github.com/coreos/go-systemd/dbus"
	"github.com/sirupsen/logrus"
)

const(
	unitName = "supervised-etcd"
)

type UnitConfig struct {
	ExistingCluster bool
	DiscoveryURL string
	Peers []string
}

type Runner struct {
	logger *logrus.Logger
	dbus *dbus.Conn
	results chan string
	config *UnitConfig
}

func NewRunner(logger *logrus.Logger) (*Runner, error) {
	runner := &Runner{
		results: make(chan string),
		logger: logger,
		config: &UnitConfig{},
	}

	conn, err := dbus.New()
	if err != nil {
		return nil, merry.Append(err, "failed to create dbus conn")
	}
	runner.dbus = conn

	return runner, nil
}

func (r *Runner)Start() error {
	jobID, err := r.dbus.StartUnit(unitName, "fail", r.results)
	r.logger.Infof("Starting Systemd unit %v as job %v", unitName, jobID)

	result := <-r.results
	if result != "done" {
		return merry.Append(err, "failed to start systemd unit")
	}

	return nil
}

func (r *Runner)Stop() error {
	jobID, err := r.dbus.StopUnit(unitName, "fail", r.results)
	r.logger.Infof("Stopping Systemd unit %v running as job %v", unitName, jobID)

	result := <-r.results
	if result != "done" {
		return merry.Append(err, "failed to stop systemd unit")
	}
	return nil
}

func (r *Runner)Configure(config *UnitConfig) error {
	tmpl, err := template.ParseFiles("supervised-etcd.service")
	if err != nil {
		return merry.Append(err, "failed to create template from file")
	}

	file, err := os.OpenFile("/var/lib/systemd/system/supervised-etcd.service", os.O_TRUNC|os.O_RDWR, 744)
	if err != nil {
		return merry.Append(err, "failed to open file for template output")
	}

	err = tmpl.Execute(file, config)
	if err != nil {
		return merry.Append(err, "failed to execute template")
	}
	r.config = config
	return nil
}