package supervisor

import (
	"github.com/ansel1/merry"
	"github.com/joshpmcghee/etcd-supervisor/generated"
	"github.com/joshpmcghee/etcd-supervisor/systemd"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

type Service struct {
	logger *logrus.Logger
	runner *systemd.Runner
}

func NewService(logger *logrus.Logger) (*Service, error) {
	runner, err := systemd.NewRunner(logger)
	if err != nil {
		return nil, merry.Append(err, "failed to created new runner")
	}

	service := &Service{
		runner: runner,
		logger: logger,
	}

	return service, nil
}
func (s *Service) Bootstrap(ctx context.Context, in *generated.BootstrapRequest) (*generated.BootstrapResponse, error) {
	resp := &generated.BootstrapResponse{}

	err := s.runner.Configure(&systemd.UnitConfig{
		DiscoveryURL: in.DiscoveryUrl,
	})
	if err != nil {
		return resp, merry.Append(err, "failed to configure unit")
	}

	err = s.runner.Start()
	if err != nil {
		return resp, merry.Append(err, "failed to start unit")
	}

	return resp, nil
}

func (s *Service) Leave(ctx context.Context, in *generated.LeaveRequest) (*generated.LeaveResponse, error) {
	resp := &generated.LeaveResponse{}

	//TODO: Implement Etcd cluster leave

	err := s.runner.Stop()
	if err != nil {
		return resp, merry.Append(err, "failed to stop unit")
	}
	return resp, nil
}

func (s *Service) Join(ctx context.Context, in *generated.JoinRequest) (*generated.JoinResponse, error) {
	return &generated.JoinResponse{}, nil
}

func (s *Service) Upgrade(ctx context.Context, in *generated.UpgradeRequest) (*generated.UpgradeResponse, error) {
	return &generated.UpgradeResponse{}, nil
}

func (s *Service) Subjugate(ctx context.Context, in *generated.SubjugateRequest) (*generated.SubjugateResponse, error) {
	return &generated.SubjugateResponse{}, nil
}