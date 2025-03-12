# Mango

[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/mango)](https://goreportcard.com/report/github.com/yourusername/mango)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![SvelteKit](https://img.shields.io/badge/SvelteKit-FF3E00?style=flat&logo=svelte&logoColor=white)](https://kit.svelte.dev/)
[![Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)](https://golang.org/)
[![JWT](https://img.shields.io/badge/JWT-black?style=flat&logo=JSON%20web%20tokens)](https://jwt.io/)
[![Website](https://img.shields.io/badge/Website-thejoyestboy.com-blue)](https://www.thejoyestboy.com)

A secure, full-stack application for hosting and streaming anime, movies, audiobooks, and manga with user progress tracking and protected authentication.

## 📋 Features

- **Secure Authentication System**
  - JWT-based authentication
  - Protected user registration
  - Role-based access control
  
- **Media Management**
  - Audiobooks, anime, movies, and manga hosting
  - Streaming capabilities
  - Media progress tracking
  
- **Performance Optimized**
  - Go backend for high performance
  - SvelteKit frontend for reactive UI
  - Efficient media streaming

## 🏗️ Architecture

The application follows a modern client-server architecture:

- **Backend**: Go server providing RESTful API endpoints
- **Frontend**: SvelteKit application with reactive components
- **Authentication**: JWT token-based auth with refresh token mechanism
- **Storage**: Media files stored with efficient retrieval system

## 🚀 Getting Started

### Prerequisites

- Go 1.18+
- Node.js 16+
- npm or yarn
- PostgreSQL 13+

### Installation

1. Clone the repository
```bash
git clone https://github.com/yourusername/mango.git
cd mango
```

2. Set up the backend
```bash
cd backend
cp .env.example .env  # Configure your environment variables
go mod download
go run main.go
```

3. Set up the frontend
```bash
cd frontend
npm install
cp .env.example .env  # Configure your environment variables
npm run dev
```

4. Navigate to `http://localhost:5173` in your browser

## 📁 Project Structure

```
mango/
├── backend/                  # Go server
│   ├── cmd/                  # Application entry points
│   ├── internal/             # Internal packages
│   │   ├── config/           # Configuration
│   │   │   └── config.go
│   │   ├── database/         # Database models and connections
│   │   ├── handlers/         # HTTP handlers
│   │   │   ├── auth.go       # Authentication handlers
│   │   │   ├── content.go    # Content delivery handlers
│   │   │   └── progress.go   # User progress tracking
│   │   └── routes/           # Route definitions
│   │       └── routes.go     # Main router setup
│   └── static/               # Static media files
│
├── frontend/                 # SvelteKit application
│   ├── src/
│   │   ├── lib/              # Reusable components
│   │   │   ├── components/   # UI components
│   │   │   ├── services/     # API services
│   │   │   └── stores/       # Svelte stores
│   │   ├── routes/           # Application routes
│   │   └── app.html          # HTML template
│   ├── static/               # Static assets
│   └── svelte.config.js      # SvelteKit configuration
```

## 🔒 Security Features

- Secure password hashing using bcrypt
- JWT token validation and expiration
- HTTPS enforcement in production
- CORS policy configuration
- Rate limiting to prevent brute force attacks
- Protected registration system

## 🛠️ API Documentation

The backend provides the following RESTful API endpoints:

### Authentication

- `POST /login` - Authenticate and receive JWT
- `POST /logout` - Invalidate current token
- `POST /signup` - Register a new user
- `GET /auth/check` - Verify authentication status

### Content

- `GET /api/content/audiobooks` - List available audiobooks
- `GET /api/content/anime` - List available anime
- `GET /api/content/movies` - List available movies
- `GET /api/content/mangas` - List available manga
- `GET /api/content/audiobooks/{id}` - Get audiobook details

### User Data

- `GET /api/progress/{id}` - Get user progress for a specific media item
- `POST /api/progress` - Save user progress for media consumption
- `GET /user/profile` - Get current user profile information

### Static Content

- `GET /static/*` - Access static media files (with CORS handling)

## 📊 Database Schema

The application's database includes the following main tables:

- `users` - User accounts and authentication data
- `audiobooks` - Audiobook metadata
- `anime` - Anime series and episodes
- `movies` - Movie metadata
- `mangas` - Manga series and chapters
- `progress` - User progress for media consumption

## 🧪 Testing

### Backend Testing
```bash
cd backend
go test ./... -v
```

### Frontend Testing
```bash
cd frontend
npm run test
```

## 🚢 Deployment

The application can be deployed using Docker:

```bash
docker-compose up -d
```

For production deployment, consider:
- Using a reverse proxy like Nginx
- Setting up SSL certificates
- Configuring a CDN for media delivery
- Implementing proper backup procedures

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 📧 Contact

Project Link: [https://github.com/yourusername/mediavault](https://github.com/yourusername/mediavault)
