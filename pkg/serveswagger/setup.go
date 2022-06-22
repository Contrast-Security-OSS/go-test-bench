package serveswagger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations"
	"github.com/Contrast-Security-OSS/go-test-bench/cmd/go-swagger/restapi/operations/swagger_server"
	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	"github.com/Contrast-Security-OSS/go-test-bench/pkg/serveswagger/shared"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	flags "github.com/jessevdk/go-flags"
)

// DefaultAddr holds default localhost info
const DefaultAddr = "localhost:8080"

// SwaggerParams holds default ConstParams for the go-swagger executable
var SwaggerParams = common.ConstParams{
	Year:      time.Now().Year(),
	Logo:      "https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png",
	Framework: "Go-Swagger",
	Addr:      DefaultAddr,
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

	// set up currently supported routes and resources
	if err := common.ParseViewTemplates(); err != nil {
		return nil, err
	}
	// add new routes to this function, which is shared with cmd/go-swagger/regen
	shared.RegisterNewRoutes()

	rmap := common.PopulateRouteMap(common.AllRoutes)
	shared.FilterInputTypes(rmap)

	// set up template files for legacy routes
	for k, v := range rmap {
		if len(v.TmplFile) == 0 {
			v.TmplFile = strings.TrimLeft(v.Base+".gohtml", "/")
			rmap[k] = v
		}
	}

	// lives in generated code. initializes all route handlers other than root.
	generatedInit(api, rmap, &SwaggerParams)

	server := restapi.NewServer(api)

	server.Port = 8080 // put it up here so it can be overridden by flag

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

	SwaggerParams.Rulebar = rmap

	return server, nil
}

// RouteHandler returns a middleware.Responder that serves our html and the vulnerable functions.
func RouteHandler(rt common.Route, pd *common.ConstParams, req *http.Request) middleware.Responder {
	if len(rt.Base) == 0 {
		log.Fatalf("missing data for route: %#v\nurl=%s", rt, req.URL)
	}
	return &responder{
		rt: rt,
		params: common.Parameters{
			ConstParams: *pd,
			Name:        rt.Base,
		},
		req: req,
	}
}

type responder struct {
	rt     common.Route
	params common.Parameters
	req    *http.Request
}

func (r *responder) WriteResponse(w http.ResponseWriter, p runtime.Producer) {
	log.Println(r.rt.Name, r.req.URL.Path)
	elems := strings.Split(strings.Trim(r.req.URL.Path, "/"), "/")
	if len(elems) < 2 {
		// main page
		tmpl, ok := common.Templates[r.rt.TmplFile]
		if !ok {
			log.Printf("missing template %s for %#v", r.rt.TmplFile, r.rt)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		err := tmpl.ExecuteTemplate(w, "layout.gohtml", &r.params)
		if err != nil {
			log.Print(err.Error())
			fmt.Fprintf(w, "template error: %s", err)
		}
		return
	}
	for _, s := range r.rt.Sinks {
		if elems[1] == s.URL {
			mode := common.Safety(elems[len(elems)-1])
			switch mode {
			case common.NOOP, common.Safe, common.Unsafe:
				// valid modes
			default:
				// invalid
				w.WriteHeader(http.StatusNotFound)
				return
			}
			if s.Handler == nil {
				var err error
				s.Handler, err = common.GenericHandler(s)
				if err != nil {
					log.Fatal(err)
				}
			}
			in := common.GetUserInput(r.req)
			data, mime, status := s.Handler(mode, in, p)
			w.WriteHeader(status)
			if len(mime) == 0 {
				mime = "text/plain"
			}
			w.Header().Set("Content-Type", mime)
			w.Header().Set("Cache-Control", "no-store") //makes development a whole lot easier
			fmt.Fprint(w, data)
			return
		}
	}
	// does not match any sink or the main page
	w.WriteHeader(http.StatusNotFound)
}
