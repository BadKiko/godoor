# 🍳 GoDoor - Tandoor Recipes on Go

[![Docker Image](https://img.shields.io/docker/v/badkiko/godoor/latest?label=Docker)](https://hub.docker.com/r/badkiko/godoor)
[![Docker Pulls](https://img.shields.io/docker/pulls/badkiko/godoor)](https://hub.docker.com/r/badkiko/godoor)
[![Docker Image Size](https://img.shields.io/docker/image-size/badkiko/godoor/latest)](https://hub.docker.com/r/badkiko/godoor)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8)](https://golang.org/)

A lightweight, fast, and fully compatible implementation of Tandoor Recipes built with Go. Perfect for self-hosting your recipe management system with Docker support.

**🚀 Ready to deploy with a single command!**

**Version: 2.3.6** (compatible with Tandoor 2.3.6)

## ✨ Features

### 🔐 Authentication & Security
- **Token-based authentication** with OAuth2 support
- **Password hashing** with bcrypt
- **Session management** with automatic token cleanup
- **Space-based access control**

### 👥 User Management
- **Multi-user support** with role-based permissions
- **Space management** for organizing recipes
- **User preferences** and customization
- **Profile management**

### 📚 Recipe Management
- **Complete CRUD operations** for recipes
- **Step-by-step instructions** with timing
- **Ingredient management** with units and amounts
- **Keyword categorization**
- **Recipe books** and collections
- **Image upload support**

### 🛒 Shopping Lists
- **Smart shopping list generation** from recipes
- **Ingredient tracking** and management
- **List sharing** and collaboration

### 📅 Meal Planning
- **Weekly meal planning** with calendar integration
- **Recipe scheduling** and organization
- **Shopping list generation** from meal plans

### 🏷️ Organization
- **Keyword system** for recipe categorization
- **Recipe books** for collections and organization
- **Advanced filtering** and search capabilities

### ✅ Feature Status

| Feature | Status | Notes |
|---------|--------|-------|
| 🔐 Authentication | ✅ Complete | OAuth2 tokens, bcrypt hashing |
| 👥 User Management | ✅ Complete | Multi-user with spaces |
| 📚 Recipes | ✅ Complete | Full CRUD with images |
| 🛒 Shopping Lists | ✅ Complete | Smart list generation |
| 📅 Meal Planning | ✅ Complete | Calendar integration |
| 🏷️ Keywords | ✅ Complete | Recipe categorization |
| 📖 Recipe Books | ✅ Complete | Collections & organization |
| ⚙️ User Preferences | ✅ Complete | Customization options |
| 🔍 Search & Filtering | ✅ Complete | Advanced query support |
| 📱 Mobile App Support | ✅ Complete | Android/iOS compatible |

### 🐳 Docker Support

- **Official Docker image:** `badkiko/godoor`
- **Multi-stage build** for optimal size
- **Volume mounting** for persistent data
- **Health checks** included
- **Docker Compose** ready

## 📖 Usage Examples

### 🔑 Authentication

```bash
# Get access token
curl -X POST http://localhost:8080/api-token-auth/ \
  -H "Content-Type: application/json" \
  -d '{"username": "admin", "password": "admin123"}'

# Use token in requests
curl -H "Authorization: Bearer YOUR_TOKEN" \
  http://localhost:8080/api/recipe/
```

### 🍳 Working with Recipes

```bash
# List recipes
curl -H "Authorization: Bearer TOKEN" \
  "http://localhost:8080/api/recipe/?page=1&page_size=10"

# Create recipe
curl -X POST -H "Authorization: Bearer TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Pasta Carbonara", "servings": 4}' \
  http://localhost:8080/api/recipe/

# Search recipes
curl -H "Authorization: Bearer TOKEN" \
  "http://localhost:8080/api/recipe/?query=pasta&keywords=pasta,italian"
```

## 🚀 Quick Start

### 🐳 Docker (Recommended)

The easiest way to get started is with Docker:

```bash
# Pull from Docker Hub
docker pull badkiko/godoor

# Run with persistent database
docker run -d \
  --name godoor \
  -p 8080:8080 \
  -v $(pwd)/godoor.db:/app/godoor.db \
  badkiko/godoor
```

Or use Docker Compose:

```bash
# Create data directory
mkdir -p data

# Start with Docker Compose
docker-compose up -d

# Check logs
docker-compose logs -f godoor

# Stop the service
docker-compose down
```

**Your data will persist** between container restarts thanks to the volume mounts.

### 📁 Data Persistence

The Docker setup includes two volume mounts:

- `godoor.db` - SQLite database file
- `data/` - Directory for uploaded files and additional data

Make sure these files/directories exist before starting:

```bash
mkdir -p data
touch godoor.db  # Optional: will be created automatically
```

### 🛠️ Local Development

For development or advanced customization:

```bash
# Prerequisites: Go 1.21+
git clone <repository-url>
cd godoor

# Install dependencies
go mod tidy

# Copy environment file
cp env.example .env

# Build and run
go build -ldflags="-s -w" -o godoor .
./godoor
```

### ⚙️ Configuration

#### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `GODOOR_DB_PATH` | `godoor.db` | SQLite database file path |

#### Default Admin Credentials

- **Username:** `admin`
- **Password:** `admin123`

**⚠️ Change these credentials after first login!**

## Configuration

Create `.env` file in the project root (see `env.example` for reference):

```bash
# Admin user
ADMIN_USERNAME=admin
ADMIN_PASSWORD=admin123

# Server port
PORT=8080
```

Default admin credentials: `admin` / `admin123`

## 🔗 API Compatibility

GoDoor is **100% compatible** with the original Tandoor Recipes API. All endpoints, request/response formats, and authentication methods work identically.

### 📋 Supported Endpoints

- **Authentication:** `/api-token-auth/`, `/api-auth/`
- **Users & Spaces:** `/api/user/`, `/api/space/`, `/api/user-space/`
- **Recipes:** `/api/recipe/`, `/api/recipe/:id/`
- **Recipe Books:** `/api/recipe-book/`, `/api/recipe-book-entry/`
- **Keywords:** `/api/keyword/`
- **Meal Planning:** `/api/meal-plan/`, `/api/meal-type/`
- **Shopping Lists:** `/api/shopping-list-entry/`
- **And more...**

## 🔒 Security

- **Token-based authentication** with automatic cleanup
- **Password hashing** using bcrypt
- **Access tokens expire** after 5 years (Tandoor compatible)
- **Space-based access control** for multi-user environments

## 🏗️ Architecture

### 🛠️ Tech Stack

- **Backend:** Go 1.21+ with Gin web framework
- **Database:** SQLite with GORM ORM
- **Authentication:** JWT tokens with bcrypt hashing
- **Containerization:** Docker with multi-stage builds

### 📊 Performance

- **Lightning fast** startup and response times
- **Minimal memory footprint** (~10MB container size)
- **SQLite database** - no external dependencies
- **Optimized for** self-hosting and small deployments

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- **Original Tandoor Recipes** project for the amazing API design
- **Go community** for excellent libraries and tooling
- **Docker** for making deployment effortless

---

<p align="center">
  <b>Made with ❤️ in Go</b><br>
  <i>A lightweight alternative to Tandoor Recipes</i>
</p>
