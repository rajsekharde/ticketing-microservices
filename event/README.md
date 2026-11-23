# Event service

## Functions

- Add a venue
    - Create entries for the venue and its seats
- Add an event
    - Create entries for an event and its event_seats
- Get event info
- Get seat info for an event
- Update event_seats
- Get all venue info
- Get all event info
- Get all seat info of an event

## Database schema

venues table:
- id
- name
- address
- capacity

seats table:
- id
- venue_id
- seat_number

events table:
- id
- name
- venue_id
- start_time
- duration
- ticket_price

event_seats table: (Junction table)
- id
- event_id
- seat_id
- status
- ticket_price

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
- For adding a new operation, like get all events:
    - create request, response messages, rpc function in proto file
    - generate code from proto file
    - create DB SQL query
    - create gRPC endpoint in server, call query function
    - create gRPC function in client, call rpc function
    - create REST endpoint, call gRPC function
    - expose REST endpoint