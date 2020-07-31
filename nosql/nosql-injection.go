package nosql

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Contrast-Security-OSS/go-test-bench/utils"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

var templates = template.Must(template.ParseFiles(
	"./views/partials/safeButtons.gohtml",
	"./views/pages/nosqlInjection.gohtml",
	"./views/partials/ruleInfo.gohtml",
))

var mongoClient *mongo.Client
var mongoDB *mongo.Database
var collection *mongo.Collection
var ctx context.Context
var cancel context.CancelFunc

// MongoInit creates a connection to MongoDB and instantiates
// the data in the database
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
		bson.D{{Key: "name", Value: "Alice"}},
		bson.D{{Key: "name", Value: "Bob"}},
	}
	_, _ = collection.InsertMany(context.TODO(), docs, opts)
}

//MongoKill cleans up the mongo database by dropping the data and
// disconnecting the client
func MongoKill() {
	_ = collection.Drop(ctx)
	_ = mongoDB.Drop(ctx)
	_ = mongoClient.Disconnect(ctx)
	cancel()
}

func mongoDBHandler(w http.ResponseWriter, r *http.Request, routeInfo utils.Route, mode string) (template.HTML, bool) {
	formValue := r.FormValue("input")

	switch mode {
	case "safe":
		opts := options.Find()
		cursor, err := collection.Find(
			context.TODO(),
			bson.D{{Key: "name", Value: formValue}},
			opts,
		)

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
		cursor, err := collection.Find(
			context.TODO(),
			bson.D{{Key: "$where", Value: "this.name == \"" + formValue + "\""}},
			opts,
		)
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

func nosqlTemplate(w http.ResponseWriter, r *http.Request, routeInfo utils.Route) (template.HTML, bool) {
	var buf bytes.Buffer

	err := templates.ExecuteTemplate(&buf, "nosql", routeInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	return template.HTML(buf.String()), true

}

// Handler is the nosql endpoint API handler
func Handler(w http.ResponseWriter, r *http.Request, pd utils.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	switch splitURL[2] {

	case "query":
		switch splitURL[3] {
		case "mongodbCollectionFind":
			return mongoDBHandler(w, r, pd.Rulebar[pd.Name], splitURL[len(splitURL)-1])
		default:
			log.Println("noSQL Injection Handler reached incorrectly") //It may not be the best to Kill here
			return "", false
		}

	case "":
		return nosqlTemplate(w, r, pd.Rulebar[pd.Name])
	default:
		log.Fatal("noSQL Injection Handler reached incorrectly") //or here
		return "", false
	}
}
