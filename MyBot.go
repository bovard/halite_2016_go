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


func move(myID int, gameMap hlt.GameMap, loc hlt.Location) hlt.Move {
	var site = gameMap.GetSite(loc, hlt.STILL)
	var allies = 0
	for _,d := range hlt.CARDINALS {
		var new_site = gameMap.GetSite(loc, d)
		if new_site.Owner != myID && new_site.Strength < site.Strength {
			return hlt.Move {
				Location: loc,
				Direction: d,
			}
		}
		if new_site.Owner == myID {
			allies += 1
		}
	}

	if site.Strength < site.Production * 5 {
		return hlt.Move {
			Location: loc,
			Direction: hlt.STILL,
		}
	}

	if allies == 4 {
		var north = findCountToFront(myID, gameMap, loc, hlt.NORTH)
		var east = findCountToFront(myID, gameMap, loc, hlt.EAST)
		var south = findCountToFront(myID, gameMap, loc, hlt.SOUTH)
		var west = findCountToFront(myID, gameMap, loc, hlt.WEST)

		if north <= east && north <= west && north <= east {
			return hlt.Move {
				Location: loc,
				Direction: hlt.NORTH,
			}
		}

		if east <= west && east <= south {
			return hlt.Move {
				Location: loc,
				Direction: hlt.EAST,
			}
		}

		if south <= west {
			return hlt.Move {
				Location: loc,
				Direction: hlt.SOUTH,
			}
		}

		return hlt.Move {
			Location: loc,
			Direction: hlt.WEST,
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
