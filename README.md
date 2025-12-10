# GoDoor - Tandoor Recipes on Go

Simple implementation of Tandoor recipe management system in Go.

**Version: 2.3.6** (compatible with Tandoor 2.3.6)

## Implemented Features

### ✅ Authentication
- `POST /api-token-auth/` - user authentication (OAuth2 tokens)
- `GET /api-auth/` - DRF authentication endpoints
- Token-based authentication with scopes
- Password hashing with bcrypt

### ✅ Users and Spaces
- `GET /api/user/` - list of users in space
- `GET /api/user/:id/` - get specific user
- `PATCH /api/user/:id/` - update user
- `GET /api/space/` - list of user's spaces
- `POST /api/space/` - create space
- `GET /api/space/:id/` - get specific space
- `PUT/PATCH /api/space/:id/` - update space
- `GET /api/space/current/` - current active space
- `GET /api/user-space/` - user-space relationships

### ✅ Server Settings
- `GET /api/server-settings/current/` - public server settings (no auth required)
- `GET /api/` - root API endpoint returns 403 (server identification)

### ✅ User Preferences
- `GET /api/user-preference/` - get user preferences
- `PATCH /api/user-preference/` - update user preferences

### ✅ Recipe Books
- `GET /api/recipe-book/` - list of recipe books
- `POST /api/recipe-book/` - create recipe book
- `GET /api/recipe-book/:id/` - get specific recipe book
- `PUT/PATCH /api/recipe-book/:id/` - update recipe book
- `DELETE /api/recipe-book/:id/` - delete recipe book

### ✅ Keywords
- `GET /api/keyword/` - list of keywords
- `POST /api/keyword/` - create keyword
- `GET /api/keyword/:id/` - get specific keyword
- `PUT/PATCH /api/keyword/:id/` - update keyword
- `DELETE /api/keyword/:id/` - delete keyword

### ✅ Meal Types
- `GET /api/meal-type/` - list of meal types
- `POST /api/meal-type/` - create meal type

### ✅ Meal Plans
- `GET /api/meal-plan/` - list of meal plans (with date filtering)
- `POST /api/meal-plan/` - create meal plan
- `GET /api/meal-plan/:id/` - get specific meal plan
- `PUT/PATCH /api/meal-plan/:id/` - update meal plan
- `DELETE /api/meal-plan/:id/` - delete meal plan

### ✅ Recipes
- `GET /api/recipe/` - list of recipes (with sorting and pagination)
- `POST /api/recipe/` - create recipe
- `GET /api/recipe/:id/` - get specific recipe
- `PUT/PATCH /api/recipe/:id/` - update recipe
- `DELETE /api/recipe/:id/` - delete recipe
- `PUT /api/recipe/:id/image/` - upload/update recipe image (multipart/form-data: image file or image_url)

## Running

```bash
# Copy environment file
cp env.example .env

# Edit .env file if needed
# nano .env

# To start with clean database, remove the database file:
# rm godoor.db

# Run
go mod tidy
go build -ldflags="-s -w" -o godoor .
./godoor
```

Server starts on port 8080 (configurable via PORT environment variable).

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

## API Compatibility

API is fully compatible with original Tandoor Recipes - all responses have identical structure and fields.

## Security

- Access tokens expire after 5 years (matching Tandoor)
- Expired tokens are automatically cleaned up on server startup
- Invalid/expired tokens return appropriate error codes

## Features Status

- ✅ Authentication and user management
- ✅ Spaces and user-space relationships
- ✅ User preferences
- ✅ Recipe books and entries
- ✅ Keywords
- ✅ Meal types and meal plans
- ✅ Recipes with steps and ingredients
- ✅ Basic image support (URL-based)
- ✅ Basic shopping list support
- ✅ Server settings
- ✅ Recipe image upload endpoint

## Technologies

- **Go** - main language
- **Gin** - web framework
- **GORM** - ORM for database operations
- **SQLite** - database
- **bcrypt** - password hashing
- **godotenv** - environment variables support
