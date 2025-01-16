package terrain

import (
	"fmt"
	"github.com/ojrac/opensimplex-go"
	"go.uber.org/zap"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/screen"
)

type Props struct {
	width       int
	height      int
	scale       float64
	lacunarity  float64
	persistence float64
	octaves     int
}

type Terrain struct {
	Logger   *zap.SugaredLogger
	props    Props
	World    MapView
	MiniMap  MapView
	mapItems []MapItem
}

type MapItem interface {
	GetPos() common.Coordinates
	GetTerrainType() TerrainType
}

func (t *Terrain) SetMapItem(m MapItem) {
	t.mapItems = append(t.mapItems, m)
}

func Init(logger *zap.SugaredLogger) *Terrain {
	logger.Info(fmt.Sprintf(fmt.Sprintf("Initializing terrain - height: %v width: %v", common.WorldHeight, common.WorldWidth)))
	//default values for terrain map generation
	worldGrid := make([][]TerrainType, common.WorldHeight)
	for i := range worldGrid {
		worldGrid[i] = make([]TerrainType, common.WorldHeight)
	}

	return &Terrain{
		Logger: logger,
		props: Props{
			width:       common.WorldWidth,
			height:      common.WorldHeight,
			scale:       60,
			lacunarity:  2.0,
			persistence: 0.5,
			octaves:     5,
		},
		World: MapView{
			isMiniMap: false,
			logger:    logger,
			grid:      worldGrid,
		},
		MiniMap: MapView{
			isMiniMap: true,
			logger:    logger,
			grid:      [][]TerrainType{},
		},
		mapItems: []MapItem{},
	}
}

func (t *Terrain) GenerateWorld() {
	t.Logger.Info("Initializing world...")
	noise := opensimplex.New(rand.Int63())

	for x := 0; x < t.props.width; x++ {
		for y := 0; y < t.props.height; y++ {
			//sample x and y and apply scale
			xFloat := float64(x) / t.props.scale
			yFloat := float64(y) / t.props.scale

			//init values for octave calculation
			frequency := 1.0
			amplitude := 1.0
			normalizeOctaves := 0.0
			total := 0.0

			//octave calculation
			for i := 0; i < t.props.octaves; i++ {
				total += noise.Eval2(xFloat*frequency, yFloat*frequency) * amplitude
				normalizeOctaves += amplitude
				amplitude *= t.props.persistence
				frequency *= t.props.lacunarity
			}

			c := common.Coordinates{
				X: x,
				Y: y,
			}
			//normalize to -1 to 1, and then from 0 to 1 (this is for the ability to use grayscale, if using colors could keep from -1 to 1)
			var s = (total/normalizeOctaves + 1) / 2
			if s > 0.59 {
				t.World.SetPositionType(c, TypeDeepWater)
			} else if s > 0.44 {
				t.World.SetPositionType(c, TypeOpenWater)
			} else if s > 0.42 {
				t.World.SetPositionType(c, TypeShallowWater)
			} else if s > 0.40 {
				t.World.SetPositionType(c, TypeBeach)
			} else if s > 0.31 {
				t.World.SetPositionType(c, TypeLowland)
			} else if s > 0.26 {
				t.World.SetPositionType(c, TypeHighland)
			} else if s > 0.21 {
				t.World.SetPositionType(c, TypeRock)
			} else {
				t.World.SetPositionType(c, TypePeak)
			}
		}
	}
}

func (t *Terrain) GenerateMiniMap() {
	// Calculate MiniMap dimensions
	height := len(t.World.grid) / screen.MiniMapFactor
	width := len(t.World.grid[0]) / screen.MiniMapFactor

	// Create new 2D slice
	t.MiniMap.grid = make([][]TerrainType, height+1)
	for i := range t.MiniMap.grid {
		t.MiniMap.grid[i] = make([]TerrainType, width+1)
	}

	// Down-sample
	for i, row := range t.World.grid {
		for j, val := range row {
			// Calculate corresponding index in new slice
			c := common.Coordinates{
				X: i / screen.MiniMapFactor,
				Y: j / screen.MiniMapFactor,
			}
			// Assign original TerrainType value
			t.MiniMap.SetPositionType(c, val)
		}
	}
	for _, m := range t.mapItems {
		t.MiniMap.SetPositionType(common.GetMiniMapScale(m.GetPos()), m.GetTerrainType())
	}
}
