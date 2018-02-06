# API Documentation

[Live Demo](https://shonei-comp3011.herokuapp.com/)

## /addpayment [POST]
It expects a JSON that hols the url that needs to be added.

Example of the body:
```json
{
  "url" : "a valid url"
}
```

## /removepayment [DELETE]
It removes a given url. The url that is to be removed is passed as a query parameter. The parameter needs to be named 'url'.

Example of the request:
```
/removepayment?url=http://www.google.com
```

## /getpayments [GET]
Takes no paramaters and returns an array of all the urls.
