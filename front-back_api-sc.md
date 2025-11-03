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
    "consigner_id": "int",
    "recipient_id": "int",
    "from_station": "long",
    "to_station": "long",
    "deadline": "iso8601",
    "max_price": "decimal"
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
  "cargo_type": "int",
  "worth": "int",
  "cargo_request_id": "int"
}
````

Http status: 200
````
{
  "id: "int"
}
````

### Создание получателя

Post /recipients

````
{
  "first_name": "string",
  "second_name": "string",
  "third_name": "string",
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

### Получение поездок для заявки
Get /trips/cargo_request/{id}?page_number=int&page_size=int

response

Http status: 200
````
{
  "trips": [
    {
      "id": "uuid",
      "from_station": "long",
      "to_station": "long",
      "started_at": "long, #utc timestamp
      "calculated_end_at": "long",
      "actual_end_at": "long",
      "price": "decimal",
      "status": "string",
      "carrier_id": "int",
      "car_id": "int"
    }
  ]
}
````

### Выбор подходящей поездки

Post /cargo_request/{id}/trip/{id}

Http status: 200

### получение поездки для заявки
Get /cargo_request/{id}/trip?page_number=int&page_size=int

Http status: 200

````
{
  "id": "uuid",
  "from_station": "long",
  "to_station": "long",
  "started_at": "long",
  "calculated_end_at": "long",
  "actual_end_at": "long",
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
  "stationSequencs": [
    {
      "station_id": "long",
      "distance": "int",
      "order_num": "int",
      "arrival_at": "long",
      "departure_time": "long"
    }
  ]
}

````

## Форматы

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
