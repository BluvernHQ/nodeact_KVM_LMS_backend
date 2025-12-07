package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"firebase.google.com/go/v4/auth"
)

type RequestPayload struct {
	Query      map[string]interface{} `json:"query"`
	Projection map[string]interface{} `json:"projection"`
	Collection string                 `json:"collection"`
}

func MapToBsonM(input map[string]interface{}) (bson.M, error) {
	data, err := bson.Marshal(input)
	if err != nil {
		return nil, err
	}
	var result bson.M
	err = bson.Unmarshal(data, &result)
	return result, err
}

func convertDatesDeep(data interface{}) interface{} {
	switch val := data.(type) {
	case map[string]interface{}:
		for k, v := range val {
			if k == "$date" {
				if s, ok := v.(string); ok {
					if t, err := time.Parse(time.RFC3339, s); err == nil {
						return t
					}
				}
			} else {
				val[k] = convertDatesDeep(v)
			}
		}
		return val
	case []interface{}:
		for i, v := range val {
			val[i] = convertDatesDeep(v)
		}
		return val
	default:
		return val
	}
}

func ExportDocs(w http.ResponseWriter, r *http.Request, db *mongo.Client, authClient *auth.Client) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
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
	if Role["Role"] == nil {
		http.Error(w, "Role not found", http.StatusNotFound)
		return
	}
	if Role["Role"] != "admin" && Role["Role"] != "staff" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload RequestPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	proj, err := MapToBsonM(payload.Projection)
	if err != nil {
		http.Error(w, "Invalid projection", http.StatusBadRequest)
		return
	}

	findOptions := options.Find().
		SetProjection(proj)
	collection = db.Database("KVM").Collection(payload.Collection)
	filter := convertDatesDeep(payload.Query)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Error querying database", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	var results []map[string]interface{}
	if err := cursor.All(ctx, &results); err != nil {
		http.Error(w, "Error reading results", http.StatusInternalServerError)
		return
	}

	// Convert results to Excel format
	filename := fmt.Sprintf("%s_export_%d.xlsx", payload.Collection, time.Now().Unix())

	// Create Excel file in memory
	file := excelize.NewFile()
	sheetName := "Sheet1"

	if len(results) > 0 {
		// Get headers from first document
		headers := make([]string, 0)
		for key := range results[0] {
			headers = append(headers, key)
		}

		// Write headers
		for i, header := range headers {
			cell := fmt.Sprintf("%c1", 'A'+i)
			file.SetCellValue(sheetName, cell, header)
		}

		// Write data rows
		for rowIdx, result := range results {
			for colIdx, header := range headers {
				cell := fmt.Sprintf("%c%d", 'A'+colIdx, rowIdx+2)
				value := result[header]
				file.SetCellValue(sheetName, cell, value)
			}
		}
	}

	// Save Excel file to /var/www/objects directory
	filePath := fmt.Sprintf("/var/www/objectfiles/%s", filename)
	if err := file.SaveAs(filePath); err != nil {
		http.Error(w, "Error saving Excel file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Start goroutine to delete file after 2 minutes
	go func() {
		time.Sleep(2 * time.Minute)
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Error deleting file %s: %v\n", filePath, err)
		} else {
			fmt.Printf("Successfully deleted file: %s\n", filePath)
		}
	}()

	// Return success response with file path
	response := map[string]interface{}{
		"success":    true,
		"message":    "File exported successfully",
		"url":        fmt.Sprintf("https://kvmtcc.org/api/documents/%s", filename),
		"expires_in": "2 minutes",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
