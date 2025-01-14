package errorpresenter

import (
	"context"
	"errors"
	"fmt"

	"github.com/kyma-incubator/compass/components/director/pkg/apperrors"

	"github.com/kyma-incubator/compass/components/director/pkg/log"

	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// UUIDService missing godoc
type UUIDService interface {
	Generate() string
}

type presenter struct {
	uuidService UUIDService
}

// NewPresenter missing godoc
func NewPresenter(service UUIDService) *presenter {
	return &presenter{uuidService: service}
}

// Do missing godoc
func (p *presenter) Do(ctx context.Context, err error) *gqlerror.Error {
	customErr := apperrors.Error{}
	errID := p.uuidService.Generate()

	if found := errors.As(err, &customErr); !found {
		log.C(ctx).WithField("errorID", errID).WithError(err).Errorf("Unknown error: %v", err)
		return newGraphqlErrorResponse(ctx, apperrors.InternalError, "Internal Server Error [errorID=%s]", errID)
	}

	if apperrors.ErrorCode(customErr) == apperrors.InternalError {
		log.C(ctx).WithField("errorID", errID).WithError(err).Errorf("Internal Server Error: %v", err)
		return newGraphqlErrorResponse(ctx, apperrors.InternalError, "Internal Server Error [errorID=%s]", errID)
	}

	log.C(ctx).WithField("errorID", errID).WithError(err).Error()
	return newGraphqlErrorResponse(ctx, apperrors.ErrorCode(customErr), customErr.Error())
}

func newGraphqlErrorResponse(ctx context.Context, errCode apperrors.ErrorType, msg string, args ...interface{}) *gqlerror.Error {
	return &gqlerror.Error{
		Message:    fmt.Sprintf(msg, args...),
		Path:       graphql.GetFieldContext(ctx).Path(),
		Extensions: map[string]interface{}{"error_code": errCode, "error": errCode.String()},
	}
}
