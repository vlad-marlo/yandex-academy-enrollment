package production

import (
	"context"
	"github.com/asaskevich/govalidator"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"strconv"
)

func (srv *Service) GetCourierByID(ctx context.Context, id string) (courier *model.CourierDTO, err error) {
	var courierID int
	courierID, err = strconv.Atoi(id)
	if err != nil {
		return nil, ErrBadRequest
	}
	courier, err = srv.storage.GetCourierByID(ctx, courierID)
	if err != nil {
		return nil, ErrNotFound
	}
	return courier, nil
}

func (srv *Service) CreateCouriers(ctx context.Context, req *model.CreateCourierRequest) (resp *model.CouriersCreateResponse, err error) {
	if ok, _ := govalidator.ValidateStruct(req); !ok || req == nil {
		return nil, ErrBadRequest
	}
	var couriers []model.CourierDTO
	couriers, err = srv.storage.CreateCouriers(ctx, req.Couriers)
	if err != nil {
		return nil, ErrBadRequest
	}
	if len(couriers) == 0 || len(req.Couriers) == 0 {
		return nil, ErrBadRequest
	}
	return &model.CouriersCreateResponse{Couriers: couriers}, nil
}

func (srv *Service) GetCouriers(ctx context.Context, opts model.PaginationOpts) (*model.GetCouriersResponse, error) {
	couriers, err := srv.storage.GetCouriers(ctx, opts.Limit(), opts.Offset())
	if err != nil {
		return nil, ErrBadRequest
	}
	return &model.GetCouriersResponse{
		Couriers: couriers,
		Limit:    opts.Limit(),
		Offset:   opts.Offset(),
	}, nil
}

func (srv *Service) GetCourierMetaInfo(ctx context.Context, req *model.GetCourierMetaInfoRequest) (resp *model.GetCourierMetaInfoResponse, err error) {
	_ = ctx.Err()
	_ = req.CourierID
	return nil, ErrNotImplemented
}
