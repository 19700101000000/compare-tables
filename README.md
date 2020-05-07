# Compare Tables

## Goal
compare 2 tables

## Spec

### env.yml
```
driver: <DRIVER_NAME>[:<DRIVER_NAME_2>]
host: <HOST_NAME>[:<HOST_NAME_2>]
port: <PORT_NUMBER>[:<PORT_NUMBER_2>]
username: <USER_NAME>[:<USER_NAME_2>]
password: <USER_PASSWORD>[:<USER_PASSWORD_2>]
database: <DATABASE_NAME>[:<DATABASE_NAME_2>]
[not_same: true | false]
```

#### Driver Name
- mysql
- postgres

### *.yml
```
- table: <TABLE_NAME>[ <OMIT_NAME>][:<TABLE_NAME_2>[ <OMIT_NAME_2>]]
  columns:
    - target: <COLUMN_NAME>[:<COLUMN_NAME_2>]
      [disable_match: true | false]
      [distinct: true | false[:true | false]]
  join_on:
    - and | or: <JOIN_CONDITION>[:<JOIN_CONDITION_2>]
  where:
    - and | or: <WHERE_CONDITION>[:<WHERE_CONDITION_2>]
```
