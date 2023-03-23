package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"math/rand"
	"time"

	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"github.com/liyue201/gostl/ds/queue"
	"github.com/liyue201/gostl/ds/stack"
	"github.com/liyue201/gostl/ds/vector"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type Player struct {
	image       pixel.Picture
	playerPos   pixel.Vec
	playerSpeed float64
	points      int
}

type Customer struct {
	order      int
	eatingTime int
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}

func createCustomer(customers *queue.Queue[Customer]) {
	var customer Customer
	customer.order = rand.Intn(2)
	customer.eatingTime = rand.Intn(8)
	customers.Push(customer)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "El Restorante",
		Bounds: pixel.R(0, 0, 640, 480),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	score := text.New(pixel.V(0, 0), basicAtlas)

	var player Player

	player.playerPos = pixel.V(50, 200)
	player.playerSpeed = 130
	player.image, err = loadPicture("player.png")
	if err != nil {
		panic(err)
	}
	playerSprite := pixel.NewSprite(player.image, player.image.Bounds())
	player.points = 0

	plates := stack.New[int]()
	plateImg, err := loadPicture("plate.png")
	if err != nil {
		panic(err)
	}
	plateSprite := pixel.NewSprite(plateImg, plateImg.Bounds())

	pizzaImg, err := loadPicture("pizza.png")
	if err != nil {
		panic(err)
	}
	pizzaSprite := pixel.NewSprite(pizzaImg, pizzaImg.Bounds())
	pizzaPos := pixel.V(32, 400)

	friesImg, err := loadPicture("fries.png")
	if err != nil {
		panic(err)
	}
	friesSprite := pixel.NewSprite(friesImg, friesImg.Bounds())
	friesPos := pixel.V(25, 300)

	burgerImg, err := loadPicture("burger.png")
	if err != nil {
		panic(err)
	}
	burgerSprite := pixel.NewSprite(burgerImg, burgerImg.Bounds())
	burgerPos := pixel.V(32, 130)

	sinkImg, err := loadPicture("sink.png")
	if err != nil {
		panic(err)
	}
	sinkSprite := pixel.NewSprite(sinkImg, sinkImg.Bounds())
	sinkPos := pixel.V(32, 20)

	customerOrderRect := pixel.R(0, 180, 70, 232)
	//Creating array of customer sprites
	tempImg, err := loadPicture("customer.png")
	if err != nil {
		panic(err)
	}
	customerSprite := pixel.NewSprite(tempImg, tempImg.Bounds())
	//Creating customers Queue
	customers := queue.New[Customer]()
	var (
		spawnTime int
		spawning  bool
	)

	ui := imdraw.New(nil)
	sep := imdraw.New(nil)

	sep.Color = colornames.Burlywood
	sep.Push(pixel.R(64, 0, 74, 480).Min, pixel.R(64, 0, 74, 480).Max)
	sep.Rectangle(0)
	ui.Color = color.Black
	ui.Push(pixel.R(74, 0, 650, 100).Min, pixel.R(74, 0, 650, 100).Max)
	ui.Rectangle(0)

	eating := vector.New[int]()
	ordered := Customer{5, 5}

	backgroundImg, err := loadPicture("background-01.png")
	if err != nil {
		panic(err)
	}
	backgroundSprite := pixel.NewSprite(backgroundImg, backgroundImg.Bounds())
	start := false
	gameOver := false
	rules := "How to play\n1. Take an order from the customer (Press Space near customer)\n2. Walk to the correct food icon and press Space to complete the order\n3. Don't forget to wash the dishes by going to the sink and pressing Space\n\nUp and Down Arrow keys move the player \n\nGame ends when there are more than 5 unwashed plates\nor more than 8 customers waiting\n\n Press Enter to start the game"
	fmt.Fprintln(score, rules)
	last := time.Now()
	//RUNTIME
	for !win.Closed() {
		dt := time.Since(last).Seconds()
		last = time.Now()
		win.Clear(color.Black)

		if start {
			backgroundSprite.Draw(win, pixel.IM.ScaledXY(pixel.V(348, 258), pixel.V(0.3, 0.3)))
			pizzaSprite.Draw(win, pixel.IM.Moved(pixel.V(pizzaPos.X, pizzaPos.Y+10)))
			friesSprite.Draw(win, pixel.IM.Moved(pixel.V(friesPos.X, friesPos.Y+10)))
			burgerSprite.Draw(win, pixel.IM.Moved(pixel.V(burgerPos.X, burgerPos.Y+10)))
			sinkSprite.Draw(win, pixel.IM.Moved(pixel.V(sinkPos.X, sinkPos.Y+10)))

			if win.Pressed(pixelgl.KeyUp) {
				player.playerPos.Y += player.playerSpeed * dt
			}
			if win.Pressed(pixelgl.KeyDown) {
				player.playerPos.Y -= player.playerSpeed * dt
			}
			if win.JustPressed(pixelgl.KeySpace) {
				if !customers.Empty() {
					if pizzaSprite.Frame().Moved(pizzaPos).Contains(player.playerPos) {
						fmt.Println("PIZZA ORDERED")
						if ordered.order == 0 {
							player.points++
							score.Clear()
							fmt.Fprintf(score, "Score\n%v", player.points)
							eating.PushBack(ordered.eatingTime * 60)
						} else {
							player.points--
							score.Clear()
							fmt.Fprintf(score, "Score\n%v", player.points)
						}
						customers.Pop()
						ordered = Customer{5, 5}
					}
					if burgerSprite.Frame().Moved(burgerPos).Contains(player.playerPos) {
						fmt.Println("BURGER ORDERED")
						if ordered.order == 1 {
							player.points++
							score.Clear()
							fmt.Fprintf(score, "Score\n%v", player.points)
							eating.PushBack(ordered.eatingTime * 60)
						} else {
							player.points--
							score.Clear()
							fmt.Fprintf(score, "Score\n%v", player.points)
						}
						customers.Pop()
						ordered = Customer{5, 5}
					}
					if friesSprite.Frame().Moved(friesPos).Contains(player.playerPos) {
						fmt.Println("FRIES ORDERED")
						if ordered.order == 2 {
							player.points++
							score.Clear()
							fmt.Fprintf(score, "Score\n%v", player.points)
							eating.PushBack(ordered.eatingTime * 60)
						} else {
							player.points--
							score.Clear()
							fmt.Fprintf(score, "Score\n%v", player.points)
						}
						customers.Pop()
						ordered = Customer{5, 5}
					}
					if sinkSprite.Frame().Moved(sinkPos).Contains(player.playerPos) {
						fmt.Println("A PLATE WAS WASHED")
						plates.Pop()
					}
					if customerOrderRect.Contains(player.playerPos) {
						fmt.Println("ORDER TAKEN")
						ordered = customers.Front().(Customer)
					}
				}

			}
			playerSprite.Draw(win, pixel.IM.Moved(player.playerPos))
			if !spawning {
				spawnTime = rand.Intn(8) * 60
				spawning = true
			} else if spawning {
				if spawnTime != 0 {
					spawnTime -= 1
				} else {
					createCustomer(customers)
					spawning = false
				}
			}

			if customers.Size() != 0 {
				for i := 1; i <= customers.Size(); i++ {
					if i == 1 {
						customerSprite.Draw(win, pixel.IM.Moved(pixel.V(100, 200)))
					} else {
						customerSprite.Draw(win, pixel.IM.Moved(pixel.V(float64(i*32+100), 200)))
					}

				}
			}
			sep.Draw(win)
			ui.Draw(win)
			if !plates.Empty() {
				for i := 1; i <= plates.Size(); i++ {
					j := 16 * i
					plateSprite.Draw(win, pixel.IM.Moved(pixel.V(600, float64(j))))
				}
			}

			if ordered.order == 0 {
				pizzaSprite.Draw(win, pixel.IM.Moved(pixel.V(400, 30)))
			} else if ordered.order == 1 {
				burgerSprite.Draw(win, pixel.IM.Moved(pixel.V(500, 20)))
			} else if ordered.order == 2 {
				friesSprite.Draw(win, pixel.IM.Moved(pixel.V(500, 20)))
			}

			score.Draw(win, pixel.IM.Moved(pixel.V(120, 40)))

			if !eating.Empty() {
				for i := 0; i < eating.Size(); i++ {
					if eating.At(i) == 0 {
						plates.Push(1)
						eating.SetAt(i, eating.At(i)-1)
					} else if eating.At(i) > 0 {
						eating.SetAt(i, eating.At(i)-1)
					} else if eating.At(i) < 0 {
						eating.EraseAt(i)
					}
				}
			}
			if plates.Size() > 5 {
				fmt.Println("Game Over")
				gameOver = true
				start = false
				rules = "GAME OVER\n\nSink was overfilled!!!\n\n Enter to try again\nEsc to exit"
				fmt.Fprintln(score, rules)
			}
			if customers.Size() > 8 {
				fmt.Println("Game Over")
				gameOver = true
				start = false
				rules = "GAME OVER\n\nToo many customers in line!!!\n\n Enter to try again\nEsc to exit"
				fmt.Fprintln(score, rules)
			}
		} else if gameOver {
			score.Draw(win, pixel.IM.Moved(pixel.V(50, 300)))
			if win.Pressed(pixelgl.KeyEnter) {
				start = true
				score.Clear()
				customers.Clear()
				plates.Clear()
				player.points = 0
			}
			if win.Pressed(pixelgl.KeyEscape) {
				break
			}
		} else {
			score.Draw(win, pixel.IM.Moved(pixel.V(50, 300)))
			if win.Pressed(pixelgl.KeyEnter) {
				start = true
				score.Clear()
			}
		}

		win.Update()
	}
}

func main() {
	pixelgl.Run(run)
}
