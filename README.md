<div align="center">

# 🧬 EnzymeML Go

**A comprehensive Go implementation of the EnzymeML data exchange format for enzymatic data**

[![Go Version](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)

---

</div>

## 📋 Table of Contents

- [🧬 EnzymeML Go](#-enzymeml-go)
  - [📋 Table of Contents](#-table-of-contents)
  - [🌟 Overview](#-overview)
    - [Why EnzymeML Go?](#why-enzymeml-go)
  - [✨ Features](#-features)
  - [🚀 Installation](#-installation)
    - [Requirements](#requirements)
  - [📖 Usage](#-usage)
    - [Quick Start](#quick-start)
    - [🌐 Start v2 REST API](#-start-v2-rest-api)
    - [🗄️ Database Example](#️-database-example)
  - [🤝 Contributing](#-contributing)
  - [📄 License](#-license)

## 🌟 Overview

Welcome to **EnzymeML Go** - a robust Go interface for working with EnzymeML documents. EnzymeML is a comprehensive data exchange format specifically designed for enzymatic data, enabling researchers to document, share, and analyze enzymatic experiments in a standardized way.

### Why EnzymeML Go?

🎯 **Perfect for Database & Web Infrastructure** - The key strength of EnzymeML Go lies in powering databases and web applications with enzymatic data while maintaining full compliance with the EnzymeML specification.

## ✨ Features

EnzymeML Go implements the complete **EnzymeML v2 specification**, providing you with:

- 📊 **Experimental Documentation**: Standardized format for enzymatic experiments
- ⚗️ **Reaction Conditions**: Capture detailed reaction conditions and measurements
- 📈 **Kinetic Models**: Define and work with kinetic models and parameters
- 🔄 **Data Exchange**: Seamless data exchange between platforms and modeling tools
- 🌐 **FAIR Compliance**: Ensures Findable, Accessible, Interoperable, Reusable data principles
- 🗄️ **Database Ready**: Optimized for database storage and web applications

## 🚀 Installation

Get started with EnzymeML Go in seconds:

```bash
go get github.com/EnzymeML/enzymeml-go
```

### Requirements

- **Go 1.25+** - Ensure you have a recent version of Go installed
- **Modern Go modules** - This package uses Go modules for dependency management

## 📖 Usage

### Quick Start

```go
package main

import (
    "fmt"
    enzymeml_v2 "github.com/EnzymeML/enzymeml-go/src"
)

func main() {
    doc := enzymeml_v2.EnzymeMLDocument{Name: "Example"}
    fmt.Println(doc.Name)
}
```

### 🌐 Start v2 REST API

Create a `main.go`:

```go
package main

import (
 "log"
 "net/http"

 apiv2 "github.com/EnzymeML/enzymeml-go/src/api/v2"
 enzymeml_v2 "github.com/EnzymeML/enzymeml-go/src"
 database "github.com/EnzymeML/enzymeml-go/src/database/v2"
)

func main() {
 models := []interface{}{
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

 dbManager, err := database.NewDBManager("enzymeml.db", models)
 if err != nil {
  log.Fatalf("failed to init DB: %v", err)
 }
 defer dbManager.Close()

 handler, err := apiv2.NewHandler(dbManager)
 if err != nil {
  log.Fatalf("failed to build API handler: %v", err)
 }

 log.Println("v2 API listening on :8080")
 log.Fatal(http.ListenAndServe(":8080", handler))
}
```

Run:

```bash
go run .
```

Example CRUD for documents:

```bash
curl -X POST http://localhost:8080/v2/documents \
  -H "Content-Type: application/json" \
  -d '{"name":"my-doc","version":"2.0.0"}'

curl http://localhost:8080/v2/documents
curl http://localhost:8080/v2/documents/1

curl -X PUT http://localhost:8080/v2/documents/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"my-doc-updated","version":"2.0.1"}'

curl -X DELETE http://localhost:8080/v2/documents/1
```

The API exposes CRUD for all v2 resources under `/v2`:
`documents`, `creators`, `vessels`, `proteins`, `complexes`, `small-molecules`, `reactions`, `reaction-elements`, `modifier-elements`, `equations`, `variables`, `parameters`, `measurements`, `measurement-data`, `unit-definitions`, `base-units`.

With Huma, docs/spec are available automatically:

- `http://localhost:8080/v2/docs`
- `http://localhost:8080/v2/openapi.json`
- `http://localhost:8080/v2/openapi.yaml`

### 🗄️ Database Example

For a comprehensive example of how to use EnzymeML Go to create an enzymatic data database and power web applications, check out our detailed example:

**👉 [Database Example](./examples/database_example/)**

This example demonstrates:

- Setting up an EnzymeML-compliant database
- Storing and retrieving enzymatic data
- Web application integration
- Best practices for data management

## 🤝 Contributing

We welcome contributions from the community! Please feel free to submit a pull request.

## 📄 License

This project is licensed under the **MIT License** - see the [LICENSE](LICENSE) file for details.

---

<div align="center">
<strong>Made with ❤️ by the EnzymeML Team</strong>
</div>
