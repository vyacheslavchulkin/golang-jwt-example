echo "Run init script"
mongo admin -u "$MONGO_INITDB_ROOT_USERNAME" -p "$MONGO_INITDB_ROOT_PASSWORD" --eval "db.getSiblingDB('$MONGO_APPLICATION_DATABASE').createUser({user: '$MONGO_APPLICATION_USERNAME', pwd: '$MONGO_APPLICATION_PASSWORD', roles: [{role: 'readWrite', db: '$MONGO_APPLICATION_DATABASE'}]});"
echo "End init script"