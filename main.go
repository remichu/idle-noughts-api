package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"fmt"
	"encoding/json"
	"io/ioutil"
)

type profile struct {
	Tics string `json:"tics"`
	BrowserId string `json:"browserId"`
	Username string `json:"username"`
	Toes string `json:"toes"`
	Tacs string `json:"tacs"`
}


func main() {
allowList := map[string]bool{
	"http://idlenoughts.tk":  true,
    "https://aeolus-1.github.io":  true,
	"https://idlenoughts.tk":  true,
}


	url := os.Getenv("url")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(url))
	if err != nil {
		panic(err)
	}
	col := client.Database("Cluster1").Collection("idle-nought")
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.GET("/get/leaderboard", func(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")

		cur, err := col.Find(context.TODO(), bson.D{{}})
		if err != nil {
			panic(err)
		}
		var lbs []primitive.M

		for cur.Next(context.TODO()) {
			var lb bson.M
			err := cur.Decode(&lb)
			if err != nil {
				panic(err)
			}
			lbs = append(lbs, lb)
		}
		defer cur.Close(context.TODO())

		c.JSON(200, lbs)
	})

	r.POST("/post/update", func(c *gin.Context) {
		if origin := c.Request.Header.Get("Origin"); allowList[origin] {
        c.Header("Access-Control-Allow-Origin", origin)
    }
	var updateLb profile
	jsonData, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
    	panic(err)
	}
	if err := json.Unmarshal(jsonData, &updateLb); err != nil {
        c.JSON(200, gin.H{})
    }
	fmt.Println("Id: ", updateLb.BrowserId)
	fmt.Println("Tics: ", updateLb.Tics)
	fmt.Println("Username: ", updateLb.Username)
	
	filter :=  bson.D{
			{
				Key: "browserId",
				Value: updateLb.BrowserId,
			},
	}
count, err := col.CountDocuments(context.TODO(), filter)

if count >= 1 {
    fmt.Println("Documents exist in this collection!")
	update := bson.D{{"$set",
        bson.D{
            {"tics", updateLb.Tics},
			{"browserId", updateLb.BrowserId},
			{"username", updateLb.Username},
			{"tacs", updateLb.Tacs},
			{"toes", updateLb.Toes},
        },
    }}
	
result, errr := col.UpdateOne(context.TODO(), filter, update)
if errr != nil {
	panic(errr)
}
	fmt.Println(result)
} else {
	res, err := col.InsertOne(context.TODO(), bson.D{
     {"tics", updateLb.Tics},
	{"browserId", updateLb.BrowserId},
	{"username", updateLb.Username},
	{"tacs", updateLb.Tacs},
	{"toes", updateLb.Toes},
})
if err != nil {
	panic(err)
}
fmt.Println(res)
}


	})
	r.Run(":8080")
}
