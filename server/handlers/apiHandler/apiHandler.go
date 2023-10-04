package apiHandler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"peargram/database"
	"peargram/models"
	"peargram/notifications"
	"peargram/users"
	"regexp"
	"strings"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/square/go-jose.v2/jwt"
)

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.Redirect(http.StatusTemporaryRedirect, "/")
}

func GETSearchUsers(c *gin.Context) {
	// Implement this internally ?

	var usernameCoincidences []string
	var userCoincidences []models.UserDetails
	userStr := c.Param("username")

	DB := database.ConnectDB()
	err := DB.Select(&usernameCoincidences, "SELECT username FROM users WHERE username LIKE ?;", "%"+userStr+"%") // TODO: Concatenate with userDetails
	if err != nil {
		fmt.Println(err)
	}

	for _, username := range usernameCoincidences {
		fmt.Println(username)
		userDetails := users.GetUserDetails(username)
		userCoincidences = append(userCoincidences, userDetails)
	}

	c.JSON(http.StatusOK, gin.H{"Amount": len(userCoincidences), "Coincidences": userCoincidences})
}

func GETToggleFollow(c *gin.Context) {
	session := sessions.Default(c)
	currentUsername := session.Get("Username").(string)

	userStr := c.Param("username")

	currentlyFollowing := false
	DB := database.ConnectDB()
	err := DB.QueryRow("SELECT EXISTS (SELECT * FROM follows WHERE actor=? AND target=?)", currentUsername, userStr).Scan(&currentlyFollowing)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if currentlyFollowing {
		_, err := DB.Exec("DELETE FROM follows WHERE actor=? AND target=?", currentUsername, userStr)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	} else {
		date := time.Now().Unix()
		_, err := DB.Exec("INSERT INTO follows (actor, target, date) VALUES (?, ?, ?)", currentUsername, userStr, date)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"NowFollowing": !currentlyFollowing})

	if !currentlyFollowing {
		notifications.SendNotification(models.NOTIFICATION_FOLLOW, currentUsername, userStr, -1)
	}
}

var googleOauthConfigSignIn = &oauth2.Config{
	RedirectURL:  "http://localhost/api/signinEndpoint",
	ClientID:     "871553305028-f0jmpj0ve485brejh0deg92td4pdgool.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}
var googleOauthConfigSignUp = &oauth2.Config{
	RedirectURL:  "http://localhost/api/signupEndpoint",
	ClientID:     "871553305028-f0jmpj0ve485brejh0deg92td4pdgool.apps.googleusercontent.com",
	ClientSecret: "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",
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
	Endpoint:     google.Endpoint,
}

func GETSignin(c *gin.Context) {
	redirectURL := googleOauthConfigSignIn.AuthCodeURL("randomized")
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}
func GETSignup(c *gin.Context) {
	redirectURL := googleOauthConfigSignUp.AuthCodeURL("randomized")
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func POSTSigninEndpoint(c *gin.Context) {
	// fmt.Println(c.Request.ParseForm())
	// fmt.Println(c.PostForm("code"))
	// fmt.Println(c.Request.Form)
	// fmt.Println(c.Request.Form.Get("code"))

	credential := c.PostForm("credential")
	// g_csrf_token := c.PostForm("g_csrf_token")
	userInfo := decodeToken(credential)
	userEmail := strings.ToLower(userInfo["email"].(string))
	// userAvatar := userInfo["picture"]

	DB := database.ConnectDB()

	// Check if user is in pending sign up
	userExists := false
	DB.QueryRow("SELECT EXISTS(SELECT * FROM pendingSignups WHERE lower(email)=?)", userEmail).Scan(&userExists)
	if userExists {
		session := sessions.Default(c)
		session.Clear()
		session.Set("Username", "UNDEFINED")
		session.Set("Email", userEmail)
		session.Save()
		c.Redirect(http.StatusFound, "/auth/endSignup")
		return
	}

	// Check DB
	userExists = false
	DB.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE lower(email)=?)", userEmail).Scan(&userExists)

	if !userExists {
		c.Redirect(http.StatusFound, "/auth/signin?error=notexists")
		return
	}

	username := ""
	err := DB.QueryRow("SELECT username FROM users WHERE lower(email)=?", userEmail).Scan(&username)
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	if username == "" {
		fmt.Println("####### GOT EMPTY USERNAME FROM DB, ERROR")
		c.Status(http.StatusInternalServerError)
		return
	}

	session := sessions.Default(c)
	session.Clear()
	session.Set("Username", username)
	session.Set("Email", userEmail)
	session.Save()
	c.Redirect(http.StatusFound, "/")
}

func POSTSignupEndpoint(c *gin.Context) {
	var userAvatar string
	credential := c.PostForm("credential")

	userInfo := decodeToken(credential)

	userEmail := strings.ToLower(userInfo["email"].(string))
	userAvatarUrl := userInfo["picture"].(string)

	avatarSizePattern := `=s\d+-c`
	regex := regexp.MustCompile(avatarSizePattern)

	userAvatarUrl = regex.ReplaceAllString(userAvatarUrl, "=s200-c")

	response, err := http.Get(userAvatarUrl)
	if err == nil {
		userAvatarBytes, err := ioutil.ReadAll(response.Body)
		defer response.Body.Close()

		if err == nil {
			userAvatar = string(userAvatarBytes)
		}
	}

	// Check DB
	userExists := false
	DB := database.ConnectDB()
	DB.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE lower(email)=?)", userEmail).Scan(&userExists)

	if userExists {
		c.Redirect(http.StatusFound, "/auth/signup?error=exists")
		return
	}

	// Check if user is in pending sign up
	userExists = false
	DB.QueryRow("SELECT EXISTS(SELECT * FROM pendingSignups WHERE lower(googleEmail)=?)", userEmail).Scan(&userExists)

	session := sessions.Default(c)
	session.Clear()

	if !userExists {
		date := time.Now().Unix()
		_, err := DB.Query("INSERT INTO pendingSignups (googleEmail, date, avatar) VALUES(?, ?, ?)", userEmail, date, userAvatar)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		fmt.Printf("####### USER PENDING SIGN UP (%s)\n", userEmail)
	}

	session.Set("Username", "UNDEFINED")
	session.Set("Email", userEmail)
	session.Save()
	c.Redirect(http.StatusFound, "/auth/endSignup")

}

const (
	ERROR_LENGTH_LESS   = 0
	ERROR_LENGTH_MORE   = 1
	ERROR_INVALID_CHAR  = 2
	ERROR_UDSC_PRD_MORE = 3
	ERROR_IN_USE        = 4
	ERROR_UNEXPECTED    = 5
)

type EndSignupBody struct {
	Username string
}

func POSTEndSignup(c *gin.Context) {
	session := sessions.Default(c)
	currentEmail := session.Get("Email")

	pendingSignupExists := false
	DB := database.ConnectDB()
	DB.QueryRow("SELECT EXISTS(SELECT * FROM pendingSignups WHERE lower(email)=?)", currentEmail).Scan(&pendingSignupExists)

	if pendingSignupExists { // TODO: Handle if exists on error side. Log in? Clear session? Apply session?
		c.JSON(http.StatusBadRequest, gin.H{"errorID": ERROR_UNEXPECTED})
		return
	}

	var requestBody EndSignupBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errorID": ERROR_UNEXPECTED})
	}

	desiredUsnm := requestBody.Username

	errorID := ValidateUsername(desiredUsnm)

	if errorID != -1 {
		httpStatus := http.StatusBadRequest
		if errorID == ERROR_UNEXPECTED {
			httpStatus = http.StatusInternalServerError
			return
		}
		c.JSON(httpStatus, gin.H{"errorID": errorID})
		return
	}

	userExists := false
	DB.QueryRow("SELECT EXISTS(SELECT * FROM users WHERE lower(username)=?)", strings.ToLower(desiredUsnm)).Scan(&userExists)

	if userExists {
		c.JSON(http.StatusBadRequest, gin.H{"errorID": ERROR_IN_USE})
		return
	}

	// Then user can be registered:
	// Get his avatar
	userAvatar := ""
	err := DB.QueryRow("SELECT avatar FROM pendingSignups WHERE lower(googleEmail)=?", currentEmail).Scan(&userAvatar)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userAvatar)

	if pendingSignupExists { // TODO: Handle if exists on error side. Log in? Clear session? Apply session?
		c.JSON(http.StatusBadRequest, gin.H{"errorID": ERROR_UNEXPECTED})
		return
	}

	_, err = DB.Exec("DELETE FROM pendingSignups WHERE lower(googleEmail)=?", currentEmail)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"errorID": ERROR_UNEXPECTED})
		return
	}

	date := time.Now().Unix()
	_, err = DB.Exec("INSERT INTO users (username, email, regDate) VALUES (?, ?, ?)", desiredUsnm, currentEmail, date)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"errorID": ERROR_UNEXPECTED})
		return
	}

	_, err = DB.Exec("INSERT INTO userDetails (username, name, description, avatar) VALUES (?, ?, ?, ?)", desiredUsnm, desiredUsnm, "", userAvatar)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"errorID": ERROR_UNEXPECTED})
		return
	}
	// TODO: Must query the 3 or none

	session.Clear()
	session.Set("Username", desiredUsnm)
	session.Set("Email", currentEmail)
	session.Save()

	c.Status(http.StatusOK)
}

func ValidateUsername(usnm string) int {

	// LENGTH RELATED
	if len(usnm) < 4 {
		return ERROR_LENGTH_LESS
	}
	if len(usnm) > 14 {
		return ERROR_LENGTH_MORE
	}

	// INVALID CHAR
	valid, _ := regexp.MatchString("(?i)[a-z0-9._]+$", usnm)
	if !valid {
		return ERROR_INVALID_CHAR
	}

	// UNDERSCORE & PERIOD
	regex, err := regexp.Compile("[._]")
	if err != nil {
		return ERROR_UNEXPECTED
	}
	regMatch := regex.FindAllString(usnm, -1)

	if len(regMatch) > 1 {
		return ERROR_UDSC_PRD_MORE
	}

	return -1
}

func decodeToken(tokenString string) map[string]interface{} {
	var claims map[string]interface{}

	// decode JWT token without verifying the signature
	token, _ := jwt.ParseSigned(tokenString)
	_ = token.UnsafeClaimsWithoutVerification(&claims)

	return claims
}
