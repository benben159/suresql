#!/bin/bash

# manage_users.sh - Script for managing users in the SureSQL database
# Usage:
#   ./manage_users.sh add <username> <password> [<role_name>]
#   ./manage_users.sh update <username> [<new_username>] [<new_password>] [<new_role_name>]
#   ./manage_users.sh delete <username>

# Example:

# List all users
# ./manage_users.sh list

# List users with usernames containing "admin"
# ./manage_users.sh list admin

# Add User
# ./manage_users.sh add johndoe securepass admin

# Update password
# ./manage_users.sh update johndoe --new-password newpass123

# Update username
# ./manage_users.sh update johndoe --new-username john.doe

# Update role
# ./manage_users.sh update johndoe --new-role user

# Update multiple attributes at once
# ./manage_users.sh update johndoe --new-username john.doe --new-password newpass123 --new-role user

# Delete user
# ./manage_users.sh delete johndoe

# Source environment variables
if [ -f ./.env.dev ]; then
    source ./.env.suresql
else
    echo "Error: .env.dev file not found"
    exit 1
fi

if [ -f ./.env.simplehttp ]; then
    source ./.env.simplehttp
else
    echo "Error: .env.simplehttp file not found"
    exit 1
fi

# Check if required environment variables are set
if [ -z "$DBMS_USERNAME" ] || [ -z "$DBMS_PASSWORD" ]; then
    echo "Error: DBMS_USERNAME or DBMS_PASSWORD environment variables not set"
    exit 1
fi

if [ -z "$SURESQL_HOST" ] || [ -z "$SURESQL_PORT" ]; then
    echo "Error: SURESQL_HOST or SURESQL_PORT environment variables not set"
    exit 1
fi

# Set the server URL
SERVER_URL="http://$SURESQL_HOST:$SURESQL_PORT"
# SERVER_URL=${SERVER_URL:-"http://localhost:5130"}

# Create basic auth header
AUTH_HEADER=$(echo -n "$DBMS_USERNAME:$DBMS_PASSWORD" | base64)

# Show usage information
show_usage() {
    echo "Usage:"
    echo "  $0 add <username> <password> [<role_name>]"
    echo "  $0 update <username> [--new-username <new_username>] [--new-password <new_password>] [--new-role <new_role_name>]"
    echo "  $0 delete <username>"
    echo "  $0 list [<username_filter>]"
    exit 1
}

# Function to check HTTP response
check_response() {
    if [[ $1 -ge 200 ]] && [[ $1 -lt 300 ]]; then
        echo "Success: $2"
    else
        echo "Error ($1): $3"
        exit 1
    fi
}

# Parse command line arguments
if [ $# -lt 1 ]; then
    show_usage
fi

ACTION=$1
shift

case $ACTION in
"list")
    USERNAME_FILTER=""
    if [ $# -eq 1 ]; then
        USERNAME_FILTER=$1
        echo "Listing users matching: $USERNAME_FILTER"
    else
        echo "Listing all users"
    fi
    
    # Prepare URL
    REQUEST_URL="${SERVER_URL}${SURESQL_INTERNAL_API}/iusers"
    if [ ! -z "$USERNAME_FILTER" ]; then
        REQUEST_URL="$REQUEST_URL?username=$USERNAME_FILTER"
    fi
    
    # Make API request
    RESPONSE=$(curl -s -w "\n%{http_code}" -X GET "$REQUEST_URL" \
        -H "Authorization: Basic $AUTH_HEADER")
    
    # Extract status code and response body
    HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
    RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')
    echo "$RESPONSE_BODY"
    
    check_response $HTTP_STATUS "Users retrieved successfully" "Retreive users failed"
    
    # Format the JSON response for better readability
    if command -v jq &> /dev/null; then
        # If jq is installed, use it for pretty printing
        echo "$RESPONSE_BODY" | jq '.data[] | {id, username, role_name, created_at}'
    else
        # Otherwise just print the raw response
        echo "$RESPONSE_BODY"
    fi
    ;;

    "add")
        if [ $# -lt 2 ]; then
            echo "Error: 'add' requires at least username and password"
            show_usage
        fi
        
        USERNAME=$1
        PASSWORD=$2
        ROLE_NAME=${3:-""}
        
        echo "Adding user: $USERNAME"
        
        # Prepare JSON data
        JSON_DATA="{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\""
        if [ ! -z "$ROLE_NAME" ]; then
            JSON_DATA="$JSON_DATA,\"role_name\":\"$ROLE_NAME\""
        fi
        JSON_DATA="$JSON_DATA}"
        # echo "DEBUG; data: ${JSON_DATA}"
        # Make API request
        RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "${SERVER_URL}${SURESQL_INTERNAL_API}/iusers" \
            -H "Authorization: Basic $AUTH_HEADER" \
            -H "Content-Type: application/json" \
            -d "$JSON_DATA")
        
        # Extract status code and response body
        HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
        RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')
        
        check_response $HTTP_STATUS "User $USERNAME added successfully" "User $USERNAME insert failed"
        echo "$RESPONSE_BODY"
        ;;
        
    "update")
        if [ $# -lt 1 ]; then
            echo "Error: 'update' requires a username"
            show_usage
        fi
        
        USERNAME=$1
        shift
        
        # Initialize update data
        NEW_USERNAME=""
        NEW_PASSWORD=""
        NEW_ROLE=""
        
        # Parse optional parameters
        while [ $# -gt 0 ]; do
            case "$1" in
                --new-username)
                    NEW_USERNAME="$2"
                    shift 2
                    ;;
                --new-password)
                    NEW_PASSWORD="$2"
                    shift 2
                    ;;
                --new-role)
                    NEW_ROLE="$2"
                    shift 2
                    ;;
                *)
                    echo "Unknown option: $1"
                    show_usage
                    ;;
            esac
        done
        
        # Check if at least one update parameter is provided
        if [ -z "$NEW_USERNAME" ] && [ -z "$NEW_PASSWORD" ] && [ -z "$NEW_ROLE" ]; then
            echo "Error: At least one of --new-username, --new-password, or --new-role must be specified"
            show_usage
        fi
        
        echo "Updating user: $USERNAME"
        
        # Prepare JSON data
        JSON_DATA="{\"username\":\"$USERNAME\""
        
        if [ ! -z "$NEW_USERNAME" ]; then
            JSON_DATA="$JSON_DATA,\"new_username\":\"$NEW_USERNAME\""
        fi
        
        if [ ! -z "$NEW_PASSWORD" ]; then
            JSON_DATA="$JSON_DATA,\"new_password\":\"$NEW_PASSWORD\""
        fi
        
        if [ ! -z "$NEW_ROLE" ]; then
            JSON_DATA="$JSON_DATA,\"new_role_name\":\"$NEW_ROLE\""
        fi
        
        JSON_DATA="$JSON_DATA}"
        
        # Make API request
        RESPONSE=$(curl -s -w "\n%{http_code}" -X PUT "${SERVER_URL}${SURESQL_INTERNAL_API}/iusers" \
            -H "Authorization: Basic $AUTH_HEADER" \
            -H "Content-Type: application/json" \
            -d "$JSON_DATA")
        
        # Extract status code and response body
        HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
        RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')
        
        check_response $HTTP_STATUS "User $USERNAME updated successfully" "User $USERNAME update failed"
        echo "$RESPONSE_BODY"
        ;;
        
    "delete")
        if [ $# -ne 1 ]; then
            echo "Error: 'delete' requires a username"
            show_usage
        fi
        
        USERNAME=$1
        
        echo "Deleting user: $USERNAME"
        
        # Make API request
        RESPONSE=$(curl -s -w "\n%{http_code}" -X DELETE "${SERVER_URL}${SURESQL_INTERNAL_API}/iusers?username=$USERNAME" \
            -H "Authorization: Basic $AUTH_HEADER")
        
        # Extract status code and response body
        HTTP_STATUS=$(echo "$RESPONSE" | tail -n1)
        RESPONSE_BODY=$(echo "$RESPONSE" | sed '$d')
        
        check_response $HTTP_STATUS "User $USERNAME deleted successfully" "User $USERNAME delete failed"
        echo "$RESPONSE_BODY"
        ;;
        
    *)
        echo "Error: Unknown action '$ACTION'"
        show_usage
        ;;
esac

exit 0
