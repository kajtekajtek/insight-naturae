# scripts/dump_table.sh - dump a chosen table from the database
#!/bin/sh

# check if the user has provided the table name
if [ $# -eq 0 ]; then 
    echo "Usage: $0 <table_name>";
    exit 1;
fi

TABLE_NAME=$1

if [ -z "$DB_PATH" ]; then
    echo "DB_NAME is not set. Continuing with the default value.";
    DB_PATH="insight-naturae.db";
fi

if [ ! -f $DB_PATH ]; then
    echo "Database file not found. Exiting...";
    exit 1;
fi

sqlite3 "$DB_PATH" <<EOF
.headers on
.mode column
SELECT * FROM $TABLE_NAME;
EOF