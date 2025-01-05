package terrain

import (
	"fmt"
	"github.com/ojrac/opensimplex-go"
	"go.uber.org/zap"
	"math/rand"
	"pirate-wars/cmd/common"
	"pirate-wars/cmd/town"
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
	logger  *zap.SugaredLogger
	props   Props
	World   MapView
	MiniMap MapView
}

func Init(logger *zap.SugaredLogger) *Terrain {
	logger.Info("Initializing terrain")
	//default values for terrain map generation
	worldGrid := make([][]TerrainType, common.WorldHeight)
	for i := range worldGrid {
		worldGrid[i] = make([]TerrainType, common.WorldHeight)
	}

	// Calculate MiniMap dimensions
	height := len(worldGrid) / common.MiniMapFactor
	width := len(worldGrid[0]) / common.MiniMapFactor

	// Create new 2D slice
	miniMap := make([][]TerrainType, height+1)
	for i := range miniMap {
		miniMap[i] = make([]TerrainType, width+1)
	}

	return &Terrain{
		logger: logger,
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
			grid:      miniMap,
		},
	}
}

func (t *Terrain) genTownCoords() common.Coordinates {
	coords := common.Coordinates{X: rand.Intn(common.WorldWidth), Y: rand.Intn(common.WorldHeight)}
	t.logger.Debug(fmt.Sprintf("Generating random town coordinates: %v", coords))
	return coords
}

func (t *Terrain) generateTowns(fn func() common.Coordinates) {
	t.logger.Info("Initializing %v towns", common.TotalTowns)
	for i := 0; i <= common.TotalTowns; i++ {
		for {
			coords := fn()
			if coords.X > 1 && coords.Y > 1 &&
				coords.X < common.WorldWidth-1 && coords.Y < common.WorldHeight &&
				t.World.grid[coords.X][coords.Y] == TypeBeach {

				if t.World.isAdjacentToWater(coords) {
					town.Create(coords, '⩎')
					t.World.grid[coords.X][coords.Y] = TypeTown
					// grow towns
					for _, a := range t.World.GetAdjacentCoords(coords) {
						if (t.World.grid[a.X][a.Y] == TypeLowland || t.World.grid[a.X][a.Y] == TypeBeach) && t.World.isAdjacentToWater(a) {
							t.World.grid[a.X][a.Y] = TypeTown
						}
					}
					break
				}
			}
		}
	}
}

func (t *Terrain) GenerateTowns() {
	t.generateTowns(t.genTownCoords)
}

func (t *Terrain) GenerateWorld() {
	t.logger.Info("Initializing world")
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

			//normalize to -1 to 1, and then from 0 to 1 (this is for the ability to use grayscale, if using colors could keep from -1 to 1)
			var s = (total/normalizeOctaves + 1) / 2
			if s > 0.60 {
				t.World.grid[x][y] = TypeDeepWater
			} else if s > 0.46 {
				t.World.grid[x][y] = TypeOpenWater
			} else if s > 0.42 {
				t.World.grid[x][y] = TypeShallowWater
			} else if s > 0.40 {
				t.World.grid[x][y] = TypeBeach
			} else if s > 0.31 {
				t.World.grid[x][y] = TypeLowland
			} else if s > 0.26 {
				t.World.grid[x][y] = TypeHighland
			} else if s > 0.21 {
				t.World.grid[x][y] = TypeRock
			} else {
				t.World.grid[x][y] = TypePeak
			}
		}
	}
	t.GenerateMiniMap()
}

func (t *Terrain) GenerateMiniMap() {
	// Down-sample
	for i, row := range t.World.grid {
		for j, val := range row {
			// Calculate corresponding index in new slice
			newI := i / common.MiniMapFactor
			newJ := j / common.MiniMapFactor

			// Assign original value
			t.MiniMap.grid[newI][newJ] = val
		}
	}
}

func (t *Terrain) RandomPositionDeepWater() common.Coordinates {
	for {
		coords := common.Coordinates{X: rand.Intn(common.WorldWidth), Y: rand.Intn(common.WorldHeight)}
		if t.World.grid[coords.X][coords.Y] == TypeDeepWater {
			return coords
		}
	}
}
