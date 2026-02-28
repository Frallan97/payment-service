#!/usr/bin/env bash
set -uo pipefail

AUTH_URL="${AUTH_URL:-http://localhost:8080}"
PAYMENTS_URL="${PAYMENTS_URL:-http://localhost:8082}"
TEST_EMAIL="${TEST_EMAIL:-test@example.com}"
TEST_NAME="${TEST_NAME:-Smoke Test User}"

PASS=0
FAIL=0

pass() { echo "  PASS: $1"; PASS=$((PASS+1)); }
fail() { echo "  FAIL: $1"; FAIL=$((FAIL+1)); }

echo "=== Payment Service Smoke Test ==="
echo "Auth:     $AUTH_URL"
echo "Payments: $PAYMENTS_URL"
echo ""

# 1. Get dev token
echo "--- Step 1: Get dev token ---"
TOKEN_RESP=$(curl -sf -X POST "$AUTH_URL/api/dev/token" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"name\":\"$TEST_NAME\"}" 2>&1) || {
  fail "Could not get dev token from $AUTH_URL/api/dev/token (is ENV=development?)"
  echo "Response: $TOKEN_RESP"
  echo ""
  echo "=== Results: $PASS passed, $FAIL failed ==="
  exit 1
}
TOKEN=$(echo "$TOKEN_RESP" | jq -r .token)
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  fail "Token is empty or null"
  echo "$TOKEN_RESP" | jq .
  echo "=== Results: $PASS passed, $FAIL failed ==="
  exit 1
fi
pass "Got dev token"

# 2. Health check
echo "--- Step 2: Health check ---"
HEALTH=$(curl -sf "$PAYMENTS_URL/health" 2>&1) || { fail "Health endpoint unreachable"; }
if [ -n "$HEALTH" ]; then
  pass "Health check OK ($HEALTH)"
else
  fail "Health endpoint returned empty response"
fi

# 3. Unauthenticated request should return 401
echo "--- Step 3: Auth check (expect 401) ---"
HTTP_CODE=$(curl -so /dev/null -w "%{http_code}" "$PAYMENTS_URL/api/payments" 2>&1)
if [ "$HTTP_CODE" = "401" ]; then
  pass "Unauthenticated request returns 401"
else
  fail "Expected 401, got $HTTP_CODE"
fi

# 4. Create a payment
echo "--- Step 4: Create payment ---"
CREATE_RESP=$(curl -s -X POST "$PAYMENTS_URL/api/payments" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"provider":"stripe","amount":10000,"currency":"SEK","description":"Smoke test payment"}' 2>&1)
PAYMENT_ID=$(echo "$CREATE_RESP" | jq -r .id 2>/dev/null)
if [ -n "$PAYMENT_ID" ] && [ "$PAYMENT_ID" != "null" ]; then
  pass "Created payment: $PAYMENT_ID"
else
  fail "Create payment failed: $(echo "$CREATE_RESP" | jq -r .error.message 2>/dev/null || echo "$CREATE_RESP")"
fi

# 5. List payments
echo "--- Step 5: List payments ---"
LIST_RESP=$(curl -s "$PAYMENTS_URL/api/payments?limit=5" \
  -H "Authorization: Bearer $TOKEN" 2>&1)
TOTAL=$(echo "$LIST_RESP" | jq -r .total 2>/dev/null)
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ]; then
  pass "Listed payments (total: $TOTAL)"
else
  fail "Could not parse payment list"
  echo "$LIST_RESP" | jq . 2>/dev/null || echo "$LIST_RESP"
fi

# Summary
echo ""
echo "=== Results: $PASS passed, $FAIL failed ==="
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
