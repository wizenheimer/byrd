# File: user_and_token_tests.hurl

# Test 1: Get Token Validation
GET {{host}}/api/public/v1/token
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Clerk token is valid"

# Test 2: Get User Account
GET {{host}}/api/public/v1/users
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"

# Test 3: Delete User Account (Note: This should typically be the last test run)
DELETE {{host}}/api/public/v1/users
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "User deleted successfully"

# Optional: Verify account deletion by attempting to get user info
GET {{host}}/api/public/v1/users
Authorization: Bearer {{auth_token}}
HTTP 404
