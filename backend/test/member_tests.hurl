# File: member_tests.hurl

# Test 1: List workspace users with pagination
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/users?_page=1&_limit=10
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Listed workspace users successfully"
jsonpath "$.data.users" isCollection
jsonpath "$.data.hasMore" isBoolean

# Optional: Capture a user ID for later tests if the list is not empty
[Captures]
existing_user_id: jsonpath "$.data.users[0].user_id"

# Test 2: Add new users to workspace
POST {{host}}/api/public/v1/workspace/{{workspace_id}}/users
Authorization: Bearer {{auth_token}}
{
    "emails": [
        "test@test.com"
    ]
}
HTTP 201
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Invited users to workspace successfully"
jsonpath "$.data[0].email" == "test@test.com"
jsonpath "$.data[0].workspace_role" == "user"
jsonpath "$.data[0].membership_status" == "pending"
[Captures]
new_user_id: jsonpath "$.data[0].user_id"

# Test 3: Update user role to admin
PUT {{host}}/api/public/v1/workspace/{{workspace_id}}/users/{{new_user_id}}
Authorization: Bearer {{auth_token}}
{
    "role": "admin"
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Updated user role in workspace successfully"
jsonpath "$.data.role" == "admin"

# Test 4: Verify user list after update (with role filter)
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/users?role=admin&_page=1&_limit=10
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Listed workspace users successfully"
jsonpath "$.data.users[*].workspace_role" includes "admin"

# Test 5: Remove user from workspace
DELETE {{host}}/api/public/v1/workspace/{{workspace_id}}/users/{{new_user_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Removed user from workspace successfully"
jsonpath "$.data.userId" == "{{new_user_id}}"
jsonpath "$.data.workspaceId" == "{{workspace_id}}"

# Test 6: Verify user removal
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/users?_page=1&_limit=10
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
# Verify the removed user is not in the list
jsonpath "$.data.users[*].user_id" not includes "{{new_user_id}}"
