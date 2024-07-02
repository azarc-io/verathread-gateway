// Package http handles GraphQL HTTP Requests including WebSocket Upgrades.
package http

import (
	"bytes"
	"net/http"

	"github.com/wundergraph/graphql-go-tools/v2/pkg/graphqlerrors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/wundergraph/graphql-go-tools/execution/engine"
	"github.com/wundergraph/graphql-go-tools/execution/graphql"
	"github.com/wundergraph/graphql-go-tools/v2/pkg/engine/resolve"
)

const (
	httpHeaderContentType          string = "Content-Type"
	httpContentTypeApplicationJSON string = "application/json"
	bufferSize                            = 4096
)

func (g *GraphQLHTTPRequestHandler) handleHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	var gqlRequest graphql.Request
	if err = graphql.UnmarshalHttpRequest(r, &gqlRequest); err != nil {
		g.log.Error().Err(err).Msgf("UnmarshalHttpRequest")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s := trace.SpanFromContext(r.Context())
	s.SetAttributes(attribute.String("operation", gqlRequest.OperationName))

	var opts []engine.ExecutionOptions

	if g.enableART {
		tracingOpts := resolve.TraceOptions{
			Enable:                                 false,
			ExcludePlannerStats:                    false,
			ExcludeRawInputData:                    true,
			ExcludeInput:                           true,
			ExcludeOutput:                          true,
			ExcludeLoadStats:                       false,
			EnablePredictableDebugTimings:          false,
			ExcludeParseStats:                      true,
			ExcludeValidateStats:                   false,
			IncludeTraceOutputInResponseExtensions: true,
		}

		opts = append(opts, engine.WithRequestTraceOptions(tracingOpts))
	}

	buf := bytes.NewBuffer(make([]byte, 0, bufferSize))
	resultWriter := graphql.NewEngineResultWriterFromBuffer(buf)
	if err = g.engine.Execute(r.Context(), &gqlRequest, &resultWriter, opts...); err != nil {
		g.log.Error().Err(err).Msgf("engine.Execute")

		errs := graphqlerrors.RequestErrorsFromError(err)
		if _, err = errs.WriteResponse(w); err != nil {
			g.log.Error().Err(err).Msgf("write response")
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		return
	}

	w.Header().Add(httpHeaderContentType, httpContentTypeApplicationJSON)
	w.WriteHeader(http.StatusOK)

	if _, err = w.Write(buf.Bytes()); err != nil {
		g.log.Error().Err(err).Msgf("write response")
		return
	}
}
