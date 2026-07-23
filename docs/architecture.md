# Architecture Decisions

- Backend: Go
- Frontend: React
- Database: MongoDB
- Authentication: JWT
- JWT Storage: HttpOnly Cookie
- API: REST
- Initial Deployment: Local network
- Future Deployment: Cloud


React
   │
REST API
   │
Gin
   │
Services
   │
Repositories
   │
MongoDB


Handler
↓

Service

↓

Repository

↓

Database