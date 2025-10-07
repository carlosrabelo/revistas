package renderer

import (
	"testing"

	"representacao-figuras/pkg/types"
)

func TestDefaultRenderConfig(t *testing.T) {
	config := DefaultRenderConfig()

	// Verifica valores padrão
	if config.Background.R != 1 || config.Background.G != 1 || config.Background.B != 1 {
		t.Errorf("Expected white background (1,1,1), got (%f,%f,%f)",
			config.Background.R, config.Background.G, config.Background.B)
	}

	if config.LineColor.R != 0 || config.LineColor.G != 0 || config.LineColor.B != 0 {
		t.Errorf("Expected black lines (0,0,0), got (%f,%f,%f)",
			config.LineColor.R, config.LineColor.G, config.LineColor.B)
	}

	if config.LineWidth != 1.0 {
		t.Errorf("Expected line width=1.0, got %f", config.LineWidth)
	}

	if config.ShowVertices {
		t.Error("Expected ShowVertices=false by default")
	}

	if config.ShowLabels {
		t.Error("Expected ShowLabels=false by default")
	}
}

func TestConfigFromFigure_Nil(t *testing.T) {
	config, err := ConfigFromFigure(nil)
	if err != nil {
		t.Errorf("ConfigFromFigure with nil figure should not error: %v", err)
	}

	// Deve retornar configuração padrão
	defaultConfig := DefaultRenderConfig()
	if config.LineWidth != defaultConfig.LineWidth {
		t.Error("Should return default config for nil figure")
	}
}

func TestConfigFromFigure_NoRenderSettings(t *testing.T) {
	figure := &types.Figure{
		Nome:   "test",
		Render: nil,
	}

	config, err := ConfigFromFigure(figure)
	if err != nil {
		t.Errorf("ConfigFromFigure should not error with no render settings: %v", err)
	}

	// Deve retornar configuração padrão
	defaultConfig := DefaultRenderConfig()
	if config.LineWidth != defaultConfig.LineWidth {
		t.Error("Should return default config when no render settings")
	}
}

func TestConfigFromFigure_WithSettings(t *testing.T) {
	showVertices := true
	showLabels := false

	figure := &types.Figure{
		Nome: "test",
		Render: &types.RenderSettings{
			Background:   "black",
			LineColor:    "#ff0000",
			LineWidth:    2.5,
			VertexColor:  "#0000ff",
			ShowVertices: &showVertices,
			ShowLabels:   &showLabels,
		},
	}

	config, err := ConfigFromFigure(figure)
	if err != nil {
		t.Errorf("ConfigFromFigure failed: %v", err)
	}

	// Verifica se as configurações foram aplicadas
	if config.Background.R != 0 || config.Background.G != 0 || config.Background.B != 0 {
		t.Errorf("Expected black background, got (%f,%f,%f)",
			config.Background.R, config.Background.G, config.Background.B)
	}

	if config.LineColor.R != 1 || config.LineColor.G != 0 || config.LineColor.B != 0 {
		t.Errorf("Expected red lines, got (%f,%f,%f)",
			config.LineColor.R, config.LineColor.G, config.LineColor.B)
	}

	if config.LineWidth != 2.5 {
		t.Errorf("Expected line width=2.5, got %f", config.LineWidth)
	}

	if !config.ShowVertices {
		t.Error("Expected ShowVertices=true")
	}

	if config.ShowLabels {
		t.Error("Expected ShowLabels=false")
	}
}

func TestParseColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected colorRGB
		wantErr  bool
	}{
		{
			name:     "named color white",
			input:    "white",
			expected: colorRGB{R: 1, G: 1, B: 1},
			wantErr:  false,
		},
		{
			name:     "named color black",
			input:    "black",
			expected: colorRGB{R: 0, G: 0, B: 0},
			wantErr:  false,
		},
		{
			name:     "hex color with #",
			input:    "#ff0000",
			expected: colorRGB{R: 1, G: 0, B: 0},
			wantErr:  false,
		},
		{
			name:     "hex color without #",
			input:    "00ff00",
			expected: colorRGB{R: 0, G: 1, B: 0},
			wantErr:  false,
		},
		{
			name:     "short hex color",
			input:    "#f0f",
			expected: colorRGB{R: 1, G: 0, B: 1},
			wantErr:  false,
		},
		{
			name:     "case insensitive",
			input:    "WHITE",
			expected: colorRGB{R: 1, G: 1, B: 1},
			wantErr:  false,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "invalid hex",
			input:   "#xyz",
			wantErr: true,
		},
		{
			name:    "invalid length",
			input:   "#12345",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseColor(tt.input)

			if tt.wantErr && err == nil {
				t.Error("Expected error, got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}

			if !tt.wantErr {
				tolerance := 0.001
				if abs(result.R-tt.expected.R) > tolerance ||
					abs(result.G-tt.expected.G) > tolerance ||
					abs(result.B-tt.expected.B) > tolerance {
					t.Errorf("Expected (%f,%f,%f), got (%f,%f,%f)",
						tt.expected.R, tt.expected.G, tt.expected.B,
						result.R, result.G, result.B)
				}
			}
		})
	}
}

func TestParseHexComponent(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
		wantErr  bool
	}{
		{"00", 0.0, false},
		{"ff", 1.0, false},
		{"80", 0.502, false}, // 128/255 ≈ 0.502
		{"zz", 0.0, true},    // Invalid hex
	}

	for _, tt := range tests {
		result, err := parseHexComponent(tt.input)

		if tt.wantErr && err == nil {
			t.Errorf("Input '%s': expected error, got nil", tt.input)
		}

		if !tt.wantErr && err != nil {
			t.Errorf("Input '%s': expected no error, got: %v", tt.input, err)
		}

		if !tt.wantErr {
			tolerance := 0.01
			if abs(result-tt.expected) > tolerance {
				t.Errorf("Input '%s': expected %f, got %f", tt.input, tt.expected, result)
			}
		}
	}
}

// Helper function for absolute value
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}