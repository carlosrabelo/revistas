package core

import (
	"os"
	"path/filepath"
	"testing"

	"representacao-figuras/pkg/types"
)

func TestLoadFigureFromYAML(t *testing.T) {
	// Cria um arquivo YAML temporário para teste
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_figure.yaml")

	yamlContent := `nome: test_cube
pontos:
  - {x: -1, y: 5, z: -1, nome: "A"}
  - {x:  1, y: 5, z: -1, nome: "B"}
  - {x:  1, y: 5, z:  1, nome: "C"}
  - {x: -1, y: 5, z:  1, nome: "D"}

linhas:
  - {p1: 0, p2: 1}
  - {p1: 1, p2: 2}
  - {p1: 2, p2: 3}
  - {p1: 3, p2: 0}

camera:
  observador: {x: 0, y: 0, z: 0}
  distancia: 8
  largura: 12.8
  altura: 9.6`

	err := os.WriteFile(testFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Testa carregamento
	figure, err := LoadFigureFromYAML(testFile)
	if err != nil {
		t.Fatalf("LoadFigureFromYAML failed: %v", err)
	}

	// Verifica dados carregados
	if figure.Nome != "test_cube" {
		t.Errorf("Expected nome='test_cube', got '%s'", figure.Nome)
	}

	if len(figure.Pontos) != 4 {
		t.Errorf("Expected 4 points, got %d", len(figure.Pontos))
	}

	if len(figure.Linhas) != 4 {
		t.Errorf("Expected 4 lines, got %d", len(figure.Linhas))
	}

	// Verifica primeiro ponto
	p0 := figure.Pontos[0]
	if p0.X != -1 || p0.Y != 5 || p0.Z != -1 || p0.Nome != "A" {
		t.Errorf("Point 0 incorrect: got (%f,%f,%f,'%s'), expected (-1,5,-1,'A')",
			p0.X, p0.Y, p0.Z, p0.Nome)
	}

	// Verifica câmera
	if figure.Camera.Distance != 8 {
		t.Errorf("Expected camera distance=8, got %f", figure.Camera.Distance)
	}
}

func TestLoadFigureFromYAML_WithDefaultCamera(t *testing.T) {
	// Testa arquivo sem especificação de câmera (deve usar padrão)
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test_no_camera.yaml")

	yamlContent := `nome: simple_test
pontos:
  - {x: 0, y: 5, z: 0}

linhas:
  - {p1: 0, p2: 0}`

	err := os.WriteFile(testFile, []byte(yamlContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	figure, err := LoadFigureFromYAML(testFile)
	if err != nil {
		t.Fatalf("LoadFigureFromYAML failed: %v", err)
	}

	// Deve usar câmera padrão
	defaultCam := types.DefaultCamera()
	if figure.Camera.Distance != defaultCam.Distance {
		t.Errorf("Expected default camera distance=%f, got %f",
			defaultCam.Distance, figure.Camera.Distance)
	}
}

func TestLoadFigureFromYAML_FileNotFound(t *testing.T) {
	_, err := LoadFigureFromYAML("nonexistent_file.yaml")
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestLoadFigureFromYAML_InvalidYAML(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "invalid.yaml")

	// YAML inválido
	err := os.WriteFile(testFile, []byte("invalid: yaml: content: ["), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err = LoadFigureFromYAML(testFile)
	if err == nil {
		t.Error("Expected error for invalid YAML, got nil")
	}
}

func TestValidateFigure(t *testing.T) {
	tests := []struct {
		name    string
		figure  types.Figure
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid figure",
			figure: types.Figure{
				Nome: "valid",
				Pontos: []types.Point3D{
					{X: 0, Y: 5, Z: 0},
					{X: 1, Y: 5, Z: 1},
				},
				Linhas: []types.Line{
					{P1: 0, P2: 1},
				},
			},
			wantErr: false,
		},
		{
			name: "no points",
			figure: types.Figure{
				Nome:   "no_points",
				Pontos: []types.Point3D{},
				Linhas: []types.Line{{P1: 0, P2: 1}},
			},
			wantErr: true,
			errMsg:  "pelo menos um ponto",
		},
		{
			name: "no lines",
			figure: types.Figure{
				Nome:   "no_lines",
				Pontos: []types.Point3D{{X: 0, Y: 5, Z: 0}},
				Linhas: []types.Line{},
			},
			wantErr: true,
			errMsg:  "pelo menos uma linha",
		},
		{
			name: "invalid line reference",
			figure: types.Figure{
				Nome: "invalid_line",
				Pontos: []types.Point3D{
					{X: 0, Y: 5, Z: 0},
				},
				Linhas: []types.Line{
					{P1: 0, P2: 5}, // P2 inválido
				},
			},
			wantErr: true,
			errMsg:  "ponto P2 inválido",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFigure(&tt.figure)

			if tt.wantErr && err == nil {
				t.Errorf("Expected error containing '%s', got nil", tt.errMsg)
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if tt.wantErr && err != nil {
				// Verifica se a mensagem de erro contém o texto esperado
				if tt.errMsg != "" && !containsString(err.Error(), tt.errMsg) {
					t.Errorf("Expected error message to contain '%s', got '%s'",
						tt.errMsg, err.Error())
				}
			}
		})
	}
}

// Helper function to check if string contains substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) &&
		   (s == substr || (len(s) > len(substr) &&
		   	findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}