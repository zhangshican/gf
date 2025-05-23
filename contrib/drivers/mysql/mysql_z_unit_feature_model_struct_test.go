// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql_test

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Model_Embedded_Insert(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Id         int    `json:"id"`
			Uid        int    `json:"uid"`
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		result, err := db.Model(table).Data(User{
			Passport: "john-test",
			Password: "123456",
			Nickname: "John",
			Base: Base{
				Id:         100,
				Uid:        100,
				CreateTime: gtime.Now().String(),
			},
		}).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)
		value, err := db.Model(table).Fields("passport").Where("id=100").Value()
		t.AssertNil(err)
		t.Assert(value.String(), "john-test")
	})
}

func Test_Model_Embedded_MapToStruct(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		type Ids struct {
			Id  int `json:"id"`
			Uid int `json:"uid"`
		}
		type Base struct {
			Ids
			CreateTime string `json:"create_time"`
		}
		type User struct {
			Base
			Passport string `json:"passport"`
			Password string `json:"password"`
			Nickname string `json:"nickname"`
		}
		data := g.Map{
			"id":          100,
			"uid":         101,
			"passport":    "t1",
			"password":    "123456",
			"nickname":    "T1",
			"create_time": gtime.Now().String(),
		}
		result, err := db.Model(table).Data(data).Insert()
		t.AssertNil(err)
		n, _ := result.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Model(table).Where("id=100").One()
		t.AssertNil(err)

		user := new(User)

		t.Assert(one.Struct(user), nil)
		t.Assert(user.Id, data["id"])
		t.Assert(user.Passport, data["passport"])
		t.Assert(user.Password, data["password"])
		t.Assert(user.Nickname, data["nickname"])
		t.Assert(user.CreateTime, data["create_time"])
	})
}

func Test_Struct_Pointer_Attribute(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       *int
		Passport *string
		Password *string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).WherePri(1).One()
		t.AssertNil(err)
		user := new(User)
		err = one.Struct(user)
		t.AssertNil(err)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.Model(table).Scan(user, "id=1")
		t.AssertNil(err)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).Scan(&user, "id=1")
		t.AssertNil(err)
		t.Assert(*user.Id, 1)
		t.Assert(*user.Passport, "user_1")
		t.Assert(*user.Password, "pass_1")
		t.Assert(user.Nickname, "name_1")
	})
}

func Test_Structs_Pointer_Attribute(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       *int
		Passport *string
		Password *string
		Nickname string
	}
	// All
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		users := make([]User, 0)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		users := make([]*User, 0)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		one, err := db.Model(table).All("id < 3")
		t.AssertNil(err)
		err = one.Structs(&users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	// Structs
	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		users := make([]*User, 0)
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(table).Scan(&users, "id < 3")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(*users[0].Id, 1)
		t.Assert(*users[0].Passport, "user_1")
		t.Assert(*users[0].Password, "pass_1")
		t.Assert(users[0].Nickname, "name_1")
	})
}

func Test_Struct_Empty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		err := db.Model(table).Where("id=100").Scan(user)
		t.Assert(err, sql.ErrNoRows)
		t.AssertNE(user, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		one, err := db.Model(table).Where("id=100").One()
		t.AssertNil(err)
		var user *User
		t.Assert(one.Struct(&user), nil)
		t.Assert(user, nil)
	})

	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).Where("id=100").Scan(&user)
		t.AssertNil(err)
		t.Assert(user, nil)
	})
}

func Test_Structs_Empty(t *testing.T) {
	table := createTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
		Password string
		Nickname string
	}

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]User, 0)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]User, 10)
		t.Assert(all.Structs(&users), sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		var users []User
		t.Assert(all.Structs(&users), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]*User, 0)
		t.Assert(all.Structs(&users), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		users := make([]*User, 10)
		t.Assert(all.Structs(&users), sql.ErrNoRows)
	})
	gtest.C(t, func(t *gtest.T) {
		all, err := db.Model(table).Where("id>100").All()
		t.AssertNil(err)
		var users []*User
		t.Assert(all.Structs(&users), nil)
	})
}

type MyTime struct {
	gtime.Time
}

type MyTimeSt struct {
	CreateTime MyTime
}

func (st *MyTimeSt) UnmarshalValue(v interface{}) error {
	m := gconv.Map(v)
	t, err := gtime.StrToTime(gconv.String(m["create_time"]))
	if err != nil {
		return err
	}
	st.CreateTime = MyTime{*t}
	return nil
}

func Test_Model_Scan_CustomType_Time(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		st := new(MyTimeSt)
		err := db.Model(table).Fields("create_time").Scan(st)
		t.AssertNil(err)
		t.Assert(st.CreateTime.String(), "2018-10-24 10:00:00")
	})
	gtest.C(t, func(t *gtest.T) {
		var stSlice []*MyTimeSt
		err := db.Model(table).Fields("create_time").Scan(&stSlice)
		t.AssertNil(err)
		t.Assert(len(stSlice), TableSize)
		t.Assert(stSlice[0].CreateTime.String(), "2018-10-24 10:00:00")
		t.Assert(stSlice[9].CreateTime.String(), "2018-10-24 10:00:00")
	})
}

func Test_Model_Scan_CustomType_String(t *testing.T) {
	type MyString string

	type MyStringSt struct {
		Passport MyString
	}

	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		st := new(MyStringSt)
		err := db.Model(table).Fields("Passport").WherePri(1).Scan(st)
		t.AssertNil(err)
		t.Assert(st.Passport, "user_1")
	})
	gtest.C(t, func(t *gtest.T) {
		var sts []MyStringSt
		err := db.Model(table).Fields("Passport").Order("id asc").Scan(&sts)
		t.AssertNil(err)
		t.Assert(len(sts), TableSize)
		t.Assert(sts[0].Passport, "user_1")
	})
}

type User struct {
	Id         int
	Passport   string
	Password   string
	Nickname   string
	CreateTime *gtime.Time
}

func (user *User) UnmarshalValue(value interface{}) error {
	if record, ok := value.(gdb.Record); ok {
		*user = User{
			Id:         record["id"].Int(),
			Passport:   record["passport"].String(),
			Password:   "",
			Nickname:   record["nickname"].String(),
			CreateTime: record["create_time"].GTime(),
		}
		return nil
	}
	return gerror.NewCodef(gcode.CodeInvalidParameter, `unsupported value type for UnmarshalValue: %v`, reflect.TypeOf(value))
}

func Test_Model_Scan_UnmarshalValue(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)
	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[0].Passport, "user_1")
		t.Assert(users[0].Password, "")
		t.Assert(users[0].Nickname, "name_1")
		t.Assert(users[0].CreateTime.String(), CreateTime)

		t.Assert(users[9].Id, 10)
		t.Assert(users[9].Passport, "user_10")
		t.Assert(users[9].Password, "")
		t.Assert(users[9].Nickname, "name_10")
		t.Assert(users[9].CreateTime.String(), CreateTime)
	})
}

func Test_Model_Scan_Map(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	gtest.C(t, func(t *gtest.T) {
		var users []*User
		err := db.Model(table).Order("id asc").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
		t.Assert(users[0].Passport, "user_1")
		t.Assert(users[0].Password, "")
		t.Assert(users[0].Nickname, "name_1")
		t.Assert(users[0].CreateTime.String(), CreateTime)

		t.Assert(users[9].Id, 10)
		t.Assert(users[9].Passport, "user_10")
		t.Assert(users[9].Password, "")
		t.Assert(users[9].Nickname, "name_10")
		t.Assert(users[9].CreateTime.String(), CreateTime)
	})
}

func Test_Scan_AutoFilteringByStructAttributes(t *testing.T) {
	table := createInitTable()
	defer dropTable(table)

	type User struct {
		Id       int
		Passport string
	}
	// db.SetDebug(true)
	gtest.C(t, func(t *gtest.T) {
		var user *User
		err := db.Model(table).OrderAsc("id").Scan(&user)
		t.AssertNil(err)
		t.Assert(user.Id, 1)
	})
	gtest.C(t, func(t *gtest.T) {
		var users []User
		err := db.Model(table).OrderAsc("id").Scan(&users)
		t.AssertNil(err)
		t.Assert(len(users), TableSize)
		t.Assert(users[0].Id, 1)
	})
}
