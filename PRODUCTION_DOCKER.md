# 🏭 Production Docker Configuration

Simple production-ready Docker configuration with multi-stage builds for cloud deployment.

## 📁 Production Files

- **`Dockerfile.prod`** - Multi-stage production build (~20MB image)
- **`docker-compose.prod.yml`** - Production services (app, db, nginx, redis)
- **`.env.prod.example`** - Environment variables template
- **`nginx/nginx.conf`** - Nginx reverse proxy with SSL

## 🚀 Quick Usage

### Local Production Testing
```bash
# Copy environment template
cp .env.prod.example .env.prod
# Edit .env.prod with your values

# Build and run
docker build -f Dockerfile.prod -t pandoragym-api:prod .
docker-compose -f docker-compose.prod.yml up -d
```

### Cloud Deployment
Use these files with your cloud provider:
- **AWS ECS/Fargate** - Use `Dockerfile.prod`
- **Google Cloud Run** - Use `Dockerfile.prod`
- **Azure Container Instances** - Use `Dockerfile.prod`
- **Kubernetes** - Use `Dockerfile.prod` + create K8s manifests
- **Docker Swarm** - Use `docker-compose.prod.yml`

## 🔧 Production Features

### Multi-Stage Build
- **Stage 1**: Build with full Go toolchain
- **Stage 2**: Minimal Alpine runtime (~20MB)
- **Security**: Non-root user, static binary

### Services Included
- **App**: Optimized Go binary with health checks
- **Database**: PostgreSQL with persistent storage
- **Nginx**: Reverse proxy with SSL and rate limiting
- **Redis**: Optional caching layer

### Security & Performance
- ✅ Non-root user execution
- ✅ SSL/HTTPS ready
- ✅ Rate limiting
- ✅ Security headers
- ✅ Health checks
- ✅ Resource limits
- ✅ Log rotation

## 🌐 Environment Variables

Copy `.env.prod.example` to `.env.prod` and configure:

```bash
# Required
DATABASE_PASSWORD=your-secure-password
JWT_SECRET=your-jwt-secret-min-32-chars

# Optional
SUPABASE_URL=https://your-project.supabase.co
APP_PORT=3333
```

## 📊 Development vs Production

| Feature | Development | Production |
|---------|-------------|------------|
| **Image Size** | ~800MB | ~20MB |
| **Hot Reload** | ✅ | ❌ |
| **SSL** | ❌ | ✅ |
| **User** | root | non-root |
| **Build** | Single-stage | Multi-stage |

---

**Ready for cloud deployment! 🚀**
