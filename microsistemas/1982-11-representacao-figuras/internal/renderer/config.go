// Package renderer/config implementa o sistema de configuração visual
// para renderização de figuras 3D.
//
// Este módulo estende as capacidades do artigo original, permitindo
// personalização de cores, espessuras e outros aspectos visuais
// que não eram possíveis no HP-85 de 1982.
//
// Responsabilidades:
// - Conversão de configurações YAML para estruturas internas
// - Parsing de cores em diferentes formatos (nomes, hex, etc.)
// - Fornecimento de valores padrão sensíveis
// - Validação de parâmetros visuais
package renderer

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"representacao-figuras/pkg/types"
)

// colorRGB representa uma cor no espaço RGB com valores de 0.0 a 1.0.
//
// Esta representação é compatível com a biblioteca gráfica gg
// e permite precisão superior aos 16 ou 256 cores do hardware de 1982.
type colorRGB struct {
	R float64 // Componente vermelho (0.0 = sem vermelho, 1.0 = vermelho total)
	G float64 // Componente verde (0.0 = sem verde, 1.0 = verde total)
	B float64 // Componente azul (0.0 = sem azul, 1.0 = azul total)
}

// RenderConfig encapsula todas as opções visuais aplicadas pelo renderizador.
//
// Esta estrutura permite controle fino sobre a aparência das figuras,
// indo muito além das capacidades limitadas do HP-85 original que
// tinha apenas algumas cores básicas e resolução fixa.
type RenderConfig struct {
	Background   colorRGB // Cor de fundo da imagem
	LineColor    colorRGB // Cor das linhas (arestas) da figura
	LineWidth    float64  // Espessura das linhas em pixels
	VertexColor  colorRGB // Cor dos vértices (pontos)
	ShowVertices bool     // Se deve mostrar círculos nos vértices
	ShowLabels   bool     // Se deve mostrar nomes dos pontos
}

// DefaultRenderConfig retorna a configuração visual padrão.
//
// Os valores padrão são inspirados na estética do artigo original:
// - Fundo branco (como papel)
// - Linhas pretas (como tinta)
// - Estilo minimalista sem elementos visuais extras
//
// Esta abordagem mantém a clareza e legibilidade dos exemplos
// apresentados na MICRO SISTEMAS de 1982.
func DefaultRenderConfig() RenderConfig {
	return RenderConfig{
		// Fundo branco (RGB: 255,255,255) - estética clássica
		Background: colorRGB{R: 1, G: 1, B: 1},

		// Linhas pretas (RGB: 0,0,0) - máximo contraste
		LineColor: colorRGB{R: 0, G: 0, B: 0},

		// Linha fina padrão (1 pixel)
		LineWidth: 1.0,

		// Vértices em vermelho escuro para destaque quando ativados
		VertexColor: colorRGB{R: 0.8, G: 0, B: 0},

		// Por padrão, apenas as linhas são visíveis (como no artigo)
		ShowVertices: false,
		ShowLabels:   false,
	}
}

// ConfigFromFigure converte configurações YAML para estrutura interna.
//
// Esta função faz a ponte entre as configurações declarativas
// (definidas no arquivo YAML) e as estruturas otimizadas para
// renderização interna.
//
// Processo:
// 1. Começa com configurações padrão
// 2. Aplica sobreposições definidas no YAML
// 3. Valida valores fornecidos
// 4. Retorna configuração final ou erro
//
// Parâmetros:
//   fig: figura contendo configurações opcionais de renderização
//
// Retorna:
//   RenderConfig: configuração final validada
//   error: erro se houver valores inválidos
func ConfigFromFigure(fig *types.Figure) (RenderConfig, error) {
	// Começa com valores padrão seguros
	cfg := DefaultRenderConfig()

	// Se não há configurações customizadas, usa padrões
	if fig == nil || fig.Render == nil {
		return cfg, nil
	}

	settings := fig.Render

	// === PROCESSAMENTO DE CORES ===

	// Cor de fundo (background)
	if settings.Background != "" {
		col, err := parseColor(settings.Background)
		if err != nil {
			return cfg, fmt.Errorf("cor de fundo inválida: %w", err)
		}
		cfg.Background = col
	}

	// Cor das linhas
	if settings.LineColor != "" {
		col, err := parseColor(settings.LineColor)
		if err != nil {
			return cfg, fmt.Errorf("cor da linha inválida: %w", err)
		}
		cfg.LineColor = col
	}

	// Cor dos vértices
	if settings.VertexColor != "" {
		col, err := parseColor(settings.VertexColor)
		if err != nil {
			return cfg, fmt.Errorf("cor dos vértices inválida: %w", err)
		}
		cfg.VertexColor = col
	}

	// === CONFIGURAÇÕES NUMÉRICAS ===

	// Espessura das linhas (deve ser positiva)
	if settings.LineWidth > 0 {
		cfg.LineWidth = settings.LineWidth
	}

	// === CONFIGURAÇÕES BOOLEANAS ===
	// Usa ponteiros para distinguir entre "não especificado" e "false"

	if settings.ShowVertices != nil {
		cfg.ShowVertices = *settings.ShowVertices
	}

	if settings.ShowLabels != nil {
		cfg.ShowLabels = *settings.ShowLabels
	}

	return cfg, nil
}

// namedColors contém cores pré-definidas por nome para conveniência.
//
// Permite uso de nomes intuitivos em vez de códigos hexadecimais,
// tornando os arquivos YAML mais legíveis. Inclui variações de
// grafia (gray/grey) para flexibilidade.
var namedColors = map[string]colorRGB{
	// Cores básicas
	"white": {R: 1, G: 1, B: 1},       // Branco puro
	"black": {R: 0, G: 0, B: 0},       // Preto puro

	// Tons de cinza (ambas grafias aceitas)
	"gray":      {R: 0.5, G: 0.5, B: 0.5},     // Cinza médio
	"grey":      {R: 0.5, G: 0.5, B: 0.5},     // Cinza médio (grafia britânica)
	"lightgray": {R: 0.82, G: 0.82, B: 0.82},  // Cinza claro
	"lightgrey": {R: 0.82, G: 0.82, B: 0.82},  // Cinza claro (grafia britânica)
	"darkgray":  {R: 0.25, G: 0.25, B: 0.25},  // Cinza escuro
	"darkgrey":  {R: 0.25, G: 0.25, B: 0.25},  // Cinza escuro (grafia britânica)
}

// parseColor converte uma string de cor para colorRGB.
//
// Suporta múltiplos formatos de entrada:
// 1. Nomes de cores ("white", "black", "red", etc.)
// 2. Códigos hexadecimais completos ("#ff0000", "ff0000")
// 3. Códigos hexadecimais curtos ("#f00" → "#ff0000")
//
// Todos os formatos são case-insensitive para conveniência.
//
// Parâmetros:
//   value: string representando uma cor
//
// Retorna:
//   colorRGB: cor convertida para formato interno
//   error: erro se o formato for inválido
func parseColor(value string) (colorRGB, error) {
	// Normalização: remove espaços e converte para minúsculas
	v := strings.TrimSpace(strings.ToLower(value))
	if v == "" {
		return colorRGB{}, fmt.Errorf("valor vazio")
	}

	// === TENTATIVA 1: CORES NOMEADAS ===
	// Verifica se é uma cor pré-definida
	if col, ok := namedColors[v]; ok {
		return col, nil
	}

	// === TENTATIVA 2: CÓDIGO HEXADECIMAL ===

	// Remove prefixo '#' se presente
	if strings.HasPrefix(v, "#") {
		v = v[1:]
	}

	// Expande formato curto (#rgb → #rrggbb)
	if len(v) == 3 {
		// Cada caractere é duplicado: "f0a" → "ff00aa"
		var sb strings.Builder
		for _, ch := range v {
			sb.WriteRune(ch) // Primeiro
			sb.WriteRune(ch) // Duplica
		}
		v = sb.String()
	}

	// Valida comprimento final
	if len(v) != 6 {
		return colorRGB{}, fmt.Errorf("formato de cor inválido: %s", value)
	}

	// === CONVERSÃO DOS COMPONENTES RGB ===

	// Converte cada par de caracteres hex para float64 (0.0-1.0)
	r, err := parseHexComponent(v[0:2]) // Vermelho
	if err != nil {
		return colorRGB{}, err
	}

	g, err := parseHexComponent(v[2:4]) // Verde
	if err != nil {
		return colorRGB{}, err
	}

	b, err := parseHexComponent(v[4:6]) // Azul
	if err != nil {
		return colorRGB{}, err
	}

	return colorRGB{R: r, G: g, B: b}, nil
}

// parseHexComponent converte um componente hexadecimal (00-FF) para float64 (0.0-1.0).
//
// Transforma valores de cor do formato hexadecimal (0-255) para o formato
// de ponto flutuante usado pela biblioteca gráfica (0.0-1.0).
//
// Exemplos:
//   "00" → 0.0 (sem cor)
//   "80" → ~0.5 (meio tom)
//   "FF" → 1.0 (cor total)
//
// Parâmetros:
//   component: string de 2 caracteres hexadecimais
//
// Retorna:
//   float64: valor normalizado entre 0.0 e 1.0
//   error: erro se não for hexadecimal válido
func parseHexComponent(component string) (float64, error) {
	// Converte hex (base 16) para inteiro de 8 bits (0-255)
	v, err := strconv.ParseUint(component, 16, 8)
	if err != nil {
		return 0, err
	}

	// Normaliza para 0.0-1.0 e arredonda para 3 casas decimais
	// para evitar imprecisões de ponto flutuante
	return math.Round((float64(v)/255.0)*1000) / 1000, nil
}
