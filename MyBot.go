package main

import (
	"hlt"
)

func findCountToFront(myID int, gameMap hlt.GameMap, loc hlt.Location, d hlt.Direction) int {
	var maxNum = gameMap.Height
	if (d == hlt.EAST || d == hlt.WEST) {
		maxNum = gameMap.Width
	}
	var current = loc;
	var site = gameMap.GetSite(current, hlt.STILL)
	for i := 0; i < maxNum; i++ {
		current = gameMap.GetLocation(current, d)
		site = gameMap.GetSite(current, hlt.STILL)
		if site.Owner != myID {
			return i
		}
	}
	return maxNum
}


func canNeighborCaptureWithOurHelp(myID int, gameMap hlt.GameMap, loc hlt.Location, d hlt.Direction) float64 {
	var ourSite = gameMap.GetSite(loc, hlt.STILL)
	var theirLocation = gameMap.GetLocation(loc, d)
	var theirSite = gameMap.GetSite(theirLocation, hlt.STILL)
	var value = 10000.0
	if theirSite.Owner != myID {
		return value
	}
	for _,d := range hlt.CARDINALS {
		var l = gameMap.GetLocation(theirLocation, d)
		var s = gameMap.GetSite(l, hlt.STILL)
		if s.Owner != myID {
			if ourSite.Strength + theirSite.Strength + theirSite.Production > s.Strength && theirSite.Strength + theirSite.Production < s.Strength {
				var v = float64(s.Strength) / float64(s.Production)
				if v < value {
					value = v
				}
			}
		}
	}
	return value
}


func move(myID int, gameMap hlt.GameMap, loc hlt.Location) hlt.Move {
	var site = gameMap.GetSite(loc, hlt.STILL)
	var allies = 0
	var value = 999999999.0
	var dir = hlt.STILL
	for _,d := range hlt.CARDINALS {
		var new_site = gameMap.GetSite(loc, d)
		if new_site.Owner != myID && new_site.Strength < site.Strength {
			var v = float64(new_site.Strength) / float64(new_site.Production)
			if v < value {
				value = v
				dir = d
			}
		}
		if new_site.Owner == myID {
			allies += 1
		}
	}

	if (dir != hlt.STILL) {
		return hlt.Move {
			Location: loc,
			Direction: dir,
		}
	}

	// fix for null times and 255 walls
	if allies < 4 && site.Strength == 255 {
		for _,d := range hlt.CARDINALS {
			if new_site.Owner != myID && new_site.Strength <= site.Strength {
				return hlt.Move {
					Location: loc,
					Direction: d,
				}
			}
		}
	}

	if site.Strength < site.Production * 5 {
		return hlt.Move {
			Location: loc,
			Direction: hlt.STILL,
		}
	}

	// see if we can help any of our allies capture
	if allies != 4 && (loc.X + loc.Y) % 2 != 0 {
		var best = 9999.0 
		var dir = hlt.STILL
		for _,d := range hlt.CARDINALS {
			var value = canNeighborCaptureWithOurHelp(myID, gameMap, loc, d)
			if value < best {
				best = value
				dir = d
			}
		}
		if best < 999.0 {
			return hlt.Move {
				Location: loc,
				Direction: dir,
			}
		}
	}

	// if we are surrounded by allies, move toward the nearest front
	if allies == 4 {
		var best = 100000
		var dir = hlt.STILL
		for _,d := range hlt.CARDINALS {
			var dist = findCountToFront(myID, gameMap, loc, d)
			if dist < best {
				best = dist
				dir = d
			}
		}
		return hlt.Move {
			Location: loc,
			Direction: dir,
		}
	}

	return hlt.Move {
		Location: loc,
		Direction: hlt.STILL,
	}
}

func main () {
	conn, gameMap := hlt.NewConnection("bovard")
	for {
		var moves hlt.MoveSet
		gameMap = conn.GetFrame()
		for y := 0; y < gameMap.Height; y++ {
			for x := 0; x < gameMap.Width; x++ {
				loc := hlt.NewLocation(x,y)
				if gameMap.GetSite(loc, hlt.STILL).Owner == conn.PlayerTag {
					moves = append(moves, move(conn.PlayerTag, gameMap, loc))
				}
			}
		}
		conn.SendFrame(moves)

	}
}
