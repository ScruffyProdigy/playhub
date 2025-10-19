# API Documentation

This document describes the GraphQL API for PlayHub.

## Implementation Status

### ✅ Implemented
- **System Queries**: `version`, `healthz` - Basic system information
- **GraphQL Schema**: Complete schema definition with all types
- **Mock Resolvers**: All resolvers return mock data for testing

### 🚧 In Development
- **Authentication**: JWT-based authentication system
- **Database Integration**: Real data persistence
- **Business Logic**: Actual game management and queuing

### 📋 Planned
- **Real-time Subscriptions**: WebSocket support for live updates
- **File Uploads**: Support for game assets and user avatars
- **Rate Limiting**: API rate limiting and throttling

## Base URL

- **Development**: `http://localhost:8080/query`
- **Production**: `https://api.playhub.com/query`

## Authentication

> **Note**: Authentication is currently in development. The API currently accepts requests without authentication for testing purposes.

PlayHub will use JWT-based authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Queries

### System Queries

#### `version` ✅
Returns the current API version.

```graphql
query {
  version
}
```

**Response:**
```json
{
  "data": {
    "version": "1.0.0"
  }
}
```

#### `healthz` ✅
Health check endpoint.

```graphql
query {
  healthz
}
```

**Response:**
```json
{
  "data": {
    "healthz": "ok"
  }
}
```

### User Queries

#### `me` 🚧
Get current user information. *Returns mock data - authentication in development*

```graphql
query {
  me {
    id
    email
    displayName
    createdAt
  }
}
```

**Response:**
```json
{
  "data": {
    "me": {
      "id": "user-123",
      "email": "user@example.com",
      "displayName": "Test User",
      "createdAt": "2024-01-01T00:00:00Z"
    }
  }
}
```

### Game Queries

#### `games` 🚧
List available games with pagination. *Returns mock data - database integration in development*

```graphql
query {
  games(limit: 10, offset: 0) {
    id
    name
    description
    maxPlayers
    status
  }
}
```

**Response:**
```json
{
  "data": {
    "games": [
      {
        "id": "game-1",
        "name": "Example Game",
        "description": "A fun example game",
        "maxPlayers": 4,
        "status": "ACTIVE"
      }
    ]
  }
}
```

#### `game` 🚧
Get a specific game by ID. *Returns mock data - database integration in development*

```graphql
query {
  game(id: "game-1") {
    id
    name
    description
    maxPlayers
    status
    currentPlayers
  }
}
```

### Session Queries

#### `session`
Get session information.

```graphql
query {
  session(id: "session-123") {
    id
    gameId
    userId
    status
    joinedAt
  }
}
```

### Digital Goods Queries

#### `goods`
List digital goods, optionally filtered by game.

```graphql
query {
  goods(gameId: "game-1") {
    id
    name
    description
    price
    gameId
    type
  }
}
```

#### `myInventory`
Get user's digital goods inventory.

```graphql
query {
  myInventory(gameId: "game-1") {
    id
    goodId
    userId
    gameId
    acquiredAt
    status
  }
}
```

## Mutations

> **Note**: All mutations currently return mock data. Real database integration is in development.

### Game Management

#### `createGame` 🚧
Create a new game. *Returns mock data - database integration in development*

```graphql
mutation {
  createGame(input: {
    name: "New Game"
    description: "A new game to play"
    maxPlayers: 4
  }) {
    id
    name
    description
    maxPlayers
    status
  }
}
```

### Queue Management

#### `joinQueue` 🚧
Join a game queue. *Returns mock data - queue system in development*

```graphql
mutation {
  joinQueue(gameId: "game-1") {
    id
    gameId
    userId
    status
    joinedAt
  }
}
```

#### `leaveQueue`
Leave a game queue.

```graphql
mutation {
  leaveQueue(sessionId: "session-123")
}
```

### Trading

#### `purchaseGood`
Purchase a digital good.

```graphql
mutation {
  purchaseGood(input: {
    goodId: "good-1"
    gameId: "game-1"
  }) {
    id
    goodId
    userId
    gameId
    acquiredAt
    status
  }
}
```

#### `tradeGood`
Trade a digital good with another user.

```graphql
mutation {
  tradeGood(input: {
    goodId: "good-1"
    fromUserId: "user-1"
    toUserId: "user-2"
    gameId: "game-1"
  }) {
    id
    goodId
    fromUserId
    toUserId
    gameId
    tradedAt
    status
  }
}
```

## Error Handling

GraphQL returns errors in a standardized format:

```json
{
  "data": null,
  "errors": [
    {
      "message": "Game not found",
      "path": ["game"],
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ]
    }
  ]
}
```

### Common Error Codes

- `GAME_NOT_FOUND`: The specified game doesn't exist
- `SESSION_NOT_FOUND`: The specified session doesn't exist
- `UNAUTHORIZED`: Authentication required
- `FORBIDDEN`: Insufficient permissions
- `VALIDATION_ERROR`: Input validation failed
- `RATE_LIMITED`: Too many requests

## Rate Limiting

API requests are rate limited to prevent abuse:

- **Authenticated users**: 1000 requests per hour
- **Unauthenticated users**: 100 requests per hour

Rate limit headers are included in responses:

```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1640995200
```

## Examples

### Complete User Workflow

1. **Get user info**
```graphql
query {
  me {
    id
    displayName
  }
}
```

2. **Browse games**
```graphql
query {
  games(limit: 5) {
    id
    name
    description
    maxPlayers
  }
}
```

3. **Join a game queue**
```graphql
mutation {
  joinQueue(gameId: "game-1") {
    id
    status
  }
}
```

4. **Check session status**
```graphql
query {
  session(id: "session-123") {
    id
    status
    gameId
  }
}
```

5. **Browse digital goods**
```graphql
query {
  goods(gameId: "game-1") {
    id
    name
    price
  }
}
```

6. **Purchase a good**
```graphql
mutation {
  purchaseGood(input: {
    goodId: "good-1"
    gameId: "game-1"
  }) {
    id
    status
  }
}
```
