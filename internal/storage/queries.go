package storage

const (
	checkUserRegister  = "SELECT id FROM users WHERE email = $1"
	registrationInsert = "INSERT INTO users_hash(email,pass_hash,code) VALUES ($1,$2,$3) returning id;"
	loginSelect        = "SELECT * FROM users WHERE email = $1"
	appSelect          = "SELECT id,secret FROM app WHERE id = $1"
	updateChangePass   = `UPDATE users
	SET pass_hash = (SELECT pass_hash FROM users_hash WHERE email = $1)
	WHERE email = $1;
	`
	updateChangePassDeleteHash = `DELETE FROM users_hash WHERE email = $1`
	acceptCodeSelect           = "SELECT code FROM users_hash WHERE email = $1"
	relocateUser               = `INSERT INTO users (id, email, pass_hash)
	SELECT id, email, pass_hash FROM users_hash WHERE email = $1;`
	deleteUserHashRelocate = `DELETE FROM users_hash WHERE email = $1`
	updateUsercode         = "UPDATE users_hash SET code = $1 WHERE email = $2"
	selectIdAdmin
)
