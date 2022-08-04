package common

import (
	"fmt"
	"net/http"
)

// HandlerFn is a framework-agnostic function to handle a vulnerable endpoint.
// `opaque` can be set to some framework-specific struct - for example, gin.Context.
//
// Prefer statuses 200 (success), 400 (generic, expected error), and 500 (generic, unexpected error).
//
// If a HandlerFn returns empty data, drivers should not write any data to the response.
type HandlerFn func(safety Safety, payload string, opaque interface{}) (data, mime string, status int)

// VulnerableFnWrapper is a function wrapping something vulnerable. Used
// to adapt things for use with GenericHandler. 'raw' indicates data
// should be sent verbatim, not decorated.
type VulnerableFnWrapper func(opaque interface{}, payload string) (data string, raw bool, err error)

// GenericHandler returns a generic replacement for HandlerFn. It requires VulnerableFnWrapper and Sanitize to be set.
func GenericHandler(s *Sink) (HandlerFn, error) {
	if s.Sanitize == nil {
		return nil, fmt.Errorf("sink %#v: internal error - Sanitizer cannot be nil", s)
	}
	if s.VulnerableFnWrapper == nil {
		return nil, fmt.Errorf("sink %#v: internal error - VulnerableFnWrapper cannot be nil", s)
	}
	return func(safety Safety, payload string, opaque interface{}) (data, mime string, status int) {
		mime = "text/plain"
		switch safety {
		case Unsafe:
			// nothing to do here
		case Safe:
			payload = s.Sanitize(payload)
		case NOOP:
			return "NOOP", mime, http.StatusOK
		default:
			msg := "expect one of 'unsafe', 'safe', 'noop' - instead got " + string(safety)
			return msg, mime, http.StatusBadRequest
		}
		res, raw, err := s.VulnerableFnWrapper(opaque, payload)
		if raw {
			if len(s.RawMime) > 0 {
				mime = s.RawMime
			}
			return res, mime, http.StatusOK
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
		return data, mime, status
	}, nil
}
