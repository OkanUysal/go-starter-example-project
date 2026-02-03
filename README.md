# Go Starter Example Project

A production-ready Go REST API starter template with authentication, authorization, real-time WebSocket communication, caching, metrics, and comprehensive developer tooling.

## ✨ Features

### Authentication & Authorization
- 🔐 **JWT Authentication** - Access & refresh token flow
- 🚫 **Token Blacklisting** - Family-based token invalidation
- 👤 **Guest Login** - Anonymous user support with optional ID reuse
- 🛡️ **Role-Based Access** - Admin middleware and user roles
- 🔄 **Token Refresh** - Secure token rotation with automatic blacklisting

### Real-Time Communication
- 🌐 **WebSocket Support** - Real-time bidirectional communication
- 🏠 **Public Lobby** - Always-open room for all users
- 🎮 **Dynamic Game Rooms** - Admin-created rooms with player limits
- 👥 **Room Management** - Join, leave, invite, and broadcast messages
- 📨 **Message Types** - Chat, game events, room notifications

### Performance & Caching
- ⚡ **Multi-Backend Cache** - Memory and Redis support
- 🚀 **Smart Caching** - Token blacklist, user data, and statistics
- 📊 **Cache Hit/Miss Logging** - Performance monitoring
- ⏱️ **Configurable TTL** - Environment-based cache duration

### API & Documentation
- 📚 **Swagger/OpenAPI** - Auto-generated interactive API docs
- ✅ **Standardized Responses** - Consistent API response format
- 🎯 **Error Handling** - Structured error codes and messages

### Observability
- 📊 **Prometheus Metrics** - HTTP metrics with automatic collection
- ☁️ **Grafana Cloud** - Optional metrics push integration
- 🏥 **Health Checks** - `/health` and `/metrics` endpoints
- 📝 **Structured Logging** - JSON logs with go-logger

### Database
- 🗄️ **PostgreSQL + GORM** - Production-ready ORM
- 🔄 **Migrations** - SQL migration files included
- 🔧 **Dynamic Table Names** - Environment-based table configuration

### Developer Experience
- 🔧 **Environment Config** - `.env` file support
- 📦 **Modular Structure** - Clean separation of concerns
- 🚀 **Railway Ready** - One-click deployment configuration

## 🚀 Quick Start

### Creating a New Project

You can use the included PowerShell script to create a new project based on this template:

```powershell
# Run from the parent directory of this repository
.\go-starter-example-project\create-new-project.ps1 -ProjectName "my-new-api" -ModulePath "github.com/username/my-new-api"
```

This will:
- Create a new project directory
- Copy all template files
- Update module paths and references
- Initialize git repository
- Set up a fresh project ready for development

### Manual Installation

### Prerequisites
- Go 1.21+
- PostgreSQL database

### Installation

1. **Clone the repository**
```bash
git clone https://github.com/OkanUysal/go-starter-example-project.git
cd go-starter-example-project
```

2. **Install dependencies**
```bash
go mod download
```

3. **Configure environment**
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. **Run migrations**
```bash
# Connect to your PostgreSQL database and run migration files in migrations/ folder
psql $DATABASE_URL_LOCAL -f migrations/001_create_users_table.up.sql
psql $DATABASE_URL_LOCAL -f migrations/002_create_token_blacklist.up.sql
psql $DATABASE_URL_LOCAL -f migrations/003_add_family_id_to_blacklist.up.sql
```

5. **Start the server**
```bash
go run main.go
```

Server will start on `http://localhost:8080`

## 📖 API Documentation

Once the server is running, visit:
- **Swagger UI**: http://localhost:8080/swagger/index.html
- **API Spec**: http://localhost:8080/swagger.json

### Available Endpoints

#### Public Endpoints
- `GET /api/hello` - Simple hello endpoint
- `POST /api/auth/guest-login` - Guest user login
- `POST /api/auth/refresh` - Refresh access token

#### Protected Endpoints (Requires Authentication)
- `GET /api/auth/me` - Get current user info

#### Admin Endpoints (Requires Admin Role)
- `GET /api/admin/dashboard` - Admin dashboard with statistics
- `GET /api/admin/users` - List all users

#### WebSocket Endpoints (Requires Authentication)
- `GET /api/ws?room_id=lobby` - Connect to WebSocket (room_id optional, defaults to lobby)
- `GET /api/ws/rooms` - Get all active rooms
- `GET /api/ws/rooms/:room_id` - Get room information

#### WebSocket Admin Endpoints (Requires Admin Role)
- `POST /api/ws/rooms` - Create a new game room
- `DELETE /api/ws/rooms/:room_id` - Close a game room
- `POST /api/ws/invite` - Invite users to a room

#### Monitoring
- `GET /health` - Health check
- `GET /metrics` - Prometheus metrics

## 🔐 Authentication Flow

### 1. Guest Login
```bash
curl -X POST http://localhost:8080/api/auth/guest-login \
  -H "Content-Type: application/json" \
  -d '{}'
```

Response:
```json
{
  "success": true,
  "data": {
    "access_token": "eyJhbGc...",
    "refresh_token": "eyJhbGc...",
    "user": {
      "id": "uuid",
      "display_name": "Guest1234",
      "role": "USER",
      "is_guest": true
    }
  },
  "message": "Guest login successful"
}
```

### 2. Access Protected Endpoint
```bash
curl http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN"
```

### 3. Refresh Token
```bash
curl -X POST http://localhost:8080/api/auth/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token": "YOUR_REFRESH_TOKEN"}'
```

**Note**: When you refresh, both old access and refresh tokens are invalidated (family-based blacklisting).

## ⚙️ Configuration

### Environment Variables

```bash
# Application
APP_NAME=go-starter-example-project
PORT=8080
ENVIRONMENT=development

# Database
DATABASE_URL_LOCAL=postgresql://user:pass@host:port/db?sslmode=require

# JWT
JWT_SECRET=your-secret-key-change-in-production
ACCESS_TOKEN_DURATION=24    # hours
REFRESH_TOKEN_DURATION=168  # hours (7 days)

# Database Tables
USER_TABLE=example_user
TOKEN_BLACKLIST_TABLE=example_token_blacklist

# Cache Configuration
CACHE_TYPE=memory           # or "redis"
CACHE_TTL=300              # seconds (5 minutes)
REDIS_URL=                 # redis://localhost:6379 (if using Redis)

# Metrics
SERVICE_NAME=go-starter-example-project
METRICS_ENABLED=true

# Grafana Cloud (Optional)
GRAFANA_CLOUD_URL=https://prometheus-prod-XX-prod-XX.grafana.net/api/prom/push
GRAFANA_CLOUD_USER=123456
GRAFANA_CLOUD_KEY=glc_xxxxx
```

## 🌐 WebSocket Usage

### Connect to WebSocket

```javascript
// Connect to public lobby
const ws = new WebSocket('ws://localhost:8080/api/ws?room_id=lobby', {
  headers: {
    'Authorization': 'Bearer YOUR_ACCESS_TOKEN'
  }
});

// Or connect to a specific game room
const ws = new WebSocket('ws://localhost:8080/api/ws?room_id=ROOM_ID', {
  headers: {
    'Authorization': 'Bearer YOUR_ACCESS_TOKEN'
  }
});

ws.onopen = () => {
  console.log('Connected to WebSocket');
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log('Received:', message);
  
  // Handle different message types
  switch(message.type) {
    case 'join':
      console.log(`${message.data.username} joined`);
      break;
    case 'leave':
      console.log(`${message.data.username} left`);
      break;
    case 'chat':
      console.log(`${message.data.username}: ${message.data.content}`);
      break;
    case 'invite':
      console.log('You have been invited to:', message.data.room);
      break;
    case 'room_created':
      console.log('New room created:', message.data.room);
      break;
    case 'room_closed':
      console.log('Room closed:', message.data.room_id);
      break;
  }
};
```

### Admin: Create a Game Room

```bash
curl -X POST http://localhost:8080/api/ws/rooms \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Game Room 1",
    "max_players": 4
  }'
```

### Admin: Invite Users to Room

```bash
curl -X POST http://localhost:8080/api/ws/invite \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "room_id": "ROOM_ID",
    "user_ids": ["user-id-1", "user-id-2"]
  }'
```

### Admin: Close a Game Room

```bash
curl -X DELETE http://localhost:8080/api/ws/rooms/ROOM_ID \
  -H "Authorization: Bearer YOUR_ADMIN_TOKEN"
```

### Get All Rooms

```bash
curl http://localhost:8080/api/ws/rooms \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### WebSocket Message Types

- **join**: User joined a room
- **leave**: User left a room
- **chat**: Chat message
- **game_event**: Game-specific events
- **room_created**: New room created (broadcast to lobby)
- **room_closed**: Room closed by admin
- **invite**: User invited to a room
- **error**: Error message

### Use Cases

1. **Queue System**: 
   - Users join the public lobby
   - Admin creates game rooms when enough players are ready
   - Admin invites specific players to game rooms

2. **Game Rooms**:
   - Dynamic room creation with player limits
   - Room-based communication
   - Automatic room cleanup when game ends

3. **Chat System**:
   - Public lobby for general chat
   - Private game rooms for team communication

## 📊 Monitoring with Grafana Cloud (Optional)

This project includes built-in support for pushing metrics to Grafana Cloud:

1. **Sign up** at [grafana.com](https://grafana.com) (free tier available)
2. **Get your credentials** from your Grafana Cloud stack
3. **Add to .env**:
```bash
GRAFANA_CLOUD_URL=https://prometheus-prod-XX-prod-XX.grafana.net/api/prom/push
GRAFANA_CLOUD_USER=123456
GRAFANA_CLOUD_KEY=glc_xxxxx
```
4. **Restart** the application - metrics will automatically push every 15 seconds

### Local Metrics

Even without Grafana Cloud, metrics are available at:
- http://localhost:8080/metrics (Prometheus format)

## 🏗️ Project Structure

```
.
├── auth/                    # Authentication & authorization
│   ├── jwt.go              # JWT token generation & validation
│   ├── service.go          # Auth business logic
│   ├── middleware.go       # JWT middleware
│   ├── admin_middleware.go # Admin access middleware
│   └── blacklist.go        # Token blacklist operations
├── config/                  # Configuration
│   └── database.go         # Database connection & helpers
├── docs/                    # Swagger documentation (auto-generated)
├── handlers/                # HTTP handlers
│   ├── auth.go             # Auth endpoints
│   ├── admin.go            # Admin endpoints
│   └── hello.go            # Example endpoint
├── migrations/              # Database migrations
├── models/                  # Database models
│   ├── user.go             # User model
│   ├── token_blacklist.go  # Token blacklist model
│   └── helpers.go          # Model helpers
├── main.go                  # Application entry point
├── .env.example             # Example environment variables
└── .gitignore
```

## 🚀 Deployment

### Railway

This project is ready for Railway deployment:

1. Push to GitHub
2. Create new project on [Railway](https://railway.app)
3. Connect your GitHub repository
4. Add PostgreSQL database
5. Set environment variables
6. Deploy!

Railway will automatically detect `railway.json` configuration.

### Docker

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./main"]
```

## 🛠️ Development

### Generate Swagger Docs

After modifying API endpoints:
```bash
swag init
```

### Database Migrations

Create new migration:
```sql
-- migrations/004_your_migration.up.sql
CREATE TABLE ...;

-- migrations/004_your_migration.down.sql
DROP TABLE ...;
```

## 📦 Used Libraries

- [gin-gonic/gin](https://github.com/gin-gonic/gin) - HTTP web framework
- [gorm.io/gorm](https://github.com/go-gorm/gorm) - ORM library
- [golang-jwt/jwt](https://github.com/golang-jwt/jwt) - JWT implementation
- [swaggo/swag](https://github.com/swaggo/swag) - Swagger documentation
- [@OkanUysal/go-logger](https://github.com/OkanUysal/go-logger) - Structured logging
- [@OkanUysal/go-metrics](https://github.com/OkanUysal/go-metrics) - Prometheus metrics
- [@OkanUysal/go-swagger](https://github.com/OkanUysal/go-swagger) - Swagger helpers
- [@OkanUysal/go-response](https://github.com/OkanUysal/go-response) - API responses

## 🔒 Security Features

- ✅ JWT token-based authentication
- ✅ Refresh token rotation
- ✅ Token family blacklisting (invalidates both access & refresh)
- ✅ Role-based authorization
- ✅ Secure password hashing (ready for password auth)
- ✅ Environment-based secrets

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 🎯 Creating a New Project from This Template

When starting a new project based on this template, follow these steps:

### Option 1: Automated Setup (Recommended)

**Step 1: Copy the project to a new location**
```powershell
# Windows
Copy-Item -Recurse go-starter-example-project C:\Projects\my-awesome-project
cd C:\Projects\my-awesome-project

# Linux/Mac
cp -r go-starter-example-project ~/projects/my-awesome-project
cd ~/projects/my-awesome-project
```

**Step 2: Run the setup script**
```powershell
.\setup-new-project.ps1 -ProjectName "my-awesome-project" -ModulePath "github.com/yourusername/my-awesome-project"
```

The script will ask for confirmation and then automatically:
- Update `go.mod` with your new module path
- Replace all import paths in Go files
- Update migration files with your table names (e.g., `example_user` → `myawesomeproject_user`)
- Update Swagger documentation
- Configure `.env` files with your project name
- Update README with your project details
- Run `go mod tidy`

**Step 3: Update environment variables**
Edit `.env` file with your actual values (see Step 5 below)

**Step 4: Done!**
Your project is ready to customize and extend.

### Option 2: Manual Setup

If you prefer manual setup, follow these steps:

#### 1. Copy the Project
```bash
# Windows
Copy-Item -Recurse go-starter-example-project C:\Projects\your-new-project
cd C:\Projects\your-new-project

# Linux/Mac
cp -r go-starter-example-project your-new-project
cd your-new-project
```

#### 2. Update Go Module Name
Edit `go.mod`:
```go
module github.com/yourusername/your-new-project

go 1.24
// ... rest of the file
```

#### 3. Update All Import Paths
Replace all occurrences of `github.com/OkanUysal/go-starter-example-project` with your new module path:

**Files to update:**
- `main.go`
- `handlers/*.go`
- `auth/*.go`
- `config/*.go`
- `websocket/*.go`
- `models/*.go`

**Find and replace:**
```bash
# Linux/Mac
find . -type f -name "*.go" -exec sed -i 's|github.com/OkanUysal/go-starter-example-project|github.com/yourusername/your-new-project|g' {} +

# Windows PowerShell
Get-ChildItem -Recurse -Filter *.go | ForEach-Object {
    (Get-Content $_.FullName) -replace 'github.com/OkanUysal/go-starter-example-project', 'github.com/yourusername/your-new-project' | Set-Content $_.FullName
}
```

#### 4. Update Swagger Documentation
Edit the `@title` and `@description` in `main.go`:
```go
// @title           Your New Project API
// @version         1.0.0
// @description     Your project description here
```

Then regenerate Swagger docs:
```bash
swag init
```

#### 5. Update Environment Variables
Copy and customize `.env.example`:
```bash
cp .env.example .env
```

Update these values:
- `DATABASE_URL_LOCAL` - Your database connection
- `JWT_SECRET` - Generate a new secret: `openssl rand -base64 32`
- `SERVICE_NAME` - Your new project name
- `USER_TABLE` - Your user table name (e.g., `your_project_user`)
- `TOKEN_BLACKLIST_TABLE` - Your blacklist table name (e.g., `your_project_token_blacklist`)

#### 6. Update Database Table Names
The project uses environment-based table names to avoid conflicts:

**Default tables:**
- `example_user`
- `example_token_blacklist`

**Change to your project-specific names:**
```bash
USER_TABLE=myapp_user
TOKEN_BLACKLIST_TABLE=myapp_token_blacklist
```

#### 7. Run Database Migrations
```bash
# Install go-migrate if you haven't
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Run migrations
migrate -path migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up
```

#### 8. Clean Up Git History
```bash
rm -rf .git
git init
git add .
git commit -m "Initial commit: Project based on go-starter-example-project"
```

#### 9. Verify Everything Works
```bash
# Install dependencies
go mod tidy

# Run the server
go run main.go

# Test endpoints
curl http://localhost:8080/health
curl http://localhost:8080/metrics
```

#### 10. Customize Further

Now you're ready to build your application:
- Add new models in `models/`
- Create handlers in `handlers/`
- Add business logic in service files
- Extend WebSocket functionality in `websocket/`
- Add new migrations as needed

### 🎨 Recommended Customizations

**For Your Specific Use Case:**
1. **Remove unused features** - If you don't need WebSocket, remove the `websocket/` package
2. **Add new middleware** - CORS, rate limiting, request ID tracking
3. **Extend authentication** - Add OAuth providers, email/password login
4. **Custom business logic** - Add domain-specific services and repositories
5. **Enhanced monitoring** - Add custom Prometheus metrics for your features

**Project Structure Best Practices:**
```
your-new-project/
├── cmd/                    # Multiple entry points (optional)
├── internal/              # Private application code
│   ├── domain/           # Business logic and entities
│   ├── repository/       # Data access layer
│   └── service/          # Business services
├── pkg/                   # Public libraries (optional)
├── api/                   # API-specific code
│   ├── handlers/         # HTTP handlers
│   └── middleware/       # Middleware
└── migrations/           # Database migrations
```

## 📝 License

MIT License - feel free to use this project for your own purposes.

## 🙏 Acknowledgments

Built with Go and modern best practices for production-ready REST APIs.

---

**Happy Coding! 🚀**
