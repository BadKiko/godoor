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

## Running

```bash
go mod tidy
go build -o godoor .
./godoor
```

Server starts on port 8080.

## Test User

- **Username:** admin
- **Password:** admin123

## API Compatibility

API is fully compatible with original Tandoor Recipes - all responses have identical structure and fields.

## Next Steps

- [ ] Recipe models
- [ ] Recipe API (CRUD)
- [ ] Meal plan models
- [ ] Meal plan API
- [ ] Recipe import/export

## Technologies

- **Go** - main language
- **Gin** - web framework
- **GORM** - ORM for database operations
- **SQLite** - database
- **bcrypt** - password hashing
