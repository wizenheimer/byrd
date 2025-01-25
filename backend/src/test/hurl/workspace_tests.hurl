# File: workspace_tests.hurl

# Test 1: Create a new workspace
POST {{host}}/api/public/v1/workspace
Authorization: Bearer {{auth_token}}
{
    "competitors": [
        "https://www.commonroom.io/"
    ],
    "team": [
        "xnayankumar@gmail.com",
        "x4rc3us@gmail.com"
    ],
    "profiles": [
        "product"
    ],
    "features": [
        "page"
    ]
}
HTTP 201
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Created workspace successfully"
jsonpath "$.data.workspace_status" == "active"
[Captures]
workspace_id: jsonpath "$.data.id"

# Test 2: List workspaces
GET {{host}}/api/public/v1/workspace
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Listed workspaces successfully"
jsonpath "$.data[*]" count > 0
jsonpath "$.data[0].workspace_status" == "active"

# Test 3: Get workspace by ID
GET {{host}}/api/public/v1/workspace/{{workspace_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Fetched workspace successfully"
jsonpath "$.data.id" == "{{workspace_id}}"
jsonpath "$.data.workspace_status" == "active"

# Test 4: Update workspace by ID
PUT {{host}}/api/public/v1/workspace/{{workspace_id}}
Authorization: Bearer {{auth_token}}
{
    "billing_email": "billing@byrdhq.com",
    "name": "Byrd's Workspace"
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Updated workspace successfully"
jsonpath "$.data.billingEmail" == "billing@byrdhq.com"
jsonpath "$.data.name" == "Byrd's Workspace"
jsonpath "$.data.workspaceId" == "{{workspace_id}}"

# Test 5: Delete workspace
DELETE {{host}}/api/public/v1/workspace/{{workspace_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Deleted workspace successfully"
jsonpath "$.data.status" == "inactive"
jsonpath "$.data.workspaceId" == "{{workspace_id}}"
