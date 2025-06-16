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
	"firebase.google.com/go/v4/auth"
)

type Class struct {
	Name			   string   `json:"Name" bson:"Name"`
	BatchId			   string   `json:"BatchId" bson:"BatchId"`
}

func CreateClass(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
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

	var data Class
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	collection = db.Database("KVM").Collection("Batches")
	idHex := data.BatchId
	if idHex == "" {
		http.Error(w, "Missing id parameter", http.StatusBadRequest)
		return
	}
	objID, err := primitive.ObjectIDFromHex(idHex)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	err = collection.FindOne(ctx, bson.M{"_id": objID}).Err()
	if err == mongo.ErrNoDocuments {
		http.Error(w, "Batch not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}


	collection = db.Database("KVM").Collection("Classes")
	err = collection.FindOne(ctx, bson.M{"Name": data.Name}).Decode(&Role)
	if err == mongo.ErrNoDocuments{
		result, err := collection.InsertOne(ctx, data)
		if err != nil {
			http.Error(w, "Insert error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "Class created successfully",
		"ID": result.InsertedID,
	})
	}else{
		http.Error(w, "Error creating Class", http.StatusConflict)
		return
	}
}
