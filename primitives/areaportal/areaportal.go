package areaportal

type AreaPortal struct {
	PortalKey           uint16 // Entities have a key called portal number (and in vbsp a variable called areaportalnum) to bind them to the area portals by comparing this value
	OtherArea           uint16 // The area this portal looks into
	FirstClipPortalVert uint16 // Portal geometry
	NumClipPortalVerts  uint16

	PlaneNum int32
}
