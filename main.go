package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"plugin"
	"strings"
	"time"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"

	"google.golang.org/api/option"
)

var funcMap = make(map[string]func(http.ResponseWriter, *http.Request, *mongo.Client, *auth.Client))
var funcMapMutex sync.RWMutex

var mongoClient *mongo.Client

var authClient *auth.Client

// func auth(token string) bool {
// 	//!!!!!!!!!!!!!!!
// 	return true
// }

func initMongo() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://rajanand:machine2003@localhost:27017/admin"))
	if err != nil {
		log.Fatal("Mongo Connect error:", err)
	}

	mongoClient = client
	fmt.Println("MongoDB connected!")
}

func loadFunction(path string, functionName string) (func(http.ResponseWriter, *http.Request, *mongo.Client, *auth.Client), error) {
	// Open the plugin file
	p, err := plugin.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error loading plugin: %v", err)
	}

	// Look up the function
	sym, err := p.Lookup(functionName)
	if err != nil {
		return nil, fmt.Errorf("Function not found: %v", err)
	}

	// Assert function type
	fn, ok := sym.(func(http.ResponseWriter, *http.Request, *mongo.Client, *auth.Client))
	if !ok {
		return nil, fmt.Errorf("Invalid function signature")
	}

	return fn, nil
}

func loadEndpoints() {
	files, err := os.ReadDir("modules_bin")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".so" {
			continue
		}
		funcName := file.Name()[:len(file.Name())-3]

		funcPnt, err := loadFunction("modules_bin/"+file.Name(), funcName)
		if err != nil {
			fmt.Println("Invalid function signature")
			return
		}
		funcMap[funcName] = funcPnt
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {

	path := r.URL.Path
	segments := strings.Split(strings.Trim(path, "/"), "/")

	CallType := segments[0]

	if CallType != "" {

		funcMapMutex.RLock()
		tryFunc, ok := funcMap[CallType]
		funcMapMutex.RUnlock()

		if !ok {
			w.Header().Set("Content-Type", "text/plain")
			fmt.Fprintln(w, "Not registered function")
			return
		}
		tryFunc(w, r, mongoClient, authClient)
	} else {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprintln(w, "invalid")
		return
	}
}

func re_load_func(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	token, err := authClient.VerifyIDToken(ctx, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Token verified:", token.UID)

	collection := mongoClient.Database("KVM").Collection("Users")
	var Role bson.M
	err = collection.FindOne(ctx, bson.M{"UID": token.UID}, options.FindOne().SetProjection(bson.M{"Role": 1, "_id": 0})).Decode(&Role)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if Role["Role"] == nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}
	if Role["Role"] != "DEVELOPER" {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}

	funcName := r.URL.Query().Get("name")
	if funcName == "" {
		http.Error(w, "Missing function name", http.StatusBadRequest)
		return
	}

	pluginPath := "modules_bin/" + funcName + ".so"
	newFunc, err := loadFunction(pluginPath, funcName)
	if err != nil {
		http.Error(w, "Failed to reload function: "+err.Error(), http.StatusInternalServerError)
		return
	}

	funcMapMutex.Lock()
	funcMap[funcName] = newFunc
	funcMapMutex.Unlock()

	fmt.Fprintf(w, "Function '%s' reloaded successfully\n", funcName)
}

func main() {

	loadEndpoints()

	ctx := context.Background()
	opt := option.WithCredentialsFile("nodeact-kvm-firebase-adminsdk-fbsvc-c9fe0118fb.json")
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		log.Fatalf("error initializing auth: %v", err)
	}
	authClient, err = app.Auth(ctx)
	if err != nil {
		log.Fatalf("Firebase Auth error: %v", err)
	}

	initMongo()

	port := "0.0.0.0:5600"
	http.HandleFunc("/re_load_func", re_load_func)
	http.HandleFunc("/", httpHandler)

	fmt.Println("Server running on http://" + port)
	http.ListenAndServe(port, nil)
}
