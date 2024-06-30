# X-Clone-Backend
## Building the Environment
To set up the development environment, follow these steps:

1. Build Docker image  

Make sure you have Docker installed on your machine.
```
make build
```

2. Start Docker Containers

Start the Docker containers using docker-compose:
```
make up
```

3. Access the Application Container  

Enter the application container:
```
make exec_app
```

## Checking PostgreSQL State
To check the state of the PostgreSQL database:

1. Access the PostgreSQL Container  

Enter the PostgreSQL container:
```
make exec_pg
```

2. Connect to the PostgreSQL Database  

Once inside the PostgreSQL container, connect to the database:
```
psql
```

3. Check the Database State  

You can now run SQL commands to check the state of the database. For example, to list all tables:
```
\dt
```
To view the contents of a table:
```
SELECT * FROM [table_name];
```

## Migration Guide
To add a new schema migration:

1. Create Migration Files

Run the following command inside the application container to create new migration files:

```
migrate create -ext sql -seq -dir ./db/migrations [migration_name]
```
This will generate two SQL files: one for the up migration and one for the down migration.

2. Edit Migration Files

Write the SQL statements in the generated files. For example, to create a new table, add the SQL to the up migration file.
