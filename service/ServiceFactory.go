package service

const (
	RECIEVE_DATA = 0
	SEND_DATA    = 1
)

func GetDataService(flag int) IDataService {
	var ids IDataService
	if flag == RECIEVE_DATA {
		ids = &RecivieveDataService{}
	} else if flag == SEND_DATA {
		ids = &SendDataService{}
	}
	return ids
}
