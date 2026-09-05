package main

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const SCREEN_WIDTH int32 = 800
const SCREEN_HEIGHT int32 = 450

const GRAVITATIONAL_CONSTANT float32 = 15
const STEP float32 = 8.0
const PHYSICS_DT float32 = 1.0 / (60.0 * STEP)

type particle struct {
	pos  rl.Vector2
	vel  rl.Vector2
	acc  rl.Vector2
	mass float32
}

func new_particle(pos rl.Vector2, mass float32) particle {
	return particle{
		pos:  pos,
		vel:  rl.Vector2Zero(),
		acc:  rl.Vector2Zero(),
		mass: mass,
	}
}

// update pos and vel using velocity verlet integration
func (p *particle) update_pos(dt float32) {
	p.pos = p.pos.Add(p.vel.Scale(dt)).Add(p.acc.Scale(dt * dt * 0.5))
}

func (p *particle) update_vel(dt float32) {
	p.vel = p.vel.Add(p.acc.Scale(dt * 0.5))
}

func (p particle) draw_particle() {
	rl.DrawCircle(int32(p.pos.X), int32(p.pos.Y), p.mass, rl.RayWhite)
}

func new_pos_rand() rl.Vector2 {
	pos_x := rand.Float32() * float32(SCREEN_WIDTH)
	pos_y := rand.Float32() * float32(SCREEN_HEIGHT)
	return rl.NewVector2(pos_x, pos_y)
}

func create_particles(count uint32) []particle {
	particles := make([]particle, 0, count)

	for range count {
		p := new_particle(new_pos_rand(), rand.Float32()*10+5)
		particles = append(particles, p)
	}

	return particles
}

func calc_acceleration(particles []particle) {
	for i := range particles {
		particles[i].acc = rl.Vector2Zero()
	}

	for i := 0; i < len(particles); i++ {
		for j := i + 1; j < len(particles); j++ {
			p1 := &particles[i]
			p2 := &particles[j]

			dx := p2.pos.X - p1.pos.X
			dy := p2.pos.Y - p1.pos.Y
			distance_sqrd := dx*dx + dy*dy

			if distance_sqrd == 0 {
				continue
			}
			distance := math.Sqrt(float64(distance_sqrd))

			inv_distance := 1.0 / distance
			inv_distance_cubed := float32(inv_distance * inv_distance * inv_distance)

			p1.acc.X += GRAVITATIONAL_CONSTANT * p2.mass * dx * inv_distance_cubed
			p1.acc.Y += GRAVITATIONAL_CONSTANT * p2.mass * dy * inv_distance_cubed

			p2.acc.X -= GRAVITATIONAL_CONSTANT * p1.mass * dx * inv_distance_cubed
			p2.acc.Y -= GRAVITATIONAL_CONSTANT * p1.mass * dy * inv_distance_cubed
		}
	}
}

func update(particles []particle, dt float32) {
	for i := range particles {
		particles[i].update_pos(dt)
	}

	calc_acceleration(particles)

	for i := range particles {
		particles[i].update_vel(dt)
	}
}

func draw_particles(particles []particle) {
	for _, p := range particles {
		p.draw_particle()
	}
}

func main() {
	rl.InitWindow(SCREEN_WIDTH, SCREEN_HEIGHT, "Gravity Simulation - AoS")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	var count uint32 = 10
	particles := create_particles(count)

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

		draw_particles(particles)

		rl.EndDrawing()
	}
}
