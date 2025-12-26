package main

import (
	"dmensions/internal/database"
	"dmensions/internal/math"
	"dmensions/internal/model"
	"dmensions/internal/utils"
	"fmt"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type AppState struct {
	app    fyne.App
	window fyne.Window
	store  *database.Storage

	// UI Elements
	skyContainer *fyne.Container
	statusLabel  *widget.Label
	inputEntry   *widget.Entry

	// Navigation State
	zoomLevel float64
	panOffset fyne.Position

	// Data
	currentMap map[int64]model.Point2D
	allWords   []model.WordData
}

type GestureDetector struct {
	widget.BaseWidget
	onDrag   func(id *fyne.DragEvent)
	onScroll func(id *fyne.ScrollEvent)
}

func (g *GestureDetector) Dragged(e *fyne.DragEvent)    { g.onDrag(e) }
func (g *GestureDetector) DragEnd()                     {}
func (g *GestureDetector) Scrolled(e *fyne.ScrollEvent) { g.onScroll(e) }
func (g *GestureDetector) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewRectangle(color.Transparent))
}

func main() {
	utils.Splash()

	dbConn, err := database.InitDB()
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer dbConn.Close()

	store := database.NewStorage(dbConn)

	myApp := app.New()
	state := &AppState{
		app:       myApp,
		window:    myApp.NewWindow("Dmensions: Vector Galaxy"),
		store:     store,
		zoomLevel: 0.8,
		panOffset: fyne.NewPos(0, 0),
	}

	state.window.Resize(fyne.NewSize(1200, 800))
	state.buildUI()
	go state.refreshData()
	state.window.ShowAndRun()
}

func (s *AppState) buildUI() {
	s.skyContainer = container.NewWithoutLayout()
	skyBG := canvas.NewRectangle(color.RGBA{15, 15, 25, 255})

	gestureLayer := &GestureDetector{
		onDrag: func(e *fyne.DragEvent) {
			s.panOffset.X += e.Dragged.DX
			s.panOffset.Y += e.Dragged.DY
			s.drawSky()
		},
		onScroll: func(e *fyne.ScrollEvent) {
			mousePos := e.Position

			zoomStep := 1.1
			if e.Scrolled.DY < 0 {
				zoomStep = 0.9
			}

			s.panOffset.X = mousePos.X - (mousePos.X-s.panOffset.X)*float32(zoomStep)
			s.panOffset.Y = mousePos.Y - (mousePos.Y-s.panOffset.Y)*float32(zoomStep)

			s.zoomLevel *= zoomStep

			if s.zoomLevel < 0.01 {
				s.zoomLevel = 0.01
			}
			if s.zoomLevel > 120.0 {
				s.zoomLevel = 120.0
			}

			s.drawSky()
		},
	}
	gestureLayer.ExtendBaseWidget(gestureLayer)

	skyStack := container.NewStack(skyBG, s.skyContainer, gestureLayer)

	// Sidebar
	s.statusLabel = widget.NewLabel("Ready.")
	s.inputEntry = widget.NewEntry()
	s.inputEntry.SetPlaceHolder("Add concept...")
	s.inputEntry.OnSubmitted = func(str string) { s.addWord(str) }

	resetBtn := widget.NewButton("Reset View", func() {
		s.zoomLevel = 1.0
		s.panOffset = fyne.NewPos(0, 0)
		s.drawSky()
	})

	sidebar := container.NewVBox(
		widget.NewLabelWithStyle("DMENSIONS", fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewSeparator(),
		s.inputEntry,
		widget.NewButton("Inject", func() { s.addWord(s.inputEntry.Text) }),
		layout.NewSpacer(),
		resetBtn,
		widget.NewSeparator(),
		s.statusLabel,
	)

	split := container.NewHSplit(container.NewPadded(sidebar), skyStack)
	split.SetOffset(0.2)
	s.window.SetContent(split)
}

func (s *AppState) refreshData() {
	fyne.Do(func() { s.statusLabel.SetText("Simulating t-SNE...") })
	words, _ := s.store.GetAllVectors()
	s.allWords = words

	mapping := math.ProjectTSNE(words, 500, 20.0)

	fyne.Do(func() {
		s.currentMap = mapping
		s.drawSky()
		s.statusLabel.SetText(fmt.Sprintf("%d Entities Active", len(words)))
	})
}

func (s *AppState) drawSky() {
	s.skyContainer.Objects = nil
	if len(s.currentMap) == 0 {
		return
	}

	minX, maxX, minY, maxY := s.getBounds()

	w := float64(s.skyContainer.Size().Width)
	h := float64(s.skyContainer.Size().Height)

	for id, point := range s.currentMap {
		rx, ry := maxX-minX, maxY-minY
		if rx == 0 {
			rx = 1
		}
		if ry == 0 {
			ry = 1
		}

		nx := (point.X - minX) / rx
		ny := (point.Y - minY) / ry

		screenX := float32(nx*w*s.zoomLevel) + s.panOffset.X
		screenY := float32(ny*h*s.zoomLevel) + s.panOffset.Y

		if screenX < -100 || screenX > s.skyContainer.Size().Width+100 ||
			screenY < -100 || screenY > s.skyContainer.Size().Height+100 {
			continue
		}

		star := canvas.NewCircle(color.RGBA{80, 200, 255, 255})
		star.Resize(fyne.NewSize(8, 8))
		star.Move(fyne.NewPos(screenX-3, screenY-3))

		txt := canvas.NewText(s.getWordLabel(id), color.White)
		txt.TextSize = 10
		txt.Move(fyne.NewPos(screenX+10, screenY-5))

		s.skyContainer.Add(star)
		s.skyContainer.Add(txt)
	}
	s.skyContainer.Refresh()
}

func (s *AppState) addWord(word string) {
	if word == "" {
		return
	}
	s.inputEntry.Disable()
	go func() {
		s.store.SaveWord(word)
		s.refreshData()
		fyne.Do(func() { s.inputEntry.Enable(); s.inputEntry.SetText("") })
	}()
}

func (s *AppState) getBounds() (minX, maxX, minY, maxY float64) {
	first := true
	for _, p := range s.currentMap {
		if first {
			minX, maxX, minY, maxY = p.X, p.X, p.Y, p.Y
			first = false
		} else {
			if p.X < minX {
				minX = p.X
			}
			if p.X > maxX {
				maxX = p.X
			}
			if p.Y < minY {
				minY = p.Y
			}
			if p.Y > maxY {
				maxY = p.Y
			}
		}
	}
	return
}

func (s *AppState) getWordLabel(id int64) string {
	for _, w := range s.allWords {
		if w.ID == id {
			return w.Word
		}
	}
	return "?"
}
