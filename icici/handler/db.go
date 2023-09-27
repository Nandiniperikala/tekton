package handler

import (
	config "com.nandini.icici/db"
	pb "com.nandini.icici/proto"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/micro/micro/v3/service/logger"
	"gorm.io/gorm"
	"net/http"
)

var tableName = "event"
var dbClient *gorm.DB

func InitializeDb() {
	dbClient = config.GetClient()
}

func AddEvent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var event *pb.Event
	_ = json.NewDecoder(request.Body).Decode(&event)
	e := dbClient.Table(tableName).Create(&event)
	if e.Error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + e.Error.Error() + `" }`))
		return
	}
	logger.Infof("Inserted data with Id:" + event.Id)
	json.NewEncoder(response).Encode(event)

}

func ReadEventById(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(r)
	id := params["id"]
	var event *pb.Event
	e := dbClient.Table(tableName).First(&event, id)
	if e.Error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + e.Error.Error() + `" }`))
		return
	}
	logger.Infof("Fetched Event With Id:" + id)
	json.NewEncoder(response).Encode(event)
}

func GetEvents(response http.ResponseWriter, r *http.Request) {
	response.Header().Set("content-type", "application/json")
	var events []*pb.Event
	e := dbClient.Table(tableName).Find(&events)
	if e.Error != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`{ "message": "` + e.Error.Error() + `" }`))
		return
	}
	logger.Infof("Fetched all Events")
	json.NewEncoder(response).Encode(events)
}

func UpdateEvent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	var event *pb.Event
	var ev *pb.Event
	_ = json.NewDecoder(request.Body).Decode(&event)
	dbClient.Table(tableName).First(&ev, event.Id)
	dbClient.Table(tableName).Model(&ev).Updates(event)
	logger.Infof("Updated Event with Id:" + event.Id)
	json.NewEncoder(response).Encode("Updated Event with Id:" + event.Id)
}

func DeleteEvent(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	params := mux.Vars(request)
	var event *pb.Event
	id := params["id"]
	result := dbClient.Table(tableName).Where("id = ?", id).Delete(&event)
	if result.RowsAffected == 0 {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte(`Unable to delete Please Verify`))
		return
	}
	logger.Infof("Deleted Event with Id:" + event.Id)
	json.NewEncoder(response).Encode(result.RowsAffected)
}
