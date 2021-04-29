package cmd

import (
	"github.com/ArcCS/Nevermore/data"
	"github.com/ArcCS/Nevermore/permissions"
	"github.com/ArcCS/Nevermore/spells"
	"github.com/ArcCS/Nevermore/stats"
	"github.com/ArcCS/Nevermore/utils"
	"log"
	"strconv"
	"strings"
)

func init() {
	addHandler(build_copy{}, "Usage: copy (room|mob|item) (SubjectID) \n \n Use this to copy an existing item in the database \n",
		permissions.Builder,
		"copy", "duplicate")
}

type build_copy cmd

func (build_copy) process(s *state) {
	// TODO: Finish this up to copy something from the database to make a new db entry
	return
	// Check arguments
	if len(s.words) < 3 {
		s.msg.Actor.SendInfo("Edit what, how?")
		return
	}

	//log.Println("Trying to edit: " + strings.ToLower(s.words[0]))
	switch strings.ToLower(s.words[0]) {
	// Handle Rooms
	case "room":
		// Toggle Flags
		if strings.ToLower(s.words[1]) == "toggle" {
			for _, flag := range s.input[2:] {
				if (s.actor.Permission.HasFlags(permissions.Builder, permissions.Dungeonmaster)) || flag != "active" {
					if s.where.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}
			}

			// Set a variable
		} else {
			switch strings.ToLower(s.words[1]) {
			case "description":
				s.where.Description = strings.Join(s.input[2:], " ")
				s.msg.Actor.SendGood("Description changed.")
			case "name":
				s.where.Name = strings.Join(s.input[2:], " ")
				s.msg.Actor.SendGood("Name changed.")
			default:
				s.msg.Actor.SendBad("Property not found.")
			}
		}
		s.where.Save()
		return

	// Handle Exits
	case "exit":
		// Toggle Flags
		exitName := s.input[2]
		log.Println("Attempting to edit ", exitName)
		objectRef := strings.ToLower(exitName)
		if !utils.StringIn(strings.ToUpper(objectRef), directionals) {
			for txtE := range s.where.Exits {
				if strings.Contains(txtE, objectRef) {
					objectRef = txtE
				}
			}
		}
		if exit, exists := s.where.Exits[objectRef]; exists {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if exit.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					exit.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					oldName := exit.Name
					exit.Name = strings.Join(s.input[3:], " ")
					s.where.Exits[strings.ToLower(strings.Join(s.input[3:], " "))] = exit
					delete(s.where.Exits, oldName)
					data.RenameExit(exit.Name, oldName, exit.ParentId, exit.ToId)
					s.msg.Actor.SendGood("Name changed.")
				case "key_id":
					intKey, _ := strconv.Atoi(s.words[3])
					exit.KeyId = intKey
					s.msg.Actor.SendGood("Change Key Id")
				case "placement":
					intKey, _ := strconv.Atoi(s.words[3])
					if intKey >= 1 && intKey <= 5 {
						exit.KeyId = intKey
						s.msg.Actor.SendGood("Changed placement")
					} else {
						s.msg.Actor.SendBad("Placement Id not valid. ")
					}
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			exit.Save()
		} else {
			s.msg.Actor.SendBad("Exit not found.")
		}

		return

	// Handle Items
	case "item":
		// Toggle Flags
		itemName := s.input[2]
		item := s.actor.Inventory.Search(itemName, 1)

		if item != nil {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if item.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					item.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					item.Name = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Name changed.")
				case "spell":
					if _, ok := spells.Spells[s.input[3]]; ok {
						item.Spell = s.input[3]
						s.msg.Actor.SendGood("Spell changed.")
					} else {
						s.msg.Actor.SendBad("Spell not found.")
					}
				case "weight":
					weight, _ := strconv.Atoi(s.words[3])
					item.Weight = weight
					s.msg.Actor.SendGood("Change weight")
				case "type":
					types, err := strconv.Atoi(s.words[3])
					if err != nil {
						s.msg.Actor.SendBad("Type must be an integer, use command 'types' to print types.")
						return
					}
					item.ItemType = types
					s.msg.Actor.SendGood("Changed types.")
				case "value":
					value, _ := strconv.Atoi(s.words[3])
					item.Value = value
					s.msg.Actor.SendGood("Changed value")
				case "ndice":
					value, _ := strconv.Atoi(s.words[3])
					item.NumDice = value
					s.msg.Actor.SendGood("Changed number of dice")
				case "armor":
					value, _ := strconv.Atoi(s.words[3])
					item.Armor = value
					s.msg.Actor.SendGood("Changed armor value")
				case "pdice":
					value, _ := strconv.Atoi(s.words[3])
					item.PlusDice = value
					s.msg.Actor.SendGood("Changed plus dice")
				case "sdice":
					value, _ := strconv.Atoi(s.words[3])
					item.SidesDice = value
					s.msg.Actor.SendGood("Changed sides of dice")
				case "max_uses":
					value, _ := strconv.Atoi(s.words[3])
					item.MaxUses = value
					s.msg.Actor.SendGood("Changed max_uses")
				/*case "placement":
				intKey, _ :=  strconv.Atoi(s.words[3])
				if intKey >= 1 && intKey <= 5 {
					exit.KeyId = intKey
					s.msg.Actor.SendGood("Changed placement")
				}else{
					s.msg.Actor.SendBad("Placement Id not valid. ")
				}*/
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			item.Save()
		} else {
			s.msg.Actor.SendBad("Item not found.")
		}

		return

	// Handle Mobs
	case "mob":
		// Toggle Flags
		mobName := s.input[2]
		mob := s.where.Mobs.Search(mobName, 1, true)

		if mob != nil {
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if mob.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					mob.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					mob.Name = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Name changed.")
				case "level":
					value, _ := strconv.Atoi(s.words[3])
					mob.Level = value
					s.msg.Actor.SendGood("Change level")
				case "experience":
					types, _ := strconv.Atoi(s.words[3])
					mob.Experience = types
					s.msg.Actor.SendGood("Changed experience value.")
				case "gold":
					value, _ := strconv.Atoi(s.words[3])
					mob.Gold = value
					s.msg.Actor.SendGood("Changed amount of gold dropped.")
				case "con":
					value, _ := strconv.Atoi(s.words[3])
					mob.Con.Current = value
					s.msg.Actor.SendGood("Changed constitution")
				case "int":
					value, _ := strconv.Atoi(s.words[3])
					mob.Int.Current = value
					s.msg.Actor.SendGood("Changed intelligence")
				case "str":
					value, _ := strconv.Atoi(s.words[3])
					mob.Str.Current = value
					s.msg.Actor.SendGood("Changed strength")
				case "dex":
					value, _ := strconv.Atoi(s.words[3])
					mob.Dex.Current = value
					s.msg.Actor.SendGood("Changed dexterity")
				case "pie":
					value, _ := strconv.Atoi(s.words[3])
					mob.Pie.Current = value
					s.msg.Actor.SendGood("Changed piety")
				case "mana":
					value, _ := strconv.Atoi(s.words[3])
					mob.Mana.Max = value
					s.msg.Actor.SendGood("Changed mana")
				case "stam":
					value, _ := strconv.Atoi(s.words[3])
					mob.Stam.Max = value
					s.msg.Actor.SendGood("Changed stam")
				case "ndice":
					value, _ := strconv.Atoi(s.words[3])
					mob.NumDice = value
					s.msg.Actor.SendGood("Changed number of dice")
				case "armor":
					value, _ := strconv.Atoi(s.words[3])
					mob.Armor = value
					s.msg.Actor.SendGood("Changed armor value")
				case "pdice":
					value, _ := strconv.Atoi(s.words[3])
					mob.PlusDice = value
					s.msg.Actor.SendGood("Changed plus dice")
				case "sdice":
					value, _ := strconv.Atoi(s.words[3])
					mob.SidesDice = value
					s.msg.Actor.SendGood("Changed sides of dice")
				case "chancecast":
					value, _ := strconv.Atoi(s.words[3])
					mob.ChanceCast = value
					s.msg.Actor.SendGood("Changed chance to cast")
				case "numwander":
					value, _ := strconv.Atoi(s.words[3])
					mob.NumWander = value
					s.msg.Actor.SendGood("Changed amount of ticks to wander")
				case "wimpyvalue":
					value, _ := strconv.Atoi(s.words[3])
					mob.WimpyValue = value
					s.msg.Actor.SendGood("Changed value that mob tries to flee")
				/*case "placement":
				intKey, _ :=  strconv.Atoi(s.words[3])
				if intKey >= 1 && intKey <= 5 {
					exit.KeyId = intKey
					s.msg.Actor.SendGood("Changed placement")
				}else{
					s.msg.Actor.SendBad("Placement Id not valid. ")
				}*/
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			mob.Save()
		} else {
			s.msg.Actor.SendBad("Mob not found.")
		}

		return
		// Handle Mobs
	case "char":
		// Toggle Flags
		charName := s.input[2]

		character := stats.ActiveCharacters.Find(charName)

		if character != nil {
			stats.ActiveCharacters.Lock()
			if strings.ToLower(s.input[1]) == "toggle" {
				for _, flag := range s.input[3:] {
					if character.ToggleFlag(strings.ToLower(flag)) {
						s.msg.Actor.SendGood("Toggled " + flag)
					} else {
						s.msg.Actor.SendBad("Failed to toggle " + flag + ".  Is it an actual flag?")
					}
				}

				// Set a variable
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					character.Description = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					character.Name = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Name changed.")
				case "title":
					character.Title = strings.Join(s.input[3:], " ")
					s.msg.Actor.SendGood("Title changed.")
				case "tier":
					value, _ := strconv.Atoi(s.words[3])
					character.Tier = value
					s.msg.Actor.SendGood("Changed Tier")
				case "experience":
					types, _ := strconv.Atoi(s.words[3])
					character.Experience.Value = types
					s.msg.Actor.SendGood("Changed amount of experience.")
				case "gold":
					value, _ := strconv.Atoi(s.words[3])
					character.Gold.Value = value
					s.msg.Actor.SendGood("Changed amount of gold on character")
				case "bankgold":
					value, _ := strconv.Atoi(s.words[3])
					character.BankGold.Value = value
					s.msg.Actor.SendGood("Changed amount of gold in bank.")
				case "passages":
					value, _ := strconv.Atoi(s.words[3])
					character.Passages.Value = value
					s.msg.Actor.SendGood("Changed number of passages")
				case "bonuspoints":
					value, _ := strconv.Atoi(s.words[3])
					character.BonusPoints.Value = value
					s.msg.Actor.SendGood("Changed number of bonus points")
				case "broadcasts":
					value, _ := strconv.Atoi(s.words[3])
					character.Broadcasts = value
					s.msg.Actor.SendGood("Changed broadcasts.")
				case "evals":
					value, _ := strconv.Atoi(s.words[3])
					character.Broadcasts = value
					s.msg.Actor.SendGood("Changed evals")
				case "concur":
					value, _ := strconv.Atoi(s.words[3])
					character.Con.Current = value
					s.msg.Actor.SendGood("Changed constitution")
				case "intcur":
					value, _ := strconv.Atoi(s.words[3])
					character.Int.Current = value
					s.msg.Actor.SendGood("Changed intelligence")
				case "strcur":
					value, _ := strconv.Atoi(s.words[3])
					character.Str.Current = value
					s.msg.Actor.SendGood("Changed strength")
				case "dexcur":
					value, _ := strconv.Atoi(s.words[3])
					character.Dex.Current = value
					s.msg.Actor.SendGood("Changed dexterity")
				case "piecur":
					value, _ := strconv.Atoi(s.words[3])
					character.Pie.Current = value
					s.msg.Actor.SendGood("Changed piety")
				case "stamcur":
					value, _ := strconv.Atoi(s.words[3])
					character.Stam.Current = value
					s.msg.Actor.SendGood("Changed current stamina")
				case "stammax":
					value, _ := strconv.Atoi(s.words[3])
					character.Stam.Max = value
					s.msg.Actor.SendGood("Changed stamina maximum")
				case "manamax":
					value, _ := strconv.Atoi(s.words[3])
					character.Mana.Max = value
					s.msg.Actor.SendGood("Changed mana maximum")
				case "manacur":
					value, _ := strconv.Atoi(s.words[3])
					character.Mana.Current = value
					s.msg.Actor.SendGood("Changed mana current")
				case "vitmax":
					value, _ := strconv.Atoi(s.words[3])
					character.Vit.Max = value
					s.msg.Actor.SendGood("Changed vitality max")
				case "vitcur":
					value, _ := strconv.Atoi(s.words[3])
					character.Vit.Current = value
					s.msg.Actor.SendGood("Changed vit current")
				case "sharpexp":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[0].Value = value
					s.msg.Actor.SendGood("Changed sharp exp")
				case "thrustexp":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[1].Value = value
					s.msg.Actor.SendGood("Changed thrust exp")
				case "bluntexp":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[2].Value = value
					s.msg.Actor.SendGood("Changed blunt exp")
				case "poleexp":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[3].Value = value
					s.msg.Actor.SendGood("Changed pole exp")
				case "missileexp":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[4].Value = value
					s.msg.Actor.SendGood("Changed missile exp")

				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
			character.Save()
			stats.ActiveCharacters.Unlock()
		} else {
			if strings.ToLower(s.input[1]) == "toggle" {
				s.msg.Actor.SendBad("Cannot toggle offline character.")
			} else {
				switch strings.ToLower(s.input[1]) {
				case "description":
					data.SaveCharField(charName, "description", strings.Join(s.input[3:], " "))
					s.msg.Actor.SendGood("Description changed.")
				case "name":
					data.SaveCharField(charName, "name", strings.Join(s.input[3:], " "))
					s.msg.Actor.SendGood("Name changed.")
				case "title":
					data.SaveCharField(charName, "title", strings.Join(s.input[3:], " "))
					s.msg.Actor.SendGood("Title changed.")
				case "tier":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "tier", value)
					s.msg.Actor.SendGood("Changed Tier")
				case "experience":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "experience", value)
					s.msg.Actor.SendGood("Changed amount of experience.")
				case "gold":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "gold", value)
					s.msg.Actor.SendGood("Changed amount of gold on character")
				case "bankgold":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "bankgold", value)
					s.msg.Actor.SendGood("Changed amount of gold in bank.")
				case "passages":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "passages", value)
					s.msg.Actor.SendGood("Changed number of passages")
				case "bonuspoints":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "bonuspoints", value)
					s.msg.Actor.SendGood("Changed number of bonus points")
				case "broadcasts":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "broadcasts", value)
					s.msg.Actor.SendGood("Changed broadcasts.")
				case "evals":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "evals", value)
					s.msg.Actor.SendGood("Changed evals")
				case "concur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "concur", value)
					s.msg.Actor.SendGood("Changed constitution")
				case "intcur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "intcur", value)
					s.msg.Actor.SendGood("Changed intelligence")
				case "strcur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "strcur", value)
					s.msg.Actor.SendGood("Changed strength")
				case "dexcur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "dexcur", value)
					s.msg.Actor.SendGood("Changed dexterity")
				case "piecur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "piecur", value)
					s.msg.Actor.SendGood("Changed piety")
				case "stamcur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "curr_stam", value)
					s.msg.Actor.SendGood("Changed current stamina")
				case "stammax":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "max_stam", value)
					s.msg.Actor.SendGood("Changed stamina maximum")
				case "manamax":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "max_mana", value)
					s.msg.Actor.SendGood("Changed mana maximum")
				case "manacur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "curr_Mana", value)
					s.msg.Actor.SendGood("Changed mana current")
				case "vitmax":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "max_vit", value)
					s.msg.Actor.SendGood("Changed vitality max")
				case "vitcur":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "curr_vit", value)
					s.msg.Actor.SendGood("Changed vit current")
				case "sharpexp":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "sharpexp", value)
					s.msg.Actor.SendGood("Changed sharp exp")
				case "thrustexp":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "thrustexp", value)
					s.msg.Actor.SendGood("Changed thrust exp")
				case "bluntexp":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "bluntexp", value)
					s.msg.Actor.SendGood("Changed blunt exp")
				case "poleexp":
					value, _ := strconv.Atoi(s.words[3])
					data.SaveCharField(charName, "poleexp", value)
					s.msg.Actor.SendGood("Changed pole exp")
				case "missileexp":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[4].Value = value
					s.msg.Actor.SendGood("Changed missile exp")
				case "parent_id":
					value, _ := strconv.Atoi(s.words[3])
					character.Skills[4].Value = value
					s.msg.Actor.SendGood("Changed room")
				default:
					s.msg.Actor.SendBad("Property not found.")
				}
			}
		}

		return
	default:
		s.msg.Actor.SendBad("Not an object that can be edited.")
	}

	s.ok = true
	return
}