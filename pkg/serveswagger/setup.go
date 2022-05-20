package serveswagger

import (
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/cmd_injection"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/go-openapi/runtime"
	"log"
	"os"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/injection/cmdi"
	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"
)

// DefaultAddr holds default localhost info
const DefaultAddr = "localhost:8080"

// SwaggerParams holds default ConstParams for the go-swagger executable
var SwaggerParams = common.ConstParams{
	Year:      2022,
	Logo:      "https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png",
	Framework: "Go-Swagger",
	Addr: 		DefaultAddr,
}

// Setup sets up the configuration for the go-swagger server
func Setup() (*restapi.Server, error) {
	// load up the swagger spec.
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	// set up the handlers for the api
	api := operations.NewSwaggerBenchAPI(swaggerSpec)

	api.HTMLProducer = runtime.ProducerFunc(HTMLProducer)

	api.SwaggerServerRootHandler = swagger_server.RootHandlerFunc(SwaggerRootHandler)

	api.CmdInjectionCmdInjectionFrontHandler = cmd_injection.CmdInjectionFrontHandlerFunc(CmdInjectionFront)

	api.CmdInjectionGetQueryCommandHandler = cmd_injection.GetQueryCommandHandlerFunc(GetQueryCommand)

	api.CmdInjectionGetQueryCommandContextHandler = cmd_injection.GetQueryCommandContextHandlerFunc(GetQueryCommandContext)

	server := restapi.NewServer(api)

	parser := flags.NewParser(server, flags.Default)
	parser.ShortDescription = "go swagger server"
	parser.LongDescription = "an intentionally vulnerable app built with go-swagger"
	server.ConfigureFlags()
	for _, optsGroup := range api.CommandLineOptionsGroups {
		_, err := parser.AddGroup(optsGroup.ShortDescription, optsGroup.LongDescription, optsGroup.Options)
		if err != nil {
			log.Fatalln(err)
		}
	}

	if _, err := parser.Parse(); err != nil {
		code := 1
		if fe, ok := err.(*flags.Error); ok {
			if fe.Type == flags.ErrHelp {
				code = 0
			}
		}
		os.Exit(code)
	}

	server.ConfigureAPI()
	server.Port = 8080

	// set up currently supported routes and resources
	if err := common.ParseViewTemplates(); err != nil {
		return nil, err
	}

	cmdi.RegisterRoutes("go-swagger")
	SwaggerParams.Rulebar = common.PopulateRouteMap(common.AllRoutes)

	return server, nil
}
