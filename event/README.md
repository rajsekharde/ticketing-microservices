# Event service

## Database schema

events table:
- id
- venue_id
- name
- time
- duration

venues table:
- id
- name
- address
- capacity

seats table:
- id
- venue_id
- seat_number

event_seats table: (Junction table)
- id
- event_id
- seat_id
- status
- price