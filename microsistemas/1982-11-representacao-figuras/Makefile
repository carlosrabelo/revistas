# Makefile para Representação de Figuras 3D
# MICRO SISTEMAS - Novembro/1982

.PHONY: build run clean test ascii viewer help

# Variáveis
BINARY_NAME=figuras3d
CMD_PATH=./cmd/figuras3d
BUILD_DIR=./bin

# Targets principais
help:
	@echo "Representação de Figuras por Computador - MICRO SISTEMAS Nov/1982"
	@echo ""
	@echo "Comandos disponíveis:"
	@echo "  build         - Compila o binário"
	@echo "  generate FILE - Gera PNG do arquivo YAML"
	@echo "  view FILE     - Abre viewfinder interativo"
	@echo "  clean         - Remove binários"
	@echo "  test          - Executa testes"
	@echo ""
	@echo "Exemplos:"
	@echo "  make generate FILE=modelos/cubo.yaml"
	@echo "  make view FILE=modelos/casa.yaml"

build:
	@echo "Compilando representacao-figuras..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)
	@echo "Binário criado: $(BUILD_DIR)/$(BINARY_NAME)"

generate:
	@if [ -z "$(FILE)" ]; then \
		echo "Erro: especifique FILE=arquivo.yaml"; \
		echo "   Exemplo: make generate FILE=modelos/cubo.yaml"; \
		exit 1; \
	fi
	@mkdir -p output
	@go run $(CMD_PATH) generate $(FILE)

view:
	@if [ -z "$(FILE)" ]; then \
		echo "Erro: especifique FILE=arquivo.yaml"; \
		echo "   Exemplo: make view FILE=modelos/casa.yaml"; \
		exit 1; \
	fi
	@go run $(CMD_PATH) view $(FILE)

test:
	@echo "Executando testes..."
	@go test ./...

clean:
	@echo "Limpando binários e arquivos gerados..."
	@rm -rf $(BUILD_DIR)
	@rm -rf output/
	@rm -f output_*.png
	@echo "Limpeza concluída"


# Instalação de dependências
deps:
	@echo "Baixando dependências..."
	@go mod tidy
	@go mod download

# Verificação de qualidade
lint:
	@echo "Verificando código..."
	@gofmt -l .
	@go vet ./...

# Informações do projeto
info:
	@echo "Representação de Figuras por Computador"
	@echo "Baseado no artigo de Luiz Antonio Pereira"
	@echo "MICRO SISTEMAS - Novembro/1982"
	@echo ""
	@echo "Estrutura do projeto:"
	@tree -I 'bin|*.png|go.sum' || ls -la