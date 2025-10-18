# Architecture Overview

PlayHub is a gaming lobby platform built with a modern microservices architecture.

## System Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React Frontend │    │  Go GraphQL API │    │   PostgreSQL    │
│   (Port 5173)   │◄──►│   (Port 8080)   │◄──►│   (Port 5432)   │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Vite Dev      │    │   gqlgen        │    │   Database      │
│   Server        │    │   Code Gen      │    │   Migrations    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Component Overview

### Frontend (React + Vite)
- **Technology**: React 19, Vite, Vitest, Playwright
- **Purpose**: User interface for gaming lobby
- **Features**: Game queuing, user management, digital goods trading
- **Port**: 5173 (development)

### Backend (Go + GraphQL)
- **Technology**: Go 1.25, gqlgen, GraphQL
- **Purpose**: API server and business logic
- **Features**: User authentication, game management, trading system
- **Port**: 8080

### Database (PostgreSQL)
- **Technology**: PostgreSQL
- **Purpose**: Data persistence
- **Features**: User data, game sessions, trading history

## Data Flow

### User Authentication Flow
```
User → Frontend → GraphQL API → Database
     ← JWT Token ←
```

### Game Queue Flow
```
User → Join Queue → GraphQL API → Database
     ← Queue Status ←
```

### Trading Flow
```
User → Trade Request → GraphQL API → Database
     ← Trade Confirmation ←
```

## API Design

### GraphQL Schema Structure
```
Query {
  version: String!
  healthz: String!
  me: User
  games(limit: Int, offset: Int): [Game!]!
  game(id: ID!): Game
  session(id: ID!): Session
  goods(gameId: ID): [DigitalGood!]!
  myInventory(gameId: ID): [Entitlement!]!
}

Mutation {
  createGame(input: CreateGameInput!): Game!
  joinQueue(gameId: ID!): Session!
  leaveQueue(sessionId: ID!): Boolean!
  purchaseGood(input: PurchaseGoodInput!): Entitlement!
  tradeGood(input: TradeGoodInput!): Trade!
}
```

## Security

### Authentication
- JWT-based authentication
- Secure token storage
- Session management

### Authorization
- Role-based access control
- Resource-level permissions
- API rate limiting

## Deployment

### Development
- Local development servers
- Hot reloading enabled
- Test databases

### Production
- Kubernetes deployment
- Containerized services
- Load balancing
- SSL/TLS termination

## Monitoring & Observability

### Logging
- Structured logging (JSON)
- Log levels (DEBUG, INFO, WARN, ERROR)
- Request tracing

### Metrics
- Application metrics
- Performance monitoring
- Error tracking

### Health Checks
- `/healthz` endpoint
- Database connectivity
- External service status

## Scalability Considerations

### Horizontal Scaling
- Stateless API design
- Database connection pooling
- Load balancer ready

### Performance
- GraphQL query optimization
- Database indexing
- Caching strategies

### Future Enhancements
- Microservices architecture
- Event-driven communication
- Message queues
- Caching layers
