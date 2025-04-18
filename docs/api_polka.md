#Webhook Endpoint for Polka

Exposes webhook endpoint for use by polka payment processor.

All requests from polka expected to have the following header field:
Authorization: Apikey string_representing_api_key

**POST /api/polka/webhook**
Expects json:
```
{
    "event": "event string"
    "data" : {
        _key / value based on event_
    }
}
```

Only one event is currently support, `user.upgraded`, which expects the following data payload:
```
{
    "UserID": "userid string"
}
```

Based on UserID and valid API key from polka, will update a user record so `is_chirpy_red` is `true`.
If record exists, returns HTTP\204
If API Key is invalid, returns HTTP\403
If record does not exist, returns HTTP\404
If internal error occurs, returns HTTP\500

Returns will have no data in body.
