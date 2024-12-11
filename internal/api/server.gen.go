// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/oapi-codegen/oapi-codegen/v2 version v2.4.1 DO NOT EDIT.
package api

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	strictnethttp "github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

const (
	BearerScopes = "Bearer.Scopes"
)

// Error defines model for Error.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Success defines model for Success.
type Success struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BadRequest defines model for BadRequest.
type BadRequest = Error

// InternalError defines model for InternalError.
type InternalError = Error

// NotFound defines model for NotFound.
type NotFound = Error

// PingResponse defines model for PingResponse.
type PingResponse struct {
	Data *string `json:"data,omitempty"`
}

// TemporaryRedirect defines model for TemporaryRedirect.
type TemporaryRedirect = Error

// Unauthorized defines model for Unauthorized.
type Unauthorized = Error

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Проверка жизни сервиса
	// (GET /ping)
	Ping(w http.ResponseWriter, r *http.Request)
}

// Unimplemented server implementation that returns http.StatusNotImplemented for each endpoint.

type Unimplemented struct{}

// Проверка жизни сервиса
// (GET /ping)
func (_ Unimplemented) Ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandlerFunc   func(w http.ResponseWriter, r *http.Request, err error)
}

type MiddlewareFunc func(http.Handler) http.Handler

// Ping operation middleware
func (siw *ServerInterfaceWrapper) Ping(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	ctx = context.WithValue(ctx, BearerScopes, []string{})

	r = r.WithContext(ctx)

	handler := http.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		siw.Handler.Ping(w, r)
	}))

	for _, middleware := range siw.HandlerMiddlewares {
		handler = middleware(handler)
	}

	handler.ServeHTTP(w, r)
}

type UnescapedCookieParamError struct {
	ParamName string
	Err       error
}

func (e *UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter '%s'", e.ParamName)
}

func (e *UnescapedCookieParamError) Unwrap() error {
	return e.Err
}

type UnmarshalingParamError struct {
	ParamName string
	Err       error
}

func (e *UnmarshalingParamError) Error() string {
	return fmt.Sprintf("Error unmarshaling parameter %s as JSON: %s", e.ParamName, e.Err.Error())
}

func (e *UnmarshalingParamError) Unwrap() error {
	return e.Err
}

type RequiredParamError struct {
	ParamName string
}

func (e *RequiredParamError) Error() string {
	return fmt.Sprintf("Query argument %s is required, but not found", e.ParamName)
}

type RequiredHeaderError struct {
	ParamName string
	Err       error
}

func (e *RequiredHeaderError) Error() string {
	return fmt.Sprintf("Header parameter %s is required, but not found", e.ParamName)
}

func (e *RequiredHeaderError) Unwrap() error {
	return e.Err
}

type InvalidParamFormatError struct {
	ParamName string
	Err       error
}

func (e *InvalidParamFormatError) Error() string {
	return fmt.Sprintf("Invalid format for parameter %s: %s", e.ParamName, e.Err.Error())
}

func (e *InvalidParamFormatError) Unwrap() error {
	return e.Err
}

type TooManyValuesForParamError struct {
	ParamName string
	Count     int
}

func (e *TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("Expected one value for %s, got %d", e.ParamName, e.Count)
}

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{})
}

type ChiServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	Middlewares      []MiddlewareFunc
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// HandlerFromMux creates http.Handler with routing matching OpenAPI spec based on the provided mux.
func HandlerFromMux(si ServerInterface, r chi.Router) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseRouter: r,
	})
}

func HandlerFromMuxWithBaseURL(si ServerInterface, r chi.Router, baseURL string) http.Handler {
	return HandlerWithOptions(si, ChiServerOptions{
		BaseURL:    baseURL,
		BaseRouter: r,
	})
}

// HandlerWithOptions creates http.Handler with additional options
func HandlerWithOptions(si ServerInterface, options ChiServerOptions) http.Handler {
	r := options.BaseRouter

	if r == nil {
		r = chi.NewRouter()
	}
	if options.ErrorHandlerFunc == nil {
		options.ErrorHandlerFunc = func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	}
	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandlerFunc:   options.ErrorHandlerFunc,
	}

	r.Group(func(r chi.Router) {
		r.Get(options.BaseURL+"/ping", wrapper.Ping)
	})

	return r
}

type BadRequestJSONResponse Error

type InternalErrorJSONResponse Error

type NotFoundJSONResponse Error

type PingResponseJSONResponse struct {
	Data *string `json:"data,omitempty"`
}

type SuccessJSONResponse Success

type TemporaryRedirectJSONResponse Error

type UnauthorizedJSONResponse Error

type PingRequestObject struct {
}

type PingResponseObject interface {
	VisitPingResponse(w http.ResponseWriter) error
}

type Ping200JSONResponse struct{ PingResponseJSONResponse }

func (response Ping200JSONResponse) VisitPingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)

	return json.NewEncoder(w).Encode(response)
}

type Ping401JSONResponse struct{ UnauthorizedJSONResponse }

func (response Ping401JSONResponse) VisitPingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(401)

	return json.NewEncoder(w).Encode(response)
}

type Ping500JSONResponse struct{ InternalErrorJSONResponse }

func (response Ping500JSONResponse) VisitPingResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(500)

	return json.NewEncoder(w).Encode(response)
}

// StrictServerInterface represents all server handlers.
type StrictServerInterface interface {
	// Проверка жизни сервиса
	// (GET /ping)
	Ping(ctx context.Context, request PingRequestObject) (PingResponseObject, error)
}

type StrictHandlerFunc = strictnethttp.StrictHTTPHandlerFunc
type StrictMiddlewareFunc = strictnethttp.StrictHTTPMiddlewareFunc

type StrictHTTPServerOptions struct {
	RequestErrorHandlerFunc  func(w http.ResponseWriter, r *http.Request, err error)
	ResponseErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

func NewStrictHandler(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: StrictHTTPServerOptions{
		RequestErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
		ResponseErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		},
	}}
}

func NewStrictHandlerWithOptions(ssi StrictServerInterface, middlewares []StrictMiddlewareFunc, options StrictHTTPServerOptions) ServerInterface {
	return &strictHandler{ssi: ssi, middlewares: middlewares, options: options}
}

type strictHandler struct {
	ssi         StrictServerInterface
	middlewares []StrictMiddlewareFunc
	options     StrictHTTPServerOptions
}

// Ping operation middleware
func (sh *strictHandler) Ping(w http.ResponseWriter, r *http.Request) {
	var request PingRequestObject

	handler := func(ctx context.Context, w http.ResponseWriter, r *http.Request, request interface{}) (interface{}, error) {
		return sh.ssi.Ping(ctx, request.(PingRequestObject))
	}
	for _, middleware := range sh.middlewares {
		handler = middleware(handler, "Ping")
	}

	response, err := handler(r.Context(), w, r, request)

	if err != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, err)
	} else if validResponse, ok := response.(PingResponseObject); ok {
		if err := validResponse.VisitPingResponse(w); err != nil {
			sh.options.ResponseErrorHandlerFunc(w, r, err)
		}
	} else if response != nil {
		sh.options.ResponseErrorHandlerFunc(w, r, fmt.Errorf("unexpected response type: %T", response))
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/7xV32oTTxR+leX8fpdLN/XPzd4pKBRBpNWrkotx9zSZujsznZktxrKQpiBFhYIPoBe+",
	"QBqMTStpX+HMG8nMponRigq2N8vM2Z1zvvm+75zdg0yWSgoU1kC6BxqNksJg2Nxn+TruVGis32VSWBRh",
	"yZQqeMYslyLZNlL4mMm6WDK/+l/jFqTwX7JInTRvTfJAa6mhrusYcjSZ5songRTogzukCR3TGQ0jGkV0",
	"QkO6cH06d/s0hjqGNWFRC1Y0Ka4f0HuaugM3cH0a05Sm7sgdRXS+QOkxPZb2oaxEfhP80LF7Q2M6c4OI",
	"pjT2jyGd0mcPz2N5wkVnfSbfX+FRWirUljei58yGKL5kpSoQUpAvIAbbU35trOaiAx7fLCKfb2Nmr4K8",
	"UWUZGhOpcCK+3P8zri7zXcXWJ7dPFzR2hzR1b+k0CiqeuAP66t65AQ3dwCN6iqWSmuneOuZc+2tcu47z",
	"ktG8Zh3DM8Eq25Wav8L8pnttSCM3oHPXp4lvOveaJjQJCs8S+Trzpls2Sybz4LUf3BFDicawzlXv6hg0",
	"7lRc+5tuNhkW37d/8tWSb268umcBs0pz29vwbMwGIzKNgQ7u+ewiy1FDDIKV/vC9mZZBtEXvMMUfYa8R",
	"hIstGeBxG3rM26JgFqMN1Ls886B2UZtGr9WV1krL30sqFExxSOF2CMWgmO0GSEnosnQPOhhs43kK9ddy",
	"SMNsgHh5ut9qtX7loPl3ydJQqWO401r9/aElN9cx3P2TSsvj/XviId1cUL7ZrtsxmKosme55L3/0/wga",
	"0dj1Gz9/CT6e0iTyfw7XpxFN3D4NvRCsY7zslyS366aO9lyHMpUuvKDWqjRJCpmxoiuNhbpdfwsAAP//",
	"kHSZ5SoHAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %w", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %w", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	res := make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	resolvePath := PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		pathToFile := url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}