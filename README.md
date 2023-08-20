# apicurio-schema-registry-rabbitmq

## build

```shell
go build ./cmd/producer -o bin/producer 
```

```shell
go build ./cmd/consumer -o bin/consumer 
```

## Szenarien

### Szenario 1: optionales Feld ergänzen

Hier wird das optionale Feld `occupation` ergänzt. *Der Validator akzeptiert die Nachricht*

```shell
 go run cmd/producer/main.go --payload='{"name":"Max Weis","age":24, "occupation":"software engineer"}' --validation=true
```

## Szenario 2: Pflichtfeld geändert

Hier wird das Pflichtfeld `name` durch ein String Array ersäzt. *Der Validator akzeptiert die Nachricht nicht*

```shell
 go run cmd/producer/main.go --payload='{"name":["Max", "Weis"],"age":24}' --validation=true
```

## Szenario 3: Semantik eines Pflichtfeldes ändern
