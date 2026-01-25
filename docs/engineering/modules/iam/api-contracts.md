# IAM Module API Contracts

This document specifies the API contracts for the Identity and Access Management (IAM) module.

---

## 1. Register User

- **Endpoint:** `POST /auth/register`
- **Description:** Registers a new user in the system.
- **Request Body:**
  ```json
  {
    "email": "new.user@example.com",
    "password": "a-strong-password"
  }
  ```
- **Success Response (201 Created):**
  ```json
  {
    "id": "user-uuid-123",
    "email": "new.user@example.com"
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: If email is invalid or password is too weak.
  - `409 Conflict`: If a user with that email already exists.

---

## 2. Login User

- **Endpoint:** `POST /auth/login`
- **Description:** Authenticates a user and returns JWT access and refresh tokens.
- **Request Body:**
  ```json
  {
    "email": "registered.user@example.com",
    "password": "user-password"
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "access_token": "ey...",
    "refresh_token": "ey...",
    "user": {
        "id": "user-uuid-123",
        "email": "registered.user@example.com"
    }
  }
  ```
- **Error Responses:**
  - `400 Bad Request`: If request format is invalid.
  - `401 Unauthorized`: If credentials are incorrect.

---

## 3. Refresh Token

- **Endpoint:** `POST /auth/refresh`
- **Description:** Issues a new access token using a valid refresh token.
- **Request Body:**
  ```json
  {
    "refresh_token": "a-valid-refresh-token"
  }
  ```
- **Success Response (200 OK):**
  ```json
  {
    "access_token": "ey..."
  }
  ```
- **Error Responses:**
  - `401 Unauthorized`: If the refresh token is invalid or expired.

---

## 4. Logout User

- **Endpoint:** `POST /auth/logout`
- **Description:** Invalidates the user's session (e.g., by blacklisting the token).
- **Request Body:** (Requires authentication)
  ```json
  {}
  ```
- **Success Response (204 No Content):**
- **Error Responses:**
  - `401 Unauthorized`: If no valid access token is provided.
