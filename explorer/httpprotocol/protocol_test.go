package httpprotocol

import (
	"errors"
	"net/http"
	"testing"
)

func TestFailed(t *testing.T) {
	result := Failed(errors.New("test"))
	if result.Code != http.StatusInternalServerError {
		t.Errorf("test case TestFailed failed")
	}

	if "test" != result.Result.(error).Error() {
		t.Errorf("test case TestSucceed failed")
	}
}

func TestSucceed(t *testing.T) {
	result := Succeed()
	if result.Code != http.StatusOK {
		t.Errorf("test case TestSucceed failed")
	}

	if nil != result.Result {
		t.Errorf("test case TestSucceed failed")
	}
}

func TestSucceedWithResultNil(t *testing.T) {
	result := SucceedWithResult(nil)
	if result.Code != http.StatusOK {
		t.Errorf("test case TestSucceed failed")
	}

	if nil != result.Result {
		t.Errorf("test case TestSucceed failed")
	}
}

func TestSucceedWithResult(t *testing.T) {
	result := SucceedWithResult(string("test"))
	if result.Code != http.StatusOK {
		t.Errorf("test case TestSucceed failed")
	}

	if "test" != result.Result.(string) {
		t.Errorf("test case TestSucceed failed")
	}
}
