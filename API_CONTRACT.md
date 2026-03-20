# API Contract
## Coffee Shop Point of Sales System

---

## 1. Base URL dan Versioning

```
Base URL: /api/v1
```

Semua endpoint API menggunakan prefix `/api/v1`. Versioning dilakukan melalui URL path untuk memudahkan backward compatibility di masa depan.

**Contoh:**
```
GET /api/v1/products
POST /api/v1/auth/login
```

---

## 2. Authentication

Semua endpoint protected memerlukan authentication header dengan format Bearer Token.

### Header Format

```http
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Token Info

- **Tipe:** JWT (JSON Web Token)
- **Expiry:** 24 jam dari waktu login
- **Payload:**
  ```json
  {
    "sub": "user_uuid",
    "username": "cashier01",
    "role": "cashier",
    "iat": 1711000000,
    "exp": 1711086400
  }
  ```

### Public Endpoints (Tidak Perlu Auth)

- `POST /api/v1/auth/login`

---

## 3. Standard Response Format

Semua endpoint API mengikuti format response yang konsisten.

### Response Sukses (Single Resource)

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "Example Resource"
  }
}
```

### Response Sukses dengan List (Pagination)

```json
{
  "success": true,
  "message": "Operation successful",
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "name": "Example Resource 1"
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Example Resource 2"
      }
    ],
    "meta": {
      "page": 1,
      "limit": 20,
      "total": 100,
      "total_pages": 5
    }
  }
}
```

### Response Error Validasi

```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "username": "Username is required",
    "email": "Invalid email format"
  }
}
```

### Response Error Umum

```json
{
  "success": false,
  "message": "Resource not found"
}
```

### HTTP Status Codes

| Code | Keterangan |
|------|------------|
| 200 | OK - Request sukses (GET, PUT) |
| 201 | Created - Resource berhasil dibuat (POST) |
| 400 | Bad Request - Request tidak valid |
| 401 | Unauthorized - Tidak ada token atau token expired |
| 403 | Forbidden - Token valid tapi tidak punya akses |
| 404 | Not Found - Data tidak ditemukan |
| 422 | Unprocessable Entity - Validasi gagal |
| 429 | Too Many Requests - Rate limit exceeded |
| 500 | Internal Server Error - Server error |

### Error Response Berdasarkan Tipe

| Tipe Error | HTTP Code | Format Response |
|------------|-----------|-----------------|
| Validasi | 422 | `{ "success": false, "message": "...", "errors": {...} }` |
| Tidak terautentikasi | 401 | `{ "success": false, "message": "..." }` |
| Tidak punya akses | 403 | `{ "success": false, "message": "..." }` |
| Data tidak ditemukan | 404 | `{ "success": false, "message": "..." }` |
| Request tidak valid | 400 | `{ "success": false, "message": "..." }` |
| Rate limit | 429 | `{ "success": false, "message": "..." }` |
| Server error | 500 | `{ "success": false, "message": "..." }` |

---

## 4. Daftar Endpoint

### 4.1 Auth

#### Login

```http
POST /api/v1/auth/login
```

**Akses:** Public

**Request Body:**
```json
{
  "username": "cashier01",
  "password": "securepassword123"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_at": "2026-03-21T10:00:00Z",
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "username": "cashier01",
      "full_name": "John Barista",
      "role": "cashier"
    }
  }
}
```

**Response Error:**
```json
{
  "success": false,
  "message": "Invalid username or password"
}
```

---

#### Get Current User

```http
GET /api/v1/auth/me
```

**Akses:** Owner, Cashier (Protected)

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "username": "cashier01",
    "full_name": "John Barista",
    "role": "cashier",
    "is_active": true,
    "created_at": "2026-01-15T08:00:00Z"
  }
}
```

---

#### Logout

```http
POST /api/v1/auth/logout
```

**Akses:** Owner, Cashier (Protected)

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Logout successful"
}
```

---

### 4.2 Categories (Owner)

#### Get All Categories

```http
GET /api/v1/categories
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page |
| `search` | string | - | Search by name |
| `include_deleted` | bool | false | Include soft deleted |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "Coffee",
        "description": "Hot and cold coffee beverages",
        "created_at": "2026-01-10T08:00:00Z",
        "updated_at": "2026-01-10T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

#### Get Category by ID

```http
GET /api/v1/categories/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "name": "Coffee",
    "description": "Hot and cold coffee beverages",
    "created_at": "2026-01-10T08:00:00Z",
    "updated_at": "2026-01-10T08:00:00Z"
  }
}
```

**Response Error (404):**
```json
{
  "success": false,
  "message": "Category not found"
}
```

---

#### Create Category

```http
POST /api/v1/categories
```

**Akses:** Owner

**Request Body:**
```json
{
  "name": "Pastry",
  "description": "Fresh baked pastries"
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Category created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Pastry",
    "description": "Fresh baked pastries",
    "created_at": "2026-03-20T10:00:00Z",
    "updated_at": "2026-03-20T10:00:00Z"
  }
}
```

**Response Error (409):**
```json
{
  "success": false,
  "message": "Category name already exists"
}
```

---

#### Update Category

```http
PUT /api/v1/categories/:id
```

**Akses:** Owner

**Request Body:**
```json
{
  "name": "Premium Pastry",
  "description": "Fresh baked premium pastries"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Category updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440002",
    "name": "Premium Pastry",
    "description": "Fresh baked premium pastries",
    "created_at": "2026-03-20T10:00:00Z",
    "updated_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Delete Category (Soft Delete)

```http
DELETE /api/v1/categories/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Category deleted successfully"
}
```

**Response Error (409):**
```json
{
  "success": false,
  "message": "Cannot delete category with active products"
}
```

---

### 4.3 Products (Owner)

#### Get All Products

```http
GET /api/v1/products
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page |
| `search` | string | - | Search by name/SKU |
| `category_id` | string | - | Filter by category |
| `is_active` | bool | - | Filter by status |
| `include_deleted` | bool | false | Include soft deleted |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440010",
        "sku": "COF-LAT-001",
        "name": "Caffe Latte",
        "description": "Espresso with steamed milk",
        "category_id": "550e8400-e29b-41d4-a716-446655440001",
        "category_name": "Coffee",
        "price": 3500000,
        "stock": 100,
        "min_stock": 10,
        "is_active": true,
        "images": [
          {
            "id": "img-001",
            "file_path": "/products/latte.jpg",
            "is_primary": true
          }
        ],
        "created_at": "2026-01-10T08:00:00Z",
        "updated_at": "2026-03-15T10:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 50,
      "total_pages": 3
    }
  }
}
```

---

#### Get Product by ID

```http
GET /api/v1/products/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "sku": "COF-LAT-001",
    "name": "Caffe Latte",
    "description": "Espresso with steamed milk",
    "category_id": "550e8400-e29b-41d4-a716-446655440001",
    "category_name": "Coffee",
    "price": 3500000,
    "stock": 100,
    "min_stock": 10,
    "is_active": true,
    "images": [
      {
        "id": "img-001",
        "file_path": "/products/latte.jpg",
        "file_name": "latte.jpg",
        "file_size": 245000,
        "mime_type": "image/jpeg",
        "is_primary": true
      }
    ],
    "created_at": "2026-01-10T08:00:00Z",
    "updated_at": "2026-03-15T10:00:00Z"
  }
}
```

---

#### Create Product

```http
POST /api/v1/products
```

**Akses:** Owner

**Request Body:**
```json
{
  "sku": "COF-LAT-001",
  "name": "Caffe Latte",
  "description": "Espresso with steamed milk",
  "category_id": "550e8400-e29b-41d4-a716-446655440001",
  "price": 3500000,
  "stock": 100,
  "min_stock": 10,
  "is_active": true,
  "images": [
    {
      "file_path": "/products/latte.jpg",
      "file_name": "latte.jpg",
      "file_size": 245000,
      "mime_type": "image/jpeg",
      "is_primary": true
    }
  ]
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Product created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "sku": "COF-LAT-001",
    "name": "Caffe Latte",
    "price": 3500000,
    "stock": 100,
    "created_at": "2026-03-20T10:00:00Z"
  }
}
```

**Response Error (422):**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "sku": "SKU is required",
    "price": "Price must be greater than 0"
  }
}
```

---

#### Update Product

```http
PUT /api/v1/products/:id
```

**Akses:** Owner

**Request Body:**
```json
{
  "name": "Caffe Latte Premium",
  "price": 4000000,
  "description": "Premium espresso with steamed milk"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Product updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "name": "Caffe Latte Premium",
    "price": 4000000,
    "updated_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Toggle Product Status

```http
PATCH /api/v1/products/:id/status
```

**Akses:** Owner

**Request Body:**
```json
{
  "is_active": false
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Product status updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440010",
    "is_active": false,
    "updated_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Delete Product (Soft Delete)

```http
DELETE /api/v1/products/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Product deleted successfully"
}
```

---

#### Upload Product Image

```http
POST /api/v1/products/:id/images
Content-Type: multipart/form-data
```

**Akses:** Owner

**Request Body (FormData):**
| Field | Type | Keterangan |
|-------|------|------------|
| `image` | file | Image file (max 2MB, JPG/PNG) |
| `is_primary` | bool | Set as primary image |

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Image uploaded successfully",
  "data": {
    "id": "img-002",
    "product_id": "550e8400-e29b-41d4-a716-446655440010",
    "file_path": "/products/latte-2.jpg",
    "file_name": "latte-premium.jpg",
    "file_size": 320000,
    "mime_type": "image/jpeg",
    "is_primary": false,
    "created_at": "2026-03-20T10:10:00Z"
  }
}
```

**Response Error (422):**
```json
{
  "success": false,
  "message": "Invalid file. Max size 2MB, allowed formats: JPG, PNG",
  "errors": {
    "image": "Invalid file. Max size 2MB, allowed formats: JPG, PNG"
  }
}
```

---

### 4.4 Stock (Owner)

#### Get Stock Overview

```http
GET /api/v1/stock
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `low_stock` | bool | false | Show only low stock items |
| `search` | string | - | Search by product name/SKU |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440010",
        "product_name": "Caffe Latte",
        "sku": "COF-LAT-001",
        "current_stock": 8,
        "min_stock": 10,
        "is_low_stock": true
      }
    ],
    "summary": {
      "total_products": 50,
      "low_stock_count": 5,
      "out_of_stock_count": 2
    }
  }
}
```

---

#### Get Stock Movements

```http
GET /api/v1/stock/movements
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page |
| `product_id` | string | - | Filter by product |
| `movement_type` | string | - | Filter: in/out/adjustment |
| `source` | string | - | Filter: purchase/waste/sale/etc |
| `start_date` | date | - | Filter from date |
| `end_date` | date | - | Filter to date |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440100",
        "product_id": "550e8400-e29b-41d4-a716-446655440010",
        "product_name": "Caffe Latte",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "user_name": "Owner Admin",
        "movement_type": "in",
        "source": "purchase",
        "quantity": 50,
        "previous_stock": 50,
        "current_stock": 100,
        "notes": "Restock from supplier",
        "created_at": "2026-03-20T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 150,
      "total_pages": 8
    }
  }
}
```

---

#### Add Stock (Stock In)

```http
POST /api/v1/stock/in
```

**Akses:** Owner

**Request Body:**
```json
{
  "product_id": "550e8400-e29b-41d4-a716-446655440010",
  "quantity": 50,
  "source": "purchase",
  "notes": "Restock from supplier ABC"
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Stock added successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440100",
    "product_id": "550e8400-e29b-41d4-a716-446655440010",
    "movement_type": "in",
    "source": "purchase",
    "quantity": 50,
    "previous_stock": 50,
    "current_stock": 100,
    "created_at": "2026-03-20T08:00:00Z"
  }
}
```

**Response Error (422):**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "quantity": "Quantity must be greater than 0"
  }
}
```

---

#### Reduce Stock (Stock Out)

```http
POST /api/v1/stock/out
```

**Akses:** Owner

**Request Body:**
```json
{
  "product_id": "550e8400-e29b-41d4-a716-446655440010",
  "quantity": 5,
  "source": "waste",
  "notes": "Spilled during preparation"
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Stock reduced successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440101",
    "product_id": "550e8400-e29b-41d4-a716-446655440010",
    "movement_type": "out",
    "source": "waste",
    "quantity": -5,
    "previous_stock": 100,
    "current_stock": 95,
    "created_at": "2026-03-20T09:00:00Z"
  }
}
```

**Response Error (422):**
```json
{
  "success": false,
  "message": "Insufficient stock"
}
```

---

### 4.5 Tables (Owner)

#### Get All Tables

```http
GET /api/v1/tables
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `status` | string | - | Filter: available/occupied/reserved |
| `include_deleted` | bool | false | Include soft deleted |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440050",
        "table_number": "T01",
        "capacity": 4,
        "location": "indoor",
        "status": "available",
        "created_at": "2026-01-10T08:00:00Z",
        "updated_at": "2026-03-20T10:00:00Z"
      },
      {
        "id": "550e8400-e29b-41d4-a716-446655440051",
        "table_number": "T02",
        "capacity": 2,
        "location": "outdoor",
        "status": "occupied",
        "created_at": "2026-01-10T08:00:00Z",
        "updated_at": "2026-03-20T11:00:00Z"
      }
    ]
  }
}
```

---

#### Create Table

```http
POST /api/v1/tables
```

**Akses:** Owner

**Request Body:**
```json
{
  "table_number": "T10",
  "capacity": 6,
  "location": "vip"
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Table created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440060",
    "table_number": "T10",
    "capacity": 6,
    "location": "vip",
    "status": "available",
    "created_at": "2026-03-20T10:00:00Z"
  }
}
```

---

#### Update Table

```http
PUT /api/v1/tables/:id
```

**Akses:** Owner

**Request Body:**
```json
{
  "table_number": "T10A",
  "capacity": 8,
  "location": "vip"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Table updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440060",
    "table_number": "T10A",
    "capacity": 8,
    "updated_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Delete Table (Soft Delete)

```http
DELETE /api/v1/tables/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Table deleted successfully"
}
```

**Response Error (409):**
```json
{
  "success": false,
  "message": "Cannot delete table that is currently occupied"
}
```

---

### 4.6 Users/Cashier Management (Owner)

#### Get All Cashiers

```http
GET /api/v1/users/cashiers
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page |
| `is_active` | bool | - | Filter by status |
| `include_deleted` | bool | false | Include soft deleted |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440000",
        "username": "cashier01",
        "full_name": "John Barista",
        "role": "cashier",
        "is_active": true,
        "created_at": "2026-01-15T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 10,
      "total_pages": 1
    }
  }
}
```

---

#### Create Cashier

```http
POST /api/v1/users/cashiers
```

**Akses:** Owner

**Request Body:**
```json
{
  "username": "cashier02",
  "password": "securepassword123",
  "full_name": "Jane Barista"
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Cashier created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "username": "cashier02",
    "full_name": "Jane Barista",
    "role": "cashier",
    "is_active": true,
    "created_at": "2026-03-20T10:00:00Z"
  }
}
```

**Response Error (409):**
```json
{
  "success": false,
  "message": "Username already exists"
}
```

---

#### Update Cashier

```http
PUT /api/v1/users/cashiers/:id
```

**Akses:** Owner

**Request Body:**
```json
{
  "full_name": "Jane Barista Senior",
  "is_active": true
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Cashier updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440001",
    "full_name": "Jane Barista Senior",
    "is_active": true,
    "updated_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Reset Cashier Password

```http
POST /api/v1/users/cashiers/:id/reset-password
```

**Akses:** Owner

**Request Body:**
```json
{
  "new_password": "newsecurepassword456"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Password reset successfully"
}
```

---

#### Delete Cashier (Soft Delete)

```http
DELETE /api/v1/users/cashiers/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Cashier deleted successfully"
}
```

**Response Error (409):**
```json
{
  "success": false,
  "message": "Cannot delete cashier with active shift"
}
```

---

### 4.7 Promos (Owner)

#### Get All Promos

```http
GET /api/v1/promos
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page |
| `is_active` | bool | - | Filter by status |
| `include_deleted` | bool | false | Include soft deleted |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440200",
        "name": "Happy Hour 20%",
        "description": "20% off all drinks",
        "type": "percentage",
        "value": 2000,
        "min_transaction": 5000000,
        "max_discount": 5000000,
        "started_at": "2026-03-01T00:00:00Z",
        "ended_at": "2026-03-31T23:59:59Z",
        "is_active": true,
        "created_at": "2026-02-25T08:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 5,
      "total_pages": 1
    }
  }
}
```

---

#### Get Active Promos (For Cashier)

```http
GET /api/v1/promos/active
```

**Akses:** Owner, Cashier

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440200",
        "name": "Happy Hour 20%",
        "description": "20% off all drinks",
        "type": "percentage",
        "value": 2000,
        "min_transaction": 5000000,
        "max_discount": 5000000,
        "started_at": "2026-03-01T00:00:00Z",
        "ended_at": "2026-03-31T23:59:59Z"
      }
    ]
  }
}
```

---

#### Create Promo

```http
POST /api/v1/promos
```

**Akses:** Owner

**Request Body:**
```json
{
  "name": "Weekend Special",
  "description": "Rp 10.000 off for weekend purchases",
  "type": "nominal",
  "value": 1000000,
  "min_transaction": 3000000,
  "started_at": "2026-03-22T00:00:00Z",
  "ended_at": "2026-03-24T23:59:59Z",
  "is_active": true
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Promo created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440201",
    "name": "Weekend Special",
    "type": "nominal",
    "value": 1000000,
    "created_at": "2026-03-20T10:00:00Z"
  }
}
```

---

#### Update Promo

```http
PUT /api/v1/promos/:id
```

**Akses:** Owner

**Request Body:**
```json
{
  "name": "Weekend Special Extended",
  "ended_at": "2026-03-31T23:59:59Z"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Promo updated successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440201",
    "name": "Weekend Special Extended",
    "updated_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Delete Promo (Soft Delete)

```http
DELETE /api/v1/promos/:id
```

**Akses:** Owner

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Promo deleted successfully"
}
```

---

### 4.8 Reports (Owner)

#### Get Revenue Report

```http
GET /api/v1/reports/revenue
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `period` | string | daily | daily/weekly/monthly |
| `start_date` | date | - | Start date (required if custom) |
| `end_date` | date | - | End date (required if custom) |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "period": "daily",
    "summary": {
      "total_revenue": 150000000,
      "total_transactions": 45,
      "average_transaction": 3333333
    },
    "daily_breakdown": [
      {
        "date": "2026-03-20",
        "revenue": 150000000,
        "transactions": 45
      }
    ],
    "trend": {
      "previous_period_revenue": 140000000,
      "growth_percentage": 7.14
    }
  }
}
```

---

#### Get Best Selling Products Report

```http
GET /api/v1/reports/best-sellers
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `limit` | int | 10 | Number of products |
| `start_date` | date | - | Start date |
| `end_date` | date | - | End date |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "period": {
      "start_date": "2026-03-01",
      "end_date": "2026-03-20"
    },
    "items": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440010",
        "product_name": "Caffe Latte",
        "total_quantity": 150,
        "total_revenue": 52500000,
        "rank": 1
      },
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440011",
        "product_name": "Cappuccino",
        "total_quantity": 120,
        "total_revenue": 42000000,
        "rank": 2
      }
    ]
  }
}
```

---

#### Get Cashier Performance Report

```http
GET /api/v1/reports/cashier-performance
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `start_date` | date | - | Start date |
| `end_date` | date | - | End date |
| `cashier_id` | string | - | Filter by specific cashier |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "period": {
      "start_date": "2026-03-01",
      "end_date": "2026-03-20"
    },
    "items": [
      {
        "cashier_id": "550e8400-e29b-41d4-a716-446655440000",
        "cashier_name": "John Barista",
        "total_transactions": 120,
        "total_revenue": 45000000,
        "average_transaction": 375000
      }
    ]
  }
}
```

---

#### Export Transactions to CSV

```http
GET /api/v1/reports/transactions/export
```

**Akses:** Owner

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `start_date` | date | - | Start date (required) |
| `end_date` | date | - | End date (required) |
| `status` | string | - | Filter by status |

**Response Sukses (200):**
```
Content-Type: text/csv
Content-Disposition: attachment; filename="transactions_2026-03-01_2026-03-20.csv"

Transaction Number,Date,Cashier,Table,Subtotal,Discount,Total,Status
TRX-20260320-0001,2026-03-20 10:00:00,John Barista,T01,35000,0,35000,paid
TRX-20260320-0002,2026-03-20 10:15:00,Jane Barista,T02,50000,5000,45000,paid
```

---

### 4.9 Shifts (Cashier)

#### Get Current Shift

```http
GET /api/v1/shifts/current
```

**Akses:** Cashier

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440300",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "opening_balance": 50000000,
    "status": "open",
    "opened_at": "2026-03-20T08:00:00Z"
  }
}
```

**Response Error (404):**
```json
{
  "success": false,
  "message": "No active shift found. Please open a shift first."
}
```

---

#### Open Shift

```http
POST /api/v1/shifts/open
```

**Akses:** Cashier

**Request Body:**
```json
{
  "opening_balance": 50000000
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Shift opened successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440300",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "opening_balance": 50000000,
    "status": "open",
    "opened_at": "2026-03-20T08:00:00Z"
  }
}
```

**Response Error (409):**
```json
{
  "success": false,
  "message": "You already have an active shift. Please close it first."
}
```

---

#### Close Shift

```http
POST /api/v1/shifts/:id/close
```

**Akses:** Cashier

**Request Body:**
```json
{
  "closing_balance": 125000000,
  "closing_notes": "All transactions completed successfully"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Shift closed successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440300",
    "opening_balance": 50000000,
    "closing_balance": 125000000,
    "expected_balance": 127000000,
    "balance_difference": -2000000,
    "status": "closed",
    "closed_at": "2026-03-20T20:00:00Z",
    "summary": {
      "total_transactions": 25,
      "total_sales": 77000000
    }
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Cannot close shift with no transactions"
}
```

---

#### Get Shift History

```http
GET /api/v1/shifts/history
```

**Akses:** Cashier

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `page` | int | 1 | Page number |
| `limit` | int | 20 | Items per page |
| `start_date` | date | - | Filter from date |
| `end_date` | date | - | Filter to date |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440300",
        "opening_balance": 50000000,
        "closing_balance": 125000000,
        "balance_difference": -2000000,
        "status": "closed",
        "opened_at": "2026-03-20T08:00:00Z",
        "closed_at": "2026-03-20T20:00:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "limit": 20,
      "total": 30,
      "total_pages": 2
    }
  }
}
```

---

### 4.10 Orders (Cashier)

#### Get Active Orders

```http
GET /api/v1/orders
```

**Akses:** Cashier

**Query Params:**
| Param | Type | Default | Keterangan |
|-------|------|---------|------------|
| `status` | string | pending | Filter by status |
| `table_id` | string | - | Filter by table |
| `shift_id` | string | - | Filter by shift (default: current shift) |

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "items": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440400",
        "transaction_number": "TRX-20260320-0001",
        "user_id": "550e8400-e29b-41d4-a716-446655440000",
        "table_id": "550e8400-e29b-41d4-a716-446655440050",
        "table_number": "T01",
        "is_takeaway": false,
        "subtotal": 3500000,
        "discount_amount": 0,
        "total": 3500000,
        "status": "pending",
        "items": [
          {
            "product_id": "550e8400-e29b-41d4-a716-446655440010",
            "product_name": "Caffe Latte",
            "quantity": 1,
            "price": 3500000
          }
        ],
        "created_at": "2026-03-20T10:00:00Z"
      }
    ]
  }
}
```

---

#### Create Order

```http
POST /api/v1/orders
```

**Akses:** Cashier

**Request Body:**
```json
{
  "table_id": "550e8400-e29b-41d4-a716-446655440050",
  "is_takeaway": false,
  "items": [
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440010",
      "quantity": 2,
      "notes": "Less sugar"
    },
    {
      "product_id": "550e8400-e29b-41d4-a716-446655440011",
      "quantity": 1,
      "notes": "Extra hot"
    }
  ],
  "customer_notes": "For here"
}
```

**Response Sukses (201):**
```json
{
  "success": true,
  "message": "Order created successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "transaction_number": "TRX-20260320-0001",
    "subtotal": 7000000,
    "discount_amount": 0,
    "total": 7000000,
    "status": "pending",
    "items": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440010",
        "product_name": "Caffe Latte",
        "quantity": 2,
        "price": 3500000,
        "subtotal": 7000000
      }
    ],
    "created_at": "2026-03-20T10:00:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Cannot create order. Please open a shift first."
}
```

**Response Error (422):**
```json
{
  "success": false,
  "message": "Validation failed",
  "errors": {
    "items": "At least one item is required"
  }
}
```

---

#### Get Order by ID

```http
GET /api/v1/orders/:id
```

**Akses:** Cashier

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "transaction_number": "TRX-20260320-0001",
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "table_id": "550e8400-e29b-41d4-a716-446655440050",
    "table_number": "T01",
    "is_takeaway": false,
    "subtotal": 3500000,
    "discount_amount": 0,
    "total": 3500000,
    "status": "pending",
    "items": [
      {
        "product_id": "550e8400-e29b-41d4-a716-446655440010",
        "product_name": "Caffe Latte",
        "quantity": 1,
        "price": 3500000,
        "notes": "Less sugar"
      }
    ],
    "created_at": "2026-03-20T10:00:00Z"
  }
}
```

---

#### Apply Promo to Order

```http
POST /api/v1/orders/:id/promo
```

**Akses:** Cashier

**Request Body:**
```json
{
  "promo_id": "550e8400-e29b-41d4-a716-446655440200"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Promo applied successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "subtotal": 3500000,
    "discount_amount": 700000,
    "discount_type": "percentage",
    "promo_id": "550e8400-e29b-41d4-a716-446655440200",
    "promo_name": "Happy Hour 20%",
    "total": 2800000
  }
}
```

**Response Error (422):**
```json
{
  "success": false,
  "message": "Promo is not valid",
  "errors": {
    "promo_id": "Promo has expired"
  }
}
```

---

#### Remove Promo from Order

```http
DELETE /api/v1/orders/:id/promo
```

**Akses:** Cashier

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Promo removed successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "subtotal": 3500000,
    "discount_amount": 0,
    "discount_type": "none",
    "total": 3500000
  }
}
```

---

#### Cancel Order

```http
POST /api/v1/orders/:id/cancel
```

**Akses:** Cashier

**Request Body:**
```json
{
  "reason": "Customer changed mind"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Order cancelled successfully",
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440400",
    "status": "cancelled",
    "cancelled_at": "2026-03-20T10:30:00Z"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Cannot cancel order that has been paid"
}
```

---

### 4.11 Payments

#### Initiate Payment (Midtrans Snap)

```http
POST /api/v1/orders/:id/pay
```

**Akses:** Cashier

**Request Body:**
```json
{
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "customer_phone": "081234567890"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Payment initiated successfully",
  "data": {
    "transaction_id": "550e8400-e29b-41d4-a716-446655440400",
    "midtrans_order_id": "TRX-20260320-0001-PAY",
    "snap_token": "66e4fa1e-fdac-4329-9091-7bec49c09120",
    "snap_redirect_url": "https://app.sandbox.midtrans.com/snap/v2/vtweb/66e4fa1e-fdac-4329-9091-7bec49c09120",
    "amount": 3500000,
    "status": "pending"
  }
}
```

**Response Error (400):**
```json
{
  "success": false,
  "message": "Order status must be pending to initiate payment"
}
```

---

#### Get Payment Status

```http
GET /api/v1/payments/:transaction_id
```

**Akses:** Cashier

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440500",
    "transaction_id": "550e8400-e29b-41d4-a716-446655440400",
    "midtrans_order_id": "TRX-20260320-0001-PAY",
    "midtrans_transaction_id": "b54a14d7-60a7-4ca2-ae35-d0802d063ae5",
    "payment_method": "gopay",
    "payment_amount": 3500000,
    "payment_status": "success",
    "paid_at": "2026-03-20T10:15:00Z",
    "created_at": "2026-03-20T10:05:00Z"
  }
}
```

---

#### Midtrans Webhook Handler

```http
POST /api/v1/payments/webhook
```

**Akses:** Public (Midtrans only)

**Headers:**
```http
Content-Type: application/json
X-Midtrans-Signature: <sha512_signature>
```

**Request Body:**
```json
{
  "order_id": "TRX-20260320-0001-PAY",
  "transaction_id": "b54a14d7-60a7-4ca2-ae35-d0802d063ae5",
  "transaction_time": "2026-03-20 10:15:00",
  "transaction_status": "capture",
  "fraud_status": "accept",
  "payment_type": "gopay",
  "gross_amount": "35000.00"
}
```

**Response Sukses (200):**
```json
{
  "success": true,
  "message": "Webhook received"
}
```

---

#### Get Transaction Receipt

```http
GET /api/v1/payments/:transaction_id/receipt
```

**Akses:** Cashier

**Response Sukses (200):**
```json
{
  "success": true,
  "data": {
    "transaction_number": "TRX-20260320-0001",
    "date": "2026-03-20T10:15:00Z",
    "cashier_name": "John Barista",
    "table_number": "T01",
    "items": [
      {
        "name": "Caffe Latte",
        "quantity": 1,
        "price": 35000,
        "subtotal": 35000
      }
    ],
    "subtotal": 35000,
    "discount": 0,
    "total": 35000,
    "payment_method": "GoPay",
    "payment_status": "Paid"
  }
}
```

---

## 5. Error Handling Best Practices

### Client-Side Error Handling

1. **401 Unauthorized:** Redirect to login page
2. **403 Forbidden:** Show "Access Denied" message
3. **404 Not Found:** Show "Page Not Found" page
4. **422 Validation Error:** Display field-specific error messages
5. **500 Server Error:** Show generic error message, log error

### Retry Logic

- **Idempotent operations (GET, PUT, DELETE):** Retry up to 3 times with exponential backoff
- **Non-idempotent operations (POST):** Do not retry automatically

---

## Version History

| Version | Date | Author | Description |
|---------|------|--------|-------------|
| 1.0 | 20 Maret 2026 | API Team | Initial API contract |
