// Package agent Agent metric collector
package agent

import (
	"crypto/rsa"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/HAGIT4/go-middle/pb"
	"github.com/HAGIT4/go-middle/pkg/agent/config"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type agent struct {
	serverAddr     string
	pollInterval   time.Duration
	reportInterval time.Duration
	httpClient     *http.Client
	grpcClient     pb.MetricServiceClient
	hashKey        string
	batch          bool
	grpcPort       int
	logger         *zerolog.Logger
	publicKey      *rsa.PublicKey
	localIPstring  string
}

var _ AgentInterface = (*agent)(nil)

func NewAgent(cfg *config.AgentConfig) (a *agent, err error) {
	httpClient := &http.Client{}

	grpcAddr := fmt.Sprintf(":%d", cfg.GrpcPort)
	grpcConn, err := grpc.Dial(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	grpcClient := pb.NewMetricServiceClient(grpcConn)

	logger, err := NewAgentLogger()
	if err != nil {
		return nil, err
	}

	pub, err := a.GetPublicKeyFromPem(cfg.CryptoKey)
	if err != nil {
		return nil, err
	}

	conn, err := net.Dial("tcp", cfg.ServerAddr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localIP := conn.LocalAddr()
	fmt.Println(localIP)

	a = &agent{
		serverAddr:     cfg.ServerAddr,
		pollInterval:   cfg.PollInterval,
		reportInterval: cfg.ReportInterval,
		httpClient:     httpClient,
		grpcClient:     grpcClient,
		hashKey:        cfg.HashKey,
		batch:          cfg.Batch,
		grpcPort:       cfg.GrpcPort,
		logger:         logger,
		publicKey:      pub,
		localIPstring:  localIP.String(),
	}
	return a, nil
}

func (a *agent) CollectMetrics() *agentData {
	data := newAgentData()
	return data
}
