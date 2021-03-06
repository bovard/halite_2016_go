package main

import (
	"hlt"
)

const NEUTRAL = 0

func findCountToFront(myID int, gameMap hlt.GameMap, loc hlt.Location, d hlt.Direction) int {
	var maxNum = gameMap.Height
	if d == hlt.EAST || d == hlt.WEST {
		maxNum = gameMap.Width
	}
	var current = loc
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
	var value = 100000.0
	if theirSite.Owner != myID {
		return value
	}
	for _, d := range hlt.CARDINALS {
		l := gameMap.GetLocation(theirLocation, d)
		s := gameMap.GetSite(l, hlt.STILL)
		if s.Owner != myID {
			if ourSite.Strength+theirSite.Strength+theirSite.Production > s.Strength && theirSite.Strength+theirSite.Production < s.Strength {
				v := siteValue(s)
				if v < value {
					value = v
				}
			}
		}
	}
	return value
}

func siteValue(site hlt.Site) float64 {
	if site.Production == 0 {
		return 1000.0
	} else {
		return float64(site.Strength) / float64(site.Production)
	}
}

func nearEnemey(myID int, gameMap hlt.GameMap, loc hlt.Location) bool {
	var sites []hlt.Site
	for _, d := range hlt.CARDINALS {
		l := gameMap.GetLocation(loc, d)
		sites = append(sites, gameMap.GetSite(loc, d))
		sites = append(sites, gameMap.GetSite(l, d))
		if d == hlt.NORTH || d == hlt.SOUTH {
			l := gameMap.GetLocation(loc, d)
			sites = append(sites, gameMap.GetSite(l, hlt.RotateDirectionLeft(d)))
			sites = append(sites, gameMap.GetSite(l, hlt.RotateDirectionRight(d)))
		}
	}
	for _, s := range sites {
		if s.Owner != NEUTRAL && s.Owner != myID {
			return true
		}
	}
	return false
}

func damageToEnemyFromMove(myID int, gameMap hlt.GameMap, loc hlt.Location, dir hlt.Direction) (int, int, int) {
	l := gameMap.GetLocation(loc, dir)
	straight := gameMap.GetSite(l, dir)
	left := gameMap.GetSite(l, hlt.RotateDirectionLeft(dir))
	right := gameMap.GetSite(l, hlt.RotateDirectionRight(dir))
	lDamage := 0
	sDamage := 0
	rDamage := 0
	if straight.Owner != NEUTRAL && straight.Owner != myID {
		sDamage = straight.Strength
	}
	if left.Owner != NEUTRAL && left.Owner != myID {
		lDamage = left.Strength
	}
	if right.Owner != NEUTRAL && right.Owner != myID {
		rDamage = right.Strength
	}
	return sDamage, lDamage, rDamage
}

func move(myID int, gameMap hlt.GameMap, loc hlt.Location) hlt.Move {
	site := gameMap.GetSite(loc, hlt.STILL)
	allies := 0
	enemies := 0
	value := 999999999.0
	dir := hlt.STILL
	for _, d := range hlt.CARDINALS {
		s := gameMap.GetSite(loc, d)
		if s.Owner != myID && s.Strength < site.Strength {
			v := siteValue(s)
			if v < value {
				value = v
				dir = d
			}
		}
		if s.Owner == myID {
			allies += 1
		} else if s.Owner != NEUTRAL {
			enemies += 1
		}
	}

	if nearEnemey(myID, gameMap, loc) {
		maxDmg := 0
		dir := hlt.STILL
		for _, d := range hlt.CARDINALS {
			s, l, r := damageToEnemyFromMove(myID, gameMap, loc, d)
			dmg := s + l + r
			if dmg > maxDmg {
				maxDmg = dmg
				dir = d
			}
		}

		if maxDmg > 0 {
			return hlt.Move{
				Location:  loc,
				Direction: dir,
			}
		}

	}

	if dir != hlt.STILL {
		return hlt.Move{
			Location:  loc,
			Direction: dir,
		}
	}

	// fix for null times and 255 walls
	if allies < 4 && site.Strength == 255 {
		for _, d := range hlt.CARDINALS {
			s := gameMap.GetSite(loc, d)
			if s.Owner != myID && s.Strength == site.Strength {
				return hlt.Move{
					Location:  loc,
					Direction: d,
				}
			}
		}
	}

	distToClosestFront := 100000
	dirToClosestFont := hlt.STILL
	for _, d := range hlt.CARDINALS {
		var dist = findCountToFront(myID, gameMap, loc, d)
		if dist < distToClosestFront {
			distToClosestFront = dist
			dirToClosestFont = d
		}
	}

	// see if we can help any of our allies capture by moving to their square
	if distToClosestFront <= 1 {
		theirBest := 9999.0
		toThem := hlt.STILL
		for _, d := range hlt.CARDINALS {
			value := canNeighborCaptureWithOurHelp(myID, gameMap, loc, d)
			if value < theirBest {
				theirBest = value
				toThem = d
			}
		}

		ourBest := 9999.0
		for _, d := range hlt.CARDINALS {
			site := gameMap.GetSite(loc, d)
			if site.Owner != myID {
				v := siteValue(site)
				if v < ourBest {
					ourBest = v
				}
			}
		}
		// if they have a better spot than we do, move there!
		if theirBest < 9999.0 && theirBest < ourBest {
			return hlt.Move{
				Location:  loc,
				Direction: toThem,
			}
		}
	}

	// see if anyone else is waiting to take the same square as we are
	if allies == 3 {
		theirLoc := loc
		strength := 0
		toThem := hlt.STILL
		for _, d := range hlt.CARDINALS {
			site := gameMap.GetSite(loc, d)
			if site.Owner != myID {
				theirLoc = gameMap.GetLocation(loc, d)
				strength = site.Strength
				toThem = d
			}
		}

		ourStr := 0
		for _, d := range hlt.CARDINALS {
			site := gameMap.GetSite(theirLoc, d)
			if site.Owner == myID {
				ourStr += site.Strength
			}
		}
		if ourStr > strength {
			return hlt.Move{
				Location:  loc,
				Direction: toThem,
			}
		}

	}

	if site.Strength < site.Production*5 {
		return hlt.Move{
			Location:  loc,
			Direction: hlt.STILL,
		}
	}

	// if we are surrounded by allies, move toward the nearest front
	if distToClosestFront >= 2 {
		return hlt.Move{
			Location:  loc,
			Direction: dirToClosestFont,
		}
	}

	return hlt.Move{
		Location:  loc,
		Direction: hlt.STILL,
	}
}

func locIsAdjacentToOwner(gameMap hlt.GameMap, loc hlt.Location, owner int) bool {
	return gameMap.GetSite(loc, hlt.NORTH).Owner == owner || gameMap.GetSite(loc, hlt.SOUTH).Owner == owner || gameMap.GetSite(loc, hlt.EAST).Owner == owner || gameMap.GetSite(loc, hlt.WEST).Owner == owner
}

func moveOrReserveToCaptureLoc(gameMap hlt.GameMap, loc hlt.Location, myID int, allies []hlt.Location) ([]hlt.Location, hlt.MoveSet) {
	target := gameMap.GetSite(loc, hlt.STILL)
	var distances [10000]int
	var strengths [10]int
	var productions [10]int
	for i, al := range allies {
		dist := gameMap.GetManDistance(al, loc)
		distances[i] = dist
		if dist < 10 {
			site := gameMap.GetSite(al, hlt.STILL)
			strengths[dist] += site.Strength
			productions[dist] += site.Production
		}
	}
	needed := 10
	for i := 1; i < 10; i++ {
		total := 0
		for j := 1; j <= i; j++ {
			total += (strengths[j] + (i-j)*productions[j])
		}
		if total > target.Strength {
			needed = i
			break
		}
	}
	//gameMap.LogMessage("we need ")
	//gameMap.LogMessage(strconv.Itoa(needed))
	var toRemove []hlt.Location
	var moves hlt.MoveSet
	for i := 1; i <= needed; i++ {
		if i < needed {
			for j := 0; j < len(allies); j++ {
				if distances[j] == i {
					moves = append(moves, hlt.Move{
						Location:  allies[j],
						Direction: hlt.STILL,
					})
					toRemove = append(toRemove, allies[j])
				}
			}
		} else {
			for j := 0; j < len(allies); j++ {
				if distances[j] == i {
					if true {
						dir, _ := gameMap.GetDirectionTo(allies[j], loc)
						moves = append(moves, hlt.Move{
							Location:  allies[j],
							Direction: dir,
						})
						toRemove = append(toRemove, allies[j])
					} else if i == 1 {
						dir, _ := gameMap.GetDirectionTo(allies[j], loc)
						moves = append(moves, hlt.Move{
							Location:  allies[j],
							Direction: dir,
						})
						toRemove = append(toRemove, allies[j])
					} else {
						dist := 1000
						var dir hlt.Direction
						for _, d := range hlt.CARDINALS {
							dis := gameMap.GetManDistance(loc, allies[j])
							if dis < dist && gameMap.GetSite(loc, d).Owner == myID {
								dist = dis
								dir = d
							}
						}
						if dist < 1000 {
							moves = append(moves, hlt.Move{
								Location:  allies[j],
								Direction: dir,
							})
							toRemove = append(toRemove, allies[j])
						}

					}

				}
			}
		}
	}
	var remaining []hlt.Location
	for _, loc := range allies {
		remove := false
		for _, rem := range toRemove {
			if loc.X == rem.X && loc.Y == rem.Y {
				remove = true
				break
			}
		}
		if !remove {
			remaining = append(remaining, loc)
		}
	}
	return remaining, moves
}

func main() {
	conn, gameMap := hlt.NewConnection("bovard")
	//for turn := 0; turn < 5; turn++ {
	for {
		//gameMap.LogMessage("NEW TURN")
		//gameMap.LogMessage(strconv.Itoa(turn))
		var moves hlt.MoveSet
		strength := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		production := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		territory := [10]int{0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		var allies []hlt.Location
		var adjacent []hlt.Location
		gameMap = conn.GetFrame()
		for y := 0; y < gameMap.Height; y++ {
			for x := 0; x < gameMap.Width; x++ {
				loc := hlt.NewLocation(x, y)
				site := gameMap.GetSite(loc, hlt.STILL)
				strength[site.Owner] += site.Strength
				production[site.Owner] += site.Production
				territory[site.Owner] += 1
				if site.Owner == conn.PlayerTag {
					allies = append(allies, loc)
				} else if locIsAdjacentToOwner(gameMap, loc, conn.PlayerTag) {
					adjacent = append(adjacent, loc)
				}

			}
		}

		for _, loc := range allies {
			moves = append(moves, move(conn.PlayerTag, gameMap, loc))
		}

		var smarterMoves hlt.MoveSet;
		var result [50][50]int
		for _, m :=  range moves {
			if m.Direction == hlt.STILL {
				result[m.Location.Y][m.Location.X] += gameMap.GetSite(m.Location, hlt.STILL).Strength
				smarterMoves = append(smarterMoves, m)
			}
		}
		for _, m := range moves {
			if m.Direction != hlt.STILL {
				str := gameMap.GetSite(m.Location, hlt.STILL).Strength
				dest := gameMap.GetLocation(m.Location, m.Direction)
				if result[dest.Y][dest.X] + str > 300 {
					result[m.Location.Y][m.Location.X] += str
					smarterMoves = append(smarterMoves, hlt.Move{m.Location, hlt.STILL})
				} else {
					result[dest.Y][dest.X] += str
					smarterMoves = append(smarterMoves, m)
				}
			}
		}
		conn.SendFrame(smarterMoves)

	}
}
