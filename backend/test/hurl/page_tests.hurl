# File: page_tests.hurl

# Test 1: Add a new page to a competitor
POST {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages
Authorization: Bearer {{auth_token}}
[
    {
        "url": "https://www.commonroom.io/customers/"
    }
]
HTTP 201
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Added page to competitor successfully"
jsonpath "$.data[0].status" == "active"
jsonpath "$.data[0].url" == "https://www.commonroom.io/customers/"
[Captures]
page_id: jsonpath "$.data[0].id"

# Verify capture profile properties
[Asserts]
jsonpath "$.data[0].capture_profile" exists
jsonpath "$.data[0].capture_profile.capture_beyond_viewport" == true
jsonpath "$.data[0].capture_profile.full_page" == true
jsonpath "$.data[0].capture_profile.format" == "png"
jsonpath "$.data[0].capture_profile.block_ads" == true
jsonpath "$.data[0].capture_profile.wait_until" isCollection

# Test 2: List pages for competitor with pagination
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages?_page=1&_limit=10
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Listed page for competitor successfully"
jsonpath "$.data.hasMore" isBoolean
jsonpath "$.data.pages" isCollection
jsonpath "$.data.pages[*].url" includes "https://www.commonroom.io/customers/"

# Test 3: Get page by ID
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages/{{page_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Fetched page for competitor successfully"
jsonpath "$.data.id" == "{{page_id}}"
jsonpath "$.data.url" == "https://www.commonroom.io/customers/"
jsonpath "$.data.status" == "active"

# Test 4: Get page history
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages/{{page_id}}/history
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Listed page history successfully"
jsonpath "$.data.hasMore" isBoolean
jsonpath "$.data.history" isCollection

# Test 5: Update page URL
PUT {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages/{{page_id}}
Authorization: Bearer {{auth_token}}
{
    "url": "https://www.google.com"
}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Updated page in competitor successfully"
jsonpath "$.data.id" == "{{page_id}}"
jsonpath "$.data.url" == "https://www.google.com"
jsonpath "$.data.status" == "active"

# Test 6: Verify update with get request
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages/{{page_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
jsonpath "$.data.url" == "https://www.google.com"

# Test 7: Remove page
DELETE {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages/{{page_id}}
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Removed page from competitor successfully"

# Test 8: Verify page removal
GET {{host}}/api/public/v1/workspace/{{workspace_id}}/competitors/{{competitor_id}}/pages
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
jsonpath "$.data.pages[*].id" not includes "{{page_id}}"
