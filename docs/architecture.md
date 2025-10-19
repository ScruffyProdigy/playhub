# Architecture Overview

PlayHub is a gaming lobby platform built with a modern microservices architecture.

## System Architecture

### Development Environment
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

### Production Environment (Kubernetes)
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   React Frontend │    │  Go GraphQL API │    │   PostgreSQL    │
│   (Nginx)       │◄──►│   (Container)   │◄──►│   (Container)   │
│   Port 80       │    │   Port 8080     │    │   Port 5432     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                       │                       │
         │                       │                       │
         ▼                       ▼                       ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Environment   │    │   Kubernetes    │    │   Persistent    │
│   ConfigMaps    │    │   Services      │    │   Volumes       │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Component Overview

### Frontend (React + Vite) ✅
- **Technology**: React 19, Vite, Vitest, Playwright
- **Purpose**: User interface for gaming lobby
- **Status**: Foundation implemented with testing infrastructure
- **Features**: Basic UI, environment configuration, comprehensive testing
- **Port**: 5173 (development), 80 (production)

### Backend (Go + GraphQL) 🚧
- **Technology**: Go 1.25, gqlgen, GraphQL
- **Purpose**: API server and business logic
- **Status**: Schema and mock resolvers implemented
- **Features**: GraphQL API with mock data, health checks, version endpoint
- **Port**: 8080
- **In Progress**: Database integration, authentication, real business logic

### Database (PostgreSQL) 📋
- **Technology**: PostgreSQL
- **Purpose**: Data persistence
- **Status**: Planned - not yet implemented
- **Features**: User data, game sessions, trading history
- **Next Steps**: Schema design, migration system, connection pooling

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

## Environment Configuration System

PlayHub uses a Docker-based environment configuration system that allows the same Docker image to work across different environments.

### How It Works

1. **Docker Entrypoint Script**: The frontend Docker image includes a script that generates `env.js` at runtime
2. **Kubernetes ConfigMaps**: Each environment has its own ConfigMap with environment-specific values
3. **Runtime Injection**: When the container starts, it reads environment variables and creates the `env.js` file
4. **Frontend Access**: The frontend loads `/env.js` and accesses variables via `window.env`

### Environment Configurations

#### Local Development
- **API URL**: `http://localhost:8081`
- **Environment**: `local`
- **Deployment**: Kubernetes with port-forwarding

#### Staging
- **API URL**: `https://api-staging.playhub.com`
- **Environment**: `staging`
- **Deployment**: Kubernetes cluster with staging ConfigMaps

#### Production
- **API URL**: `https://api.playhub.com`
- **Environment**: `production`
- **Deployment**: Kubernetes cluster with production ConfigMaps

### Benefits

- **Single Docker Image**: Same image works in all environments
- **Runtime Configuration**: No need to rebuild for different environments
- **Kubernetes Native**: Uses ConfigMaps for environment-specific values
- **Secure**: Sensitive values can be stored in Secrets
- **Easy Deployment**: Simple scripts for each environment

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
