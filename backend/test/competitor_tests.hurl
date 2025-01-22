# File: competitor_tests.hurl

# Test 1: List competitors with pagination and pages included
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors?_page=1&_limit=10&includePages=true
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Listed workspace competitors successfully"
jsonpath "$.data.competitors" isCollection
jsonpath "$.data.hasMore" isBoolean

# Test 2: Create a new competitor
POST {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors
Authorization: Bearer {{auth_token}}
[
    {
        "url": "https://about.scarf.sh/"
    }
]
HTTP 201
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Created competitor successfully"
jsonpath "$.data.name" == "SCARF"
jsonpath "$.data.status" == "active"
[Captures]
competitor_id: jsonpath "$.data.id"

# Test 3: Get competitor by ID including pages
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}?includePages=true
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Retrieved competitor successfully"
jsonpath "$.data.competitor.id" == "{{competitor_id}}"
jsonpath "$.data.competitor.status" == "active"
jsonpath "$.data.pages" isCollection
jsonpath "$.data.competitor.workspace_id" == "{{workspace_id}}"

# Verify page details if they exist
[Asserts]
jsonpath "$.data.pages[0].capture_profile" exists
jsonpath "$.data.pages[0].diff_profile" isCollection
jsonpath "$.data.pages[0].status" == "active"

# Test 4: Update competitor name
PUT {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}
Authorization: Bearer {{auth_token}}
{
    "name": "beta"
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Updated competitor successfully"
jsonpath "$.data.id" == "{{competitor_id}}"
jsonpath "$.data.name" == "beta"
jsonpath "$.data.status" == "active"
jsonpath "$.data.workspace_id" == "{{workspace_id}}"

# Test 5: Verify update with get request
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
jsonpath "$.data.competitor.name" == "beta"

# Test 6: Remove competitor
DELETE {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Removed competitor from workspace successfully"

# Test 7: Verify competitor is removed
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
jsonpath "$.data.competitors[*].competitor.id" not includes "{{competitor_id}}"
