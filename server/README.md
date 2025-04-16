# SureSQL Backend

SureSQL is a robust SQL database abstraction service that provides a RESTful API for accessing and manipulating SQL databases. It's designed to be a secure, reliable, and easy-to-use interface for applications that need to interact with SQL databases.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Features](#features)
- [Configuration](#configuration)
- [Authentication](#authentication)
- [API Endpoints](#api-endpoints)
  - [Authentication and Connection](#authentication-and-connection)
  - [Database Operations](#database-operations)
- [Usage Examples](#usage-examples)
  - [Connect to the Database](#connect-to-the-database)
  - [Refresh Token](#refresh-token)
  - [Execute SQL Statements](#execute-sql-statements)
  - [Query Data](#query-data)
  - [SQL Query](#sql-query)
  - [Insert Data](#insert-data)
  - [Get Database Status](#get-database-status)
  - [Get Schema](#get-schema)
- [Internal API](#internal-api)
- [Error Handling](#error-handling)

## Architecture Overview

SureSQL acts as a middleware layer between client applications and the underlying SQL database system. The current implementation supports RQLite as the database backend, with potential for future support of additional database systems.

Key components:
- RESTful API server with JWT-based authentication
- Connection pooling for improved performance
- Middleware for security and logging
- ORM (Object-Relational Mapping) for simplified data operations

## Features

- **Secure Authentication**: API key and token-based authentication
- **Connection Pooling**: Efficient management of database connections
- **Parameterized Queries**: Protection against SQL injection
- **Token Refresh**: Support for token refresh to maintain session validity
- **Detailed Logging**: Comprehensive logging of all operations
- **Database Status**: Get information about the database status and node configuration
- **Structured Response Format**: Consistent response structure for all API calls

## Configuration

SureSQL is configured via environment files:
- `.env.dev` - Development configuration
- `.env.simplehttp` - Server configuration

The main configuration options, this is used to connect to the DBMS include:
- `DB_HOST`, `DB_PORT`: Database server connection details
- `DB_USERNAME`, `DB_PASSWORD`: Database credentials
- `DB_SSL`: Whether to use SSL for database connections
- `DB_API_KEY`, `DB_CLIENT_ID`: API key and client ID for authentication
- `DB_CONSISTENCY`: Consistency level for distributed database operations
- `DB_OPTIONS`: Options for the DBMS
- `DB_HTTP_TIMEOUT`, `DB_RETRY_TIMEOUT`, `DB_MAX_RETRIES`: Connection parameters

Information regarding SureSQL service that will be returned to the client is in the DB itself.
These settings are also in the environment:
- `SURESQL_HOST`, `SURESQL_PORT`: SureSQL server connection details
- `SURESQL_IP`: SureSQL server IP (which are not used at this moment)
- `SURESQL_SSL`: Whether to use SSL for database connections (always true)
- `SURESQL_DBMS`: The DBMS used by SureSQL (default is RQLite)
Currently the environment takes the precedence, especially if the settings in DB table value is empty. Some of the boolean settings definitely overwritten by environment variables.

## Authentication

SureSQL uses a two-level authentication system:

1. **API Key Authentication**: Every request must include the API key and client ID in the headers
2. **Token Authentication**: After initial authentication, operations use a token-based system

The API key and client ID must be included in the headers for all requests:
```
API_KEY: your-api-key
CLIENT_ID: your-client-id
```

For authenticated endpoints, the token must be included in the Authorization header:
```
Authorization: Bearer your-token
```

## API Endpoints

### Authentication and Connection

#### POST /db/connect

Authenticates a user and creates a new database connection.

**Request Body**:
```json
{
  "username": "your-username",
  "password": "your-password"
}
```

**Response**:
```json
{
  "status": 200,
  "message": "Authentication successful",
  "data": {
    "token": "your-auth-token",
    "refresh_token": "your-refresh-token",
    "token_expired_at": "2023-01-01T12:00:00Z",
    "refresh_expired_at": "2023-01-02T12:00:00Z",
    "user_id": "1"
  }
}
```

#### POST /db/refresh

Refreshes an authentication token.

**Request Body**:
```json
{
  "refresh": "your-refresh-token"
}
```

**Response**:
```json
{
  "status": 200,
  "message": "Token refreshed successfully",
  "data": {
    "token": "your-new-auth-token",
    "refresh_token": "your-new-refresh-token",
    "token_expired_at": "2023-01-01T12:00:00Z",
    "refresh_expired_at": "2023-01-02T12:00:00Z",
    "user_id": "1"
  }
}
```

### Database Operations

All database operation endpoints require a valid authentication token.

#### POST /db/api/sql

Executes one or more SQL statements.

**Request Body**:
```json
{
  "statements": ["SQL statement 1", "SQL statement 2"],
  "param_sql": [
    {
      "query": "INSERT INTO table (column1, column2) VALUES (?, ?)",
      "values": ["value1", 2]
    }
  ]
}
```

**Response**:
```json
{
  "status": 200,
  "message": "SQL executed successfully",
  "data": {
    "results": [
      {
        "error": null,
        "timing": 0.005,
        "rows_affected": 1,
        "last_insert_id": 123
      }
    ],
    "execution_time": 0.005,
    "rows_affected": 1
  }
}
```

#### POST /db/api/query

Queries data from a table with optional conditions.

**Request Body**:
```json
{
  "table": "users",
  "condition": {
    "field": "age",
    "operator": ">",
    "value": 18,
    "order_by": ["name ASC"],
    "limit": 10
  },
  "single_row": false
}
```

**Response**:
```json
{
  "status": 200,
  "message": "Query executed successfully",
  "data": {
    "records": [
      {
        "table_name": "users",
        "data": {
          "id": 1,
          "name": "John Doe",
          "age": 30
        }
      }
    ],
    "execution_time": 0.003,
    "count": 1
  }
}
```

#### POST /db/api/querysql

Executes SQL queries and returns the results.

**Request Body**:
```json
{
  "statements": ["SELECT * FROM users WHERE age > 18"],
  "param_sql": [
    {
      "query": "SELECT * FROM users WHERE age > ?",
      "values": [18]
    }
  ],
  "single_row": false
}
```

**Response**:
```json
{
  "status": 200,
  "message": "SQL executed successfully",
  "data": [
    {
      "records": [
        {
          "table_name": "users",
          "data": {
            "id": 1,
            "name": "John Doe",
            "age": 30
          }
        }
      ],
      "execution_time": 0.003,
      "count": 1
    }
  ]
}
```

#### POST /db/api/insert

Inserts one or more records into the database.

**Request Body**:
```json
{
  "records": [
    {
      "table_name": "users",
      "data": {
        "name": "Jane Smith",
        "age": 25
      }
    }
  ],
  "queue": false,
  "same_table": true
}
```

**Response**:
```json
{
  "status": 200,
  "message": "Successfully inserted 1 records",
  "data": {
    "results": [
      {
        "error": null,
        "timing": 0.004,
        "rows_affected": 1,
        "last_insert_id": 124
      }
    ],
    "execution_time": 0.004,
    "rows_affected": 1
  }
}
```

#### GET /db/api/status

Retrieves the status of the database connection.

**Response**:
```json
{
  "status": 200,
  "message": "Status peers vs config matched",
  "data": {
    "url": "http://localhost:4001",
    "version": "0.0.1",
    "dbms": "rqlite",
    "dbms_driver": "direct-rqlite",
    "start_time": "2023-01-01T00:00:00Z",
    "uptime": "24h0m0s",
    "dir_size": 1024,
    "db_size": 2048,
    "node_id": "1",
    "is_leader": true,
    "leader": "http://localhost:4001",
    "mode": "rw",
    "nodes": 1,
    "node_number": 1,
    "max_pool": 25,
    "peers": {
      "1": {
        "url": "http://localhost:4001",
        "is_leader": true,
        "mode": "rw",
        "nodes": 1,
        "node_number": 1
      }
    }
  }
}
```

#### POST /db/api/getschema

Retrieves the database schema information.

**Response**:
```json
{
  "status": 200,
  "message": "Schema get successfully",
  "data": [
    {
      "type": "table",
      "name": "users",
      "tbl_name": "users",
      "rootpage": 2,
      "sql": "CREATE TABLE users (id INTEGER PRIMARY KEY, name TEXT, age INTEGER)",
      "hidden": false
    }
  ]
}
```

## Usage Examples

### Connect to the Database

```javascript
const response = await fetch('http://your-suresql-server/db/connect', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id'
  },
  body: JSON.stringify({
    username: 'your-username',
    password: 'your-password'
  })
});

const data = await response.json();
const token = data.data.token;
// Save token for subsequent requests
```

### Refresh Token

```javascript
const response = await fetch('http://your-suresql-server/db/refresh', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id'
  },
  body: JSON.stringify({
    refresh: 'your-refresh-token'
  })
});

const data = await response.json();
const newToken = data.data.token;
// Update the saved token
```

### Execute SQL Statements

Basic example with a single SQL statement:

```javascript
const response = await fetch('http://your-suresql-server/db/api/sql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    statements: ["CREATE TABLE if not exists test (id INTEGER PRIMARY KEY, value TEXT)"]
  })
});

const data = await response.json();
console.log(data);
```

Example with parameterized SQL:

```javascript
const response = await fetch('http://your-suresql-server/db/api/sql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    param_sql: [
      {
        query: "INSERT INTO test (value) VALUES (?)",
        values: ["test value"]
      }
    ]
  })
});

const data = await response.json();
console.log(data);
```

### Query Data

Simple query to get all records from a table:

```javascript
const response = await fetch('http://your-suresql-server/db/api/query', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    table: "test"
  })
});

const data = await response.json();
console.log(data);
```

Query with conditions:

```javascript
const response = await fetch('http://your-suresql-server/db/api/query', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    table: "test",
    condition: {
      field: "id",
      operator: ">",
      value: 5,
      order_by: ["id DESC"],
      limit: 10
    }
  })
});

const data = await response.json();
console.log(data);
```

Complex query with nested conditions:

```javascript
const response = await fetch('http://your-suresql-server/db/api/query', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    table: "users",
    condition: {
      logic: "OR",
      nested: [
        {
          field: "age",
          operator: ">",
          value: 30
        },
        {
          logic: "AND",
          nested: [
            {
              field: "status",
              operator: "=",
              value: "active"
            },
            {
              field: "role",
              operator: "=",
              value: "admin"
            }
          ]
        }
      ],
      order_by: ["name ASC"],
      limit: 20
    }
  })
});

const data = await response.json();
console.log(data);
```

### SQL Query

Execute a SQL query and get the results:

```javascript
const response = await fetch('http://your-suresql-server/db/api/querysql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    statements: ["SELECT * FROM users WHERE role = 'admin'"]
  })
});

const data = await response.json();
console.log(data);
```

Parameterized SQL query:

```javascript
const response = await fetch('http://your-suresql-server/db/api/querysql', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    param_sql: [
      {
        query: "SELECT * FROM users WHERE role = ? AND status = ?",
        values: ["admin", "active"]
      }
    ]
  })
});

const data = await response.json();
console.log(data);
```

### Insert Data

Insert a single record:

```javascript
const response = await fetch('http://your-suresql-server/db/api/insert', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    records: [
      {
        table_name: "users",
        data: {
          name: "John Smith",
          age: 35,
          role: "user",
          status: "active"
        }
      }
    ]
  })
});

const data = await response.json();
console.log(data);
```

Insert multiple records into the same table:

```javascript
const response = await fetch('http://your-suresql-server/db/api/insert', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    records: [
      {
        table_name: "users",
        data: {
          name: "Alice Johnson",
          age: 28,
          role: "user",
          status: "active"
        }
      },
      {
        table_name: "users",
        data: {
          name: "Bob Williams",
          age: 42,
          role: "admin",
          status: "active"
        }
      }
    ],
    same_table: true
  })
});

const data = await response.json();
console.log(data);
```

Insert records into different tables:

```javascript
const response = await fetch('http://your-suresql-server/db/api/insert', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  },
  body: JSON.stringify({
    records: [
      {
        table_name: "users",
        data: {
          name: "Charlie Brown",
          age: 32,
          role: "user",
          status: "active"
        }
      },
      {
        table_name: "posts",
        data: {
          title: "My First Post",
          content: "Hello, world!",
          user_id: 1
        }
      }
    ],
    same_table: false
  })
});

const data = await response.json();
console.log(data);
```

### Get Database Status

```javascript
const response = await fetch('http://your-suresql-server/db/api/status', {
  method: 'GET',
  headers: {
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  }
});

const data = await response.json();
console.log(data);
```

### Get Schema

```javascript
const response = await fetch('http://your-suresql-server/db/api/getschema', {
  method: 'POST',
  headers: {
    'Content-Type': 'application/json',
    'API_KEY': 'your-api-key',
    'CLIENT_ID': 'your-client-id',
    'Authorization': `Bearer ${token}`
  }
});

const data = await response.json();
console.log(data);
```

## Internal API

SureSQL also provides an internal API accessible only with basic authentication using the internal configuration credentials. This API is intended for administrative purposes.

Internal API endpoints:
- `/suresql/iusers` (GET, POST, PUT, DELETE) - Manage users
- `/suresql/schema` (GET) - Get database schema information
- `/suresql/dbms_status` (GET) - Get DBMS status information

## Error Handling

All endpoints return a consistent error response format:

```json
{
  "status": 400, // HTTP status code
  "message": "Error message",
  "data": null // or error details
}
```

Common error status codes:
- `400`: Bad Request - Invalid input or parameters
- `401`: Unauthorized - Missing or invalid authentication
- `404`: Not Found - Resource not found
- `500`: Internal Server Error - Server-side error

Each error response includes a descriptive message to help diagnose the issue.