package lib

func GetLevels() []string {
	return []string{
		// Rebirth
		"Movement", "Pummel", "Gunner", "Cascade", "Elevate", "Bounce", "Purify", "Climb", "Fasttrack", "Glass Port",
		// Killer Inside
		"Take Flight", "Godspeed", "Dasher", "Thrasher", "Outstretched", "Smackdown", "Catwalk", "Fastlane", "Distinguish", "Dancer",
		// Only Shallow
		"Guardian", "Stomp", "Jumper", "Dash Tower", "Descent", "Driller", "Canals", "Sprint", "Mountain", "Superkinetic",
		// The Old City
		"Arrival", "Forgotten City", "The Clocktower",
		// The Burn That Cures
		"Fireball", "Ringer", "Cleaner", "Warehouse", "Boom", "Streets", "Steps", "Demolition", "Arcs", "Apartment",
		// Covenant
		"Hanging Gardens", "Tangled", "Waterworks", "Killswitch", "Falling", "Shocker", "Bouquet", "Prepare", "Triptrack", "Race",
		// Reckoning
		"Bubble", "Shield", "Overlook", "Pop", "Minefield", "Mimic", "Trigger", "Greenhouse", "Sweep", "Fuse",
		// Benediction
		"Heaven's Edge", "Zipline", "Swing", "Chute", "Crash", "Ascent", "Straightaway", "Firecracker", "Streak", "Mirror",
		// Apocrypha
		"Escalation", "Bolt", "Godstreak", "Plunge", "Mayhem", "Barrage", "Estate", "Trapwire", "Ricochet", "Fortress",
		// The Third Temple
		"Holy Ground", "The Third Temple",
		// Thousand Pound Butterfly
		"Spree", "Breakthrough", "Glide", "Closer", "Hike", "Switch", "Access", "Congregation", "Sequence", "Marathon",
		// Hand of God
		"Sacrifice", "Absolution",
		// Red Sidequests
		"Elevate Traversal I", "Elevate Traversal II", "Purify Traversal", "Godspeed Traversal", "Stomp Traversal", "Fireball Traversal", "Dominion Traversal", "Book of Life Traversal",
		// Yellow Sidequests
		"Sunset Flip Powerbomb", "Balloon Mountain", "Climbing Gym", "Fisherman Suplex", "STF", "Arena", "Attitude Adjustment", "Rocket",
		// Violet Sidequests
		"Doghouse", "Choker", "Chain", "Hellevator", "Razor", "All Seeing Eye", "Resident Saw I", "Resident Saw II",
	}
}

func GetXorKey() [16]byte {
	return [16]byte{0x82, 0xca, 0x81, 0xa9, 0x96, 0x86, 0x97, 0xc6,
		0xb9, 0xd9, 0xd2, 0xa7, 0x8d, 0xc9, 0x87, 0xa8}
}
