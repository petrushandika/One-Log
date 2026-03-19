#!/bin/bash

echo "🔍 ULAM Troubleshooting Script"
echo "================================"

# Check if backend is running
echo ""
echo "1. Checking Backend Server..."
if lsof -Pi :8080 -sTCP:LISTEN -t >/dev/null ; then
    echo "   ✅ Backend is running on port 8080"
else
    echo "   ❌ Backend is NOT running"
    echo "   💡 Start with: cd Backend && go run cmd/api/main.go"
fi

# Check if frontend is running
echo ""
echo "2. Checking Frontend Server..."
if lsof -Pi :5173 -sTCP:LISTEN -t >/dev/null ; then
    echo "   ✅ Frontend is running on port 5173"
else
    echo "   ❌ Frontend is NOT running"
    echo "   💡 Start with: cd Frontend && npm run dev"
fi

# Check database connection
echo ""
echo "3. Checking Database..."
if command -v psql &> /dev/null; then
    if PGPASSWORD=${DB_PASSWORD:-postgres} psql -h ${DB_HOST:-localhost} -U ${DB_USER:-postgres} -d ${DB_NAME:-ulam} -c "SELECT 1;" >/dev/null 2>&1; then
        echo "   ✅ Database connection successful"
        
        # Check if tables exist
        echo ""
        echo "4. Checking Tables..."
        TABLES=("activity_feeds" "compliance_exports" "status_page_configs" "sources" "log_entries")
        for table in "${TABLES[@]}"; do
            if PGPASSWORD=${DB_PASSWORD:-postgres} psql -h ${DB_HOST:-localhost} -U ${DB_USER:-postgres} -d ${DB_NAME:-ulam} -c "SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = '$table');" 2>/dev/null | grep -q "t"; then
                echo "   ✅ Table '$table' exists"
            else
                echo "   ❌ Table '$table' MISSING"
            fi
        done
    else
        echo "   ❌ Cannot connect to database"
        echo "   💡 Check your .env file for DB_HOST, DB_USER, DB_PASSWORD, DB_NAME"
    fi
else
    echo "   ⚠️  psql not installed, skipping database check"
fi

echo ""
echo "5. Quick Fixes:"
echo "   - Clear browser cache & cookies"
echo "   - Logout and login again"
echo "   - Run: cd Backend && goose up"
echo ""

# Test API endpoints
echo "6. Testing API Endpoints..."

# Test health endpoint
HEALTH=$(curl -s http://localhost:8080/health 2>/dev/null | grep -o '"status":"success"' || echo "FAILED")
if [ "$HEALTH" = '"status":"success"' ]; then
    echo "   ✅ /health endpoint working"
else
    echo "   ❌ /health endpoint FAILED"
fi

echo ""
echo "📋 Common Issues & Solutions:"
echo ""
echo "Issue: Activity Feed 401 Error"
echo "Solution: JWT token expired. Logout and login again."
echo ""
echo "Issue: Config Page Not Found"
echo "Solution: Run 'goose up' to create missing tables"
echo ""
echo "Issue: Database connection failed"
echo "Solution: Check .env file and ensure PostgreSQL is running"
echo ""
