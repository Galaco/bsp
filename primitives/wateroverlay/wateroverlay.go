package wateroverlay

/*

#define WATEROVERLAY_BSP_FACE_COUNT				256
#define WATEROVERLAY_RENDER_ORDER_NUM_BITS		2
#define WATEROVERLAY_NUM_RENDER_ORDERS			(1<<WATEROVERLAY_RENDER_ORDER_NUM_BITS)
#define WATEROVERLAY_RENDER_ORDER_MASK			0xC000	// top 2 bits set
struct dwateroverlay_t
{
	DECLARE_BYTESWAP_DATADESC();
	int				nId;
	short			nTexInfo;

	// Accessors..
	void			SetFaceCount( unsigned short count );
	unsigned short	GetFaceCount() const;
	void			SetRenderOrder( unsigned short order );
	unsigned short	GetRenderOrder() const;

private:

	unsigned short	m_nFaceCountAndRenderOrder;

public:

	int				aFaces[WATEROVERLAY_BSP_FACE_COUNT];
	float			flU[2];
	float			flV[2];
	Vector			vecUVPoints[4];
	Vector			vecOrigin;
	Vector			vecBasisNormal;
};
*/

// WaterOverlay
type WaterOverlay struct {
}
