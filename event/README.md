# Event service

## Functions

- Add a venue
    - Create entries for the venue and its seats
- Add an event
    - Create entries for an event and its event_seats
- Get event info
- Get seat info for an event
- Update event_seats

## Database schema

venues table:
- id
- name
- address

seats table:
- id
- venue_id
- seat_number

events table:
- id
- venue_id
- name
- time
- duration

event_seats table: (Junction table)
- id
- event_id
- seat_id
- status
- price

## Steps
- Design database schema
- Determine scope of API for the service
- Create proto file for the gRPC API
- Generate go code from proto file
- Create gRPC server and start listening for connections
- Connect to database
- Initialize db tables
- Create query, grpc endpoint for venue creation
- Create grpc client for event service in gateway
- Create handler, rest endpoint for venue creation in gateway