/*
* CMPE 273 Assignment 3 - Trip planner using Uber API
* Neha Viswanathan
* 010029097
*/

package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Input struct {
		Name string `json:"name"`
		Address string `json:"address"`
		City string	`json:"city"`
		State string `json:"state"`
		Zip string `json:"zip"`
	}

type Output struct {
		Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
		Name string `json:"name"`
		Address string `json:"address"`
		City string	`json:"city" `
		State string `json:"state"`
		Zip string `json:"zip"`
		Coordinates struct {
			Latitude string `json:"latitude"`
			Longitude string `json:"longitude"`
		}
	}

//Google Maps Response struct -- start
type GoogleResponse struct {
	Result []GoogleResult
}

type GoogleResult struct {
	Address string `json:"formatted_address"`
	AddressParts []GoogleAddressPart `json:"address_components"`
	Geometry Geometry
	Types []string
}

type GoogleAddressPart struct {
	Name string `json:"long_name"`
	ShortName string `json:"short_name"`
	Types []string
}

type Geometry struct {
	Bounds   Bounds
	Location Point
	Type     string
	Viewport Bounds
}
type Bounds struct {
	NorthEast, SouthWest Point
}

type Point struct {
	Lat float64
	Lng float64
}
//Google Maps Response struct -- end
