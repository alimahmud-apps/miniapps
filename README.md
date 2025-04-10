## running

1. run the docker compose :
    - `docker-compose up -d`

## List Endpoints

1. POST `localhost:8080/api/users` :
    - create new user

2. GET `localhost:8080/api/users/<userid>`
    - Retreive user by id

3. POST `localhost:8080/api/transactions/credit`:
    - Add amount balance for user

4. POST `localhost:8080/api/transactions/devit`:
    - Deducted amount balance for user

#### detail payload with sample request in postman collection
