package hlt

import (
	"encoding/json"
)

var Constants map[string]interface{}

func init() {
	Constants = make(map[string]interface{})
	err := json.Unmarshal([]byte(`{
    "ADDITIONAL_PRODUCTIVITY": 6,
    "BASE_PRODUCTIVITY": 6,
    "BASE_SHIP_HEALTH": 255,
    "DOCKED_SHIP_REGENERATION": 0,
    "DOCK_RADIUS": 4.0,
    "DOCK_TURNS": 5,
    "DRAG": 10.0,
    "EXPLOSION_RADIUS": 10.0,
    "EXTRA_PLANETS": 4,
    "INFINITE_RESOURCES": true,
    "MAX_ACCELERATION": 7.0,
    "MAX_SHIP_HEALTH": 255,
    "MAX_SPEED": 7.0,
    "MAX_TURNS": 300,
    "PLANETS_PER_PLAYER": 6,
    "PRODUCTION_PER_SHIP": 72,
    "RESOURCES_PER_RADIUS": 144,
    "SHIPS_PER_PLAYER": 3,
    "SHIP_RADIUS": 0.5,
    "SPAWN_RADIUS": 2,
    "WEAPON_COOLDOWN": 1,
    "WEAPON_DAMAGE": 64,
    "WEAPON_RADIUS": 5.0
}`), &Constants)
	if err != nil {
		panic(err)
	}

}
