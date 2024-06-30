package resolvers

import (
	"context"
	"errors"
	"fmt"
	"github.com/azarc-io/verathread-gateway/internal/config"
	"github.com/azarc-io/verathread-next-common/common/genericdb"
	gqlutil "github.com/azarc-io/verathread-next-common/util/gql"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidServiceResponse = errors.New("service returned an invalid response")
)

func doGenericPagedQuery[T any](
	ctx context.Context, opts *config.APIGatewayOptions, table string, query *genericdb.GenericPagedQuery,
) (T, *genericdb.PageInfo, bool) {
	var (
		db     = opts.MongoUseCase.GenericClient()
		result T
	)

	pi, err := db.PagedQuery(ctx, table, query, &result)

	if err != nil {
		log.Error().Err(err).Msgf("error caught while executing generic query")
		gqlutil.AddGeneralError(ctx, err, 500)
		return result, nil, false
	}

	if pi == nil {
		gqlutil.AddGeneralError(ctx, fmt.Errorf("something went wrong, page info was not generated"), 500)
		return result, nil, false
	}

	return result, &genericdb.PageInfo{
		Total:     pi.Total,
		Next:      pi.Next,
		Prev:      pi.Prev,
		Page:      pi.Page,
		PerPage:   pi.PerPage,
		TotalPage: pi.TotalPage,
	}, true
}
