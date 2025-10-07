package renderer

import (
	"math"
	"testing"

	"representacao-figuras/pkg/types"
)

func TestNew(t *testing.T) {
	width, height := 800, 600
	renderer := New(width, height)

	if renderer.width != width {
		t.Errorf("Expected width=%d, got %d", width, renderer.width)
	}

	if renderer.height != height {
		t.Errorf("Expected height=%d, got %d", height, renderer.height)
	}

	expectedCenterX := float64(width) / 2
	if renderer.centerX != expectedCenterX {
		t.Errorf("Expected centerX=%f, got %f", expectedCenterX, renderer.centerX)
	}

	expectedCenterY := float64(height) / 2
	if renderer.centerY != expectedCenterY {
		t.Errorf("Expected centerY=%f, got %f", expectedCenterY, renderer.centerY)
	}
}

func TestSetCamera(t *testing.T) {
	renderer := New(800, 600)
	camera := types.Camera{
		Observer: types.Point3D{X: 1, Y: 2, Z: 3},
		Distance: 10,
		Width:    12.8,
		Height:   9.6,
	}

	renderer.SetCamera(camera)

	if renderer.camera.Observer.X != 1 || renderer.camera.Observer.Y != 2 || renderer.camera.Observer.Z != 3 {
		t.Errorf("Camera observer not set correctly")
	}

	if renderer.camera.Distance != 10 {
		t.Errorf("Expected camera distance=10, got %f", renderer.camera.Distance)
	}
}

func TestProjectPoint(t *testing.T) {
	renderer := New(800, 600)

	// Câmera na origem olhando para frente
	camera := types.Camera{
		Observer: types.Point3D{X: 0, Y: 0, Z: 0},
		Distance: 10,
		Width:    12.8,
		Height:   9.6,
	}
	renderer.SetCamera(camera)

	tests := []struct {
		name     string
		point3D  types.Point3D
		expected types.Point2D
		tolerance float64
	}{
		{
			name:    "point at center depth",
			point3D: types.Point3D{X: 0, Y: 5, Z: 0}, // Centro, profundidade 5
			expected: types.Point2D{X: 400, Y: 300}, // Centro da tela
			tolerance: 1.0,
		},
		{
			name:    "point to the right",
			point3D: types.Point3D{X: 1, Y: 5, Z: 0}, // Direita, profundidade 5
			expected: types.Point2D{X: 525, Y: 300}, // Direita do centro (corrigido)
			tolerance: 2.0,
		},
		{
			name:    "point above center",
			point3D: types.Point3D{X: 0, Y: 5, Z: 1}, // Acima, profundidade 5
			expected: types.Point2D{X: 400, Y: 175}, // Acima do centro (corrigido)
			tolerance: 2.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := renderer.ProjectPoint(tt.point3D)

			if math.Abs(result.X-tt.expected.X) > tt.tolerance {
				t.Errorf("X projection: expected %f±%f, got %f",
					tt.expected.X, tt.tolerance, result.X)
			}

			if math.Abs(result.Y-tt.expected.Y) > tt.tolerance {
				t.Errorf("Y projection: expected %f±%f, got %f",
					tt.expected.Y, tt.tolerance, result.Y)
			}
		})
	}
}

func TestProjectPoint_PerspectiveFormula(t *testing.T) {
	// Testa se a fórmula de perspectiva cônica está correta
	// Baseado nas equações do artigo: x = Px * R/Pz, y = Py * R/Pz

	renderer := New(800, 600)
	camera := types.Camera{
		Observer: types.Point3D{X: 0, Y: 0, Z: 0},
		Distance: 10, // R
		Width:    12.8,
		Height:   9.6,
	}
	renderer.SetCamera(camera)

	// Ponto de teste
	point := types.Point3D{X: 2, Y: 5, Z: 1} // Px=2, Pz=5, Py=1

	result := renderer.ProjectPoint(point)

	// Calcula manualmente conforme as equações do artigo
	px := point.X - camera.Observer.X // = 2
	py := point.Z - camera.Observer.Z // = 1 (nota: Z é altura no nosso sistema)
	pz := point.Y - camera.Observer.Y // = 5 (nota: Y é profundidade)

	expectedProjX := px * camera.Distance / pz // = 2 * 10 / 5 = 4
	expectedProjY := py * camera.Distance / pz // = 1 * 10 / 5 = 2

	// Converte para coordenadas de tela
	scaleX := float64(renderer.width) / camera.Width   // 800/12.8 = 62.5
	scaleY := float64(renderer.height) / camera.Height // 600/9.6 = 62.5

	expectedScreenX := renderer.centerX + (expectedProjX * scaleX) // 400 + (4 * 62.5) = 650
	expectedScreenY := renderer.centerY - (expectedProjY * scaleY) // 300 - (2 * 62.5) = 175

	tolerance := 1.0
	if math.Abs(result.X-expectedScreenX) > tolerance {
		t.Errorf("Perspective formula X: expected %f, got %f (diff: %f)",
			expectedScreenX, result.X, math.Abs(result.X-expectedScreenX))
	}

	if math.Abs(result.Y-expectedScreenY) > tolerance {
		t.Errorf("Perspective formula Y: expected %f, got %f (diff: %f)",
			expectedScreenY, result.Y, math.Abs(result.Y-expectedScreenY))
	}
}

func TestProjectPoint_BehindCamera(t *testing.T) {
	// Testa pontos atrás da câmera (pz <= 0)
	renderer := New(800, 600)
	camera := types.Camera{
		Observer: types.Point3D{X: 0, Y: 0, Z: 0},
		Distance: 10,
		Width:    12.8,
		Height:   9.6,
	}
	renderer.SetCamera(camera)

	// Ponto atrás da câmera
	point := types.Point3D{X: 1, Y: -1, Z: 1} // Y negativo = atrás

	result := renderer.ProjectPoint(point)

	// Deve usar pz=0.1 para evitar divisão por zero
	// A projeção deve ser extrema mas finita
	if math.IsInf(result.X, 0) || math.IsNaN(result.X) {
		t.Error("X projection should not be infinite or NaN for points behind camera")
	}

	if math.IsInf(result.Y, 0) || math.IsNaN(result.Y) {
		t.Error("Y projection should not be infinite or NaN for points behind camera")
	}
}

func TestRenderFigure_EmptyPoints(t *testing.T) {
	renderer := New(800, 600)

	figure := &types.Figure{
		Nome:   "empty",
		Pontos: []types.Point3D{},
		Linhas: []types.Line{},
		Camera: types.DefaultCamera(),
	}

	err := renderer.RenderFigure(figure)
	if err == nil {
		t.Error("Expected error for figure with no points")
	}
}

func TestRenderFigure_ValidFigure(t *testing.T) {
	renderer := New(800, 600)

	figure := &types.Figure{
		Nome: "simple_line",
		Pontos: []types.Point3D{
			{X: 0, Y: 5, Z: 0},
			{X: 1, Y: 5, Z: 1},
		},
		Linhas: []types.Line{
			{P1: 0, P2: 1},
		},
		Camera: types.DefaultCamera(),
	}

	renderer.SetCamera(figure.Camera)
	err := renderer.RenderFigure(figure)
	if err != nil {
		t.Errorf("RenderFigure failed for valid figure: %v", err)
	}
}

func TestRenderFigure_InvalidLineReference(t *testing.T) {
	renderer := New(800, 600)

	figure := &types.Figure{
		Nome: "invalid_line",
		Pontos: []types.Point3D{
			{X: 0, Y: 5, Z: 0},
		},
		Linhas: []types.Line{
			{P1: 0, P2: 5}, // P2 inválido
		},
		Camera: types.DefaultCamera(),
	}

	renderer.SetCamera(figure.Camera)
	err := renderer.RenderFigure(figure)
	// Não deve dar erro - linhas inválidas são ignoradas
	if err != nil {
		t.Errorf("RenderFigure should handle invalid line references gracefully: %v", err)
	}
}

func TestAddGrid(t *testing.T) {
	renderer := New(200, 150)

	// Testa se AddGrid não causa panic
	renderer.AddGrid()

	// Não há muito o que testar além de não dar panic
	// A funcionalidade visual seria testada manualmente
}