# Refresh Token API Documentation

This document describes the refresh token functionality implemented in the PandoraGym Go API. The refresh token system provides secure, long-lived authentication that allows users to obtain new access tokens without re-entering their credentials.

## 🔐 Overview

The refresh token system implements the following security features:

- **Token Rotation**: Each refresh generates a new refresh token and revokes the old one
- **Long-lived Sessions**: Refresh tokens expire after 30 days (configurable)
- **Device Tracking**: Optional device information and IP address tracking
- **Secure Storage**: Tokens are stored securely in the database with proper indexing
- **Automatic Cleanup**: Expired tokens are automatically cleaned up

## 📊 Database Schema

The refresh token system adds a new `refresh_tokens` table:

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL UNIQUE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP NULL,
    device_info TEXT NULL,
    ip_address INET NULL
);
```

## 🚀 API Endpoints

### 1. Authentication (Updated)

**Endpoint:** `POST /auth/session`

**Description:** Authenticates a user and returns both access and refresh tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "device_info": "Optional device information",
  "ip_address": "Optional IP address (auto-detected if not provided)"
}
```

**Response:**
```json
{
  "user": {
    "id": "uuid",
    "name": "User Name",
    "email": "user@example.com",
    "role": "STUDENT",
    // ... other user fields
  },
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "base64-encoded-refresh-token",
  "expires_in": 604800000000000
}
```

### 2. Refresh Token

**Endpoint:** `POST /auth/refresh`

**Description:** Exchanges a valid refresh token for a new access token and refresh token.

**Request Body:**
```json
{
  "refresh_token": "base64-encoded-refresh-token"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "new-base64-encoded-refresh-token",
  "expires_in": 604800
}
```

**Error Responses:**
- `400 Bad Request`: Invalid request format
- `401 Unauthorized`: Invalid or expired refresh token

### 3. Revoke Token

**Endpoint:** `POST /auth/revoke`

**Description:** Revokes a specific refresh token, effectively logging out that session.

**Request Body:**
```json
{
  "refresh_token": "base64-encoded-refresh-token"
}
```

**Response:**
```json
{
  "message": "Token revoked successfully"
}
```

## 🔧 Implementation Details

### Token Generation

- **Access Tokens**: JWT tokens with 7-day expiration
- **Refresh Tokens**: 32-byte random tokens encoded in base64, 30-day expiration
- **Token Rotation**: Each refresh operation generates new tokens and revokes old ones

### Security Features

1. **Automatic IP Detection**: If not provided, IP address is extracted from request
2. **Device Tracking**: Optional device information for session management
3. **Token Revocation**: Immediate token invalidation capability
4. **Cleanup Process**: Automatic removal of expired tokens

### Database Operations

The system provides the following database operations:

- `CreateRefreshToken`: Create a new refresh token
- `GetRefreshToken`: Retrieve and validate a refresh token
- `GetRefreshTokensByUserID`: Get all active tokens for a user
- `RevokeRefreshToken`: Revoke a specific token
- `RevokeAllUserRefreshTokens`: Revoke all tokens for a user
- `CleanupExpiredRefreshTokens`: Remove expired tokens
- `UpdateRefreshTokenLastUsed`: Update token usage timestamp

## 📝 Usage Examples

### JavaScript/Frontend Example

```javascript
class AuthService {
  constructor() {
    this.baseURL = 'http://localhost:3333';
    this.accessToken = localStorage.getItem('access_token');
    this.refreshToken = localStorage.getItem('refresh_token');
  }

  async login(email, password) {
    const response = await fetch(`${this.baseURL}/auth/session`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        email,
        password,
        device_info: navigator.userAgent,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      this.accessToken = data.token;
      this.refreshToken = data.refresh_token;
      
      localStorage.setItem('access_token', this.accessToken);
      localStorage.setItem('refresh_token', this.refreshToken);
      
      return data;
    }
    
    throw new Error('Login failed');
  }

  async refreshAccessToken() {
    if (!this.refreshToken) {
      throw new Error('No refresh token available');
    }

    const response = await fetch(`${this.baseURL}/auth/refresh`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        refresh_token: this.refreshToken,
      }),
    });

    if (response.ok) {
      const data = await response.json();
      this.accessToken = data.access_token;
      this.refreshToken = data.refresh_token;
      
      localStorage.setItem('access_token', this.accessToken);
      localStorage.setItem('refresh_token', this.refreshToken);
      
      return data;
    }
    
    throw new Error('Token refresh failed');
  }

  async makeAuthenticatedRequest(url, options = {}) {
    const headers = {
      'Authorization': `Bearer ${this.accessToken}`,
      'Content-Type': 'application/json',
      ...options.headers,
    };

    let response = await fetch(url, {
      ...options,
      headers,
    });

    // If token expired, try to refresh
    if (response.status === 401) {
      try {
        await this.refreshAccessToken();
        
        // Retry the request with new token
        headers['Authorization'] = `Bearer ${this.accessToken}`;
        response = await fetch(url, {
          ...options,
          headers,
        });
      } catch (error) {
        // Refresh failed, redirect to login
        this.logout();
        throw new Error('Authentication failed');
      }
    }

    return response;
  }

  async logout() {
    if (this.refreshToken) {
      try {
        await fetch(`${this.baseURL}/auth/revoke`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            refresh_token: this.refreshToken,
          }),
        });
      } catch (error) {
        console.error('Failed to revoke token:', error);
      }
    }

    this.accessToken = null;
    this.refreshToken = null;
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
  }
}
```

### cURL Examples

```bash
# 1. Login and get tokens
curl -X POST http://localhost:3333/auth/session \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "device_info": "curl/7.68.0"
  }'

# 2. Refresh access token
curl -X POST http://localhost:3333/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token-here"
  }'

# 3. Make authenticated request
curl -X GET http://localhost:3333/api/users/profile \
  -H "Authorization: Bearer your-access-token-here"

# 4. Revoke refresh token
curl -X POST http://localhost:3333/auth/revoke \
  -H "Content-Type: application/json" \
  -d '{
    "refresh_token": "your-refresh-token-here"
  }'
```

## 🛡️ Security Considerations

1. **Token Storage**: Store refresh tokens securely (HttpOnly cookies recommended for web apps)
2. **Token Rotation**: Always rotate refresh tokens to prevent replay attacks
3. **Expiration**: Set appropriate expiration times based on security requirements
4. **Revocation**: Implement proper token revocation for logout scenarios
5. **Rate Limiting**: Consider implementing rate limiting on refresh endpoints
6. **HTTPS**: Always use HTTPS in production to protect tokens in transit

## 🔄 Token Lifecycle

1. **Login**: User provides credentials, receives access + refresh token
2. **API Calls**: Use access token for authenticated requests
3. **Token Expiry**: When access token expires, use refresh token to get new tokens
4. **Token Rotation**: Old refresh token is revoked, new tokens are issued
5. **Logout**: Refresh token is revoked, ending the session

## 🧹 Maintenance

### Cleanup Process

Run the cleanup process periodically to remove expired tokens:

```go
// In your application startup or cron job
err := authService.CleanupExpiredTokens(ctx)
if err != nil {
    log.Error("Failed to cleanup expired tokens", "error", err)
}
```

### Monitoring

Monitor the following metrics:
- Active refresh token count per user
- Token refresh frequency
- Failed refresh attempts
- Token cleanup statistics

## 🚨 Error Handling

The API returns appropriate HTTP status codes:

- `200 OK`: Successful operation
- `400 Bad Request`: Invalid request format or missing required fields
- `401 Unauthorized`: Invalid or expired tokens
- `500 Internal Server Error`: Server-side errors

Error responses include descriptive messages:

```json
{
  "error": "Invalid or expired refresh token"
}
```

## 📋 Testing

Use the provided test script to verify functionality:

```bash
./test_refresh_token.sh
```

This script tests:
- User registration
- Authentication with token generation
- Access token usage
- Token refresh
- Token revocation
- Revoked token rejection

## 🔮 Future Enhancements

Potential improvements to consider:

1. **Multiple Device Management**: UI for managing active sessions
2. **Push Notifications**: Notify users of new login sessions
3. **Geolocation Tracking**: Track login locations for security
4. **Token Fingerprinting**: Additional security through device fingerprinting
5. **Refresh Token Families**: Implement token families for better security
6. **Rate Limiting**: Implement rate limiting on refresh endpoints

---

This refresh token implementation provides a robust, secure authentication system that enhances user experience while maintaining strong security practices.
