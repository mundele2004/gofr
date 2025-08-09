# LDAP Authentication Demo

This is a Go application that demonstrates LDAP authentication with JWT token issuance. The application addresses all linting issues and reviewer feedback.

## Features

- ✅ **Interface-based JWT issuance** - JWT tokens are issued through a pluggable interface
- ✅ **Separate error handling** - User not found (401) vs multiple users found (500)  
- ✅ **Comprehensive user information** - Retrieves full user details from LDAP
- ✅ **All linting issues fixed** - err113, gci, gocritic, gocyclo, revive, whitespace

## Project Structure

```
ldap-auth-demo/
├── main.go                                 # HTTP server with login/userinfo endpoints
├── internal/auth/
│   ├── interfaces/interfaces.go           # Shared interfaces (UserInfo, TokenIssuer)
│   ├── ldapauth/ldapauth.go               # LDAP authentication logic
│   └── jwt/jwt.go                         # JWT token issuer implementation
├── go.mod                                 # Go module definition
└── README.md                              # This file
```

## Building and Running

1. **Build the application:**
   ```bash
   cd ldap-auth-demo
   go build -o ldap-auth-server
   ```

2. **Run the server:**
   ```bash
   ./ldap-auth-server
   ```
   The server will start on port 8080.

## API Endpoints

### POST /login
Authenticate user and return JWT token.

**Request:**
```json
{
  "username": "john.doe",
  "password": "secretpassword"
}
```

**Success Response (200):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error Responses:**
- `401 Unauthorized` - Invalid credentials or user not found
- `500 Internal Server Error` - Multiple users found (LDAP misconfiguration)

### GET /userinfo?username=john.doe
Retrieve user information from LDAP (without authentication).

**Success Response (200):**
```json
{
  "dn": "uid=john.doe,dc=example,dc=com",
  "username": "john.doe",
  "email": "john.doe@example.com",
  "full_name": "John Doe",
  "groups": ["cn=developers,dc=example,dc=com"],
  "department": "Engineering",
  "attributes": {
    "telephoneNumber": "+1-555-0123"
  }
}
```

## Configuration

Update the LDAP configuration in `main.go`:

```go
cfg := ldapauth.Config{
    Addr:         "localhost:389",                    // LDAP server
    BaseDN:       "dc=example,dc=com",               // Search base
    BindUserDN:   "cn=admin,dc=example,dc=com",      // Service account DN
    BindPassword: "admin",                           // Service account password
    TokenIssuer:  jwtIssuer,                        // JWT issuer implementation
    UserAttributes: []string{                       // LDAP attributes to retrieve
        "dn", "uid", "mail", "cn", "ou", "memberOf", "telephoneNumber",
    },
    UsernameAttr:   "uid",                          // Username attribute
    EmailAttr:      "mail",                         // Email attribute
    FullNameAttr:   "cn",                           // Full name attribute
    DepartmentAttr: "ou",                           // Department attribute
}
```

## Testing Without LDAP Server

To test the HTTP endpoints without an LDAP server:

1. **Test server startup:**
   ```bash
   curl http://localhost:8080/login -d '{"username":"test","password":"test"}' -H "Content-Type: application/json"
   ```
   
   Expected: Connection error (since no LDAP server is running)

2. **Test with mock LDAP (for development):**
   The code includes interfaces that make it easy to inject mock implementations for testing.

## Dependencies

- `github.com/go-ldap/ldap/v3` - LDAP client library
- `github.com/golang-jwt/jwt/v4` - JWT library

## Architecture Decisions

1. **Interface-based design** - TokenIssuer interface allows different token implementations
2. **Separate concerns** - LDAP authentication is separate from JWT issuance  
3. **Error differentiation** - Different HTTP status codes for different error types
4. **User information** - Comprehensive user data retrieval from LDAP directories
5. **No circular imports** - Shared interfaces package prevents circular dependencies