// Package viewer implementa interface gráfica interativa para visualização
// de figuras 3D em tempo real.
//
// Esta é uma extensão significativa do artigo original de 1982, que só
// podia mostrar imagens estáticas na tela do HP-85. A interface permite:
//
// - Manipulação em tempo real dos parâmetros da câmera
// - Visualização imediata das mudanças de perspectiva
// - Controles intuitivos para posição do observador
// - Ajuste dinâmico da distância e dimensões
//
// A implementação usa a biblioteca Fyne para criar uma interface
// moderna e responsiva, mantendo os cálculos matemáticos originais
// da perspectiva cônica.
package viewer

import (
	"fmt"
	"image"
	"strconv"

	"representacao-figuras/internal/core"
	"representacao-figuras/internal/renderer"
	"representacao-figuras/pkg/types"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

// GUI gerencia a interface gráfica interativa.
//
// Esta estrutura encapsula todos os elementos necessários para
// criar uma experiência de visualização 3D interativa, permitindo
// ao usuário explorar figuras de diferentes ângulos e distâncias
// em tempo real.
type GUI struct {
	app          fyne.App
	window       fyne.Window
	figura       *types.Figure
	filename     string
	renderCfg    renderer.RenderConfig
	canvasWidth  int
	canvasHeight int

	// Controles da câmera
	camXEntry *widget.Entry
	camYEntry *widget.Entry
	camZEntry *widget.Entry
	distEntry *widget.Entry

	// Área de visualização
	imageCanvas *canvas.Image
	statusLabel *widget.Label
}

// NewGUI cria uma nova instância do visualizador GUI
func NewGUI(filename string) *GUI {
	myApp := app.New()

	window := myApp.NewWindow("MICRO SISTEMAS - Representação de Figuras 3D")
	window.Resize(fyne.NewSize(1200, 800))
	window.CenterOnScreen()

	// Permite fechar a janela normalmente
	window.SetOnClosed(func() {
		myApp.Quit()
	})

	viewer := &GUI{
		app:          myApp,
		window:       window,
		filename:     filename,
		canvasWidth:  800,
		canvasHeight: 600,
		renderCfg:    renderer.DefaultRenderConfig(),
	}

	viewer.setupUI()
	viewer.loadFigure()

	return viewer
}

// setupUI configura a interface do usuário
func (v *GUI) setupUI() {
	// Título estilo anos 80
	title := widget.NewLabelWithStyle("REPRESENTAÇÃO DE FIGURAS 3D", fyne.TextAlignCenter, fyne.TextStyle{Bold: true})
	subtitle := widget.NewLabelWithStyle("Baseado no artigo da MICRO SISTEMAS - Nov/1982", fyne.TextAlignCenter, fyne.TextStyle{Italic: true})

	// Área de visualização
	v.imageCanvas = canvas.NewImageFromImage(image.NewRGBA(image.Rect(0, 0, v.canvasWidth, v.canvasHeight)))
	v.imageCanvas.FillMode = canvas.ImageFillOriginal

	// Controles de câmera
	v.camXEntry = widget.NewEntry()
	v.camXEntry.SetText("1")
	v.camYEntry = widget.NewEntry()
	v.camYEntry.SetText("1")
	v.camZEntry = widget.NewEntry()
	v.camZEntry.SetText("0")
	v.distEntry = widget.NewEntry()
	v.distEntry.SetText("4")

	// Labels e controles
	cameraForm := container.NewGridWithColumns(2,
		widget.NewLabel("Observador X:"), v.camXEntry,
		widget.NewLabel("Observador Y:"), v.camYEntry,
		widget.NewLabel("Observador Z:"), v.camZEntry,
		widget.NewLabel("Distância:"), v.distEntry,
	)

	// Botões
	renderBtn := widget.NewButton("🔄 Renderizar", v.renderFigure)
	reloadBtn := widget.NewButton("📁 Recarregar", v.loadFigure)
	saveBtn := widget.NewButton("💾 Salvar PNG", v.savePNG)

	buttonBox := container.NewHBox(renderBtn, reloadBtn, saveBtn)

	// Status
	v.statusLabel = widget.NewLabel("Carregando...")

	// Painel de controles
	controlPanel := container.NewVBox(
		title,
		subtitle,
		widget.NewSeparator(),
		widget.NewLabelWithStyle("CONTROLES DE CÂMERA", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		cameraForm,
		buttonBox,
		widget.NewSeparator(),
		v.statusLabel,
	)

	// Layout principal
	content := container.NewHSplit(
		container.NewScroll(v.imageCanvas),
		controlPanel,
	)
	content.SetOffset(0.7) // 70% para imagem, 30% para controles

	v.window.SetContent(content)
}

// loadFigure carrega a figura do arquivo YAML
func (v *GUI) loadFigure() {
	figura, err := core.LoadFigureFromYAML(v.filename)
	if err != nil {
		v.statusLabel.SetText(fmt.Sprintf("Erro: %v", err))
		dialog.ShowError(err, v.window)
		return
	}

	v.figura = figura

	// Configura dimensões do canvas com base na figura
	v.canvasWidth = 800
	v.canvasHeight = 600
	if figura.Render != nil {
		if figura.Render.CanvasWidth > 0 {
			v.canvasWidth = figura.Render.CanvasWidth
		}
		if figura.Render.CanvasHeight > 0 {
			v.canvasHeight = figura.Render.CanvasHeight
		}
	}

	v.imageCanvas.Image = image.NewRGBA(image.Rect(0, 0, v.canvasWidth, v.canvasHeight))
	v.imageCanvas.Refresh()

	cfg, err := renderer.ConfigFromFigure(figura)
	if err != nil {
		v.statusLabel.SetText(fmt.Sprintf("Configuração inválida: %v", err))
		dialog.ShowError(err, v.window)
		cfg = renderer.DefaultRenderConfig()
	}
	v.renderCfg = cfg
	v.updateCameraControls()
	v.renderFigure()

	v.statusLabel.SetText(fmt.Sprintf("Figura: %s | Pontos: %d | Linhas: %d",
		figura.Nome, len(figura.Pontos), len(figura.Linhas)))
}

// updateCameraControls atualiza os controles com os valores da câmera
func (v *GUI) updateCameraControls() {
	if v.figura == nil {
		return
	}

	cam := v.figura.Camera
	v.camXEntry.SetText(fmt.Sprintf("%.1f", cam.Observer.X))
	v.camYEntry.SetText(fmt.Sprintf("%.1f", cam.Observer.Y))
	v.camZEntry.SetText(fmt.Sprintf("%.1f", cam.Observer.Z))
	v.distEntry.SetText(fmt.Sprintf("%.1f", cam.Distance))
}

// getCameraFromControls lê os valores dos controles
func (v *GUI) getCameraFromControls() types.Camera {
	cam := v.figura.Camera

	if x, err := strconv.ParseFloat(v.camXEntry.Text, 64); err == nil {
		cam.Observer.X = x
	}
	if y, err := strconv.ParseFloat(v.camYEntry.Text, 64); err == nil {
		cam.Observer.Y = y
	}
	if z, err := strconv.ParseFloat(v.camZEntry.Text, 64); err == nil {
		cam.Observer.Z = z
	}
	if d, err := strconv.ParseFloat(v.distEntry.Text, 64); err == nil {
		cam.Distance = d
	}

	return cam
}

// renderFigure renderiza a figura com os parâmetros atuais
func (v *GUI) renderFigure() {
	if v.figura == nil {
		return
	}

	// Atualiza câmera com valores dos controles
	v.figura.Camera = v.getCameraFromControls()

	// Cria renderizador
	r := renderer.New(v.canvasWidth, v.canvasHeight)
	r.SetCamera(v.figura.Camera)

	// Renderiza
	err := r.RenderFigureWithConfig(v.figura, v.renderCfg)
	if err != nil {
		v.statusLabel.SetText(fmt.Sprintf("Erro na renderização: %v", err))
		return
	}

	// Converte para imagem Fyne
	if img, ok := r.GetImage().(image.Image); ok {
		v.imageCanvas.Image = img
		v.imageCanvas.Refresh()
	}

	v.statusLabel.SetText(fmt.Sprintf(
		"Renderizado! | Obs: (%.1f,%.1f,%.1f) | Dist: %.1f | Canvas: %dx%d",
		v.figura.Camera.Observer.X,
		v.figura.Camera.Observer.Y,
		v.figura.Camera.Observer.Z,
		v.figura.Camera.Distance,
		v.canvasWidth,
		v.canvasHeight,
	))
}

// savePNG salva a imagem atual como PNG
func (v *GUI) savePNG() {
	if v.figura == nil {
		return
	}

	outputFile := fmt.Sprintf("output/%s.png", v.figura.Nome)

	// Cria novo renderizador para salvar
	r := renderer.New(v.canvasWidth, v.canvasHeight)
	r.SetCamera(v.figura.Camera)
	r.RenderFigureWithConfig(v.figura, v.renderCfg)

	err := r.SaveImage(outputFile)
	if err != nil {
		dialog.ShowError(err, v.window)
		return
	}

	dialog.ShowInformation("Salvo!", fmt.Sprintf("Imagem salva como %s", outputFile), v.window)
}

// Run inicia o aplicativo
func (v *GUI) Run() {
	v.window.ShowAndRun()
}
