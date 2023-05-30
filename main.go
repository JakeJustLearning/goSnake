package main

import (
	"fmt"
	"image/color"
	"math/rand"

	ebiten "github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

const (
	screenWidth  = 320
	screenHeight = 240
	tileSize     = 5
)

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body        []Point
	Direction   Point
	GrowCounter int
}

type Food struct {
	Position Point
}

type Game struct {
	snake    *Snake
	food     *Food
	score    int
	gameOver bool
	// ticks         int
	updateCounter int
	speed         int
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func NewFood() *Food {
	return &Food{
		Position: Point{
			X: rand.Intn(screenWidth / tileSize),
			Y: rand.Intn(screenHeight / tileSize),
		},
	}
}

func NewSnake() *Snake {
	return &Snake{Body: []Point{
		{X: screenWidth / tileSize / 2, Y: screenHeight / tileSize / 2},
	},
		Direction: Point{X: 1, Y: 0},
	}
}

func (snake *Snake) Move() {
	newHead := Point{
		X: snake.Body[0].X + snake.Direction.X,
		Y: snake.Body[0].Y + snake.Direction.Y,
	}

	snake.Body = append([]Point{newHead}, snake.Body...)

	if snake.GrowCounter > 0 {
		snake.GrowCounter--
	} else {
		snake.Body = snake.Body[:len(snake.Body)-1]
	}

}

func (g *Game) restart() {
	g.snake = NewSnake()
	g.score = 0
	g.gameOver = false
	g.food = NewFood()
	g.speed = 10
}

func (g *Game) UpdateSnakeDirection() {
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: -1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyRight) && g.snake.Direction.X == 0 {
		g.snake.Direction = Point{X: 1, Y: 0}
	} else if ebiten.IsKeyPressed(ebiten.KeyUp) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: -1}
	} else if ebiten.IsKeyPressed(ebiten.KeyDown) && g.snake.Direction.Y == 0 {
		g.snake.Direction = Point{X: 0, Y: 1}
	}

}

// initiates a gameOver state
func (g *Game) MakeGameOver() {
	g.gameOver = true
	g.speed = 10
}

func (g *Game) CheckFoodConsumption(head Point) {
	if head.X == g.food.Position.X && head.Y == g.food.Position.Y {
		g.score++
		g.snake.GrowCounter += 1
		g.food = NewFood()

		if g.speed > 2 {
			g.speed--
		}
	}
}

func (g *Game) CheckCollisions(head Point) {
	if head.X < 0 || head.Y < 0 || head.X >= screenWidth/tileSize || head.Y >= screenHeight/tileSize {
		g.MakeGameOver()
	}
	for _, part := range g.snake.Body[1:] {
		if head.X == part.X && head.Y == part.Y {
			g.MakeGameOver()
		}
	}

	g.CheckFoodConsumption(head)

}

func (g *Game) Update() error {
	if g.gameOver {
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.restart()
		}
		return nil
	}

	g.updateCounter++
	if g.updateCounter < g.speed {
		return nil
	}
	g.updateCounter = 0

	g.snake.Move()

	g.UpdateSnakeDirection()

	head := g.snake.Body[0]
	g.CheckCollisions(head)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 255})

	for _, section := range g.snake.Body {
		vector.DrawFilledRect(
			screen,
			float32(section.X*tileSize),
			float32(section.Y*tileSize),
			tileSize,
			tileSize,
			color.RGBA{0, 255, 0, 255},
			false,
		)
	}

	vector.DrawFilledCircle(
		screen,
		float32(g.food.Position.X*tileSize),
		float32(g.food.Position.Y*tileSize),
		tileSize,
		color.RGBA{244, 12, 0, 23},
		false,
	)

	face := basicfont.Face7x13

	if g.gameOver {
		text.Draw(screen, "Game Over", face, screenWidth/2-40, screenHeight/2, color.White)
		text.Draw(screen, "Press 'R' to restart", face, screenWidth/2-60, screenHeight/2+16, color.White)
	}

	scoreText := fmt.Sprintf("Score: %d", g.score)
	text.Draw(screen, scoreText, face, 5, screenHeight-5, color.White)
}

func main() {
	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Snake Extra?")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
