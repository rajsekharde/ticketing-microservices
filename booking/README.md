# Booking service

## Database schema

bookings table:
- id
- event_id
- user_id
- status
- booked_at

booking_items table:
- id
- booking_id
- event_seat_id


## Steps for operations

Booking an event ticket
- Accept the request with event_seat id, user id
- Create an entry in bookings table with status = "PENDING"
- Set the status of the event_seat to LOCKED in an atomic operation using the condition WHERE status = "AVAILABLE"
- Send payment request to payment service and wait for a response
- If payment fails or timeouts, set event_seat status to AVAILABLE with WHERE status = "LOCKED" check
- If payment succeeds, set the event_seat status to BOOKED with WHERE status = "LOCKED" check
- Record the booking details in database and set booking status = CONFIRMED or FAILED
- Return response