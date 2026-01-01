# RentalFlow Deployment Guide

## 🚀 Quick Start with Docker

### Prerequisites
- Docker & Docker Compose installed
- Chapa account (for payments)
- Gmail account or SMTP server (for emails)
- Cloudinary account (for image uploads)

### 1. Clone & Configure

```bash
git clone <repository-url>
cd RentalFlow

# Copy environment template
cp .env.example .env

# Edit .env with your credentials
nano .env
```

### 2. Required Environment Variables

#### Critical (Must Configure):
```bash
# Database
POSTGRES_PASSWORD=<strong-password>

# JWT Secret (min 32 characters)
JWT_SECRET=<your-secret-key>

# Chapa Payment
CHAPA_SECRET_KEY=<from-chapa-dashboard>
CHAPA_PUBLIC_KEY=<from-chapa-dashboard>
CALLBACK_URL=https://yourdomain.com/payment/callback

# Email (Gmail example)
SMTP_USERNAME=your-email@gmail.com
SMTP_PASSWORD=<app-password>

# Cloudinary
CLOUDINARY_CLOUD_NAME=<your-cloud-name>
CLOUDINARY_UPLOAD_PRESET=<your-preset>

# CORS Origins (Comma separated list of frontend domains)
ALLOWED_ORIGINS=https://your-frontend.vercel.app,http://localhost:3000
```

### 3. Start Services

```bash
# Build and start all services
docker-compose up -d

# Check status
docker-compose ps

# View logs
docker-compose logs -f
```

### 4. Run Database Migrations

```bash
# Auth service
docker-compose exec auth-service migrate -path migrations -database "postgresql://rentalflow:password@postgres:5432/rentalflow?sslmode=disable" up

# Review service
docker-compose exec review-service migrate -path migrations -database "postgresql://rentalflow:password@postgres:5432/rentalflow?sslmode=disable" up
```

### 5. Access Application

- **Frontend**: http://localhost:3001
- **Auth API**: http://localhost:8081
- **Inventory API**: http://localhost:8082
- **Booking API**: http://localhost:8083
- **Payment API**: http://localhost:8084
- **Review API**: http://localhost:8085
- **Notification API**: http://localhost:8086

---

## 🔧 Service Configuration Details

### Auth Service (Port 8081)
- Handles user authentication
- JWT token generation
- User profile management
- Required ENV: `JWT_SECRET`, `JWT_EXPIRY`

### Inventory Service (Port 8082)
- Item management (CRUD)
- Search and filtering
- Image URL storage

### Booking Service (Port 8083)
- Booking creation & management
- Status tracking
- Date validation

### Payment Service (Port 8084)
- Chapa integration
- Payment initialization
- Webhook handling
- Required ENV: `CHAPA_SECRET_KEY`, `CHAPA_PUBLIC_KEY`, `CALLBACK_URL`

### Review Service (Port 8085)
- Reviews and ratings
- Review moderation

### Notification Service (Port 8086)
- Email notifications
- SMTP integration
- Required ENV: `SMTP_HOST`, `SMTP_USERNAME`, `SMTP_PASSWORD`

---

## 📧 Email Configuration

### Gmail Setup:
1. Enable 2FA on Gmail account
2. Generate App Password: https://myaccount.google.com/apppasswords
3. Use app password in `SMTP_PASSWORD`

### SendGrid Setup:
```bash
SMTP_HOST=smtp.sendgrid.net
SMTP_PORT=587
SMTP_USERNAME=apikey
SMTP_PASSWORD=<your-sendgrid-api-key>
```

---

## 💳 Chapa Payment Setup

1. Create account: https://chapa.co
2. Get API keys from dashboard
3. Set webhook URL: `https://yourdomain.com/api/payments/webhook`
4. Configure callback URL in `.env`

---

## 🖼️ Cloudinary Setup

1. Create account: https://cloudinary.com
2. Get cloud name from dashboard
3. Create upload preset:
   - Settings → Upload → Add upload preset
   - Mode: Unsigned
   - Folder: rentalflow
4. Add credentials to `.env`

---

## 🔒 Security Checklist

- [ ] Change default PostgreSQL password
- [ ] Generate strong JWT secret (32+ chars)
- [ ] Use app passwords for email
- [ ] Enable HTTPS in production
- [ ] Set secure CORS origins
- [ ] Enable rate limiting
- [ ] Regular security updates

---

## 🌐 Production Deployment

### Option 1: Docker on VPS

```bash
# On your server
git clone <repo>
cd RentalFlow
cp .env.example .env
# Edit .env with production values
docker-compose up -d
```

### Option 2: Kubernetes

```bash
# Create configmaps and secrets
kubectl create secret generic rentalflow-secrets \
  --from-env-file=.env

# Apply deployments
kubectl apply -f k8s/
```

### Option 3: Cloud Platforms

#### AWS ECS/Fargate
- Use docker-compose with ECS CLI
- Configure RDS for PostgreSQL
- Use ElastiCache for Redis

#### Google Cloud Run
- Deploy each service separately
- Use Cloud SQL for database
- Use Memorystore for Redis

#### Digital Ocean App Platform
- Connect GitHub repository
- Auto-deploy from docker-compose.yml
- Use managed PostgreSQL

---

## 🔄 Updating Services

```bash
# Pull latest changes
git pull

# Rebuild and restart
docker-compose up -d --build

# Run new migrations if any
docker-compose exec auth-service migrate ...
```

---

## 🐛 Troubleshooting

### Services won't start
```bash
# Check logs
docker-compose logs <service-name>

# Check database connection
docker-compose exec postgres psql -U rentalflow -d rentalflow
```

### Frontend can't connect to backend
- Check `VITE_API_URL` in your frontend environment variables.
- Verify that your frontend domain is added to `ALLOWED_ORIGINS` in the backend (API Gateway).
- Ensure all services are running.

### Email not sending
- Verify SMTP credentials
- Check Gmail app password
- Review notification service logs

### Payment not working
- Verify Chapa credentials
- Check callback URL is accessible
- Review payment service logs

---

## 📊 Monitoring

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f auth-service

# Since timestamp
docker-compose logs --since 2023-01-01 auth-service
```

### Health checks
```bash
curl http://localhost:8081/health
curl http://localhost:8082/health
# ... for all services
```

---

## 🔄 Backup & Restore

### Backup Database
```bash
docker-compose exec postgres pg_dump -U rentalflow rentalflow > backup.sql
```

### Restore Database
```bash
docker-compose exec -T postgres psql -U rentalflow rentalflow < backup.sql
```

---

## 🎯 Next Steps

1. Configure all environment variables
2. Start services with `docker-compose up -d`
3. Run database migrations
4. Test all features
5. Set up domain and SSL
6. Configure monitoring
7. Deploy to production!

---

## 📞 Support

For issues or questions, check:
- Service logs: `docker-compose logs <service>`
- Database status: `docker-compose exec postgres pg_isready`
- Network connectivity: `docker-compose exec <service> ping postgres`

## 🎉 You're Ready!

RentalFlow is now configured for deployment. Just add your credentials and launch! 🚀
