*** Settings ***
Library  RequestsLibrary
Library    Collections
Library   MongoDBLibrary.py    connection_string=mongodb://localhost:27017

*** Variables ***
${Base_URL}  http://localhost:9010
*** Test Cases ***

TC_001_Register
    [Documentation]  Test Register success
    Create Session  Backend_API  ${Base_URL}
    ${uuid}=  Generate Random UUID
    ${username}=    Set Variable  username+${uuid}
    ${name}=    Set Variable  name+${uuid}
    ${response}=  Post Register  ${username}  1234  ${name}
    Should Be Equal As Strings   ${response.status_code}   200
    #check database
    Connect to MongoDB    test
    ${query}=  Create Dictionary  username=${Username}
    ${databaseuser}=  Execute MongoDB Query   users   ${query}
    Should Be Equal As Strings   ${databaseuser[0]}[username]  ${username}
    Should Be Equal As Strings   ${databaseuser[0]}[name]  ${name}

TC_002_Login
    [Documentation]  Test Login success
    Create Session  Backend_API  ${Base_URL}
    ${uuid}=  Generate Random UUID
    ${username}=    Set Variable  username+${uuid}
    ${name}=    Set Variable  name+${uuid}
    ${response}=  Post Register  ${Username}  1234  ${name}
    Should Be Equal As Strings   ${response.status_code}   200
    ${response}=  Post Login  ${Username}  1234
    Should Be Equal As Strings   ${response.status_code}   200

TC_003_GetIdentity
    [Documentation]  Test get identity success
    Create Session  Backend_API  ${Base_URL}
    ${uuid}=  Generate Random UUID 
    ${username}=    Set Variable  username+${uuid}
    ${name}=    Set Variable  name+${uuid}
    ${response}=  Post Register  ${username}  1234  ${name}
    Should Be Equal As Strings   ${response.status_code}   200
    ${response}=  Post Login  ${Username}  1234
    Should Be Equal As Strings   ${response.status_code}   200
    ${token}=   Set Variable   ${response.json()["message"]}
    ${response}=  Get Identity  ${token}
    Should Be Equal As Strings   ${response.status_code}   200
    #check database
    Connect to MongoDB    test
    ${query}=  Create Dictionary  username=${Username}
    ${databaseuser}=  Execute MongoDB Query   users   ${query}
    Should Be Equal As Strings   ${databaseuser[0]}[username]  ${response.json()}[username]
    Should Be Equal As Strings   ${databaseuser[0]}[name]  ${response.json()}[name]
    Should Be Equal As Strings   ${databaseuser[0]}[_id]  ${response.json()}[id]
TC_004_GetUserList
    [Documentation]  Test get user list success
    Create Session  Backend_API  ${Base_URL}
    ${uuid}=  Generate Random UUID 
    ${username}=    Set Variable  username+${uuid}
    ${name}=    Set Variable  name+${uuid}
    ${response}=  Post Register  ${username}  1234  ${name}
    Should Be Equal As Strings   ${response.status_code}   200
    ${response}=  Post Login  ${Username}  1234
    Should Be Equal As Strings   ${response.status_code}   200
    ${token}=   Set Variable   ${response.json()["message"]}
    ${response}=  Get User List  ${token}
    Should Be Equal As Strings   ${response.status_code}   200

TC_005_CreateConversation
    [Documentation]  Test new conversation success
    Create Session  Backend_API  ${Base_URL}
    ${uuid}=  Generate Random UUID 
    ${username}=    Set Variable  username+${uuid}
    ${name}=    Set Variable  name+${uuid}
    ${response}=  Post Register  ${username}  1234  ${name}
    Should Be Equal As Strings   ${response.status_code}   200
    ${response}=  Post Login  ${Username}  1234
    Should Be Equal As Strings   ${response.status_code}   200
    ${token}=   Set Variable   ${response.json()["message"]}
    ${response}=  Get User List  ${token}
    Should Be Equal As Strings   ${response.status_code}   200              
    ${user_list}=  Set Variable  ${response.json()}
    ${user_ids}=    Create List
    FOR    ${user}    IN    @{user_list}
        ${user_id}=    Set Variable    ${user['id']}
        Append To List    ${user_ids}    ${user_id}
    END
    ${user_ids}=  Evaluate   json.dumps(${user_ids})
    ${conversationname}=    Set Variable  conversation_name+${uuid}
    ${response}=  Post Conversation   ${token}  ${conversationname}  ${user_ids}
    Should Be Equal As Strings   ${response.status_code}   200
    #check database
    Connect to MongoDB    test
    ${query}=  Create Dictionary  name=${conversationname}
    ${dataconversation}=  Execute MongoDB Query   conversations   ${query}
    Should Be Equal As Strings   ${dataconversation[0]}[name]  ${response.json()}[name]
    Should Be Equal As Strings   ${dataconversation[0]}[_id]  ${response.json()}[id]


TC_007_CreateMessage
    [Documentation]  Test get conversations success
    Create Session  Backend_API  ${Base_URL}
      ${uuid}=  Generate Random UUID 
    ${username}=    Set Variable  username+${uuid}
    ${name}=    Set Variable  name+${uuid}
    ${response}=  Post Register  ${username}  1234  ${name}
    Should Be Equal As Strings   ${response.status_code}   200
    ${response}=  Post Login  ${Username}  1234
    Should Be Equal As Strings   ${response.status_code}   200
    ${token}=   Set Variable   ${response.json()["message"]}
    ${response}=  Get User List  ${token}
    Should Be Equal As Strings   ${response.status_code}   200 
    ${user_list}=  Set Variable  ${response.json()}
    ${user_ids}=    Create List
    FOR    ${user}    IN    @{user_list}
        ${user_id}=    Set Variable    ${user['id']}
        Append To List    ${user_ids}    ${user_id}
    END
    ${user_ids}=  Evaluate   json.dumps(${user_ids})
    ${response}=  Post Conversation   ${token}  test conversation  ${user_ids}
    Should Be Equal As Strings   ${response.status_code}   200
    ${conversation_id}=  Set Variable  ${response.json()["id"]}
    ${response}=  Create Message  ${token}   ${conversation_id}  "test messages"
    Should Be Equal As Strings   ${response.status_code}   200

*** Keywords ***
Generate Random UUID
    ${uuid}=    Evaluate    str(uuid.uuid4())    modules=uuid
    [Return]    ${uuid}

Post Register
    [Arguments]  ${username}  ${password}  ${name}
    ${response}=  POST on session  Backend_API  /register  data={"name":"${name}", "username": "${username}", "password": "${password}"} 
    [Return]     ${response}

Post Login
    [Arguments]  ${username}  ${password}
    ${response}=  POST on session  Backend_API  /login  data={"username": "${username}", "password": "${password}"} 
    [Return]  ${response}
    
Get Identity
    [Arguments]  ${token}
    ${headers}=       Create Dictionary   Authorization=Bearer ${token}
    ${response}=  GET on session  Backend_API  /identity  headers=${headers}
    [Return]   ${response}

Get User List
    [Arguments]  ${token}
    ${headers}=       Create Dictionary   Authorization=Bearer ${token}
    ${response}=  GET on session  Backend_API   url=/users  headers=${headers}
    [Return]   ${response}

Post Conversation
    [Arguments]  ${token}  ${name}  ${user_ids}
    ${headers}=  Create Dictionary   Authorization=Bearer ${token}
    ${request}=  Set Variable   {"name":"${name}","members":${user_ids}}
    ${response}=  POST on session  Backend_API  /conversations  data=${request}    headers=${headers}
    [Return]  ${response}

Get Conversation
    [Arguments]  ${token} 
    ${headers}=  Create Dictionary   Authorization=Bearer ${token}
    ${response}=  GET on session  Backend_API  /conversations  headers=${headers}
    [Return]  ${response}

Create Message
    [Arguments]  ${token}  ${conversation_id}   ${content} 
    ${headers}=  Create Dictionary   Authorization=Bearer ${token}
    ${request}=  Set Variable   {"conversationId":"${conversation_id}","content":${content}}
    ${response}=  POST on session  Backend_API  /messages   data=${request}  headers=${headers}
    [Return]  ${response}