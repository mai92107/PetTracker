package account

func (account *AccountServiceImpl)Register(ip, username, password, email, lastName, firstName, nickName string, deviceName string)(map[string]interface{},error){

	accountUuid,err := account.accountRepo.CreateAccount(username,password,email)
	if err != nil{
		return nil, err
	}

	memberUuid,err := account.memberRepo.CreateMember(accountUuid,lastName,firstName,nickName,email)
	if err != nil{
		return nil, err
	}

	err = account.passwordRepo.CreateHistory(accountUuid,password)
	if err != nil{
		return nil, err
	}

	deviceId,err := account.deviceRepo.CreateDevice(deviceName,memberUuid)
	if err != nil {
		return nil, err
	}
	return account.Login(ip,deviceId,username,password)
}