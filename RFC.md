# RFC: From REST to GraphQL Federation - A Modern API Strategy

- [RFC: From REST to GraphQL Federation - A Modern API Strategy](#rfc-from-rest-to-graphql-federation---a-modern-api-strategy)
  - [1. Executive Summary](#1-executive-summary)
  - [2. GraphQL vs REST: A Comparative Analysis](#2-graphql-vs-rest-a-comparative-analysis)
    - [2.1 Data Fetching Patterns](#21-data-fetching-patterns)
    - [2.3 Development Experience](#23-development-experience)
  - [3. The Microservices Challenge](#3-the-microservices-challenge)
    - [3.1 REST API Gateway Pattern](#31-rest-api-gateway-pattern)
    - [3.2 GraphQL Federation Approach](#32-graphql-federation-approach)
  - [4. GraphQL Federation in Action](#4-graphql-federation-in-action)
    - [4.1 Federation Concepts \& Directives](#41-federation-concepts--directives)
    - [4.2 Schema Stitching Example](#42-schema-stitching-example)
      - [User Service Schema (users/schema.graphql)](#user-service-schema-usersschemagraphql)
      - [Order Service Schema (orders/schema.graphql)](#order-service-schema-ordersschemagraphql)
      - [Product Service Schema (products/schema.graphql)](#product-service-schema-productsschemagraphql)
    - [4.3 Cross-Service Query Resolution](#43-cross-service-query-resolution)
    - [4.4 Performance Comparison](#44-performance-comparison)
  - [5. Implementation with Go](#5-implementation-with-go)
    - [5.1 Technology Stack](#51-technology-stack)
      - [Why Go for GraphQL?](#why-go-for-graphql)
    - [5.2 Go GraphQL Tooling](#52-go-graphql-tooling)
      - [gqlgen Configuration](#gqlgen-configuration)
      - [Apollo Router Configuration (router.yaml)](#apollo-router-configuration-routeryaml)
      - [Go GraphQL Development Tools](#go-graphql-development-tools)
    - [5.3 Example: User Service Implementation](#53-example-user-service-implementation)
      - [Schema Definition](#schema-definition)
      - [Project Structure](#project-structure)
      - [Go Implementation](#go-implementation)
      - [DataLoader Implementation for N+1 Query Prevention](#dataloader-implementation-for-n1-query-prevention)
  - [6. Frontend Integration](#6-frontend-integration)
    - [6.1 Code Generation with urql for React](#61-code-generation-with-urql-for-react)
  - [7. Challenges and Considerations](#7-challenges-and-considerations)
    - [7.1 Performance Concerns](#71-performance-concerns)
    - [7.2 Security](#72-security)
    - [7.3 Learning Curve](#73-learning-curve)
  - [8. Conclusion](#8-conclusion)
    - [8.1 GraphQL Federation Benefits for Our New Project](#81-graphql-federation-benefits-for-our-new-project)
    - [8.2 Success Criteria](#82-success-criteria)
    - [8.3 Key Benefits Realized](#83-key-benefits-realized)


## 1. Executive Summary
This document outlines the architecture and implementation approach for our new Go-based project using GraphQL and GraphQL Federation. It details the technical advantages of GraphQL over REST, with special focus on how GraphQL Federation addresses microservices architecture challenges. By adopting GraphQL Federation from the project's inception, we establish a future-proof, type-safe API layer that scales with our business needs while leveraging Go's performance characteristics and maintaining clear service boundaries.

## 2. GraphQL vs REST: A Comparative Analysis

### 2.1 Data Fetching Patterns

**REST API Example**:
```
# Multiple endpoints for related data
GET /users/123
GET /users/123/orders
GET /orders/456/items

# Response from /users/123
{
  "id": 123,
  "name": "John Doe",
  "email": "john@example.com",
  "address": {"street": "...", "city": "..."},  // Maybe not needed
  "preferences": {...}  // Maybe not needed
}
```

**GraphQL Equivalent**:
```graphql
query GetUserWithOrders($userId: ID!) {
  user(id: $userId) {
    id
    name
    email
    orders {
      id
      total
      items {
        name
        price
      }
    }
  }
}
```

*Key Differences*:
- **Over-fetching**: REST returns complete resources, GraphQL fetches only requested fields
- **Under-fetching**: REST needs multiple requests, GraphQL gets everything in one
- **Versioning**: REST uses URL versions, GraphQL evolves the schema

### 2.3 Development Experience

**REST Challenges**:
- Frontend teams depend on backend changes
- Documentation often lags implementation
- Multiple round trips affect performance
- No type safety between frontend and backend

**GraphQL Advantages**:
- Frontend-driven development
- Self-documenting with type system
- Single request for all data needs
- Strong typing across the stack

## 3. The Microservices Challenge

### 3.1 REST API Gateway Pattern
```
┌─────────────┐     ┌─────────────┐     ┌─────────────────┐
│   Client    │────▶│ API Gateway │────▶│   Microservice  │
└─────────────┘     └─────┬───────┘     └─────────────────┘
                          │
                          ▼
                 ┌─────────────────┐     ┌─────────────────┐
                 │  BFF Service    │────▶│  Microservice   │
                 │  (per client)   │     └─────────────────┘
                 └─────────────────┘
```
*Issues*:
- Complex orchestration logic in Gateway/BFF
- Multiple hops increase latency
- Tight coupling between services
- Hard to maintain as system grows

### 3.2 GraphQL Federation Approach
```
┌─────────────┐     ┌───────────────────────────┐     ┌─────────────────┐
│   Client    │────▶│     GraphQL Gateway       │     │  User Service   │
└─────────────┘     │  (Schema Composition)     │◀───▶│  (gqlgen)       │
                    └───────────┬───────────────┘     └─────────────────┘
                                │
                    ┌───────────▼───────────────┐     ┌─────────────────┐
                    │     Query Planning        │     │  Order Service  │
                    │  & Execution Engine       │◀───▶│  (gqlgen)       │
                    └───────────────────────────┘     └─────────────────┘
```
*Advantages*:
- Single request for cross-service data
- Services own their domain
- No need for orchestration layer
- Better developer experience

## 4. GraphQL Federation in Action

GraphQL Federation is an architecture that allows you to build a unified GraphQL API from multiple underlying services. Instead of building a monolithic GraphQL server or manually stitching schemas together, federation provides a declarative composition model where each service defines its own schema. These schemas are then automatically combined into a single graph by the federation gateway.

Key characteristics of GraphQL Federation:

- **Distributed Schema**: Each service defines its portion of the graph
- **Type Extensions**: Services can extend types defined in other services
- **Entity Resolution**: Services can reference and resolve entities owned by other services
- **Composition**: The gateway composes individual schemas into a unified graph
- **Transparent Execution**: Clients query the gateway as if it were a single GraphQL service


### 4.1 Federation Concepts & Directives

GraphQL Federation uses special directives to compose schemas across services:

- **@key**: Identifies entities that can be referenced across services
- **@external**: Marks fields from another service that are referenced locally
- **@requires**: Indicates fields needed from another service to resolve a field
- **@provides**: Indicates that a field normally resolved by another service can be resolved locally
- **@shareable**: Indicates a field that multiple services can resolve
- **@override**: Specifies that this service's field definition takes precedence over another's
- **@inaccessible**: Marks fields that should not be accessible in the gateway schema

### 4.2 Schema Stitching Example

#### User Service Schema (users/schema.graphql)
```graphql
extend schema @link(url: "https://specs.apollo.dev/federation/v2.0", import: ["@key", "@shareable"])

type User @key(fields: "id") {
  id: ID!
  name: String!
  email: String! @shareable
  createdAt: Time!
  profile: Profile
}

type Profile {
  bio: String
  avatarUrl: String
}

type Query {
  user(id: ID!): User
  users: [User!]!
}

scalar Time
```

#### Order Service Schema (orders/schema.graphql)
```graphql
extend schema @link(url: "https://specs.apollo.dev/federation/v2.0", 
                 import: ["@key", "@external", "@provides", "@requires"])

type Order @key(fields: "id") {
  id: ID!
  orderNumber: String!
  total: Float!
  status: OrderStatus!
  createdAt: Time!
  items: [OrderItem!]!
  user: User! @provides(fields: "name")
}

type OrderItem {
  product: Product!
  quantity: Int!
  unitPrice: Float!
}

type Product @key(fields: "id") {
  id: ID!
  name: String!
  price: Float!
}

enum OrderStatus {
  PENDING
  PROCESSING
  SHIPPED
  DELIVERED
  CANCELED
}

extend type User @key(fields: "id") {
  id: ID! @external
  email: String! @external
  name: String @external
  orders: [Order!]! @requires(fields: "id email")
}

type Query {
  order(id: ID!): Order
  ordersByUser(userId: ID!): [Order!]!
}

scalar Time
```

#### Product Service Schema (products/schema.graphql)
```graphql
extend schema @link(url: "https://specs.apollo.dev/federation/v2.0", import: ["@key"])

type Product @key(fields: "id") {
  id: ID!
  name: String!
  description: String
  price: Float!
  inventory: Int!
  category: Category
  images: [Image!]!
}

type Category {
  id: ID!
  name: String!
  description: String
}

type Image {
  url: String!
  alt: String
  isPrimary: Boolean!
}

type Query {
  product(id: ID!): Product
  products(limit: Int = 10, offset: Int = 0): [Product!]!
  productsByCategory(categoryId: ID!): [Product!]!
}
```

### 4.3 Cross-Service Query Resolution

With Federation, clients can write queries that span multiple services:

```graphql
# This query touches all three services but is resolved as one
query GetUserWithOrdersAndProducts {
  user(id: "user-123") {
    name
    email
    profile {
      avatarUrl
    }
    orders {
      orderNumber
      total
      status
      items {
        quantity
        product {
          name
          price
          inventory
          images {
            url
          }
        }
      }
    }
  }
}
```

This query would be processed as follows:

1. Gateway receives query and creates a query plan
2. User service resolves basic user fields and profile
3. Order service resolves user.orders field using the user's ID
4. Product service resolves product details for items in orders
5. Gateway composes the complete response

All this happens transparently to the client, who sees only a single GraphQL API.

### 4.4 Performance Comparison

**REST API Gateway Flow**:
1. Client → Gateway: Request user and orders
2. Gateway → User Service: Get user
3. Gateway → Order Service: Get user's orders
4. Gateway aggregates responses
5. Gateway → Client: Combined response

**GraphQL Federation Flow**:
1. Client → Gateway: Single GraphQL query
2. Gateway analyzes query
3. Parallel requests to services
4. Gateway combines results
5. Client gets exactly what was requested

## 5. Implementation with Go

### 5.1 Technology Stack
- **GraphQL Server**: gqlgen (Go)
- **Federation**: Apollo Federation 2
- **Gateway**: Apollo Router or Go-based alternative like Bramble
- **Service Discovery**: Consul
- **Containerization**: Docker & Kubernetes
- **Database Access**: GORM/SQLx/Ent
- **Authentication**: JWT with Go middleware
- **Metrics**: Prometheus with GraphQL instrumentation
- **Tracing**: OpenTelemetry

#### Why Go for GraphQL?
- **Performance**: Go's lightweight concurrency model (goroutines) aligns perfectly with GraphQL's parallel execution needs
- **Type Safety**: Go's strong typing pairs naturally with GraphQL's type system
- **Memory Efficiency**: Low overhead for high-throughput API services
- **Compilation**: Catches errors at build time rather than runtime
- **Simplicity**: Straightforward learning curve for developers

### 5.2 Go GraphQL Tooling

#### gqlgen Configuration

**gqlgen.yml**:
```yaml
schema:
  - "*.graphql"

federation:
  filename: graph/federation.go
  package: graph
  version: 2

model:
  filename: graph/models/models.go
  package: models

resolver:
  layout: follow-schema
  dir: graph
  package: graph
  filename_template: "{name}.resolvers.go"

autobinding:
  - "github.com/your-org/user-service/internal/model"

models:
  ID:
    model:
      - github.com/99designs/gqlgen/graphql.ID
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Int:
    model:
      - github.com/99designs/gqlgen/graphql.Int
      - github.com/99designs/gqlgen/graphql.Int64
      - github.com/99designs/gqlgen/graphql.Int32
  Time:
    model: 
      - github.com/99designs/gqlgen/graphql.Time
```

#### Apollo Router Configuration (router.yaml)
```yaml
supergraph:
  listen: 0.0.0.0:4000

federation:
  enabled: true
  endpoints:
    - name: users
      url: http://user-service:8080/query
    - name: orders
      url: http://order-service:8080/query
    - name: products
      url: http://product-service:8080/query

headers:
  all:
    request:
      - propagate:
          named: "Authorization"

tracing:
  enabled: true
  provider: opentelemetry

caching:
  enabled: true
  ttl: 60s
```

#### Go GraphQL Development Tools

- **GraphQL Playground**: Interactive API explorer with documentation
- **VS Code GraphQL plugins**: Schema validation and autocomplete
- **gqlgen-lint**: Validate GraphQL schema against best practices
- **graphql-inspector**: Detect schema changes and potential breaks
- **GraphQL Voyager**: Visualize GraphQL schema relationships

### 5.3 Example: User Service Implementation

#### Schema Definition
**Schema (user/schema.graphql)**
```graphql
extend schema @link(url: "https://specs.apollo.dev/federation/v2.0", 
                   import: ["@key", "@shareable"])

type User @key(fields: "id") {
  id: ID!
  name: String!
  email: String! @shareable
  createdAt: Time!
  # User profile information
  profile: Profile
  # Referenced by other services
  orders: [Order] @provides(fields: "status")
}

type Profile {
  bio: String
  avatarUrl: String
  phoneNumber: String
}

type Query {
  user(id: ID!): User
  users(limit: Int = 10, offset: Int = 0): [User!]!
  searchUsers(term: String!): [User!]!
}

type Mutation {
  createUser(input: UserInput!): UserPayload!
  updateUser(id: ID!, input: UserInput!): UserPayload!
}

input UserInput {
  name: String!
  email: String!
  password: String!
  profile: ProfileInput
}

input ProfileInput {
  bio: String
  avatarUrl: String
  phoneNumber: String
}

type UserPayload {
  user: User
  errors: [Error!]
}

type Error {
  field: String
  message: String!
}

scalar Time
```

#### Project Structure
```
/user-service
├── cmd
│   └── server
│       └── main.go         # Entry point
├── graph
│   ├── generated           # Generated code by gqlgen
│   ├── model               # GraphQL models
│   ├── resolver.go         # Main resolver
│   ├── schema.graphql      # Schema definition
│   └── schema.resolvers.go # Resolver implementations
├── internal
│   ├── auth               # Authentication logic
│   ├── database           # Database access layer
│   └── loader             # DataLoader implementations
└── gqlgen.yml            # gqlgen configuration
```

#### Go Implementation
**Go Resolver (graph/schema.resolvers.go)**
```go
package graph

import (
    "context"
    "time"
    "github.com/99designs/gqlgen/graphql"
    "github.com/your-org/user-service/graph/model"
    "github.com/your-org/user-service/internal/database"
    "github.com/your-org/user-service/internal/loader"
)

// Entity resolver for Federation - critical for service composition
func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
    // Leverage context-based caching to optimize entity resolution
    if cached, ok := loader.UserFromContext(ctx, id); ok {
        return cached, nil
    }
    
    // Fall back to database if not in cache
    user, err := r.DB.Users.FindByID(ctx, id)
    if err != nil {
        r.Logger.WithError(err).Error("Failed to find user by ID")
        return nil, err
    }
    
    return &model.User{
        ID:        user.ID,
        Name:      user.Name,
        Email:     user.Email,
        CreatedAt: user.CreatedAt,
        Profile:   mapProfileToGraphQL(user.Profile),
    }, nil
}

// Query resolvers with pagination
func (r *queryResolver) Users(ctx context.Context, limit *int, offset *int) ([]*model.User, error) {
    // Authorization check using context
    if !hasPermission(ctx, "users:read") {
        return nil, fmt.Errorf("unauthorized")
    }
    
    l := 10
    if limit != nil {
        l = *limit
    }
    
    o := 0
    if offset != nil {
        o = *offset
    }
    
    users, err := r.DB.Users.List(ctx, l, o)
    if err != nil {
        return nil, err
    }
    
    result := make([]*model.User, len(users))
    for i, u := range users {
        result[i] = mapUserToGraphQL(u)
    }
    
    return result, nil
}

// Mutation with input validation
func (r *mutationResolver) CreateUser(ctx context.Context, input model.UserInput) (*model.UserPayload, error) {
    // Validate input
    errors := validateUserInput(input)
    if len(errors) > 0 {
        return &model.UserPayload{
            Errors: errors,
        }, nil
    }
    
    // Hash password and create user
    hashedPassword, err := hashPassword(input.Password)
    if err != nil {
        return nil, err
    }
    
    user, err := r.DB.Users.Create(ctx, database.CreateUserParams{
        Name:     input.Name,
        Email:    input.Email,
        Password: hashedPassword,
        Profile:  mapProfileFromGraphQL(input.Profile),
    })
    
    if err != nil {
        // Handle unique constraint errors
        if isUniqueViolation(err, "email") {
            return &model.UserPayload{
                Errors: []*model.Error{
                    {Field: "email", Message: "Email already in use"},
                },
            }, nil
        }
        return nil, err
    }
    
    // Return the created user
    return &model.UserPayload{
        User: mapUserToGraphQL(user),
    }, nil
}
```

#### DataLoader Implementation for N+1 Query Prevention
**internal/loader/user.go**
```go
package loader

import (
    "context"
    "time"
    "github.com/graph-gophers/dataloader"
    "github.com/your-org/user-service/graph/model"
    "github.com/your-org/user-service/internal/database"
)

type ctxKey string

const loaderKey = ctxKey("dataloader")

type Loaders struct {
    UserByID       *dataloader.Loader
    OrdersByUserID *dataloader.Loader
}

// Middleware adds dataloaders to the request context
func Middleware(db *database.DB) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            loaders := &Loaders{
                UserByID: dataloader.NewBatchedLoader(userBatchFn(db), 
                    dataloader.WithCache(&dataloader.NoCache{}),
                    dataloader.WithWait(1*time.Millisecond),
                ),
                OrdersByUserID: dataloader.NewBatchedLoader(ordersBatchFn(db),
                    dataloader.WithCache(&dataloader.NoCache{}),
                    dataloader.WithWait(1*time.Millisecond),
                ),
            }
            
            ctx := context.WithValue(r.Context(), loaderKey, loaders)
            next.ServeHTTP(w, r.WithContext(ctx))
        })
    }
}

// BatchFunction for loading users efficiently
func userBatchFn(db *database.DB) dataloader.BatchFunc {
    return func(ctx context.Context, keys dataloader.Keys) []*dataloader.Result {
        // Get all IDs to fetch at once
        ids := make([]string, len(keys))
        for i, key := range keys {
            ids[i] = key.String()
        }
        
        // Single query to fetch all users
        users, err := db.Users.FindByIDs(ctx, ids)
        if err != nil {
            return makeBatchError(keys, err)
        }
        
        // Build a map for O(1) lookup
        userMap := make(map[string]*model.User, len(users))
        for _, user := range users {
            userMap[user.ID] = mapUserToGraphQL(user)
        }
        
        // Ensure results match requested keys order
        results := make([]*dataloader.Result, len(keys))
        for i, key := range keys {
            id := key.String()
            if user, ok := userMap[id]; ok {
                results[i] = &dataloader.Result{Data: user, Error: nil}
            } else {
                results[i] = &dataloader.Result{Data: nil, Error: fmt.Errorf("user not found: %s", id)}
            }
        }
        
        return results
    }
}
```

## 6. Frontend Integration

### 6.1 Code Generation with urql for React

For frontend development, a GraphQL client like `urql` provides a seamless, type-safe experience by generating code directly from the backend's schema. This process is backend-agnostic; it only requires access to the GraphQL endpoint to introspect the schema.

**Key Benefits of this Approach:**

-   **Type-Safe Code Generation**: Using `graphql-code-generator` with `urql`, we can automatically generate React hooks and components from our GraphQL queries. This eliminates manual typing and reduces runtime errors.
-   **Improved Developer Experience**: Frontend developers get autocompletion for queries, mutations, and variables, which speeds up development and reduces bugs.
-   **Decoupling**: The frontend can develop independently of the backend. As long as the schema contract is respected, teams can work in parallel.

**Example Workflow:**

1.  **Introspect Schema**: The frontend tooling points to the GraphQL gateway's endpoint to fetch the latest schema.
2.  **Define Queries**: Frontend developers write GraphQL queries in `.graphql` files.
3.  **Run Code Generator**: A script runs `graphql-code-generator` using the introspected schema and the defined queries.
4.  **Use Generated Hooks**: The generator outputs custom React hooks (e.g., `useGetUserQuery`, `useCreateUserMutation`) that are fully typed and ready to be used in components.


```graphql
query GetUser($id: ID!) {
  user(id: $id) {
    id
    name
    email
  }
}
```

```jsx
// Example of using a generated hook in a React component
import { useGetUserQuery } from './generated/graphql';

function UserProfile({ id }) {
  const [result] = useGetUserQuery({ variables: { id } });

  const { data, fetching, error } = result;

  if (fetching) return <p>Loading...</p>;
  if (error) return <p>Oh no... {error.message}</p>;

  return (
    <div>
      <h2>{data.user.name}</h2>
      <p>{data.user.email}</p>
    </div>
  );
}
```

This approach ensures that our frontend and backend are always in sync, and any breaking changes in the API are caught at compile time on the frontend.

## 7. Challenges and Considerations
### 7.1 Performance Concerns
- N+1 query problem
- Complex queries impacting database performance
- Caching complexity

### 7.2 Security
- Rate limiting implementation
- Query depth and complexity limits
- Authentication and authorization

### 7.3 Learning Curve
- Team training requirements
- New debugging patterns
- Tooling and infrastructure changes

## 8. Conclusion

### 8.1 GraphQL Federation Benefits for Our New Project
- Scalable architecture from day one
- Clear service boundaries promoting team autonomy
- Strong typing system and schema validation
- Flexible client data fetching
- Self-documenting API
- Superior developer experience

### 8.2 Success Criteria
- Reduced development time for new features
- Improved frontend-backend collaboration
- Decreased API maintenance overhead
- Better performance through optimized data fetching
- Increased developer satisfaction

### 8.3 Key Benefits Realized
- Improved developer experience
- Better performance
- Stronger type safety
- More flexible architecture