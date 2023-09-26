package bsp

const (
	ContentsEmpty = 0 // No contents

	ContentsSolid       = 0x1 // an eye is never valid in a solid
	ContentsWindow      = 0x2 // translucent, but not watery (glass)
	ContentsAux         = 0x4
	ContentsGrate       = 0x8 // alpha-tested "grate" textures.  Bullets/sight pass through, but solids don't
	ContentsSlime       = 0x10
	ContentsWater       = 0x20
	ContentsBlockLOS    = 0x40 // block AI line of sight
	ContentsOpaque      = 0x80 // things that cannot be seen through (may be non-solid though)
	LastVisibleContents = 0x80

	AllVisibleContents = LastVisibleContents | (LastVisibleContents - 1)

	ContentsTestFogVolume = 0x100
	ContentsUnused        = 0x200

	// ContentsUnused6 is unused.
	// NOTE: If it's visible, grab from the top + update LastVisibleContents
	// if not visible, then grab from the bottom.
	ContentsUnused6 = 0x400

	ContentsTeam1 = 0x800  // per team contents used to differentiate collisions
	ContentsTeam2 = 0x1000 // between players and objects on different teams

	// ContentsIgnoreNodrawOpaque ignores ContentsOpaque on surfaces that have SurfNodraw.
	ContentsIgnoreNodrawOpaque = 0x2000

	// ContentsMoveable hits entities which are MOVETYPE_PUSH (doors, plats, etc.)
	ContentsMoveable = 0x4000

	// Remaining contents are non-visible, and don't eat brushes.

	ContentsAreaPortal = 0x8000

	ContentsPlayerClip  = 0x10000
	ContentsMonsterClip = 0x20000

	// Currents can be added to any other contents, and may be mixed.

	ContentsCurrent0    = 0x40000
	ContentsCurrent90   = 0x80000
	ContentsCurrent180  = 0x100000
	ContentsCurrent270  = 0x200000
	ContentsCurrentUp   = 0x400000
	ContentsCurrentDown = 0x800000

	ContentsOrigin = 0x1000000 // removed before bsping an entity

	ContentsMonster     = 0x2000000 // should never be on a brush, only in game
	ContentsDebris      = 0x4000000
	ContentsDetail      = 0x8000000  // brushes to be added after vis leafs
	ContentsTranslucent = 0x10000000 // auto set if any surface has trans
	ContentsLadder      = 0x20000000
	ContentsHitbox      = 0x40000000 // use accurate hitboxes on trace

	// NOTE: These are stored in a short in the engine now.  Don't use more than 16 bits.

	SurfLight    = 0x0001 // value will hold the light strength
	SurfSky2D    = 0x0002 // don't draw, indicates we should skylight + draw 2d sky but not draw the 3D skybox
	SurfSky      = 0x0004 // don't draw, but add to skybox
	SurfWarp     = 0x0008 // turbulent water warp
	SurfTrans    = 0x0010
	SurfNoPortal = 0x0020 // the surface can not have a portal placed on it
	SurfTrigger  = 0x0040 // FIXME: This is an xbox hack to work around elimination of trigger surfaces, which breaks occluders
	SurfNodraw   = 0x0080 // don't bother referencing the texture

	SurfHint = 0x0100 // make a primary bsp splitter

	SurfSkip      = 0x0200 // completely ignore, allowing non-closed brushes
	SurfNoLight   = 0x0400 // Don't calculate light
	SurfBumpLight = 0x0800 // calculate three lightmaps for the surface for bumpmapping
	SurfNoShadows = 0x1000 // Don't receive shadows
	SurfNoDecals  = 0x2000 // Don't receive decals
	SurfNoChop    = 0x4000 // Don't subdivide patches on this surface
	SurfHitbox    = 0x8000 // surface is part of a hitbox

	// Spacial content masks - used for spacial queries (traceline,etc.).

	// MaskAll is everything.
	MaskAll = 0xFFFFFFFF

	// MaskSolid is everything that is normally solid.
	MaskSolid = ContentsSolid | ContentsMoveable | ContentsWindow | ContentsMonster | ContentsGrate

	// MaskPlayerSolid is everything that blocks player movement.
	MaskPlayerSolid = ContentsSolid | ContentsMoveable | ContentsPlayerClip | ContentsWindow | ContentsMonster | ContentsGrate

	// MaskNPCSolid blocks npc movement.
	MaskNPCSolid = ContentsSolid | ContentsMoveable | ContentsMonsterClip | ContentsWindow | ContentsMonster | ContentsGrate

	// MaskWater is water physics in these contents.
	MaskWater = ContentsWater | ContentsMoveable | ContentsSlime

	// MaskOpaque is everything that blocks lighting.
	MaskOpaque = ContentsSolid | ContentsMoveable | ContentsOpaque

	// MaskOpaqueAndNPCs is everything that blocks lighting, but with monsters added.
	MaskOpaqueAndNPCs = MaskOpaque | ContentsMonster

	// MaskBlockLOS is everything that blocks line of sight for AI.
	MaskBlockLOS = ContentsSolid | ContentsMoveable | ContentsBlockLOS

	// MaskBlockLOSAndNPCs is everything that blocks line of sight for AI plus NPCs.
	MaskBlockLOSAndNPCs = MaskBlockLOS | ContentsMonster

	// MaskVisible is everything that blocks line of sight for players.
	MaskVisible = MaskOpaque | ContentsIgnoreNodrawOpaque

	// MaskVisibleAndNPCs is everything that blocks line of sight for players, but with monsters added.
	MaskVisibleAndNPCs = MaskOpaqueAndNPCs | ContentsIgnoreNodrawOpaque

	// MaskShot is bullets see these as solid.
	MaskShot = ContentsSolid | ContentsMoveable | ContentsMonster | ContentsWindow | ContentsDebris | ContentsHitbox

	// MaskShotHull is non-raycasted weapons see this as solid (includes grates).
	MaskShotHull = ContentsSolid | ContentsMoveable | ContentsMonster | ContentsWindow | ContentsDetail | ContentsGrate

	// MaskShotPortal hits solids (not grates) and passes through everything else.
	MaskShotPortal = ContentsSolid | ContentsMoveable | ContentsWindow | ContentsMonster

	// MaskSolidBrushOnly is everything normally solid, except monsters (world+brush only).
	MaskSolidBrushOnly = ContentsSolid | ContentsMoveable | ContentsWindow | ContentsGrate

	// MaskPlayerSolidBrushOnly is everything normally solid for player movement, except monsters (world+brush only).
	MaskPlayerSolidBrushOnly = ContentsSolid | ContentsMoveable | ContentsWindow | ContentsPlayerClip | ContentsGrate

	// MaskNPCSolidBrushOnly is everything normally solid for npc movement, except monsters (world+brush only).
	MaskNPCSolidBrushOnly = ContentsSolid | ContentsMoveable | ContentsWindow | ContentsMonsterClip | ContentsGrate

	// MaskNPCWorldStatic is just the world, used for route rebuilding.
	MaskNPCWorldStatic = ContentsSolid | ContentsWindow | ContentsMonsterClip | ContentsGrate

	// MaskSplitAreaPortal are things that can split areaportals.
	MaskSplitAreaPortal = ContentsWater | ContentsSlime

	// MaskCurrent is moving water.
	// UNDONE: This is untested.
	MaskCurrent = ContentsCurrent0 | ContentsCurrent90 | ContentsCurrent180 | ContentsCurrent270 | ContentsCurrentUp | ContentsCurrentDown

	// MaskDeadSolid is everything that blocks corpse movement.
	// UNDONE: Not used yet / may be deleted.
	MaskDeadSolid = ContentsSolid | ContentsPlayerClip | ContentsWindow | ContentsGrate
)
