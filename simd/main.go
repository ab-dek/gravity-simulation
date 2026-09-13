package main

import (
	"flag"
	"fmt"
	"math"
	"simd"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	G                   float32 = 1000.0 // gravitational constant
	GRAVITY_SOFTENINING float32 = 5.0
	PHYSICS_DT          float32 = 1.0 / 120.0
	RADIUS_SCALE        float32 = 0.25
	SPAWN_RADIUS        float32 = 350.0
)

var (
	countPtr *int     = flag.Int("count", 1000, "Particle count")
	ringPtr  *int     = flag.Int("ring", 333, "number of concentric circles")
	speedPtr *float64 = flag.Float64("speed", 0.3, "speed of particles")
)

type particles struct {
	posX, posY []float32
	velX, velY []float32
	accX, accY []float32
	mass       []float32
}

func allocParticles(count int) particles {
	return particles{
		posX: make([]float32, 0, count),
		posY: make([]float32, 0, count),
		velX: make([]float32, 0, count),
		velY: make([]float32, 0, count),
		accX: make([]float32, 0, count),
		accY: make([]float32, 0, count),
		mass: make([]float32, 0, count),
	}
}

func createParticles() particles {
	count := *countPtr
	speed := float32(*speedPtr)
	ringCount := *ringPtr

	ps := allocParticles(count)

	deltaDeg := 360.0 / float64(count/ringCount) * math.Pi / 180.0

	for i := range count {
		degree := (float64(i/ringCount) + float64(i%ringCount)/float64(ringCount)) * deltaDeg
		rad := SPAWN_RADIUS - SPAWN_RADIUS*float32(i%ringCount)/float32(ringCount)
		dirX := rad * float32(math.Cos(degree))
		dirY := rad * float32(math.Sin(degree))

		posX := dirX + float32(rl.GetScreenWidth())/2.0
		ps.posX = append(ps.posX, posX)
		posY := dirY + float32(rl.GetScreenHeight())/2.0
		ps.posY = append(ps.posY, posY)

		velX := dirY * speed
		ps.velX = append(ps.velX, velX)
		velY := -dirX * speed
		ps.velY = append(ps.velY, velY)

		mass := float32(5)
		ps.mass = append(ps.mass, mass)
		degree += deltaDeg
	}

	return ps
}

func update(ps particles, lanes int, dt float32) {

}

func main() {
	flag.Parse()
	rl.InitWindow(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), "Gravity Simulation - SoA")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	particles := createParticles()

	var accumulator float32 = 0.0
	var vec simd.Float32s
	lanes := vec.Len()

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		accumulator += dt

		for accumulator >= PHYSICS_DT {
			update(particles, lanes, PHYSICS_DT)

			accumulator -= PHYSICS_DT
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		// drawParticles(particles)

		// rl.DrawCircleLines(int32(rl.GetScreenWidth()/2.0), int32(rl.GetScreenHeight()/2.0), SPAWN_RADIUS, rl.RayWhite)

		rl.DrawText(fmt.Sprintf("lanes: %d", lanes), 10, 30, 20, rl.RayWhite)
		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
}
