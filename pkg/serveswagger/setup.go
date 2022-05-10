package serveswagger

import (
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
)

var SwaggerParams = common.ConstParams{
	Year:      2022,
	Logo:      "https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png",
	Framework: "Go-Swagger",
	Addr: 		"localhost:8080",
}

func Setup() error {
	if err := common.ParseViewTemplates(); err != nil {
		return err
	}

	cmdi.RegisterRoutes("go-swagger")
	SwaggerParams.Rulebar = common.PopulateRouteMap(common.AllRoutes)
	return nil
}
