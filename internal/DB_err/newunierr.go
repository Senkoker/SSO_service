package DB_err

import "errors"

var (
	Dublicate_name   = errors.New("Dublicate name")
	Does_not_exist   = errors.New("Does not exist")
	Data_not_equel   = errors.New("Data not equel")
	Invalid_password = errors.New("Invalid password")
	Already_exists   = errors.New("Already exists")
	Invalid_usercode = errors.New("Invalid usercode")
)
