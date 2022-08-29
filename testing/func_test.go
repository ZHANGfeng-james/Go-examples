package testing_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func TestMain(main *testing.M) {
	log.Println("setup")

	main.Run()

	log.Println("teardown")
}

func TestSomeOutput(t *testing.T) {
	type Input struct {
		param string
		want  string
	}

	inputs := []Input{
		{"i", "I"},
	}

	for _, input := range inputs {
		if got := someOutput(input.param); got != input.want {
			t.Fail()
		}
	}
}

// 成功的单测
func TestShouldUpdateStats(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	// 执行update ...的时候返回值
	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	// 执行insert ...的返回值
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	// now we execute our method
	if err = recordStats(db, 2, 3); err != nil {
		t.Errorf("error was not expected while updating stats: %s", err)
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

// 失败的单测
func TestShouldRollbackStatUpdatesOnFailure(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	mock.ExpectExec("UPDATE products").WillReturnResult(sqlmock.NewResult(1, 1))
	// 返回错误
	mock.ExpectExec("INSERT INTO product_viewers").WithArgs(2, 3).WillReturnError(fmt.Errorf("some error"))

	mock.ExpectRollback()

	// now we execute our method
	if err = recordStats(db, 2, 3); err == nil {
		t.Errorf("was expecting an error, but there was none")
	}
	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>")
	}
	req := httptest.NewRequest("GET", "http://example.com/foo", nil)
	w := httptest.NewRecorder()

	handler(w, req) // handler is func(w http.ResponseWriter, r *http.Request)

	resp := w.Result()

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(resp.StatusCode)
	fmt.Println(resp.Header.Get("Content-Type"))
	fmt.Println(string(body))
}

func TestHTTPServer(t *testing.T) {
	ts := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s\n", r.Proto)
	}))
	ts.EnableHTTP2 = true
	ts.StartTLS()

	defer ts.Close()

	fmt.Println(ts.URL)
	res, err := ts.Client().Get(ts.URL) // https://127.0.0.1:57408
	if err != nil {
		log.Fatal(err)
	}
	greeting, err := io.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", greeting)
}
