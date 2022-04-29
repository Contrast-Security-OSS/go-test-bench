package serveswagger

import (
	"html/template"
	"net/http"
	"log"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/servestd"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
)

func SwaggerRootHandler(params swagger_server.RootParams) middleware.Responder {

	return CustomResponder(func(w http.ResponseWriter, producer runtime.Producer) {

		var t *template.Template
		//t = servestd.Templates["underConstruction.gohtml"]
		t = servestd.Templates["underConstruction.gohtml"]

		//err := t.ExecuteTemplate(w, "layout.gohtml", Pd)
		err := t.ExecuteTemplate(w, "underConstruction.gohtml", Pd)

		if err != nil {
			log.Print(err.Error())
		}

		payload := "Payload for the swagger root"
		if err := producer.Produce(w, payload); err != nil {
			panic(err) // let the recovery middleware deal with this
		}
	})
}
