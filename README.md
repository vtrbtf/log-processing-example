# Log processor demo

> Demo simples de um parser de multiplos logs (estilo apache) que separa e unifica os logs por `userid` ordenando por data.

## Table of Contents

- [Background](#background)
- [Environments](#environments)
- [Resources](#resources)
- [Docs](#related-documentation)
- [Maintainers](#maintainers)

## Background

O propósito da demo é construir uma solução simples que separe logs por uma determinada chave, no caso um cookie `userid`. Esse log inicial pode ser gerado usando o apache-fake-log-gen, que é um fork de um projeto OS. O fork foi necessário para poder modificar o formato do log para incluir userids gerados aleatoriamente. Já o collector faz o parse dos logs gerados escrevendo o log unificado por `userid` no diretorio `/tmp` ( `tmp/${userid}.log` )

### Disclaimer
O projeto foi feito somente para fins de processo seletivo. __Para problemas reais com logs usaria uma ferramenta especifica para logs ( e.g. ELK ), ou para um problema de grande volume de dados, usaria Spark, Pig ou semelhantes.__

### Input 

Projeto python3 forkado de [kiritbasu/Fake-Apache-Log-Generator](https://github.com/kiritbasu/Fake-Apache-Log-Generator). As modificações no fork foram para incluir o `userid` no formato de log e para incluir a opção do numero de usuários. Todas as dependencias do script estão no `requirements.txt` e podem ser instaladas com: 

> Lembrando que o script foi feito para python 3.6.3

```
   $ pip install -r requirements.txt
```

Para executar o script: 
```
   $ python apache-fake-log-gen.py -n ${numero_de_linhas_no_log} -o LOG --prefix "${prefixo_do_arquivo_de_log}" -u ${numero_de_usuarios}
```

### Collect
Base da demo que foi desenvolvida em Go. Todas as dependencias do projeto (lib para progressbar) já estao versionadas dentro do repositorio. 

#### Arquitetura
A demo foi escrita para poder operar em multiplos arquivos independente de tamanho. Para conseguir isso usei alguns dicts, fazendo IO de forma gradual. Certas operações de ( como IO ) não foram otimizadas ao maximo para a plataforma ( Poderia ter usado goroutines mais efetivamente, mas ainda nao tenho maturidade suficiente na plataforma pra fazer de forma segura). Também não tive tempo o suficiente para availiar ao melhor algortimo para o probema, sendo que o arquivo de log e o número de userids tem um tamanho indefinido. O footprint de memória também poderia ser melhorando, porém não tive tempo o suficiente para fazer o profiling completo da aplicação. Meu tempo de experiencia na plataforma é relativamente pequeno (+- 1 mes), então a organização não ficou algo que eu considerei bom. Também não consegui ter uma cobertuda de testes boa ( por conta de tempo ).

O projeto foi feito tendo em mente que seria um CLI, então tambem esta incluido algumas saidas no console indicando o progresso do collector.

##### Parser
O parser é simples (via regex) e é feito de forma concorrente (entre todos os arquivos de input). Ele extrai a timestamp e userid de cada linha do arquivo e indexa essas informações for userid

##### Reader
Cria um indice de conteudo para todas as linhas de determinado `userid`, sendo a chave do indice um hash de `${input_file};${file_row_number}`. Um dos problemas atuais da solução é que o algoritimo precisa de multiplos lookups para conseguir construir o indice de forma efetiva. Esses lookups podem ser calculados com a formula `numero_de_usuarios * numero_de_arquivos`

##### Writer 
Escreve no arquivo alvo (`userid.log`) de forma ordenada por timestamp, considerando todos os logs parseados

----

Para fazer o build do projeto é necessario somente o go 1.9 e setar o GOPATH para o diretório.

```
   $ cd collect
   $ export GOPATH=`pwd`
   $ go install collect
```

Para rodar os testes

```
   $ cd collect/src/collect
   $ go test -v
```

Para iniciar o CLI ( depois de ter feito o build ): 
```
   $ cd bin
   $ ./collect --logfile arquivo_de_log1.log --logfile arquivo_de_log2.log
```

## Environments

#### Scripts
Os scripts automatizam toda a configuração necessaria para a execução dos projetos. 

##### Usage
Como esperado, em uma box limpa executar o [remote install](#remote-install).  
Depois da instalação de todas as dependencias, executar o script de start.  

```
    $ chmod +x start.sh
    $ ./start.sh
```

O start.sh já esta pré configurado com uns valores default. Para alterar esses valores é só exportar as seguintes variaveis:
```
    $ export LOGPROCESSOR_NUM_OF_SERVERS= ... # numero de servers simulados (default: 4)
    $ export LOGPROCESSOR_NUM_OF_USERS= ... # numero de usuarios (default: 1000)
    $ export LOGPROCESSOR_NUM_OF_LINES_FOR_EACH_FILE= ... # numero de linhas por arquivo de log (default: 1000)
```

#### Docker
Para facilitar o isolamento da execução do projeto eu crei 2 Dockerfiles, cada um responsavel pela a execução de 1 script. As 2 imagens usadas são da [standard library](https://github.com/docker-library/official-images), portando são seguras e pequenas.

O repositório contem scripts para automatizar a instalação e execução do docker e outras dependencias (`install.sh` e `start.sh`).
Um possivél downside é relacionado a como a arquitetura se comporata rodando em Docker, sendo que o scripts são IO bound (Talvez seja necessario setar `ulimit`)

Os arquivos de input serão gerados no docker host em `$HOME/data` e os arquivos de output (unificados por `userid`) serão gerados em `/tmp/` no docker host

##### Remote install
Instalar com somente um curl ( somente para testes )
```
    $ bash < <(curl -s -S -L https://raw.githubusercontent.com/vtrbtf/log-processing-example/master/install.sh)
```

###### Rebuilding 
Para rebuildar as imagens:
```
    # input
    $ cd input && sudo docker build -t accesslog-generator . && cd ..

    # collect
    $ cd collect && sudo docker build -t accesslog-collector . && cd ..
```


#### Vagrant
Para conseguir testar a execução correta dos scripts auxiliares (`install.sh` e `start.sh`) eu criei um `Vagrantfile` que tem uma box configurada com o Ubuntu Trusty Tahr (14.04 64bits)

## Resources
- [Docker](https://www.docker.com/)
- [Go](https://golang.org/)
- [Python](https://www.python.org/)
- [Vagrant](https://www.vagrantup.com/)


## Related documentation
- [How to install Vagrant](https://www.vagrantup.com/docs/installation/)
- [How to install Docker](https://docs.docker.com/engine/installation/)
- [How to install Go](https://golang.org/doc/install)

## Maintainers

[@vtrbtf](https://github.com/vtrbtf)  

