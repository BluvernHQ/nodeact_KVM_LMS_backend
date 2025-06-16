package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	// "io"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"

	"firebase.google.com/go/v4/auth"
)

type RequestPayload struct {
	Id string   `json:"Id" bson:"Id"`
	Collection string `json:"Collection" bson:"Collection"`
	set string `json:"set" bson:"set"`
}

func DeleteBatch(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
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

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

collection = db.Database("KVM").Collection(payload.Collection)
objectID, err := primitive.ObjectIDFromHex(payload.Id)
    if err != nil {
        log.Fatal("Invalid ObjectID:", err)
    }

var update bson.M
if err := json.Unmarshal([]byte(payload.set), &update); err != nil {
	http.Error(w, "Invalid set payload: "+err.Error(), http.StatusBadRequest)
	return
}

result, err := collection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": update})
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if result.ModifiedCount == 0 {
		http.Error(w, "No documents were updated", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	response := map[string]string{"message": "Doc Updated successfully"}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Content-Type", "application/json")
}
