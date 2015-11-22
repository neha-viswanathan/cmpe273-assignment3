/*
* CMPE 273 Assignment 3 - Trip planner using Uber API
* Neha Viswanathan
* 010029097
*/

package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
  "strconv"
  "io/ioutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
  "github.com/julienschmidt/httprouter"
  "github.com/neha-viswanathan/cmpe273-assignment3/uberService"
)

// LocationController to control Input resource
type LocationController struct {
		session *mgo.Session
	}

	type Input struct {
			Name string `json:"name"`
			Address string `json:"address"`
			City string	`json:"city"`
			State string `json:"state"`
			Zipcode string `json:"zip"`
		}

	type Output struct {
			Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
			Name string `json:"name"`
			Address string `json:"address"`
			City string	`json:"city" `
			State string `json:"state"`
			Zipcode string `json:"zip"`
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

//Output struct for Uber
type UberOut struct {
	Distance float64
  Cost int
	Duration int
}

//Struct for Uber Trip - Input for POST request
type TripPostInput struct{
	Starting_from_location_id string `json:"starting_from_location_id"`
	Location_ids []string
}

//Struct for Uber Trip - Output for POST request
type TripPostOutput struct{
	Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Status string `json:"status"`
	Starting_from_location_id string `json:"starting_from_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int `json:"total_uber_costs"`
	Total_uber_duration int	`json:"total_uber_duration"`
	Total_distance float64 `json:"total_distance"`
}

//Struct for Uber Trip - Output for PUT request
type TripPutOutput struct{
	Id bson.ObjectId `json:"_id" bson:"_id,omitempty"`
	Status string `json:"status"`
	Starting_from_location_id string `json:"starting_from_location_id"`
	Next_destination_location_id string `json:"next_destination_location_id"`
	Best_route_location_ids []string
	Total_uber_costs int `json:"total_uber_costs"`
	Total_uber_duration int	`json:"total_uber_duration"`
	Total_distance float64 `json:"total_distance"`
	Uber_wait_time_eta int `json:"uber_wait_time_eta"`
}

type putStruct struct{
	trip_route []string
	trip_visits map[string]int
}

type Final_struct struct{
	result map[string]putStruct
}

// NewLocationController for providing reference to LocationController with provided mongo session
func NewLocationController(s *mgo.Session) *LocationController {
	return &LocationController{s}
}

//Accessing Google Maps API to retrieve co-ordinates
func getGeoLocation(addr string) Output{
	//create a client to get data from Google Maps
	client := &http.Client{}

	req_URL := "http://maps.google.com/maps/api/geocode/json?address="
	req_URL += url.QueryEscape(addr)

	req_URL += "&sensor=false";
	fmt.Println("Request URL :: ", req_URL)

	request, err := http.NewRequest("GET", req_URL , nil)

	response, err := client.Do(request)
	if err != nil {
		fmt.Println("Error while calling Google Maps API :: ", err);
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error while reading response :: ", err);
	}

	//parsing the json response
	var geoRes GoogleResponse
	err = json.Unmarshal(body, &geoRes)
	if err != nil {
		fmt.Println("Error while unmarshalling response :: ", err);
	}

	//retrieve latitude and longitude from GoogleResponse
	var val Output
	val.Coordinates.Latitude = strconv.FormatFloat(geoRes.Result[0].Geometry.Location.Lat,'f',7,64)
	val.Coordinates.Longitude = strconv.FormatFloat(geoRes.Result[0].Geometry.Location.Lng,'f',7,64)
	fmt.Println("Retrieved co-ordinates :: ", val.Coordinates)
	return val;
}

// GetLocation retrieves an individual location resource
func (uc LocationController) GetLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("location_id")
	// fmt.Println(id)
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    // Grab id
    oid := bson.ObjectIdHex(id)
	var o Output
	if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }
	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateLocation creates a new Location resource
func (uc LocationController) CreateLocation(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var u Input
	var oA Output

	json.NewDecoder(r.Body).Decode(&u)
	googResCoor := getGeoLocation(u.Address + "+" + u.City + "+" + u.State + "+" + u.Zipcode);
    fmt.Println("resp is :: ", googResCoor.Coordinates.Latitude, googResCoor.Coordinates.Longitude);

	// oA.Id = bson.NewObjectId()
	oA.Name = u.Name
	oA.Address = u.Address
	oA.City= u.City
	oA.State= u.State
	oA.Zipcode = u.Zipcode
	oA.Coordinates.Latitude = googResCoor.Coordinates.Latitude
	oA.Coordinates.Longitude = googResCoor.Coordinates.Longitude

	// Write the user to mongo
	uc.session.DB("cmpe273").C("GeoLocations").Insert(oA)

	// Marshal provided interface into JSON structure
	uj, _ := json.Marshal(oA)
	// Write content-type, statuscode, payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)


	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

// CreateTrip creates a new Trip
func (uc LocationController) CreateTrip(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var tI TripPostInput
	var tO TripPostOutput
	var cost_array []int
	var duration_array []int
	var distance_array []float64
	cost_total := 0
	duration_total := 0
	distance_total := 0.0

	json.NewDecoder(r.Body).Decode(&tI)

	starting_id:= bson.ObjectIdHex(tI.Starting_from_location_id)
	var start Output
	if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(starting_id).One(&start); err != nil {
       	w.WriteHeader(404)
        return
    }
    start_Lat := start.Coordinates.Latitude
    start_Lang := start.Coordinates.Longitude

    for len(tI.Location_ids)>0 {

			for _, loc := range tI.Location_ids {
				id := bson.ObjectIdHex(loc)
				var o Output
				if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(id).One(&o); err != nil {
		       		w.WriteHeader(404)
		        	return
		    	}
		    	loc_Lat := o.Coordinates.Latitude
		    	loc_Lang := o.Coordinates.Longitude

		    	getUberResponse := uberService.GetUberRate(start_Lat, start_Lang, loc_Lat, loc_Lang)
					distance_array = append(distance_array, getUberResponse.Distance)
					cost_array = append(cost_array, getUberResponse.Cost)
		    	duration_array = append(duration_array, getUberResponse.Duration)
			}
			fmt.Println("Cost Array", cost_array)

			min_cost:= cost_array[0]
			var indexNeeded int
			for index, value := range cost_array {
		        if value < min_cost {
		            min_cost = value
		            indexNeeded = index
		        }
		    }

			cost_total += min_cost
			duration_total += duration_array[indexNeeded]
			distance_total += distance_array[indexNeeded]

			tO.Best_route_location_ids = append(tO.Best_route_location_ids, tI.Location_ids[indexNeeded])

			starting_id = bson.ObjectIdHex(tI.Location_ids[indexNeeded])
			if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(starting_id).One(&start); err != nil {
       			w.WriteHeader(404)
        		return
    		}
    		tI.Location_ids = append(tI.Location_ids[:indexNeeded], tI.Location_ids[indexNeeded+1:]...)

    		start_Lat = start.Coordinates.Latitude
    		start_Lang = start.Coordinates.Longitude

    		cost_array = cost_array[:0]
    		duration_array = duration_array[:0]
    		distance_array = distance_array[:0]
	}

	Last_loc_id := bson.ObjectIdHex(tO.Best_route_location_ids[len(tO.Best_route_location_ids)-1])
	var o2 Output
	if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(Last_loc_id).One(&o2); err != nil {
		w.WriteHeader(404)
		return
	}
	last_loc_Lat := o2.Coordinates.Latitude
	last_loc_Lang := o2.Coordinates.Longitude

	ending_id:= bson.ObjectIdHex(tI.Starting_from_location_id)
	var end Output
	if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(ending_id).One(&end); err != nil {
       	w.WriteHeader(404)
        return
    }
    end_Lat := end.Coordinates.Latitude
    end_Lang := end.Coordinates.Longitude

	getUberResponse_last := uberService.GetUberRate(last_loc_Lat, last_loc_Lang, end_Lat, end_Lang)


	tO.Id = bson.NewObjectId()
	tO.Status = "planning"
	tO.Starting_from_location_id = tI.Starting_from_location_id
	tO.Total_uber_costs = cost_total + getUberResponse_last.Cost
	tO.Total_distance = distance_total + getUberResponse_last.Distance
	tO.Total_uber_duration = duration_total + getUberResponse_last.Duration

	uc.session.DB("cmpe273").C("UberTrips").Insert(tO)

	uj, _ := json.Marshal(tO)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}

type Internal_data struct{
	Id string `json:"_id" bson:"_id,omitempty"`
	Trip_visited []string  `json:"trip_visited"`
	Trip_not_visited []string  `json:"trip_not_visited"`
	Trip_completed int `json:"trip_completed"`
}

// GetTrip retrieves an individual trip resource
func (uc LocationController) GetTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	// Grab id
	id := p.ByName("trip_id")

	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }

    oid := bson.ObjectIdHex(id)
	var tO TripPostOutput
	if err := uc.session.DB("cmpe273").C("UberTrips").FindId(oid).One(&tO); err != nil {
        w.WriteHeader(404)
        return
    }

	uj, _ := json.Marshal(tO)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}

//UpdateTrip updates an existing location resource
func (uc LocationController) UpdateTrip(w http.ResponseWriter, r *http.Request, p httprouter.Params) {

	var theStruct putStruct
	var final Final_struct
	final.result = make(map[string]putStruct)

	var tPO TripPutOutput
	var internal Internal_data

	id := p[0].Value
	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)
	if err := uc.session.DB("cmpe273").C("UberTrips").FindId(oid).One(&tPO); err != nil {
        w.WriteHeader(404)
        return
    }

	theStruct.trip_route = tPO.Best_route_location_ids
    theStruct.trip_route = append([]string{tPO.Starting_from_location_id}, theStruct.trip_route...)
    fmt.Println("The route array is: ", theStruct.trip_route)
    theStruct.trip_visits = make(map[string]int)

    var trip_visited []string
    var trip_not_visited []string

  	if err := uc.session.DB("cmpe273").C("Trip_Temp").FindId(id).One(&internal); err != nil {
    	for index, loc := range theStruct.trip_route{
    		if index == 0 {
    			theStruct.trip_visits[loc] = 1
    			trip_visited = append(trip_visited, loc)
    		}else{
    			theStruct.trip_visits[loc] = 0
    			trip_not_visited = append(trip_not_visited, loc)
    		}
    	}
    	internal.Id = id
    	internal.Trip_visited = trip_visited
    	internal.Trip_not_visited = trip_not_visited
    	internal.Trip_completed = 0
    	uc.session.DB("cmpe273").C("Trip_Temp").Insert(internal)

    }else {
    	for _, loc_id := range internal.Trip_visited {
    		theStruct.trip_visits[loc_id] = 1
    	}
    	for _, loc_id := range internal.Trip_not_visited {
    		theStruct.trip_visits[loc_id] = 0
    	}
    }


  	fmt.Println("Trip visit map ", theStruct.trip_visits)
  	final.result[id] = theStruct


  	last_index := len(theStruct.trip_route) - 1
  	trip_completed := internal.Trip_completed
  	if trip_completed == 1 {
  		fmt.Println("Entering the trip completed if statement")
  		tPO.Status = "completed"

		uj, _ := json.Marshal(tPO)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		fmt.Fprintf(w, "%s", uj)
		return
	}

	for i, location := range theStruct.trip_route{
	  	if  (theStruct.trip_visits[location] == 0){
	  		tPO.Next_destination_location_id = location
	  		nextoid := bson.ObjectIdHex(location)
			var o Output
			if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(nextoid).One(&o); err != nil {
        		w.WriteHeader(404)
        		return
    		}
    		nlat := o.Coordinates.Latitude
    		nlang:= o.Coordinates.Longitude

	  		if i == 0 {
	  			starting_point := theStruct.trip_route[last_index]
	  			startingoid := bson.ObjectIdHex(starting_point)
				var o Output
				if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(startingoid).One(&o); err != nil {
        			w.WriteHeader(404)
        			return
    			}
    			slat := o.Coordinates.Latitude
    			slang:= o.Coordinates.Longitude


	  			eta := uberService.GetUberETA(slat, slang, nlat, nlang)
	  			tPO.Uber_wait_time_eta = eta
	  			trip_completed = 1
	  		}else {
	  			starting_point2 := theStruct.trip_route[i-1]
	  			startingoid2 := bson.ObjectIdHex(starting_point2)
				var o Output
				if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(startingoid2).One(&o); err != nil {
        			w.WriteHeader(404)
        			return
    			}
    			slat := o.Coordinates.Latitude
    			slang:= o.Coordinates.Longitude
	  			eta := uberService.GetUberETA(slat, slang, nlat, nlang)
	  			tPO.Uber_wait_time_eta = eta
	  		}

	  		fmt.Println("Starting Location :: ", tPO.Starting_from_location_id)
	  		fmt.Println("Next destination :: ", tPO.Next_destination_location_id)
	  		theStruct.trip_visits[location] = 1
	  		if i == last_index {
	  			theStruct.trip_visits[theStruct.trip_route[0]] = 0
	  		}
	  		break
	  	}
	}

	trip_visited  = trip_visited[:0]
	trip_not_visited  = trip_not_visited[:0]
	for location, visit := range theStruct.trip_visits{
		if visit == 1 {
			trip_visited = append(trip_visited, location)
		}else {
			trip_not_visited = append(trip_not_visited, location)
		}
	}

	internal.Id = id
	internal.Trip_visited = trip_visited
	internal.Trip_not_visited = trip_not_visited
	fmt.Println("Trip Visited", internal.Trip_visited)
	fmt.Println("Trip Not Visited", internal.Trip_not_visited)
	internal.Trip_completed = trip_completed

	c := uc.session.DB("cmpe273").C("Trip_Temp")
	id2 := bson.M{"_id": id}
	err := c.Update(id2, internal)
	if err != nil {
		panic(err)
	}

    uj, _ := json.Marshal(tPO)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)

}

// RemoveLocation removes an existing location resource
func (uc LocationController) RemoveLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id := p.ByName("location_id")

	if !bson.IsObjectIdHex(id) {
		w.WriteHeader(404)
		return
	}

	oid := bson.ObjectIdHex(id)

	if err := uc.session.DB("cmpe273").C("GeoLocations").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	w.WriteHeader(200)
}

//UpdateLocation updates an existing location resource
func (uc LocationController) UpdateLocation(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var i Input
	var o Output

	id := p.ByName("location_id")

	if !bson.IsObjectIdHex(id) {
        w.WriteHeader(404)
        return
    }
    oid := bson.ObjectIdHex(id)

	if err := uc.session.DB("cmpe273").C("GeoLocations").FindId(oid).One(&o); err != nil {
        w.WriteHeader(404)
        return
    }

	json.NewDecoder(r.Body).Decode(&i)

	googResCoor := getGeoLocation(i.Address + "+" + i.City + "+" + i.State + "+" + i.Zipcode);
    fmt.Println("Response is :: ", googResCoor.Coordinates.Latitude, googResCoor.Coordinates.Longitude);


	o.Address = i.Address
	o.City = i.City
	o.State = i.State
	o.Zipcode = i.Zipcode
	o.Coordinates.Latitude = googResCoor.Coordinates.Latitude
	o.Coordinates.Longitude = googResCoor.Coordinates.Longitude

	c := uc.session.DB("cmpe273").C("GeoLocations")

	id2 := bson.M{"_id": oid}
	err := c.Update(id2, o)
	if err != nil {
		panic(err)
	}

	uj, _ := json.Marshal(o)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)
}
