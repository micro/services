package handler_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/micro/services/users/handler"
	pb "github.com/micro/services/users/proto"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func testHandler(t *testing.T) *handler.Users {
	// connect to the database
	db, err := gorm.Open(postgres.Open("postgresql://postgres@localhost:5432/postgres?sslmode=disable"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Error connecting to database: %v", err)
	}

	// clean any data from a previous run
	if err := db.Exec("DROP TABLE IF EXISTS users, tokens CASCADE").Error; err != nil {
		t.Fatalf("Error cleaning database: %v", err)
	}

	// migrate the database
	if err := db.AutoMigrate(&handler.User{}, &handler.Token{}); err != nil {
		t.Fatalf("Error migrating database: %v", err)
	}

	return &handler.Users{DB: db, Time: time.Now}
}

func TestCreate(t *testing.T) {
	tt := []struct {
		Name      string
		FirstName string
		LastName  string
		Email     string
		Password  string
		Error     error
	}{
		{
			Name:     "MissingFirstName",
			LastName: "Doe",
			Email:    "john@doe.com",
			Password: "password",
			Error:    handler.ErrMissingFirstName,
		},
		{
			Name:      "MissingLastName",
			FirstName: "John",
			Email:     "john@doe.com",
			Password:  "password",
			Error:     handler.ErrMissingLastName,
		},
		{
			Name:      "MissingEmail",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
			Error:     handler.ErrMissingEmail,
		},
		{
			Name:      "InvalidEmail",
			FirstName: "John",
			LastName:  "Doe",
			Password:  "password",
			Email:     "foo.foo.foo",
			Error:     handler.ErrInvalidEmail,
		},
		{
			Name:      "InvalidPassword",
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "pwd",
			Error:     handler.ErrInvalidPassword,
		},
	}

	// test the validations
	h := testHandler(t)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			err := h.Create(context.TODO(), &pb.CreateRequest{
				FirstName: tc.FirstName,
				LastName:  tc.LastName,
				Email:     tc.Email,
				Password:  tc.Password,
			}, &pb.CreateResponse{})
			assert.Equal(t, tc.Error, err)
		})
	}

	t.Run("Valid", func(t *testing.T) {
		var rsp pb.CreateResponse
		req := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &req, &rsp)

		assert.NoError(t, err)
		u := rsp.User
		if u == nil {
			t.Fatalf("No user returned")
		}
		assert.NotEmpty(t, u.Id)
		assert.Equal(t, req.FirstName, u.FirstName)
		assert.Equal(t, req.LastName, u.LastName)
		assert.Equal(t, req.Email, u.Email)
		assert.NotEmpty(t, rsp.Token)
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		var rsp pb.CreateResponse
		req := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &req, &rsp)
		assert.Equal(t, handler.ErrDuplicateEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("DifferentEmail", func(t *testing.T) {
		var rsp pb.CreateResponse
		req := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "johndoe@gmail.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &req, &rsp)

		assert.NoError(t, err)
		u := rsp.User
		if u == nil {
			t.Fatalf("No user returned")
		}
		assert.NotEmpty(t, u.Id)
		assert.Equal(t, req.FirstName, u.FirstName)
		assert.Equal(t, req.LastName, u.LastName)
		assert.Equal(t, req.Email, u.Email)
	})
}

func TestRead(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingIDs", func(t *testing.T) {
		var rsp pb.ReadResponse
		err := h.Read(context.TODO(), &pb.ReadRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingIDs, err)
		assert.Nil(t, rsp.Users)
	})

	t.Run("NotFound", func(t *testing.T) {
		var rsp pb.ReadResponse
		err := h.Read(context.TODO(), &pb.ReadRequest{Ids: []string{"foo"}}, &rsp)
		assert.Nil(t, err)
		if rsp.Users == nil {
			t.Fatal("Expected the users object to not be nil")
		}
		assert.Nil(t, rsp.Users["foo"])
	})

	// create some mock data
	var rsp1 pb.CreateResponse
	req1 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(context.TODO(), &req1, &rsp1)
	assert.NoError(t, err)
	if rsp1.User == nil {
		t.Fatal("No user returned")
		return
	}

	var rsp2 pb.CreateResponse
	req2 := pb.CreateRequest{
		FirstName: "Apple",
		LastName:  "Tree",
		Email:     "apple@tree.com",
		Password:  "passwordabc",
	}
	err = h.Create(context.TODO(), &req2, &rsp2)
	assert.NoError(t, err)
	if rsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	// test the read
	var rsp pb.ReadResponse
	err = h.Read(context.TODO(), &pb.ReadRequest{
		Ids: []string{rsp1.User.Id, rsp2.User.Id},
	}, &rsp)
	assert.NoError(t, err)

	if rsp.Users == nil {
		t.Fatal("Users not returned")
		return
	}
	assert.NotNil(t, rsp.Users[rsp1.User.Id])
	assert.NotNil(t, rsp.Users[rsp2.User.Id])

	// check the users match
	if u := rsp.Users[rsp1.User.Id]; u != nil {
		assert.Equal(t, rsp1.User.Id, u.Id)
		assert.Equal(t, rsp1.User.FirstName, u.FirstName)
		assert.Equal(t, rsp1.User.LastName, u.LastName)
		assert.Equal(t, rsp1.User.Email, u.Email)
	}
	if u := rsp.Users[rsp2.User.Id]; u != nil {
		assert.Equal(t, rsp2.User.Id, u.Id)
		assert.Equal(t, rsp2.User.FirstName, u.FirstName)
		assert.Equal(t, rsp2.User.LastName, u.LastName)
		assert.Equal(t, rsp2.User.Email, u.Email)
	}
}

func TestUpdate(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{}, &rsp)
		assert.Equal(t, handler.ErrMissingID, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("NotFound", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{Id: "foo"}, &rsp)
		assert.Equal(t, handler.ErrNotFound, err)
		assert.Nil(t, rsp.User)
	})

	// create some mock data
	var cRsp1 pb.CreateResponse
	cReq1 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(context.TODO(), &cReq1, &cRsp1)
	assert.NoError(t, err)
	if cRsp1.User == nil {
		t.Fatal("No user returned")
		return
	}

	var cRsp2 pb.CreateResponse
	cReq2 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@gmail.com",
		Password:  "passwordabc",
	}
	err = h.Create(context.TODO(), &cReq2, &cRsp2)
	assert.NoError(t, err)
	if cRsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	t.Run("BlankFirstName", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, FirstName: &wrapperspb.StringValue{},
		}, &rsp)
		assert.Equal(t, handler.ErrMissingFirstName, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("BlankLastName", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, LastName: &wrapperspb.StringValue{},
		}, &rsp)
		assert.Equal(t, handler.ErrMissingLastName, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("BlankLastName", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, LastName: &wrapperspb.StringValue{},
		}, &rsp)
		assert.Equal(t, handler.ErrMissingLastName, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("BlankEmail", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, Email: &wrapperspb.StringValue{},
		}, &rsp)
		assert.Equal(t, handler.ErrMissingEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("InvalidEmail", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, Email: &wrapperspb.StringValue{Value: "foo.bar"},
		}, &rsp)
		assert.Equal(t, handler.ErrInvalidEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("EmailAlreadyExists", func(t *testing.T) {
		var rsp pb.UpdateResponse
		err := h.Update(context.TODO(), &pb.UpdateRequest{
			Id: cRsp1.User.Id, Email: &wrapperspb.StringValue{Value: cRsp2.User.Email},
		}, &rsp)
		assert.Equal(t, handler.ErrDuplicateEmail, err)
		assert.Nil(t, rsp.User)
	})

	t.Run("Valid", func(t *testing.T) {
		uReq := pb.UpdateRequest{
			Id:        cRsp1.User.Id,
			Email:     &wrapperspb.StringValue{Value: "foobar@gmail.com"},
			FirstName: &wrapperspb.StringValue{Value: "Foo"},
			LastName:  &wrapperspb.StringValue{Value: "Bar"},
		}
		var uRsp pb.UpdateResponse
		err := h.Update(context.TODO(), &uReq, &uRsp)
		assert.NoError(t, err)
		if uRsp.User == nil {
			t.Error("No user returned")
			return
		}
		assert.Equal(t, cRsp1.User.Id, uRsp.User.Id)
		assert.Equal(t, uReq.Email.Value, uRsp.User.Email)
		assert.Equal(t, uReq.FirstName.Value, uRsp.User.FirstName)
		assert.Equal(t, uReq.LastName.Value, uRsp.User.LastName)
	})
}

func TestDelete(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingID", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{}, &pb.DeleteResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	// create some mock data
	var cRsp pb.CreateResponse
	cReq := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(context.TODO(), &cReq, &cRsp)
	assert.NoError(t, err)
	if cRsp.User == nil {
		t.Fatal("No user returned")
		return
	}

	t.Run("Valid", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{
			Id: cRsp.User.Id,
		}, &pb.DeleteResponse{})
		assert.NoError(t, err)

		// check it was actually deleted
		var rsp pb.ReadResponse
		err = h.Read(context.TODO(), &pb.ReadRequest{
			Ids: []string{cRsp.User.Id},
		}, &rsp)
		assert.NoError(t, err)
		assert.Nil(t, rsp.Users[cRsp.User.Id])
	})

	t.Run("Retry", func(t *testing.T) {
		err := h.Delete(context.TODO(), &pb.DeleteRequest{
			Id: cRsp.User.Id,
		}, &pb.DeleteResponse{})
		assert.NoError(t, err)
	})
}

func TestList(t *testing.T) {
	h := testHandler(t)

	// create some mock data
	var cRsp1 pb.CreateResponse
	cReq1 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(context.TODO(), &cReq1, &cRsp1)
	assert.NoError(t, err)
	if cRsp1.User == nil {
		t.Fatal("No user returned")
		return
	}

	var cRsp2 pb.CreateResponse
	cReq2 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "johndoe@gmail.com",
		Password:  "passwordabc",
	}
	err = h.Create(context.TODO(), &cReq2, &cRsp2)
	assert.NoError(t, err)
	if cRsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	var rsp pb.ListResponse
	err = h.List(context.TODO(), &pb.ListRequest{}, &rsp)
	assert.NoError(t, err)
	if rsp.Users == nil {
		t.Error("No users returned")
		return
	}

	var u1Found, u2Found bool
	for _, u := range rsp.Users {
		switch u.Id {
		case cRsp1.User.Id:
			assertUsersMatch(t, cRsp1.User, u)
			u1Found = true
		case cRsp2.User.Id:
			assertUsersMatch(t, cRsp2.User, u)
			u2Found = true
		default:
			t.Fatal("Unexpected user returned")
			return
		}
	}
	assert.True(t, u1Found)
	assert.True(t, u2Found)
}

func TestLogin(t *testing.T) {
	h := testHandler(t)

	// create some mock data
	var cRsp pb.CreateResponse
	cReq := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(context.TODO(), &cReq, &cRsp)
	assert.NoError(t, err)
	if cRsp.User == nil {
		t.Fatal("No user returned")
		return
	}

	tt := []struct {
		Name     string
		Email    string
		Password string
		Error    error
		User     *pb.User
	}{
		{
			Name:     "MissingEmail",
			Password: "passwordabc",
			Error:    handler.ErrMissingEmail,
		},
		{
			Name:  "MissingPassword",
			Email: "john@doe.com",
			Error: handler.ErrInvalidPassword,
		},
		{
			Name:     "UserNotFound",
			Email:    "foo@bar.com",
			Password: "passwordabc",
			Error:    handler.ErrNotFound,
		},
		{
			Name:     "IncorrectPassword",
			Email:    "john@doe.com",
			Password: "passwordabcdef",
			Error:    handler.ErrIncorrectPassword,
		},
		{
			Name:     "Valid",
			Email:    "john@doe.com",
			Password: "passwordabc",
			User:     cRsp.User,
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			var rsp pb.LoginResponse
			err := h.Login(context.TODO(), &pb.LoginRequest{
				Email: tc.Email, Password: tc.Password,
			}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.User != nil {
				assertUsersMatch(t, tc.User, rsp.User)
				assert.NotEmpty(t, rsp.Token)
			} else {
				assert.Nil(t, tc.User)
			}
		})
	}
}

func TestLogout(t *testing.T) {
	h := testHandler(t)

	t.Run("MissingUserID", func(t *testing.T) {
		err := h.Logout(context.TODO(), &pb.LogoutRequest{}, &pb.LogoutResponse{})
		assert.Equal(t, handler.ErrMissingID, err)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		err := h.Logout(context.TODO(), &pb.LogoutRequest{Id: uuid.New().String()}, &pb.LogoutResponse{})
		assert.Equal(t, handler.ErrNotFound, err)
	})

	t.Run("Valid", func(t *testing.T) {
		// create some mock data
		var cRsp pb.CreateResponse
		cReq := pb.CreateRequest{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john@doe.com",
			Password:  "passwordabc",
		}
		err := h.Create(context.TODO(), &cReq, &cRsp)
		assert.NoError(t, err)
		if cRsp.User == nil {
			t.Fatal("No user returned")
			return
		}

		err = h.Logout(context.TODO(), &pb.LogoutRequest{Id: cRsp.User.Id}, &pb.LogoutResponse{})
		assert.NoError(t, err)

		err = h.Validate(context.TODO(), &pb.ValidateRequest{Token: cRsp.Token}, &pb.ValidateResponse{})
		assert.Error(t, err)
	})
}

func TestValidate(t *testing.T) {
	h := testHandler(t)

	// create some mock data
	var cRsp1 pb.CreateResponse
	cReq1 := pb.CreateRequest{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john@doe.com",
		Password:  "passwordabc",
	}
	err := h.Create(context.TODO(), &cReq1, &cRsp1)
	assert.NoError(t, err)
	if cRsp1.User == nil {
		t.Fatal("No user returned")
		return
	}

	var cRsp2 pb.CreateResponse
	cReq2 := pb.CreateRequest{
		FirstName: "Barry",
		LastName:  "Doe",
		Email:     "barry@doe.com",
		Password:  "passwordabc",
	}
	err = h.Create(context.TODO(), &cReq2, &cRsp2)
	assert.NoError(t, err)
	if cRsp2.User == nil {
		t.Fatal("No user returned")
		return
	}

	tt := []struct {
		Name  string
		Token string
		Time  func() time.Time
		Error error
		User  *pb.User
	}{
		{
			Name:  "MissingToken",
			Error: handler.ErrMissingToken,
		},
		{
			Name:  "InvalidToken",
			Error: handler.ErrInvalidToken,
			Token: uuid.New().String(),
		},
		{
			Name:  "ExpiredToken",
			Error: handler.ErrTokenExpired,
			Token: cRsp1.Token,
			Time:  func() time.Time { return time.Now().Add(time.Hour * 24 * 8) },
		},
		{
			Name:  "ValidToken",
			User:  cRsp2.User,
			Token: cRsp2.Token,
			Time:  func() time.Time { return time.Now().Add(time.Hour * 24 * 3) },
		},
		{
			Name:  "RefreshedToken",
			User:  cRsp2.User,
			Token: cRsp2.Token,
			Time:  func() time.Time { return time.Now().Add(time.Hour * 24 * 8) },
		},
	}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Time == nil {
				h.Time = time.Now
			} else {
				h.Time = tc.Time
			}

			var rsp pb.ValidateResponse
			err := h.Validate(context.TODO(), &pb.ValidateRequest{Token: tc.Token}, &rsp)
			assert.Equal(t, tc.Error, err)

			if tc.User != nil {
				assertUsersMatch(t, tc.User, rsp.User)
			} else {
				assert.Nil(t, tc.User)
			}
		})
	}
}

func assertUsersMatch(t *testing.T, exp, act *pb.User) {
	if act == nil {
		t.Error("No user returned")
		return
	}
	assert.Equal(t, exp.Id, act.Id)
	assert.Equal(t, exp.FirstName, act.FirstName)
	assert.Equal(t, exp.LastName, act.LastName)
	assert.Equal(t, exp.Email, act.Email)
}
