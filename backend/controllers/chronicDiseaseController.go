package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ahmedgulabkhan/ChronicAid/database"
	"github.com/ahmedgulabkhan/ChronicAid/models"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var chronicDiseaseCollection *mongo.Collection = database.OpenCollection(database.Client, "chronic-disease")
var diseases = [11]string{"OSTEOPOROSIS", "STROKE", "DIABETES", "ASTHMA", "CHRONIC-KIDNEY-DISEASE",
	"CYSTIC-FIBROSIS", "SLEEP-APNEA", "OBESITY", "LUPUS", "HASHIMOTO'S-DISEASE", "HIV"}
var criticalityRating = [5]string{"CRITICAL", "BAD", "MODERATE", "GOOD", "IDEAL"}

// Gets the list of all the chronic disease names the user has
func GetUserChronicDiseases() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userName := c.Param("user_name")
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		var chronicDiseaseDetails models.ChronicDisease
		err := chronicDiseaseCollection.FindOne(ctx, bson.M{"user_name": userName}).Decode(&chronicDiseaseDetails)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("The userName does not have any Chronic Diseases.")
				c.JSON(http.StatusOK, []string{})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		var chronicDiseasesList []string

		for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
			chronicDiseasesList = append(chronicDiseasesList, chronicDiseaseDetails.Diseases[i].Disease_name)
		}

		c.JSON(http.StatusOK, chronicDiseasesList)
	}
}

func GetMetricNamesForDisease() gin.HandlerFunc {
	return func(c *gin.Context) {
		diseaseName := c.Param("disease_name")
		diseaseName = strings.ToUpper(diseaseName)

		for i := 0; i < len(diseases); i++ {
			if diseaseName == diseases[i] {
				break
			} else if i == len(diseases)-1 {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Disease name is invalid"})
				return
			}
		}

		var response []string

		if diseaseName == "OSTEOPOROSIS" {
			response = append(response, "Bone Density")
		} else if diseaseName == "STROKE" {
			response = append(response, "Low Density Lipoprotein", "Very Low Density Lipoprotein",
				"High Density Lipoprotein", "Total Cholesterol", "Pulse", "Systolic Blood Pressure", "Diastolic Blood Pressure")
		} else if diseaseName == "DIABETES" {
			response = append(response, "Blood Sugar Level")
		} else if diseaseName == "ASTHMA" {
			response = append(response, "FVC/FEV-1 Ratio", "FeNo Levels")
		} else if diseaseName == "CHRONIC-KIDNEY-DISEASE" {
			response = append(response, "Glomerular Filtration Rate")
		} else if diseaseName == "CYSTIC-FIBROSIS" {
			response = append(response, "Chloride Amount in Sweat")
		} else if diseaseName == "SLEEP-APNEA" {
			response = append(response, "Apnea-Hypopnea Index")
		} else if diseaseName == "OBESITY" {
			response = append(response, "Body Mass Index")
		} else if diseaseName == "LUPUS" {
			response = append(response, "Anti Nuclear Anitbody Test", "Red Blood Cells Count", "White Blood Cells Count",
				"Platelets Count", "C3 Levels", "C4 Levels")
		} else if diseaseName == "HASHIMOTO'S-DISEASE" {
			response = append(response, "Thyroid Peroxidase Antibodies", "Thyroid-Stimulating Hormone",
				"Free T-4")
		} else {
			response = append(response, "Antibody Test", "Nucleic Acid Test")
		}

		c.JSON(http.StatusOK, response)
	}
}

func GetAllUserMetricsForLast7Days() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		response := make(map[string]map[string]string)

		userName := c.Param("user_name")
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		var chronicDiseaseDetails models.ChronicDisease
		err := chronicDiseaseCollection.FindOne(ctx, bson.M{"user_name": userName}).Decode(&chronicDiseaseDetails)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("No chronic disease details found for the user.")
				c.JSON(http.StatusOK, gin.H{"message": "No chronic disease details found for the user."})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		diseases := chronicDiseaseDetails.Diseases
		for i := 0; i < len(diseases); i++ {
			metrics := diseases[i].Metrics
			for j := 0; j < len(metrics); j++ {
				metric_name := metrics[j].Metric_name
				metric_values := metrics[j].Metric_Values

				response[metric_name] = make(map[string]string)
				for key, value := range metric_values {
					epochDay, _ := strconv.Atoi(key)
					if (time.Now().Unix()/86400)-7 <= int64(epochDay) {
						response[metric_name][key] = value
					}
				}
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// Gets User Chronic Disease Metrics for Today based on Epoch Day
func GetUserChronicDiseaseMetricsForToday() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		response := map[string]string{}

		userName := c.Param("user_name")
		diseaseName := c.Param("disease_name")
		diseaseName = strings.ToUpper(diseaseName)
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		chronicDiseaseDetails, isValid := isDiseaseNameValid(c, ctx, userName, diseaseName)
		if !isValid {
			return
		}

		for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
			if chronicDiseaseDetails.Diseases[i].Disease_name == diseaseName {
				metrics := chronicDiseaseDetails.Diseases[i].Metrics
				for j := 0; j < len(metrics); j++ {
					for key, value := range metrics[j].Metric_Values {
						epochDay, _ := strconv.Atoi(key)
						if (time.Now().Unix() / 86400) == int64(epochDay) {
							response[metrics[j].Metric_name] = value
						}
					}
				}
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// Gets User Chronic Disease Metrics for Last 7 Days based on Epoch Day
func GetUserChronicDiseaseMetricsForLast7Days() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		response := make(map[string]map[string]string)

		userName := c.Param("user_name")
		diseaseName := c.Param("disease_name")
		diseaseName = strings.ToUpper(diseaseName)
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		chronicDiseaseDetails, isValid := isDiseaseNameValid(c, ctx, userName, diseaseName)
		if !isValid {
			return
		}

		for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
			if chronicDiseaseDetails.Diseases[i].Disease_name == diseaseName {
				metrics := chronicDiseaseDetails.Diseases[i].Metrics
				for j := 0; j < len(metrics); j++ {
					response[metrics[j].Metric_name] = make(map[string]string)
					for key, value := range metrics[j].Metric_Values {
						epochDay, _ := strconv.Atoi(key)
						if (time.Now().Unix()/86400)-7 <= int64(epochDay) {
							response[metrics[j].Metric_name][key] = value
						}
					}
				}
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// Gets User Chronic Disease Metrics for Last 28 Days based on Epoch Day
func GetUserChronicDiseaseMetricsForLast28Days() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		response := make(map[string]map[string]string)

		userName := c.Param("user_name")
		diseaseName := c.Param("disease_name")
		diseaseName = strings.ToUpper(diseaseName)
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		chronicDiseaseDetails, isValid := isDiseaseNameValid(c, ctx, userName, diseaseName)
		if !isValid {
			return
		}

		for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
			if chronicDiseaseDetails.Diseases[i].Disease_name == diseaseName {
				metrics := chronicDiseaseDetails.Diseases[i].Metrics
				for j := 0; j < len(metrics); j++ {
					response[metrics[j].Metric_name] = make(map[string]string)
					for key, value := range metrics[j].Metric_Values {
						epochDay, _ := strconv.Atoi(key)
						if (time.Now().Unix()/86400)-28 <= int64(epochDay) {
							response[metrics[j].Metric_name][key] = value
						}
					}
				}
			}
		}

		c.JSON(http.StatusOK, response)
	}
}

// Adds a new chronic disease with metrics to the user in the chronic-disease collection. If the document
// for user does not exist, a new document is created and then the chronic disease is added
func AddNewUserChronicDisease() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userName := c.Param("user_name")
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		var disease models.Disease
		var chronicDiseaseDetails models.ChronicDisease

		if err := c.BindJSON(&disease); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(disease)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		disease.Disease_name = strings.ToUpper(disease.Disease_name)
		for i := 0; i < len(diseases); i++ {
			if diseases[i] == disease.Disease_name {
				break
			} else if i == len(diseases)-1 {
				errorMessage := fmt.Sprintf(`The disease that you want to add 
					is currently not supported. Only the diseases: %s are supported currently`, diseases)
				errorMessage = strings.ReplaceAll(errorMessage, "\n", "")
				errorMessage = strings.ReplaceAll(errorMessage, "\t", "")
				c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
				return
			}
		}

		err := chronicDiseaseCollection.FindOne(ctx, bson.M{"user_name": userName}).Decode(&chronicDiseaseDetails)
		if err != nil {
			if err == mongo.ErrNoDocuments {
				fmt.Println("No chronic disease details found for the user.")
				chronicDiseaseDetails.User_name = &userName
				chronicDiseaseDetails.Diseases = append(chronicDiseaseDetails.Diseases, disease)
				chronicDiseaseCollection.InsertOne(ctx, chronicDiseaseDetails)
				c.JSON(http.StatusOK, gin.H{"message": "Chronic Disease added to the user."})
				return
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
			if chronicDiseaseDetails.Diseases[i].Disease_name == disease.Disease_name {
				errorMessage := fmt.Sprintf("The chronic disease: %s is already present for the user with userName: %s", disease.Disease_name, userName)
				c.JSON(http.StatusBadRequest, gin.H{"error": errorMessage})
				return
			}
		}

		var updatedDiseasesList []models.Disease = append(chronicDiseaseDetails.Diseases, disease)

		update := bson.D{{Key: "$set", Value: bson.D{{Key: "diseases", Value: updatedDiseasesList}}}}
		_, err = chronicDiseaseCollection.UpdateOne(ctx, bson.M{"user_name": userName}, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Chronic Disease added to the user."})
	}
}

// Adds a new metric name along with it's values for the given date, if the metrics for the given metric name
// and the given date are already present, this API updates them
func AddNewChronicDiseaseMetric() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var diseaseMetric models.DiseaseMetric
		if err := c.BindJSON(&diseaseMetric); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(diseaseMetric)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		userName := c.Param("user_name")
		diseaseName := c.Param("disease_name")
		diseaseName = strings.ToUpper(diseaseName)
		if _, isValid := isUserIdValid(c, ctx, userName); !isValid {
			return
		}

		chronicDiseaseDetails, isValid := isDiseaseNameValid(c, ctx, userName, diseaseName)
		if !isValid {
			return
		}

		var metricValue string
		for _, val := range diseaseMetric.Metric_Values {
			metricValue = val
		}
		content := fmt.Sprintf(`The result for the metric %s is %s. Now based on this metric and value, give
							me two answers. The first anwer should be where you tell me how good or bad 
							this metric and value is and give some advice to improve this metric value; and 
							in the second answer rate this metric's value on a scale of 1 to 5 where 1 is 
							Critical and 5 is Ideal for a person. Give me the second answer's response in 
							only one one digit from 1 to 5, without adding any more details. Add the label 
							'first' to the first answer and 'second' to the second answer`,
			diseaseMetric.Metric_name, metricValue)
		response, err := callChatGpt(content)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		adviceFromOpenAi := response[7 : len(response)-11]
		ratingFromOpenAi := response[len(response)-1:]
		ratingFromOpenAiInt, _ := strconv.Atoi(ratingFromOpenAi)
		ratingFromOpenAi = criticalityRating[ratingFromOpenAiInt-1]

		for key, val := range diseaseMetric.Metric_Values {
			diseaseMetric.Metric_Values[key] = val + "-" + ratingFromOpenAi
		}

		var diseaseIndex int
		for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
			if chronicDiseaseDetails.Diseases[i].Disease_name == diseaseName {
				diseaseIndex = i
				break
			}
		}

		metrics := chronicDiseaseDetails.Diseases[diseaseIndex].Metrics

		// Length of the metrics is zero, so we need to add a new metric to the empty metrics list
		if len(metrics) == 0 {
			keyToUpdate := fmt.Sprintf("diseases.%d.metrics", diseaseIndex)
			updatedMetrics := append(metrics, diseaseMetric)
			update := bson.D{{Key: "$set", Value: bson.D{{Key: keyToUpdate, Value: updatedMetrics}}}}
			_, err := chronicDiseaseCollection.UpdateOne(ctx, bson.M{"user_name": userName}, update)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		} else {
			for i := 0; i < len(metrics); i++ {
				if metrics[i].Metric_name == diseaseMetric.Metric_name {
					for key, value := range diseaseMetric.Metric_Values {
						keyToUpdate := fmt.Sprintf("diseases.%d.metrics", diseaseIndex)
						metrics[i].Metric_Values[key] = value
						update := bson.D{{Key: "$set", Value: bson.D{{Key: keyToUpdate, Value: metrics}}}}
						_, err := chronicDiseaseCollection.UpdateOne(ctx, bson.M{"user_name": userName}, update)
						if err != nil {
							c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
							return
						}
						break
					}
					break
				} else if i == len(metrics)-1 {
					keyToUpdate := fmt.Sprintf("diseases.%d.metrics", diseaseIndex)
					updatedMetrics := append(metrics, diseaseMetric)
					update := bson.D{{Key: "$set", Value: bson.D{{Key: keyToUpdate, Value: updatedMetrics}}}}
					_, err := chronicDiseaseCollection.UpdateOne(ctx, bson.M{"user_name": userName}, update)
					if err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
						return
					}
					break
				}
			}
		}

		responseMessage := fmt.Sprintf("Successfully added new metric to the chronic disease for the user. General Advice: %s", adviceFromOpenAi)
		c.JSON(http.StatusOK, gin.H{"message": responseMessage})
	}
}

// Send symptoms as text for a particular user_name and using OpenAI API, get back what might be the
// chronic disease out of the supported 11 chronic diseases
func CheckSymptomsForChronicDisease() gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestBody models.OpenAiRequest
		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(requestBody)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		response, err := callChatGpt(requestBody.Message)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": response})
	}
}

func isUserIdValid(c *gin.Context, ctx context.Context, userName string) (models.User, bool) {
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"user_name": userName}).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			errorMessage := "The userName does not exist."
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": errorMessage})
			return models.User{}, false
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return models.User{}, false
		}
	}

	return user, true
}

func isDiseaseNameValid(c *gin.Context, ctx context.Context, userName string, diseaseName string) (models.ChronicDisease, bool) {
	var chronicDiseaseDetails models.ChronicDisease
	err := chronicDiseaseCollection.FindOne(ctx, bson.M{"user_name": userName}).Decode(&chronicDiseaseDetails)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The user does not have this disease."})
			return models.ChronicDisease{}, false
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return models.ChronicDisease{}, false
		}
	}

	for i := 0; i < len(chronicDiseaseDetails.Diseases); i++ {
		if chronicDiseaseDetails.Diseases[i].Disease_name == diseaseName {
			break
		} else if i == len(chronicDiseaseDetails.Diseases)-1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "The user does not have this disease."})
			return models.ChronicDisease{}, false
		}
	}
	return chronicDiseaseDetails, true
}

func callChatGpt(content string) (string, error) {
	body := map[string]interface{}{
		"model": "gpt-4-turbo",
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": content,
			},
		},
	}

	byteArray, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(byteArray))
	if err != nil {
		return "", err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer <key>")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	openAiResponse := &models.OpenAiResponse{}
	decodeErr := json.NewDecoder(res.Body).Decode(openAiResponse)
	if decodeErr != nil {
		return "", nil
	}

	if openAiResponse.Choices == nil {
		errorMessage := fmt.Sprintf("Error from OpenAI API - %s", openAiResponse.Error["message"])
		return "", errors.New(errorMessage)
	}

	return openAiResponse.Choices[0].Message.Content, nil
}
