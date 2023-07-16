package main
import (
	"encoding/json"
	"fmt"
	"github.com/tomchavakis/turf-go"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func intersectionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received POST request to /intersection")

	// Perform header-based authentication check
	authHeader := r.Header.Get("Authorization")
	expectedAuth := "your_expected_auth_value"
	if authHeader != expectedAuth {
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	// Parse the request body to extract the long linestring
	var longLinestring turf.Feature
	err := json.NewDecoder(r.Body).Decode(&longLinestring)
	if err != nil {
		http.Error(w, "Malformed request body", http.StatusBadRequest)
		return
	}

	// Generate random lines and check for intersections
	intersections := make(map[string][][2]float64)
	for i := 0; i < 50; i++ {
		// Generate random start and end points for the lines
		startPoint := turf.Point{Coordinates: [2]float64{randomFloat(), randomFloat()}}
		endPoint := turf.Point{Coordinates: [2]float64{randomFloat(), randomFloat()}}

		// Create the line feature
		line := turf.LineString{Coordinates: [][2]float64{{startPoint.Coordinates[0], startPoint.Coordinates[1]}, {endPoint.Coordinates[0], endPoint.Coordinates[1]}}}
		lineFeature := turf.Feature{Geometry: &line, Properties: map[string]interface{}{"id": fmt.Sprintf("L%02d", i+1)}}

		// Check if the long linestring intersects with the line
		if turf.Intersects(&longLinestring, &lineFeature) {
			intersections[lineFeature.Properties["id"].(string)] = line.Coordinates
		}
	}

	// Return the intersections or empty array
	if len(intersections) == 0 {
		fmt.Fprint(w, "[]")
	} else {
		json.NewEncoder(w).Encode(intersections)
	}
}

func randomFloat() float64 {
	return rand.Float64() * 180.0 - 90.0
}

func main() {
	rand.Seed(time.Now().UnixNano())

	http.HandleFunc("/intersection", intersectionHandler)

	log.Println("Starting server on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
