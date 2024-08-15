package main

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Armature struct {
	name   string
	bones  []Bone
	base   rl.Vector2
	locked bool
}

func NewArmature(name string, base rl.Vector2) Armature {
	return Armature{
		name:   name,
		bones:  []Bone{},
		base:   base,
		locked: true,
	}
}

func (a *Armature) AddBone(b Bone) {
	a.bones = append(a.bones, b)
}

func (a *Armature) Last() *Bone {
	return &a.bones[len(a.bones)-1]
}

func (a *Armature) follow(target rl.Vector2) {
	a.Last().follow(target)

	for i := len(a.bones) - 2; i >= 0; i-- {
		a.bones[i].follow(a.bones[i+1].start)
	}
}

func (a *Armature) draw() {
	for i := 0; i < len(a.bones); i++ {
		bone := a.bones[i]
		rl.DrawLineEx(bone.start, bone.end, 10, rl.White)
		rl.DrawCircleV(bone.start, 10, rl.Red)
		rl.DrawCircleV(bone.end, 10, rl.Yellow)
	}
}

func (a *Armature) lock() {
	if !a.locked {
		return
	}
	a.bones[0].start = a.base
	a.bones[0].calculateEnd()

	for i := 1; i < len(a.bones); i++ {
		a.bones[i].start = a.bones[i-1].end
		a.bones[i].calculateEnd()
	}
}

type Bone struct {
	start rl.Vector2
	end   rl.Vector2
	len   float64
	angle float32
}

func NewBone(start rl.Vector2, len float64) Bone {
	return Bone{
		start: start,
		len:   len,
	}
}

func (s *Bone) follow(target rl.Vector2) {
	dir := rl.Vector2Normalize(rl.Vector2Subtract(target, s.start))
	s.angle = rl.Vector2Angle(
		rl.NewVector2(1, 0), dir,
	)
	s.start = rl.Vector2Add(target, rl.Vector2Scale(dir, -float32(s.len)))
	s.calculateEnd()
}

func (s *Bone) calculateEnd() {
	dx := s.len * math.Cos(float64(s.angle))
	dy := s.len * math.Sin(float64(s.angle))

	s.end = rl.NewVector2(
		s.start.X+float32(dx), s.start.Y+float32(dy),
	)
}

func (s *Bone) translate(reference rl.Vector2) {
	s.start = rl.Vector2Subtract(s.start, reference)
	s.end = rl.Vector2Subtract(s.end, reference)
}

func findAngle(s1, s2 *Bone) float32 {
	v1 := rl.Vector2Normalize(rl.Vector2Subtract(s1.end, s1.start))
	v2 := rl.Vector2Normalize(rl.Vector2Subtract(s2.start, s2.end))

	dot := rl.Vector2DotProduct(v1, v2)
	l1 := rl.Vector2Length(v1)
	l2 := rl.Vector2Length(v2)

	degree := math.Acos(float64(dot/(l1*l2))) * rl.Rad2deg

	return float32(degree)
}
