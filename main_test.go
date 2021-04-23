package main_test

import "testing"

func TestMain(t *testing.T){

}

func TestSum(t *testing.T){
  total := Sum(5, 5)
    if total != 10 {
       t.Errorf("Sum was incorrect, got: %d, want: %d.", total, 10)
    }
}
