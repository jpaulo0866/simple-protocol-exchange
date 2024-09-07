## Simple Protocolo Exchange

Serviço serve para receber dados em um protocolo, mapear e enviar para outro local

- TCP -> TCP
- UDP -> UDP
- TCP -> UDP
- UDP -> TCP
- UDP -> HTTP
- TCP -> HTTP
- HTTP -> HTTP
- etc...

Além disso, possui uma funcionalidade de transformação de JSON com estruturas simples:

- Adicionar valores estáticos
- Remover campos
- Remapear campos

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