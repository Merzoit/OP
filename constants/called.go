package constants

const (
	//REPOSITORY
	//user
	CallRepoCreateUser  = "Called [CreateUser] from user tid: %v"
	CallRepoGetUserById = "Called [GetUserById] from user tid: %v"
	//code
	CallRepoDeleteCode = "Called [DeleteCode]: %v"
	//sponsor
	CallRepoCreateSponsor = "Called [CreateSponsor] from sponsor tid: %v"
	CallRepoGetSponsor    = "Called [GetSponsor] from sponsor tid: %v"
	CallRepoDeleteSponsor = "Called [DeleteSponsor] from sponsor tid: %v"
	CallRepoGetSponsors   = "Called [GetSponsors]"

	//CONTROLLER
	CallControllerCreateUser    = "Called [CreateUser] from controller"
	CallControllerGetUserById   = "Called [GetUserById] from controller"
	CallControllerSponsorDelete = "Called [DeleteSponsor] from controller"
	CallControllerSponsorsGet   = "Called [GetSponsors] from controller"
	//OTHER
	CallExtractId = "Called [ExtractId] from controller"
)
