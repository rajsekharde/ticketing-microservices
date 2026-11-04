# User service

## gRPC endpoints
- Register user
- Login user
- Get user info

## Database schema

users table:
- id: int, Primary key, auto-incrementing
- email: varchar(255)
- name: varchar(100)
- password_hash: varchar(255)
- role: varchar(20)
- created_at: timestamp, current time as default
