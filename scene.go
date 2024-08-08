package robb

type Scene struct {
	Camera  Camera
	objects []Drawable
}

func NewScene(camera Camera) *Scene {
	return &Scene{
		Camera: camera,
	}
}

func (s *Scene) Draw() {
	view := s.Camera.ViewMatrix()
	projection := s.Camera.ProjectionMatrix()

	for _, obj := range s.objects {
		obj.Draw(view, projection)
	}
}

func (s *Scene) AddObject(obj ...Drawable) {
	s.objects = append(s.objects, obj...)
}
