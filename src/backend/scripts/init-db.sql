-- Create databases for each service
-- Postgres does not support IF NOT EXISTS for CREATE DATABASE
-- Since this runs only on fresh initialization, we can assume they don't exist

CREATE DATABASE auth_db;
CREATE DATABASE inventory_db;
CREATE DATABASE booking_db;
CREATE DATABASE payment_db;
CREATE DATABASE review_db;
CREATE DATABASE notification_db;

-- In production, use the same database with schema separation
-- This is for development with docker-compose
