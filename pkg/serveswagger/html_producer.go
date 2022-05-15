package serveswagger

import (
	"log"
	"io"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
)

func HTMLProducer(w io.Writer, data interface{}) error {
	var (
		t                  = common.Templates["underConstruction.gohtml"]
		params interface{} = SwaggerParams
	)
	if str, ok := data.(string); ok {
		for _, r := range common.AllRoutes {
			log.Println("loading template file:", r.TmplFile)
			log.Println("route Base:", r.Base)
			if str != r.Base {
				continue
			}
			tmpl, ok := common.Templates[r.TmplFile]
			if !ok {
				break
			}
			t = tmpl
			params = common.Parameters{
				ConstParams: SwaggerParams,
				Name:        r.Base,
			}
		}
	}

	err := t.ExecuteTemplate(w, "layout.gohtml", params)
	if err != nil {
		log.Println(err.Error())
	}
	return nil
}