# Representação de Figuras por Computador

**Artigo Original:** Luiz Antonio Pereira
**Revista:** MICRO SISTEMAS - Novembro/1982
**Reimplementação:** Go 1.23 (2024)

## 📖 Sobre o Artigo Original

Este artigo, publicado na MICRO SISTEMAS edição 014 de novembro de 1982, apresenta os fundamentos matemáticos para representação gráfica de figuras tridimensionais usando **perspectiva cônica** em microcomputadores.

### Conceitos Abordados

- **Perspectiva Cônica**: Simula a visão real do objeto através de projeção matemática
- **Sistema de Coordenadas XYZ**: Representação tridimensional de pontos no espaço
- **Transformações Geométricas**: Conversão de coordenadas 3D para coordenadas 2D da tela
- **Projeção por Segmentos**: Figuras definidas por linhas conectando vértices

### Implementação Original

O artigo foi desenvolvido para o **HP-85** com:
- Resolução de 256×192 pixels
- Linguagem BASIC
- Sistema de plotagem por pontos
- Funções matemáticas para transformação de coordenadas

### Fórmulas Matemáticas

O artigo apresenta as equações fundamentais da projeção cônica:

```
Para cada ponto P(X,Y,Z):
- Translação: P' = P - V (onde V é o observador)
- Projeção: x = P'x * R/P'z, y = P'y * R/P'z
- Conversão: coordenadas de tela baseadas nas dimensões L1 e L2
```

## 🚀 Implementação Moderna em Go

Esta reimplementação moderniza os conceitos do artigo usando:

- **Go 1.23**: Linguagem moderna e performática
- **Arquivos YAML**: Definição declarativa de figuras 3D
- **Biblioteca gg**: Renderização gráfica de alta qualidade
- **PNG Export**: Saída em formato moderno
- **Estruturas Tipadas**: Segurança de tipos e clareza de código

### Estrutura do Projeto

```
microsistemas/1982-11-representacao-figuras/
├── cmd/figuras3d/main.go  # Ponto de entrada do programa
├── internal/              # Lógica interna da aplicação
│   ├── core/             # Carregamento de modelos
│   ├── renderer/         # Engine de renderização 3D
│   └── viewer/           # Interface gráfica
├── pkg/types/            # Definições de tipos (Point3D, Figure, Camera)
├── modelos/              # Modelos 3D de exemplo
│   ├── cubo.yaml        # Cubo 3D simples
│   ├── casa.yaml        # Casa com telhado, porta e janela
│   ├── piramide.yaml    # Pirâmide triangular
│   ├── estrela.yaml     # Estrela 3D
│   └── escada.yaml      # Escada em degraus
├── go.mod               # Dependências do projeto
├── Makefile             # Comandos de build e execução
└── README.md            # Este arquivo
```

## 🎯 Como Usar

### Instalação

```bash
cd microsistemas/1982-11-representacao-figuras
go mod tidy
```

### Execução

```bash
# Renderizar o cubo de exemplo
make generate FILE=modelos/cubo.yaml

# Renderizar a casa de exemplo
make generate FILE=modelos/casa.yaml

# Ou usando go run diretamente
go run cmd/figuras3d/main.go generate modelos/cubo.yaml
```

### Criar Suas Próprias Figuras

Crie um arquivo YAML seguindo a estrutura:

```yaml
nome: minha_figura
pontos:
  - {x: 0, y: 0, z: 0, nome: "origem"}
  - {x: 2, y: 0, z: 0, nome: "P1"}
  # ... mais pontos

linhas:
  - {p1: 0, p2: 1}  # Conecta ponto 0 ao ponto 1
  # ... mais linhas

camera:
  observador: {x: 0, y: 0, z: 0}
  distancia: 15
  largura: 12.8
  altura: 9.6
```

## 📊 Exemplos Incluídos

### Cubo (`modelos/cubo.yaml`)
- Cubo 3D básico com 8 vértices
- Demonstra faces, arestas e perspectiva
- Ideal para entender os conceitos fundamentais

### Casa (`modelos/casa.yaml`)
- Casa com telhado, porta e janela
- Estrutura mais complexa inspirada nas figuras do artigo
- Mostra diferentes tipos de formas geométricas

## 🔧 Parâmetros da Câmera

- **observador**: Posição do observador no espaço 3D
- **distancia**: Distância R do plano projetante (afeta perspectiva)
- **largura/altura**: Dimensões da "tela virtual" (baseadas no HP-85 original)

## 📐 Diferenças da Implementação Original

| Aspecto | Original (1982) | Moderno (2024) |
|---------|----------------|-----------------|
| Linguagem | BASIC | Go 1.23 |
| Entrada | Código hardcoded | Arquivos YAML |
| Saída | Tela 256×192 | PNG de alta resolução |
| Estruturas | Arrays simples | Tipos estruturados |
| Validação | Manual | Automática |

## 🎓 Valor Educacional

Esta implementação preserva os **conceitos matemáticos fundamentais** do artigo original enquanto demonstra:

- **Evolução das linguagens**: De BASIC interpretado para Go compilado
- **Melhores práticas modernas**: Separação de responsabilidades, tipagem forte
- **Flexibilidade**: Sistema configurável vs código hardcoded
- **Qualidade gráfica**: Alta resolução vs limitações de hardware dos anos 80

## 🏗️ Extensões Possíveis

- [ ] Suporte a cores diferentes para linhas
- [ ] Animações rotacionando objetos
- [ ] Exportação para formatos SVG
- [ ] Interface web interativa
- [ ] Lighting e shading básicos
- [ ] Importação de modelos 3D simples

## 📚 Referências

- **Artigo original**: "Representação de figuras por computador" - Luiz Antonio Pereira, MICRO SISTEMAS, Nov/1982
- **HP-85**: Computador pessoal da HP usado na implementação original
- **Perspectiva Cônica**: Técnica fundamental de computação gráfica

---

**Homenagem ao conhecimento atemporal da computação brasileira dos anos 80** 🇧🇷