# Package `jwt`

This package provides functions to create, parse, and verify JWT (JSON Web Tokens) using the `github.com/golang-jwt/jwt/v5` library.

## Types

### 1. Payload
A struct representing the payload claims for a JWT token.

#### Fields:
- `PlayerID string`: The player ID associated with the token.
- `ServerGroup string`: The server group associated with the token.

## Functions

### 1. CreateToken
Creates and signs a new JWT token using the provided key and payload.

#### Parameters:
- `key []byte`: Secret key used for signing the token.
- `payload Payload`: The payload containing player ID and server group information.

#### Returns:
- A string representing the signed JWT token.
- An error if any occurs during token creation or signing.

#### Example Usage:
```go
key := []byte("my_secret_key")
playerID := "12345"
serverGroup := "groupA"

payload := Payload{
    PlayerID:    playerID,
    ServerGroup: serverGroup,
}

tokenString, err := CreateToken(key, payload)
if err != nil {
    log.Fatalf("Failed to create token: %v", err)
}
fmt.Println(tokenString)
```

### 2. CheckSignature
Verifies the signature of a JWT token using the provided key.

#### Parameters:
- `tokenString string`: The JWT token to verify.
- `key []byte`: Secret key used for signing the token.

#### Returns:
- A boolean indicating whether the token's signature is valid.
- An error if any occurs during verification.

#### Example Usage:
```go
tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
key := []byte("my_secret_key")

valid, err := CheckSignature(tokenString, key)
if err != nil {
    log.Fatalf("Failed to check token signature: %v", err)
}
fmt.Println(valid) // true or false
```

### 3. ParseUnverified
Parses a JWT token without verifying its signature.

#### Parameters:
- `tokenString string`: The JWT token to parse.

#### Returns:
- A `Payload` struct containing the parsed claims.
- An error if any occurs during parsing.

#### Example Usage:
```go
tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."

payload, err := ParseUnverified(tokenString)
if err != nil {
    log.Fatalf("Failed to parse unverified token: %v", err)
}
fmt.Println(payload.PlayerID) // "player_id_value"
```

### 4. ParseToken
Parses and verifies a JWT token using the provided key.

#### Parameters:
- `tokenString string`: The JWT token to parse.
- `key []byte`: Secret key used for signing the token.

#### Returns:
- A `Payload` struct containing the parsed claims.
- An error if any occurs during parsing or verification.

#### Example Usage:
```go
tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
key := []byte("my_secret_key")

payload, err := ParseToken(tokenString, key)
if err != nil {
    log.Fatalf("Failed to parse token: %v", err)
}
fmt.Println(payload.PlayerID) // "player_id_value"
```

## Notes:

- **CreateToken**: This function creates a new JWT token with the specified payload and signs it using HMAC SHA256.
  
- **CheckSignature**: It verifies if the provided JWT token is signed correctly without extracting any claims from it.

- **ParseUnverified**: Parses a JWT token but does not verify its signature. Use this for scenarios where you only need to extract claims without verifying them.

- **ParseToken**: Parses and verifies a JWT token using the specified key, ensuring that the token's signature matches before returning the payload claims.
