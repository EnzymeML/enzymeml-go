package database

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	enzymeml_v2 "github.com/EnzymeML/enzymeml-go/src"
	"gorm.io/gorm"
)

func TestNewDBManager_CreatesPathAndDB(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "nested", "dir", "enzymeml.db")

	manager, err := NewDBManager(dbPath, allV2Models())
	if err != nil {
		t.Fatalf("NewDBManager returned error: %v", err)
	}
	t.Cleanup(func() { _ = manager.Close() })

	if manager.DB() == nil {
		t.Fatalf("expected DB() to return initialized *gorm.DB")
	}

	if _, err := os.Stat(filepath.Dir(dbPath)); err != nil {
		t.Fatalf("expected database directory to exist: %v", err)
	}
	if _, err := os.Stat(dbPath); err != nil {
		t.Fatalf("expected database file to exist: %v", err)
	}
}

func TestDBManager_DocumentAndCreatorOperations(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "enzymeml.db")
	manager, err := NewDBManager(dbPath, allV2Models())
	if err != nil {
		t.Fatalf("NewDBManager returned error: %v", err)
	}
	t.Cleanup(func() { _ = manager.Close() })

	docOne := &enzymeml_v2.EnzymeMLDocument{
		Name:    "doc-1",
		Version: "2.0.0",
		Creators: []enzymeml_v2.Creator{
			{GivenName: "Ada", FamilyName: "Lovelace", Mail: "ada@example.com"},
		},
	}
	if err := manager.SaveEnzymeMLDocument(docOne); err != nil {
		t.Fatalf("SaveEnzymeMLDocument(docOne) error: %v", err)
	}

	docTwo := &enzymeml_v2.EnzymeMLDocument{
		Name:    "doc-2",
		Version: "2.0.0",
	}
	if err := manager.SaveEnzymeMLDocument(docTwo); err != nil {
		t.Fatalf("SaveEnzymeMLDocument(docTwo) error: %v", err)
	}

	got, err := manager.GetEnzymeMLDocumentByID(uint(docOne.Id))
	if err != nil {
		t.Fatalf("GetEnzymeMLDocumentByID error: %v", err)
	}
	if got.Name != "doc-1" {
		t.Fatalf("expected doc name doc-1, got %q", got.Name)
	}
	if len(got.Creators) != 1 {
		t.Fatalf("expected 1 creator preloaded, got %d", len(got.Creators))
	}

	allDocs, err := manager.GetAllEnzymeMLDocuments()
	if err != nil {
		t.Fatalf("GetAllEnzymeMLDocuments error: %v", err)
	}
	if len(allDocs) != 2 {
		t.Fatalf("expected 2 documents, got %d", len(allDocs))
	}

	limitedDocs, err := manager.GetEnzymeMLDocuments(1)
	if err != nil {
		t.Fatalf("GetEnzymeMLDocuments error: %v", err)
	}
	if len(limitedDocs) != 1 {
		t.Fatalf("expected 1 limited document, got %d", len(limitedDocs))
	}

	creator := &enzymeml_v2.Creator{
		GivenName:  "Grace",
		FamilyName: "Hopper",
		Mail:       "grace@example.com",
	}
	if err := manager.SaveCreator(creator); err != nil {
		t.Fatalf("SaveCreator error: %v", err)
	}

	gotCreator, err := manager.GetCreatorByID(uint(creator.Id))
	if err != nil {
		t.Fatalf("GetCreatorByID error: %v", err)
	}
	if gotCreator.FamilyName != "Hopper" {
		t.Fatalf("expected creator family name Hopper, got %q", gotCreator.FamilyName)
	}

	limitedCreators, err := manager.GetCreators(1)
	if err != nil {
		t.Fatalf("GetCreators error: %v", err)
	}
	if len(limitedCreators) != 1 {
		t.Fatalf("expected 1 limited creator, got %d", len(limitedCreators))
	}
}

func TestGetEnzymeMLDocumentByID_NotFound(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), "enzymeml.db")
	manager, err := NewDBManager(dbPath, allV2Models())
	if err != nil {
		t.Fatalf("NewDBManager returned error: %v", err)
	}
	t.Cleanup(func() { _ = manager.Close() })

	_, err = manager.GetEnzymeMLDocumentByID(9999)
	if err == nil {
		t.Fatalf("expected not-found error")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Fatalf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func allV2Models() []interface{} {
	return []interface{}{
		&enzymeml_v2.EnzymeMLDocument{},
		&enzymeml_v2.Creator{},
		&enzymeml_v2.Vessel{},
		&enzymeml_v2.Protein{},
		&enzymeml_v2.Complex{},
		&enzymeml_v2.SmallMolecule{},
		&enzymeml_v2.Reaction{},
		&enzymeml_v2.ReactionElement{},
		&enzymeml_v2.ModifierElement{},
		&enzymeml_v2.Equation{},
		&enzymeml_v2.Variable{},
		&enzymeml_v2.Parameter{},
		&enzymeml_v2.Measurement{},
		&enzymeml_v2.MeasurementData{},
		&enzymeml_v2.UnitDefinition{},
		&enzymeml_v2.BaseUnit{},
	}
}
