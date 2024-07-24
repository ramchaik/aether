#!/bin/bash

# Check for the correct number of arguments
if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <vpc-id> <region>"
  exit 1
fi

VPC_ID=$1
AWS_REGION=$2

echo "Starting cleanup for VPC: $VPC_ID in region: $AWS_REGION"

# Function to delete NAT Gateways
delete_nat_gateways() {
  nat_gateways=$(aws ec2 describe-nat-gateways --filter "Name=vpc-id,Values=$VPC_ID" --query 'NatGateways[*].NatGatewayId' --output text --region "$AWS_REGION")
  for nat in $nat_gateways; do
    echo "Deleting NAT Gateway: $nat"
    aws ec2 delete-nat-gateway --nat-gateway-id $nat --region "$AWS_REGION"
  done

  echo "Waiting for NAT Gateways to be deleted..."
  aws ec2 wait nat-gateway-deleted --filter "Name=vpc-id,Values=$VPC_ID" --region "$AWS_REGION"
}

# Function to release Elastic IPs
release_eips() {
  eips=$(aws ec2 describe-addresses --filter "Name=domain,Values=vpc" --query 'Addresses[*].AllocationId' --output text --region "$AWS_REGION")
  for eip in $eips; do
    echo "Releasing Elastic IP: $eip"
    aws ec2 release-address --allocation-id $eip --region "$AWS_REGION"
  done
}

# Function to detach and delete Internet Gateways
delete_internet_gateways() {
  igws=$(aws ec2 describe-internet-gateways --filters "Name=attachment.vpc-id,Values=$VPC_ID" --query 'InternetGateways[*].InternetGatewayId' --output text --region "$AWS_REGION")
  for igw in $igws; do
    echo "Detaching Internet Gateway: $igw"
    aws ec2 detach-internet-gateway --internet-gateway-id $igw --vpc-id $VPC_ID --region "$AWS_REGION"
    
    echo "Deleting Internet Gateway: $igw"
    aws ec2 delete-internet-gateway --internet-gateway-id $igw --region "$AWS_REGION"
  done
}

# Function to delete subnets
delete_subnets() {
  subnets=$(aws ec2 describe-subnets --filters "Name=vpc-id,Values=$VPC_ID" --query 'Subnets[*].SubnetId' --output text --region "$AWS_REGION")
  for subnet in $subnets; do
    echo "Deleting subnet: $subnet"
    aws ec2 delete-subnet --subnet-id $subnet --region "$AWS_REGION"
  done
}

# Main cleanup process
echo "Deleting NAT Gateways..."
delete_nat_gateways

echo "Releasing Elastic IPs..."
release_eips

echo "Deleting Internet Gateways..."
delete_internet_gateways

echo "Deleting Subnets..."
delete_subnets

# Function to delete VPC
delete_vpc() {
  echo "Deleting VPC: $VPC_ID"
  aws ec2 delete-vpc --vpc-id $VPC_ID --region "$AWS_REGION"
}

delete_vpc

echo "Cleanup process completed for VPC: $VPC_ID in region: $AWS_REGION"
