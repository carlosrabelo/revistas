// Package types define as estruturas de dados fundamentais para representação
// de figuras tridimensionais conforme o artigo "Representação de figuras por
// computador" de Luiz Antonio Pereira (MICRO SISTEMAS, Nov/1982).
//
// Este pacote implementa os conceitos matemáticos descritos no artigo original,
// adaptando-os para uma linguagem moderna mantendo fidelidade aos fundamentos.
package types

// Point3D representa um ponto no espaço tridimensional.
//
// Sistema de coordenadas conforme o artigo:
// - X: coordenada horizontal (largura)
// - Y: coordenada de profundidade (distância do observador)
// - Z: coordenada vertical (altura)
//
// Este sistema segue a convenção do artigo original onde Y representa
// a profundidade, sendo fundamental para os cálculos de perspectiva cônica.
type Point3D struct {
	X, Y, Z float64 // Coordenadas espaciais em unidades arbitrárias
	Nome    string  `yaml:"nome,omitempty"` // Nome opcional para identificação
}

// Point2D representa um ponto projetado na tela (resultado da projeção 3D→2D).
//
// Após aplicar as fórmulas de perspectiva cônica, os pontos 3D são convertidos
// em coordenadas de tela para renderização. As coordenadas são em pixels.
type Point2D struct {
	X, Y float64 // Coordenadas na tela em pixels
}

// Line representa uma linha conectando dois pontos da figura.
//
// Conforme descrito no artigo, as figuras são definidas por vértices
// conectados por segmentos de reta. Esta estrutura armazena os índices
// dos pontos que devem ser conectados.
type Line struct {
	P1, P2 int // Índices dos pontos na lista (base 0)
}

// RenderSettings controla opções visuais de renderização da figura.
//
// Estas configurações permitem personalizar a aparência da imagem gerada,
// indo além das capacidades do HP-85 original mas mantendo a essência
// da representação gráfica descrita no artigo.
type RenderSettings struct {
	// Dimensões da tela de saída (em pixels)
	CanvasWidth  int `yaml:"largura_canvas,omitempty"`  // Largura da imagem
	CanvasHeight int `yaml:"altura_canvas,omitempty"`   // Altura da imagem

	// Configurações de cores (nomes ou códigos hex)
	Background  string `yaml:"fundo,omitempty"`       // Cor de fundo
	LineColor   string `yaml:"cor_linha,omitempty"`   // Cor das linhas
	VertexColor string `yaml:"cor_vertices,omitempty"` // Cor dos vértices

	// Configurações de desenho
	LineWidth float64 `yaml:"espessura_linha,omitempty"` // Espessura das linhas

	// Opções de visualização (ponteiros permitem nil = usar padrão)
	ShowVertices *bool `yaml:"mostrar_vertices,omitempty"` // Mostrar pontos dos vértices
	ShowLabels   *bool `yaml:"mostrar_nomes,omitempty"`    // Mostrar nomes dos pontos
}

// Camera representa os parâmetros da câmera virtual conforme o artigo.
//
// Implementa o sistema de projeção cônica descrito nas páginas 6-7 do artigo,
// onde o observador (V) está posicionado no espaço e observa objetos através
// de um plano projetante a uma distância R.
//
// Parâmetros fundamentais da perspectiva cônica:
// - V (Observer): posição do observador no espaço 3D
// - R (Distance): distância do observador ao plano projetante
// - L1,L2 (Width,Height): dimensões do "retângulo de visualização"
type Camera struct {
	// Posição do observador no espaço 3D (ponto V do artigo)
	Observer Point3D `yaml:"observador"`

	// Distância R do plano projetante (fundamental para perspectiva)
	// Valores maiores = menos perspectiva, valores menores = mais perspectiva
	Distance float64 `yaml:"distancia"`

	// Dimensões da "tela virtual" (L1 e L2 do artigo)
	// Baseadas nas dimensões do HP-85: proporção 4:3
	Width  float64 `yaml:"largura"` // L1: largura da tela virtual
	Height float64 `yaml:"altura"`  // L2: altura da tela virtual
}

// Figure representa uma figura tridimensional completa.
//
// Esta estrutura encapsula todos os elementos necessários para definir
// e renderizar uma figura 3D conforme a metodologia do artigo:
// 1. Pontos no espaço (vértices)
// 2. Linhas conectando os pontos (arestas)
// 3. Parâmetros da câmera (observador e projeção)
// 4. Configurações de renderização (opcionais)
type Figure struct {
	Nome   string          `yaml:"nome"`    // Nome identificador da figura
	Pontos []Point3D       `yaml:"pontos"`  // Lista de vértices 3D
	Linhas []Line          `yaml:"linhas"`  // Lista de arestas (segmentos)
	Camera Camera          `yaml:"camera"`  // Parâmetros de visualização
	Render *RenderSettings `yaml:"render,omitempty"` // Configurações visuais opcionais
}

// DefaultCamera retorna uma câmera com configuração padrão baseada no artigo.
//
// Os valores padrão são derivados das especificações do HP-85 mencionadas
// no artigo original:
// - Resolução: 256×192 pixels
// - Proporção: 4:3 (1.33:1)
// - Valores L1=12.8 e L2=9.6 mantêm esta proporção em unidades virtuais
//
// A distância R=10 oferece uma perspectiva moderada, adequada para
// visualização geral de objetos tridimensionais.
func DefaultCamera() Camera {
	return Camera{
		// Observador posicionado na origem do sistema de coordenadas
		Observer: Point3D{X: 0, Y: 0, Z: 0},

		// Distância moderada para perspectiva equilibrada
		Distance: 10,

		// Dimensões baseadas no HP-85 (proporção 4:3)
		// L1 = 12.8 unidades (largura)
		Width: 12.8,
		// L2 = 9.6 unidades (altura, mantém proporção 4:3)
		Height: 9.6,
	}
}
