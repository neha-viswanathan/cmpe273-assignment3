/*
* CMPE 273 Assignment 3 - Trip planner using Uber API
* Neha Viswanathan
* 010029097
*/

package main

//import statements
import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"gopkg.in/mgo.v2"
  "github.com/neha-viswanathan/cmpe273-assignment3/controller"
)

//main function
func main() {
	// Instantiate a new router
	route := httprouter.New()

	// Create a new LocationController instance
	locCtrl := controller.NewLocationController(getSession())

	// Retrieve a location
	route.GET("/locations/:location_id", locCtrl.GetLocation)

	// Retrieve a trip
	route.GET("/trips/:trip_id", locCtrl.GetTrip)

	// Create a new location
	route.POST("/locations", locCtrl.CreateLocation)

	// Create a new trip
	route.POST("/trips", locCtrl.CreateTrip)

	// Update a location
	route.PUT("/locations/:location_id", locCtrl.UpdateLocation)

	// Update a trip
	route.PUT("/trips/:trip_id/request", locCtrl.UpdateTrip)

	// Remove an existing location
	route.DELETE("/locations/:location_id", locCtrl.RemoveLocation)

	// Start server
	http.ListenAndServe("localhost:1111", route)
}

// getSession function creates a new mongo session and panics if connection error occurs
func getSession() *mgo.Session {
  //Connect to MongoDB
	sess, err := mgo.Dial("mongodb://admin:root123@ds043324.mongolab.com:43324/cmpe273")

  //panic if connection failed
  if err != nil {
    panic(err)
  }

  //to make data reading consistent across sequential queries in same session
  sess.SetMode(mgo.Monotonic, true)

  return sess
  }
