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

	"firebase.google.com/go/v4/auth"
)

type Staff struct {
	UID                 string   `json:"UID" bson:"UID"`
	Name                string   `json:"Name" bson:"Name"`
	Role                string   `json:"Role" bson:"Role"`
	TimeStamp           string   `json:"TimeStamp" bson:"TimeStamp"`
	DOB                 string   `json:"DOB" bson:"DOB"`
	Qualification       string   `json:"Qualification" bson:"Qualification"`
	Subjects            []string `json:"Subjects" bson:"Subjects"`
	Experience          string   `json:"Experience" bson:"Experience"`
	Phone               string   `json:"Phone" bson:"Phone"`
	WorkingAt           string   `json:"WorkingAt" bson:"WorkingAt"`
	OtherSpecialization string   `json:"OtherSpecialization" bson:"OtherSpecialization"`
	BatchId             []string `json:"Batch" bson:"Batch"`
	ProfilePic          string   `json:"ProfilePic" bson:"ProfilePic"`
}

func CreateStaff(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
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

	var data Staff
	err = json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	params := (&auth.UserToCreate{}).
		Email(r.Header.Get("Add-User-Name")).
		Password(r.Header.Get("Add-User-Pwd"))
	userRecord, err := authClient.CreateUser(context.Background(), params)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating UID: %v", err), http.StatusInternalServerError)
		return
	}
	data.UID = userRecord.UID

	data.Role = "staff"

	result, err := collection.InsertOne(ctx, data)
	if err != nil {
		http.Error(w, "Insert error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":    "User created successfully",
		"insertedID": result.InsertedID,
		"UID":        data.UID,
	})
}
