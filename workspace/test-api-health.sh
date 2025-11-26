#!/bin/bash

# Test API Health Endpoint
# Usage: ./test-api-health.sh

API_URL="${API_URL:-http://localhost:8088}"
HEALTH_ENDPOINT="${API_URL}/health"
RANKINGS_ENDPOINT="${API_URL}/api/v1/rankings"

echo "🔍 Testing API Health..."
echo "================================"
echo ""

# Test Health Endpoint
echo "1. Testing Health Endpoint: ${HEALTH_ENDPOINT}"
HEALTH_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "${HEALTH_ENDPOINT}" 2>&1)
HTTP_STATUS=$(echo "$HEALTH_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
BODY=$(echo "$HEALTH_RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$HTTP_STATUS" = "200" ]; then
    echo "✅ Health check passed!"
    echo "Response: $BODY"
else
    echo "❌ Health check failed!"
    echo "HTTP Status: $HTTP_STATUS"
    echo "Response: $BODY"
    echo ""
    echo "💡 Make sure the API is running:"
    echo "   cd workspace/apps/api && go run cmd/api/main.go"
    exit 1
fi

echo ""
echo "2. Testing Rankings Endpoint: ${RANKINGS_ENDPOINT}?type=team&category=attacking"
RANKINGS_RESPONSE=$(curl -s -w "\nHTTP_STATUS:%{http_code}" "${RANKINGS_ENDPOINT}?type=team&category=attacking" 2>&1)
RANKINGS_STATUS=$(echo "$RANKINGS_RESPONSE" | grep "HTTP_STATUS" | cut -d: -f2)
RANKINGS_BODY=$(echo "$RANKINGS_RESPONSE" | sed '/HTTP_STATUS/d')

if [ "$RANKINGS_STATUS" = "200" ]; then
    echo "✅ Rankings endpoint working!"
    echo "Response preview:"
    echo "$RANKINGS_BODY" | head -20
else
    echo "❌ Rankings endpoint failed!"
    echo "HTTP Status: $RANKINGS_STATUS"
    echo "Response: $RANKINGS_BODY"
fi

echo ""
echo "================================"
echo "✅ API Health Test Complete!"
