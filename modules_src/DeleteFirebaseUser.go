package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"firebase.google.com/go/v4/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type RequestDeleteFirebaseUser struct {
	Email string `json:"Email" bson:"Email"`
}

func DeleteFirebaseUser(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Verify Firebase token
	token, err := authClient.VerifyIDToken(ctx, r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
		return
	}
	fmt.Println("Token verified:", token.UID)

	// Check if user has admin or developer role
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

	// Parse request body
	var payload RequestDeleteFirebaseUser
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Validate email
	if payload.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Get user by email from Firebase
	userRecord, err := authClient.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if auth.IsUserNotFound(err) {
			http.Error(w, "User not found with this email", http.StatusNotFound)
			return
		}
		http.Error(w, "Firebase error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete user from Firebase
	err = authClient.DeleteUser(ctx, userRecord.UID)
	if err != nil {
		http.Error(w, "Error deleting user from Firebase: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Also delete from MongoDB Users collection if exists
	collection = db.Database("KVM").Collection("Users")
	result, err := collection.DeleteOne(ctx, bson.M{"UID": userRecord.UID})
	if err != nil {
		// Log the error but don't fail the request since Firebase user was deleted
		fmt.Printf("Warning: Could not delete user from MongoDB: %v\n", err)
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]interface{}{
		"message":      "Firebase user deleted successfully",
		"email":        payload.Email,
		"uid":          userRecord.UID,
		"mongoDeleted": result.DeletedCount > 0,
	}
	json.NewEncoder(w).Encode(response)
}
