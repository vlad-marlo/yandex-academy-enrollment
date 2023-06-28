package production

import (
	"errors"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/fielderr"
	"github.com/vlad-marlo/yandex-academy-enrollment/pkg/model"
)

var (
	ErrNilReference   = errors.New("unexpectedly reference to nil object")
	ErrNotImplemented = fielderr.New("not implemented", model.BadRequestResponse{}, fielderr.CodeBadRequest)
	ErrBadRequest     = fielderr.New("bad request", model.BadRequestResponse{}, fielderr.CodeBadRequest)
	ErrNotFound       = fielderr.New("not found", model.BadRequestResponse{}, fielderr.CodeNotFound)
	ErrNoContent      = fielderr.New("no content to return", model.GetCourierMetaInfoResponse{}, fielderr.CodeOK)
)
