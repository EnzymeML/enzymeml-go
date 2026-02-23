package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	enzymeml_v2 "github.com/EnzymeML/enzymeml-go/src"
	database "github.com/EnzymeML/enzymeml-go/src/database/v2"
	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"
	"gorm.io/gorm"
)

type modelRoute struct {
	Path     string
	Model    interface{}
	Preloads []string
}

type crudHandler struct {
	db            *gorm.DB
	resourcePath  string
	modelType     reflect.Type
	primaryKeyIdx int
	primaryKeyTyp reflect.Type
	preloads      []string
}

type createInput struct {
	Body map[string]interface{}
}

type idInput struct {
	ID string `path:"id"`
}

type updateInput struct {
	ID   string `path:"id"`
	Body map[string]interface{}
}

type createdOutput struct {
	Status int                    `status:"201"`
	Body   map[string]interface{} `json:"body"`
}

type entityOutput struct {
	Body map[string]interface{} `json:"body"`
}

type listOutput struct {
	Body []map[string]interface{} `json:"body"`
}

type deleteOutput struct {
	Status int `status:"204"`
}

// NewHandler returns a REST handler with CRUD endpoints for all EnzymeML v2
// tables/models under the /v2 prefix, including Huma-generated docs/spec.
func NewHandler(manager *database.DBManager) (http.Handler, error) {
	mux := http.NewServeMux()
	if err := RegisterRoutes(mux, manager, "/v2"); err != nil {
		return nil, err
	}
	return mux, nil
}

// RegisterRoutes registers CRUD endpoints for all EnzymeML v2 models.
func RegisterRoutes(mux *http.ServeMux, manager *database.DBManager, basePath string) error {
	if mux == nil {
		return fmt.Errorf("mux cannot be nil")
	}
	if manager == nil {
		return fmt.Errorf("manager cannot be nil")
	}

	base := normalizeResourcePath(basePath)
	cfg := huma.DefaultConfig("EnzymeML v2 API", "v2")
	cfg.OpenAPIPath = base + "/openapi"
	cfg.DocsPath = base + "/docs"
	cfg.SchemasPath = base + "/schemas"
	api := humago.New(mux, cfg)

	for _, route := range v2ModelRoutes(base) {
		handler, err := newCRUDHandler(manager.DB(), route.Path, route.Model, route.Preloads...)
		if err != nil {
			return fmt.Errorf("failed registering %s: %w", route.Path, err)
		}
		registerCRUD(api, handler)
	}

	return nil
}

func registerCRUD(api huma.API, h *crudHandler) {
	name := h.modelType.Name()
	tag := strings.ToLower(name)
	modelSchema := h.modelSchema(api)
	entitySchema := h.entitySchema(api)
	entityListSchema := &huma.Schema{
		Type:  "array",
		Items: entitySchema,
	}

	huma.Register(api, huma.Operation{
		OperationID: "create" + name,
		Method:      http.MethodPost,
		Path:        h.resourcePath,
		Summary:     "Create " + name,
		Tags:        []string{tag},
		SkipValidateBody: true,
		RequestBody: &huma.RequestBody{
			Required: true,
			Content: map[string]*huma.MediaType{
				"application/json": {
					Schema: modelSchema,
				},
			},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusCreated): {
				Description: http.StatusText(http.StatusCreated),
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: entitySchema,
					},
				},
			},
		},
	}, func(ctx context.Context, input *createInput) (*createdOutput, error) {
		entity := h.newEntity()
		if err := decodeEntityMap(input.Body, entity); err != nil {
			return nil, huma.Error400BadRequest("invalid JSON payload", err)
		}

		if err := h.db.Session(&gorm.Session{FullSaveAssociations: true}).Create(entity).Error; err != nil {
			return nil, huma.Error500InternalServerError("failed to create resource", err)
		}

		return &createdOutput{Status: http.StatusCreated, Body: h.wrapEntity(entity)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "list" + name,
		Method:      http.MethodGet,
		Path:        h.resourcePath,
		Summary:     "List " + name,
		Tags:        []string{tag},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: http.StatusText(http.StatusOK),
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: entityListSchema,
					},
				},
			},
		},
	}, func(ctx context.Context, input *struct{}) (*listOutput, error) {
		slicePtr := h.newSlice()
		if err := h.withPreloads(h.db).Find(slicePtr).Error; err != nil {
			return nil, huma.Error500InternalServerError("failed to list resources", err)
		}

		slice := reflect.ValueOf(slicePtr).Elem()
		result := make([]map[string]interface{}, 0, slice.Len())
		for i := 0; i < slice.Len(); i++ {
			item := slice.Index(i)
			itemPtr := reflect.New(h.modelType)
			itemPtr.Elem().Set(item)
			result = append(result, h.wrapEntity(itemPtr.Interface()))
		}
		return &listOutput{Body: result}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "get" + name,
		Method:      http.MethodGet,
		Path:        h.resourcePath + "/{id}",
		Summary:     "Get " + name + " by ID",
		Tags:        []string{tag},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: http.StatusText(http.StatusOK),
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: entitySchema,
					},
				},
			},
		},
	}, func(ctx context.Context, input *idInput) (*entityOutput, error) {
		id, err := h.parseID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id", err)
		}

		entity := h.newEntity()
		if err := h.withPreloads(h.byID(h.db, id)).First(entity).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, huma.Error404NotFound("resource not found")
			}
			return nil, huma.Error500InternalServerError("failed to fetch resource", err)
		}

		return &entityOutput{Body: h.wrapEntity(entity)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "update" + name,
		Method:      http.MethodPut,
		Path:        h.resourcePath + "/{id}",
		Summary:     "Update " + name,
		Tags:        []string{tag},
		SkipValidateBody: true,
		RequestBody: &huma.RequestBody{
			Required: true,
			Content: map[string]*huma.MediaType{
				"application/json": {
					Schema: modelSchema,
				},
			},
		},
		Responses: map[string]*huma.Response{
			strconv.Itoa(http.StatusOK): {
				Description: http.StatusText(http.StatusOK),
				Content: map[string]*huma.MediaType{
					"application/json": {
						Schema: entitySchema,
					},
				},
			},
		},
	}, func(ctx context.Context, input *updateInput) (*entityOutput, error) {
		id, err := h.parseID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id", err)
		}

		current := h.newEntity()
		if err := h.byID(h.db, id).First(current).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, huma.Error404NotFound("resource not found")
			}
			return nil, huma.Error500InternalServerError("failed to fetch resource", err)
		}

		updated := h.newEntity()
		if err := decodeEntityMap(input.Body, updated); err != nil {
			return nil, huma.Error400BadRequest("invalid JSON payload", err)
		}
		h.copyPrimaryKey(current, updated)

		if err := h.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(updated).Error; err != nil {
			return nil, huma.Error500InternalServerError("failed to update resource", err)
		}

		reloaded := h.newEntity()
		if err := h.withPreloads(h.byID(h.db, id)).First(reloaded).Error; err != nil {
			return nil, huma.Error500InternalServerError("failed to reload resource", err)
		}

		return &entityOutput{Body: h.wrapEntity(reloaded)}, nil
	})

	huma.Register(api, huma.Operation{
		OperationID: "delete" + name,
		Method:      http.MethodDelete,
		Path:        h.resourcePath + "/{id}",
		Summary:     "Delete " + name,
		Tags:        []string{tag},
	}, func(ctx context.Context, input *idInput) (*deleteOutput, error) {
		id, err := h.parseID(input.ID)
		if err != nil {
			return nil, huma.Error400BadRequest("invalid id", err)
		}

		result := h.byID(h.db, id).Delete(h.newEntity())
		if result.Error != nil {
			return nil, huma.Error500InternalServerError("failed to delete resource", result.Error)
		}
		if result.RowsAffected == 0 {
			return nil, huma.Error404NotFound("resource not found")
		}

		return &deleteOutput{Status: http.StatusNoContent}, nil
	})
}

func v2ModelRoutes(base string) []modelRoute {
	return []modelRoute{
		{
			Path:  base + "/documents",
			Model: &enzymeml_v2.EnzymeMLDocument{},
			Preloads: []string{
				"Creators",
				"Vessels.Unit.BaseUnits",
				"Proteins",
				"Complexes",
				"SmallMolecules",
				"Reactions.KineticLaw.Variables",
				"Reactions.Reactants",
				"Reactions.Products",
				"Reactions.Modifiers",
				"Measurements.SpeciesData.DataUnit.BaseUnits",
				"Measurements.SpeciesData.TimeUnit.BaseUnits",
				"Measurements.TemperatureUnit.BaseUnits",
				"Equations.Variables",
				"Parameters.Unit.BaseUnits",
			},
		},
		{
			Path:  base + "/creators",
			Model: &enzymeml_v2.Creator{},
		},
		{
			Path:     base + "/vessels",
			Model:    &enzymeml_v2.Vessel{},
			Preloads: []string{"Unit.BaseUnits"},
		},
		{
			Path:  base + "/proteins",
			Model: &enzymeml_v2.Protein{},
		},
		{
			Path:  base + "/complexes",
			Model: &enzymeml_v2.Complex{},
		},
		{
			Path:  base + "/small-molecules",
			Model: &enzymeml_v2.SmallMolecule{},
		},
		{
			Path:     base + "/reactions",
			Model:    &enzymeml_v2.Reaction{},
			Preloads: []string{"KineticLaw.Variables", "Reactants", "Products", "Modifiers"},
		},
		{
			Path:  base + "/reaction-elements",
			Model: &enzymeml_v2.ReactionElement{},
		},
		{
			Path:  base + "/modifier-elements",
			Model: &enzymeml_v2.ModifierElement{},
		},
		{
			Path:     base + "/equations",
			Model:    &enzymeml_v2.Equation{},
			Preloads: []string{"Variables"},
		},
		{
			Path:  base + "/variables",
			Model: &enzymeml_v2.Variable{},
		},
		{
			Path:     base + "/parameters",
			Model:    &enzymeml_v2.Parameter{},
			Preloads: []string{"Unit.BaseUnits"},
		},
		{
			Path:  base + "/measurements",
			Model: &enzymeml_v2.Measurement{},
			Preloads: []string{
				"SpeciesData.DataUnit.BaseUnits",
				"SpeciesData.TimeUnit.BaseUnits",
				"TemperatureUnit.BaseUnits",
			},
		},
		{
			Path:     base + "/measurement-data",
			Model:    &enzymeml_v2.MeasurementData{},
			Preloads: []string{"DataUnit.BaseUnits", "TimeUnit.BaseUnits"},
		},
		{
			Path:     base + "/unit-definitions",
			Model:    &enzymeml_v2.UnitDefinition{},
			Preloads: []string{"BaseUnits"},
		},
		{
			Path:  base + "/base-units",
			Model: &enzymeml_v2.BaseUnit{},
		},
	}
}

func newCRUDHandler(db *gorm.DB, resource string, model interface{}, preloads ...string) (*crudHandler, error) {
	if db == nil {
		return nil, fmt.Errorf("db cannot be nil")
	}

	modelType := reflect.TypeOf(model)
	if modelType == nil {
		return nil, fmt.Errorf("model cannot be nil")
	}
	if modelType.Kind() == reflect.Pointer {
		modelType = modelType.Elem()
	}
	if modelType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("model must be a struct or pointer to struct")
	}

	pkIdx, pkType, err := findPrimaryKey(modelType)
	if err != nil {
		return nil, err
	}

	return &crudHandler{
		db:            db,
		resourcePath:  normalizeResourcePath(resource),
		modelType:     modelType,
		primaryKeyIdx: pkIdx,
		primaryKeyTyp: pkType,
		preloads:      preloads,
	}, nil
}

func (h *crudHandler) withPreloads(db *gorm.DB) *gorm.DB {
	for _, preload := range h.preloads {
		db = db.Preload(preload)
	}
	return db
}

func (h *crudHandler) byID(db *gorm.DB, id interface{}) *gorm.DB {
	fieldName := h.modelType.Field(h.primaryKeyIdx).Name
	return db.Where(map[string]interface{}{fieldName: id})
}

func (h *crudHandler) newEntity() interface{} {
	return reflect.New(h.modelType).Interface()
}

func (h *crudHandler) newSlice() interface{} {
	return reflect.New(reflect.SliceOf(h.modelType)).Interface()
}

func (h *crudHandler) wrapEntity(entity interface{}) map[string]interface{} {
	value := reflect.ValueOf(entity)
	if value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	return map[string]interface{}{
		"id":   value.Field(h.primaryKeyIdx).Interface(),
		"data": value.Interface(),
	}
}

func (h *crudHandler) copyPrimaryKey(source interface{}, target interface{}) {
	src := reflect.ValueOf(source)
	if src.Kind() == reflect.Pointer {
		src = src.Elem()
	}
	dst := reflect.ValueOf(target)
	if dst.Kind() == reflect.Pointer {
		dst = dst.Elem()
	}
	dst.Field(h.primaryKeyIdx).Set(src.Field(h.primaryKeyIdx))
}

func (h *crudHandler) parseID(raw string) (interface{}, error) {
	if strings.TrimSpace(raw) == "" {
		return nil, fmt.Errorf("invalid id")
	}

	switch h.primaryKeyTyp.Kind() {
	case reflect.String:
		return raw, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		parsed, err := strconv.ParseInt(raw, 10, 64)
		if err != nil {
			return nil, err
		}
		id := reflect.New(h.primaryKeyTyp).Elem()
		id.SetInt(parsed)
		return id.Interface(), nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		parsed, err := strconv.ParseUint(raw, 10, 64)
		if err != nil {
			return nil, err
		}
		id := reflect.New(h.primaryKeyTyp).Elem()
		id.SetUint(parsed)
		return id.Interface(), nil
	default:
		return raw, nil
	}
}

func (h *crudHandler) modelSchema(api huma.API) *huma.Schema {
	return api.OpenAPI().Components.Schemas.Schema(h.modelType, true, h.modelType.Name())
}

func (h *crudHandler) entitySchema(api huma.API) *huma.Schema {
	return &huma.Schema{
		Type: "object",
		Properties: map[string]*huma.Schema{
			"id": api.OpenAPI().Components.Schemas.Schema(h.primaryKeyTyp, false, h.modelType.Name()+"ID"),
			"data": h.modelSchema(api),
		},
		Required: []string{"id", "data"},
	}
}

func findPrimaryKey(modelType reflect.Type) (int, reflect.Type, error) {
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("gorm")
		if strings.Contains(tag, "primaryKey") {
			return i, field.Type, nil
		}
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Name == "ID" || field.Name == "Id" {
			return i, field.Type, nil
		}
	}

	return -1, nil, fmt.Errorf("no primary key found on model %s", modelType.Name())
}

func normalizeResourcePath(resource string) string {
	resource = strings.TrimSpace(resource)
	if resource == "" {
		return "/"
	}
	if !strings.HasPrefix(resource, "/") {
		resource = "/" + resource
	}
	if len(resource) > 1 {
		resource = strings.TrimRight(resource, "/")
	}
	return resource
}

func decodeEntityMap(raw map[string]interface{}, out interface{}) error {
	b, err := json.Marshal(raw)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, out)
}
