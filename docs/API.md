# 📡 API Documentation

## Base URL
- Development: `http://localhost:8080`
- Production: `https://api.boltvisa.com`

## Authentication

Most endpoints require authentication via JWT token in the Authorization header:

```
Authorization: Bearer <token>
```

## Endpoints

### Authentication

#### Register User
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "first_name": "John",
    "last_name": "Doe",
    "role": "applicant",
    "active": true
  },
  "expires_at": "2024-01-02T12:00:00Z"
}
```

#### Login
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:** Same as register

#### Refresh Token
```http
POST /api/v1/auth/refresh
Authorization: Bearer <token>
```

**Response:** New token and user data

### Users

#### Get Current User
```http
GET /api/v1/users/me
Authorization: Bearer <token>
```

#### Update Current User
```http
PUT /api/v1/users/me
Authorization: Bearer <token>
Content-Type: application/json

{
  "first_name": "Jane",
  "last_name": "Smith"
}
```

### Visa Categories

#### List Categories
```http
GET /api/v1/visa-categories
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Tourist Visa",
    "description": "Short-term tourist visa",
    "country": "USA",
    "duration": "90 days",
    "price": 160.00,
    "active": true
  }
]
```

#### Get Category
```http
GET /api/v1/visa-categories/:id
Authorization: Bearer <token>
```

### Visa Applications

#### List Applications
```http
GET /api/v1/applications
Authorization: Bearer <token>
```

**Response:**
```json
[
  {
    "id": 1,
    "user_id": 1,
    "category_id": 1,
    "status": "draft",
    "passport_number": "AB123456",
    "created_at": "2024-01-01T12:00:00Z",
    "category": { ... }
  }
]
```

#### Create Application
```http
POST /api/v1/applications
Authorization: Bearer <token>
Content-Type: application/json

{
  "category_id": 1,
  "passport_number": "AB123456",
  "date_of_birth": "1990-01-01",
  "nationality": "US",
  "travel_date": "2024-06-01"
}
```

#### Update Application
```http
PUT /api/v1/applications/:id
Authorization: Bearer <token>
Content-Type: application/json

{
  "status": "submitted",
  "notes": "Ready for review"
}
```

#### Delete Application
```http
DELETE /api/v1/applications/:id
Authorization: Bearer <token>
```

### Documents

#### Upload Document
```http
POST /api/v1/applications/:id/documents
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: <file>
type: passport
```

#### List Documents
```http
GET /api/v1/applications/:id/documents
Authorization: Bearer <token>
```

#### Delete Document
```http
DELETE /api/v1/documents/:id
Authorization: Bearer <token>
```

### Notifications

#### List Notifications
```http
GET /api/v1/notifications
Authorization: Bearer <token>
```

#### Mark Notification Read
```http
PUT /api/v1/notifications/:id/read
Authorization: Bearer <token>
```

### Admin Endpoints

#### List Users (Admin Only)
```http
GET /api/v1/admin/users
Authorization: Bearer <admin_token>
```

#### Create Visa Category (Admin Only)
```http
POST /api/v1/admin/visa-categories
Authorization: Bearer <admin_token>
Content-Type: application/json

{
  "name": "Business Visa",
  "description": "Business travel visa",
  "country": "USA",
  "duration": "1 year",
  "price": 200.00
}
```

### Consultant Endpoints

#### Get Consultant Clients
```http
GET /api/v1/consultant/clients
Authorization: Bearer <consultant_token>
```

#### Get Consultant Applications
```http
GET /api/v1/consultant/applications
Authorization: Bearer <consultant_token>
```

## Error Responses

All errors follow this format:

```json
{
  "error": "Error message description"
}
```

### Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `409` - Conflict
- `500` - Internal Server Error

## Rate Limiting

Rate limiting will be implemented in Phase 4. Current limits:
- 100 requests per minute per IP (planned)
- 1000 requests per hour per authenticated user (planned)

