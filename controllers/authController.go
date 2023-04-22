package controllers

import (
	"github.com/CSC4990-Project/CSC4990BackEnd/database"
	"github.com/CSC4990-Project/CSC4990BackEnd/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
)

const SecretKey = "secret"

func Register(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
	uType, _ := strconv.ParseUint(data["type"], 10, 64)
	User := models.User{
		Type:     int(uType),
		Email:    data["email"],
		Password: password,
	}
	stmt, err := database.DB.Prepare("INSERT INTO user (email,password,type) VALUES (?,?,?)")
	if err != nil {
		panic(err)
	}
	_, err = stmt.Exec(User.Email, User.Password, User.Type)
	if err != nil {
		panic(err)
	}
	return c.JSON(User)
}

func Login(c *fiber.Ctx) error {
	var data map[string]string

	if err := c.BodyParser(&data); err != nil {
		return err
	}

	var user models.User

	err := database.DB.QueryRow("SELECT email, password FROM user WHERE email =?", data["email"]).Scan(&user.Email, &user.Password)

	if err != nil {
		c.Status(fiber.StatusNotFound)
		return c.JSON(fiber.Map{
			"message": "user not found",
		})
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{
			"message": "incorrect password",
		})
	}
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    user.Email,
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})
	token, err := claims.SignedString([]byte("secret"))
	if err != nil {
		c.Status(fiber.StatusInternalServerError)
		return c.JSON(fiber.Map{
			"message": "could not login",
		})
	}
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func User(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	claims := token.Claims.(*jwt.StandardClaims)
	var user models.EmailType
	database.DB.QueryRow("SELECT user.email,usertype.type from user,usertype WHERE email =? AND usertype.id = user.type", claims.Issuer).Scan(&user.Email, &user.UserType)
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}

func AdminTicketView(c *fiber.Ctx) error {
	var tickets []models.Ticket
	result, err := database.DB.Query("SELECT t.id,b.buildingName,c.type,p.type,r.roomNumber,t.timeSubmitted from ticket t, building b, " +
		"category c, progress p, roomnumber r WHERE b.id = t.building AND t.category=c.id AND t.progress=p.id AND t.roomNumber=r.id")
	if err != nil {
		return c.JSON(err)
	}
	for result.Next() {
		var ticket models.Ticket
		result.Scan(&ticket.Id, &ticket.Building, &ticket.Category, &ticket.Progress, &ticket.RoomNum, &ticket.TimeSubmit)
		tickets = append(tickets, ticket)
	}
	return c.JSON(tickets)
}
func UserTicketView(c *fiber.Ctx) error {
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
	var tickets []models.Ticket
	claims := token.Claims.(*jwt.StandardClaims)
	result, err := database.DB.Query("SELECT t.id,b.buildingName,c.type,p.type,r.roomNumber,t.timeSubmitted, i.issue from ticket t, building b, "+
		"category c, progress p, roomnumber r, issue i WHERE t.user = ? AND b.id = t.building AND t.category=c.id AND t.progress=p.id AND t.roomNumber=r.id AND t.issue = i.id", claims.Issuer)
	if err != nil {
		return c.JSON(err)
	}
	for result.Next() {
		var ticket models.Ticket
		result.Scan(&ticket.Id, &ticket.Building, &ticket.Category, &ticket.Progress, &ticket.RoomNum, &ticket.TimeSubmit, &ticket.Issue)
		tickets = append(tickets, ticket)
	}
	return c.JSON(tickets)
}

func DetailedView(c *fiber.Ctx) error {
	id := c.Params("id")
	result, err := database.DB.Query("SELECT t.id,s.severity,t.user,b.buildingName,c.type,p.type,r.roomNumber,i.issue,t.timeSubmitted,t.image,t.userComments,t.internalComments,t.timeFinished from ticket t, building b, "+
		"category c, progress p, roomnumber r,severity s, issue i WHERE t.id = ? AND b.id = t.building AND t.category=c.id AND t.progress=p.id AND t.roomNumber=r.id AND t.severity=s.id AND t.issue = i.id", id)
	if err != nil {
		return c.JSON(err)
	}
	var im []byte
	var ticket models.TicketDetails
	for result.Next() {
		result.Scan(&ticket.Id, &ticket.Severity, &ticket.User, &ticket.Building, &ticket.Category, &ticket.Progress, &ticket.RoomNum, &ticket.Issue, &ticket.TimeSubmit, im, &ticket.UserComments, &ticket.InternalComments, &ticket.TimeFinished)
		//ticket.Image = b64.StdEncoding.EncodeToString(im)
	}
	return c.JSON(ticket)

	return c.JSON(ticket)
}

func UpdateTicket(c *fiber.Ctx) error {
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	id := c.Params("id")

	stmt, err := database.DB.Prepare("UPDATE ticket set internalComments = ?,timeFinished =?, severity = ?, progress =? Where id = ?")
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(err)
	}
	_, err = stmt.Exec(data["internalComments"], data["timeUpdated"], data["severity"], data["progress"], id)
	if err != nil {
		return c.JSON(err)
	}
	return c.JSON(data)
}
func SubmitTicket(c *fiber.Ctx) error {
	var data map[string]string
	cookie := c.Cookies("jwt")
	token, err := jwt.ParseWithClaims(cookie, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}

	claims := token.Claims.(*jwt.StandardClaims)
	if err := c.BodyParser(&data); err != nil {
		return err
	}
	subTicket := models.SubmitTicket{Image: data["image"]}
	stmt, err := database.DB.Prepare("INSERT INTO ticket (user,building,category,issue,roomNumber,userComments,image) VALUES(?,?,?,?,?,?,?) ")
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "unable to submit ticket"})
	}
	_, err = stmt.Exec(claims.Issuer, data["building"], data["category"], data["issue"], data["roomNumber"], data["userComments"], subTicket.Image)
	if err != nil {
		c.Status(fiber.StatusBadRequest)
		return c.JSON(fiber.Map{"message": "unable to submit ticket"})
	}
	return c.JSON(data)
}

//building roomNumber image comments user category
