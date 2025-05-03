package main

import (
	"bytes"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/text/language"
)

const (
	screenWidth  = 800
	screenHeight = 400
	defaultBarW  = 20
	defaultBarH  = 100
)

type Ball struct {
	Radius     float32
	PosX       float32
	PosY       float32
	DirectionX int
	DirectionY int
	Speed      int
}

func NewBall() *Ball {
	return &Ball{
		Radius:     10,
		PosX:       screenWidth / 2,
		PosY:       screenHeight / 2,
		Speed:      5,
		DirectionX: 1,
		DirectionY: 1,
	}
}

func (b *Ball) Draw(screen *ebiten.Image) {
	vector.DrawFilledCircle(screen, b.PosX, b.PosY, b.Radius, color.White, false)
}

func (b *Ball) IsBallCollidingBorder() bool {

	// Going down: 1
	// Going up: -1
	return (b.DirectionY == 1 && b.PosY+b.Radius >= screenHeight) || (b.DirectionY == -1 && b.PosY-b.Radius <= 0)
}

func (b *Ball) Update() {
	b.PosX += float32(b.DirectionX) * float32(b.Speed)
	b.PosY += float32(b.DirectionY) * float32(b.Speed)
}

func (b *Ball) IsCollidingBar(bar *Bar, barSide int) bool {
	var ballClosePosX float32
	var barSidePosX float32
	// Left side
	if barSide == 1 {
		ballClosePosX = b.PosX - b.Radius
		barSidePosX = bar.PosX + defaultBarW
	}

	if barSide == -1 {
		ballClosePosX = b.PosX + b.Radius
		barSidePosX = bar.PosX
	}
	hitCond := (bar.PosY <= b.PosY && bar.PosY+defaultBarH >= b.PosY) &&
		(barSidePosX == ballClosePosX)
	return hitCond
}

type Bar struct {
	W           float32
	H           float32
	PosX        float32
	PosY        float32
	Speed       float32
	MoveUpKey   ebiten.Key
	MoveDownKey ebiten.Key
}

func (b *Bar) Draw(screen *ebiten.Image) {
	vector.DrawFilledRect(screen, b.PosX, b.PosY, b.W, b.H, color.White, false)
}

func (b *Bar) Update() {
	if ebiten.IsKeyPressed(b.MoveDownKey) {
		if (b.PosY + defaultBarH + b.Speed) <= screenHeight {
			b.PosY += b.Speed
		}
	}
	if ebiten.IsKeyPressed(b.MoveUpKey) {
		if b.PosY-b.Speed >= 0 {
			b.PosY -= b.Speed
		}
	}
}

type Game struct {
	Ball      *Ball
	LeftBar   *Bar
	RightBar  *Bar
	PlayerWin int
	Font      text.Face
}

func (g *Game) Update() error {
	if g.PlayerWin != 0 {
		return nil
	}

	g.LeftBar.Update()
	g.RightBar.Update()

	if g.Ball.IsCollidingBar(g.LeftBar, 1) {
		g.Ball.DirectionX = -g.Ball.DirectionX
	}

	if g.Ball.IsCollidingBar(g.RightBar, -1) {
		g.Ball.DirectionX = -g.Ball.DirectionX
	}

	if g.Ball.IsBallCollidingBorder() {
		g.Ball.DirectionY = -g.Ball.DirectionY
	}

	g.Ball.Update()

	// Check for condition
	if g.Ball.PosX-g.Ball.Radius == 0 {
		// Right hand wins
		g.PlayerWin = 1
	}
	if g.Ball.PosX+g.Ball.Radius == screenWidth {
		// Left hand wins
		g.PlayerWin = -1
	}

	return nil
}
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ebiten.WindowSize()
}

func NewBar(posX float32, posY float32, upKey ebiten.Key, downKey ebiten.Key) *Bar {
	return &Bar{
		W:           defaultBarW,
		H:           defaultBarH,
		PosX:        posX,
		PosY:        posY,
		MoveUpKey:   upKey,
		Speed:       10,
		MoveDownKey: downKey,
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{16, 12, 122, 255})
	g.Ball.Draw(screen)
	g.LeftBar.Draw(screen)
	g.RightBar.Draw(screen)

	if g.PlayerWin != 0 {
		var victoryString string
		if g.PlayerWin == -1 {
			victoryString = "Player 1 wins"
		}
		if g.PlayerWin == 1 {
			victoryString = "Player 2 wins"
		}
		textOpts := text.DrawOptions{}
		text.Draw(screen, victoryString, g.Font, &textOpts)
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ball := NewBall()
	defaultBarYPos := float32((screenHeight - defaultBarH) / 2)
	leftBar := NewBar(20, defaultBarYPos, ebiten.KeyW, ebiten.KeyS)
	rightBar := NewBar(screenWidth-20-defaultBarW, defaultBarYPos, ebiten.KeyUp, ebiten.KeyDown)

	var textSource *text.GoTextFaceSource
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}

	textSource = s

	font := text.GoTextFace{
		Size:     24,
		Source:   textSource,
		Language: language.English,
	}
	game := Game{
		Ball:      ball,
		LeftBar:   leftBar,
		RightBar:  rightBar,
		PlayerWin: 0,
		Font:      &font,
	}

	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
