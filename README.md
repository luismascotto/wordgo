# WordGo

Uma aplicação Go para encontrar palavras grandes em uma matriz de letras usando processamento paralelo com goroutines.

## Descrição

WordGo é uma ferramenta de linha de comando que carrega uma matriz de letras e um dicionário de palavras, permitindo buscar palavras de forma eficiente usando múltiplas goroutines para processamento paralelo.

## Funcionalidades

- ✅ Carregamento de matriz de letras a partir de arquivo de texto
- ✅ Carregamento de dicionário de palavras
- ✅ Validação de dados de entrada
- ✅ Estrutura preparada para busca paralela com goroutines
- 🔄 Algoritmos de busca (em desenvolvimento)
- 🔄 Processamento paralelo (em desenvolvimento)

## Estrutura do Projeto

```
wordgo/
├── main.go          # Arquivo principal da aplicação
├── go.mod           # Módulo Go
├── README.md        # Este arquivo
└── res/             # Recursos
    ├── example.txt  # Matriz de letras de exemplo
    └── words.txt    # Dicionário de palavras
```

## Formato dos Arquivos

### Matriz de Letras (example.txt)
- Cada linha representa uma linha da matriz
- Todas as linhas devem ter o mesmo comprimento
- Caracteres vazios são ignorados

### Dicionário (words.txt)
- Uma palavra por linha
- Palavras são convertidas para maiúsculas automaticamente
- Linhas vazias são ignoradas

## Como Executar

1. Certifique-se de ter Go 1.21+ instalado
2. Clone o repositório
3. Execute o projeto:

```bash
go run main.go
```

## Próximos Passos

- [ ] Implementar algoritmos de busca em todas as direções
- [ ] Adicionar processamento paralelo com goroutines
- [ ] Implementar busca de palavras de tamanho mínimo
- [ ] Adicionar estatísticas de performance
- [ ] Criar interface de linha de comando configurável

## Tecnologias

- **Go 1.21+** - Linguagem principal
- **Goroutines** - Para processamento paralelo
- **Sync package** - Para sincronização thread-safe

## Licença

Este projeto está sob licença MIT.
