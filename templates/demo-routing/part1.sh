#!/bin/bash

## SETUP

RAND="$RANDOM"

if [ "${NEXTMV_API_KEY}" = "" ]; then
  echo "You must set the NEXTMV_API_KEY environment variable to your API key."
  exit 1
fi

printf "\nEnter your initials (One to three letters, lowercase letters only): "
read -r ID

APPID="demo-routing-app-$ID-$RAND"

## CREATE AN APP

cat << EOF
====== Create An Application ======
This script will create an app for you to use in the workshop. 
An application acts as a container for all things related to the problem you
are solving, such as binaries, versions, instances, and runs. 
The following command will be run:
  nextmv app create \\
    --app-id "$APPID"  \\
    --name "Farm share delivery app - $ID" \\
    --description "Routing app for user $ID" \\
    --json;
Press return to create your application.
EOF
read -r

nextmv app create \
--app-id "$APPID"  \
--name "Farm share delivery app - $ID" \
--description "Routing app for user $ID" \
--json;

cat << EOF
====== Push the Application Binary ======
In order to run your app remotely, you'll
need to first push it to Nextmv Cloud.
The following command will be run:
  nextmv app push --app-id "$APPID"
Press return to push your app binary.
EOF
read -r

nextmv app push --app-id "$APPID"

cat << EOF
====== Create an Application Version ======
Once you are happy with the current app you are working with in the cloud, 
you can create a version. The version represents a specific executable app 
binary (e.g., to use for testing or to run in production).
The following command will be run:
  nextmv app version create \\
    --app-id "$APPID"  \\
    --name "v1 release" \\
    --description "Routing demo v1 app" \\
    --version-id demo-routing-app-v1 \\
    --json
Press return to create a version.
EOF
read -r

nextmv app version create \
--app-id "$APPID"  \
--name "v1 release" \
--description "Routing demo v1 app" \
--version-id demo-routing-app-v1 \
--json

cat << EOF
====== Create an Application Instance ======
Finally we'll create an application instance. An application instance is a 
representation of a version and optional configuration that you want to use in 
some context. The same version can be used by multiple application instances. 
For example, you might have a configuration that you run for each farm share 
delivery region, like the Northeast, Midwest, etc or one for different 
environments, like staging and production.
The following command will be run:
  nextmv app instance create \\
    --app-id "$APPID" \\
    --version-id demo-routing-app-v1 \\
    --instance-id prod \\
    --name "Production instance" \\
    --description "Production instance of routing demo app" \\
    --json
Press return to create the production instance.
EOF
read -r


nextmv app instance create \
--app-id "$APPID" \
--version-id demo-routing-app-v1 \
--instance-id prod \
--name "Production instance" \
--description "Production instance of routing demo app" \
--json


printf "\n ============= Completed ================"
printf "\n\nYour assigned app ID is %s." "$APPID"

printf "\n\nRun export APPID=%s to set it in your environment.\n"  "$APPID"

cat << "EOF"
Finally, execute the following command to run your application from 
its assigned endpoint:
  curl -X POST \
    -H "Content-Type: application/json" \
    -H "Authorization: Bearer $NEXTMV_API_KEY" \
    -d "{\"input\": $(cat data/denv_s.json)}" \
    "https://api.cloud.nextmv.io/v1/applications/$APPID/runs?instance_id=prod"
EOF

printf "\n\nYour assigned app ID is $APPID \n"