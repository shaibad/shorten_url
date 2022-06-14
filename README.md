## Shorten URL system

### Descreption
Basic prototype of shorten URL system (not persistent)

### Build
```bash
docker build .
```

### Run
```bash
docker run -it -p 10000:10000 <container id>
```

### Endpoints

/get_url (GET)

/shorten_url (POST)

### Examples
Shorten URL:
```curl -X POST -d '{"Url":"https://recolabs.dev"}' localhost:10000/shorten_url```

Output:
```json
{"Message":"7krJCv7nn5V","Status":"OK"}
```

Get URL:
```curl -X GET localhost:10000/get_url?url=7krJCv7nn5V```

Output:
```json
{"Url":"https://recolabs.dev"}
```

Response for URL that doesn't exist:
```json
{"Message":"Url not found","Status":"Error"}
```

