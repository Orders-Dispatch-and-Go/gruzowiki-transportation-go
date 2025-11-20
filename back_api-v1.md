## Api грузовиков

### получение заявки, если оставить параметр null, то он не будет использоваться в фильтре

Post /search/cargo_request?page_number=int&page_size=int
````
{
    "id": "uuid",
    "consignerId": int,
    "recipientId": int,
    "status": "string",
    "createdFrom": "iso8601",
    "createdTo": "iso8601"
}
````

Http status: 200
````
{
  "cargoRequests": [
    {
        "id": "uuid",
        "consignerId": int,
        "recipientId": int,
        "createdAt": long,
        "deadline": long,
        "calculatedTripId": "uuid",
        "actualTripId": "uuid",
        "fromStation": {
            address: string,
            coords: {
                lat: float,
                lon: float
            }
        },
        "toStation": {
            address: string,
            coords: {
                lat: float,
                lon: float
            }
        },
        "maxPrice": "decimal(10, 2)",
        "status": "string"
    }
  ]
}

````

### создание заявки

Post /cargo_request
````
{
    "consignerId": int,
    "recipientId": int,
    "fromStation": {
        address: string,
        coords: {
            lat: float,
            lon: float
        }
    },
    "toStation":  {
        address: string,
        coords: {
            lat: float,
            lon: float
        }
    },
    "deadline": "iso8601",
    "maxPrice": "decimal(10, 2)"
}
````

Http status: 200
````
{
  "id": "uuid"
}
````

### получение типов грузов
Get /cargo/types

Http status: 200
````
{
  "cargoTypes": [
    {
      "id": int,
      "type": "string",
      "fragile": "boolean"
    }
  ]
}
````

### создание груза

Post /cargo
````
{
  cargo: [
    {
      "length": int,
      "height": int,
      "width": int,
      "weight": int,
      "cargoType": int,
      "description": "string",
      "worth": int,
      "cargoRequestId": int
    }
  ]
}
````

Http status: 200
````
{
  ids: [
    int,
    int
  ]
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
  "id: int
}
````

### получение поездок для заявки

Get /trips/cargo_request/{id}?page_number=int&page_size=int

Http status: 200
````
{
  "trips": [
    {
      "id": "uuid",
      "fromStation": {
        address: string,
        coords: {
            lat: float,
            lon: float
        }
      },
      "toStation": {
        address: string,
        coords: {
            lat: float,
            lon: float
        }
      },
      "startedAt": long,
      "calculatedEndAt": long,
      "actualEndAt": long,
      "price": "decimal(10, 2)",
      "status": "string",
      "carrierId": int,
      "carId": int
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
  "fromStation": {
    address: string,
    coords: {
        lat: float,
        lon: float
    }
  },
  "toStation": {
    address: string,
    coords: {
        lat: float,
        lon: float
    }
  },
  "startedAt": long,
  "calculatedEndAt": long,
  "actualEndAt": long,
  "price": "decimal(10, 2)",
  "status": "string",
  "carrierId": int,
  "carId": int
}
````

### получение маршрута для поездки
Get /route/trip/{uuid}

Http status: 200
````
{
  "stations": [
    {
      "station": {
        address: string,
        coords: {
            lat: float,
            lon: float
        }
      },
      "distance": int,
      "orderNum": int,
      "arrivalAt": long,
      "departureTime": long
    }
  ]
}
````

## Внутренние запросы (фронту не нужны)

### сохранение отправителя в системе

Post /consigner
````
{
  "id": int
}

````

Http status: 200

### сохранение водителя в системе

Post /carrier
````
{
  "id": int,
  "driver_category": "string"
}
````

Http status: 200

## статусы
статусы для заявок, грузов, поездок - констаны в string, сообщим позже 
````
WAITING_TRIP_CHOICE, WAITING_DRIVER_APPROVAL, APPROVED_BY_DRIVER
REJECTED_BY_DRIVER, IN_DELIVERY, COMPLETED, CANCELED
````
## форматы
### порядок сортировки
````
ASC
DESC
````

### decimal(10, 2)
````
предполагается, что будет строка типа "7.47", 
надо договориться, как правильно передавать цены
````

### timestamp (long)
время возвращается в секундах UTC https://www.unixtimestamp.com

### iso8601
````
2025-11-03T05:16:38+00:00
````

## ошибки
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
