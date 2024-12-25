package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Ship struct {
	x, y                 float64
	angle                float64
	velocityX, velocityY float64
}

type Bullet struct {
	x, y, dx, dy float64
}

type Asteroid struct {
	x, y, dx, dy, radius float64
}

type Game struct {
	ship         Ship
	bullets      []Bullet
	asteroids    []Asteroid
	score        int
	isGameOver   bool
	screenWidth  int
	screenHeight int
	lastBullet   time.Time
	lastAsteroid time.Time
}

func (g *Game) lastBulletTime() time.Time {
	return g.lastBullet
}

func (g *Game) Update() error {
	if g.isGameOver {
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.resetGame() // reset the game when R is pressed
		}
		return nil
	}

	shipSpeedMultiplier := 0.2 // default game value is 0.2

	// ship controls
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.ship.angle -= 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.ship.angle += 0.1
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		g.ship.velocityX += math.Cos(g.ship.angle) * shipSpeedMultiplier
		g.ship.velocityY += math.Sin(g.ship.angle) * shipSpeedMultiplier
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		g.ship.velocityX -= math.Cos(g.ship.angle) * 0.2
		g.ship.velocityY -= math.Sin(g.ship.angle) * 0.2
	}
	g.ship.x += g.ship.velocityX
	g.ship.y += g.ship.velocityY
	g.ship.velocityX *= 0.99 // friction
	g.ship.velocityY *= 0.99

	// using modulo for screen wrap for the ship
	g.ship.x = math.Mod(g.ship.x+float64(g.screenWidth), float64(g.screenWidth))
	g.ship.y = math.Mod(g.ship.y+float64(g.screenHeight), float64(g.screenHeight))

	// shoot bullets
	bulletSpeedMultiplier := 5.0                  // default game value is 5.0
	timeSinceLastBullet := 300 * time.Millisecond // rate of fire // default game value is 300 * time.Millisecond
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if len(g.bullets) == 0 || time.Since(g.lastBulletTime()) > timeSinceLastBullet {
			g.bullets = append(g.bullets, Bullet{
				x:  g.ship.x,
				y:  g.ship.y,
				dx: math.Cos(g.ship.angle) * bulletSpeedMultiplier,
				dy: math.Sin(g.ship.angle) * bulletSpeedMultiplier,
			})
			g.lastBullet = time.Now()
		}
	}

	// update bullets
	newBullets := []Bullet{}
	for _, bullet := range g.bullets {
		bullet.x += bullet.dx
		bullet.y += bullet.dy
		if bullet.x >= 0 && bullet.x < float64(g.screenWidth) && bullet.y >= 0 && bullet.y < float64(g.screenHeight) {
			newBullets = append(newBullets, bullet)
		}
	}
	g.bullets = newBullets

	// generate new asteroids
	asteroidSpawnInterval := 3 * time.Second // default game value is 3 * time.Second
	if time.Since(g.lastAsteroid) > asteroidSpawnInterval {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*2 + 1
		g.asteroids = append(g.asteroids, Asteroid{
			x:      rand.Float64() * float64(g.screenWidth),
			y:      rand.Float64() * float64(g.screenHeight),
			dx:     speed * math.Cos(angle),
			dy:     speed * math.Sin(angle),
			radius: 30,
		})
		g.lastAsteroid = time.Now()
	}

	// update asteroids
	for i := range g.asteroids {
		g.asteroids[i].x += g.asteroids[i].dx
		g.asteroids[i].y += g.asteroids[i].dy
		g.asteroids[i].x = math.Mod(g.asteroids[i].x+float64(g.screenWidth), float64(g.screenWidth))
		g.asteroids[i].y = math.Mod(g.asteroids[i].y+float64(g.screenHeight), float64(g.screenHeight))
	}

	// check collisions
	for i := len(g.bullets) - 1; i >= 0; i-- {
		for j := len(g.asteroids) - 1; j >= 0; j-- {
			if g.checkCollision(g.bullets[i].x, g.bullets[i].y, g.asteroids[j].x, g.asteroids[j].y, g.asteroids[j].radius) {
				g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
				g.splitAsteroid(j)
				g.score++
				break
			}
		}
	}

	// game over if ship collides with an asteroid
	for _, asteroid := range g.asteroids {
		if g.checkCollision(g.ship.x, g.ship.y, asteroid.x, asteroid.y, asteroid.radius) {
			g.isGameOver = true
			break
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)

	if g.isGameOver {
		face := basicfont.Face7x13
		text.Draw(screen, "Game Over", face, g.screenWidth/2-40, g.screenHeight/2, color.White)
		text.Draw(screen, "Press R to Restart", face, g.screenWidth/2-60, g.screenHeight/2+20, color.White)
		return
	}

	// draw ship
	x1 := g.ship.x + math.Cos(g.ship.angle)*10
	y1 := g.ship.y + math.Sin(g.ship.angle)*10
	x2 := g.ship.x + math.Cos(g.ship.angle+math.Pi*2/3)*10
	y2 := g.ship.y + math.Sin(g.ship.angle+math.Pi*2/3)*10
	x3 := g.ship.x + math.Cos(g.ship.angle-math.Pi*2/3)*10
	y3 := g.ship.y + math.Sin(g.ship.angle-math.Pi*2/3)*10

	ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)

	ebitenutil.DrawLine(screen, x2, y2, x3, y3, color.RGBA{R: 100, G: 100, B: 0, A: 100})

	ebitenutil.DrawLine(screen, x3, y3, x1, y1, color.White)

	// draw bullets
	for _, b := range g.bullets {
		ebitenutil.DrawRect(screen, b.x-2, b.y-2, 4, 4, color.RGBA{R: 100, G: 0, B: 0, A: 100})
	}

	// draw asteroids
	for _, a := range g.asteroids {
		ebitenutil.DrawCircle(screen, a.x, a.y, a.radius, color.White)
	}

	// draw score
	face := basicfont.Face7x13
	text.Draw(screen, "Score: "+fmt.Sprintf("%d", g.score), face, 10, 20, color.White)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return g.screenWidth, g.screenHeight
}

func (g *Game) resetGame() {
	g.ship = Ship{
		x: float64(g.screenWidth) / 2,
		y: float64(g.screenHeight) / 2,
	}
	g.bullets = []Bullet{}
	g.asteroids = []Asteroid{}
	g.score = 0
	g.isGameOver = false

	// Reinitialize asteroids
	for i := 0; i < 10; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*2 + 1
		g.asteroids = append(g.asteroids, Asteroid{
			x:      rand.Float64() * float64(g.screenWidth),
			y:      rand.Float64() * float64(g.screenHeight),
			dx:     speed * math.Cos(angle),
			dy:     speed * math.Sin(angle),
			radius: 30,
		})
	}
}

func (g *Game) checkCollision(x1, y1, x2, y2, radius float64) bool {
	dx := x1 - x2
	dy := y1 - y2
	distance := math.Sqrt(dx*dx + dy*dy)
	return distance < radius
}

func (g *Game) splitAsteroid(index int) {
	asteroid := g.asteroids[index]
	g.asteroids = append(g.asteroids[:index], g.asteroids[index+1:]...)
	if asteroid.radius > 10 {
		speed := math.Sqrt(asteroid.dx*asteroid.dx+asteroid.dy*asteroid.dy) + 0.5
		angle1 := rand.Float64() * 2 * math.Pi
		angle2 := rand.Float64() * 2 * math.Pi
		g.asteroids = append(g.asteroids, Asteroid{
			x: asteroid.x, y: asteroid.y,
			dx: speed * math.Cos(angle1), dy: speed * math.Sin(angle1),
			radius: asteroid.radius / 2,
		}, Asteroid{
			x: asteroid.x, y: asteroid.y,
			dx: speed * math.Cos(angle2), dy: speed * math.Sin(angle2),
			radius: asteroid.radius / 2,
		})
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	game := &Game{
		screenWidth:  1720,
		screenHeight: 1060,
		ship: Ship{
			x: 400, y: 300,
		},
	}

	// initialize asteroids
	for i := 0; i < 5; i++ {
		angle := rand.Float64() * 2 * math.Pi
		speed := rand.Float64()*2 + 1
		game.asteroids = append(game.asteroids, Asteroid{
			x:      rand.Float64() * 800,
			y:      rand.Float64() * 600,
			dx:     speed * math.Cos(angle),
			dy:     speed * math.Sin(angle),
			radius: 30,
		})
	}

	ebiten.SetWindowSize(game.screenWidth, game.screenHeight)
	ebiten.SetWindowTitle("Asteroids")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
