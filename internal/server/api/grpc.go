package api

import (
	"context"

	"github.com/HAGIT4/go-middle/internal/server/service"
	"github.com/HAGIT4/go-middle/internal/server/storage"
	"github.com/HAGIT4/go-middle/pb"
	"github.com/HAGIT4/go-middle/pkg/models"
)

type metricGrpcHandler struct {
	pb.UnimplementedMetricServiceServer
	sv service.MetricServiceInterface
}

func newGrpcMetricRouter(sv service.MetricServiceInterface, st storage.StorageInterface) (r *metricGrpcHandler, err error) {
	r = &metricGrpcHandler{
		sv: sv,
	}
	return r, nil
}

func (h *metricGrpcHandler) Update(ctx context.Context, in *pb.UpdateReq) (resp *pb.UpdateResp, err error) {
	req := &models.Metrics{
		ID:    in.ID,
		MType: in.Type,
		Delta: &in.Delta,
		Value: &in.Value,
		Hash:  in.Hash,
	}
	err = h.sv.UpdateMetric(req)
	if err != nil {
		return nil, err
	}

	resp = &pb.UpdateResp{}
	return resp, nil
}
