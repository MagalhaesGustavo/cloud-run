# cloud run weather - Golang Expert

Para buildar a aplicação use o comando docker abaixo:

```
docker-compose up -d
```

A Aplicação estará disponível em:

```
http://localhost:8080
```

A API foi publicada usando o Google Cloud Run, para acessar basta usar o endereço substituindo o {cep} pelo número do cep desejado, sem espaços, hifem ou pontuação:

```
https://cloudrun-goexpert-fvm5rtvuiq-uc.a.run.app/{cep}
```

Exemplo:

```
https://cloudrun-goexpert-fvm5rtvuiq-uc.a.run.app/13380001
```