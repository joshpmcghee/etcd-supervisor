package main

import (
	"net"

	"github.com/joshpmcghee/etcd-supervisor/generated"
	"github.com/joshpmcghee/etcd-supervisor/supervisor"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "superetcd <subcommand>",
	}

	serverCmd := &cobra.Command{
		Use: "supervise",
		Run: supervise,
	}

	rootCmd.AddCommand(serverCmd)

	rootCmd.Execute()
}

func supervise(_ *cobra.Command, _ []string) {
	s := grpc.NewServer()
	supervisor := &supervisor.Service{}
	generated.RegisterSupervisorServiceServer(s, supervisor)

	lis, err := net.Listen("tcp", "0.0.0.0:8484")
	if err != nil {
		logrus.WithError(err).Error("failed to listen on port")
	}

	err = s.Serve(lis)
	if err != nil {
		logrus.WithError(err).Error("failed to start gRPC server")
	}
}
