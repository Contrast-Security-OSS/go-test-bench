package main

var rules = map[string]Route{
	"cmdInjection": {
		Base:     "/cmdInjection",
		Name:     "Command Injection",
		Link:     "https://www.owasp.org/index.php/Command_Injection",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "cookies"},
		Sinks: []Sink{
			{
				Name:   "os.Exec",
				URL:    "/cmdInjection/osExec",
				Method: "GET",
			},
		},
	},

	// "nosqlInjection": {
	// 	Base:     "/nosqlInjection",
	// 	Name:     "NoSQL Injection",
	// 	Link:     "https://www.owasp.org/index.php/Testing_for_NoSQL_injection",
	// 	Products: []string{"Assess", "Protect"},
	// 	Inputs:   []string{"query"},
	// 	Sinks: []Sink{{
	// 		Name:   "mongodb.Collection.Find",
	// 		URL:    "/nosqlInjection/query/mongodbCollectionFind",
	// 		Method: "GET",
	// 	}},
	// },

	"pathTraversal": {
		Base:     "/pathTraversal",
		Name:     "Path Traversal",
		Link:     "https://owasp.org/www-community/attacks/Path_Traversal",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "headers", "body"},
		Sinks: []Sink{
			{
				Name:   "gin.File",
				URL:    "/pathTraversal",
				Method: "GET",
			},
			{
				Name:   "ioutil.ReadFile",
				URL:    "/pathTraversal",
				Method: "GET",
			},
			{
				Name:   "ioutil.WriteFile",
				URL:    "/pathTraversal",
				Method: "GET",
			},
		},
	},

	"sqlInjection": {
		Base:     "/sqlInjection",
		Name:     "SQL Injection",
		Link:     "https://www.owasp.org/index.php/SQL_Injection",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "headers-json", "body"},
		Sinks: []Sink{{
			Name:   "sqlit3.exec",
			URL:    "/query/sqlite3Exec",
			Method: "GET",
		}},
	},

	"ssrf": {
		Base:     "/ssrf",
		Name:     "Server Side Request Forgery",
		Link:     "https://owasp.org/www-community/attacks/Server_Side_Request_Forgery",
		Products: []string{"Assess"},
		Inputs:   []string{"query"},
	},

	"unvalidatedRedirect": {
		Base:     "/unvalidatedRedirect",
		Name:     "Unvalidated Redirect",
		Link:     "https://cheatsheetseries.owasp.org/cheatsheets/Unvalidated_Redirects_and_Forwards_Cheat_Sheet.html",
		Products: []string{"Assess"},
		Inputs:   []string{"query"},
		Sinks: []Sink{{
			Name:   "gin.Redirect",
			URL:    "/unvalidatedRedirect/gin.Redirect/",
			Method: "GET",
		}},
	},

	// "xpathInjection": {
	// 	Base:     "/xpathInjection",
	// 	Name:     "XPath Injection",
	// 	Link:     "https : //owasp.org/www-community/attacks/XPATH_Injection",
	// 	Products: []string{"Assess"},
	// 	Inputs:   []string{"query"},
	// },

	"xss": {
		Base:     "/xss",
		Name:     "Reflected XSS",
		Link:     "https://www.owasp.org/index.php/Cross-site_Scripting_(XSS)#Stored_and_Reflected_XSS_Attacks",
		Products: []string{"Assess", "Protect"},
		Inputs:   []string{"query", "params", "body"},
	},

	// "xssJSON": {
	// 	Base:     "/xssJSON",
	// 	Name:     "Reflected XSS JSON (Safe)",
	// 	Link:     "https://www.owasp.org/index.php/Cross-site_Scripting_(XSS)#Stored_and_Reflected_XSS_Attacks",
	// 	Products: []string{"Assess", "Protect"},
	// 	Inputs:   []string{"query", "params"},
	// },

	// "xssStealthyRequire": {
	// 	Base:     "/xssStealthyRequire",
	// 	Name:     "Reflected XSS (stealthy-require)",
	// 	Link:     "https://www.owasp.org/index.php/Cross-site_Scripting_(XSS)#Stored_and_Reflected_XSS_Attacks",
	// 	Products: []string{"Assess", "Protect"},
	// 	Inputs:   []string{"query", "params"},
	// },

	// "xxe": {
	// 	Base:     "/xxe",
	// 	Name:     "XXE Processing",
	// 	Link:     "https://www.owasp.org/index.php/XML_External_Entity_(XXE)_Processing",
	// 	Products: []string{"Assess", "Protect"},
	// 	Inputs:   []string{"query"},
	// },

	// "paramPollution": {
	// 	Base:     "/parampollution",
	// 	Name:     "HTTP Parameter Pollution / Cache Controls Missing",
	// 	Link:     "https://owasp.org/www-pdf-archive/AppsecEU09_CarettoniDiPaola_v0.8.pdf",
	// 	Products: []string{"Assess"},
	// 	Inputs: []string{"test"},
	// },
}

// Sink is a struct that identifies the name
// of the sink, the associated URL and the
// HTTP method
type Sink struct {
	Name   string
	URL    string
	Method string
}

// Route is the template information for a specific route
type Route struct {
	Base     string
	Name     string
	Link     string
	Products []string
	Inputs   []string
	Sinks    []Sink
}
