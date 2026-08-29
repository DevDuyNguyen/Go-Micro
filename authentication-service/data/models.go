package data

import (
	"authentication/internal"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type Models struct{
	User User
}

func NewModels(dbPool *sql.DB) (Models){
	db=dbPool
	return Models{
		User:User{},
	}
}

type User struct{
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"-"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (*User) GetAll() ([]*User, error){
	ctx, cancel:=internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:=`SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users ORDER BY LastName`

	rows, err:= db.QueryContext(ctx, stmt)
	if err!=nil{
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next(){
		tmpUser:=User{}
		err=rows.Scan(
			&tmpUser.ID,
			&tmpUser.Email,
			&tmpUser.FirstName,
			&tmpUser.LastName,
			&tmpUser.Password,
			&tmpUser.Active,
			&tmpUser.CreatedAt,
			&tmpUser.UpdatedAt,
		)
		if err!=nil{
			return nil,err
		}
		users=append(users, &tmpUser)
	}
	if rows.Err()!=nil{
		return nil, rows.Err()
	}
	return users,nil
}

func (*User) GetByEmail(email string) (*User, error){
	ctx,cancel:=internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:=`SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users WHERE email=$1`

	var tmpUser User
	row:=db.QueryRowContext(ctx, stmt, email)
	err:=row.Scan(
		&tmpUser.ID,
		&tmpUser.Email,
		&tmpUser.FirstName,
		&tmpUser.LastName,
		&tmpUser.Password,
		&tmpUser.Active,
		&tmpUser.CreatedAt,
		&tmpUser.UpdatedAt,
	)
	if err!=nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil,nil
		}else{
			return nil,err
		}
	}
	return &tmpUser, nil
}

func (*User) GetOne(id int) (*User, error){
	ctx,cancel:= internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:=`SELECT id, email, first_name, last_name, password, user_active, created_at, updated_at
	FROM users WHERE id=$1`

	var tmpUser User
	row:=db.QueryRowContext(ctx, stmt, id)
	err:=row.Scan(
		&tmpUser.ID,
		&tmpUser.Email,
		&tmpUser.FirstName,
		&tmpUser.LastName,
		&tmpUser.Password,
		&tmpUser.Active,
		&tmpUser.CreatedAt,
		&tmpUser.UpdatedAt,
	)
	if err!=nil{
		if errors.Is(err, sql.ErrNoRows){
			return nil, nil
		}else{
			return nil, err
		}
	}
	return &tmpUser, nil
}

func (u *User) Update() error{
	ctx, cancel:= internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:=`UPDATE Users
	SET email=$1, first_name=$2, last_name=$3, 
		active=$4, updated_at=$5
	WHERE id=$6`

	_,err:=db.ExecContext(ctx, stmt, u.Email, u.FirstName, 
		u.LastName, u.Active, time.Now(), u.ID)
	if err!=nil{
		return err
	}
	return nil
}

func (u *User) Delete() error{
	ctx, cancel:=internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:="DELETE FROM users where id=$1"
	_,err:=db.ExecContext(ctx, stmt, u.ID)
	if err!=nil{
		return err
	}
	return nil
}

func (*User) DeleteByID(id int) error{
	ctx, cancel:=internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:=`DELETE FROM users where id=$1`
	_,err:=db.ExecContext(ctx, stmt, id)
	if err!=nil{
		return err
	}
	return nil
}

func (*User) Insert(user User) (int, error){
	ctx, cancel:=internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	stmt:=`INSERT INTO usesrs(email, first_name, last_name, password, user_active, created_at)
	VALUES ($1,$2,$3,$4,$5,$6)
	RETURNING id`

	var id int
	err:=db.QueryRowContext(ctx, stmt,
		user.Email, 
		user.FirstName, 
		user.LastName, 
		user.Password,
		user.Active, 
		time.Now(),
	).Scan(&id)
	if err!=nil{
		return 0,nil
	}
	return id,nil
}

func (u *User) ResetPassword(password string) error{
	ctx, cancel:=internal.GetContextWithTimeOut(3*time.Second)
	defer cancel()

	hashedPassword, err:=bcrypt.GenerateFromPassword([]byte(password), 13)
	if err!=nil{
		return err
	}
	
	stmt:="UPDATE users SET password=$1 where id=$2"
	_,err=db.ExecContext(ctx, stmt, hashedPassword, u.ID)
	if err!=nil{
		return err
	}
	return nil	
}

func (u *User) PasswordMatches(plainText string) (bool, error){
	var hashedPassword=u.Password
	if hashedPassword==""{
		ctx, cancel:= internal.GetContextWithTimeOut(3*time.Second)
		defer cancel()

		getHashedPassword:=`SELECT password FROM users WHERE id=$1`
		err:=db.QueryRowContext(ctx, getHashedPassword, u.ID).Scan(&hashedPassword)
		if err!=nil{
			if errors.Is(err, sql.ErrNoRows){
				return false, ErrEntityNotFound
			}else{
				return false, err
			}
		}
	}

	err:=bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainText))
	if err!=nil{
		return true, nil
	} else if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword){
		return false, nil
	} else{
		return false, err
	}

}