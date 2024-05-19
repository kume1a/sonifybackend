package shared

import (
	"net/http"
)

type HttpHandlerFunc func(*http.Request, *ApiConfig) (*HttpRes, error)

func WrapHandlerFunc(
	apiCfg *ApiConfig,
	h HttpHandlerFunc,
	middlewares ...func(http.HandlerFunc) http.HandlerFunc,
) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res, err := h(r, apiCfg)

		if err != nil {
			httpErr, ok := err.(HttpError)

			if ok {
				ResHttpError(w, &httpErr)
				return
			}

			ResHttpError(w, InternalServerError(ErrInternal))
			return
		}

		ResJson(w, res)
	})
}
