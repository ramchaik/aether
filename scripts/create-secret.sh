#!/bin/bash

cat << EOF > ./k8s/aws-credentials-secret.yaml
apiVersion: v1
kind: Secret
metadata:
  name: aws-credentials
type: Opaque
stringData:
  AWS_ACCESS_KEY_ID: ${AWS_ACCESS_KEY_ID}
  AWS_SECRET_ACCESS_KEY: ${AWS_SECRET_ACCESS_KEY}
  AWS_SESSION_TOKEN: ${AWS_SESSION_TOKEN}
EOF

echo "Secret YAML file created/updated: ./k8s/aws-credentials-secret.yaml"