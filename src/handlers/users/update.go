package users

import (
	stdMsg "GO-API-template/src/helpers/stdMessages"
	"GO-API-template/src/models"
	"GO-API-template/src/utils"
	"context"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// UpdateUser update user
// @Summary      update user
// @Description  Update user info
// @Accept       json
// @Produce      json
// @security     BearerAuth
// @param	updateUserData body models.User{} true "data to update, currently only allows to update the fullName field"
// @param	uid path string true "User ID or username"
// @Success      200  {object}  interface{}
// @Failure      401  {object}  interface{}
// @Failure      422  {object}  interface{}
// @Failure      500  {object}  interface{}
// @Router       /users/{uid} [patch]
func UpdateUser(c *fiber.Ctx) error {
	//TODO: add more fields to update (maybe via model.User)

	// user update input
	var uui models.User
	if err := c.BodyParser(&uui); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  "error",
			"message": "Unpocessable Entity, Review your input",
			"data":    err,
		})
	}

	// Identity of the user to modify
	identity := c.Params("uid")
	// get the data of the user we want to modify
	var user models.User
	err := user.Fill(identity, true, true, false)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the requested user with identificator: "+identity,
				err.Error(),
			))
	}
	// Role of the user to be modified
	userRole, err := user.RoleData()
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the users's role with ID: "+user.RoleID.Hex(),
				err.Error(),
			))
	}

	// Token of the editor's user
	token := c.Locals("user").(*jwt.Token)
	editorUID := token.Claims.(jwt.MapClaims)["uid"].(string)

	// Get the editor's data
	var editorUser models.User
	err = editorUser.Fill(editorUID, true, false, false)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the editor's user with ID: "+editorUID,
				err.Error(),
			))
	}
	var editorRole models.Role
	err = editorRole.Fill(editorUser.RoleID.Hex(), true, false)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(
			stdMsg.ErrorDefault(
				"could not find the editor's role with ID: "+editorUser.RoleID.Hex(),
				err.Error(),
			))
	}

	// check if editor is authorised to do the operation
	if user.ID.Hex() != editorUser.ID.Hex() {
		// Parametrized permissons
		if !editorRole.Permissons.UsersAdmin {
			/*
				——————————No Perms?——————————————————
				⠀⣞⢽⢪⢣⢣⢣⢫⡺⡵⣝⡮⣗⢷⢽⢽⢽⣮⡷⡽⣜⣜⢮⢺⣜⢷⢽⢝⡽⣝
				⠸⡸⠜⠕⠕⠁⢁⢇⢏⢽⢺⣪⡳⡝⣎⣏⢯⢞⡿⣟⣷⣳⢯⡷⣽⢽⢯⣳⣫⠇
				⠀⠀⢀⢀⢄⢬⢪⡪⡎⣆⡈⠚⠜⠕⠇⠗⠝⢕⢯⢫⣞⣯⣿⣻⡽⣏⢗⣗⠏⠀
				⠀⠪⡪⡪⣪⢪⢺⢸⢢⢓⢆⢤⢀⠀⠀⠀⠀⠈⢊⢞⡾⣿⡯⣏⢮⠷⠁⠀⠀
				⠀⠀⠀⠈⠊⠆⡃⠕⢕⢇⢇⢇⢇⢇⢏⢎⢎⢆⢄⠀⢑⣽⣿⢝⠲⠉⠀⠀⠀⠀
				⠀⠀⠀⠀⠀⡿⠂⠠⠀⡇⢇⠕⢈⣀⠀⠁⠡⠣⡣⡫⣂⣿⠯⢪⠰⠂⠀⠀⠀⠀
				⠀⠀⠀⠀⡦⡙⡂⢀⢤⢣⠣⡈⣾⡃⠠⠄⠀⡄⢱⣌⣶⢏⢊⠂⠀⠀⠀⠀⠀⠀
				⠀⠀⠀⠀⢝⡲⣜⡮⡏⢎⢌⢂⠙⠢⠐⢀⢘⢵⣽⣿⡿⠁⠁⠀⠀⠀⠀⠀⠀⠀
				⠀⠀⠀⠀⠨⣺⡺⡕⡕⡱⡑⡆⡕⡅⡕⡜⡼⢽⡻⠏⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
				⠀⠀⠀⠀⣼⣳⣫⣾⣵⣗⡵⡱⡡⢣⢑⢕⢜⢕⡝⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
				⠀⠀⠀⣴⣿⣾⣿⣿⣿⡿⡽⡑⢌⠪⡢⡣⣣⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
				⠀⠀⠀⡟⡾⣿⢿⢿⢵⣽⣾⣼⣘⢸⢸⣞⡟⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
				⠀⠀⠀⠀⠁⠇⠡⠩⡫⢿⣝⡻⡮⣒⢽⠋⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀
				——————————————————————————————————————
			*/
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "failed to update the user, your role cant update other users.",
				"data":    nil,
			})
		}

		if editorRole.Level >= userRole.Level { // the lower the number the higher the permissons are
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "failed to update the user, Your lole permissons level is not high enouth to edit this role.",
				"data":    nil,
			})
		}
	}

	// Authenticated & autorized

	/*
	 ———————— Validate fields to be updated, writes changes into user ———————————
	 Anything not explicitly written here to user will not be passed to the db

	 							uui >> user

	*/

	// Special field role requires aditional autorization
	if uui.Role != "" {
		/*
		 The editor can only modify users of a higher level (so, less permissons) than itself
		 and can only set the user to a role with a higher level than the editor
		*/

		// get the original role
		err := user.SetRole()
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(
				stdMsg.ErrorDefault("The specified role does not exist", nil),
			)
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				stdMsg.ErrorDefault("Unknown error while retieving the specified role", nil),
			)
		}
		var role models.Role
		err = role.Fill(user.RoleID.Hex(), true, false)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(
				stdMsg.ErrorDefault(
					"could not find the original role for the user, with ID: "+user.RoleID.Hex(),
					err.Error(),
				))
		}

		// check if editor can edit this user's role by role level, the lower the more permissons
		if role.Level < editorRole.Level {
			return c.Status(fiber.StatusForbidden).JSON(
				stdMsg.ErrorDefault("You cant set the the role of this user UwU", nil),
			)
		}
		// check if editor can edit this user's role by role permissons
		if editorRole.Permissons.RolesAdmin && editorRole.Permissons.UsersAdmin {
			user.Role = uui.Role
		} else {
			return c.Status(fiber.StatusForbidden).JSON(
				stdMsg.ErrorDefault("failed to update the user, your role cant modify users roles.", nil),
			)
		}

		// check if the specified role exists and set it
		err = user.SetRole()
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusBadRequest).JSON(
				stdMsg.ErrorDefault("The specified role does not exist", nil),
			)
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(
				stdMsg.ErrorDefault("Unknown error while retieving the specified role", nil),
			)
		}
		role.Fill(user.RoleID.Hex(), true, false)
		// check if editor can set the specified role
		if role.Level <= editorRole.Level {
			return c.Status(fiber.StatusForbidden).JSON(
				stdMsg.ErrorDefault("You cant set the the role of this user", nil),
			)
		}

	}

	// user and email validation
	if uui.Username != "" {
		user.Username = uui.Username
	}
	if uui.Email != "" || uui.Email != user.Email {
		user.Email = uui.Email
	}
	unique, err := user.CheckUnique()
	if err != nil {
		return c.Status(500).JSON(stdMsg.ErrorDefault(
			"An error ocured while checking if username and email are unique",
			err,
		))
	}
	if !unique {
		return c.Status(fiber.StatusConflict).JSON(stdMsg.ErrorDefault(
			"Specified username or email is already being used",
			err,
		))
	}

	// password validation
	if uui.Password != "" {
		hash, err := utils.Hash(uui.Password)
		if err != nil {
			return c.Status(500).JSON(stdMsg.ErrorDefault("Couldn't hash password", err))

		}
		user.Password = hash
	}

	if uui.Name != "" || user.Name != uui.Name {
		user.Name = uui.Name
	}

	if uui.BlockedTokens != nil {
		user.BlockedTokens = uui.BlockedTokens
	}
	if uui.Tokens != nil {
		user.Tokens = uui.Tokens
	}

	filter := bson.M{"_id": user.ID}
	update := bson.M{"$set": user}

	_, err = models.UsersCollection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status":  "error",
			"message": "failed to delete the user",
			"data":    nil,
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "User successfully updated",
		"data":    user.Private(),
	})
}
