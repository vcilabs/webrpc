// Code generated by statik. DO NOT EDIT.

// Package contains static assets.
package embed

var Asset = "PK\x03\x04\x14\x00\x08\x00\x00\x00SN[Q\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00client.go.tmplUT\x05\x00\x01o\xed\x97_{{define \"client\"}}\x0d\n{{if or .Services .GoInterface}}\x0d\n//\x0d\n// Client\x0d\n//\x0d\n\x0d\n{{range .Services}}\x0d\nconst {{.Name | constPathPrefix}} = \"/rpc/{{.Name}}/\"\x0d\n{{end}}\x0d\n\x0d\n{{range .GoInterface}}\x0d\nconst {{.Name | constPathPrefix}} = \"/rpc/{{.Name}}/\"\x0d\n{{end}}\x0d\n \x0d\n{{range .Services}}\x0d\n  {{ $serviceName := .Name | clientServiceName}}\x0d\n  type {{$serviceName}} struct {\x0d\n    client HTTPClient\x0d\n    urls   [{{.Methods | countMethods}}]string\x0d\n  }\x0d\n\x0d\n  func {{.Name | newClientServiceName }}(addr string, client HTTPClient) {{.Name}} {\x0d\n    prefix := urlBase(addr) + {{.Name | constPathPrefix}}\x0d\n    urls := [{{.Methods | countMethods}}]string{\x0d\n      {{- range .Methods}}\x0d\n      prefix + \"{{.Name}}\",\x0d\n      {{- end}}\x0d\n    }\x0d\n    return &{{$serviceName}}{\x0d\n      client: client,\x0d\n      urls:   urls,\x0d\n    }\x0d\n  }\x0d\n\x0d\n  {{range $i, $method := .Methods}}\x0d\n    func (c *{{$serviceName}}) {{.Name}}({{.Inputs | methodInputs}}) ({{.Outputs | methodOutputs }}) {\x0d\n      {{- $inputVar := \"nil\" -}}\x0d\n      {{- $outputVar := \"nil\" -}}\x0d\n      {{- if .Inputs | len}}\x0d\n      {{- $inputVar = \"in\"}}\x0d\n      in := struct {\x0d\n        {{- range $i, $input := .Inputs}}\x0d\n          Arg{{$i}} {{$input | methodArgType}} `json:\"{{$input.Name | downcaseName}}\"`\x0d\n        {{- end}}          \x0d\n      }{ {{.Inputs | methodArgNames}} }\x0d\n      {{- end}}\x0d\n      {{- if .Outputs | len}}\x0d\n      {{- $outputVar = \"&out\"}}\x0d\n      out := struct {\x0d\n        {{- range $i, $output := .Outputs}}\x0d\n          Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name | downcaseName}}\"`\x0d\n        {{- end}}          \x0d\n      }{}\x0d\n    {{- end}}\x0d\n\x0d\n      err := doJSONRequest(ctx, c.client, c.urls[{{$i}}], {{$inputVar}}, {{$outputVar}})\x0d\n      return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} err\x0d\n    }\x0d\n  {{end}}\x0d\n{{end}}\x0d\n\x0d\n{{range .GoInterface}}\x0d\n  {{ $InterfaceName := .Name | clientServiceName}}\x0d\n  type {{$InterfaceName}} struct {\x0d\n    client HTTPClient\x0d\n    urls   [{{.Methods | countMethods}}]string\x0d\n  }\x0d\n\x0d\n  func {{.Name | newClientServiceName }}(addr string, client HTTPClient) {{.Name}} {\x0d\n    prefix := urlBase(addr) + {{.Name | constPathPrefix}}\x0d\n    urls := [{{.Methods | countMethods}}]string{\x0d\n      {{- range .Methods}}\x0d\n      prefix + \"{{.Name}}\",\x0d\n      {{- end}}\x0d\n    }\x0d\n    return &{{$InterfaceName}}{\x0d\n      client: client,\x0d\n      urls:   urls,\x0d\n    }\x0d\n  }\x0d\n\x0d\n  {{range $i, $method := .Methods}}\x0d\n    func (c *{{$InterfaceName}}) {{.Name}}({{.Inputs | methodInputs}}) ({{.Outputs | methodOutputs }}) {\x0d\n      {{- $inputVar := \"nil\" -}}\x0d\n      {{- $outputVar := \"nil\" -}}\x0d\n      {{- if .Inputs | len}}\x0d\n      {{- $inputVar = \"in\"}}\x0d\n      in := struct {\x0d\n        {{- range $i, $input := .Inputs}}\x0d\n          Arg{{$i}} {{$input | methodArgType}} `json:\"{{$input.Name | downcaseName}}\"`\x0d\n        {{- end}}          \x0d\n      }{ {{.Inputs | methodArgNames}} }\x0d\n      {{- end}}\x0d\n      {{- if .Outputs | len}}\x0d\n      {{- $outputVar = \"&out\"}}\x0d\n      out := struct {\x0d\n        {{- range $i, $output := .Outputs}}\x0d\n          Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name | downcaseName}}\"`\x0d\n        {{- end}}          \x0d\n      }{}\x0d\n    {{- end}}\x0d\n\x0d\n      err := doJSONRequest(ctx, c.client, c.urls[{{$i}}], {{$inputVar}}, {{$outputVar}})\x0d\n      return {{argsList .Outputs \"out.Ret\"}}{{commaIfLen .Outputs}} err\x0d\n    }\x0d\n  {{end}}\x0d\n{{end}}\x0d\n\x0d\n\x0d\n// HTTPClient is the interface used by generated clients to send HTTP requests.\x0d\n// It is fulfilled by *(net/http).Client, which is sufficient for most users.\x0d\n// Users can provide their own implementation for special retry policies.\x0d\ntype HTTPClient interface {\x0d\n  Do(req *http.Request) (*http.Response, error)\x0d\n}\x0d\n\x0d\n// urlBase helps ensure that addr specifies a scheme. If it is unparsable\x0d\n// as a URL, it returns addr unchanged.\x0d\nfunc urlBase(addr string) string {\x0d\n  // If the addr specifies a scheme, use it. If not, default to\x0d\n  // http. If url.Parse fails on it, return it unchanged.\x0d\n  url, err := url.Parse(addr)\x0d\n  if err != nil {\x0d\n    return addr\x0d\n  }\x0d\n  if url.Scheme == \"\" {\x0d\n    url.Scheme = \"http\"\x0d\n  }\x0d\n  return url.String()\x0d\n}\x0d\n\x0d\n// newRequest makes an http.Request from a client, adding common headers.\x0d\nfunc newRequest(ctx context.Context, url string, reqBody io.Reader, contentType string) (*http.Request, error) {\x0d\n  req, err := http.NewRequest(\"POST\", url, reqBody)\x0d\n  if err != nil {\x0d\n    return nil, err\x0d\n  }\x0d\n  req.Header.Set(\"Accept\", contentType)\x0d\n  req.Header.Set(\"Content-Type\", contentType)\x0d\n	if headers, ok := HTTPRequestHeaders(ctx); ok {\x0d\n		for k := range headers {\x0d\n			for _, v := range headers[k] {\x0d\n				req.Header.Add(k, v)\x0d\n			}\x0d\n		}\x0d\n	}\x0d\n  return req, nil\x0d\n}\x0d\n\x0d\n// doJSONRequest is common code to make a request to the remote service.\x0d\nfunc doJSONRequest(ctx context.Context, client HTTPClient, url string, in, out interface{}) error {\x0d\n	reqBody, err := json.Marshal(in)\x0d\n	if err != nil {\x0d\n		return clientError(\"failed to marshal json request\", err)\x0d\n	}\x0d\n	if err = ctx.Err(); err != nil {\x0d\n		return clientError(\"aborted because context was done\", err)\x0d\n	}\x0d\n\x0d\n	req, err := newRequest(ctx, url, bytes.NewBuffer(reqBody), \"application/json\")\x0d\n	if err != nil {\x0d\n		return clientError(\"could not build request\", err)\x0d\n	}\x0d\n	resp, err := client.Do(req)\x0d\n	if err != nil {\x0d\n		return clientError(\"request failed\", err)\x0d\n	}\x0d\n\x0d\n	defer func() {\x0d\n		cerr := resp.Body.Close()\x0d\n		if err == nil && cerr != nil {\x0d\n			err = clientError(\"failed to close response body\", cerr)\x0d\n		}\x0d\n	}()\x0d\n\x0d\n	if err = ctx.Err(); err != nil {\x0d\n		return clientError(\"aborted because context was done\", err)\x0d\n	}\x0d\n\x0d\n	if resp.StatusCode != 200 {\x0d\n		return errorFromResponse(resp)\x0d\n	}\x0d\n\x0d\n	if out != nil {\x0d\n		respBody, err := ioutil.ReadAll(resp.Body)\x0d\n		if err != nil {\x0d\n			return clientError(\"failed to read response body\", err)\x0d\n		}\x0d\n\x0d\n		err = json.Unmarshal(respBody, &out)\x0d\n		if err != nil {\x0d\n			return clientError(\"failed to unmarshal json response body\", err)\x0d\n		}\x0d\n		if err = ctx.Err(); err != nil {\x0d\n			return clientError(\"aborted because context was done\", err)\x0d\n		}\x0d\n	}\x0d\n\x0d\n	return nil\x0d\n}\x0d\n\x0d\n// errorFromResponse builds a webrpc Error from a non-200 HTTP response.\x0d\nfunc errorFromResponse(resp *http.Response) Error {\x0d\n	respBody, err := ioutil.ReadAll(resp.Body)\x0d\n	if err != nil {\x0d\n		return clientError(\"failed to read server error response body\", err)\x0d\n	}\x0d\n\x0d\n	var respErr ErrorPayload\x0d\n	if err := json.Unmarshal(respBody, &respErr); err != nil {\x0d\n		return clientError(\"failed unmarshal error response\", err)\x0d\n	}\x0d\n\x0d\n	errCode := ErrorCode(respErr.Code)\x0d\n\x0d\n	if HTTPStatusFromErrorCode(errCode) == 0 {\x0d\n		return ErrorInternal(\"invalid code returned from server error response: %s\", respErr.Code)\x0d\n	}\x0d\n\x0d\n	return &rpcErr{\x0d\n		code:  errCode,\x0d\n		msg:   respErr.Msg,\x0d\n		cause: errors.New(respErr.Cause),\x0d\n	}\x0d\n}\x0d\n\x0d\nfunc clientError(desc string, err error) Error {\x0d\n	return WrapError(ErrInternal, err, desc)\x0d\n}\x0d\n\x0d\nfunc WithHTTPRequestHeaders(ctx context.Context, h http.Header) (context.Context, error) {\x0d\n	if _, ok := h[\"Accept\"]; ok {\x0d\n		return nil, errors.New(\"provided header cannot set Accept\")\x0d\n	}\x0d\n	if _, ok := h[\"Content-Type\"]; ok {\x0d\n		return nil, errors.New(\"provided header cannot set Content-Type\")\x0d\n	}\x0d\n\x0d\n	copied := make(http.Header, len(h))\x0d\n	for k, vv := range h {\x0d\n		if vv == nil {\x0d\n			copied[k] = nil\x0d\n			continue\x0d\n		}\x0d\n		copied[k] = make([]string, len(vv))\x0d\n		copy(copied[k], vv)\x0d\n	}\x0d\n\x0d\n	return context.WithValue(ctx, HTTPClientRequestHeadersCtxKey, copied), nil\x0d\n}\x0d\n\x0d\nfunc HTTPRequestHeaders(ctx context.Context) (http.Header, bool) {\x0d\n	h, ok := ctx.Value(HTTPClientRequestHeadersCtxKey).(http.Header)\x0d\n	return h, ok\x0d\n}\x0d\n{{end}}\x0d\n{{end}}\x0d\nPK\x07\x08\xe7\x8c:\xa3\x99\x1d\x00\x00\x99\x1d\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0f\x00	\x00helpers.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"helpers\"}}\x0d\n\x0d\n//\x0d\n// Helpers\x0d\n//\x0d\n\x0d\ntype ErrorPayload struct {\x0d\n	Status int    `json:\"status\"`\x0d\n	Code   string `json:\"code\"`\x0d\n	Cause  string `json:\"cause,omitempty\"`\x0d\n	Msg    string `json:\"msg\"`\x0d\n	Error  string `json:\"error\"`\x0d\n}\x0d\n\x0d\ntype Error interface {\x0d\n	// Code is of the valid error codes\x0d\n	Code() ErrorCode\x0d\n\x0d\n	// Msg returns a human-readable, unstructured messages describing the error\x0d\n	Msg() string\x0d\n\x0d\n	// Cause is reason for the error\x0d\n	Cause() error\x0d\n\x0d\n	// Error returns a string of the form \"webrpc error <Code>: <Msg>\"\x0d\n	Error() string\x0d\n\x0d\n	// Error response payload\x0d\n	Payload() ErrorPayload\x0d\n}\x0d\n\x0d\nfunc Errorf(code ErrorCode, msgf string, args ...interface{}) Error {\x0d\n	msg := fmt.Sprintf(msgf, args...)\x0d\n	if IsValidErrorCode(code) {\x0d\n		return &rpcErr{code: code, msg: msg}\x0d\n	}\x0d\n	return &rpcErr{code: ErrInternal, msg: \"invalid error type \" + string(code)}\x0d\n}\x0d\n\x0d\nfunc WrapError(code ErrorCode, cause error, format string, args ...interface{}) Error {\x0d\n	msg := fmt.Sprintf(format, args...)\x0d\n	if IsValidErrorCode(code) {\x0d\n		return &rpcErr{code: code, msg: msg, cause: cause}\x0d\n	}\x0d\n	return &rpcErr{code: ErrInternal, msg: \"invalid error type \" + string(code), cause: cause}\x0d\n}\x0d\n\x0d\nfunc Failf(format string, args ...interface{}) Error {\x0d\n	return Errorf(ErrFail, format, args...)\x0d\n}\x0d\n\x0d\nfunc WrapFailf(cause error, format string, args ...interface{}) Error {\x0d\n	return WrapError(ErrFail, cause, format, args...)\x0d\n}\x0d\n\x0d\nfunc ErrorNotFound(format string, args ...interface{}) Error {\x0d\n	return Errorf(ErrNotFound, format, args...)\x0d\n}\x0d\n\x0d\nfunc ErrorInvalidArgument(argument string, validationMsg string) Error {\x0d\n	return Errorf(ErrInvalidArgument, argument+\" \"+validationMsg)\x0d\n}\x0d\n\x0d\nfunc ErrorRequiredArgument(argument string) Error {\x0d\n	return ErrorInvalidArgument(argument, \"is required\")\x0d\n}\x0d\n\x0d\nfunc ErrorInternal(format string, args ...interface{}) Error {\x0d\n	return Errorf(ErrInternal, format, args...)\x0d\n}\x0d\n\x0d\ntype ErrorCode string\x0d\n\x0d\nconst (\x0d\n	// Unknown error. For example when handling errors raised by APIs that do not\x0d\n	// return enough error information.\x0d\n	ErrUnknown ErrorCode = \"unknown\"\x0d\n\x0d\n	// Fail error. General failure error type.\x0d\n	ErrFail ErrorCode = \"fail\"\x0d\n\x0d\n	// Canceled indicates the operation was cancelled (typically by the caller).\x0d\n	ErrCanceled ErrorCode = \"canceled\"\x0d\n\x0d\n	// InvalidArgument indicates client specified an invalid argument. It\x0d\n	// indicates arguments that are problematic regardless of the state of the\x0d\n	// system (i.e. a malformed file name, required argument, number out of range,\x0d\n	// etc.).\x0d\n	ErrInvalidArgument ErrorCode = \"invalid argument\"\x0d\n\x0d\n	// DeadlineExceeded means operation expired before completion. For operations\x0d\n	// that change the state of the system, this error may be returned even if the\x0d\n	// operation has completed successfully (timeout).\x0d\n	ErrDeadlineExceeded ErrorCode = \"deadline exceeded\"\x0d\n\x0d\n	// NotFound means some requested entity was not found.\x0d\n	ErrNotFound ErrorCode = \"not found\"\x0d\n\x0d\n	// BadRoute means that the requested URL path wasn't routable to a webrpc\x0d\n	// service and method. This is returned by the generated server, and usually\x0d\n	// shouldn't be returned by applications. Instead, applications should use\x0d\n	// NotFound or Unimplemented.\x0d\n	ErrBadRoute ErrorCode = \"bad route\"\x0d\n\x0d\n	// AlreadyExists means an attempt to create an entity failed because one\x0d\n	// already exists.\x0d\n	ErrAlreadyExists ErrorCode = \"already exists\"\x0d\n\x0d\n	// PermissionDenied indicates the caller does not have permission to execute\x0d\n	// the specified operation. It must not be used if the caller cannot be\x0d\n	// identified (Unauthenticated).\x0d\n	ErrPermissionDenied ErrorCode = \"permission denied\"\x0d\n\x0d\n	// Unauthenticated indicates the request does not have valid authentication\x0d\n	// credentials for the operation.\x0d\n	ErrUnauthenticated ErrorCode = \"unauthenticated\"\x0d\n\x0d\n	// ResourceExhausted indicates some resource has been exhausted, perhaps a\x0d\n	// per-user quota, or perhaps the entire file system is out of space.\x0d\n	ErrResourceExhausted ErrorCode = \"resource exhausted\"\x0d\n\x0d\n	// FailedPrecondition indicates operation was rejected because the system is\x0d\n	// not in a state required for the operation's execution. For example, doing\x0d\n	// an rmdir operation on a directory that is non-empty, or on a non-directory\x0d\n	// object, or when having conflicting read-modify-write on the same resource.\x0d\n	ErrFailedPrecondition ErrorCode = \"failed precondition\"\x0d\n\x0d\n	// Aborted indicates the operation was aborted, typically due to a concurrency\x0d\n	// issue like sequencer check failures, transaction aborts, etc.\x0d\n	ErrAborted ErrorCode = \"aborted\"\x0d\n\x0d\n	// OutOfRange means operation was attempted past the valid range. For example,\x0d\n	// seeking or reading past end of a paginated collection.\x0d\n	//\x0d\n	// Unlike InvalidArgument, this error indicates a problem that may be fixed if\x0d\n	// the system state changes (i.e. adding more items to the collection).\x0d\n	//\x0d\n	// There is a fair bit of overlap between FailedPrecondition and OutOfRange.\x0d\n	// We recommend using OutOfRange (the more specific error) when it applies so\x0d\n	// that callers who are iterating through a space can easily look for an\x0d\n	// OutOfRange error to detect when they are done.\x0d\n	ErrOutOfRange ErrorCode = \"out of range\"\x0d\n\x0d\n	// Unimplemented indicates operation is not implemented or not\x0d\n	// supported/enabled in this service.\x0d\n	ErrUnimplemented ErrorCode = \"unimplemented\"\x0d\n\x0d\n	// Internal errors. When some invariants expected by the underlying system\x0d\n	// have been broken. In other words, something bad happened in the library or\x0d\n	// backend service. Do not confuse with HTTP Internal Server Error; an\x0d\n	// Internal error could also happen on the client code, i.e. when parsing a\x0d\n	// server response.\x0d\n	ErrInternal ErrorCode = \"internal\"\x0d\n\x0d\n	// Unavailable indicates the service is currently unavailable. This is a most\x0d\n	// likely a transient condition and may be corrected by retrying with a\x0d\n	// backoff.\x0d\n	ErrUnavailable ErrorCode = \"unavailable\"\x0d\n\x0d\n	// DataLoss indicates unrecoverable data loss or corruption.\x0d\n	ErrDataLoss ErrorCode = \"data loss\"\x0d\n\x0d\n	// ErrNone is the zero-value, is considered an empty error and should not be\x0d\n	// used.\x0d\n	ErrNone ErrorCode = \"\"\x0d\n)\x0d\n\x0d\nfunc HTTPStatusFromErrorCode(code ErrorCode) int {\x0d\n	switch code {\x0d\n	case ErrCanceled:\x0d\n		return 408 // RequestTimeout\x0d\n	case ErrUnknown:\x0d\n		return 400 // Bad Request\x0d\n	case ErrFail:\x0d\n		return 422 // Unprocessable Entity\x0d\n	case ErrInvalidArgument:\x0d\n		return 400 // BadRequest\x0d\n	case ErrDeadlineExceeded:\x0d\n		return 408 // RequestTimeout\x0d\n	case ErrNotFound:\x0d\n		return 404 // Not Found\x0d\n	case ErrBadRoute:\x0d\n		return 404 // Not Found\x0d\n	case ErrAlreadyExists:\x0d\n		return 409 // Conflict\x0d\n	case ErrPermissionDenied:\x0d\n		return 403 // Forbidden\x0d\n	case ErrUnauthenticated:\x0d\n		return 401 // Unauthorized\x0d\n	case ErrResourceExhausted:\x0d\n		return 403 // Forbidden\x0d\n	case ErrFailedPrecondition:\x0d\n		return 412 // Precondition Failed\x0d\n	case ErrAborted:\x0d\n		return 409 // Conflict\x0d\n	case ErrOutOfRange:\x0d\n		return 400 // Bad Request\x0d\n	case ErrUnimplemented:\x0d\n		return 501 // Not Implemented\x0d\n	case ErrInternal:\x0d\n		return 500 // Internal Server Error\x0d\n	case ErrUnavailable:\x0d\n		return 503 // Service Unavailable\x0d\n	case ErrDataLoss:\x0d\n		return 500 // Internal Server Error\x0d\n	case ErrNone:\x0d\n		return 200 // OK\x0d\n	default:\x0d\n		return 0 // Invalid!\x0d\n	}\x0d\n}\x0d\n\x0d\nfunc IsErrorCode(err error, code ErrorCode) bool {\x0d\n	if rpcErr, ok := err.(Error); ok {\x0d\n		if rpcErr.Code() == code {\x0d\n			return true\x0d\n		}\x0d\n	}\x0d\n	return false\x0d\n}\x0d\n\x0d\nfunc IsValidErrorCode(code ErrorCode) bool {\x0d\n	return HTTPStatusFromErrorCode(code) != 0\x0d\n}\x0d\n\x0d\ntype rpcErr struct {\x0d\n	code  ErrorCode\x0d\n	msg   string\x0d\n	cause error\x0d\n}\x0d\n\x0d\nfunc (e *rpcErr) Code() ErrorCode {\x0d\n	return e.code\x0d\n}\x0d\n\x0d\nfunc (e *rpcErr) Msg() string {\x0d\n	return e.msg\x0d\n}\x0d\n\x0d\nfunc (e *rpcErr) Cause() error {\x0d\n	return e.cause\x0d\n}\x0d\n\x0d\nfunc (e *rpcErr) Error() string {\x0d\n	if e.cause != nil && e.cause.Error() != \"\" {\x0d\n		if e.msg != \"\" {\x0d\n			return fmt.Sprintf(\"webrpc %s error: %s -- %s\", e.code, e.cause.Error(), e.msg)\x0d\n		} else {\x0d\n			return fmt.Sprintf(\"webrpc %s error: %s\", e.code, e.cause.Error())\x0d\n		}\x0d\n	} else {\x0d\n		return fmt.Sprintf(\"webrpc %s error: %s\", e.code, e.msg)\x0d\n	}\x0d\n}\x0d\n\x0d\nfunc (e *rpcErr) Payload() ErrorPayload {\x0d\n	statusCode := HTTPStatusFromErrorCode(e.Code())\x0d\n	errPayload := ErrorPayload{\x0d\n		Status: statusCode,\x0d\n		Code:   string(e.Code()),\x0d\n		Msg:    e.Msg(),\x0d\n		Error:  e.Error(),\x0d\n	}\x0d\n	if e.Cause() != nil {\x0d\n		errPayload.Cause = e.Cause().Error()\x0d\n	}\x0d\n	return errPayload\x0d\n}\x0d\n\x0d\ntype contextKey struct {\x0d\n	name string\x0d\n}\x0d\n\x0d\nfunc (k *contextKey) String() string {\x0d\n	return \"webrpc context value \" + k.name\x0d\n}\x0d\n\x0d\nvar (\x0d\n	// For Client\x0d\n	HTTPClientRequestHeadersCtxKey = &contextKey{\"HTTPClientRequestHeaders\"}\x0d\n\x0d\n	// For Server\x0d\n	HTTPResponseWriterCtxKey = &contextKey{\"HTTPResponseWriter\"}\x0d\n\x0d\n	HTTPRequestCtxKey = &contextKey{\"HTTPRequest\"}\x0d\n\x0d\n	ServiceNameCtxKey = &contextKey{\"ServiceName\"}\x0d\n\x0d\n	MethodNameCtxKey = &contextKey{\"MethodName\"}\x0d\n)\x0d\n\x0d\n{{end}}\x0d\nPK\x07\x08}M\xd0\x98\x02#\x00\x00\x02#\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x11\x00	\x00proto.gen.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{- define \"proto\" -}}\x0d\n// {{.Name}} {{.SchemaVersion}} {{.SchemaHash}}\x0d\n// --\x0d\n// This file has been generated by https://github.com/webrpc/webrpc using gen/golang\x0d\n// Do not edit by hand. Update your webrpc schema and re-generate.\x0d\npackage {{.TargetOpts.PkgName}}\x0d\n\x0d\nimport (\x0d\n  \"context\"\x0d\n  \"encoding/json\"\x0d\n  \"fmt\"\x0d\n  \"io/ioutil\"\x0d\n  \"net/http\"\x0d\n  \"time\"\x0d\n  \"strings\"\x0d\n  \"bytes\"\x0d\n  \"errors\"\x0d\n  \"io\"\x0d\n  \"net/url\"\x0d\n)\x0d\n\x0d\n// WebRPC description and code-gen version\x0d\nfunc WebRPCVersion() string {\x0d\n  return \"{{.WebRPCVersion}}\"\x0d\n}\x0d\n\x0d\n// Schema version of your RIDL schema\x0d\nfunc WebRPCSchemaVersion() string {\x0d\n  return \"{{.SchemaVersion}}\"\x0d\n}\x0d\n\x0d\n// Schema hash generated from your RIDL schema\x0d\nfunc WebRPCSchemaHash() string {\x0d\n  return \"{{.SchemaHash}}\"\x0d\n}\x0d\n\x0d\n{{template \"types\" .}}\x0d\n\x0d\n{{if .TargetOpts.Server}}\x0d\n  {{template \"server\" .}}\x0d\n{{end}}\x0d\n\x0d\n{{if .TargetOpts.Client}}\x0d\n  {{template \"client\" .}}\x0d\n{{end}}\x0d\n\x0d\n{{template \"helpers\" .}}\x0d\n\x0d\n{{- end}}\x0d\nPK\x07\x08\xbe\xd0P\xbd\xba\x03\x00\x00\xba\x03\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0e\x00	\x00server.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"server\"}}\x0d\n{{if .Services}}\x0d\n//\x0d\n// Server\x0d\n//\x0d\n\x0d\ntype WebRPCServer interface {\x0d\n  http.Handler\x0d\n}\x0d\n\x0d\n{{- range .Services}}\x0d\n  {{$name := .Name}}\x0d\n  {{$serviceName := .Name | serverServiceName}}\x0d\n\x0d\n  type {{$serviceName}} struct {\x0d\n    {{.Name}}\x0d\n  }\x0d\n\x0d\n  func {{ .Name | newServerServiceName }}(svc {{.Name}}) WebRPCServer {\x0d\n    return &{{$serviceName}}{\x0d\n      {{.Name}}: svc,\x0d\n    }\x0d\n  }\x0d\n\x0d\n  func (s *{{$serviceName}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {\x0d\n    ctx := r.Context()\x0d\n    ctx = context.WithValue(ctx, HTTPResponseWriterCtxKey, w)\x0d\n    ctx = context.WithValue(ctx, HTTPRequestCtxKey, r)\x0d\n    ctx = context.WithValue(ctx, ServiceNameCtxKey, \"{{.Name}}\")\x0d\n\x0d\n    if r.Method != \"POST\" {\x0d\n      err := Errorf(ErrBadRoute, \"unsupported method %q (only POST is allowed)\", r.Method)\x0d\n      RespondWithError(w, err)\x0d\n      return\x0d\n    }\x0d\n\x0d\n    switch r.URL.Path {\x0d\n    {{- range .Methods}}\x0d\n    case \"/rpc/{{$name}}/{{.Name}}\":\x0d\n      s.{{.Name | serviceMethodName}}(ctx, w, r)\x0d\n      return\x0d\n    {{- end}}\x0d\n    default:\x0d\n      err := Errorf(ErrBadRoute, \"no handler for path %q\", r.URL.Path)\x0d\n      RespondWithError(w, err)\x0d\n      return\x0d\n    }\x0d\n  }\x0d\n\x0d\n  {{range .Methods}}\x0d\n    func (s *{{$serviceName}}) {{.Name | serviceMethodName}}(ctx context.Context, w http.ResponseWriter, r *http.Request) {\x0d\n      header := r.Header.Get(\"Content-Type\")\x0d\n      i := strings.Index(header, \";\")\x0d\n      if i == -1 {\x0d\n        i = len(header)\x0d\n      }\x0d\n\x0d\n      switch strings.TrimSpace(strings.ToLower(header[:i])) {\x0d\n      case \"application/json\":\x0d\n        s.{{ .Name | serviceMethodJSONName }}(ctx, w, r)\x0d\n      default:\x0d\n        err := Errorf(ErrBadRoute, \"unexpected Content-Type: %q\", r.Header.Get(\"Content-Type\"))\x0d\n        RespondWithError(w, err)\x0d\n      }\x0d\n    }\x0d\n\x0d\n    func (s *{{$serviceName}}) {{.Name | serviceMethodJSONName}}(ctx context.Context, w http.ResponseWriter, r *http.Request) {\x0d\n      var err error\x0d\n      ctx = context.WithValue(ctx, MethodNameCtxKey, \"{{.Name}}\")\x0d\n\x0d\n      {{- if .Inputs|len}}\x0d\n      reqContent := struct {\x0d\n      {{- range $i, $input := .Inputs}}\x0d\n        Arg{{$i}} {{. | methodArgType}} `json:\"{{$input.Name | downcaseName}}\"`\x0d\n      {{- end}}\x0d\n      }{}\x0d\n\x0d\n      reqBody, err := ioutil.ReadAll(r.Body)\x0d\n      if err != nil {\x0d\n        err = WrapError(ErrInternal, err, \"failed to read request data\")\x0d\n        RespondWithError(w, err)\x0d\n        return\x0d\n      }\x0d\n      defer r.Body.Close()\x0d\n\x0d\n      err = json.Unmarshal(reqBody, &reqContent)\x0d\n      if err != nil {\x0d\n        err = WrapError(ErrInvalidArgument, err, \"failed to unmarshal request data\")\x0d\n        RespondWithError(w, err)\x0d\n        return\x0d\n      }\x0d\n      {{- end}}\x0d\n\x0d\n      // Call service method\x0d\n      {{- range $i, $output := .Outputs}}\x0d\n      var ret{{$i}} {{$output | methodArgType}}\x0d\n      {{- end}}\x0d\n      func() {\x0d\n        defer func() {\x0d\n          // In case of a panic, serve a 500 error and then panic.\x0d\n          if rr := recover(); rr != nil {\x0d\n            RespondWithError(w, ErrorInternal(\"internal service panic\"))\x0d\n            panic(rr)\x0d\n          }\x0d\n        }()\x0d\n        {{argsList .Outputs \"ret\"}}{{.Outputs | commaIfLen}} err = s.{{$name}}.{{.Name}}(ctx{{.Inputs | commaIfLen}}{{argsList .Inputs \"reqContent.Arg\"}})\x0d\n      }()\x0d\n      {{- if .Outputs | len}}\x0d\n      respContent := struct {\x0d\n      {{- range $i, $output := .Outputs}}\x0d\n        Ret{{$i}} {{$output | methodArgType}} `json:\"{{$output.Name | downcaseName}}\"`\x0d\n      {{- end}}         \x0d\n      }{ {{argsList .Outputs \"ret\"}} }\x0d\n      {{- end}}\x0d\n\x0d\n      if err != nil {\x0d\n        RespondWithError(w, err)\x0d\n        return\x0d\n      }\x0d\n\x0d\n      {{- if .Outputs | len}}\x0d\n      respBody, err := json.Marshal(respContent)\x0d\n      if err != nil {\x0d\n        err = WrapError(ErrInternal, err, \"failed to marshal json response\")\x0d\n        RespondWithError(w, err)\x0d\n        return\x0d\n      }\x0d\n      {{- end}}\x0d\n\x0d\n      w.Header().Set(\"Content-Type\", \"application/json\")\x0d\n      w.WriteHeader(http.StatusOK)\x0d\n\x0d\n      {{- if .Outputs | len}}\x0d\n      w.Write(respBody)\x0d\n      {{- end}}\x0d\n    }\x0d\n  {{end}}\x0d\n{{- end}}\x0d\n\x0d\nfunc RespondWithError(w http.ResponseWriter, err error) {\x0d\n	rpcErr, ok := err.(Error)\x0d\n	if !ok {\x0d\n		rpcErr = WrapError(ErrInternal, err, \"webrpc error\")\x0d\n	}\x0d\n\x0d\n	statusCode := HTTPStatusFromErrorCode(rpcErr.Code())\x0d\n\x0d\n	w.Header().Set(\"Content-Type\", \"application/json\")\x0d\n	w.WriteHeader(statusCode)\x0d\n\x0d\n	respBody, _ := json.Marshal(rpcErr.Payload())\x0d\n	w.Write(respBody)\x0d\n}\x0d\n{{end}}\x0d\n{{end}}\x0d\nPK\x07\x08\xaf\x91\xd0\xd2\x83\x11\x00\x00\x83\x11\x00\x00PK\x03\x04\x14\x00\x08\x00\x00\x00WQTQ\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x00\x0d\x00	\x00types.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_{{define \"types\"}}\x0d\n\x0d\n{{if .Messages}}\x0d\n//\x0d\n// Types\x0d\n//\x0d\n\x0d\n{{range .Messages}}\x0d\n  {{if .Type | isEnum}}\x0d\n    {{$enumName := .Name}}\x0d\n    {{$enumType := .EnumType}}\x0d\n    type {{$enumName}} {{$enumType}}\x0d\n\x0d\n    const (\x0d\n      {{- range .Fields}}\x0d\n        {{$enumName}}_{{.Name}} {{$enumName}} = {{.Value}}\x0d\n      {{- end}}\x0d\n    )\x0d\n\x0d\n    var {{$enumName}}_name = map[{{$enumType}}]string {\x0d\n      {{- range .Fields}}\x0d\n        {{.Value}}: \"{{.Name}}\",\x0d\n      {{- end}}\x0d\n    }\x0d\n\x0d\n    var {{$enumName}}_value = map[string]{{$enumType}} {\x0d\n      {{- range .Fields}}\x0d\n        \"{{.Name}}\": {{.Value}},\x0d\n      {{- end}}\x0d\n    }\x0d\n\x0d\n    func (x {{$enumName}}) String() string {\x0d\n      return {{$enumName}}_name[{{$enumType}}(x)]\x0d\n    }\x0d\n\x0d\n    func (x {{$enumName}}) MarshalJSON() ([]byte, error) {\x0d\n      buf := bytes.NewBufferString(`\"`)\x0d\n      buf.WriteString({{$enumName}}_name[{{$enumType}}(x)])\x0d\n      buf.WriteString(`\"`)\x0d\n      return buf.Bytes(), nil\x0d\n    }\x0d\n\x0d\n    func (x *{{$enumName}}) UnmarshalJSON(b []byte) error {\x0d\n      var j string\x0d\n      err := json.Unmarshal(b, &j)\x0d\n      if err != nil {\x0d\n        return err\x0d\n      }\x0d\n      *x = {{$enumName}}({{$enumName}}_value[j])\x0d\n      return nil\x0d\n    }\x0d\n  {{end}}\x0d\n  {{if .Type | isStruct  }}\x0d\n    type {{.Name}} struct {\x0d\n      {{- range .Fields}}\x0d\n        {{. | exportedField}} {{. | fieldOptional}}{{. | fieldTypeDef}} {{. | fieldTags}}\x0d\n      {{- end}}\x0d\n    }\x0d\n  {{end}}\x0d\n{{end}}\x0d\n{{end}}\x0d\n{{if .Services}}\x0d\n  {{range .Services}}\x0d\n    type {{.Name}} interface {\x0d\n      {{- range .Methods}}\x0d\n        {{.Name}}({{.Inputs | methodInputs}}) ({{.Outputs | methodOutputs}})\x0d\n      {{- end}}\x0d\n    }\x0d\n  {{end}}\x0d\n  var WebRPCServices = map[string][]string{\x0d\n    {{- range .Services}}\x0d\n      \"{{.Name}}\": {\x0d\n        {{- range .Methods}}\x0d\n          \"{{.Name}}\",\x0d\n        {{- end}}\x0d\n      },\x0d\n    {{- end}}\x0d\n  }\x0d\n{{end}}\x0d\n\x0d\n{{end}}\x0d\nPK\x07\x08\x0f\x80\xe8\xc2P\x07\x00\x00P\x07\x00\x00PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00SN[Q\xe7\x8c:\xa3\x99\x1d\x00\x00\x99\x1d\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\x00\x00\x00\x00client.go.tmplUT\x05\x00\x01o\xed\x97_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ}M\xd0\x98\x02#\x00\x00\x02#\x00\x00\x0f\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\xde\x1d\x00\x00helpers.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ\xbe\xd0P\xbd\xba\x03\x00\x00\xba\x03\x00\x00\x11\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81&A\x00\x00proto.gen.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ\xaf\x91\xd0\xd2\x83\x11\x00\x00\x83\x11\x00\x00\x0e\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81(E\x00\x00server.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x01\x02\x14\x03\x14\x00\x08\x00\x00\x00WQTQ\x0f\x80\xe8\xc2P\x07\x00\x00P\x07\x00\x00\x0d\x00	\x00\x00\x00\x00\x00\x00\x00\x00\x00\xb6\x81\xf0V\x00\x00types.go.tmplUT\x05\x00\x01\xa7\xb7\x8e_PK\x05\x06\x00\x00\x00\x00\x05\x00\x05\x00\\\x01\x00\x00\x84^\x00\x00\x00\x00"
