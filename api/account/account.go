package account

import "batchLog/service/account"


type AccountController struct {
	accountService account.AccountService
}

func NewAccountController(accountService account.AccountService) *AccountController {
	return &AccountController{accountService: accountService}
}