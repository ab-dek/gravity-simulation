package test

import (
	"math"
	"math/rand/v2"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	PARTICLE_COUNT         = 1000
	GRAVITATIONAL_CONSTANT = 15.0
	GRAVITY_SOFTENING      = 3.0
	MIN_PARTICLE_MASS      = 5.0
	MAX_PARTICLE_MASS      = 12.0
	PARTICLE_RADIUS_SCALE  = 0.45
	SPAWN_RADIUS           = 500.0
	PHYSICS_DT             = 1.0 / 120.0
	MAX_PHYSICS_STEPS      = 8
)

type particle struct {
	pos        rl.Vector2
	vel        rl.Vector2
	acc        rl.Vector2
	mass       float32
	brightness float32
}

func particleRadius(mass float32) float32 {
	return PARTICLE_RADIUS_SCALE * float32(math.Sqrt(float64(mass)))
}

func spawnParticle(clockwise bool) particle {
	angle := rand.Float32() * 2.0 * math.Pi

	// Area-uniform random radius with center bias
	rFraction := rand.Float32()
	uniformR := float32(math.Sqrt(float64(rFraction)))
	centerBiasedR := float32(math.Pow(float64(rFraction), 1.40))
	normRadius := uniformR + (centerBiasedR-uniformR)*0.75
	dist := normRadius * SPAWN_RADIUS

	radDir := rl.NewVector2(float32(math.Cos(float64(angle))), float32(math.Sin(float64(angle))))
	tangentDir := rl.NewVector2(-radDir.Y, radDir.X)
	if !clockwise {
		tangentDir = rl.NewVector2(radDir.Y, -radDir.X)
	}

	speedFactor := float32(math.Pow(float64(normRadius), 0.50))
	tangentSpeed := 6.0 + (20.0-6.0)*speedFactor

	return particle{
		pos:        radDir.Scale(dist),
		vel:        tangentDir.Scale(tangentSpeed),
		acc:        rl.Vector2Zero(),
		mass:       MIN_PARTICLE_MASS + rand.Float32()*(MAX_PARTICLE_MASS-MIN_PARTICLE_MASS),
		brightness: 0.55 + rand.Float32()*0.45,
	}
}

func createParticles() []particle {
	particles := make([]particle, PARTICLE_COUNT)
	for i := range particles {
		particles[i] = spawnParticle(i < PARTICLE_COUNT)
	}
	return particles
}

func calcAccelerations(particles []particle) {
	for i := range particles {
		particles[i].acc = rl.Vector2Zero()
	}

	softeningSq := float32(GRAVITY_SOFTENING * GRAVITY_SOFTENING)

	for i := 0; i < len(particles); i++ {
		for j := i + 1; j < len(particles); j++ {
			p1 := &particles[i]
			p2 := &particles[j]

			dx := p2.pos.X - p1.pos.X
			dy := p2.pos.Y - p1.pos.Y
			distSq := dx*dx + dy*dy + softeningSq

			invDist := 1.0 / float32(math.Sqrt(float64(distSq)))
			invDistCubed := invDist * invDist * invDist
			factor := GRAVITATIONAL_CONSTANT * invDistCubed

			p1.acc.X += dx * factor * p2.mass
			p1.acc.Y += dy * factor * p2.mass

			p2.acc.X -= dx * factor * p1.mass
			p2.acc.Y -= dy * factor * p1.mass
		}
	}
}

func mergeCollisions(particles *[]particle) {
	p := *particles
	i := 0
	for i < len(p) {
		merged := false
		r1 := particleRadius(p[i].mass)

		for j := i + 1; j < len(p); j++ {
			r2 := particleRadius(p[j].mass)
			colDist := r1 + r2
			dx := p[j].pos.X - p[i].pos.X
			dy := p[j].pos.Y - p[i].pos.Y

			if dx*dx+dy*dy <= colDist*colDist {
				totalMass := p[i].mass + p[j].mass

				// Position, Velocity & Brightness weighted by mass (inelastic merge)
				p[i].pos = p[i].pos.Scale(p[i].mass).Add(p[j].pos.Scale(p[j].mass)).Scale(1.0 / totalMass)
				p[i].vel = p[i].vel.Scale(p[i].mass).Add(p[j].vel.Scale(p[j].mass)).Scale(1.0 / totalMass)
				p[i].brightness = (p[i].brightness*p[i].mass + p[j].brightness*p[j].mass) / totalMass
				p[i].mass = totalMass

				// Remove second particle via swap
				p[j] = p[len(p)-1]
				p = p[:len(p)-1]
				merged = true
				break
			}
		}
		if !merged {
			i++
		}
	}
	*particles = p
}

func updateSystem(particles *[]particle) {
	p := *particles
	if len(p) == 0 {
		return
	}

	// Step 1: Initial Accelerations & Half Velocity update + Position step
	calcAccelerations(p)
	for i := range p {
		p[i].vel = p[i].vel.Add(p[i].acc.Scale(PHYSICS_DT * 0.5))
		p[i].pos = p[i].pos.Add(p[i].vel.Scale(PHYSICS_DT))
	}

	// Step 2: Inelastic Merging
	mergeCollisions(particles)
	p = *particles
	if len(p) == 0 {
		return
	}

	// Step 3: Recalculate Accelerations at new positions & Finish Velocity Step
	calcAccelerations(p)
	for i := range p {
		p[i].vel = p[i].vel.Add(p[i].acc.Scale(PHYSICS_DT * 0.5))
	}
}

func main() {
	rl.InitWindow(1280, 800, "Gravity Simulation - Velocity Verlet")
	defer rl.CloseWindow()
	rl.SetTargetFPS(60)

	particles := createParticles()
	var accumulator float32 = 0.0

	for !rl.WindowShouldClose() {
		dt := rl.GetFrameTime()
		if dt > 0.1 {
			dt = 0.1
		}
		accumulator += dt

		steps := 0
		for accumulator >= PHYSICS_DT && steps < MAX_PHYSICS_STEPS {
			updateSystem(&particles)
			accumulator -= PHYSICS_DT
			steps++
		}
		if steps == MAX_PHYSICS_STEPS {
			accumulator = 0.0
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.Color{R: 1, G: 2, B: 6, A: 255})

		center := rl.NewVector2(float32(rl.GetScreenWidth())*0.5, float32(rl.GetScreenHeight())*0.5)
		zoom := float32(rl.GetScreenHeight()) / (SPAWN_RADIUS * 2.25)

		for _, p := range particles {
			screenPos := center.Add(p.pos.Scale(zoom))
			radius := particleRadius(p.mass) * zoom
			if radius < 0.7 {
				radius = 0.7
			}
			rl.DrawCircleV(screenPos, radius, rl.RayWhite)
		}

		rl.DrawFPS(10, 10)
		rl.EndDrawing()
	}
}
