package nosqlInjection

import (
	"net/http"
	"html/template"
	"bytes"
	"strings"
	"log"
	"context"
	"fmt"
	"time"
	utils "bitbucket.org/contrastsecurity/go-test-apps/go-test-bench/utils"
	mongo "go.mongodb.org/mongo-driver/mongo" 
	bson "go.mongodb.org/mongo-driver/bson"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

var templates = template.Must(template.ParseFiles("./views/partials/safeButtons.gohtml","./views/pages/nosqlInjection.gohtml", "./views/partials/ruleInfo.gohtml"))

var mongoOptions *options.ClientOptions
var auth options.Credential
var mongoClient *mongo.Client
var mongoDB *mongo.Database
var collection *mongo.Collection
var ctx  context.Context
var cancel context.CancelFunc

func MongoInit() {
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://root:rootpassword@mongo:27017"))
	if err != nil {
		log.Fatal(err)
	}
	mongoDB = mongoClient.Database("go-test-bench")
	collection = mongoDB.Collection("colors")
	
	opts := options.InsertMany().SetOrdered(false)
	docs := []interface{}{
		bson.D{{"name", "Alice"}},
		bson.D{{"name", "Bob"}},
	}
	collection.InsertMany(context.TODO(), docs, opts)
}

func MongoKill() {
	collection.Drop(ctx)
	mongoDB.Drop(ctx)
	mongoClient.Disconnect(ctx)
	cancel()
}

func mongoDBHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, mode string) (template.HTML, bool) {
	form_value := r.FormValue("input")

	switch mode {
		case "safe": 
			opts := options.Find()
			cursor, err := collection.Find(context.TODO(), bson.D{{"name" , form_value }} , opts)
			if err != nil {
				log.Fatal(err)
			}
			var results []bson.M
			if err = cursor.All(context.TODO(), &results); err != nil {
				log.Fatal(err)
			}
			var output string
			output = fmt.Sprintln(results)

			return template.HTML(output), false

		case "unsafe":
			opts := options.Find()
			cursor, err := collection.Find(context.TODO(), bson.D{{"$where" , "this.name == \"" + form_value + "\"" }} , opts)
			if err != nil {
				log.Fatal(err)
			}
			var results []bson.M
			if err = cursor.All(context.TODO(), &results); err != nil {
				log.Fatal(err)
			}
			var output string
			output = fmt.Sprintln(results)

			return template.HTML(output), false

		case "noop":
			return template.HTML("NOOP"), false

		default:
			log.Fatal("Error running mongoDBHandler. No option passed")
	}

	return template.HTML("?"), false
}

func defaultHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer
	err := templates.ExecuteTemplate(&buf, "nosql", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true 

}


func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	split_url := strings.Split(r.URL.Path, "/")
	switch split_url[2] {
		
    case "query":
      switch split_url[3]{
        case "mongodbCollectionFind":
          return mongoDBHandler(w,r,pd.Rulebar[pd.Name],split_url[len(split_url) - 1])
        default:
          log.Println("noSQL Injection Handler reached incorrectly") //It may not be the best to Kill here
			    return "", false
      }
			
		case "":
			return defaultHandler(w, r, pd.Rulebar[pd.Name])
		default:
			log.Fatal("noSQL Injection Handler reached incorrectly") //or here
			return "", false
		}
}