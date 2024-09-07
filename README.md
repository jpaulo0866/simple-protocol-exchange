## Simple Protocolo Exchange



## SETUP

Instalar Go (se necessário):
	•	Baixe e instale o Go a partir de https://golang.org/dl/.

## RUN

```
go run .
```

## BUILD

```
go build -o dist/spe
```

## TEST IT

```
echo -n -e '{"_app_name":"myapp","_source":"localhost","_message":"application started","_timestamp":"2024-09-07T10:30:00Z","_level_name":"INFO","_environment":"production","_logger_name":"defaultLogger","_traceId":"1234567890"}'"\0" | nc -w0 localhost 12203
```