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
- table: <TGT1_TABLE_NAME>:<TGT2_TABLE_NAME> | <TABLE_NAME>
  columns:
    - target: <TGT1_COLUMN_NAME>:<TGT2_COLUMN_NAME> | <COLUMN_NAME>
  join_on:
    - and: <TGT1_JOIN_CONDITION>:<TGT2_JOIN_CONDITION> | <JOIN_CONDITION>
      | or: <TGT1_JOIN_CONDITION>:<TGT2_JOIN_CONDITION> | <JOIN_CONDITION>
  where:
    - and: <TGT1_WHERE_CONDITION>:<TGT2_WHERE_CONDITION> | <WHERE_CONDITION>
      | or: <TGT1_WHERE_CONDITION>:<TGT2_WHERE_CONDITION> | <WHERE_CONDITION>
```
