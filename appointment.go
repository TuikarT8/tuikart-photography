package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ctx = context.TODO()

func (appointment *Appointment) saveAppointmentInDb() (string, error) {
	_, err := appointments.InsertOne(ctx, appointment)
	if err != nil {
		log.Printf("error inserting appointment, %v\n", err)
		return "", err
	}
	return "", nil
}

func deleteAppointmentInDb(id string) error {
	bsonId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = appointments.DeleteOne(ctx, bson.M{
		"_id": bsonId,
	})

	if err != nil {
		return err
	}
	return nil
}

func getAppointment(pagination PageQueryParams) ([]Appointment, error) {
	apps := make([]Appointment, 0)

	opts := options.Find().SetSort(bson.D{{"date", -1}})
	opts.Skip = &pagination.startAt
	opts.Limit = &pagination.count
	result, err := appointments.Find(ctx, bson.M{}, opts)

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return apps, nil
		}

		return apps, err
	}

	defer result.Close(ctx)

	err = result.All(ctx, &apps)
	if err != nil {
		return apps, err
	}

	return apps, nil
}

func updateAppointmentInDb(appsId string, apps Appointment) (string, error) {
	hexId, _ := primitive.ObjectIDFromHex(appsId)
	_, err := appointments.UpdateOne(
		ctx,
		primitive.M{"_id": hexId},
		primitive.M{
			"$set": primitive.M{
				"title":       apps.Title,
				"description": apps.Description,
			},
		},
	)
	if err != nil {
		log.Printf("Error while updating order, error=%v", err)
		return "", err
	}

	return "", err
}

func handleCreateAppointment(w http.ResponseWriter, r *http.Request) {
	checkMethod(w, r, http.MethodPost)

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Errors lors de la lecture du corps de la requette"))
		log.Print("Errors lors de la lecture du corps de la requette", err)
		return
	}

	var appointment Appointment

	err = json.Unmarshal(body, &appointment)
	if err != nil {
		handleUnmarshallingError(err.Error(), w)
		return
	}

	_, err = appointment.saveAppointmentInDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Errors lors de la creation d'une appointment"))
		log.Print("Errors lors de la creation d'une appointment", err)
		return
	}

	jsondata, _ := json.Marshal(appointment)
	w.WriteHeader(http.StatusNoContent)
	w.Write(jsondata)
}

func handleDeleteAppointment(w http.ResponseWriter, r *http.Request) {
	checkMethod(w, r, http.MethodDelete)

	appointmentId := mux.Vars(r)["id"]

	err := deleteAppointmentInDb(appointmentId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Errors lors de la suppression d'un appointment"))
		log.Print("Errors lors de la suppression d'un appointment", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleWatchAppointment(w http.ResponseWriter, r *http.Request) {
	checkMethod(w, r, http.MethodGet)

	pageParams, err := pageQueryFromRequestQueryParams(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Errors lors de la recuperations des appointment"))
		log.Print("Errors lors de la  recuperations des appointment", err)
		return
	}

	appointments, err := getAppointment(pageParams)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Errors lors de la recuperations des appointment"))
		log.Print("Errors lors de la  recuperations des appointment", err)
		return
	}

	jsonData, err := json.Marshal(appointments)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Errors lors du Marshaling des appointment"))
		log.Printf("Errors lors du Marshaling des appointment, %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func handleUpdateAppointment(w http.ResponseWriter, r *http.Request) {

	appointmentId := mux.Vars(r)["id"]

	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Errors lors de la lecture du corps de la requette"))
		log.Print("Errors lors de la lecture du corps de la requette", err)
		return
	}

	var appointment Appointment

	err = json.Unmarshal(body, &appointment)
	if err != nil {
		handleUnmarshallingError(err.Error(), w)
		return
	}
	_, err = updateAppointmentInDb(appointmentId, appointment)
	if err != nil {
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func checkMethod(w http.ResponseWriter, r *http.Request, method string) {
	if r.Method != method {
		log.Print("Vous tentez d'utiliser un Endpoint avec la methode Inaproprie", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("L'endpoint invalid"))
		return
	}
}

func handleUnmarshallingError(err string, w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Errors lors de l'unmarshalisation du corps de la requette"))
	log.Print("Errors lors de l'unmarshalisation  du corps de la requette", err)
}
