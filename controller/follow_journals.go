package controller

import (
	"PaperSubmission/enum"
	"PaperSubmission/model"
	"PaperSubmission/service"
	"PaperSubmission/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type JournalListResponse struct {
	List []Journal `json:"list"`
	utils.Response
}

type UserListResponse struct {
	List []User `json:"list"`
	utils.Response
}

type User struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Organization string `json:"organization"`
}

func FollowJournalHandler(context *gin.Context) {
	userID, _ := context.MustGet("userID").(int64)
	journalIDStr := context.Query("journal_id")
	journalID, err := strconv.ParseInt(journalIDStr, 10, 64)
	if err != nil { // 请求参数中journalID不合法
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(int(enum.OperateFail), enum.OperateFail.String()))
		return
	}
	followJournal := model.FollowJournal{UserID: userID, JournalID: journalID}
	if err := service.NewFollowJournalService().Add(followJournal); err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(int(enum.OperateFail), err.Error()))
		return
	}
	context.JSON(http.StatusOK, utils.NewCommonResponse(int(enum.OperateOK), enum.OperateOK.String()))
}

func GetUserFollowingJournalListHandler(context *gin.Context) {
	userID, _ := context.MustGet("userID").(int64)
	pageStr := context.Query("page")
	pageSizeStr := context.Query("page_size")
	page, err1 := strconv.Atoi(pageStr)
	pageSize, err2 := strconv.Atoi(pageSizeStr)
	if err1 != nil || err2 != nil { // 请求参数无法被解析
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(int(enum.OperateFail), enum.OperateFail.String()))
		return
	}
	request := utils.ListQuery{Page: page, PageSize: pageSize}
	journals, err := service.NewFollowJournalService().GetJournalList(userID, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(int(enum.OperateFail), err.Error()))
		return
	}
	var journalList []Journal
	for _, journalModel := range journals {
		journal := Journal{
			Abbreviation: journalModel.Abbreviation,
			CCFRanking:   journalModel.CCFRanking,
			Deadline:     journalModel.Deadline,
			Description:  journalModel.Description,
			FullName:     journalModel.FullName,
			ID:           journalModel.ID,
			ImpactFactor: journalModel.ImpactFactor,
			ISSN:         journalModel.ISSN,
			Publisher:    journalModel.Publisher,
		}
		journalList = append(journalList, journal)
	}
	context.JSON(http.StatusOK, JournalListResponse{
		List:     journalList,
		Response: utils.NewCommonResponse(int(enum.OperateOK), enum.OperateOK.String()),
	})
}

func GetJournalFollowedUserListHandler(context *gin.Context) {
	journalIDStr := context.Query("journal_id")
	pageStr := context.Query("page")
	pageSizeStr := context.Query("page_size")
	page, err1 := strconv.Atoi(pageStr)
	pageSize, err2 := strconv.Atoi(pageSizeStr)
	if err1 != nil || err2 != nil { // 请求参数无法被解析
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(int(enum.OperateFail), enum.OperateFail.String()))
		return
	}
	request := utils.ListQuery{Page: page, PageSize: pageSize}
	journalID, err := strconv.ParseInt(journalIDStr, 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, utils.NewCommonResponse(int(enum.OperateFail), err.Error()))
		return
	}
	users, err := service.NewFollowJournalService().GetUserList(journalID, request)
	if err != nil {
		context.JSON(http.StatusInternalServerError, utils.NewCommonResponse(int(enum.OperateFail), err.Error()))
		return
	}
	var userList []User
	for _, userModel := range users {
		user := User{
			ID:           userModel.ID,
			Name:         userModel.Name,
			Organization: userModel.Organization,
		}
		userList = append(userList, user)
	}
	context.JSON(http.StatusOK, UserListResponse{
		List:     userList,
		Response: utils.NewCommonResponse(int(enum.OperateOK), enum.OperateOK.String()),
	})
}
