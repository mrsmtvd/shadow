package annotations

type Storage interface {
	Create(Annotation) error
}
