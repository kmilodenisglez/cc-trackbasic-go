# trackdemo

Asset
An asset is a digital representation of a real asset in the physical world. An asset records every single state or data change
(f.e. the update of metadata, the transfer of ownership, etc.)

## Sample asset structure
```json
{
"ID": "uuid3",
"assetType": "unknown",
"currentState": "available",
"data": "eyJnZW5lcmljRGV2aWNlIjoibW91c2UiLCJnZW5lcmljSW5mbyI6Im1hcmNhIGxvZ2l0ZWNoIn0=",
"docType": "org.asset",
"location": "21.1,72.1",
"manufacturer": "M1",
"owner": "x509::CN=User1@org1.example.com,OU=client,L=San Francisco,ST=California,C=US::CN=ca.org1.example.com,O=org1.example.com,L=San Francisco,ST=California,C=US",
"publicDescription": "",
}
```
