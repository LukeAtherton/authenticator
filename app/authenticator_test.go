// Copyright (c) Luke Atherton 2015

package authenticator_test

import (
	. "github.com/lukeatherton/authenticator/app"

	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
	"strings"
)

/*
Convert JSON data into a slice.
*/
func sliceFromJSON(data []byte) []interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.([]interface{})
}

/*
Convert JSON data into a map.
*/
func mapFromJSON(data []byte) map[string]interface{} {
	var result interface{}
	json.Unmarshal(data, &result)
	return result.(map[string]interface{})
}

type TestPublisher struct {
	messages []DomainEvent
}

func (publisher *TestPublisher) PublishMessage(message DomainEvent) (err error) {
	publisher.messages = append(publisher.messages, message)
	return nil
}

/*
Server unit tests.
*/
var _ = Describe("Auth Server", func() {
	var repo Repo
	var server *gin.Engine
	var request *http.Request
	var recorder *httptest.ResponseRecorder

	var testPublisher *TestPublisher
	var testAuth Authenticator
	hostString := []string{}
	hostString = append(hostString, "127.0.0.1")
	//dbHosts []string, authDb string, dbUsername string, dbPassword string
	repo = NewMongoTestRepo(hostString, "admin", "", "")

	BeforeEach(func() {
		// Set up a new server, connected to a test database,
		// before each test.
		testPublisher = &TestPublisher{}
		testAuth = BuildAuthenticator(repo, "../crypto/testKey.pem", "../crypto/testKey.pub")
		server = NewRouter(testPublisher, repo, testAuth)

		// Record HTTP responses.
		recorder = httptest.NewRecorder()
	})

	AfterEach(func() {
		// Clear the database after each test.
		//session.DB(dbName).DropDatabase()
		repo.Cleanup()
	})

	Describe("GET /status", func() {

		// Set up a new GET request before every test
		// in this describe block.
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/status", nil)
		})

		Context("when service is running", func() {
			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(200))
			})
		})
	})

	Describe("GET /emails?email=latherton@example.com", func() {

		var credentials *Credentials

		// Set up a new GET request before every test
		// in this describe block.
		BeforeEach(func() {
			request, _ = http.NewRequest("GET", "/api/emails?email=latherton@example.com", nil)
		})

		Context("when no users exist", func() {

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(200))
			})

			/*It("returns a empty body", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Body.String()).To(Equal("[]"))
				//Expect(recorder.Body.String()).To(Equal(`"Email-Available"`))
			})*/
		})

		Context("when email in use", func() {

			// Insert a valid company uri
			// before each test in this context.
			BeforeEach(func() {
				regView := gory.Build("userRegistrationDuplicate").(*UserRegistrationView)
				credentials, _ = DecodeRegistrationDetails(regView)
				credentials.Id = NewUUID()

				repo.SaveCredentials(credentials.Id, credentials)
			})

			It("returns a status code of 302", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(302))
			})

			/*It("returns the user info", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				responseJSON := mapFromJSON(recorder.Body.Bytes())
				Expect(responseJSON["user_id"]).ToNot(Equal(""))
			})*/

		})
	})

	Describe("GET /verification?email=latherton@example.com&code=", func() {

		var credentials *Credentials

		Context("when no users exist", func() {

			// Set up a new GET request before every test
			// in this describe block.
			BeforeEach(func() {
				request, _ = http.NewRequest("GET", "/api/verification?email=latherton@example.com&code=123", nil)
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(400))
			})

			/*It("returns a empty body", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Body.String()).To(Equal("[]"))
				//Expect(recorder.Body.String()).To(Equal(`"Email-Available"`))
			})*/
		})

		Context("when user not yet verified", func() {

			// Insert a valid company uri
			// before each test in this context.
			BeforeEach(func() {
				regView := gory.Build("userRegistrationDuplicate").(*UserRegistrationView)
				credentials, _ = DecodeRegistrationDetails(regView)
				credentials.Id = NewUUID()
				credentials.EmailVerificationCode = NewUUID().String()

				repo.SaveCredentials(credentials.Id, credentials)

				request, _ = http.NewRequest("GET", fmt.Sprintf("/api/verification?email=%s&code=%s", credentials.Email, credentials.EmailVerificationCode), nil)
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(200))
			})

			It("updates repo", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				savedCredentials, _ := repo.GetCredentials(credentials.Id)

				Expect(savedCredentials.IsEmailVerified).To(BeTrue())
			})

			/*It("returns the user info", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				responseJSON := mapFromJSON(recorder.Body.Bytes())
				Expect(responseJSON["user_id"]).ToNot(Equal(""))
			})*/

		})
	})

	Describe("POST /registrations", func() {

		Context("with invalid JSON", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("userRegistrationMissingEmail"))
				request, _ = http.NewRequest(
					"POST", "/api/registrations", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(400))
			})
		})

		Context("with valid JSON", func() {

			// Create a POST request with valid JSON from
			// our factory before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("userRegistration"))
				request, _ = http.NewRequest(
					"POST", "/api/registrations", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 201", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(201))
			})

			It("publishes user registered events", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				Eventually(func() []DomainEvent {
					return testPublisher.messages
				}).Should(HaveLen(2))

				for _, m := range testPublisher.messages {
					switch m.GetMessageType() {
					case "User.Registered":
						var message *UserRegistered
						message, _ = m.(*UserRegistered)
						Expect(message.Email == "latherton0@example.com")

					case "Email.Verification.Pending":
						var message *EmailVerificationPending
						message, _ = m.(*EmailVerificationPending)
						Expect(message.Email == "latherton0@example.com")
					}
				}

			})

			It("returns the user info", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				responseJSON := mapFromJSON(recorder.Body.Bytes())

				token := strings.Replace(responseJSON["token"].(string), "\"", "", -1)
				token = strings.Trim(token, "\r\n")

				Expect(testAuth.ValidateToken(token)).To(Equal(true))

				// token := strings.Replace(recorder.Body.String(), "\"", "", -1)
				// token = strings.Trim(token, "\r\n")

				// Expect(testAuth.ValidateToken(token)).To(Equal(true))

				// responseArray := strings.Split(recorder.Body.String(), ".")
				// fmt.Printf("TOKEN: %s\n", responseArray[1])

				// userInfo := responseArray[1]

				// //this fixes an issue decoding part of the token
				// if l := len(userInfo) % 4; l > 0 {
				// 	userInfo += strings.Repeat("=", 4-l)
				// }

				// response, err := base64.URLEncoding.DecodeString(userInfo)
				// if err != nil {
				// 	panic(err)
				// 	return
				// }

				// responseJSON := mapFromJSON(response)
				Expect(responseJSON["id"]).ToNot(Equal(""))
			})

			// Measure("new users should be created in under 500ms", func(b Benchmarker) {
			// 	runtime := b.Time("runtime", func() {
			// 		server.ServeHTTP(recorder, request)
			// 		Expect(recorder.Code).To(Equal(201))
			// 	})

			// 	Ω(runtime.Seconds()).Should(BeNumerically("<", 0.5), "creating users shouldn't take longer than 500ms.")

			// 	//b.RecordValue("disk usage (in MB)", HowMuchDiskSpaceDidYouUse())
			// }, 10)
		})

		Context("with JSON containing a duplicate email", func() {

			var credentials *Credentials

			BeforeEach(func() {
				regView := gory.Build("userRegistrationDuplicate").(*UserRegistrationView)
				credentials, _ = DecodeRegistrationDetails(regView)
				credentials.Id = NewUUID()

				repo.SaveCredentials(credentials.Id, credentials)

				body, _ := json.Marshal(
					gory.Build("userRegistrationDupEmail"))
				request, _ = http.NewRequest(
					"POST", "/api/registrations", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(400))
			})
		})
	})

	Describe("POST /auth", func() {

		Context("with invalid JSON", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("loginValid"))
				request, _ = http.NewRequest(
					"POST", "/api/auth", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 401", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(401))
			})
		})

		Context("with existing user JSON", func() {
			var credentials *Credentials

			// Create a POST request with valid JSON from
			// our factory before each test in this context.
			BeforeEach(func() {
				regView := gory.Build("userRegistrationDuplicate").(*UserRegistrationView)
				credentials, _ = DecodeRegistrationDetails(regView)
				credentials.Id = NewUUID()
				credentials.IsEmailVerified = true

				repo.SaveCredentials(credentials.Id, credentials)

				body, _ := json.Marshal(
					gory.Build("loginValid"))
				request, _ = http.NewRequest(
					"POST", "/api/auth", bytes.NewReader(body))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 200", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(200))
			})

			It("returns the user info", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				responseJSON := mapFromJSON(recorder.Body.Bytes())

				token := strings.Replace(responseJSON["token"].(string), "\"", "", -1)
				token = strings.Trim(token, "\r\n")

				Expect(testAuth.ValidateToken(token)).To(Equal(true))

				// token := strings.Split(responseJSON["token"], ".")

				// userInfo := token[1]

				// //this fixes an issue decoding part of the token
				// if l := len(userInfo) % 4; l > 0 {
				// 	userInfo += strings.Repeat("=", 4-l)
				// }

				// response, err := base64.URLEncoding.DecodeString(userInfo)
				// if err != nil {
				// 	panic(err)
				// 	return
				// }

				Expect(responseJSON["id"]).To(Equal(credentials.Id.String()))
			})

			// Measure("authentication should take less than 400ms", func(b Benchmarker) {
			// 	runtime := b.Time("runtime", func() {
			// 		server.ServeHTTP(recorder, request)
			// 		Expect(recorder.Code).To(Equal(200))
			// 	})

			// 	Ω(runtime.Seconds()).Should(BeNumerically("<", 0.4), "authentication shouldn't take longer than 400ms.")

			// 	//b.RecordValue("disk usage (in MB)", HowMuchDiskSpaceDidYouUse())
			// }, 10)
		})

	})

	Describe("POST /credentials/updaterequests", func() {
		var credentials *Credentials
		var token string

		BeforeEach(func() {
			regView := gory.Build("userRegistration").(*UserRegistrationView)
			credentials, _ = DecodeRegistrationDetails(regView)
			credentials.Id = NewUUID()
			credentials.IsEmailVerified = true

			repo.SaveCredentials(credentials.Id, credentials)

			token, _ = testAuth.Authenticate(regView.Email, regView.Password, "")
		})

		Context("with invalid JSON", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("passwordChangeRequestMissingInfo"))
				request, _ = http.NewRequest(
					"POST", "/api/credentials/updaterequests", bytes.NewReader(body))
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 400", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(400))
			})
		})

		Context("with invalid password", func() {

			// Create a POST request using JSON from our invalid
			// factory object before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("passwordChangeRequestWrongPassword"))
				request, _ = http.NewRequest(
					"POST", "/api/credentials/updaterequests", bytes.NewReader(body))
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 401", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(401))
			})
		})

		Context("with valid JSON", func() {

			// Create a POST request with valid JSON from
			// our factory before each test in this context.
			BeforeEach(func() {
				body, _ := json.Marshal(
					gory.Build("passwordChangeRequest"))
				request, _ = http.NewRequest(
					"POST", "/api/credentials/updaterequests", bytes.NewReader(body))
				request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
				request.Header.Set("content-type", "application/json")
			})

			It("returns a status code of 201", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)
				Expect(recorder.Code).To(Equal(201))
			})

			It("returns the user info", func() {
				server.ServeHTTP(recorder, request)
				fmt.Printf("%v\n", recorder)

				responseJSON := mapFromJSON(recorder.Body.Bytes())

				token := strings.Replace(responseJSON["token"].(string), "\"", "", -1)
				token = strings.Trim(token, "\r\n")

				Expect(testAuth.ValidateToken(token)).To(Equal(true))

				// token := strings.Replace(recorder.Body.String(), "\"", "", -1)
				// token = strings.Trim(token, "\r\n")

				// Expect(testAuth.ValidateToken(token)).To(Equal(true))

				// responseArray := strings.Split(recorder.Body.String(), ".")

				// userInfo := responseArray[1]

				// //this fixes an issue decoding part of the token
				// if l := len(userInfo) % 4; l > 0 {
				// 	userInfo += strings.Repeat("=", 4-l)
				// }

				// response, err := base64.URLEncoding.DecodeString(userInfo)
				// if err != nil {
				// 	panic(err)
				// 	return
				// }

				// responseJSON := mapFromJSON(response)
				Expect(responseJSON["id"]).ToNot(Equal(""))
			})
		})
	})
})
