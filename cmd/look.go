package cmd

import (
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"math"
	"strconv"
	"strings"
)

// Overloaded Look object for all of your looking pleasure
// Syntax: ( LOOK | L ) has.Thing
func init() {
	addHandler(look{},
           "Usage:  look [object|exit|character|mob] # \n \n Put your peepers on something. (Also can use short hand L",
           permissions.Player,
           "LOOK", "L")
}

type look cmd

func (look) process(s *state) {
	var others []string
	var mobs []string
	if len(s.input) == 0 {
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Actor.SendInfo(objects.Rooms[s.actor.ParentId].Look(true))
		}else{
			s.msg.Actor.SendInfo(objects.Rooms[s.actor.ParentId].Look(false))
		}
		// Pick whether it's a GM or a user looking and go for it.
		if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			others = objects.Rooms[s.actor.ParentId].Chars.List(true, s.actor.Name, true)
			mobs = objects.Rooms[s.actor.ParentId].Mobs.List(true, true)
		}else{
			others = objects.Rooms[s.actor.ParentId].Chars.List(false, s.actor.Name, false)
			mobs = objects.Rooms[s.actor.ParentId].Mobs.List(false, false)
		}
		if len(others) == 1 {
			s.msg.Actor.SendInfo(strings.Join(others, ", "), " is also here.")
		} else if len(others) > 1{
			s.msg.Actor.SendInfo(strings.Join(others, ", "), " are also here.")
		}
		//log.Println("Mobs here:" + strconv.Itoa(len(mobs)))
		if len(mobs) == 1 {
			s.msg.Actor.SendInfo("You see: " + strings.Join(mobs, ", "))
		} else if len(mobs) > 1{
			s.msg.Actor.SendInfo("You see: " + strings.Join(mobs, ", "))
		}
		return
	}

	name := s.input[0]
	nameNum := 1

	if len(s.words) > 1 {
		// Try to snag a number off the list
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			nameNum = val
		}
	}

	var whatChar *objects.Character
	// Check characters in the room first.
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		whatChar = s.where.Chars.Search(name, true)
	}else{
		whatChar = s.where.Chars.Search(name, false)
	}
	// It was a person!
	if whatChar != nil {
		s.msg.Actor.SendInfo(whatChar.Look())
		return
	}

	// Check exits
	whatExit := s.where.FindExit(strings.ToLower(name))

	// Nice, looking at an exit.
	if whatExit != nil {
		s.msg.Actor.SendInfo(whatExit.Look())
		return
	}

	// Check mobs
	var whatMob *objects.Mob
	if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
		whatMob = s.where.Mobs.Search(name, nameNum,true)
	}else{
		whatMob = s.where.Mobs.Search(name, nameNum,false)
	}
	// It was a mob!
	if whatMob != nil {
		s.msg.Actor.SendInfo(whatMob.Look())
		s.msg.Actor.SendInfo("It is standing" + WhereAt(whatMob.Placement, s.actor.Placement))
		_, err := whatMob.ThreatTable[s.actor.Name]; if err {
			s.msg.Actor.SendInfo("It isn't paying attention to you.")
		}else{
			s.msg.Actor.SendInfo("It appears very angry at you!")
		}
		if whatMob.CurrentTarget == s.actor.Name {
			s.msg.Actor.SendInfo("It is attacking you!")
		}else if whatMob.CurrentTarget != "" {
			s.msg.Actor.SendInfo("it is attacking " + whatMob.CurrentTarget)
		}
		return
	}

	// Check items
	what := s.where.Items.Search(name, nameNum)

	// Item in the room?
	if what != nil {
		s.msg.Actor.SendInfo(what.Look())
		s.msg.Actor.SendInfo("It is on the ground " + WhereAt(whatMob.Placement, s.actor.Placement))
		return
	}

	what = s.actor.Inventory.Search(name, nameNum)

	// It was on you the whole time
	if what != nil {
		s.msg.Actor.SendInfo(what.Look())
		return
	}else{
		s.msg.Actor.SendBad("You see no '", name, "' to examine.")
		return
	}
}

func WhereAt(subLoc int64, charLoc int64) string {
	// Moving backwards
	if subLoc == charLoc {
		return " next to you."
	}
	diff := math.Abs(float64(subLoc - charLoc))
	steps := ""
	direction := ""
	if subLoc > charLoc {
		direction = "in front of you."
	}else{
		direction = "behind you."
	}
	if diff == 1 {
		steps = " a couple steps "
	}else if diff == 2 {
		steps = " a dozen steps "
	}else if diff == 3 {
		steps = " many steps "
	}else if diff == 4 {
		steps = " at the other end of the room "
	}
	return steps + direction
}