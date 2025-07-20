# RULES-GO.md - Go Coding Standards and Best Practices

> **Purpose**: This document defines coding standards for Go projects to ensure code quality, maintainability, and consistency. These rules help code pass linting tools gracefully while following industry best practices.

## Table of Contents

1. [Core Principles](#core-principles)
   - [KISS - Keep It Simple](#kiss---keep-it-simple)
   - [DRY - Don't Repeat Yourself](#dry---dont-repeat-yourself)
   - [SOLID Principles](#solid-principles)
2. [Go Language Conventions](#go-language-conventions)
   - [Naming Conventions](#naming-conventions)
   - [Package Design](#package-design)
   - [Error Handling](#error-handling)
   - [Interface Design](#interface-design)
   - [Concurrency](#concurrency)
3. [Code Quality & Linting](#code-quality--linting)
   - [Linting Rules](#linting-rules)
   - [Performance Guidelines](#performance-guidelines)
   - [Security Best Practices](#security-best-practices)
4. [Code Organization](#code-organization)
   - [Project Structure](#project-structure)
   - [Testing Standards](#testing-standards)
   - [Documentation](#documentation)
5. [Common Patterns & Anti-patterns](#common-patterns--anti-patterns)

---

## Core Principles

### KISS - Keep It Simple

Go's philosophy emphasizes simplicity. Write code that is easy to read, understand, and maintain.

#### Rules:

**1. Prefer simple solutions over clever ones**
```go
// ❌ Bad - Clever but hard to understand
func isEven(n int) bool {
    return n&1 == 0
}

// ✅ Good - Simple and clear
func isEven(n int) bool {
    return n%2 == 0
}
```

**2. Avoid premature optimization**
```go
// ❌ Bad - Premature optimization
var bufferPool = sync.Pool{
    New: func() interface{} {
        return make([]byte, 1024)
    },
}

// ✅ Good - Start simple, optimize when proven necessary
func processData(data []byte) {
    buf := make([]byte, 1024)
    // Process...
}
```

**3. Use standard library when possible**
```go
// ❌ Bad - Reinventing the wheel
func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}

// ✅ Good - Use standard library (Go 1.21+)
import "slices"
found := slices.Contains(slice, item)
```

### DRY - Don't Repeat Yourself

Balance DRY with readability. In Go, "a little copying is better than a little dependency."

#### Rules:

**1. Extract common logic, but don't over-abstract**
```go
// ❌ Bad - Too much abstraction for simple validation
type Validator interface {
    Validate() error
}

type StringValidator struct {
    value string
    rules []func(string) error
}

// ✅ Good - Simple and direct
func validateEmail(email string) error {
    if email == "" {
        return errors.New("email is required")
    }
    if !strings.Contains(email, "@") {
        return errors.New("email is invalid")
    }
    return nil
}
```

**2. Use composition over inheritance**
```go
// ✅ Good - Composition
type Logger struct {
    output io.Writer
}

type Service struct {
    logger *Logger
    db     *sql.DB
}

func (s *Service) Process() error {
    s.logger.Write([]byte("Processing..."))
    // ...
}
```

### SOLID Principles

Adapt SOLID principles to Go's design philosophy.

#### Single Responsibility Principle
Each package, type, and function should have one reason to change.

```go
// ❌ Bad - Mixed responsibilities
type User struct {
    ID    int
    Name  string
    Email string
}

func (u *User) Save(db *sql.DB) error { /*...*/ }
func (u *User) SendEmail(msg string) error { /*...*/ }

// ✅ Good - Separated concerns
type User struct {
    ID    int
    Name  string
    Email string
}

type UserRepository struct {
    db *sql.DB
}

func (r *UserRepository) Save(u *User) error { /*...*/ }

type EmailService struct {
    smtp *SMTPClient
}

func (s *EmailService) Send(to, msg string) error { /*...*/ }
```

#### Open/Closed Principle
Design should be open for extension but closed for modification.

```go
// ✅ Good - Extensible through interfaces
type Processor interface {
    Process(data []byte) ([]byte, error)
}

type Pipeline struct {
    processors []Processor
}

func (p *Pipeline) Execute(data []byte) ([]byte, error) {
    var err error
    for _, proc := range p.processors {
        data, err = proc.Process(data)
        if err != nil {
            return nil, err
        }
    }
    return data, nil
}
```

#### Liskov Substitution Principle
Interfaces should be implemented correctly.

```go
// ✅ Good - Interface implementation
type Writer interface {
    Write([]byte) (int, error)
}

// Any type implementing Writer should behave as expected
type FileWriter struct {
    file *os.File
}

func (w *FileWriter) Write(data []byte) (int, error) {
    return w.file.Write(data) // Correct implementation
}
```

#### Interface Segregation Principle
Keep interfaces small and focused.

```go
// ❌ Bad - Fat interface
type DataStore interface {
    Connect() error
    Disconnect() error
    Query(string) ([]Row, error)
    Insert(Row) error
    Update(Row) error
    Delete(int) error
    BeginTransaction() error
    Commit() error
    Rollback() error
}

// ✅ Good - Segregated interfaces
type Reader interface {
    Read(id int) (*Record, error)
}

type Writer interface {
    Write(record *Record) error
}

type Deleter interface {
    Delete(id int) error
}

type ReadWriter interface {
    Reader
    Writer
}
```

#### Dependency Inversion Principle
Depend on interfaces, not concrete types.

```go
// ❌ Bad - Depends on concrete type
type Service struct {
    db *PostgresDB
}

// ✅ Good - Depends on interface
type Service struct {
    db Database
}

type Database interface {
    Query(string) ([]Row, error)
}
```

---

## Go Language Conventions

### Naming Conventions

Follow Go's naming conventions for consistency and readability.

#### Package Names
- Lowercase, single-word names
- No underscores or mixedCaps
- Short but descriptive

```go
// ❌ Bad
package user_service
package userService
package usersvc

// ✅ Good
package user
package auth
package httputil
```

#### Variable and Function Names
- Use camelCase for unexported
- Use PascalCase for exported
- Be descriptive but concise

```go
// ❌ Bad
var u User
var userobj User
var user_name string
func get_user() {}

// ✅ Good
var currentUser User
var userName string
func GetUser() User {}
func (s *Service) processRequest() {}
```

#### Interface Names
- Use -er suffix when possible
- Keep them small and focused

```go
// ✅ Good
type Reader interface {
    Read([]byte) (int, error)
}

type Stringer interface {
    String() string
}

type Repository interface {
    Find(id int) (*User, error)
    Save(*User) error
}
```

#### Constants
- Use PascalCase for exported constants
- Use camelCase for unexported constants
- Group related constants

```go
// ✅ Good
const (
    MaxRetries = 3
    MinTimeout = 1 * time.Second
)

const (
    statusPending = iota
    statusActive
    statusInactive
)
```

### Package Design

#### Keep packages focused
Each package should have a single, well-defined purpose.

```go
// ✅ Good package structure
myapp/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── auth/
│   ├── user/
│   └── database/
├── pkg/
│   ├── logger/
│   └── validator/
```

#### Avoid circular dependencies
Structure packages to form a directed acyclic graph.

```go
// ❌ Bad - Circular dependency
// package user imports package order
// package order imports package user

// ✅ Good - Clear dependency direction
// package models (defines User, Order types)
// package user imports models
// package order imports models
```

### Error Handling

Go's explicit error handling is a feature, not a bug. Handle errors gracefully.

#### Always check errors
```go
// ❌ Bad - Ignoring errors
result, _ := someFunction()

// ✅ Good - Handle errors explicitly
result, err := someFunction()
if err != nil {
    return fmt.Errorf("failed to do something: %w", err)
}
```

#### Wrap errors with context
```go
// ❌ Bad - No context
if err != nil {
    return err
}

// ✅ Good - Add context
if err != nil {
    return fmt.Errorf("failed to process user %d: %w", userID, err)
}
```

#### Create custom error types when needed
```go
// ✅ Good - Custom error type
type ValidationError struct {
    Field string
    Value interface{}
    Msg   string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for %s: %s", e.Field, e.Msg)
}
```

#### Use error variables for sentinel errors
```go
// ✅ Good - Sentinel errors
var (
    ErrNotFound   = errors.New("not found")
    ErrInvalidID  = errors.New("invalid ID")
    ErrPermission = errors.New("permission denied")
)

// Usage
if errors.Is(err, ErrNotFound) {
    // Handle not found case
}
```

### Interface Design

#### Accept interfaces, return structs
```go
// ✅ Good
func ProcessData(r io.Reader) (*Result, error) {
    // Process data from any reader
    return &Result{}, nil
}
```

#### Keep interfaces small
```go
// ❌ Bad - Too many methods
type Cache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
    Delete(key string) error
    Clear() error
    Size() int
    Keys() []string
    Values() []interface{}
}

// ✅ Good - Minimal interface
type Cache interface {
    Get(key string) (interface{}, error)
    Set(key string, value interface{}) error
}
```

### Concurrency

Use Go's concurrency primitives correctly and safely.

#### Don't communicate by sharing memory; share memory by communicating
```go
// ❌ Bad - Shared memory
var counter int
var mu sync.Mutex

func increment() {
    mu.Lock()
    counter++
    mu.Unlock()
}

// ✅ Good - Channel communication
type Counter struct {
    ch chan int
}

func (c *Counter) Increment() {
    c.ch <- 1
}
```

#### Always close channels from the sender
```go
// ✅ Good
func producer(ch chan<- int) {
    defer close(ch)
    for i := 0; i < 10; i++ {
        ch <- i
    }
}
```

#### Use context for cancellation
```go
// ✅ Good
func worker(ctx context.Context) error {
    for {
        select {
        case <-ctx.Done():
            return ctx.Err()
        default:
            // Do work
        }
    }
}
```

---

## Code Quality & Linting

### Linting Rules

Configure golangci-lint to catch common issues. Here's a recommended configuration:

```yaml
# .golangci.yml
run:
  timeout: 5m

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - errcheck
    - staticcheck
    - gosimple
    - ineffassign
    - unused
    - misspell
    - gosec
    - revive

linters-settings:
  gofmt:
    simplify: true
  goimports:
    local-prefixes: github.com/myorg/myproject
  errcheck:
    check-type-assertions: true
    check-blank: false
  revive:
    rules:
      - name: var-naming
      - name: package-comments
      - name: exported
```

#### Common linting fixes:

**1. Always format your code**
```bash
# Format all Go files
gofmt -w .
# Or use goimports (also adds missing imports)
goimports -w .
```

**2. Check error returns**
```go
// ❌ Bad - errcheck warning
fmt.Println("Hello")
json.Marshal(data)

// ✅ Good
fmt.Println("Hello") // fmt.Println errors are typically ignored
if _, err := json.Marshal(data); err != nil {
    return err
}
```

**3. Remove unused code**
```go
// ❌ Bad - unused variable
func process() {
    unused := 42 // This will trigger a linting error
    fmt.Println("Processing")
}

// ✅ Good
func process() {
    fmt.Println("Processing")
}
```

**4. Fix inefficient assignments**
```go
// ❌ Bad - ineffassign warning
func getValue() int {
    result := 0
    result = computeValue() // First assignment was pointless
    return result
}

// ✅ Good
func getValue() int {
    return computeValue()
}
```

### Performance Guidelines

#### Preallocate slices when size is known
```go
// ❌ Bad
var results []string
for _, item := range items {
    results = append(results, transform(item))
}

// ✅ Good
results := make([]string, 0, len(items))
for _, item := range items {
    results = append(results, transform(item))
}
```

#### Use string builder for concatenation
```go
// ❌ Bad - Multiple allocations
func concat(parts []string) string {
    result := ""
    for _, part := range parts {
        result += part
    }
    return result
}

// ✅ Good - Single allocation
func concat(parts []string) string {
    var sb strings.Builder
    for _, part := range parts {
        sb.WriteString(part)
    }
    return sb.String()
}
```

#### Avoid unnecessary conversions
```go
// ❌ Bad
data := []byte("hello")
s := string(data)
bytes := []byte(s) // Unnecessary conversion

// ✅ Good
data := []byte("hello")
// Use data directly
```

### Security Best Practices

#### Validate all inputs
```go
// ✅ Good
func processUserInput(input string) error {
    input = strings.TrimSpace(input)
    if len(input) == 0 || len(input) > 100 {
        return errors.New("invalid input length")
    }
    if !isValidInput(input) {
        return errors.New("invalid characters in input")
    }
    // Process validated input
    return nil
}
```

#### Use crypto/rand for security-sensitive randomness
```go
// ❌ Bad - Predictable
import "math/rand"
token := rand.Intn(1000000)

// ✅ Good - Cryptographically secure
import "crypto/rand"
b := make([]byte, 32)
if _, err := rand.Read(b); err != nil {
    return err
}
token := base64.URLEncoding.EncodeToString(b)
```

#### Never log sensitive data
```go
// ❌ Bad
log.Printf("User login: username=%s, password=%s", username, password)

// ✅ Good
log.Printf("User login: username=%s", username)
```

---

## Code Organization

### Project Structure

Follow the standard Go project layout:

```
myproject/
├── cmd/                    # Main applications
│   ├── api/
│   │   └── main.go
│   └── worker/
│       └── main.go
├── internal/              # Private application code
│   ├── config/
│   ├── handler/
│   ├── model/
│   └── service/
├── pkg/                   # Public libraries
│   ├── client/
│   └── utils/
├── api/                   # API definitions (OpenAPI, Proto)
├── configs/               # Configuration files
├── scripts/               # Build/install scripts
├── deployments/           # Deployment configurations
├── test/                  # Integration tests
├── docs/                  # Documentation
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

#### Key principles:
- `cmd/` contains main applications
- `internal/` prevents external imports
- `pkg/` contains exportable libraries
- Keep `main.go` minimal

### Testing Standards

#### Write table-driven tests
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {
            name:    "valid email",
            email:   "user@example.com",
            wantErr: false,
        },
        {
            name:    "empty email",
            email:   "",
            wantErr: true,
        },
        {
            name:    "missing @",
            email:   "userexample.com",
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

#### Use test helpers
```go
// testhelpers.go
func setupTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite3", ":memory:")
    if err != nil {
        t.Fatalf("Failed to create test DB: %v", err)
    }
    t.Cleanup(func() {
        db.Close()
    })
    return db
}
```

#### Test file naming
- Unit tests: `*_test.go` in same package
- Integration tests: `*_integration_test.go` with build tag
- Test data: `testdata/` directory

### Documentation

#### Package comments
```go
// Package user provides user management functionality including
// authentication, authorization, and profile management.
//
// Basic usage:
//
//	service := user.NewService(db)
//	u, err := service.GetByID(ctx, userID)
//	if err != nil {
//	    return err
//	}
package user
```

#### Function comments
```go
// ValidateEmail checks if the provided email address is valid.
// It returns an error if the email is empty or doesn't contain an @ symbol.
// This is a basic validation and doesn't guarantee email deliverability.
func ValidateEmail(email string) error {
    // Implementation
}
```

#### Example tests
```go
func ExampleValidateEmail() {
    err := ValidateEmail("user@example.com")
    if err != nil {
        fmt.Println("Invalid email:", err)
    } else {
        fmt.Println("Email is valid")
    }
    // Output: Email is valid
}
```

---

## Common Patterns & Anti-patterns

### Patterns to Follow

#### Options pattern for constructors
```go
type Option func(*Server)

func WithPort(port int) Option {
    return func(s *Server) {
        s.port = port
    }
}

func WithTimeout(timeout time.Duration) Option {
    return func(s *Server) {
        s.timeout = timeout
    }
}

func NewServer(opts ...Option) *Server {
    s := &Server{
        port:    8080,
        timeout: 30 * time.Second,
    }
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Usage
server := NewServer(
    WithPort(9090),
    WithTimeout(60 * time.Second),
)
```

#### Result pattern for multiple returns
```go
type Result[T any] struct {
    Value T
    Error error
}

func FetchUser(id int) Result[*User] {
    user, err := db.GetUser(id)
    return Result[*User]{Value: user, Error: err}
}
```

### Anti-patterns to Avoid

#### Don't panic in libraries
```go
// ❌ Bad - Library panics
func GetConfig(key string) string {
    value, ok := config[key]
    if !ok {
        panic("config key not found: " + key)
    }
    return value
}

// ✅ Good - Return error
func GetConfig(key string) (string, error) {
    value, ok := config[key]
    if !ok {
        return "", fmt.Errorf("config key not found: %s", key)
    }
    return value, nil
}
```

#### Avoid init() when possible
```go
// ❌ Bad - Hidden initialization
func init() {
    database = connectDB()
}

// ✅ Good - Explicit initialization
func Initialize() error {
    db, err := connectDB()
    if err != nil {
        return err
    }
    database = db
    return nil
}
```

#### Don't use empty interfaces unnecessarily
```go
// ❌ Bad - Loss of type safety
func Process(data interface{}) interface{} {
    // Type assertions everywhere
}

// ✅ Good - Use generics or specific types
func Process[T any](data T) (T, error) {
    // Type safe processing
}
```

---

## Summary

These rules prioritize:
1. **Simplicity** - Clear, obvious code over clever solutions
2. **Consistency** - Following Go idioms and conventions
3. **Maintainability** - Code that's easy to understand and modify
4. **Testability** - Code that's easy to test
5. **Performance** - Efficient code without premature optimization

Remember: "Clear is better than clever" - The Go Proverb

---

*Last updated: January 2025*
*Version: 1.0.0*