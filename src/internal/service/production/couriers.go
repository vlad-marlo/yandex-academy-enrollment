package production

import (
	"context"
	"errors"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/store"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/datetime"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
	"go.uber.org/zap"
	"strconv"
)

// GetCourierByID returns courier with id.
//
// Always return only nil or fielderr errors.
// It is preferred to check is error fielderr.Error after receiving non-nil error.
func (srv *Service) GetCourierByID(ctx context.Context, id string) (courier *model.CourierDTO, err error) {
	var courierID int64

	courierID, err = strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrBadRequest.With(zap.NamedError("strconv_error", err))
	}

	courier, err = srv.storage.GetCourierByID(ctx, courierID)
	if err != nil {
		return nil, ErrNotFound
	}

	return courier, nil
}

// CreateCouriers creates couriers from request.
//
// Return slice of created users with created ids.
func (srv *Service) CreateCouriers(ctx context.Context, req *model.CreateCourierRequest) (resp *model.CouriersCreateResponse, err error) {
	if !req.Valid() {
		return nil, ErrBadRequest
	}

	var couriers []model.CourierDTO
	couriers, err = srv.storage.CreateCouriers(ctx, req.Couriers)
	if err != nil {
		return nil, ErrBadRequest.With(zap.NamedError("store_error", err))
	}
	return &model.CouriersCreateResponse{Couriers: couriers}, nil
}

// GetCouriers return couriers with pagination options.
//
// If there are no couriers found by pagination opts then will be returned
// empty slice of couriers.
func (srv *Service) GetCouriers(ctx context.Context, opts model.PaginationOpts) (*model.GetCouriersResponse, error) {
	if opts == nil {
		return nil, ErrBadRequest
	}
	couriers, err := srv.storage.GetCouriers(ctx, opts.Limit(), opts.Offset())
	if err != nil {
		if !errors.Is(err, store.ErrNoContent) {
			return nil, ErrBadRequest
		}
		couriers = []model.CourierDTO{}
	}
	return &model.GetCouriersResponse{
		Couriers: couriers,
		Limit:    opts.Limit(),
		Offset:   opts.Offset(),
	}, nil
}

// GetCourierMetaInfo return meta info of courier.
//
// This method calculates courier's earned money and rating.
// It requires all fields of request to be set. If there are no completed orders in time span then rating and earned money
// will not be calculated.
func (srv *Service) GetCourierMetaInfo(ctx context.Context, req *model.GetCourierMetaInfoRequest) (resp *model.GetCourierMetaInfoResponse, err error) {
	if req == nil {
		return nil, ErrNoContent
	}

	var (
		start, end *datetime.Date
		courier    *model.CourierDTO
		orders     []int32
	)

	resp = &model.GetCourierMetaInfoResponse{
		CourierID: req.CourierID,
	}

	if start, err = datetime.ParseDate(req.StartDate); err != nil {
		return nil, ErrNoContent.WithData(resp).With(zap.NamedError("datetime_error", err))
	}

	if end, err = datetime.ParseDate(req.EndDate); err != nil {
		return nil, ErrNoContent.WithData(resp).WithData(zap.NamedError("datetime_error", err))
	}

	courier, err = srv.storage.GetCourierByID(ctx, req.CourierID)
	if err != nil {
		return nil, ErrNoContent.WithData(resp).With(zap.NamedError("storage_error", err))
	}

	resp.Regions = courier.Regions
	resp.CourierType = courier.CourierType
	resp.WorkingHours = courier.WorkingHours

	orders, err = srv.storage.GetCompletedOrdersPriceByCourier(ctx, courier.CourierID, start.Start(), end.End())
	if err != nil {
		return nil, ErrNoContent.WithData(resp).With(zap.NamedError("storage_error", err))
	}
	if len(orders) == 0 {
		return resp, nil
	}
	for _, price := range orders {
		resp.Earnings += price
	}
	resp.Rating = int32((float64(len(orders)) / end.Start().Sub(start.Start()).Hours()) * float64(courier.RatingConst()))
	resp.Earnings *= courier.EarningsConst()

	return resp, nil
}
