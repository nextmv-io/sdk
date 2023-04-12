#!/bin/sh

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
The following command will be ran:
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


cat << EOF
====== Create a input set   ======
Input sets are a way to group together previous run inputs that can be used in experiments
to compare two or more versions of an application. Input sets can be created from previous
runs by specifying a list of input IDs, or by searching for a set of inputs ran by an instance, 
optionally using date/time ranges. In this case we will use the most recent runs for the staging instance.
The following command will be ran to create an input set (by default, up to the last 20 runs in 
from the previous day will be used):
nextmv experiment input-set create \\
		--app-id $APPID \\
		--input-set-id denver-sample \\
    --instance-id staging \\
		--name "Denver sample runs" \\
		--description "low, med, and high vol for denver"
Press return to create the staging instance, with our updated version assigned.
EOF
read -r

nextmv experiment input-set create \
		--app-id $APPID \
		--input-set-id denver-sample \
    --instance-id staging \
		--name "Denver sample runs" \
		--description "low, med, and high vol for denver"

  
cat << EOF
====== Run an experiment   ======
Now lets run an experiment comparing the results of our updated version to that of the current
production instance. An experiment runs two versions (represented by instances) against the same
set of inputs, and compares the results.
The following command will be ran to compare the results of production and staging:
nextmv experiment batch start \\
      --experiment-id "compare-with-service-times" \\
      --app-id $APPID \\
			--instance-ids "prod,staging" \\
			--input-set-id denver-sample \\
      --name "Compare v1 prod to v2 staging" \\
			--description "Comparison with service times"
Press return to create a the staging instance, with our updated version assigned.
EOF
read -r

nextmv experiment batch start \
      --experiment-id "compare-with-service-times" \
		  --app-id "$APPID" \
			--instance-ids "prod,staging" \
			--input-set-id denver-sample \
      --name "Compare v1 prod to v2 staging" \
			--description "Comparison with service times"
