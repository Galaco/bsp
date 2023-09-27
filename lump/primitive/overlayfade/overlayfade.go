package overlayfade

// OverlayFade
type OverlayFade struct {
	// FadeDistMinSq is distance to start fading in
	FadeDistMinSq float32 `json:"fadeDistMinSq"`
	// FadeDistMaxSq is distance to fully fade in
	FadeDistMaxSq float32 `json:"fadeDistMaxSq"`
}
