package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/agusdwis/rest-api/api/auth"
	"github.com/agusdwis/rest-api/api/models"
	"github.com/agusdwis/rest-api/api/responses"
	"github.com/agusdwis/rest-api/api/utils/formaterror"
	"github.com/gorilla/mux"
)

func (server *Server) LikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check post if exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	// Check user if exist
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("User not found"))
		return
	}

	like := models.Like{}
	like.AuthorID = user.ID
	like.PostID = post.ID

	likeCreated, err := like.SaveLike(server.DB)

	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		responses.ERROR(w, http.StatusInternalServerError, formattedError)
		return
	}
	responses.JSON(w, http.StatusOK, likeCreated)
}

func (server *Server) GetLike(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	//Check if post exist
	post := models.Post{}
	err = server.DB.Debug().Model(models.Post{}).Where("id = ?", pid).Take(&post).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Post not found"))
		return
	}

	like := models.Like{}

	likes, err := like.GetLikesInfo(server.DB, pid)
	if err != nil {
		responses.ERROR(w, http.StatusInternalServerError, err)
		return
	}
	responses.JSON(w, http.StatusOK, likes)
}

func (server *Server) UnlikePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	pid, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}

	uid, err := auth.ExtractTokenID(r)
	if err != nil {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}

	// Check like if exist
	like := models.Like{}
	err = server.DB.Debug().Model(models.Like{}).Where("id = ?", pid).Take(&like).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("Like not found"))
		return
	}

	// Check user if exist
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		responses.ERROR(w, http.StatusNotFound, errors.New("User not found"))
		return
	}

	// Is the authenticated user, the owner of this like?
	if uid != like.AuthorID {
		responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
		return
	}
	_, err = like.DeleteLike(server.DB)
	if err != nil {
		responses.ERROR(w, http.StatusBadRequest, err)
		return
	}
	w.Header().Set("Entity", fmt.Sprintf("%d", pid))
	responses.JSON(w, http.StatusNoContent, "")


}