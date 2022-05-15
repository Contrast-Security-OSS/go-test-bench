package serveswagger

import (
	"net/http"
	"log"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

func SwaggerRootHandler(params swagger_server.RootParams) middleware.Responder {
	return middleware.ResponderFunc(func(w http.ResponseWriter, p runtime.Producer) {
		t := common.Templates["index.gohtml"]

		w.Header().Set("Application-Framework", "Go-Swagger")
		err := t.ExecuteTemplate(w, "layout.gohtml", SwaggerParams)
		if err != nil {
			log.Print(err.Error())
		}
	})
}
