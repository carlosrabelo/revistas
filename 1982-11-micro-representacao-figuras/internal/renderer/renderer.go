// Package renderer implementa o sistema de projeção cônica descrito no artigo
// "Representação de figuras por computador" de Luiz Antonio Pereira.
//
// Este pacote é o coração matemático da implementação, convertendo pontos
// tridimensionais em coordenadas 2D da tela usando as fórmulas originais
// do artigo publicado na MICRO SISTEMAS em novembro de 1982.
//
// Principais responsabilidades:
// - Aplicar transformações 3D→2D (perspectiva cônica)
// - Renderizar linhas e vértices na tela
// - Gerenciar configurações visuais (cores, espessuras, etc.)
// - Exportar imagens em alta qualidade
//
// Fundamentação matemática:
// O artigo descreve a projeção cônica onde cada ponto P(x,y,z) no espaço
// é projetado em um plano através das fórmulas:
//   x' = (Px - Vx) * R / (Pz - Vz)
//   y' = (Py - Vy) * R / (Pz - Vz)
// onde V é o observador e R é a distância do plano projetante.
package renderer

import (
	"fmt"

	"representacao-figuras/pkg/types"

	"github.com/fogleman/gg"
)

// Renderer3D implementa o sistema de projeção cônica do artigo.
//
// Esta estrutura encapsula o contexto gráfico e os parâmetros necessários
// para realizar a projeção de figuras tridimensionais em uma tela 2D,
// seguindo fielmente as equações descritas no artigo original.
type Renderer3D struct {
	context *gg.Context   // Contexto gráfico para desenho (biblioteca gg)
	width   int           // Largura da tela em pixels
	height  int           // Altura da tela em pixels
	camera  types.Camera  // Parâmetros da câmera virtual
	centerX float64       // Centro X da tela (width/2)
	centerY float64       // Centro Y da tela (height/2)
}

// New cria um novo renderizador 3D com as dimensões especificadas.
//
// Inicializa o contexto gráfico com configurações padrão que remetem
// à apresentação visual do artigo original: fundo branco e linhas pretas,
// similar aos gráficos do HP-85.
//
// Parâmetros:
//   width: largura da tela em pixels
//   height: altura da tela em pixels
//
// Retorna:
//   *Renderer3D: renderizador configurado e pronto para uso
func New(width, height int) *Renderer3D {
	// Cria contexto gráfico com as dimensões especificadas
	ctx := gg.NewContext(width, height)

	// Configuração visual padrão (similar ao artigo original)
	// Fundo branco como no HP-85 e nos exemplos do artigo
	ctx.SetRGB(1, 1, 1) // RGB(255,255,255) = branco
	ctx.Clear()

	// Linhas pretas para contraste máximo (padrão dos anos 80)
	ctx.SetRGB(0, 0, 0) // RGB(0,0,0) = preto
	ctx.SetLineWidth(1.0) // Linha fina padrão

	return &Renderer3D{
		context: ctx,
		width:   width,
		height:  height,
		// Calcula centro da tela para facilitar projeções
		centerX: float64(width) / 2,
		centerY: float64(height) / 2,
	}
}

// SetCamera define os parâmetros da câmera virtual.
//
// Configura a posição do observador (V), a distância do plano projetante (R)
// e as dimensões da "tela virtual" (L1, L2) conforme descrito no artigo.
// Estes parâmetros são fundamentais para os cálculos de perspectiva cônica.
//
// Parâmetros:
//   camera: configuração da câmera com observador, distância e dimensões
func (r *Renderer3D) SetCamera(camera types.Camera) {
	r.camera = camera
}

// ProjectPoint implementa a projeção cônica conforme o artigo original.
//
// Esta é a função central do sistema, que implementa as equações fundamentais
// da perspectiva cônica descritas nas páginas 6-7 do artigo.
//
// PROCESSO MATEMÁTICO (conforme artigo):
//
// 1. TRANSLAÇÃO: Move o ponto para o sistema de coordenadas do observador
//    P' = P - V (onde V é a posição do observador)
//
// 2. PROJEÇÃO CÔNICA: Aplica as fórmulas do artigo (equações na página 7)
//    x' = P'x * R / P'z
//    y' = P'y * R / P'z
//    onde R é a distância do plano projetante
//
// 3. NORMALIZAÇÃO: Converte para coordenadas de tela (pixels)
//    Usa as dimensões L1 e L2 para escalar proporcionalmente
//
// SISTEMA DE COORDENADAS (conforme implementação):
// - X: horizontal (largura)
// - Y: profundidade (distância do observador)
// - Z: vertical (altura)
//
// Parâmetros:
//   p: ponto 3D no espaço mundial
//
// Retorna:
//   types.Point2D: ponto projetado em coordenadas de tela (pixels)
func (r *Renderer3D) ProjectPoint(p types.Point3D) types.Point2D {
	// === ETAPA 1: TRANSLAÇÃO ===
	// Move o ponto para o sistema de coordenadas relativo ao observador
	// Conforme descrito no artigo: P' = P - V

	// Coordenada horizontal (largura)
	px := p.X - r.camera.Observer.X

	// Coordenada vertical (altura) - note o uso de Z
	py := p.Z - r.camera.Observer.Z

	// Coordenada de profundidade (distância) - note o uso de Y
	pz := p.Y - r.camera.Observer.Y

	// === PROTEÇÃO CONTRA DIVISÃO POR ZERO ===
	// Pontos atrás da câmera (pz ≤ 0) ou muito próximos causam problemas
	// na divisão. O artigo não trata deste caso, mas é necessário na prática.
	if pz <= 0.1 {
		pz = 0.1 // Valor mínimo para evitar divisão por zero
	}

	// === ETAPA 2: PROJEÇÃO CÔNICA ===
	// Aplica as fórmulas fundamentais do artigo (equações 2 da página 7)
	// x = Px * R/Pz
	// y = Py * R/Pz
	projX := px * r.camera.Distance / pz
	projY := py * r.camera.Distance / pz

	// === ETAPA 3: CONVERSÃO PARA COORDENADAS DE TELA ===
	// Escala as coordenadas projetadas para o tamanho real da tela
	// Usa as dimensões L1 (largura) e L2 (altura) da "tela virtual"
	scaleX := float64(r.width) / r.camera.Width   // pixels por unidade em X
	scaleY := float64(r.height) / r.camera.Height // pixels por unidade em Y

	// Converte para coordenadas finais de tela
	// Centro da tela + deslocamento escalado
	screenX := r.centerX + (projX * scaleX)

	// Y negativo porque em telas o eixo Y cresce para baixo
	// mas em matemática cresce para cima
	screenY := r.centerY - (projY * scaleY)

	return types.Point2D{X: screenX, Y: screenY}
}

// RenderFigure renderiza uma figura 3D usando projeção cônica com configurações padrão.
//
// Esta função é um wrapper conveniente que usa as configurações visuais padrão.
// Para controle total sobre a aparência, use RenderFigureWithConfig diretamente.
//
// Parâmetros:
//   figure: figura 3D a ser renderizada
//
// Retorna:
//   error: nil se bem-sucedido, erro caso contrário
func (r *Renderer3D) RenderFigure(figure *types.Figure) error {
	return r.RenderFigureWithConfig(figure, DefaultRenderConfig())
}

// RenderFigureWithConfig renderiza uma figura 3D com configurações visuais personalizadas.
//
// Esta é a função principal de renderização, que implementa o processo completo
// descrito no artigo:
// 1. Projeta todos os pontos 3D para coordenadas 2D
// 2. Desenha as linhas conectando os pontos
// 3. Opcionalmente mostra vértices e rótulos
//
// O resultado é uma representação 2D da figura 3D que simula a visão real
// do objeto através da perspectiva cônica.
//
// Parâmetros:
//   figure: figura 3D contendo pontos, linhas e câmera
//   cfg: configurações visuais (cores, espessuras, etc.)
//
// Retorna:
//   error: nil se bem-sucedido, erro caso a figura seja inválida
func (r *Renderer3D) RenderFigureWithConfig(figure *types.Figure, cfg RenderConfig) error {
	// === VALIDAÇÃO DE ENTRADA ===
	if len(figure.Pontos) == 0 {
		return fmt.Errorf("figura não possui pontos")
	}

	// === CONFIGURAÇÃO VISUAL ===
	// Prepara o contexto gráfico com as cores e estilos especificados

	// Define cor de fundo e limpa a tela
	r.context.SetRGB(cfg.Background.R, cfg.Background.G, cfg.Background.B)
	r.context.Clear()

	// Configura cor e espessura das linhas
	r.context.SetRGB(cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B)
	r.context.SetLineWidth(cfg.LineWidth)

	// === PROJEÇÃO 3D → 2D ===
	// Aplica a transformação de perspectiva cônica a todos os pontos
	// Esta é a etapa central que implementa as equações do artigo
	pontos2D := make([]types.Point2D, len(figure.Pontos))
	for i, ponto3D := range figure.Pontos {
		// Cada ponto 3D é projetado individualmente usando ProjectPoint
		pontos2D[i] = r.ProjectPoint(ponto3D)
	}

	// === DESENHO DAS ARESTAS ===
	// Conecta os pontos projetados conforme especificado na figura
	for _, linha := range figure.Linhas {
		// Verificação de segurança: índices válidos
		if linha.P1 >= len(pontos2D) || linha.P2 >= len(pontos2D) {
			continue // Ignora linhas com referências inválidas
		}

		// Obtém os pontos 2D projetados
		p1 := pontos2D[linha.P1]
		p2 := pontos2D[linha.P2]

		// Desenha a linha conectando os dois pontos
		r.context.MoveTo(p1.X, p1.Y)  // Move para o primeiro ponto
		r.context.LineTo(p2.X, p2.Y)  // Desenha linha até o segundo
		r.context.Stroke()            // Aplica o traço
	}

	// === DESENHO DOS VÉRTICES (OPCIONAL) ===
	if cfg.ShowVertices {
		// Muda para cor dos vértices
		r.context.SetRGB(cfg.VertexColor.R, cfg.VertexColor.G, cfg.VertexColor.B)

		for i, p2D := range pontos2D {
			// Desenha um pequeno círculo em cada vértice
			r.context.DrawCircle(p2D.X, p2D.Y, 2)
			r.context.Fill()

			// === DESENHO DOS RÓTULOS (SE ATIVADO) ===
			if cfg.ShowLabels && figure.Pontos[i].Nome != "" {
				// Muda para cor do texto
				r.context.SetRGB(cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B)
				// Desenha o nome do ponto próximo ao vértice
				r.context.DrawString(figure.Pontos[i].Nome, p2D.X+5, p2D.Y-5)
				// Volta para cor dos vértices
				r.context.SetRGB(cfg.VertexColor.R, cfg.VertexColor.G, cfg.VertexColor.B)
			}
		}

		// Restaura cor das linhas para futuras operações
		r.context.SetRGB(cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B)

	} else if cfg.ShowLabels {
		// === RÓTULOS SEM VÉRTICES ===
		// Se apenas os rótulos devem ser mostrados (sem os círculos)
		for i, p2D := range pontos2D {
			if figure.Pontos[i].Nome == "" {
				continue // Pula pontos sem nome
			}
			// Usa cor das linhas para o texto
			r.context.SetRGB(cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B)
			r.context.DrawString(figure.Pontos[i].Nome, p2D.X+5, p2D.Y-5)
		}
		// Garante que a cor das linhas permanece configurada
		r.context.SetRGB(cfg.LineColor.R, cfg.LineColor.G, cfg.LineColor.B)
	}

	return nil
}

// SaveImage salva a imagem renderizada em arquivo PNG.
//
// Exporta o resultado da renderização para um arquivo de imagem,
// permitindo preservar e compartilhar as figuras 3D geradas.
// Uma grande evolução em relação ao HP-85 original!
//
// Parâmetros:
//   filename: caminho do arquivo PNG a ser criado
//
// Retorna:
//   error: nil se bem-sucedido, erro caso haja problemas de E/S
func (r *Renderer3D) SaveImage(filename string) error {
	return r.context.SavePNG(filename)
}

// GetImage retorna a imagem renderizada como interface{}.
//
// Permite acesso direto à imagem em memória para integração
// com outros sistemas ou exibição em interfaces gráficas.
//
// Retorna:
//   interface{}: imagem renderizada (tipo image.Image)
func (r *Renderer3D) GetImage() interface{} {
	return r.context.Image()
}

// AddGrid adiciona uma grade de referência à imagem (função utilitária).
//
// Desenha uma grade de linhas finas para ajudar na visualização e
// depuração das projeções. Útil para fins educacionais e de desenvolvimento.
// Esta funcionalidade vai além do artigo original, sendo uma adição moderna.
//
// A grade é desenhada com espaçamento de 50 pixels em cor cinza claro
// para não interferir na visualização das figuras principais.
func (r *Renderer3D) AddGrid() {
	// Configuração visual da grade
	r.context.SetRGB(0.9, 0.9, 0.9) // Cinza bem claro (quase branco)
	r.context.SetLineWidth(0.5)     // Linha bem fina

	// Desenha linhas verticais (espaçadas a cada 50 pixels)
	for x := 0; x < r.width; x += 50 {
		r.context.MoveTo(float64(x), 0)                // Topo da tela
		r.context.LineTo(float64(x), float64(r.height)) // Base da tela
		r.context.Stroke()
	}

	// Desenha linhas horizontais (espaçadas a cada 50 pixels)
	for y := 0; y < r.height; y += 50 {
		r.context.MoveTo(0, float64(y))                // Esquerda da tela
		r.context.LineTo(float64(r.width), float64(y)) // Direita da tela
		r.context.Stroke()
	}

	// Restaura configurações padrão para não afetar desenhos posteriores
	r.context.SetRGB(0, 0, 0) // Volta para preto
	r.context.SetLineWidth(1.0) // Volta para espessura padrão
}
