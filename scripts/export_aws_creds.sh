#!/bin/bash

# Read AWS credentials from ~/.aws/credentials
AWS_ACCESS_KEY_ID=$(grep -A2 "\[default\]" ~/.aws/credentials | grep "aws_access_key_id" | cut -d "=" -f2 | tr -d '[:space:]')
AWS_SECRET_ACCESS_KEY=$(grep -A2 "\[default\]" ~/.aws/credentials | grep "aws_secret_access_key" | cut -d "=" -f2 | tr -d '[:space:]')
AWS_SESSION_TOKEN=$(grep -A3 "\[default\]" ~/.aws/credentials | grep "aws_session_token" | cut -d "=" -f2 | tr -d '[:space:]')

# Export the credentials as environment variables
export AWS_ACCESS_KEY_ID
export AWS_SECRET_ACCESS_KEY
export AWS_SESSION_TOKEN

echo "AWS credentials exported as environment variables."