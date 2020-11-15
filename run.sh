#!/bin/sh

[ ! -f user_language.csv ] && echo "user,language" > user_language.csv
source /app/.env
./zebra
