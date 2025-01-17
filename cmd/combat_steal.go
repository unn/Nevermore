package cmd

import (
	"github.com/ArcCS/Nevermore/config"
	"github.com/ArcCS/Nevermore/objects"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/utils"
	"strconv"
)

// Syntax: ( INVENTORY | INV )
func init() {
	addHandler(steal{},
		"Usage:  steal target item \n \n Try to steal an item from a targets inventory",
		permissions.Thief,
		"steal")
}

type steal cmd

func (steal) process(s *state) {
	if len(s.input) < 2 {
		s.msg.Actor.SendBad("Steal what from who")
		return
	}

	// Check some timers
	ready, msg := s.actor.TimerReady("steal")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	if s.actor.Flags["hidden"] != true {
		s.msg.Actor.SendBad("You can't steal while standing out in the open.")
		return
	}

	if s.actor.Tier < 5 {
		s.msg.Actor.SendBad("You must be level 5 before you can steal.")
		return
	}

	targetStr := s.words[0]
	targetNum := 1
	nameStr := ""
	nameNum := 1

	if len(s.words) > 1 {
		if val, err := strconv.Atoi(s.words[1]); err == nil {
			targetNum = val
		} else {
			nameStr = s.words[1]
		}
	}
	if len(s.words) > 2 {
		if nameStr != "" {
			if val2, err2 := strconv.Atoi(s.words[2]); err2 == nil {
				nameNum = val2
			} else {
				nameStr = s.words[2]
			}
		}
	}

	if len(s.words) > 3 {
		if val3, err3 := strconv.Atoi(s.words[3]); err3 == nil {
			nameNum = val3
		}else{
			nameStr = s.words[3]
		}
	}

	//TODO: Steal from players inventory if PvP flag is set

	var whatMob *objects.Mob
	whatMob = s.where.Mobs.Search(targetStr, targetNum, s.actor)
	if whatMob != nil {
		if whatMob.Placement != s.actor.Placement {
			s.msg.Actor.SendBad("You are too far away to steal from ", whatMob.Name)
			return
		}

		if len(whatMob.ThreatTable) > 0 && !s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
			s.msg.Actor.SendBad("This mob is already in combat, you can't get a clear access to steal from it!")
			return
		}

		what := whatMob.Inventory.Search(nameStr, nameNum)
		if what != nil {
			s.actor.SetTimer("steal", config.StealCD)
			if (s.actor.Inventory.TotalWeight + what.GetWeight()) <= s.actor.MaxWeight() {
				// base chance is 15% to hide
				curChance := config.StealChance + (config.StealChancePerLevel*(s.actor.Tier - whatMob.Level))

				if s.actor.Permission.HasAnyFlags(permissions.Builder, permissions.Dungeonmaster, permissions.Gamemaster) {
					curChance = 100
				}

				curChance += s.actor.Dex.Current * config.StealChancePerPoint

				if curChance >= 100 || utils.Roll(100, 1, 0) <= curChance {
					whatMob.Inventory.Lock()
					s.actor.Inventory.Lock()
					whatMob.Inventory.Remove(what)
					s.actor.Inventory.Add(what)
					whatMob.Inventory.Unlock()
					s.actor.Inventory.Unlock()
					s.msg.Actor.SendGood("You steal a ", what.Name, " from ", whatMob.Name, ".")
					return
				}else{
					s.msg.Actor.SendBad("You failed to steal from ", whatMob.Name, ", and stumble out of the shadows.")
					s.actor.RemoveHook("combat", "hide")
					whatMob.AddThreatDamage(whatMob.Stam.Max/4, s.actor)
					return
				}
			} else {
				s.msg.Actor.SendInfo("That item weighs too much for you to add to your inventory.")
				return
				}
		}else{
			s.msg.Actor.SendInfo("That item isn't on the target.")
			return
		}
	}else{
		s.msg.Actor.SendBad("What are you trying to steal from?")
		s.ok = true
		return
	}

}



