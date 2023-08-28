# Nextshift

## Usage

Run app remotely with latest pushed binary using the REST API endpoint (the
instance=devint query parameter tells the service to run the latest dev binary
pushed to the app).

```bash
MY_API_KEY=<your api key>
curl -X POST \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MY_API_KEY" \
  -d "{\"input\": $(cat input.json)}" \
  "https://api.cloud.nextmv.io/v1/applications/nextshift/runs?instance_id=devint"
```

Retrieve results of remote REST API run

```bash
MY_API_KEY=<your api key>
RUN_ID=<the run ID returned by posting the run to the API>
curl -X GET \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $MY_API_KEY" \
  "https://api.cloud.nextmv.io/v1/applications/nextshift/runs/$RUN_ID"
```
