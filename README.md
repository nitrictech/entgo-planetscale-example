# How and why to use entgo with planetscale

## Discussion

### Why use entgo with planetscale?

[entgo](https://entgo.io) is an ORM that provides a typed API to your DB schema. 
[We](https://nitric.io) have used gorm and were dissatisfied with the lack of typing and controlled schema migrations.

[planetscale](https://planetscale.com) is super easy to setup and has a great free tier whether or not
you are a fan of their extra features like schema branching and merging.

### To use versioned migrations or planetscale schema merges?

Right now we don't use planetscales' schema merges but we are thinking about it.
In the mean time we are using [versioned migrations](https://entgo.io/docs/versioned-migrations), our thinking
behind this is that we will need versioned migrations in our dev branch anyways so do that first.

Another question we have is what happens in the time between upgrading the software deployment and upgrading the schema?
If you do the schema upgrade first, and the old app uses the newer schema you might get failures.
If you do the deployment first then you will have a time of using the old schema.
With both of these cases, your app needs to be either forward or backwards schema compatible.

We are currently running the schema upgrade withing our app so this is happening exactly when it is required.
Pros: no need for schema compatibility, but with the
Con: of 'if the schema upgrade fails, there is downtime'.

## Howto

Steps we are going to take
1. create a repo and basic setup of entgo
1. add a migration
1. setup a planetscale DB
1. run our app

### Prerequisites
- Go
- A free PlanetScale account
- PlanetScale CLI — You can also follow this tutorial in the PlanetScale admin dashboard, but the CLI will make setup quicker.

#### Setup the base project
```bash
cd $GOPATH/src/github.com/<ghuser>
mkdir entgo-planetscale
cd entgo-planetscale
go mod init github.com/<ghuser>/entgo-planetscale
```

#### Use ent to create some entities.
```bash
go install entgo.io/ent/cmd/ent
ent init User
```

#### Change to versioned migrations
```diff
--- a/ent/generate.go
+++ b/ent/generate.go
@@ -1,3 +1,3 @@
 package ent
 
-//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate ./schema
+//go:generate go run -mod=mod entgo.io/ent/cmd/ent generate --feature sql/versioned-migration ./schema
```

Run the generate
```bash
go generate ./...
```

#### Create a migration file

Run a local mysql DB that is empty so we can create migrations against it.
Remember to run this fresh each time you generate migrations.
```bash
docker run --name migration --rm -p 3306:3306 -e MYSQL_ROOT_PASSWORD=pass -e MYSQL_DATABASE=test -d mysql
```

Generate the new migration
```bash
mkdir -p ent/migrate/migrations
go run ./cmd/migrations add-users
```

This will create the following migrations
```diff
diff --git a/ent/migrate/migrations/20221012050944_add-users.down.sql b/ent/migrate/migrations/20221012050944_add-users.down.sql
new file mode 100644
index 0000000..6a8c12c
--- /dev/null
+++ b/ent/migrate/migrations/20221012050944_add-users.down.sql
@@ -0,0 +1,2 @@
+-- reverse: create "users" table
+DROP TABLE `users`;
diff --git a/ent/migrate/migrations/20221012050944_add-users.up.sql b/ent/migrate/migrations/20221012050944_add-users.up.sql
new file mode 100644
index 0000000..ea87419
--- /dev/null
+++ b/ent/migrate/migrations/20221012050944_add-users.up.sql
@@ -0,0 +1,2 @@
+-- create "users" table
+CREATE TABLE `users` (`id` bigint NOT NULL AUTO_INCREMENT, PRIMARY KEY (`id`)) CHARSET utf8mb4 COLLATE utf8mb4_bin;
diff --git a/ent/migrate/migrations/atlas.sum b/ent/migrate/migrations/atlas.sum
new file mode 100644
index 0000000..8d400fc
--- /dev/null
+++ b/ent/migrate/migrations/atlas.sum
@@ -0,0 +1,3 @@
+h1:gMofK6wbvoWIZX3MHz8m96y9UpeW9JopKvXkof65qII=
+20221012050944_add-users.down.sql h1:xM7q8EP/VvWoWKEZEX6DLmTjGwK1B1pImDjbXqXNI+s=
+20221012050944_add-users.up.sql h1:2mXXnpykKV7RIs8kYK0ZM9Y8HtryKRAFcndW0f/6EEY=
```

## setup a planetscale DB

Taken from: https://planetscale.com/docs/tutorials/connect-go-gorm-app

```bash
pscale auth login
pscale database create <DATABASE_NAME> --region <REGION_SLUG>
pscale password create <DATABASE_NAME> <BRANCH_NAME> <PASSWORD_NAME>
```
Take note of the values returned to you, as you won't be able to see this password again.

create a .env file and set DSN to the value have from above. Below is the value for the local docker mysql.
```bash
DSN="root:pass@tcp(localhost:3306)/deploy-test"
```

## Create an example app

*TODO*

## Run the app

*TODO*
