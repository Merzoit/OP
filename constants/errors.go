package constants

const (
	//CONFIG

	ErrConfigValidate        = "configuration validation failed"
	ErrConfigValidateAppPort = "app port is not specified"
	ErrConfigValidateDbHost  = "database host or name is not specified"
	ErrConfigLoad            = "failed to load configuration"

	//DB

	ErrDbParseConfig = "failed to parse database config"
	ErrDbConnect     = "failed to connect to database"
	ErrDbPing        = "failed to ping database"

	//USER

	ErrUserNotFound     = "user not found"
	ErrUserFetching     = "failed to fetched user"
	ErrUserCreate       = "failed to create user"
	ErrUserAlreadyExist = "user already exist"
	ErrUserBanned       = "failed to ban user"

	//CODE

	ErrCodeDelete           = "failed to delete code"
	ErrCodeCreate           = "failed to create code"
	ErrCodeFetching         = "failed to fetched code"
	ErrCodeRequestCounter   = "failed to request counter"
	ErrCodeNotFound         = "code not found"
	ErrCodeAlreadyExist     = "code already exist"
	ErrCodeCheck            = "code check failed"
	ErrCodeScan             = "code scan failed"
	ErrCodeFetchingByWorker = "failed to get codes by worker"
	ErrCodesIterate         = "failed to iterate codes"

	//ROLE

	ErrRoleNotFound = "role not found"
	ErrRoleFetching = "failed to fetched role"

	//WORKER

	ErrWorkerBalanceUp         = "worker balance up"
	ErrWorkerUpdatePaymentRate = "failed to update payment rate"
	ErrWorkerDelete            = "failed to delete worker"
	ErrWorkerCreate            = "failed to create worker"
	ErrWorkerFetching          = "failed to fetched worker"
	ErrWorkerNotFound          = "worker not found"
	ErrWorkerBalanceReset      = "worker balance reset"

	//SPONSOR

	ErrSponsorFetching = "failed to fetched sponsor"
	ErrSponsorNotFound = "sponsor not found"
	ErrSponsorCreate   = "failed to create sponsor"
	ErrSponsorDelete   = "failed to delete sponsor"
	ErrSponsorsGet     = "failed to get sponsors"
	ErrSponsorScan     = "failed to scan sponsor"
	ErrSponsorsIterate = "failed to iterate sponsor"

	//SUBSCRIBE

	ErrSubscribesBySponsor = "failed to fetching subscribes by sponsor"
	ErrSubscribesByUser    = "failed to fetching subscribes by user"
	ErrSubscribeCreate     = "failed to create subscribe"
	ErrSubscribeScan       = "failed to scan subscribe"
	ErrSubscribeIterate    = "failed to iterate subscribe"

	//REFLINK

	ErrRefRegAdd   = "failed to referral registrations add"
	ErrRefClickAdd = "failed to referall click add"
	ErrRefNotFound = "Referral link not found"
	ErrRefUpdate   = "failed to update referral link"
	ErrRefGet      = "failed to fetch referral link"
	ErrRefCreate   = "failed to create referral link"

	//RESPONSE

	ErrDecodingRequestBody = "error decoding request body"
	ErrEncodingResponse    = "error encoding response"

	//OTHER

	ErrExtractId = "error extract id"
	ErrNoRows    = "no rows in result set"
)
