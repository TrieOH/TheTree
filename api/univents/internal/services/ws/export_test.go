package ws

import "github.com/google/uuid"

// Test-only accessors (export_test pattern): expose hub internals to the
// external test package without widening the real API.

// FrameShape is the testable view of one outbound frame.
type FrameShape struct {
	Type     string
	Terminal bool
}

// FramesForForTest maps a purchase status to the frames the hub would push,
// without needing a live socket. Exercises the dedupe/terminal logic.
func (o *Operations) FramesForForTest(status string, purchaseID, intentID uuid.UUID, hasIntent, statusChanged bool) []FrameShape {
	out := o.hub.framesFor(status, purchaseID, intentID, hasIntent, statusChanged)
	shapes := make([]FrameShape, 0, len(out))
	for _, f := range out {
		shapes = append(shapes, FrameShape{Type: f.typ, Terminal: f.terminal})
	}
	return shapes
}
