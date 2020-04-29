# Compare Tables

## Goal
compare 2 tables

## Spec

### env.yml
```
driver: mysql | postgres
host: <HOST_NAME>
port: <PORT_NUMBER>
username: <USER_NAME>
password: <USER_PASSWORD>
database: <DATABASE_NAME>
```

### *.yml
```
- table: <TGT1_TABLE_NAME>:<TGT2_TABLE_NAME> | <TGT_TABLE_NAME>
  columns:
    - target: <TGT1_COLUMN_NAME>:<TGT2_COLUMN_NAME> | <TGT_COLUMN_NAME>
  join_on:
    - and: <JOIN_CONDITION>
      | or: <JOIN_CONDITION>
  where:
    - and: <WHERE_CONDITION>
      | or: <WHERE_CONDITION>
```
