# PandoraGym API - Local Services Management

This guide shows how to manage your local PostgreSQL and Go application services using the new Makefile commands.

## 🚀 Quick Commands

### Start All Services
```bash
make local-start
```
**What it does:**
- ✅ Starts PostgreSQL via Homebrew
- ✅ Tests database connection
- ✅ Starts Go application in background
- ✅ Tests API health endpoint
- ✅ Shows service URLs

### Stop All Services
```bash
make local-stop
```
**What it does:**
- 🛑 Stops Go application (all processes)
- 🛑 Kills any processes using port 3333
- 🛑 Stops PostgreSQL service
- 🧹 Cleans up PID files

### Restart All Services
```bash
make local-restart
```
**What it does:**
- 🔄 Runs `local-stop` then `local-start`
- 🔄 Complete clean restart of all services

### Check Service Status
```bash
make local-status
```
**What it shows:**
- 📊 PostgreSQL status (Running/Stopped)
- 📊 Go Application status (Running/Stopped)
- 📊 Database connection status (Connected/Failed)

## 📋 Example Usage

### Daily Development Workflow
```bash
# Start your development environment
make local-start

# Check if everything is running
make local-status

# Work on your code...

# When done for the day
make local-stop
```

### When You Make Changes
```bash
# Restart to pick up changes
make local-restart

# Or just restart the app (keep database running)
pkill -f "go run cmd/server/main.go"
make run
```

### Troubleshooting
```bash
# Check what's running
make local-status

# Force stop everything
make local-stop

# Clean start
make local-start

# View application logs
tail -f app.log
```

## 🔧 Service Details

### PostgreSQL Service
- **Technology:** PostgreSQL 15 via Homebrew
- **Location:** localhost:5432
- **Database:** pandoragym_db
- **User:** pandoragym
- **Password:** password
- **Management:** `brew services start/stop postgresql@15`

### Go Application
- **Technology:** Native Go application
- **Location:** http://localhost:3333
- **Process:** `go run cmd/server/main.go`
- **Logs:** `app.log` file
- **PID:** Stored in `app.pid` file

## 🎯 Service Status Examples

### All Services Running
```
📊 Local Services Status:
=========================
PostgreSQL: ✅ Running
Go Application: ✅ Running (http://localhost:3333)
Database Connection: ✅ Connected
```

### Services Stopped
```
📊 Local Services Status:
=========================
PostgreSQL: ❌ Stopped
Go Application: ❌ Stopped
Database Connection: ❌ Failed
```

### Database Running, App Stopped
```
📊 Local Services Status:
=========================
PostgreSQL: ✅ Running
Go Application: ❌ Stopped
Database Connection: ✅ Connected
```

## 🚨 Troubleshooting

### Port 3333 Already in Use
```bash
# Kill whatever is using port 3333
lsof -ti:3333 | xargs kill -9

# Or use the stop command
make local-stop
```

### Database Connection Issues
```bash
# Check if PostgreSQL is running
brew services list | grep postgresql

# Restart PostgreSQL
brew services restart postgresql@15

# Test connection manually
psql -h localhost -U pandoragym -d pandoragym_db
```

### Application Won't Start
```bash
# Check logs
tail -20 app.log

# Check if database is accessible
make local-status

# Try manual start
go run cmd/server/main.go
```

### Clean Reset
```bash
# Stop everything
make local-stop

# Remove log files
rm -f app.log app.pid

# Start fresh
make local-start
```

## 🔗 Integration with Other Commands

### Database Operations
```bash
# Start services
make local-start

# Run migrations
make migrate-up

# Seed database
make seed

# Check status
make local-status
```

### Development Workflow
```bash
# Start services
make local-start

# Run tests
make test

# Format code
make fmt

# Restart after changes
make local-restart
```

## 📊 Comparison: Local vs Docker

| Aspect | Local Services | Docker |
|--------|---------------|---------|
| **Performance** | ⚡ Faster (native) | 🐌 Slower (containerized) |
| **Resource Usage** | 💾 Lower | 💾 Higher |
| **Setup Complexity** | 🟢 Simple | 🟡 Complex |
| **Debugging** | 🟢 Easy | 🟡 Harder |
| **Portability** | 🟡 Mac-specific | 🟢 Cross-platform |
| **Production Similarity** | 🟡 Different | 🟢 Similar |

## 🎉 Benefits of Local Services

1. **🚀 Faster Development**: No Docker overhead
2. **🔧 Easier Debugging**: Direct access to processes and logs
3. **💾 Better Resource Usage**: No container isolation overhead
4. **🎯 Simpler Workflow**: Native tools and commands
5. **⚡ Quick Restarts**: Instant application restarts

Your PandoraGym API local services are now fully manageable with simple `make` commands! 🚀
