## Api грузовиков

### получение заявки (вернутся заявки по id отправителю, id вытащится из jwt токена)
Get /cargo_request/{id}?page_number=int&page_size=int

Http status: 200

````
{
  "cargoRequests": [
    {
      "id": "uuid",
      "recipientId": "int",
      "createdAt": "long",
      "deadline": "long",
      "calculatedTripId": "uuid",
      "actualTripId": "uuid",
      "fromStation": "long",
      "toStation": "long",
      "price": "decimal",
      "status": "string"
    }
  ]
}

````

### создание заявки

Post /cargo_request
````
{
    "consignerId": "int",
    "recipientId": "int",
    "fromStation": "long",
    "toStation": "long",
    "deadline": "iso8601",
    "maxPrice": "decimal"
}
````

Http status: 200
````
{
  "id": "int"
}
````

### создание груза

Post /cargo

````
{
  "length": "int",
  "height": "int",
  "width": "int",
  "weight": "int",
  "cargoType": "int",
  "worth": "int",
  "cargoRequestId": "int"
}
````

Http status: 200
````
{
  "id: "int"
}
````

### создание получателя

Post /recipients

````
{
  "firstname": "string",
  "secondname": "string",
  "thirdname": "string",
  "phone": "string",
  "email": "string"
}
````

Http status: 200
````
{
  "id: "int"
}
````

### получение поездок для заявки
Get /trips/cargo_request/{id}?page_number=int&page_size=int

response

Http status: 200
````
{
  "trips": [
    {
      "id": "uuid",
      "fromStation": "long",
      "toStation": "long",
      "startedAt": "long,
      "calculatedEndAt": "long",
      "actualEndAt": "long",
      "price": "decimal",
      "status": "string",
      "carrierId": "int",
      "carId": "int"
    }
  ]
}
````

### выбор подходящей поездки

Post /cargo_request/{id}/trip/{id}

Http status: 200

### получение поездки для заявки
Get /cargo_request/{id}/trip?page_number=int&page_size=int

Http status: 200

````
{
  "id": "uuid",
  "fromStation": "long",
  "toStation": "long",
  "startedAt": "long",
  "calculatedEndAt": "long",
  "actualEndAt": "long",
  "price": "decimal",
  "status": "string",
  "carrier": "int",
  "car": "int"
}
````

### получение маршрута для поездки
Get /route/trip/{uuid}

Http status: 200

````
{
  "stations": [
    {
      "stationId": "long",
      "distance": "int",
      "ordernum": "int",
      "arrivalAt": "long",
      "departureTime": "long"
    }
  ]
}
````

## Форматы

### timestamp (long)
время возвращается в секундах UTC https://www.unixtimestamp.com

### iso8601
````
2025-11-03T05:16:38+00:00
````

## Ошибки
если пришла ошибка без описания - описания не будет;
http status ошибки, если повезет, 
будет адекватный

HttpStatus: xxx
````
{
  "errorType": "string",
  "message": "ebites kak hotite"
}
````

errorType - константы, сообщим позже
