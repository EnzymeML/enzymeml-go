// Package main demonstrates how to start the EnzymeML v2 REST API server.
//
// This example shows how to:
// - Set up a SQLite database with all EnzymeML v2 models
// - Initialize the database manager with proper schema migration
// - Create and configure the REST API handler
// - Start the HTTP server with full CRUD operations
//
// The API will be available at http://localhost:8080 with the following endpoints:
// - GET/POST/PUT/DELETE /v2/documents
// - GET/POST/PUT/DELETE /v2/creators
// - GET/POST/PUT/DELETE /v2/vessels
// - GET/POST/PUT/DELETE /v2/proteins
// - GET/POST/PUT/DELETE /v2/complexes
// - GET/POST/PUT/DELETE /v2/small-molecules
// - GET/POST/PUT/DELETE /v2/reactions
// - GET/POST/PUT/DELETE /v2/reaction-elements
// - GET/POST/PUT/DELETE /v2/modifier-elements
// - GET/POST/PUT/DELETE /v2/equations
// - GET/POST/PUT/DELETE /v2/variables
// - GET/POST/PUT/DELETE /v2/parameters
// - GET/POST/PUT/DELETE /v2/measurements
// - GET/POST/PUT/DELETE /v2/measurement-data
// - GET/POST/PUT/DELETE /v2/unit-definitions
// - GET/POST/PUT/DELETE /v2/base-units
//
// Additional endpoints:
// - GET /v2/docs - Interactive API documentation
// - GET /v2/openapi.json - OpenAPI specification in JSON format
// - GET /v2/openapi.yaml - OpenAPI specification in YAML format
package main

import (
	"log"
	"net/http"

	enzymeml_v2 "github.com/EnzymeML/enzymeml-go/src"
	apiv2 "github.com/EnzymeML/enzymeml-go/src/api/v2"
	database "github.com/EnzymeML/enzymeml-go/src/database/v2"
)

// main initializes and starts the EnzymeML v2 REST API server.
//
// The function performs the following steps:
// 1. Creates a new database manager with SQLite backend
// 2. Automatically migrates all EnzymeML v2 model schemas
// 3. Initializes the REST API handler with CRUD operations
// 4. Starts the HTTP server on port 8080
func main() {
	// Initialize the database manager with SQLite backend.
	// The database file "enzymeml.db" will be created in the current directory.
	// All models will be automatically migrated to create the necessary tables.
	dbManager, err := database.NewDBManager("enzymeml.db", models())
	if err != nil {
		log.Fatalf("failed to create database manager: %v", err)
	}
	// Ensure the database connection is properly closed when the program exits
	defer dbManager.Close()

	// Create the REST API handler with full CRUD operations for all EnzymeML v2 models.
	// The handler uses Huma v2 framework for automatic OpenAPI documentation generation.
	handler, err := apiv2.NewHandler(dbManager)
	if err != nil {
		log.Fatalf("failed to create v2 API handler: %v", err)
	}

	// Start the HTTP server on port 8080
	log.Println("EnzymeML v2 REST API listening on :8080")
	log.Println("API Documentation available at: http://localhost:8080/v2/docs")
	log.Println("OpenAPI Spec available at: http://localhost:8080/v2/openapi.json")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

// models returns a slice of all EnzymeML v2 model types that need to be
// registered with the database manager for schema migration and ORM operations.
//
// These models represent the complete EnzymeML v2 specification and include:
// - EnzymeMLDocument: Root document containing all experimental data
// - Creator: Information about document creators/authors
// - Vessel: Reaction vessels and their properties
// - Protein: Protein/enzyme definitions
// - Complex: Protein complexes and their compositions
// - SmallMolecule: Small molecule substrates, products, and inhibitors
// - Reaction: Chemical reactions and their participants
// - ReactionElement: Individual reactants in reactions (substrates/products)
// - ModifierElement: Reaction modifiers (inhibitors, activators, etc.)
// - Equation: Mathematical equations for kinetic models
// - Variable: Variables used in equations
// - Parameter: Parameters for kinetic models
// - Measurement: Experimental measurements and conditions
// - MeasurementData: Time-series data points from measurements
// - UnitDefinition: Custom unit definitions
// - BaseUnit: Base units for unit definitions
//
// Each model is passed as a pointer to enable proper GORM reflection
// and database schema generation.
func models() []interface{} {
	return []interface{}{
		&enzymeml_v2.EnzymeMLDocument{}, // Root document structure
		&enzymeml_v2.Creator{},          // Document authors/creators
		&enzymeml_v2.Vessel{},           // Reaction vessels
		&enzymeml_v2.Protein{},          // Proteins and enzymes
		&enzymeml_v2.Complex{},          // Protein complexes
		&enzymeml_v2.SmallMolecule{},    // Small molecules (substrates, products, etc.)
		&enzymeml_v2.Reaction{},         // Chemical reactions
		&enzymeml_v2.ReactionElement{},  // Reaction participants (substrates/products)
		&enzymeml_v2.ModifierElement{},  // Reaction modifiers (inhibitors, activators)
		&enzymeml_v2.Equation{},         // Mathematical equations for kinetic models
		&enzymeml_v2.Variable{},         // Variables in equations
		&enzymeml_v2.Parameter{},        // Model parameters
		&enzymeml_v2.Measurement{},      // Experimental measurements
		&enzymeml_v2.MeasurementData{},  // Time-series measurement data
		&enzymeml_v2.UnitDefinition{},   // Custom unit definitions
		&enzymeml_v2.BaseUnit{},         // Base units for unit definitions
	}
}
