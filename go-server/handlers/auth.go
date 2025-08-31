package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
	Name     string `json:"name"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Name     string `json:"name" binding:"required"`
}

func handleLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// Query user
	var user User
	err = db.QueryRow("SELECT id, email, password, role, name FROM users WHERE email = $1", req.Email).Scan(
		&user.ID, &user.Email, &user.Password, &user.Role, &user.Name)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Token generation failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
		"user": gin.H{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
			"name":  user.Name,
		},
	})
}

func handleRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Connect to database
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection failed"})
		return
	}
	defer db.Close()

	// Check if user already exists
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", req.Email).Scan(&exists)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Password hashing failed"})
		return
	}

	// Insert new user
	var userID int
	err = db.QueryRow(
		"INSERT INTO users (email, password, name, role) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Email, string(hashedPassword), req.Name, "user",
	).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User creation failed"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
		"user_id": userID,
	})
}
```

```go:backend/go-server/handlers/upload.go
package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleMovieImageUpload(c *gin.Context) {
	// Get user from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Check if user is admin
	userRole, _ := c.Get("user_role")
	if userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
		return
	}

	// Get the uploaded file
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
		return
	}

	// Validate file type
	allowedTypes := []string{".jpg", ".jpeg", ".png", ".webp"}
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, t := range allowedTypes {
		if t == ext {
			allowed = true
			break
		}
	}
	if !allowed {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
		return
	}

	// Generate unique filename
	filename := uuid.New().String() + ext
	uploadDir := "./uploads/movies"
	
	// Create directory if it doesn't exist
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload directory"})
		return
	}

	filepath := filepath.Join(uploadDir, filename)
	
	// Save file
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Return file URL
	fileURL := fmt.Sprintf("/uploads/movies/%s", filename)
	c.JSON(http.StatusOK, gin.H{
		"url": fileURL,
		"filename": filename,
	})
}
```

```go:backend/go-server/handlers/payment.go
package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type ChapaPaymentRequest struct {
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	TxRef       string `json:"tx_ref"`
	CallbackURL string `json:"callback_url"`
	ReturnURL   string `json:"return_url"`
}

type ChapaResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		CheckoutURL string `json:"checkout_url"`
		Reference   string `json:"reference"`
	} `json:"data"`
}

func handleInitiatePayment(c *gin.Context) {
	var req struct {
		Amount   string `json:"amount" binding:"required"`
		Currency string `json:"currency" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		TxRef    string `json:"tx_ref" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user from context
	userID, _ := c.Get("user_id")
	userName, _ := c.Get("user_name")

	// Prepare Chapa payment request
	chapaReq := ChapaPaymentRequest{
		Amount:      req.Amount,
		Currency:    req.Currency,
		Email:       req.Email,
		FirstName:   userName.(string),
		LastName:    "",
		TxRef:       req.TxRef,
		CallbackURL: "http://localhost:8080/api/payment/verify",
		ReturnURL:   "http://localhost:3000/payment/success",
	}

	// Convert to JSON
	jsonData, err := json.Marshal(chapaReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to prepare payment request"})
		return
	}

	// Make request to Chapa API
	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://api.chapa.co/v1/transaction/initialize", strings.NewReader(string(jsonData)))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment request"})
		return
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHAPA_SECRET_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Payment service unavailable"})
		return
	}
	defer resp.Body.Close()

	var chapaResp ChapaResponse
	if err := json.NewDecoder(resp.Body).Decode(&chapaResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse payment response"})
		return
	}

	if chapaResp.Status != "success" {
		c.JSON(http.StatusBadRequest, gin.H{"error": chapaResp.Message})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"checkout_url": chapaResp.Data.CheckoutURL,
		"reference":    chapaResp.Data.Reference,
	})
}
```

```go:backend/go-server/middleware/auth.go
package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token required"})
			c.Abort()
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		// Set user info in context
		c.Set("user_id", int(claims["user_id"].(float64)))
		c.Set("user_email", claims["email"].(string))
		c.Set("user_role", claims["role"].(string))
		c.Set("user_name", claims["name"].(string))

		c.Next()
	}
}
```

```go:backend/go-server/go.mod
module cinema-tickets

go 1.21

require (
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang-jwt/jwt/v5 v5.0.0
	github.com/google/uuid v1.3.0
	github.com/joho/godotenv v1.4.0
	github.com/lib/pq v1.10.9
	golang.org/x/crypto v0.12.0
)

require (
	github.com/bytedance/sonic v1.9.1 // indirect
	github.com/chenzhuoyu/base64x v0.0.0-20221115062448-fe3a3abad311 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.14.0 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.2.4 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pelletier/go-toml/v2 v2.0.8 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.11 // indirect
	golang.org/x/arch v0.3.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.11.0 // indirect
	golang.org/x/text v0.12.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
```

## 2. Database Schema (PostgreSQL)

```sql:backend/database/schema.sql
-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'admin')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Directors table
CREATE TABLE directors (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Stars (Actors/Actresses) table
CREATE TABLE stars (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    bio TEXT,
    image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Genres table
CREATE TABLE genres (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Movies table
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    duration INTEGER NOT NULL, -- in minutes
    director_id INTEGER REFERENCES directors(id),
    featured_image_url VARCHAR(500),
    trailer_url VARCHAR(500),
    release_date DATE,
    rating DECIMAL(3,1) DEFAULT 0.0,
    total_ratings INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Movie images table
CREATE TABLE movie_images (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    image_url VARCHAR(500) NOT NULL,
    is_featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Movie genres junction table
CREATE TABLE movie_genres (
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    genre_id INTEGER REFERENCES genres(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, genre_id)
);

-- Movie stars junction table
CREATE TABLE movie_stars (
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    star_id INTEGER REFERENCES stars(id) ON DELETE CASCADE,
    PRIMARY KEY (movie_id, star_id)
);

-- Cinemas table
CREATE TABLE cinemas (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    address TEXT,
    city VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Cinema halls table
CREATE TABLE cinema_halls (
    id SERIAL PRIMARY KEY,
    cinema_id INTEGER REFERENCES cinemas(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    capacity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Schedules table
CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    cinema_hall_id INTEGER REFERENCES cinema_halls(id) ON DELETE CASCADE,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    available_seats INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tickets table
CREATE TABLE tickets (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    schedule_id INTEGER REFERENCES schedules(id) ON DELETE CASCADE,
    seat_number VARCHAR(10) NOT NULL,
    price DECIMAL(10,2) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' CHECK (status IN ('pending', 'confirmed', 'cancelled')),
    payment_reference VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- User ratings table
CREATE TABLE user_ratings (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    movie_id INTEGER REFERENCES movies(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating >= 1 AND rating <= 5),
    review TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, movie_id)
);

-- User bookmarks table
CREATE TABLE user_bookmarks (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    schedule_id INTEGER REFERENCES schedules(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, schedule_id)
);

-- Triggers for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_directors_updated_at BEFORE UPDATE ON directors FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_stars_updated_at BEFORE UPDATE ON stars FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_movies_updated_at BEFORE UPDATE ON movies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_schedules_updated_at BEFORE UPDATE ON schedules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_tickets_updated_at BEFORE UPDATE ON tickets FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_ratings_updated_at BEFORE UPDATE ON user_ratings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Function to update movie rating
CREATE OR REPLACE FUNCTION update_movie_rating()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE movies 
    SET rating = (
        SELECT AVG(rating)::DECIMAL(3,1)
        FROM user_ratings 
        WHERE movie_id = NEW.movie_id
    ),
    total_ratings = (
        SELECT COUNT(*)
        FROM user_ratings 
        WHERE movie_id = NEW.movie_id
    )
    WHERE id = NEW.movie_id;
    
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_movie_rating_trigger 
AFTER INSERT OR UPDATE OR DELETE ON user_ratings 
FOR EACH ROW EXECUTE FUNCTION update_movie_rating();

-- Function to update available seats
CREATE OR REPLACE FUNCTION update_available_seats()
RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE schedules 
        SET available_seats = available_seats - 1
        WHERE id = NEW.schedule_id;
        RETURN NEW;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE schedules 
        SET available_seats = available_seats + 1
        WHERE id = OLD.schedule_id;
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_available_seats_trigger 
AFTER INSERT OR DELETE ON tickets 
FOR EACH ROW EXECUTE FUNCTION update_available_seats();

-- Insert default admin user
INSERT INTO users (email, password, name, role) VALUES 
('admin@tickent.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin User', 'admin');

-- Insert sample genres
INSERT INTO genres (name, description) VALUES 
('Action', 'High-energy films with exciting sequences'),
('Comedy', 'Humorous and entertaining films'),
('Drama', 'Serious, character-driven stories'),
('Horror', 'Scary and suspenseful films'),
('Romance', 'Love stories and romantic comedies'),
('Sci-Fi', 'Science fiction and futuristic films'),
('Thriller', 'Suspenseful and exciting films'),
('Documentary', 'Non-fiction educational films');

-- Insert sample cinemas
INSERT INTO cinemas (name, address, city, phone) VALUES 
('CineMax Downtown', '123 Main Street, Downtown', 'New York', '+1-555-0101'),
('CineMax Uptown', '456 Park Avenue, Uptown', 'New York', '+1-555-0102'),
('CineMax Westside', '789 Broadway, Westside', 'Los Angeles', '+1-555-0103');
```

## 3. Hasura Configuration

```yaml:backend/my-project/hasura.yaml
version: v3
endpoint: http://localhost:8080
metadata_directory: metadata
migrations_directory: migrations
seeds_directory: seeds
actions:
  kind: synchronous
  handler_webhook_baseurl: http://localhost:8080
```

```yaml:backend/my-project/metadata/databases/databases.yaml
- name: default
  kind: postgres
  configuration:
    connection_info:
      database_url:
        from_env: HASURA_GRAPHQL_DATABASE_URL
      pool_settings:
        max_connections: 50
        idle_timeout: 180
        retries: 1
        pool_timeout: 360
        connection_lifetime: 600
  tables: "!include default/tables/tables.yaml"
  functions: "!include default/functions/functions.yaml"
```

```yaml:backend/my-project/metadata/databases/default/tables/public_users.yaml
table:
  name: users
  schema: public
select_permissions:
- role: user
  permission:
    columns:
    - id
    - email
    - name
    - role
    - created_at
    - updated_at
    filter:
      id:
        _eq: X-Hasura-User-Id
insert_permissions:
- role: user
  permission:
    check:
      role:
        _eq: user
    columns:
    - email
    - password
    - name
    - role
update_permissions:
- role: user
  permission:
    columns:
    - name
    - email
    filter:
      id:
        _eq: X-Hasura-User-Id
- role: admin
  permission:
    columns:
    - id
    - email
    - name
    - role
    - created_at
    - updated_at
delete_permissions:
- role: admin
  permission: {}
```

```yaml:backend/my-project/metadata/databases/default/tables/public_movies.yaml
table:
  name: movies
  schema: public
select_permissions:
- role: user
  permission:
    columns:
    - id
    - title
    - description
    - duration
    - director_id
    - featured_image_url
    - trailer_url
    - release_date
    - rating
    - total_ratings
    - created_at
    - updated_at
insert_permissions:
- role: admin
  permission:
    columns:
    - title
    - description
    - duration
    - director_id
    - featured_image_url
    - trailer_url
    - release_date
update_permissions:
- role: admin
  permission:
    columns:
    - title
    - description
    - duration
    - director_id
    - featured_image_url
    - trailer_url
    - release_date
delete_permissions:
- role: admin
  permission: {}
```

```yaml:backend/my-project/metadata/actions.yaml
- name: login
  definition:
    kind: synchronous
    handler: http://localhost:8080/api/auth/login
    forward_client_headers: true
    headers:
    - name: Content-Type
      value: application/json
  type: mutation
  arguments:
  - name: email
    type: String!
  - name: password
    type: String!
  output_type: LoginResponse

- name: register
  definition:
    kind: synchronous
    handler: http://localhost:8080/api/auth/register
    forward_client_headers: true
    headers:
    - name: Content-Type
      value: application/json
  type: mutation
  arguments:
  - name: email
    type: String!
  - name: password
    type: String!
  - name: name
    type: String!
  output_type: RegisterResponse

- name: upload_movie_image
  definition:
    kind: synchronous
    handler: http://localhost:8080/api/upload/movie-images
    forward_client_headers: true
  type: mutation
  arguments:
  - name: image
    type: Upload!
  output_type: UploadResponse

- name: initiate_payment
  definition:
    kind: synchronous
    handler: http://localhost:8080/api/payment/initiate
    forward_client_headers: true
    headers:
    - name: Content-Type
      value: application/json
  type: mutation
  arguments:
  - name: amount
    type: String!
  - name: currency
    type: String!
  - name: email
    type: String!
  - name: tx_ref
    type: String!
  output_type: PaymentResponse
```

```yaml:backend/my-project/metadata/custom_types.yaml
input_objects:
- name: LoginInput
  fields:
  - name: email
    type: String!
  - name: password
    type: String!

- name: RegisterInput
  fields:
  - name: email
    type: String!
  - name: password
    type: String!
  - name: name
    type: String!

- name: PaymentInput
  fields:
  - name: amount
    type: String!
  - name: currency
    type: String!
  - name: email
    type: String!
  - name: tx_ref
    type: String!

objects:
- name: LoginResponse
  fields:
  - name: token
    type: String!
  - name: user
    type: User!

- name: RegisterResponse
  fields:
  - name: message
    type: String!
  - name: user_id
    type: Int!

- name: UploadResponse
  fields:
  - name: url
    type: String!
  - name: filename
    type: String!

- name: PaymentResponse
  fields:
  - name: checkout_url
    type: String!
  - name: reference
    type: String!
```

## 4. Docker Configuration

```yaml:docker-compose.yml
version: '3.8'

services:
  postgres:
    image: postgres:15
    container_name: cinema_postgres
    restart: always
    environment:
      POSTGRES_DB: cinema_db
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres123
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./backend/database/schema.sql:/docker-entrypoint-initdb.d/schema.sql
    networks:
      - cinema_network

  hasura:
    image: hasura/graphql-engine:latest
    container_name: cinema_hasura
    restart: always
    ports:
      - "8080:8080"
    environment:
      HASURA_GRAPHQL_DATABASE_URL: postgres://postgres:postgres123@postgres:5432/cinema_db
      HASURA_GRAPHQL_ENABLE_CONSOLE: "true"
      HASURA_GRAPHQL_DEV_MODE: "true"
      HASURA_GRAPHQL_ENABLED_LOG_TYPES: startup, http-log, webhook-log, websocket-log, query-log
      HASURA_GRAPHQL_ADMIN_SECRET: myadminsecretkey
      HASURA_GRAPHQL_JWT_SECRET: '{"type":"HS256", "key":"your-jwt-secret-key-here"}'
      HASURA_GRAPHQL_UNAUTHORIZED_ROLE: anonymous
    depends_on:
      - postgres
    networks:
      - cinema_network

  go-server:
    build:
      context: ./backend/go-server
      dockerfile: Dockerfile
    container_name: cinema_go_server
    restart: always
    ports:
      - "8081:8080"
    environment:
      DATABASE_URL: postgres://postgres:postgres123@postgres:5432/cinema_db?sslmode=disable
      JWT_SECRET: your-jwt-secret-key-here
      CHAPA_SECRET_KEY: ${CHAPA_SECRET_KEY}
      PORT: 8080
    depends_on:
      - postgres
    volumes:
      - ./backend/go-server/uploads:/app/uploads
    networks:
      - cinema_network

  frontend:
    build:
      context: ./frontend/movie-app
      dockerfile: Dockerfile
    container_name: cinema_frontend
    restart: always
    ports:
      - "3000:3000"
    environment:
      NODE_ENV: development
      HASURA_GRAPHQL_ENDPOINT: http://localhost:8080/v1/graphql
      HASURA_GRAPHQL_ADMIN_SECRET: myadminsecretkey
    depends_on:
      - hasura
    networks:
      - cinema_network

volumes:
  postgres_data:

networks:
  cinema_network:
    driver: bridge
```

```dockerfile:backend/go-server/Dockerfile
FROM golang:1.21-alpine

WORKDIR /app

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main .

# Create uploads directory
RUN mkdir -p uploads/movies

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"]
```

```dockerfile:frontend/movie-app/Dockerfile
FROM node:18-alpine

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application
RUN npm run build

# Expose port
EXPOSE 3000

# Start the application
CMD ["npm", "run", "dev"]
```

## 5. Frontend Modifications (Nuxt 3 + Apollo + GraphQL)

```typescript:frontend/movie-app/nuxt.config.ts
// https://nuxt.com/docs/api/configuration/nuxt-config
export default defineNuxtConfig({
  compatibilityDate: '2025-07-15',
  devtools: { enabled: true },
  modules: [
    '@nuxtjs/tailwindcss',
    '@vueuse/nuxt',
    '@nuxtjs/apollo'
  ],
  css: ['~~/assets/css/main.css'],
  components: [
    { path: '~/app/components', pathPrefix: false }
  ],
  apollo: {
    clients: {
      default: {
        httpEndpoint: process.env.HASURA_GRAPHQL_ENDPOINT || 'http://localhost:8080/v1/graphql',
        httpLinkOptions: {
          headers: {
            'x-hasura-admin-secret': process.env.HASURA_GRAPHQL_ADMIN_SECRET || 'myadminsecretkey'
          }
        }
      }
    }
  },
  runtimeConfig: {
    public: {
      hasuraEndpoint: process.env.HASURA_GRAPHQL_ENDPOINT || 'http://localhost:8080/v1/graphql',
      hasuraAdminSecret: process.env.HASURA_GRAPHQL_ADMIN_SECRET || 'myadminsecretkey',
      chapaPublicKey: process.env.CHAPA_PUBLIC_KEY || ''
    }
  },
  app: {
    head: {
      title: 'Tickent - Cinema Tickets & Movie Schedules',
      meta: [
        { charset: 'utf-8' },
        { name: 'viewport', content: 'width=device-width, initial-scale=1' },
        { name: 'description', content: 'Book cinema tickets and browse movie schedules at Tickent' }
      ]
    }
  }
})
```

```json:frontend/movie-app/package.json
{
  "name": "cinema-tickets-app",
  "private": true,
  "type": "module",
  "scripts": {
    "build": "nuxt build",
    "dev": "nuxt dev",
    "generate": "nuxt generate",
    "preview": "nuxt preview",
    "postinstall": "nuxt prepare"
  },
  "devDependencies": {
    "@nuxt/devtools": "latest",
    "@nuxtjs/apollo": "^5.0.0-beta.8",
    "@nuxtjs/tailwindcss": "^6.8.0",
    "@vueuse/nuxt": "^10.4.1",
    "nuxt": "^3.7.4",
    "vue": "^3.3.8",
    "vue-router": "^4.2.5"
  },
  "dependencies": {
    "@apollo/client": "^3.7.17",
    "@pinia/nuxt": "^0.4.11",
    "@vee-validate/rules": "^4.9.0",
    "@vee-validate/zod": "^4.9.0",
    "graphql": "^16.8.1",
    "graphql-tag": "^2.12.6",
    "pinia": "^2.1.7",
    "vee-validate": "^4.9.0",
    "vue-apollo": "^3.0.0-beta.4",
    "zod": "^3.22.4"
  }
}
```

```typescript:frontend/movie-app/app/plugins/apollo.ts
import { defineNuxtPlugin } from '#app'
import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client/core'
import { setContext } from '@apollo/client/link/context'

export default defineNuxtPlugin(() => {
  const config = useRuntimeConfig()
  
  const httpLink = createHttpLink({
    uri: config.public.hasuraEndpoint,
  })

  const authLink = setContext((_, { headers }) => {
    const token = useCookie('auth-token').value
    
    return {
      headers: {
        ...headers,
        authorization: token ? `Bearer ${token}` : '',
        'x-hasura-admin-secret': config.public.hasuraAdminSecret,
      }
    }
  })

  const apolloClient = new ApolloClient({
    link: authLink.concat(httpLink),
    cache: new InMemoryCache(),
  })

  return {
    provide: {
      apollo: apolloClient
    }
  }
})
```

```typescript:frontend/movie-app/app/composables/useAuth.ts
export const useAuth = () => {
  const user = useState('user', () => null)
  const token = useCookie('auth-token')
  const { $apollo } = useNuxtApp()

  const login = async (email: string, password: string) => {
    try {
      const { data } = await $apollo.mutate({
        mutation: gql`
          mutation Login($email: String!, $password: String!) {
            login(email: $email, password: $password) {
              token
              user {
                id
                email
                name
                role
              }
            }
          }
        `,
        variables: { email, password }
      })

      if (data?.login?.token) {
        token.value = data.login.token
        user.value = data.login.user
        return { success: true, user: data.login.user }
      }
    } catch (error) {
      console.error('Login error:', error)
      return { success: false, error: error.message }
    }
  }

  const register = async (email: string, password: string, name: string) => {
    try {
      const { data } = await $apollo.mutate({
        mutation: gql`
          mutation Register($email: String!, $password: String!, $name: String!) {
            register(email: $email, password: $password, name: $name) {
              message
              user_id
            }
          }
        `,
        variables: { email, password, name }
      })

      if (data?.register?.user_id) {
        return { success: true, message: data.register.message }
      }
    } catch (error) {
      console.error('Register error:', error)
      return { success: false, error: error.message }
    }
  }

  const logout = () => {
    token.value = null
    user.value = null
    navigateTo('/')
  }

  const isAuthenticated = computed(() => !!token.value)
  const isAdmin = computed(() => user.value?.role === 'admin')

  return {
    user: readonly(user),
    login,
    register,
    logout,
    isAuthenticated,
    isAdmin
  }
}
```

```vue:frontend/movie-app/app/components/AuthModal.vue
<template>
  <div v-if="isOpen" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg p-8 max-w-md w-full mx-4">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900">
          {{ isLogin ? 'Sign In' : 'Sign Up' }}
        </h2>
        <button @click="hide" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="space-y-4">
        <div v-if="!isLogin">
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input
            v-model="form.name"
            type="text"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input
            v-model="form.email"
            type="email"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Password</label>
          <input
            v-model="form.password"
            type="password"
            required
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-primary-500"
          />
        </div>

        <div v-if="error" class="text-red-600 text-sm">{{ error }}</div>

        <button
          type="submit"
          :disabled="loading"
          class="w-full bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 disabled:opacity-50"
        >
          {{ loading ? 'Loading...' : (isLogin ? 'Sign In' : 'Sign Up') }}
        </button>
      </form>

      <div class="mt-4 text-center">
        <button @click="toggleMode" class="text-primary-600 hover:text-primary-700">
          {{ isLogin ? "Don't have an account? Sign Up" : 'Already have an account? Sign In' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useAuthStore } from '~/stores/auth'

const authStore = useAuthStore()

const isOpen = ref(false)
const isLogin = ref(true)
const loading = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  email: '',
  password: ''
})

const show = (mode = 'login') => {
  isLogin.value = mode === 'login'
  isOpen.value = true
  error.value = ''
  resetForm()
}

const hide = () => {
  isOpen.value = false
  resetForm()
}

const toggleMode = () => {
  isLogin.value = !isLogin.value
  error.value = ''
  resetForm()
}

const resetForm = () => {
  form.name = ''
  form.email = ''
  form.password = ''
}

const handleSubmit = async () => {
  loading.value = true
  error.value = ''

  try {
    if (isLogin.value) {
      const result = await authStore.login(form.email, form.password)
      if (result.success) {
        hide()
      } else {
        error.value = result.error || 'Login failed'
      }
    } else {
      const result = await authStore.register(form.email, form.password, form.name)
      if (result.success) {
        // Auto-login after successful registration
        const loginResult = await authStore.login(form.email, form.password)
        if (loginResult.success) {
          hide()
        } else {
          error.value = 'Registration successful but login failed'
        }
      } else {
        error.value = result.error || 'Registration failed'
      }
    }
  } catch (err) {
    error.value = 'An unexpected error occurred'
  } finally {
    loading.value = false
  }
}

defineExpose({
  show,
  hide
})
</script>
```

```vue:frontend/movie-app/app/pages/index.vue
<template>
  <div>
    <!-- Hero Section -->
    <HeroSection />

    <!-- Featured Movies Section -->
    <section class="py-20 bg-dark-950">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between mb-12">
          <div>
            <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Featured Movies</h2>
            <p class="text-gray-400 text-lg">Hand-picked favorites just for you</p>
          </div>
          <NuxtLink to="/movies" class="btn-secondary">
            View All Movies
          </NuxtLink>
        </div>

        <!-- Movies Grid -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-8">
          <MovieCard 
            v-for="movie in featuredMovies" 
            :key="movie.id" 
            :movie="movie"
            class="animate-slide-up"
          />
        </div>
      </div>
    </section>

    <!-- Now Showing Section -->
    <section class="py-20 bg-dark-900">
      <div class="container mx-auto px-4">
        <div class="flex items-center justify-between mb-12">
          <div>
            <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Now Showing</h2>
            <p class="text-gray-400 text-lg">Currently playing in theaters</p>
          </div>
          <div class="flex space-x-4">
            <button 
              @click="scrollMovies('left')"
              class="w-10 h-10 bg-dark-800 hover:bg-primary-500 text-white rounded-full flex items-center justify-center transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
              </svg>
            </button>
            <button 
              @click="scrollMovies('right')"
              class="w-10 h-10 bg-dark-800 hover:bg-primary-500 text-white rounded-full flex items-center justify-center transition-colors"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
              </svg>
            </button>
          </div>
        </div>

        <!-- Horizontal Scroll Movies -->
        <div ref="scrollContainer" class="flex space-x-6 overflow-x-auto scrollbar-hide pb-4" style="scroll-behavior: smooth;">
          <div 
            v-for="movie in nowShowingMovies" 
            :key="movie.id"
            class="flex-none w-80"
          >
            <MovieCard :movie="movie" />
          </div>
        </div>
      </div>
    </section>

    <!-- Features Section -->
    <section class="py-20 bg-dark-950">
      <div class="container mx-auto px-4">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Why Choose Tizetabox?</h2>
          <p class="text-gray-400 text-lg max-w-2xl mx-auto">Experience the best of Ethiopian cinema with our premium booking platform</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Feature 1 -->
          <div class="text-center p-6 rounded-lg bg-dark-900 border border-dark-800 hover:border-primary-500/50 transition-all duration-300 group">
            <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:bg-primary-500/20 transition-colors">
              <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <h3 class="font-heading font-bold text-xl text-white mb-4">Secure Payments</h3>
            <p class="text-gray-400">Pay safely with Chapa, TeleBirr, or M-Pesa. Your transactions are fully secured and encrypted.</p>
          </div>

          <!-- Feature 2 -->
          <div class="text-center p-6 rounded-lg bg-dark-900 border border-dark-800 hover:border-primary-500/50 transition-all duration-300 group">
            <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:bg-primary-500/20 transition-colors">
              <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 5v2m0 4v2m0 4v2M5 5a2 2 0 00-2 2v3a2 2 0 110 4v3a2 2 0 002 2h14a2 2 0 002-2v-3a2 2 0 110-4V7a2 2 0 00-2-2H5z"></path>
              </svg>
            </div>
            <h3 class="font-heading font-bold text-xl text-white mb-4">Easy Booking</h3>
            <p class="text-gray-400">Book your tickets in just a few clicks. Select your seats, choose your showtime, and you're done!</p>
          </div>

          <!-- Feature 3 -->
          <div class="text-center p-6 rounded-lg bg-dark-900 border border-dark-800 hover:border-primary-500/50 transition-all duration-300 group">
            <div class="w-16 h-16 bg-primary-500/10 rounded-full flex items-center justify-center mx-auto mb-6 group-hover:bg-primary-500/20 transition-colors">
              <svg class="w-8 h-8 text-primary-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
              </svg>
            </div>
            <h3 class="font-heading font-bold text-xl text-white mb-4">Multiple Locations</h3>
            <p class="text-gray-400">Choose from premium cinemas across Addis Ababa. Find the perfect location and showtime for you.</p>
          </div>
        </div>
      </div>
    </section>

    <!-- CTA Section -->
    <section class="py-20 bg-gradient-to-r from-primary-600 to-primary-500">
      <div class="container mx-auto px-4 text-center">
        <h2 class="text-3xl md:text-4xl font-heading font-bold text-white mb-4">Ready for Your Next Movie Night?</h2>
        <p class="text-primary-100 text-lg mb-8 max-w-2xl mx-auto">Join thousands of movie lovers who trust Tizetabox for their cinema experience</p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <NuxtLink to="/movies" class="bg-white text-primary-600 hover:bg-gray-100 font-bold py-4 px-8 rounded-lg transition-colors">
            Browse Movies
          </NuxtLink>
          <button @click="showAuth('register')" class="bg-transparent border-2 border-white text-white hover:bg-white hover:text-primary-600 font-bold py-4 px-8 rounded-lg transition-all">
            Sign Up Today
          </button>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { useMoviesStore } from '~/stores/movies'

const moviesStore = useMoviesStore()
const scrollContainer = ref()

const featuredMovies = computed(() => moviesStore.movies.slice(0, 4))
const nowShowingMovies = computed(() => moviesStore.nowShowing)

const showAuth = (mode) => {
  const authModal = inject('authModal')
  authModal?.show(mode)
}

const scrollMovies = (direction) => {
  if (!scrollContainer.value) return
  
  const scrollAmount = 320 // Width of movie card + gap
  const currentScroll = scrollContainer.value.scrollLeft
  
  if (direction === 'left') {
    scrollContainer.value.scrollLeft = currentScroll - scrollAmount
  } else {
    scrollContainer.value.scrollLeft = currentScroll + scrollAmount
  }
}

// Page meta
useSeoMeta({
  title: 'Tizetabox - Ethiopia\'s Premier Cinema Booking Platform',
  ogTitle: 'Tizetabox - Ethiopia\'s Premier Cinema Booking Platform',
  description: 'Book movie tickets online in Ethiopia. Discover the latest movies, choose your seats, and enjoy premium cinema experiences across Addis Ababa.',
  ogDescription: 'Book movie tickets online in Ethiopia. Discover the latest movies, choose your seats, and enjoy premium cinema experiences across Addis Ababa.',
  ogImage: 'https://images.pexels.com/photos/7991579/pexels-photo-7991579.jpeg',
  twitterCard: 'summary_large_image'
})
</script>

<style>
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}

.scrollbar-hide::-webkit-scrollbar {
  display: none;
}

.animate-slide-up {
  animation: slideUp 0.6s ease-out;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
```

```typescript:frontend/movie-app/stores/auth.ts
import { defineStore } from 'pinia'

interface User {
  id: number
  email: string
  name: string
  role: string
}

interface AuthState {
  user: User | null
  token: string | null
  loading: boolean
}

export const useAuthStore = defineStore('auth', {
  state: (): AuthState => ({
    user: null,
    token: null,
    loading: false
  }),

  getters: {
    isAuthenticated: (state) => !!state.token,
    isAdmin: (state) => state.user?.role === 'admin'
  },

  actions: {
    async initializeAuth() {
      // Check for existing token in localStorage
      const token = localStorage.getItem('auth-token')
      if (token) {
        this.token = token
        // Verify token and get user info
        await this.verifyToken()
      }
    },

    async login(email: string, password: string) {
      this.loading = true
      try {
        // This would be replaced with actual API call
        const response = await $fetch('/api/auth/login', {
          method: 'POST',
          body: { email, password }
        })
        
        this.token = response.token
        this.user = response.user
        localStorage.setItem('auth-token', response.token)
        
        return { success: true, user: response.user }
      } catch (error) {
        console.error('Login error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async register(email: string, password: string, name: string) {
      this.loading = true
      try {
        const response = await $fetch('/api/auth/register', {
          method: 'POST',
          body: { email, password, name }
        })
        
        return { success: true, message: response.message }
      } catch (error) {
        console.error('Register error:', error)
        return { success: false, error: error.message }
      } finally {
        this.loading = false
      }
    },

    async verifyToken() {
      try {
        const response = await $fetch('/api/auth/verify', {
          headers: {
            Authorization: `Bearer ${this.token}`
          }
        })
        this.user = response.user
      } catch (error) {
        this.logout()
      }
    },

    logout() {
      this.user = null
      this.token = null
      localStorage.removeItem('auth-token')
    }
  }
})
```

```typescript:frontend/movie-app/stores/movies.ts
import { defineStore } from 'pinia'

interface Movie {
  id: number
  title: string
  description: string
  poster: string
  backdrop: string
  rating: number
  duration: number
  year: string
  director: string
  genres: string[]
  cast: string[]
  trailer?: string
  schedules: Schedule[]
  isFeatured: boolean
}

interface Schedule {
  id: number
  movie_id: number
  cinema_hall_id: number
  start_time: string
  end_time: string
  price: number
  available_seats: number
  cinema_hall: {
    name: string
    cinema: {
      name: string
      city: string
    }
  }
}

interface Cinema {
  id: number
  name: string
  location: string
  address: string
  phone: string
  facilities: string[]
  screens: number
  image: string
}

interface MoviesState {
  movies: Movie[]
  cinemas: Cinema[]
  loading: boolean
  searchQuery: string
  selectedGenres: string[]
  featuredMovie: Movie | null
}

export const useMoviesStore = defineStore('movies', {
  state: (): MoviesState => ({
    movies: [],
    cinemas: [],
    loading: false,
    searchQuery: '',
    selectedGenres: [],
    featuredMovie: null
  }),

  getters: {
    filteredMovies: (state) => {
      let filtered = state.movies

      if (state.searchQuery) {
        filtered = filtered.filter(movie => 
          movie.title.toLowerCase().includes(state.searchQuery.toLowerCase()) ||
          movie.description.toLowerCase().includes(state.searchQuery.toLowerCase())
        )
      }

      if (state.selectedGenres.length > 0) {
        filtered = filtered.filter(movie => 
          movie.genres.some(genre => state.selectedGenres.includes(genre))
        )
      }

      return filtered
    },

    nowShowing: (state) => state.movies.filter(movie => movie.schedules.length > 0),

    allGenres: (state) => {
      const genres = new Set<string>()
      state.movies.forEach(movie => {
        movie.genres.forEach(genre => genres.add(genre))
      })
      return Array.from(genres).sort()
    }
  },

  actions: {
    async fetchMovies() {
      this.loading = true
      try {
        // Mock data for now - replace with actual API call
        this.movies = [
          {
            id: 1,
            title: "The Lion King",
            description: "A young lion prince flees his kingdom only to learn the true meaning of responsibility and bravery.",
            poster: "/images/movie-placeholder.jpg",
            backdrop: "/images/movie-placeholder.jpg",
            rating: 4.5,
            duration: 118,
            year: "2019",
            director: "Jon Favreau",
            genres: ["Animation", "Adventure", "Drama"],
            cast: ["Donald Glover", "Beyoncé", "James Earl Jones"],
            schedules: [],
            isFeatured: true
          },
          {
            id: 2,
            title: "Avengers: Endgame",
            description: "After the devastating events of Avengers: Infinity War, the universe is in ruins.",
            poster: "/images/movie-placeholder.jpg",
            backdrop: "/images/movie-placeholder.jpg",
            rating: 4.8,
            duration: 181,
            year: "2019",
            director: "Anthony Russo",
            genres: ["Action", "Adventure", "Drama"],
            cast: ["Robert Downey Jr.", "Chris Evans", "Mark Ruffalo"],
            schedules: [],
            isFeatured: true
          }
        ]
        this.featuredMovie = this.movies[0]
      } catch (error) {
        console.error('Error fetching movies:', error)
      } finally {
        this.loading = false
      }
    },

    async fetchCinemas() {
      try {
        // Mock data for now - replace with actual API call
        this.cinemas = [
          {
            id: 1,
            name: "CineMax Downtown",
            location: "Downtown",
            address: "123 Main Street, Downtown",
            phone: "+1-555-0101",
            facilities: ["Dolby Atmos", "Recliner Seats", "Food Service"],
            screens: 8,
            image: "/images/cinema-placeholder.jpg"
          }
        ]
      } catch (error) {
        console.error('Error fetching cinemas:', error)
      }
    },

    setSearchQuery(query: string) {
      this.searchQuery = query
    },

    toggleGenreFilter(genre: string) {
      const index = this.selectedGenres.indexOf(genre)
      if (index > -1) {
        this.selectedGenres.splice(index, 1)
      } else {
        this.selectedGenres.push(genre)
      }
    },

    clearFilters() {
      this.searchQuery = ''
      this.selectedGenres = []
    }
  }
})
```

```typescript:frontend/movie-app/app/composables/useCinemas.ts
import { ref, computed } from 'vue'

interface Cinema {
  id: number
  name: string
  location: string
  address: string
  phone: string
  facilities: string[]
  screens: number
  image: string
}

export const useCinemas = () => {
  const cinemas = ref<Cinema[]>([])
  const loading = ref(false)

  const fetchCinemas = async () => {
    loading.value = true
    try {
      // Mock data for now - replace with actual API call
      cinemas.value = [
        {
          id: 1,
          name: "CineMax Downtown",
          location: "Downtown",
          address: "123 Main Street, Downtown",
          phone: "+1-555-0101",
          facilities: ["Dolby Atmos", "Recliner Seats", "Food Service"],
          screens: 8,
          image: "/images/cinema-placeholder.jpg"
        }
      ]
    } catch (error) {
      console.error('Error fetching cinemas:', error)
    } finally {
      loading.value = false
    }
  }

  return {
    cinemas: readonly(cinemas),
    loading: readonly(loading),
    fetchCinemas
  }
}
```

```css:frontend/movie-app/assets/css/main.css
@tailwind base;
@tailwind components;
@tailwind utilities;

@layer base {
  :root {
    --color-primary-50: #eff6ff;
    --color-primary-100: #dbeafe;
    --color-primary-200: #bfdbfe;
    --color-primary-300: #93c5fd;
    --color-primary-400: #60a5fa;
    --color-primary-500: #3b82f6;
    --color-primary-600: #2563eb;
    --color-primary-700: #1d4ed8;
    --color-primary-800: #1e40af;
    --color-primary-900: #1e3a8a;
    --color-primary-950: #172554;
    
    --color-dark-50: #f8fafc;
    --color-dark-100: #f1f5f9;
    --color-dark-200: #e2e8f0;
    --color-dark-300: #cbd5e1;
    --color-dark-400: #94a3b8;
    --color-dark-500: #64748b;
    --color-dark-600: #475569;
    --color-dark-700: #334155;
    --color-dark-800: #1e293b;
    --color-dark-900: #0f172a;
    --color-dark-950: #020617;
  }
}

@layer components {
  .btn-primary {
    @apply bg-primary-600 hover:bg-primary-700 text-white font-bold py-3 px-6 rounded-lg transition-colors;
  }
  
  .btn-secondary {
    @apply bg-dark-800 hover:bg-dark-700 text-white font-bold py-3 px-6 rounded-lg transition-colors;
  }
  
  .btn-outline {
    @apply bg-transparent border-2 border-primary-600 text-primary-600 hover:bg-primary-600 hover:text-white font-bold py-3 px-6 rounded-lg transition-all;
  }
  
  .container {
    @apply max-w-7xl mx-auto px-4 sm:px-6 lg:px-8;
  }
  
  .font-heading {
    font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
  }
}

@layer utilities {
  .text-shadow {
    text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .text-shadow-lg {
    text-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }
  
  .backdrop-blur {
    backdrop-filter: blur(8px);
  }
  
  .scrollbar-hide {
    -ms-overflow-style: none;
    scrollbar-width: none;
  }
  
  .scrollbar-hide::-webkit-scrollbar {
    display: none;
  }
}

/* Custom animations */
@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes scaleIn {
  from {
    opacity: 0;
    transform: scale(0.9);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

.animate-fade-in {
  animation: fadeIn 0.6s ease-out;
}

.animate-slide-up {
  animation: slideUp 0.6s ease-out;
}

.animate-scale-in {
  animation: scaleIn 0.6s ease-out;
}
```

```javascript:frontend/movie-app/tailwind.config.js
/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,vue,ts}",
    "./pages/**/*.{js,vue,ts}",
    "./components/**/*.{js,vue,ts}",
    "./layouts/**/*.{js,vue,ts}",
    "./plugins/**/*.{js,ts}",
    "./nuxt.config.{js,ts}",
  ],
  theme: {
    extend: {
      colors: {
        primary: {
          50: '#eff6ff',
          100: '#dbeafe',
          200: '#bfdbfe',
          300: '#93c5fd',
          400: '#60a5fa',
          500: '#3b82f6',
          600: '#2563eb',
          700: '#1d4ed8',
          800: '#1e40af',
          900: '#1e3a8a',
          950: '#172554',
        },
        dark: {
          50: '#f8fafc',
          100: '#f1f5f9',
          200: '#e2e8f0',
          300: '#cbd5e1',
          400: '#94a3b8',
          500: '#64748b',
          600: '#475569',
          700: '#334155',
          800: '#1e293b',
          900: '#0f172a',
          950: '#020617',
        },
      },
      fontFamily: {
        heading: ['Inter', 'system-ui', 'sans-serif'],
        sans: ['Inter', 'system-ui', 'sans-serif'],
      },
      animation: {
        'fade-in': 'fadeIn 0.6s ease-out',
        'slide-up': 'slideUp 0.6s ease-out',
        'scale-in': 'scaleIn 0.6s ease-out',
      },
      keyframes: {
        fadeIn: {
          '0%': { opacity: '0', transform: 'translateY(20px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        slideUp: {
          '0%': { opacity: '0', transform: 'translateY(30px)' },
          '100%': { opacity: '1', transform: 'translateY(0)' },
        },
        scaleIn: {
          '0%': { opacity: '0', transform: 'scale(0.9)' },
          '100%': { opacity: '1', transform: 'scale(1)' },
        },
      },
    },
  },
  plugins: [],
}
```

```
