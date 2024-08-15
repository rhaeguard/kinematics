package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type Creature struct {
	name      string
	armatures []Armature
}

func NewCreature(name string) Creature {
	return Creature{
		name:      name,
		armatures: []Armature{},
	}
}

func (a *Creature) AddArmature(b Armature) {
	a.armatures = append(a.armatures, b)
}

type Npc struct {
	position rl.Vector2
	target   rl.Vector2
}

func randomVector2WithinBoundaries() rl.Vector2 {
	x := float32(rl.GetRandomValue(50, 1200))
	y := float32(rl.GetRandomValue(50, 680))
	return rl.NewVector2(x, y)
}

func (npc *Npc) update(deltaTime float32) {
	d := rl.Vector2Distance(npc.position, npc.target)
	if d <= 100 {
		npc.target = randomVector2WithinBoundaries()
	}
	npc.position = rl.Vector2Lerp(npc.position, npc.target, deltaTime*0.5)
}

func main() {
	rl.InitWindow(1280, 720, "inverse kinematics demo")
	defer rl.CloseWindow()

	rl.SetTargetFPS(30)

	limb := NewArmature("limb", rl.NewVector2(640, 10))

	for i := 0; i < 5; i++ {
		limb.AddBone(
			NewBone(rl.NewVector2(0, 0), 100),
		)
	}

	octo := NewCreature("octo")
	octo.AddArmature(limb)

	npcs := []Npc{
		{
			position: randomVector2WithinBoundaries(),
			target:   randomVector2WithinBoundaries(),
		},
		// {
		// 	position: randomVector2WithinBoundaries(),
		// 	target:   randomVector2WithinBoundaries(),
		// },
		// {
		// 	position: randomVector2WithinBoundaries(),
		// 	target:   randomVector2WithinBoundaries(),
		// },
	}

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)

		deltaTime := rl.GetFrameTime()

		for i := range npcs {
			npcs[i].update(deltaTime)
		}

		for i := range npcs {
			rl.DrawCircleV(npcs[i].position, 50, rl.Lime)
		}

		// target := rl.GetMousePosition()
		// target := npcs[0].position
		var target rl.Vector2

		intensity := 1.0

		candidate := npcs[0].position
		// candidate = rl.GetMousePosition()

		if rl.CheckCollisionPointCircle(candidate, limb.Last().end, 500) {
			if rl.CheckCollisionPointCircle(candidate, limb.Last().end, 100) {
				intensity = 10.0
			}
			target = candidate
		} else {
			target = limb.Last().end
		}

		limb.follow(rl.Vector2Lerp(limb.Last().end, target, float32(intensity)*rl.GetFrameTime()))
		limb.lock()
		limb.draw()

		rl.DrawCircleLines(
			int32(limb.Last().end.X),
			int32(limb.Last().end.Y),
			500,
			rl.Blue,
		)

		rl.DrawCircleLines(
			int32(limb.Last().end.X),
			int32(limb.Last().end.Y),
			100,
			rl.Red,
		)

		rl.EndDrawing()
	}
}
