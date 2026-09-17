-- Creates databases for all the services
-- Runs at Postgres container start up

CREATE DATABASE user_db;
CREATE DATABASE event_db;
CREATE DATABASE booking_db;
CREATE DATABASE payment_db;