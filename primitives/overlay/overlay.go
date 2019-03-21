package overlay

/*
struct doverlay_t
{
	DECLARE_BYTESWAP_DATADESC();
	int			nId;
	short		nTexInfo;

	// Accessors..
	void			SetFaceCount( unsigned short count );
	unsigned short	GetFaceCount() const;

	void			SetRenderOrder( unsigned short order );
	unsigned short	GetRenderOrder() const;

private:
	unsigned short	m_nFaceCountAndRenderOrder;

public:
	int			aFaces[OVERLAY_BSP_FACE_COUNT];
	float		flU[2];
	float		flV[2];
	Vector		vecUVPoints[4];
	Vector		vecOrigin;
	Vector		vecBasisNormal;
};
*/
// Overlay
type Overlay struct {
	// NId
	NId int32
	// NTexInfo
	NTexInfo int16
	// AlignmentPadding
	AlignmentPadding int16
}
