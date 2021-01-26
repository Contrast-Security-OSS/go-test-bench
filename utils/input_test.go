package utils

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetParamValue(t *testing.T) {
	desiredVal := "get_value"
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/getQuery/testing?input="+desiredVal, r)
	if err != nil {
		t.Error(err.Error())
	}
	value := GetParamValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetPathValue(t *testing.T) {
	desiredVal := "<scipt>alert(1);</script>"
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/getQuery/parameter/"+desiredVal+"/testing", r)
	if err != nil {
		t.Error(err.Error())
	}

	value := GetPathValue(request, 3, 4)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetFormValue(t *testing.T) {
	desiredVal := "form_value"
	r := strings.NewReader("input=" + desiredVal)
	request, err := http.NewRequest(http.MethodPost, "/form/testing", r)
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	value := GetFormValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetCookieValue(t *testing.T) {
	desiredVal := "cookie_value"
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/cookie/testing", r)
	if err != nil {
		t.Error(err.Error())
	}
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
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/header/testing", r)
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Add(INPUT, desiredVal)

	value := GetHeaderValue(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

