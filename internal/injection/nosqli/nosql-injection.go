package nosql

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Contrast-Security-OSS/go-test-bench/internal/common"
	bson "go.mongodb.org/mongo-driver/bson"
	mongo "go.mongodb.org/mongo-driver/mongo"
	options "go.mongodb.org/mongo-driver/mongo/options"
)

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
		log.Printf("Could not connect the Mongo client: err = %s", err)
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

func mongoDBHandler(w http.ResponseWriter, r *http.Request, mode string) (template.HTML, bool) {
	formValue := common.GetUserInput(r)

	switch mode {
	case "safe":
		opts := options.Find()
		cursor, err := collection.Find(
			context.TODO(),
			bson.D{{Key: "name", Value: formValue}},
			opts,
		)

		if err != nil {
			log.Printf("Could not query Mongo in %s: err = %s", mode, err)
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Printf("Could not get cursor in %s: err = %s", mode, err)
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
			log.Printf("Could not query Mongo in %s: err = %s", mode, err)
		}
		var results []bson.M
		if err = cursor.All(context.TODO(), &results); err != nil {
			log.Printf("Could not get cursor in %s: err = %s", mode, err)
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

func nosqlTemplate(w http.ResponseWriter, r *http.Request) (template.HTML, bool) {
	return "nosqlInjection.gohtml", true
}

// Handler is the nosql endpoint API handler
func Handler(w http.ResponseWriter, r *http.Request, pd common.Parameters) (template.HTML, bool) {
	splitURL := strings.Split(r.URL.Path, "/")
	switch splitURL[2] {

	case "query":
		switch splitURL[3] {
		case "mongodbCollectionFind":
			return mongoDBHandler(w, r, splitURL[len(splitURL)-1])
		default:
			log.Println("noSQL Injection Handler reached incorrectly") //It may not be the best to Kill here
			return "", false
		}

	case "":
		return nosqlTemplate(w, r)
	default:
		log.Fatal("noSQL Injection Handler reached incorrectly") //or here
		return "", false
	}
}
