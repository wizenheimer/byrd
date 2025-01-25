# File: workspace_tests.hurl

# Test 1: Join workspace
POST {{host}}/api/public/v1/workspace/{{workspace_id}}/join
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Joined workspace successfully"
jsonpath "$.data.workspaceId" == "{{workspace_id}}"

# Test 2: Exit workspace
POST {{host}}/api/public/v1/workspace/{{workspace_id}}/exit
Authorization: Bearer {{auth_token}}
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$.message" == "Exited workspace successfully"
jsonpath "$.data.workspaceId" == "{{workspace_id}}"
