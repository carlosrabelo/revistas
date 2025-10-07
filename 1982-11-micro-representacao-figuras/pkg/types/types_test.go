package types

import (
	"testing"
)

func TestPoint3D(t *testing.T) {
	p := Point3D{X: 1.0, Y: 2.0, Z: 3.0, Nome: "test"}

	if p.X != 1.0 {
		t.Errorf("Expected X=1.0, got %f", p.X)
	}
	if p.Y != 2.0 {
		t.Errorf("Expected Y=2.0, got %f", p.Y)
	}
	if p.Z != 3.0 {
		t.Errorf("Expected Z=3.0, got %f", p.Z)
	}
	if p.Nome != "test" {
		t.Errorf("Expected Nome='test', got '%s'", p.Nome)
	}
}

func TestPoint2D(t *testing.T) {
	p := Point2D{X: 10.5, Y: 20.7}

	if p.X != 10.5 {
		t.Errorf("Expected X=10.5, got %f", p.X)
	}
	if p.Y != 20.7 {
		t.Errorf("Expected Y=20.7, got %f", p.Y)
	}
}

func TestLine(t *testing.T) {
	line := Line{P1: 0, P2: 1}

	if line.P1 != 0 {
		t.Errorf("Expected P1=0, got %d", line.P1)
	}
	if line.P2 != 1 {
		t.Errorf("Expected P2=1, got %d", line.P2)
	}
}

func TestDefaultCamera(t *testing.T) {
	camera := DefaultCamera()

	// Verifica valores padrão baseados no artigo
	if camera.Observer.X != 0 || camera.Observer.Y != 0 || camera.Observer.Z != 0 {
		t.Errorf("Expected observer at origin, got (%f, %f, %f)",
			camera.Observer.X, camera.Observer.Y, camera.Observer.Z)
	}

	if camera.Distance != 10 {
		t.Errorf("Expected distance=10, got %f", camera.Distance)
	}

	if camera.Width != 12.8 {
		t.Errorf("Expected width=12.8 (HP-85 based), got %f", camera.Width)
	}

	if camera.Height != 9.6 {
		t.Errorf("Expected height=9.6 (HP-85 based), got %f", camera.Height)
	}
}

func TestFigure(t *testing.T) {
	// Cria uma figura simples para teste
	figure := Figure{
		Nome: "test_figure",
		Pontos: []Point3D{
			{X: 0, Y: 5, Z: 0, Nome: "origin"},
			{X: 1, Y: 5, Z: 1, Nome: "corner"},
		},
		Linhas: []Line{
			{P1: 0, P2: 1},
		},
		Camera: DefaultCamera(),
	}

	if figure.Nome != "test_figure" {
		t.Errorf("Expected nome='test_figure', got '%s'", figure.Nome)
	}

	if len(figure.Pontos) != 2 {
		t.Errorf("Expected 2 points, got %d", len(figure.Pontos))
	}

	if len(figure.Linhas) != 1 {
		t.Errorf("Expected 1 line, got %d", len(figure.Linhas))
	}

	// Verifica se a linha referencia pontos válidos
	line := figure.Linhas[0]
	if line.P1 >= len(figure.Pontos) || line.P2 >= len(figure.Pontos) {
		t.Errorf("Line references invalid points: P1=%d, P2=%d, total points=%d",
			line.P1, line.P2, len(figure.Pontos))
	}
}

func TestRenderSettings(t *testing.T) {
	// Testa configurações de renderização
	showVertices := true
	showLabels := false

	settings := RenderSettings{
		CanvasWidth:  800,
		CanvasHeight: 600,
		Background:   "white",
		LineColor:    "#000000",
		LineWidth:    2.0,
		VertexColor:  "red",
		ShowVertices: &showVertices,
		ShowLabels:   &showLabels,
	}

	if settings.CanvasWidth != 800 {
		t.Errorf("Expected CanvasWidth=800, got %d", settings.CanvasWidth)
	}

	if settings.LineWidth != 2.0 {
		t.Errorf("Expected LineWidth=2.0, got %f", settings.LineWidth)
	}

	if settings.ShowVertices == nil || !*settings.ShowVertices {
		t.Error("Expected ShowVertices=true")
	}

	if settings.ShowLabels == nil || *settings.ShowLabels {
		t.Error("Expected ShowLabels=false")
	}
}