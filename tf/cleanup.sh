#!/bin/bash

VPC_ID=$1

if [ -z "$VPC_ID" ]; then
  echo "Usage: $0 <vpc-id>"
  exit 1
fi

echo "Starting cleanup for VPC: $VPC_ID"

# Function to delete NAT Gateways
delete_nat_gateways() {
  nat_gateways=$(aws ec2 describe-nat-gateways --filter "Name=vpc-id,Values=$VPC_ID" --query 'NatGateways[*].NatGatewayId' --output text)
  for nat in $nat_gateways; do
    echo "Deleting NAT Gateway: $nat"
    aws ec2 delete-nat-gateway --nat-gateway-id $nat
  done

  echo "Waiting for NAT Gateways to be deleted..."
  aws ec2 wait nat-gateway-deleted --filter "Name=vpc-id,Values=$VPC_ID"
}

# Function to release Elastic IPs
release_eips() {
  eips=$(aws ec2 describe-addresses --filter "Name=domain,Values=vpc" --query 'Addresses[*].AllocationId' --output text)
  for eip in $eips; do
    echo "Releasing Elastic IP: $eip"
    aws ec2 release-address --allocation-id $eip
  done
}

# Function to detach and delete Internet Gateways
delete_internet_gateways() {
  igws=$(aws ec2 describe-internet-gateways --filters "Name=attachment.vpc-id,Values=$VPC_ID" --query 'InternetGateways[*].InternetGatewayId' --output text)
  for igw in $igws; do
    echo "Detaching and deleting Internet Gateway: $igw"
    aws ec2 detach-internet-gateway --internet-gateway-id $igw --vpc-id $VPC_ID
    aws ec2 delete-internet-gateway --internet-gateway-id $igw
  done
}

# Function to delete subnets
delete_subnets() {
  subnets=$(aws ec2 describe-subnets --filters "Name=vpc-id,Values=$VPC_ID" --query 'Subnets[*].SubnetId' --output text)
  for subnet in $subnets; do
    echo "Deleting subnet: $subnet"
    aws ec2 delete-subnet --subnet-id $subnet
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

echo "Cleanup process completed for VPC: $VPC_ID"