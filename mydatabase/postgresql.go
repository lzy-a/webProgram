package mydatabase

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func Pq() {
	db, err := sql.Open("postgres", "user=astaxie password=astaxie dbname=test sslmode=disable")
	CheckErr(err)

	//插入数据
	stmt, err := db.Prepare("INSERT INTO userinfo(username,department,created) VALUES($1,$2,$3) RETURNING uid")
	CheckErr(err)

	// res, err := stmt.Exec("astaxie", "研发部门", "2012-12-09")
	// CheckErr(err)

	//pg不支持这个函数，因为他没有类似MySQL的自增ID
	// id, err := res.LastInsertId()
	// CheckErr(err)
	// fmt.Println(id)

	var lastInsertId int
	err = db.QueryRow("INSERT INTO userinfo(username,departname,created) VALUES($1,$2,$3) returning uid;", "astaxie", "研发部门", "2012-12-09").Scan(&lastInsertId)
	CheckErr(err)
	fmt.Println("最后插入id =", lastInsertId)

	//更新数据
	stmt, err = db.Prepare("update userinfo set username=$1 where uid=$2")
	CheckErr(err)

	res, err := stmt.Exec("astaxieupdate", 1)
	CheckErr(err)

	affect, err := res.RowsAffected()
	CheckErr(err)

	fmt.Println(affect)

	//查询数据
	rows, err := db.Query("SELECT * FROM userinfo")
	CheckErr(err)

	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		CheckErr(err)
		fmt.Println(uid)
		fmt.Println(username)
		fmt.Println(department)
		fmt.Println(created)
	}

	//删除数据
	stmt, err = db.Prepare("delete from userinfo where uid=$1")
	CheckErr(err)

	res, err = stmt.Exec(1)
	CheckErr(err)

	affect, err = res.RowsAffected()
	CheckErr(err)

	fmt.Println(affect)

	db.Close()

}
