package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	G                   float32 = 1000.0 // gravitational constant
	GRAVITY_SOFTENINING float32 = 5.0
	PHYSICS_DT          float32 = 1.0 / 120.0
	RADIUS_SCALE        float32 = 0.25
	SPAWN_RADIUS        float32 = 350.0
)

type particle struct {
	pos  rl.Vector2
	vel  rl.Vector2
	acc  rl.Vector2
	mass float32
}

func newParticle(pos, vel rl.Vector2, mass float32) particle {
	return particle{
		pos:  pos,
		vel:  vel,
		acc:  rl.Vector2Zero(),
		mass: mass,
	}
}

func (p particle) drawParticle() {
	rl.DrawCircle(int32(p.pos.X), int32(p.pos.Y), p.mass*RADIUS_SCALE, rl.RayWhite)
}

func createParticles(count uint32) []particle {
	const speed float32 = 0.3
	const ringCount uint32 = 333

	particles := make([]particle, 0, count)
	deltaDeg := 360.0 / float64(count/ringCount) * math.Pi / 180.0

	for i := range count {
		degree := (float64(i/ringCount) + float64(i%ringCount)/float64(ringCount)) * deltaDeg
		rad := SPAWN_RADIUS - SPAWN_RADIUS*float32(i%ringCount)/float32(ringCount)
		dirX := rad * float32(math.Cos(degree))
		dirY := rad * float32(math.Sin(degree))

		posX := dirX + float32(rl.GetScreenWidth())/2.0
		posY := dirY + float32(rl.GetScreenHeight())/2.0
		pos := rl.NewVector2(posX, posY)

		tangentDir := rl.NewVector2(dirY, -dirX)
		p := newParticle(pos, tangentDir.Scale(speed), 5)
		particles = append(particles, p)
		degree += deltaDeg
	}

	// posX := float32(rl.GetScreenWidth()) / 2.0
	// posY := float32(rl.GetScreenHeight()) / 2.0
	// pos := rl.NewVector2(posX, posY)
	// central_mass := newParticle(pos, rl.Vector2Zero(), 100)
	// particles = append(particles, central_mass)

	return particles
}

func calcAcceleration(particles []particle) {
	for i := range particles {
		particles[i].acc = rl.Vector2Zero()
	}

	for i := range particles {
		for j := i + 1; j < len(particles); j++ {
			p1 := &particles[i]
			p2 := &particles[j]

			dx := p2.pos.X - p1.pos.X
			dy := p2.pos.Y - p1.pos.Y
			distanceSqrd := dx*dx + dy*dy + (GRAVITY_SOFTENINING * GRAVITY_SOFTENINING)

			if distanceSqrd == 0 {
				continue
			}
			distance := math.Sqrt(float64(distanceSqrd))

			invDist := 1.0 / distance
			invDistCubed := float32(invDist * invDist * invDist)

			p1.acc.X += G * p2.mass * dx * invDistCubed
			p1.acc.Y += G * p2.mass * dy * invDistCubed

			p2.acc.X -= G * p1.mass * dx * invDistCubed
			p2.acc.Y -= G * p1.mass * dy * invDistCubed
		}
	}
}

// update pos and vel using velocity verlet integration
func update(particles []particle, dt float32) {
	for i := range particles {
		particles[i].vel = particles[i].vel.Add(particles[i].acc.Scale(dt * 0.5))
		particles[i].pos = particles[i].pos.Add(particles[i].vel.Scale(dt))
	}

	calcAcceleration(particles)

	for i := range particles {
		particles[i].vel = particles[i].vel.Add(particles[i].acc.Scale(dt * 0.5))
	}
}

func drawParticles(particles []particle) {
	for _, p := range particles {
		p.drawParticle()
	}
}

func main() {
	rl.InitWindow(int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), "Gravity Simulation - AoS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var count uint32 = 1000
	particles := createParticles(count)

	var accumulator float32 = 0.0

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		accumulator += dt

		for accumulator >= PHYSICS_DT {
			update(particles, PHYSICS_DT)

			accumulator -= PHYSICS_DT
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		drawParticles(particles)

		// rl.DrawCircleLines(int32(rl.GetScreenWidth()/2.0), int32(rl.GetScreenHeight()/2.0), SPAWN_RADIUS, rl.RayWhite)

		rl.DrawFPS(10, 10)

		rl.EndDrawing()
	}
}
