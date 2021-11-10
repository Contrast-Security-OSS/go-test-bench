package common

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetUserInput_GetParamValue(t *testing.T) {
	desiredVal := "get_value"
	request := getMockRequest(t, http.MethodGet, "/getQuery/testing?input="+desiredVal, "test")

	value := GetUserInput(request)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetUserInput_GetFormValue(t *testing.T) {
	desiredVal := "form_value"
	request := getMockRequest(t, http.MethodPost, "/form/testing", "input="+desiredVal)

	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	value := GetUserInput(request)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetUserInput_GetCookieValue(t *testing.T) {
	desiredVal := "cookie_value"
	request := getMockRequest(t, http.MethodGet, "/cookie/testing", "test")

	request.AddCookie(&http.Cookie{
		Name:  INPUT,
		Value: desiredVal,
		Path:  "/cookie/testing",
	})

	value := GetUserInput(request)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetUserInput_GetHeaderValue(t *testing.T) {
	desiredVal := "header_value"
	request := getMockRequest(t, http.MethodGet, "/header/testing", "test")

	request.Header.Add(INPUT, desiredVal)

	value := GetUserInput(request)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetParamValue(t *testing.T) {
	desiredVal := "get_value"
	request := getMockRequest(t, http.MethodGet, "/getQuery/testing?input="+desiredVal, "test")

	value := GetParamValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetPathValue(t *testing.T) {
	desiredVal := "<scipt>alert(1);</script>"
	request := getMockRequest(t, http.MethodGet, "/getQuery/parameter/"+desiredVal+"/testing", "test")

	value := GetPathValue(request, 3, 4)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetFormValue(t *testing.T) {
	desiredVal := "form_value"
	request := getMockRequest(t, http.MethodPost, "/form/testing", "input="+desiredVal)

	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	value := GetFormValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetCookieValue(t *testing.T) {
	desiredVal := "cookie_value"
	request := getMockRequest(t, http.MethodGet, "/cookie/testing", "test")

	request.AddCookie(&http.Cookie{
		Name:  INPUT,
		Value: desiredVal,
		Path:  "/cookie/testing",
	})

	value := GetCookieValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetHeaderValue(t *testing.T) {
	desiredVal := "header_value"
	request := getMockRequest(t, http.MethodGet, "/header/testing", "test")

	request.Header.Add(INPUT, desiredVal)

	value := GetHeaderValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func getMockRequest(t *testing.T, method, url, body string) *http.Request {
	r := strings.NewReader(body)
	request, err := http.NewRequest(method, url, r)
	if err != nil {
		t.Error(err.Error())
	}
	return request
}
