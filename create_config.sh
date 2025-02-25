#!/bin/sh
echo "Creating config file..."
echo "DB Path: ${DB_PATH}"
cat <<EOF > /app/config.yaml
dbfilepath: ${DB_PATH}
devmode: ${DEV_MODE}
instanceid: ${INSTANCE_ID}
logfilepath: ${LOG_PATH}
port: ${PORT}
secretkey: ${SECRET_KEY}
discordwebhookurl: ${DISCORD_WEBHOOK}
EOF