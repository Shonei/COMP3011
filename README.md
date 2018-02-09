# API Documentation

### URL: https://shonei-comp3011.herokuapp.com

## /addpayment [POST]
Adds a new payment method. The new URL needs to be passed as a parameter and be named 'url'.

Example of the request:
```
/addpayment?url=http://www.google.com
```

## /removepayment [DELETE]
It removes a given url. The url that is to be removed is passed as a query parameter. The parameter needs to be named 'url'.

Example of the request:
```
/removepayment?url=http://www.google.com
```

## /getpayments [GET]
Takes no paramaters and returns an array of all the urls.
