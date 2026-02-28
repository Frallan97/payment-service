#!/usr/bin/env bash
set -uo pipefail

AUTH_URL="${AUTH_URL:-http://localhost:8081}"
PAYMENTS_URL="${PAYMENTS_URL:-http://localhost:8082}"
TEST_EMAIL="${TEST_EMAIL:-integration-test@example.com}"
TEST_NAME="${TEST_NAME:-Integration Test User}"
TEST_EMAIL2="${TEST_EMAIL2:-integration-test2@example.com}"
TEST_NAME2="${TEST_NAME2:-Integration Test User 2}"

PASS=0
FAIL=0
SKIP=0
SECTION_PASS=0
SECTION_FAIL=0
SECTION=""

# --- Helpers ---
# LAST_BODY and LAST_HTTP_CODE are set by authed_curl/unauthed_curl.
# Never capture these functions with $() — it creates a subshell and
# LAST_HTTP_CODE won't propagate. Use LAST_BODY to read the response.
LAST_BODY=""
LAST_HTTP_CODE=""

start_section() {
  if [ -n "$SECTION" ]; then
    echo "  Section result: $SECTION_PASS passed, $SECTION_FAIL failed"
    echo ""
  fi
  SECTION="$1"
  SECTION_PASS=0
  SECTION_FAIL=0
  echo "=== $SECTION ==="
}

pass() {
  echo "  PASS: $1"
  PASS=$((PASS+1))
  SECTION_PASS=$((SECTION_PASS+1))
}

fail() {
  echo "  FAIL: $1"
  FAIL=$((FAIL+1))
  SECTION_FAIL=$((SECTION_FAIL+1))
}

skip() {
  echo "  SKIP: $1"
  SKIP=$((SKIP+1))
}

# Make an authenticated curl request; sets LAST_BODY and LAST_HTTP_CODE
authed_curl() {
  local method="$1" url="$2" body="${3:-}"
  local args=(-s -w '\n%{http_code}' -X "$method" "$url" -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")
  [ -n "$body" ] && args+=(-d "$body")
  local raw
  raw=$(curl "${args[@]}" 2>&1)
  LAST_HTTP_CODE=$(echo "$raw" | tail -n1)
  LAST_BODY=$(echo "$raw" | sed '$d')
}

# Make an unauthenticated curl request; sets LAST_BODY and LAST_HTTP_CODE
unauthed_curl() {
  local method="$1" url="$2" body="${3:-}"
  local args=(-s -w '\n%{http_code}' -X "$method" "$url" -H "Content-Type: application/json")
  [ -n "$body" ] && args+=(-d "$body")
  local raw
  raw=$(curl "${args[@]}" 2>&1)
  LAST_HTTP_CODE=$(echo "$raw" | tail -n1)
  LAST_BODY=$(echo "$raw" | sed '$d')
}

# Assert HTTP status code
expect_status() {
  local expected="$1" label="$2"
  if [ "$LAST_HTTP_CODE" = "$expected" ]; then
    pass "$label (HTTP $LAST_HTTP_CODE)"
  else
    fail "$label — expected HTTP $expected, got $LAST_HTTP_CODE"
  fi
}

echo "============================================"
echo "  Payment Service Integration Tests"
echo "============================================"
echo "Auth:     $AUTH_URL"
echo "Payments: $PAYMENTS_URL"
echo ""

# =============================================
# Auth & Infrastructure
# =============================================
start_section "Auth & Infrastructure"

# 1. Get dev token
RESP=$(curl -sf -X POST "$AUTH_URL/api/dev/token" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL\",\"name\":\"$TEST_NAME\"}" 2>&1) || {
  fail "Could not get dev token (is auth-service running with ENV=development?)"
  echo ""
  echo "=== ABORTED: $PASS passed, $FAIL failed ==="
  exit 1
}
TOKEN=$(echo "$RESP" | jq -r .token)
if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
  fail "Token is empty or null"
  echo "=== ABORTED: $PASS passed, $FAIL failed ==="
  exit 1
fi
pass "Got dev token for $TEST_EMAIL"

# 2. Health check
HEALTH=$(curl -sf "$PAYMENTS_URL/health" 2>&1) || true
if [ -n "$HEALTH" ]; then
  pass "Health check OK"
else
  fail "Health endpoint unreachable"
fi

# 3. Unauthenticated request returns 401
unauthed_curl GET "$PAYMENTS_URL/api/payments"
expect_status "401" "Unauthenticated request returns 401"

# =============================================
# Payment CRUD
# =============================================
start_section "Payment CRUD"

# Create payment (happy path)
authed_curl POST "$PAYMENTS_URL/api/payments" \
  '{"provider":"stripe","amount":10000,"currency":"SEK","description":"Integration test payment"}'
PAYMENT_ID=$(echo "$LAST_BODY" | jq -r .id 2>/dev/null)
expect_status "201" "Create payment"

# Create payment with invalid amount (0)
authed_curl POST "$PAYMENTS_URL/api/payments" \
  '{"provider":"stripe","amount":0,"currency":"SEK","description":"Bad amount"}'
expect_status "400" "Create payment with amount=0 rejected"

# Create payment with missing body
authed_curl POST "$PAYMENTS_URL/api/payments" ''
expect_status "400" "Create payment with empty body rejected"

# Get payment by ID
if [ -n "$PAYMENT_ID" ] && [ "$PAYMENT_ID" != "null" ]; then
  authed_curl GET "$PAYMENTS_URL/api/payments/$PAYMENT_ID"
  expect_status "200" "Get payment by ID"
  GOT_ID=$(echo "$LAST_BODY" | jq -r .id 2>/dev/null)
  if [ "$GOT_ID" = "$PAYMENT_ID" ]; then
    pass "Payment ID matches"
  else
    fail "Payment ID mismatch: expected $PAYMENT_ID, got $GOT_ID"
  fi
else
  fail "Get payment by ID — skipped (no payment ID)"
  fail "Payment ID matches — skipped"
fi

# Get payment with invalid UUID
authed_curl GET "$PAYMENTS_URL/api/payments/not-a-uuid"
expect_status "400" "Get payment with invalid UUID returns 400"

# Get payment with random UUID
authed_curl GET "$PAYMENTS_URL/api/payments/00000000-0000-0000-0000-000000000000"
expect_status "404" "Get payment with random UUID returns 404"

# List payments
authed_curl GET "$PAYMENTS_URL/api/payments?limit=10"
expect_status "200" "List payments"
TOTAL=$(echo "$LAST_BODY" | jq -r .total 2>/dev/null)
if [ -n "$TOTAL" ] && [ "$TOTAL" != "null" ] && [ "$TOTAL" -ge 1 ] 2>/dev/null; then
  pass "List payments total >= 1 (got $TOTAL)"
else
  fail "List payments total should be >= 1, got: $TOTAL"
fi

# List payments with pagination
authed_curl GET "$PAYMENTS_URL/api/payments?limit=1&offset=0"
expect_status "200" "List payments with pagination"
COUNT=$(echo "$LAST_BODY" | jq '.payments | length' 2>/dev/null)
if [ "$COUNT" -le 1 ] 2>/dev/null; then
  pass "Pagination limit=1 respected (got $COUNT items)"
else
  fail "Pagination limit=1 not respected, got $COUNT items"
fi

# =============================================
# Get Customer
# =============================================
start_section "Customer"

authed_curl GET "$PAYMENTS_URL/api/customers/me"
expect_status "200" "Get /api/customers/me returns 200"

# Second user (no payments yet) should get 404
TOKEN2_RESP=$(curl -sf -X POST "$AUTH_URL/api/dev/token" \
  -H "Content-Type: application/json" \
  -d "{\"email\":\"$TEST_EMAIL2\",\"name\":\"$TEST_NAME2\"}" 2>&1) || true
TOKEN2=$(echo "$TOKEN2_RESP" | jq -r .token 2>/dev/null)
if [ -n "$TOKEN2" ] && [ "$TOKEN2" != "null" ]; then
  OLD_TOKEN="$TOKEN"
  TOKEN="$TOKEN2"
  authed_curl GET "$PAYMENTS_URL/api/customers/me"
  expect_status "404" "Second user (no customer) gets 404"
  TOKEN="$OLD_TOKEN"
else
  fail "Could not get token for second user"
fi

# =============================================
# Subscription CRUD
# =============================================
start_section "Subscription CRUD"

# Probe: try creating a subscription to detect if Stripe is configured
authed_curl POST "$PAYMENTS_URL/api/subscriptions" \
  '{"provider":"stripe","amount":5000,"currency":"SEK","interval":"month","interval_count":1,"product_name":"Integration Test Plan","product_description":"Test subscription"}'
SUBSCRIPTION_ID=$(echo "$LAST_BODY" | jq -r .id 2>/dev/null)

if [ "$LAST_HTTP_CODE" = "502" ] || [ "$LAST_HTTP_CODE" = "500" ]; then
  # Stripe not configured — skip CRUD tests, still test validation
  skip "Create subscription (Stripe not configured, got HTTP $LAST_HTTP_CODE)"
  skip "Get subscription by ID (Stripe not configured)"
  skip "List subscriptions (Stripe not configured)"
  skip "Update subscription (Stripe not configured)"
  skip "Cancel subscription (Stripe not configured)"
else
  expect_status "201" "Create subscription"

  # Get subscription by ID
  if [ -n "$SUBSCRIPTION_ID" ] && [ "$SUBSCRIPTION_ID" != "null" ]; then
    authed_curl GET "$PAYMENTS_URL/api/subscriptions/$SUBSCRIPTION_ID"
    expect_status "200" "Get subscription by ID"
  else
    fail "Get subscription by ID — skipped (no subscription ID)"
  fi

  # List subscriptions
  authed_curl GET "$PAYMENTS_URL/api/subscriptions?limit=10"
  expect_status "200" "List subscriptions"
  SUB_TOTAL=$(echo "$LAST_BODY" | jq -r .total 2>/dev/null)
  if [ -n "$SUB_TOTAL" ] && [ "$SUB_TOTAL" != "null" ] && [ "$SUB_TOTAL" -ge 1 ] 2>/dev/null; then
    pass "List subscriptions total >= 1 (got $SUB_TOTAL)"
  else
    fail "List subscriptions total should be >= 1, got: $SUB_TOTAL"
  fi

  # Update subscription
  if [ -n "$SUBSCRIPTION_ID" ] && [ "$SUBSCRIPTION_ID" != "null" ]; then
    authed_curl PATCH "$PAYMENTS_URL/api/subscriptions/$SUBSCRIPTION_ID" \
      '{"cancel_at_period_end":true}'
    expect_status "200" "Update subscription (cancel_at_period_end=true)"
  else
    fail "Update subscription — skipped"
  fi

  # Cancel subscription
  if [ -n "$SUBSCRIPTION_ID" ] && [ "$SUBSCRIPTION_ID" != "null" ]; then
    authed_curl DELETE "$PAYMENTS_URL/api/subscriptions/$SUBSCRIPTION_ID?immediate=true"
    expect_status "200" "Cancel subscription (immediate)"
  else
    fail "Cancel subscription — skipped"
  fi
fi

# Validation tests work regardless of Stripe
authed_curl POST "$PAYMENTS_URL/api/subscriptions" \
  '{"provider":"stripe","amount":5000,"currency":"SEK","interval":"month"}'
expect_status "400" "Create subscription missing product_name rejected"

# =============================================
# Refund Flow
# =============================================
start_section "Refund Flow"

# Create a new payment for refund testing
authed_curl POST "$PAYMENTS_URL/api/payments" \
  '{"provider":"stripe","amount":10000,"currency":"SEK","description":"Payment for refund test"}'
REFUND_PAYMENT_ID=$(echo "$LAST_BODY" | jq -r .id 2>/dev/null)
expect_status "201" "Create payment for refund test"

if [ -n "$REFUND_PAYMENT_ID" ] && [ "$REFUND_PAYMENT_ID" != "null" ]; then
  # Probe: try creating a refund to see if payment is in succeeded state
  # Without Stripe webhooks, payments stay "pending" and refunds are rejected
  authed_curl POST "$PAYMENTS_URL/api/refunds" \
    "{\"payment_id\":\"$REFUND_PAYMENT_ID\",\"amount\":5000,\"reason\":\"Partial refund test\"}"

  if [ "$LAST_HTTP_CODE" = "400" ]; then
    ERROR_CODE=$(echo "$LAST_BODY" | jq -r '.error.code' 2>/dev/null)
    if [ "$ERROR_CODE" = "payment_failed" ]; then
      # Payment is pending (no Stripe webhook to confirm) — skip refund CRUD
      skip "Create partial refund (payment is pending — no Stripe webhook to confirm)"
      skip "List refunds for payment (payment is pending)"
      skip "Create second partial refund (payment is pending)"
      skip "Get refund by ID (payment is pending)"
      # Validation: refunding a pending payment should be rejected
      pass "Refund of pending payment correctly rejected (HTTP 400)"
    else
      # Genuinely bad request for another reason
      fail "Create partial refund (5000 of 10000) — expected HTTP 201, got $LAST_HTTP_CODE"
    fi
  else
    REFUND_ID=$(echo "$LAST_BODY" | jq -r .id 2>/dev/null)
    expect_status "201" "Create partial refund (5000 of 10000)"

    # List refunds for payment
    authed_curl GET "$PAYMENTS_URL/api/payments/$REFUND_PAYMENT_ID/refunds"
    expect_status "200" "List refunds for payment"
    REFUND_COUNT=$(echo "$LAST_BODY" | jq '.refunds | length' 2>/dev/null)
    if [ "$REFUND_COUNT" -ge 1 ] 2>/dev/null; then
      pass "Payment has >= 1 refund (got $REFUND_COUNT)"
    else
      fail "Expected >= 1 refund, got: $REFUND_COUNT"
    fi

    # Second partial refund (other half)
    authed_curl POST "$PAYMENTS_URL/api/refunds" \
      "{\"payment_id\":\"$REFUND_PAYMENT_ID\",\"amount\":5000,\"reason\":\"Second partial refund\"}"
    expect_status "201" "Create second partial refund (5000 of 10000)"

    # Over-refund should fail
    authed_curl POST "$PAYMENTS_URL/api/refunds" \
      "{\"payment_id\":\"$REFUND_PAYMENT_ID\",\"amount\":1,\"reason\":\"Over-refund attempt\"}"
    expect_status "400" "Over-refund rejected"

    # Get refund by ID
    if [ -n "$REFUND_ID" ] && [ "$REFUND_ID" != "null" ]; then
      authed_curl GET "$PAYMENTS_URL/api/refunds/$REFUND_ID"
      expect_status "200" "Get refund by ID"
    else
      fail "Get refund by ID — skipped (no refund ID)"
    fi
  fi
else
  fail "Partial refund — skipped (no payment ID)"
  fail "List refunds — skipped"
  fail "Second partial refund — skipped"
  fail "Over-refund — skipped"
  fail "Get refund by ID — skipped"
fi

# =============================================
# Cross-User Isolation
# =============================================
start_section "Cross-User Isolation"

if [ -n "$TOKEN2" ] && [ "$TOKEN2" != "null" ] && [ -n "$PAYMENT_ID" ] && [ "$PAYMENT_ID" != "null" ]; then
  OLD_TOKEN="$TOKEN"
  TOKEN="$TOKEN2"
  authed_curl GET "$PAYMENTS_URL/api/payments/$PAYMENT_ID"
  expect_status "404" "Second user cannot access first user's payment"
  TOKEN="$OLD_TOKEN"
else
  fail "Cross-user isolation — skipped (missing token or payment ID)"
fi

# =============================================
# Summary
# =============================================
# Print last section result
if [ -n "$SECTION" ]; then
  echo "  Section result: $SECTION_PASS passed, $SECTION_FAIL failed"
fi

echo ""
echo "============================================"
if [ "$SKIP" -gt 0 ]; then
  echo "  TOTAL: $PASS passed, $FAIL failed, $SKIP skipped (no Stripe)"
else
  echo "  TOTAL: $PASS passed, $FAIL failed"
fi
echo "============================================"
[ "$FAIL" -eq 0 ] && exit 0 || exit 1
