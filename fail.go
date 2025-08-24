package fail

import (
	"maps"
	"slices"
)

type Fail struct {
	msg     string
	userMsg string

	code           string
	exitCode       int
	httpStatusCode int

	causes     []error
	associated []error

	tags  map[string]struct{}
	attrs map[string]any
}

func newFail(msg string) Fail {
	if msg == "" {
		panic("fail message must not be an empty string")
	}

	return Fail{
		msg:            msg,
		code:           DefaultErrorCode,
		exitCode:       DefaultExitCode,
		httpStatusCode: DefaultHttpStatusCode,
		tags:           make(map[string]struct{}),
		attrs:          make(map[string]any),
	}
}

func newFailFrom(err error) Fail {
	if err == nil {
		panic("cannot create a Fail from a nil error")
	}

	if f, ok := err.(Fail); ok {
		return f
	}

	attrs := make(map[string]struct{})
	for _, t := range Tags(err) {
		attrs[t] = struct{}{}
	}

	return Fail{
		msg:            Message(err),
		userMsg:        UserMessage(err),
		code:           Code(err),
		exitCode:       ExitCode(err),
		httpStatusCode: HttpStatusCode(err),
		causes:         Causes(err),
		associated:     Associated(err),
		tags:           attrs,
		attrs:          Attributes(err),
	}
}

func (f Fail) Error() string {
	return f.msg
}

func (f Fail) ErrorCauses() []error {
	return f.causes
}

func (f Fail) ErrorAssociated() []error {
	return slices.Clone(f.associated)
}

func (f Fail) ErrorCode() string {
	return f.code
}

func (f Fail) ErrorExitCode() int {
	return f.exitCode
}

func (f Fail) ErrorHttpStatusCode() int {
	return f.httpStatusCode
}

func (f Fail) ErrorMessage() string {
	return f.msg
}

func (f Fail) ErrorUserMessage() string {
	return f.userMsg
}

func (f Fail) ErrorTags() []string {
	return slices.Collect(maps.Keys(f.tags))
}

func (f Fail) ErrorAttributes() map[string]any {
	return maps.Clone(f.attrs)
}
