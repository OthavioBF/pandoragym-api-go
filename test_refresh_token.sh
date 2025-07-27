#!/bin/bash

# Test script for refresh token functionality
BASE_URL="http://localhost:3333"

echo "🚀 Testing Refresh Token API"
echo "================================"

# First, let's create a test user (student)
echo "1. Creating test student..."
STUDENT_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/register/student" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Student",
    "email": "test.student@example.com",
    "phone": "+1234567890",
    "password": "password123",
    "born_date": "1995-01-01T00:00:00Z",
    "age": 28,
    "weight": 70.5,
    "objective": "Build muscle",
    "training_frequency": "3x per week",
    "did_bodybuilding": false,
    "medical_condition": null,
    "physical_activity_level": "Intermediate",
    "observations": "Test user for refresh token"
  }')

echo "Student creation response: $STUDENT_RESPONSE"
echo ""

# Now let's authenticate to get tokens
echo "2. Authenticating user..."
AUTH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/session" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test.student@example.com",
    "password": "password123",
    "device_info": "Test Device - Chrome Browser",
    "ip_address": "127.0.0.1"
  }')

echo "Authentication response: $AUTH_RESPONSE"
echo ""

# Extract tokens from response
ACCESS_TOKEN=$(echo $AUTH_RESPONSE | jq -r '.token // empty')
REFRESH_TOKEN=$(echo $AUTH_RESPONSE | jq -r '.refresh_token // empty')

if [ -z "$ACCESS_TOKEN" ] || [ -z "$REFRESH_TOKEN" ]; then
  echo "❌ Failed to get tokens from authentication response"
  exit 1
fi

echo "✅ Got access token: ${ACCESS_TOKEN:0:20}..."
echo "✅ Got refresh token: ${REFRESH_TOKEN:0:20}..."
echo ""

# Test using the access token
echo "3. Testing access token..."
PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/profile" \
  -H "Authorization: Bearer $ACCESS_TOKEN")

echo "Profile response: $PROFILE_RESPONSE"
echo ""

# Test refresh token endpoint
echo "4. Testing refresh token..."
REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$REFRESH_TOKEN\"
  }")

echo "Refresh response: $REFRESH_RESPONSE"
echo ""

# Extract new tokens
NEW_ACCESS_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.access_token // empty')
NEW_REFRESH_TOKEN=$(echo $REFRESH_RESPONSE | jq -r '.refresh_token // empty')

if [ -z "$NEW_ACCESS_TOKEN" ] || [ -z "$NEW_REFRESH_TOKEN" ]; then
  echo "❌ Failed to get new tokens from refresh response"
  exit 1
fi

echo "✅ Got new access token: ${NEW_ACCESS_TOKEN:0:20}..."
echo "✅ Got new refresh token: ${NEW_REFRESH_TOKEN:0:20}..."
echo ""

# Test using the new access token
echo "5. Testing new access token..."
NEW_PROFILE_RESPONSE=$(curl -s -X GET "$BASE_URL/api/users/profile" \
  -H "Authorization: Bearer $NEW_ACCESS_TOKEN")

echo "Profile response with new token: $NEW_PROFILE_RESPONSE"
echo ""

# Test revoking the refresh token
echo "6. Testing token revocation..."
REVOKE_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/revoke" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$NEW_REFRESH_TOKEN\"
  }")

echo "Revoke response: $REVOKE_RESPONSE"
echo ""

# Try to use the revoked refresh token (should fail)
echo "7. Testing revoked refresh token (should fail)..."
REVOKED_REFRESH_RESPONSE=$(curl -s -X POST "$BASE_URL/auth/refresh" \
  -H "Content-Type: application/json" \
  -d "{
    \"refresh_token\": \"$NEW_REFRESH_TOKEN\"
  }")

echo "Revoked refresh response: $REVOKED_REFRESH_RESPONSE"
echo ""

echo "🎉 Refresh token test completed!"
echo "================================"
