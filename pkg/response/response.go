package response

type Code int64

const (
	CodeSuccess     Code = 0
	CodeFailUnknown Code = 1 + iota
)

type Option func(opts *Response)

type Response struct {
	Code    int64  `json:"code"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func WithError(code int64, message string) Option {
	return func(opts *Response) {
		opts.Code = code
		opts.Message = message
	}
}

func WithData(data any) Option {
	return func(opts *Response) {
		opts.Data = data
	}
}

func defaultSuccessOptions() *Response {
	return &Response{
		Code: int64(CodeSuccess),
		Data: nil,
	}
}

func defaultFailOptions() *Response {
	return &Response{
		Code:    int64(CodeFailUnknown),
		Message: "unknown error",
	}
}

func Success() *Response {
	return defaultSuccessOptions()
}

func (o *Response) WithData(data any) *Response {
	o.Data = data
	return o
}

func Fail() *Response {
	return defaultFailOptions()
}

func (o *Response) WithError(code int64, message string) *Response {
	o.Code = code
	o.Message = message
	return o
}
