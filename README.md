# Loadtest CLI

Ferramenta de linha de comando para testes de carga em serviços web.

## Características

- Definição de URL, total de requisições e concorrência via flags.
- Pool de workers para execução paralela.
- Relatório final com:
  - Tempo total de execução.
  - Total de requisições feitas.
  - Requisições bem-sucedidas (HTTP 200).
  - Distribuição de códigos de status HTTP (404, 500, etc).
- Logs detalhados com logrus.

## Pré-requisitos

- Go 1.22 ou superior.
- Docker (opcional, para empacotamento).

## Instalação

```bash
# Clone o repositório
git clone <https://github.com/hydde7/goexpert-final-challenge-2>
cd goexpert-final-challenge-2
go mod download
```

## Compilação

```bash
go build -o loadtest main.go
```

## Uso

```bash
./loadtest --url=<URL> --requests=<TOTAL> --concurrency=<CONCORRÊNCIA>
```

### Flags

| Flag          | Descrição                             |
|---------------|---------------------------------------|
| --url         | URL do serviço a ser testado.         |
| --requests    | Número total de requisições.          |
| --concurrency | Número de chamadas simultâneas.       |

### Exemplo

```bash
./loadtest --url=http://example.com --requests=1000 --concurrency=10
```

## Docker

```bash
docker build -t loadtest .
docker run --rm loadtest --url=http://example.com --requests=1000 --concurrency=10
```

## Saída

O log exibirá algo como:

```
time="2025-05-04T13:00:00Z" level=info msg="Iniciando teste de carga: url=http://example.com | total_requests=1000 | concurrency=10"
time="2025-05-04T13:00:10Z" level=info msg="======== Relatório de Carga ========"
time="2025-05-04T13:00:10Z" level=info msg="URL testada: http://example.com"
time="2025-05-04T13:00:10Z" level=info msg="Tempo total de execução: 10.123s"
time="2025-05-04T13:00:10Z" level=info msg="Total de requests feitos: 1000"
time="2025-05-04T13:00:10Z" level=info msg="Requests bem-sucedidos: 980 (200)"
time="2025-05-04T13:00:10Z" level=info msg="Distribuição de códigos HTTP:"
time="2025-05-04T13:00:10Z" level=info msg="  OK (200): 980"
time="2025-05-04T13:00:10Z" level=info msg="  Not Found (404): 15"
time="2025-05-04T13:00:10Z" level=info msg="  Internal Server Error (500): 5"
```
