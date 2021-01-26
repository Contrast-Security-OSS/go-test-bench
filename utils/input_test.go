package utils

import (
	"net/http"
	"strings"
	"testing"
)

func TestGetQueryInput(t *testing.T) {
	desiredVal := "get_value"
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/getQuery/testing?input="+desiredVal, r)
	if err != nil {
		t.Error(err.Error())
	}
	value := GetQueryInput(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestGetParameterInput(t *testing.T) {
	desiredVal := "<scipt>alert(1);</script>"
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/getQuery/parameter/"+desiredVal+"/testing", r)
	if err != nil {
		t.Error(err.Error())
	}

	value := GetParameterInput(request, 3, 4)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestPostInput(t *testing.T) {
	desiredVal := "post_value"
	r := strings.NewReader("input=" + desiredVal)
	request, err := http.NewRequest(http.MethodPost, "/postQuery/testing", r)
	if err != nil {
		t.Error(err.Error())
	}
	value, _ := PostInput(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestFormValueInput(t *testing.T) {
	desiredVal := "form_value"
	r := strings.NewReader("input=" + desiredVal)
	request, err := http.NewRequest(http.MethodPost, "/form/testing", r)
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Set("Content-type", "application/x-www-form-urlencoded")

	value := FormValueInput(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestCookieInput(t *testing.T) {
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

	value := CookieInput(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

func TestHeaderInput(t *testing.T) {
	desiredVal := "header_value"
	r := strings.NewReader("test")
	request, err := http.NewRequest(http.MethodGet, "/header/testing", r)
	if err != nil {
		t.Error(err.Error())
	}
	request.Header.Add(INPUT, desiredVal)

	value := HeaderInput(request, INPUT)
	if value != desiredVal {
		t.Error("got:", value, ",want:", desiredVal)
	}
}

