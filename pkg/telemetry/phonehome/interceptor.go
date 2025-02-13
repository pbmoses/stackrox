package phonehome

import (
	"context"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	erroxGRPC "github.com/stackrox/rox/pkg/errox/grpc"
	"github.com/stackrox/rox/pkg/grpc/authn"
	grpcError "github.com/stackrox/rox/pkg/grpc/errors"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"google.golang.org/grpc"
)

const grpcGatewayUserAgentHeader = runtime.MetadataPrefix + "User-Agent"

func (cfg *Config) track(rp *RequestParams) {
	id := cfg.HashUserAuthID(rp.UserID)
	cfg.interceptorsLock.RLock()
	defer cfg.interceptorsLock.RUnlock()
	for event, funcs := range cfg.interceptors {
		props := map[string]any{}
		ok := true
		for _, interceptor := range funcs {
			if ok = interceptor(rp, props); !ok {
				break
			}
		}
		if ok {
			cfg.telemeter.Track(event, id, props)
		}
	}
}

func getUserAgent[getter func(string) []string](headers getter) string {
	// By default, all permanent HTTP headers in grpc-gateway are added grpcgateway- prefix:
	// https://github.com/grpc-ecosystem/grpc-gateway/blob/8952e38d5addd28308e29c272c696a578aa8ace8/runtime/mux.go#L106-L114
	// User-Agent header is occupied with internal grpc-go value:
	// https://github.com/grpc/grpc-go/blob/0238b6e1cec37b55820b461d3d30652c54efe2c4/clientconn.go#L211-L215
	userAgentValues := headers(grpcGatewayUserAgentHeader)
	// If endpoint is accessed not via grpc-gateway, extract from User-Agent header.
	// If endpoint is accessed via grpc-gateway, append grpc-go value to the resultinguser agent.
	userAgentValues = append(userAgentValues, headers("User-Agent")...)
	return strings.Join(userAgentValues, " ")
}

func getGRPCRequestDetails(ctx context.Context, err error, info *grpc.UnaryServerInfo, req any) *RequestParams {
	id, iderr := authn.IdentityFromContext(ctx)
	if iderr != nil {
		log.Debug("Cannot identify user from context: ", iderr)
	}

	// Use the wrapped HTTP request details if provided:
	ri := requestinfo.FromContext(ctx)
	if ri.HTTPRequest != nil && ri.HTTPRequest.URL != nil {
		return &RequestParams{
			UserAgent: getUserAgent(ri.Metadata.Get),
			UserID:    id,
			Method:    ri.HTTPRequest.Method,
			Path:      ri.HTTPRequest.URL.Path,
			Code:      grpcError.ErrToHTTPStatus(err),
			GRPCReq:   req,
		}
	}

	return &RequestParams{
		UserAgent: getUserAgent(ri.Metadata.Get),
		UserID:    id,
		Method:    info.FullMethod,
		Path:      info.FullMethod,
		Code:      int(erroxGRPC.RoxErrorToGRPCCode(err)),
		GRPCReq:   req,
	}
}

func getHTTPRequestDetails(ctx context.Context, r *http.Request, err error) *RequestParams {
	id, iderr := authn.IdentityFromContext(ctx)
	if iderr != nil {
		log.Debug("Cannot identify user from context: ", iderr)
	}

	return &RequestParams{
		UserAgent: getUserAgent(r.Header.Values),
		UserID:    id,
		Method:    r.Method,
		Path:      r.URL.Path,
		Code:      grpcError.ErrToHTTPStatus(err),
		HTTPReq:   r,
	}
}
