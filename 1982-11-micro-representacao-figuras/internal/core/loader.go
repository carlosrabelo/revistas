// Package core implementa as funcionalidades centrais de carregamento
// e validação de figuras tridimensionais.
//
// Este pacote é responsável por:
// - Carregar definições de figuras a partir de arquivos YAML
// - Validar a consistência dos dados carregados
// - Aplicar configurações padrão quando necessário
//
// O formato YAML foi escolhido para substituir o código hardcoded
// do BASIC original, oferecendo uma forma declarativa e legível
// de definir figuras tridimensionais.
package core

import (
	"fmt"
	"os"

	"representacao-figuras/pkg/types"

	"gopkg.in/yaml.v3"
)

// LoadFigureFromYAML carrega e valida uma figura tridimensional a partir de um arquivo YAML.
//
// Esta função substitui a necessidade de definir figuras diretamente no código
// (como era feito no BASIC original), permitindo que usuários criem suas próprias
// figuras em um formato declarativo e intuitivo.
//
// Processo de carregamento:
// 1. Lê o arquivo YAML do sistema de arquivos
// 2. Faz o parse dos dados para a estrutura Figure
// 3. Aplica configurações padrão se necessário (ex: câmera)
// 4. Valida a consistência dos dados
// 5. Retorna a figura pronta para renderização
//
// Parâmetros:
//   filename: caminho para o arquivo YAML contendo a definição da figura
//
// Retorna:
//   *types.Figure: figura carregada e validada
//   error: erro caso haja problemas na leitura, parse ou validação
func LoadFigureFromYAML(filename string) (*types.Figure, error) {
	// Etapa 1: Leitura do arquivo
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo: %w", err)
	}

	// Etapa 2: Parse do YAML para estrutura Go
	var figure types.Figure
	err = yaml.Unmarshal(data, &figure)
	if err != nil {
		return nil, fmt.Errorf("erro ao parsear YAML: %w", err)
	}

	// Etapa 3: Aplicação de padrões
	// Se a câmera não foi especificada (Distance == 0), usa configuração padrão
	// baseada no HP-85 original
	if figure.Camera.Distance == 0 {
		figure.Camera = types.DefaultCamera()
	}

	// Etapa 4: Validação da consistência
	err = validateFigure(&figure)
	if err != nil {
		return nil, fmt.Errorf("figura inválida: %w", err)
	}

	return &figure, nil
}

// validateFigure verifica se a figura está bem formada e consistente.
//
// Realiza verificações essenciais para garantir que a figura possa ser
// renderizada corretamente, evitando erros durante a projeção 3D→2D.
//
// Validações realizadas:
// 1. Presença de pelo menos um ponto (vértice)
// 2. Presença de pelo menos uma linha (aresta)
// 3. Consistência das referências de índices nas linhas
//
// Parâmetros:
//   figure: ponteiro para a figura a ser validada
//
// Retorna:
//   error: nil se válida, ou descrição do problema encontrado
func validateFigure(figure *types.Figure) error {
	// Verificação 1: Deve ter pelo menos um vértice
	// Uma figura sem pontos não pode ser representada
	if len(figure.Pontos) == 0 {
		return fmt.Errorf("figura deve ter pelo menos um ponto")
	}

	// Verificação 2: Deve ter pelo menos uma aresta
	// Linhas conectam os pontos para formar a figura visível
	if len(figure.Linhas) == 0 {
		return fmt.Errorf("figura deve ter pelo menos uma linha")
	}

	// Verificação 3: Consistência das referências de índices
	// Cada linha deve referenciar índices válidos na lista de pontos
	// Índices devem estar no intervalo [0, len(pontos)-1]
	for i, linha := range figure.Linhas {
		// Verifica o primeiro ponto da linha
		if linha.P1 < 0 || linha.P1 >= len(figure.Pontos) {
			return fmt.Errorf("linha %d referencia ponto P1 inválido: %d (deve estar entre 0 e %d)",
				i, linha.P1, len(figure.Pontos)-1)
		}

		// Verifica o segundo ponto da linha
		if linha.P2 < 0 || linha.P2 >= len(figure.Pontos) {
			return fmt.Errorf("linha %d referencia ponto P2 inválido: %d (deve estar entre 0 e %d)",
				i, linha.P2, len(figure.Pontos)-1)
		}
	}

	// Se chegou até aqui, a figura é válida
	return nil
}