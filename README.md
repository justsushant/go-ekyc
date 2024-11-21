# eKYC Exercise

A fully featured REST API for an online KYC (know your customer) system. This will include REST APIs, relational DB, Redis for caching, RabbitMQ as a queue, async workers for processing, etc.


## Usage

1. Set the required environment variables in .env file
2. Run the command below to launch PostgreSQL container
   ```
   docker compose up
   ```
3. Run the command below to apply the necessary migrations
   ```
   make create-migrate
4. Run the command below to build and run the binary:
   ```
   make create-migrate && bin/migrate -m=up
   ```
5. Connect to the server using any HTTP client 

## Current Progress Feat

- Health check endpoint working
- Migrate up implemented successfully


# Bugs
- for invalid access key, im not getting 401
- in getmetadatfromuuid store function, "sql: no rows in result set" is returned. we need some apt error when no record if found for that particular uuid
