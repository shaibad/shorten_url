## Shorten URL system

### Descreption
A persistent URL shortener system


### Run
From the main directory, run:
```bash
docker-compose up --build
```

### Endpoints

/shorten_url (POST) - Generates a short url from a given URL and saves it

/get_url (GET) - Redirects to the original URL

### Examples
Shorten URL:
```curl -X POST -d '{"Url":"https://www.google.com"}' localhost:5000/shorten_url```

Output:
```json
{
    "Message": "http://localhost:8080/7PQTO2E",
    "Status": "OK"
}
```

Get URL:
Navigate to: http://localhost:8080/7PQTO2E

Response for a URL that doesn't exist:
```json
{"Message":"Url not found","Status":"Error"}
```

### Testing
From tests directory, run:
```bash
go test -v
```
