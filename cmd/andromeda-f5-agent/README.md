# `andromeda-f5-agent` app

## Sample AS3 API request: create/update partition

```sh
PAYLOAD_FILE="$(mktemp)"
trap rm '$PAYLOAD_FILE' EXIT

cat "$PAYLOAD_FILE" <<EOF
{
  "class": "ADC",
  "schemaVersion": "3.0.0",
  "tenant": {
    "class": "Tenant",
    "application": {
      "class": "Application",
      "service": {
        "class": "Service_HTTP",
        "virtualAddresses": [
          "192.0.2.0"
        ],
        "pool": "pool"
      },
      "pool": {
        "class": "Pool",
        "members": [
          {
            "servicePort": 80,
            "serverAddresses": [
              "192.0.2.1",
              "192.0.2.2"
            ]
          }
        ]
      }
    }
  }
}
EOF

F5_USER='...'
F5_PASS='...'
F5_HOST='...'
OS_PROJECT_ID='...'
BASE_URL="/mgmt/shared/appsvcs/declare/${OS_PROJECT_ID}"

curl -X POST -H 'Content-Type: application/json' -d "@${PAYLOAD_FILE}" \
    "https://${F5_USER}:${F5_PASS}@${F5_HOST}${BASE_URL}"
```
