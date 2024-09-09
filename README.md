# ChaosFiles

ChaosFiles is a scalable and secure file sharing platform built using modern web technologies and AWS services. It allows users to easily upload, manage, and share files with others.

## Features

- User authentication with Google Sign-In via AWS Cognito
- Secure file upload and download with pre-signed URLs
- Support for large file uploads using multipart upload
- File metadata storage in DynamoDB
- Scalable file storage using Amazon S3
- Serverless backend using AWS Lambda and API Gateway
- React-based responsive frontend

## Architecture

The ChaosFiles architecture leverages several AWS services to provide a scalable and robust file sharing solution:

- Amazon Cognito for user authentication and authorization
- Amazon S3 for file storage
- Amazon DynamoDB for storing file metadata
- AWS Lambda for serverless backend processing
- Amazon API Gateway for RESTful APIs
- Amazon CloudFront for content delivery

## Architecture Diagram

```mermaid
graph LR

subgraph Frontend
A[React.js + Vite + TypeScript] --> B[S3 + CloudFront]
end

subgraph Backend
C[API Gateway] --> D[Lambda Functions]
end

subgraph Database
F[DynamoDB]
end

subgraph Storage
H[S3]
end

subgraph Authentication
I[Amazon Cognito]
end

A --> I
A --> C
D --> F
D --> H
C --> I
```

## Architecture Flow Diagram

```mermaid
sequenceDiagram
    participant User
    participant Frontend
    participant Cognito
    participant API Gateway
    participant Lambda
    participant DynamoDB
    participant S3

    User->>Frontend: Access ChaosFiles web app
    Frontend->>Cognito: Initiate authentication
    Cognito->>Frontend: Return user token
    Frontend->>API Gateway: Request user files (with token)
    API Gateway->>Lambda: Invoke Lambda function
    Lambda->>Cognito: Verify user token
    Cognito->>Lambda: Token valid
    Lambda->>DynamoDB: Query user files metadata
    DynamoDB->>Lambda: Return files metadata
    Lambda->>S3: Get file preview URLs
    S3->>Lambda: Return preview URLs
    Lambda->>API Gateway: Return files metadata and preview URLs
    API Gateway->>Frontend: Return response
    Frontend->>User: Display files and previews

    User->>Frontend: Upload new file
    Frontend->>API Gateway: Initiate upload (with token)
    API Gateway->>Lambda: Invoke upload Lambda
    Lambda->>S3: Initiate multipart upload
    S3->>Lambda: Return upload ID
    Lambda->>API Gateway: Return upload ID and signed URLs
    API Gateway->>Frontend: Return upload info
    Frontend->>S3: Upload file chunks
    Frontend->>API Gateway: Complete upload
    API Gateway->>Lambda: Invoke completion Lambda
    Lambda->>S3: Complete multipart upload
    Lambda->>DynamoDB: Store file metadata
    Lambda->>API Gateway: Confirm upload complete
    API Gateway->>Frontend: Upload successful
    Frontend->>User: Display uploaded file

    User->>Frontend: Search files
    Frontend->>API Gateway: Send search query (with token)
    API Gateway->>Lambda: Invoke search Lambda
    Lambda->>DynamoDB: Get file metadata for search results
    DynamoDB->>Lambda: Return file metadata
    Lambda->>API Gateway: Return search results with metadata
    API Gateway->>Frontend: Return search results
    Frontend->>User: Display search results

    User->>Frontend: Delete file
    Frontend->>API Gateway: Send delete request (with token and file ID)
    API Gateway->>Lambda: Invoke delete Lambda
    Lambda->>S3: Delete file
    Lambda->>DynamoDB: Delete file metadata
    Lambda->>API Gateway: Confirm file deletion
    API Gateway->>Frontend: Deletion successful
    Frontend->>User: Remove deleted file from view
```

## Getting Started

To use ChaosFiles, simply visit our website at [https://d358wcpg4x8g95.cloudfront.net](https://d358wcpg4x8g95.cloudfront.net).

## Future Work
- Implement use case for network disconnects while uploading
- Add pagination for when users have a large amount of files
- Optimize file upload and download performance (more efficient file chunking strategy, client-side file compression)
- Implement AWS OpenSearch for powerful search
- Enhance file preview feature

## License

This project is licensed under the [MIT License](LICENSE).
