package main

import (
    "fmt"
    "io"
    "net/http"
	"time"
    "os"

    "go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"context"
	"encoding/json"
	


	"firebase.google.com/go/v4/auth"
)

func Upload(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
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

	collection := db.Database("KVM").Collection("Users")
	var Role bson.M
	err = collection.FindOne(ctx, bson.M{"UID": token.UID}, options.FindOne().SetProjection(bson.M{"Role": 1, "_id": 0})).Decode(&Role)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Println(Role)
	if Role["Role"] == nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}
	if Role["Role"] != "admin" && Role["Role"] != "DEVELOPER" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

    // Parse file
    err = r.ParseMultipartForm(100 << 20) // 100 MB
    if err != nil {
        http.Error(w, "File too big", http.StatusBadRequest)
        return
    }

    file, handler, err := r.FormFile("file")
    if err != nil {
        http.Error(w, "Error retrieving file", http.StatusBadRequest)
        return
    }
    defer file.Close()


    randomName := fmt.Sprintf("%d_%s", time.Now().UnixNano(), handler.Filename)
    f, err := os.Create("/var/www/objectfiles" + randomName)
    if err != nil {
        http.Error(w, "Failed to save file", http.StatusInternalServerError)
        return
    }
    defer f.Close()
    io.Copy(f, file)

	details := bson.M{
        "UID": token.UID,
		"FileName": randomName,
		"TimeStamp": time.Now(),
	}


	collection = db.Database("KVM").Collection("Documents")
	_, err = collection.InsertOne(ctx, details)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}

    w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Upload successful",
		"document URL": "http://api.kvmtcc.org/documents/objectfiles/" + randomName,
	})
}