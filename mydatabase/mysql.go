package mydatabase

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func Mysql() {
	db, err := sql.Open("mysql", "root:mysj113598@/mytest?charset=utf8")
	CheckErr(err)
	defer db.Close()

	//插入
	stmt, err := db.Prepare("INSERT INTO userinfo SET username=?,department=?,created=?")
	CheckErr(err)
	res, err := stmt.Exec("Louis", "休息部门", time.Now())
	CheckErr(err)
	id, err := res.LastInsertId()
	CheckErr(err)
	fmt.Println(id)

	//更新
	stmt, err = db.Prepare("update userinfo set username=? where uid=?")
	CheckErr(err)
	res, err = stmt.Exec("louisupdate", id)
	CheckErr(err)
	affect, err := res.RowsAffected()
	CheckErr(err)
	fmt.Println(affect)

	//查询
	rows, err := db.Query("select * from userinfo")
	CheckErr(err)
	for rows.Next() {
		var uid int
		var username string
		var department string
		var created string
		err = rows.Scan(&uid, &username, &department, &created)
		CheckErr(err)
		fmt.Println(uid, username, department, created)
	}

	//删除
	stmt, err = db.Prepare("delete from userinfo where uid=?")
	CheckErr(err)
	res, err = stmt.Exec(id)
	CheckErr(err)
	affect, err = res.RowsAffected()
	CheckErr(err)
	fmt.Println(affect)
}
