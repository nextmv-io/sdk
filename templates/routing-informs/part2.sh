#!/bin/bash

if [ "${APPID}" = "" ]; then
  printf "\nPlease set the APPID environment variable to the value created in part 1.\n\n"
  exit 1
fi

cat << EOF
====== Create a new Application Version ======
Once you are happy with the updates in the current app, create a new version. 
This version represents the binary containing the updates you made.
The following command will be run:
  nextmv app version create \\
    --app-id $APPID \\
    --version-id informs-routing-app-v2 \\
    --name "release candidate v2" \\
    --description "INFORMS workshop v2 app" \\
    --json
Press return to create the v2 version.
EOF
read -r

nextmv app version create \
--app-id "$APPID" \
--version-id informs-routing-app-v2 \
--name "release candidate v2" \
--description "INFORMS workshop v2 app" \
--json

cat << EOF
====== Create a staging instance   ======
We now want to create a staging instance, and assign our v2 version to that instance.
We can then use the staging instance to run further tests against
our production release.
The following command will be run:
  nextmv app instance create \\
    --app-id  $APPID \\
    --version-id informs-routing-app-v2 \\
    --instance-id staging \\
    --name "Staging instance" \\
    --description "Staging instance of informs routing app
Press return to create the staging instance, with our updated version assigned.
EOF
read -r

nextmv app instance create \
--app-id "$APPID" \
--version-id informs-routing-app-v2 \
--instance-id staging \
--name "Staging instance" \
--description "Staging instance of informs routing app"


cat << EOF
====== Perform some runs against the staging instance  ======
We want to perform some runs against the staging instance, to generate some data
that we can use in testing. 
Press return to perform runs.
EOF
read -r

nextmv app run \
    --app-id "$APPID"  \
    --input "data/denv_s.json" \
		--instance staging

nextmv app run \
    --app-id "$APPID"  \
    --input "data/denv_m.json" \
		--instance staging

nextmv app run \
    --app-id "$APPID"  \
    --input "data/denv_l.json" \
		--instance staging

printf "\n ============= Completed ================"
printf "\n\nIt should now be possible to run experiments on application with ID:"
printf "\n$APPID\n"
