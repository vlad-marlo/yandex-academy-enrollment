package production

import (
	"errors"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/model"
	"github.com/vlad-marlo/yandex-academy-enrollment/internal/pkg/fielderr"
)

var (
	ErrNilReference   = errors.New("unexpectedly reference to nil object")
	ErrNotImplemented = fielderr.New("not implemented", model.BadRequestResponse{}, fielderr.CodeBadRequest)
	ErrBadRequest     = fielderr.New("bad request", model.BadRequestResponse{}, fielderr.CodeBadRequest)
	ErrNotFound       = fielderr.New("not found", model.BadRequestResponse{}, fielderr.CodeNotFound)
)
