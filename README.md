# Database Initialization

Make sure MySQL is installed. You may refer to [this documentation](https://dev.mysql.com/doc/refman/8.0/en/installing.html) on how to install MySQL locally. The MySQL version used to develop this app is 8.0.36

Create the article database using the following command
```
mysql -u <USER> -p -e 'CREATE DATABASE IF NOT EXISTS article;'
```

I used [golang-migrate](https://github.com/golang-migrate/migrate) as the migration tool. You may follow [this documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) on how to install it.

After installation, run the migration file(s) using the following command
```
migrate -path ./database/migrations -database "mysql://<USER>:<PASSWORD>@tcp(<host>:<port>)/article" -verbose up
```

This should create a table with the following schema
```
Table posts {
  id integer [pk, increment]
  title varchar(200)
  content text [not null]
  category varchar(100) [not null]
  created_at timestamp
  updated_at timestamp
  status enum('Publish', 'Draft', 'Trash') [not null, default: 'Draft']
}
```

However, if you prefer to run the sql query manually, you may use the `./database/migrations/000001_init_mg.up.sql` file.

# Running The Server

I used [Air](https://github.com/cosmtrek/air) for hot reloading during development. If it is not yet installed locally, you may do the following command to install it
```
go install github.com/cosmtrek/air@latest
```

Make sure `.env` file exists. You may follow the format in `.env.example`

Simply run the following command to run a development server
```
air
```

# API Postman Collection

The Postman collection for this service is available [here](https://www.postman.com/security-technologist-46731524/workspace/sharing-vision/collection/29879623-38212276-c836-4e2f-93fb-f9a9b73915aa?action=share&creator=29879623)