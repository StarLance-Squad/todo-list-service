# Database Migrations with `golang-migrate/migrate`

## Overview

This document provides quick instructions on how to use `golang-migrate/migrate` for managing database migrations. For
comprehensive documentation, visit
the [official GitHub repository](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

## Basic Commands

### Applying All Available Migrations (Up)

To apply all available migrations to your database:

```shell
migrate -path migrations -database "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" up
```

This command will sequentially apply all migration files that haven't been applied yet in the `migrations` directory.

### Reverting All Migrations (Down)

To revert all applied migrations (be cautious with this in production):

```shell
migrate -path migrations -database "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" down
```

This command will sequentially revert all applied migrations, potentially resulting in data loss.

## Creating New Migrations

### Adding New Changes to the Database

To create a new set of migrations (up and down) for a database change:

```shell
migrate create -ext sql -dir migrations -seq add_new_feature_to_users_table 
```

This command creates two new files in the `migrations` folder:

- `XXXXXX_add_new_feature_to_users_table.up.sql` for the new changes.
- `XXXXXX_add_new_feature_to_users_table.down.sql` for reverting the changes.

Replace `add_new_feature_to_users_table` with a descriptive name for your migration.

### Writing Migration Scripts

- Edit the `.up.sql` and `.down.sql` files to include your SQL commands.
- Ensure that the `.down.sql` script accurately undoes the changes made in the `.up.sql` script.

## Advanced Usage

### Reverting the Last Applied Migration

To revert only the most recently applied migration:

```shell
migrate -path migrations -database "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" down 1
```

The number `1` indicates the number of last migrations to revert. Change this number to revert multiple migrations.

### Checking Current Migration Version

To check the current version of the database schema:

```shell
migrate -path migrations -database "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" version
```

### Forcing a Migration Version

In case of a failed migration leaving the database in a 'dirty' state:

```shell
migrate -path migrations -database "postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable" force VERSION
```

Replace `VERSION` with the version number to which you want to reset.

## Best Practices

- Always backup your database before running migrations, especially in production environments.
- Test migrations in a staging environment before applying them to production.
- Keep migration scripts in version control for tracking and collaboration.
- Document any complex migrations or special cases for future reference.
