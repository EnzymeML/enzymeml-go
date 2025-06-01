<div align="center">

# 🧬 EnzymeML Go

**A comprehensive Go implementation of the EnzymeML data exchange format for enzymatic data**

[![Go Version](https://img.shields.io/badge/Go-1.19+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
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
go get github.com/enzymeml-go
```

### Requirements

- **Go 1.19+** - Ensure you have a recent version of Go installed
- **Modern Go modules** - This package uses Go modules for dependency management

## 📖 Usage

### Quick Start

```go
package main

import (
    "fmt"
    "github.com/enzymeml-go"
)

func main() {
    // Your EnzymeML Go code here
    fmt.Println("Welcome to EnzymeML Go!")
}
```

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