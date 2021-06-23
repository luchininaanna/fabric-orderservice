package adapter

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"orderservice/api/storeservice"
	"orderservice/pkg/common/infrastructure"
	"orderservice/pkg/order/application/adapter"
)

type storeAdapter struct {
	api storeservice.StoreServiceClient
}

func NewStoreAdapter(api storeservice.StoreServiceClient) adapter.StoreAdapter {
	return &storeAdapter{api: api}
}

func (s *storeAdapter) GetFabrics() ([]adapter.Fabric, error) {
	fabrics, err := s.api.GetFabrics(context.Background(), &empty.Empty{})
	if err != nil {
		return nil, infrastructure.InternalError(err)
	}

	var ff []adapter.Fabric
	for _, fabric := range fabrics.Fabrics {
		ff = append(ff, adapter.Fabric{
			Id:       fabric.FabricId,
			Name:     fabric.Name,
			Cost:     fabric.Cost,
			Quantity: fabric.Amount,
		})
	}

	return ff, nil
}
