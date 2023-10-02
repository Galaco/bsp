package areaportal

// AreaPortal
type AreaPortal struct {
	// PortalKey - entities have a key called portal number (and in vbsp a variable called areaportalnum) to bind
	// them to the area portals by comparing this value
	PortalKey uint16 `json:"portalKey"`
	// OtherArea is the area this portal looks into
	OtherArea uint16 `json:"otherArea"`
	// FirstClipPortalVert Portal geometry
	FirstClipPortalVert uint16 `json:"firstClipPortalVert"`
	// NumClipPortalVerts number of geometry verts
	NumClipPortalVerts uint16 `json:"numClipPortalVerts"`
	// PlaneNum
	PlaneNum int32 `json:"planeNum"`
}
