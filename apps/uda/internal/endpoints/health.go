package endpoints

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.com/startupbuilder/startupbuilder/pkg/tracer"
)

type healthResponse struct {
	Msg string `json:"message"`
}

func MakeHealthEndpoint2(
	tracer tracer.Tracer,
) func(ctx context.Context, request interface{}) (interface{}, error) {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		var err error

		_, span := tracer.StartSpan(ctx, "health", nil)
		defer span.End(err)

		return healthEndpoint(ctx, request)
	}
}

func MakeHealthEndpoint(tracer tracer.Tracer) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// decodeRequest()
		// function()
		// encode()
		//
		// here you should call request encoder
		// then response decoder
		response, err := healthEndpoint(r.Context(), r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		respJson, err := json.Marshal(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.Write([]byte(respJson))
	}
}

func healthEndpoint(_ context.Context, _ interface{}) (interface{}, error) {
	return &healthResponse{
		Msg: "ok",
	}, nil
}
