package pathtraversal

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

// RegisterRoutes is to be called to add the routes in this package to common.AllRoutes.
func RegisterRoutes(frameworkSinks ...*common.Sink) {
	sinks := []*common.Sink{
		{
			Name:                "os.ReadFile",
			Sanitize:            url.QueryEscape,
			VulnerableFnWrapper: osReadFile,
		},
		{
			Name:                "os.Open",
			Sanitize:            url.QueryEscape,
			VulnerableFnWrapper: osOpen,
		},
		{
			Name:                 "os.WriteFile",
			Sanitize:             url.QueryEscape,
			VulnerableFnWrapper:  osWriteFile,
			ExpectedUnsafeStatus: http.StatusBadRequest,
		},
		{
			Name:                 "os.Create",
			Sanitize:             url.QueryEscape,
			VulnerableFnWrapper:  osCreate,
			ExpectedUnsafeStatus: http.StatusBadRequest,
		},
	}
	if len(frameworkSinks) > 0 {
		sinks = append(sinks, frameworkSinks...)
	}
	payload := "../../../../../../../../../../../../etc/passwd"
	if runtime.GOOS == "windows" {
		views, err := common.FindViewsDir()
		if err != nil {
			log.Fatalf("finding path to file for path traversal: %s", err)
		}
		// we know this file exists
		payload = filepath.Clean(views + `\..\internal\pathtraversal\secrets.txt`)
		for i := range sinks {
			if sinks[i].Name == "os.WriteFile" || sinks[i].Name == "os.Create" {
				//because we're using a writeable file, these will succeed.
				sinks[i].ExpectedUnsafeStatus = http.StatusOK
			}
		}
	}
	common.Register(common.Route{
		Name:     "Path Traversal",
		Link:     "https://owasp.org/www-community/attacks/Path_Traversal",
		Base:     "pathTraversal",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "buffered-query", "headers", "body"},
		Sinks:    sinks,
		Payload:  payload,
	})
}

// read the given file using os.ReadFile
func osReadFile(_ interface{}, payload string) (data string, raw bool, err error) {
	var content []byte
	content, err = os.ReadFile(payload)
	if err != nil {
		log.Println(err)
		return "", false, err
	}
	if len(content) == 0 {
		return fmt.Sprintf("successfully read from %s; 0 bytes returned", payload), false, nil
	}
	return string(content), false, nil
}

// read the given file using os.Open and bytes.Buffer
func osOpen(_ interface{}, payload string) (data string, raw bool, err error) {
	fr, err := os.Open(payload)
	if err != nil {
		return "", false, fmt.Errorf("os.Open: error %w", err)
	}
	defer fr.Close()

	var buf bytes.Buffer
	_, err = buf.ReadFrom(fr)
	if err != nil {
		return "", false, fmt.Errorf("bytes.(Buffer).ReadFrom: error %w", err)
	}
	return buf.String(), false, nil
}

// write to the given file using os.WriteFile
func osWriteFile(_ interface{}, payload string) (data string, raw bool, err error) {
	return "", false, os.WriteFile(payload, []byte("writing to file via os.WriteFile"), 0644)
}

// write to the given file using os.Create
func osCreate(_ interface{}, payload string) (data string, raw bool, err error) {
	buf := bytes.NewBufferString("writing to file via os.Create")

	fr, err := os.Create(payload)
	if err != nil {
		return "", false, err
	}
	defer fr.Close()

	_, err = buf.WriteTo(fr)
	return "", false, err
}
