package cmd

import (
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/text"
)

func init() {
	addHandler(haste{},
		"Usage:  haste \n\n Hasten your actions temporarily",
		permissions.Ranger,
		"haste")
}

type haste cmd

func (haste) process(s *state) {
	if s.actor.Tier < 10 {
		s.msg.Actor.SendBad("You aren't high enough level to perform that skill.")
		return
	}
	berz, ok := s.actor.Flags["haste"]
	if ok {
		if berz {
			s.msg.Actor.SendBad("You've already hastened your actions!!")
			return
		}
	}
	ready, msg := s.actor.TimerReady("combat_haste")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}
	ready, msg = s.actor.TimerReady("combat")
	if !ready {
		s.msg.Actor.SendBad(msg)
		return
	}

	s.actor.ApplyEffect("haste", "60", "0",
		func() {
			s.actor.ToggleFlagAndMsg("berserk", text.Red+"Your muscles tighten and your reflexes hasten!!!\n")
			s.actor.Dex.Current += 5
		},
		func() {
			s.actor.ToggleFlagAndMsg("haste", text.Cyan+"Your reflexes return to normal.\n")
			s.actor.Dex.Current -= 5
		})
	s.msg.Observers.SendInfo(s.actor.Name + " begins moving faster!")
	s.actor.SetTimer("combat_berserk", 60*10)

	s.ok = true
}
