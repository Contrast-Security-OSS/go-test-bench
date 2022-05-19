package common

import (
	"fmt"
	"net/http"
)

// HandlerFn is a framework-agnostic function to handle a vulnerable endpoint.
// `opaque` can be set to some framework-specific struct - for example, gin.Context.
//
// Prefer statuses 200 (success), 400 (generic, expected error), and 500 (generic, unexpected error).
type HandlerFn func(safety Safety, in string, opaque interface{}) (data string, status int)

// VulnerableFnWrapper is a function wrapping something vulnerable. Used
// to adapt things for use with GenericHandler. 'raw' indicates data
// should be sent verbatim, not decorated.
type VulnerableFnWrapper func(opaque interface{}, payload string) (data string, raw bool, err error)

// GenericHandler returns a generic replacement for HandlerFn. It requires VulnerableFnWrapper and Sanitize to be set.
func GenericHandler(s Sink) func(safety Safety, payload string, opaque interface{}) (data string, status int) {
	return func(safety Safety, payload string, opaque interface{}) (data string, status int) {
		if s.Sanitize == nil {
			return fmt.Sprintf("sink %#v: internal error - Sanitizer cannot be nil", s), http.StatusInternalServerError
		}
		if s.VulnerableFnWrapper == nil {
			return fmt.Sprintf("sink %#v: internal error - VulnerableFnWrapper cannot be nil", s), http.StatusInternalServerError
		}
		switch safety {
		case Unsafe:
			// nothing to do here
		case Safe:
			payload = s.Sanitize(payload)
		case NOOP:
			return "NOOP", http.StatusOK
		default:
			msg := "expect one of 'unsafe', 'safe', 'noop' - instead got " + string(safety)
			return msg, http.StatusBadRequest
		}
		res, raw, err := s.VulnerableFnWrapper(opaque, payload)
		if raw {
			return res, http.StatusOK
		}

		status = http.StatusOK
		e := "(no error)"
		if err != nil {
			e = err.Error()
			status = http.StatusBadRequest
		}
		if len(res) == 0 {
			res = "(no data returned)"
		}
		data = fmt.Sprintf("%q: %s action with payload=%q resulted in err=%s\nand data=\\\n%s", s.Name, safety, payload, e, res)
		return data, status
	}
}
