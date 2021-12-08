package grpcerror

import (
	"context"
	"encoding/json"
	"strconv"

	"google.golang.org/grpc"

	"google.golang.org/grpc/metadata"

	"github.com/ignishub/terr"
)

const (
	debugMetadata   = "terr.debug"
	detailsMetadata = "terr.details"
	codeMetadata    = "terr.code"
)

func decodeError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return terr.InternalServerError("INTERNAL_SERVER_ERROR", err.Error())
	}
	var e terr.Error

	status := md.Get(codeMetadata)
	if len(status) == 2 {
		e.Code = status[0]
		e.HTTPStatusCode, err = strconv.Atoi(status[1])
		if err != nil {
			panic(err)
		}
	}

	debug := md.Get(debugMetadata)
	if len(debug) > 0 {
		for i := 0; i < len(debug); i = i + 2 {

			e.WithDebug(debug[i], json.RawMessage(debug[i+1]))
		}
	}

	details := md.Get(detailsMetadata)
	if len(details) > 0 {
		for i := 0; i < len(details); i = i + 2 {
			e.WithDebug(details[i], json.RawMessage(debug[i+1]))
		}
	}

	return &e
}

func encodeError(ctx context.Context, err error, details, debug bool) error {
	if err == nil {
		return nil
	}
	e := terr.From(err)
	md := make(metadata.MD)

	md.Append(codeMetadata, e.Code, strconv.Itoa(e.HTTPStatusCode))

	if details {
		for _, v := range e.Details {
			data, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			md.Append(detailsMetadata, string(data))
		}
	}

	if debug {
		for _, v := range e.Debug {
			data, err := json.Marshal(v)
			if err != nil {
				panic(err)
			}
			md.Append(debugMetadata, string(data))
		}
	}
	grpc.SetTrailer(ctx, md)
	return e
}

func UnaryServerInterceptor(details, debug bool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		response, err := handler(ctx, req)
		return response, encodeError(ctx, err, details, debug)
	}
}

func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req, reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	return decodeError(ctx, err)
}
