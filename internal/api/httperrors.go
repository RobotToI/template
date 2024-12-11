package api

import (
	"fmt"
	"net/http"
	"net/url"
	"runtime"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-pkgz/rest"
)

// All error codes for UI mapping and translation
const (
	// ErrInternal is a generic internal error
	ErrInternal = iota // any internal error
	// ErrDecode is a generic decode error
	ErrDecode // failed to unmarshal incoming request
	// ErrAssetNotFound is a generic not found error
	ErrAssetNotFound // requested file not found
	// ErrNoAccess is a generic access error
	ErrNoAccess // rejected by auth
	// ErrUnauthorization is a generic unauthorization error
	ErrUnauthorization // rejected by auth
	// ErrActionRejected is a generic error for rejected actions
	ErrActionRejected // general error for rejected actions
	// ErrQuery is a generic error for rejected actions
	ErrQuery
	// ErrValidator is a generic error for rejected actions
	ErrValidator
)

// SendErrorJSON makes {error: blah, details: blah} json body and responds with error code
func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, details string, errCode int) {
	// ошибка логируется в ResponseLoggerMiddleware
	render.Status(r, httpStatusCode)

	errorResponse := BaseResponse{
		Data: nil,
		Error: &BaseError{
			Code:    stringPtr(fmt.Sprintf("%d", errCode)),
			Debug:   stringPtr(details),
			Message: stringPtr(err.Error()),
		},
	}

	render.JSON(w, r, errorResponse)
}

// stringPtr is a helper function to create a pointer to a string
func stringPtr(s string) *string {
	return &s
}

// RenderUnauthorized renders 401 error
func RenderUnauthorized(w http.ResponseWriter, r *http.Request) {
	render.Status(r, http.StatusUnauthorized)
	render.JSON(w, r, rest.JSON{"code": http.StatusUnauthorized, "data": "Unauthorized."})
}

// RenderSuccessCreated renders 201 success
func RenderSuccessCreated(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, rest.JSON{"code": http.StatusCreated, "data": data})
}

// ErrDetailsMsg old detail msg
func ErrDetailsMsg(r *http.Request, httpStatusCode int, err error, details string, errCode int) string {
	uinfoStr := ""

	q := r.URL.String()
	if qun, e := url.QueryUnescape(q); e == nil {
		q = qun
	}

	srcFileInfo := ""
	if pc, file, line, ok := runtime.Caller(2); ok {
		fnameElems := strings.Split(file, "/")
		funcNameElems := strings.Split(runtime.FuncForPC(pc).Name(), "/")
		srcFileInfo = fmt.Sprintf(" [caused by %s:%d %s]", strings.Join(fnameElems[len(fnameElems)-3:], "/"),
			line, funcNameElems[len(funcNameElems)-1])
	}

	remoteIP := r.RemoteAddr
	if pos := strings.Index(remoteIP, ":"); pos >= 0 {
		remoteIP = remoteIP[:pos]
	}
	return fmt.Sprintf("%s - %v - %d (%d) - %s%s - %s%s",
		details, err, httpStatusCode, errCode, uinfoStr, remoteIP, q, srcFileInfo)
}
