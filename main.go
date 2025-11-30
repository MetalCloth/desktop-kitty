package main

// GoCat: Animated cat desktop overlay that follows mouse, can be picked up and thrown.


import (
	"fmt"
	"math/rand"
	"math"
	"time"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	frameWidth  = 101
	frameHeight = 104
)

var catFrames []*ebiten.Image
var img *ebiten.Image
var choice int

type Game struct {
	x, y            float64
	frameCount      int
	hold            bool
	frameIndex      int
	isHover         bool
	idle_time       float64
	reactionTimer   int
	facingRight     bool
	direction       int
	vx, vy          float64 // velocity of the kitty
	lastMx, lastMy  int 
}

func (g *Game) IsHovering() bool {
	// check if cursor is precisely between the cat head
	mx, my := ebiten.CursorPosition()
	if (float64(mx) >= g.x && float64(mx) <= g.x+frameWidth) && (float64(my) >= g.y && float64(my) <= g.y+frameHeight) {
		if (float64(mx) >= (27+g.x) && float64(mx) < (68+g.x)) && (float64(my) < (40+g.y)) {
			return true
		}
	}
	return false
}

func (g *Game) Update() error {
	mx, my := ebiten.CursorPosition()
	sw, sh := ebiten.Monitor().Size()

	catCenterX := g.x + (float64(frameWidth) / 2.0)
	catCenterY := g.y + (float64(frameHeight) / 2.0)

	isClicking := ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft)

	if isClicking && g.IsHovering() {
		g.hold = true
	}
	if !isClicking {
		g.hold = false
	}

	if g.hold {
		g.idle_time = 0
		ebiten.SetWindowMousePassthrough(false)

		g.x = float64(mx) - 47.5
		g.y = float64(my) - 25
		g.vx = float64(mx - g.lastMx) * 1.5
		g.vy = float64(my - g.lastMy) * 1.5
		g.lastMx = mx
		g.lastMy = my

		g.reactionTimer = 0

		return nil
	}

	if math.Abs(g.vx) > 1 || math.Abs(g.vy) > 1 {

		// meths for physics
		g.x += g.vx
		g.y += g.vy

		g.vy += 0.5
		g.vx *= 0.95
		g.vy *= 0.95

		if g.y > float64(sh)-float64(frameHeight) {
			g.y = float64(sh) - float64(frameHeight)
			g.vy = -g.vy * 0.75
		}
		if g.y < 0 {
			g.y = 0
			g.vy = -g.vy * 0.75
		}
		if g.x < 0 {
			g.x = 0
			g.vx = -g.vx * 0.75
		}
		if g.x > float64(sw)-float64(frameWidth) {
			g.x = float64(sw) - float64(frameWidth)
			g.vx = -g.vx * 0.75
		}

		g.lastMx = mx
		g.lastMy = my
		return nil
	}

	dx := float64(mx) - catCenterX
	dy := float64(my) - catCenterY
	if math.Abs(dx) > math.Abs(dy) {
		if dx > 0 {
			g.facingRight = true
			g.direction = 1
		} else {
			g.facingRight = false
			g.direction = 0
		}
	} else {
		if dy > 0 {
			g.direction = 3
		} else {
			g.direction = 2
		}
	}

	g.frameCount++
	if g.frameCount > 8 { // fps
		g.frameIndex++
		if g.frameIndex >= 4 {
			g.frameIndex = 0
		}
		g.frameCount = 0
	}

	if g.IsHovering() {
		g.isHover = true
		ebiten.SetWindowMousePassthrough(false)
	} else {
		g.isHover = false
		ebiten.SetWindowMousePassthrough(true)
	}

	distance := math.Sqrt(dx*dx + dy*dy)

	if distance > 60 {
		g.idle_time = 0
		g.reactionTimer++
		if g.reactionTimer > 20 {
			speed := 4.1
			g.x += (dx / distance) * speed
			g.y += (dy / distance) * speed
		}
	} else {
		g.reactionTimer = 0
	}

	if g.x < 0 {
		g.x = 0
	} else if g.x > float64(sw)-float64(frameWidth) {
		g.x = float64(sw) - float64(frameWidth)
	}
	if g.y < 0 {
		g.y = 0
	} else if g.y > float64(sh)-float64(frameHeight) {
		g.y = float64(sh) - float64(frameHeight)
	}

	g.lastMx = mx
	g.lastMy = my

	return nil
}

func init() {
	// loading the images in the memory could to prevent constantly extracting from folder 
	for i := 0; i < 41; i++ {
		path := fmt.Sprintf("kitties/tile%03d.png", i)

		loadedImg, _, err := ebitenutil.NewImageFromFile(path)
		if err != nil {
			fmt.Println("Warning: Could not load", path)
			continue
		}

		catFrames = append(catFrames, loadedImg)
	}

}

func (g *Game) Draw(screen *ebiten.Image) {
	tilt := 0.0

	op := &ebiten.DrawImageOptions{}

	if math.Abs(g.vx) > 1.0 || math.Abs(g.vy) > 1.0 {
		if len(catFrames) > 28 {
			// hanging cat
			img = catFrames[28]
		}
		tilt = g.vx * 0.05
	} else if g.reactionTimer > 20 {
		x := g.frameIndex
		switch g.direction {
		case 0:
			x += 12
		case 1:
			x += 4
		case 2:
			x += 8
		case 3:
			x += 0
		}
		if x < len(catFrames) {
			img = catFrames[x]
		}
	} else if g.isHover {
		g.idle_time = 0
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {

			img = catFrames[28]
		} else {
			x := g.frameIndex + 37
			if x < len(catFrames) {
				img = catFrames[x]
			}
		}
	} else {
		if g.idle_time == 0 {
			rand.Seed(time.Now().UnixNano())
			choice = rand.Intn(2)
		}
		g.idle_time += 0.01
		if g.idle_time > 3.5 {
			x := g.frameIndex + 29
			if choice == 1 {
				x = g.frameIndex + 33
			}
			if x < len(catFrames) {
				img = catFrames[x]
			}
		} else {
			img = catFrames[17]
		}
	}

	// adding physics effect when hanging
	op.GeoM.Scale(1, 1)
	op.GeoM.Translate(-float64(frameWidth)/2, -float64(frameHeight)/2)
	op.GeoM.Rotate(tilt)
	op.GeoM.Translate(float64(frameWidth)/2, float64(frameHeight)/2)
	op.GeoM.Translate(g.x, g.y)

	screen.DrawImage(img, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	sw, sh := ebiten.Monitor().Size()
	ebiten.SetWindowSize(sw, sh-1)
	ebiten.SetWindowPosition(0, 0)
	ebiten.SetWindowTitle("GoCat")
	ebiten.SetWindowDecorated(false)
	ebiten.SetWindowFloating(true)
	ebiten.SetWindowMousePassthrough(true)

	opts := &ebiten.RunGameOptions{
		ScreenTransparent: true,
		SkipTaskbar:      true,
	}

	if err := ebiten.RunGameWithOptions(&Game{}, opts); err != nil {
		log.Fatal(err)
	}
}