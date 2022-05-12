package pathtraversal

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"os"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes(frameworkSinks []common.Sink) {
	sinks := []common.Sink{
		{
			Name:                "os.ReadFile",
			Method:              "GET",
			Sanitize:            url.QueryEscape,
			VulnerableFnWrapper: osReadFile,
		},
		{
			Name:                "os.Open",
			Method:              "GET",
			Sanitize:            url.QueryEscape,
			VulnerableFnWrapper: osOpen,
		},
		{
			Name:                "os.WriteFile",
			Method:              "GET",
			Sanitize:            url.QueryEscape,
			VulnerableFnWrapper: osWriteFile,
		},
		{
			Name:                "os.Create",
			Method:              "GET",
			Sanitize:            url.QueryEscape,
			VulnerableFnWrapper: osCreate,
		},
	}
	if len(frameworkSinks) > 0 {
		sinks = append(sinks, frameworkSinks...)
	}
	common.Register(common.Route{
		Name:     "Path Traversal",
		Link:     "https://owasp.org/www-community/attacks/Path_Traversal",
		Base:     "pathTraversal",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "headers", "body"},
		Sinks:    sinks,
		Payload:  "../../../../../../../../../../../../etc/passwd",
	})
}

// read the given file using os.ReadFile
func osReadFile(_ interface{}, payload string) (data template.HTML, err error) {
	var content []byte
	content, err = os.ReadFile(payload)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if len(content) == 0 {
		return template.HTML(fmt.Sprintf("successfully read from %s; 0 bytes returned", payload)), nil
	}
	return template.HTML(content), nil
}

// read the given file using os.Open and bytes.Buffer
func osOpen(_ interface{}, payload string) (data template.HTML, err error) {
	fr, err := os.Open(payload)
	if err != nil {
		return "", fmt.Errorf("os.Open: error %w", err)
	}
	defer fr.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(fr)
	if err != nil {
		return "", fmt.Errorf("bytes.(Buffer).ReadFrom: error %w", err)
	}
	return template.HTML(buf.String()), nil
}

// write to the given file using os.WriteFile
func osWriteFile(_ interface{}, payload string) (data template.HTML, err error) {
	return "", os.WriteFile(payload, []byte("writing to file via os.WriteFile"), 0644)
}

// write to the given file using os.Create
func osCreate(_ interface{}, payload string) (data template.HTML, err error) {
	var buf bytes.Buffer
	fmt.Fprint(&buf, "writing to file via os.Create")

	fr, err := os.Create(payload)
	if err != nil {
		return "", err
	}
	defer fr.Close()

	_, err = buf.WriteTo(fr)
	return "", err
}
