// Package main implementa a aplicação de linha de comando para renderizar
// figuras tridimensionais conforme o artigo "Representação de figuras por
// computador" de Luiz Antonio Pereira (MICRO SISTEMAS, Nov/1982).
//
// Esta aplicação moderniza os conceitos do artigo original, oferecendo:
// - Carregamento de figuras via arquivos YAML (vs. código hardcoded)
// - Geração de imagens PNG de alta qualidade (vs. tela 256×192)
// - Interface interativa para visualização (vs. estática)
// - Linguagem compilada moderna (vs. BASIC interpretado)
//
// A implementação mantém fidelidade às fórmulas matemáticas originais
// de perspectiva cônica enquanto oferece uma experiência de usuário
// contemporânea.
package main

import (
	"fmt"
	"log"
	"os"

	"representacao-figuras/internal/core"
	"representacao-figuras/internal/renderer"
	"representacao-figuras/internal/viewer"
)

// main é o ponto de entrada da aplicação.
//
// Implementa uma interface de linha de comando que oferece diferentes
// modos de operação para trabalhar com figuras 3D:
//
// 1. generate: Cria imagens PNG estáticas
// 2. view: Abre interface interativa
// 3. help: Exibe informações de uso
//
// A aplicação também oferece compatibilidade com uso direto
// (sem especificar comando) para facilidade de uso.
func main() {
	// === VALIDAÇÃO DE ARGUMENTOS ===
	if len(os.Args) < 2 {
		showHelp()
		os.Exit(1)
	}

	// Primeiro argumento é o comando (ou nome do arquivo)
	command := os.Args[1]

	// === PROCESSAMENTO DE COMANDOS ===
	switch command {
	// Comando para geração de imagens PNG
	case "generate", "gen", "png":
		if len(os.Args) < 3 {
			fmt.Println("Erro: especifique o arquivo YAML")
			fmt.Println("Uso: figuras3d generate <arquivo.yaml>")
			os.Exit(1)
		}
		// Executa geração de PNG estático
		generatePNG(os.Args[2])

	// Comando para visualização interativa
	case "view", "viewer", "show":
		if len(os.Args) < 3 {
			fmt.Println("Erro: especifique o arquivo YAML")
			fmt.Println("Uso: figuras3d view <arquivo.yaml>")
			os.Exit(1)
		}
		// Abre interface gráfica interativa
		openViewer(os.Args[2])

	// Comando de ajuda
	case "help", "--help", "-h":
		showHelp()

	// === MODO COMPATIBILIDADE ===
	// Se não é um comando reconhecido, tenta interpretar como arquivo
	default:
		// Caso especial: --viewer como primeiro argumento
		if command == "--viewer" && len(os.Args) >= 3 {
			openViewer(os.Args[2])
		} else {
			// Assume que o primeiro argumento é um arquivo YAML
			// Comportamento padrão: gera PNG
			generatePNG(command)
		}
	}
}

// showHelp exibe informações de uso da aplicação.
//
// Apresenta os comandos disponíveis, exemplos de uso e créditos
// ao artigo original de 1982, mantendo a conexão histórica.
func showHelp() {
	// Cabeçalho com créditos ao artigo original
	fmt.Println("Representação de Figuras por Computador")
	fmt.Println("Baseado no artigo de Luiz Antonio Pereira")
	fmt.Println("MICRO SISTEMAS - Novembro/1982")
	fmt.Println("")

	// Lista de comandos principais
	fmt.Println("Comandos:")
	fmt.Println("  generate <arquivo.yaml>    Gera imagem PNG (salva em output/)")
	fmt.Println("  view <arquivo.yaml>        Abre viewfinder interativo")
	fmt.Println("  help                       Mostra esta ajuda")
	fmt.Println("")

	// Exemplos práticos de uso
	fmt.Println("Exemplos:")
	fmt.Println("  figuras3d generate samples/cubo.yaml")
	fmt.Println("  figuras3d view samples/casa.yaml")
	fmt.Println("")

	// Atalhos e conveniências
	fmt.Println("Atalhos:")
	fmt.Println("  figuras3d gen samples/cubo.yaml       # Mesmo que generate")
	fmt.Println("  figuras3d samples/cubo.yaml           # Gera PNG (padrão)")
}

// openViewer inicia a interface gráfica interativa.
//
// Permite visualizar e manipular figuras 3D em tempo real,
// oferecendo uma experiência muito superior ao HP-85 original
// que só podia mostrar imagens estáticas.
//
// Parâmetros:
//   yamlFile: caminho para o arquivo de definição da figura
func openViewer(yamlFile string) {
	fmt.Printf("Abrindo viewfinder para: %s\n", yamlFile)

	// Cria e executa a interface gráfica
	gui := viewer.NewGUI(yamlFile)
	gui.Run()
}

// generatePNG executa o processo completo de geração de imagem estática.
//
// Esta função implementa o pipeline completo descrito no artigo:
// 1. Carregamento da figura (substitui arrays hardcoded do BASIC)
// 2. Configuração da câmera (observador, distância, dimensões)
// 3. Aplicação da projeção cônica
// 4. Renderização em alta qualidade
// 5. Export para arquivo moderno (PNG vs. tela do HP-85)
//
// Parâmetros:
//   yamlFile: caminho para o arquivo de definição da figura
func generatePNG(yamlFile string) {
	fmt.Printf("Gerando PNG para: %s\n", yamlFile)

	// === ETAPA 1: CARREGAMENTO DA FIGURA ===
	// Substitui a definição hardcoded do BASIC original
	figura, err := core.LoadFigureFromYAML(yamlFile)
	if err != nil {
		log.Fatalf("Erro ao carregar arquivo YAML: %v", err)
	}

	// Informações sobre a figura carregada
	fmt.Printf("Renderizando figura: %s\n", figura.Nome)
	fmt.Printf("Pontos 3D: %d\n", len(figura.Pontos))
	fmt.Printf("Linhas: %d\n", len(figura.Linhas))

	// === ETAPA 2: CONFIGURAÇÃO DE DIMENSÕES ===
	// Define tamanho da tela de saída (muito superior ao HP-85: 256×192)
	width, height := 800, 600 // Resolução padrão moderna

	// Permite customização via configurações no YAML
	if figura.Render != nil {
		if figura.Render.CanvasWidth > 0 {
			width = figura.Render.CanvasWidth
		}
		if figura.Render.CanvasHeight > 0 {
			height = figura.Render.CanvasHeight
		}
	}

	// === ETAPA 3: CONFIGURAÇÃO VISUAL ===
	// Converte configurações YAML para formato interno do renderizador
	renderCfg, err := renderer.ConfigFromFigure(figura)
	if err != nil {
		log.Fatalf("Erro na configuração de renderização: %v", err)
	}

	// === ETAPA 4: INICIALIZAÇÃO DO RENDERIZADOR ===
	// Cria o contexto gráfico com a resolução especificada
	r := renderer.New(width, height)

	// === ETAPA 5: CONFIGURAÇÃO DA CÂMERA ===
	// Define os parâmetros fundamentais da perspectiva cônica
	// (observador V, distância R, dimensões L1 e L2)
	r.SetCamera(figura.Camera)

	// === ETAPA 6: RENDERIZAÇÃO ===
	// Aplica as transformações 3D→2D e desenha a figura
	err = r.RenderFigureWithConfig(figura, renderCfg)
	if err != nil {
		log.Fatalf("Erro ao renderizar figura: %v", err)
	}

	// === ETAPA 7: EXPORT ===
	// Salva o resultado em arquivo PNG (tecnologia inexistente em 1982!)
	outputFile := fmt.Sprintf("output/%s.png", figura.Nome)
	err = r.SaveImage(outputFile)
	if err != nil {
		log.Fatalf("Erro ao salvar imagem: %v", err)
	}

	// Confirmação de sucesso e dica de uso
	fmt.Printf("Imagem salva: %s\n", outputFile)
	fmt.Println("Dica: Use 'figuras3d view' para visualizar interativo!")
}
