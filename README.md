# GoBlog API


## Starting Application

Run the following command to start your container up 
`docker-compose up`

To run the container in the background you need to use
`docker-compose up -d`

## Services

| Services                      | Container Name    | Exposed Port | Username       | Password |
|-------------------------------|-------------------|--------------|----------------|----------|
| Database                      | goblog-postgres   | 5050         | N/A            | N/A      |
| API                           | blog_app          | 8080         | N/A            | N/A      |
| Database Interfac e (pgAdmin4) | pgadmin_container | 5050         | live@admin.com | password |

## Closing the application

To close the application

## Running Tests

docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
