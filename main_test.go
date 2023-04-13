package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func Test_isPrime(t *testing.T) {
	primeTests := []struct {
		name     string
		testNum  int
		expected bool
		msg      string
	}{
		{"prime", 7, true, "7 is a prime number!"},
		{"not prime", 8, false, "8 is not a prime number because it is divisible by 2!"},
		{"zero", 0, false, "0 is not prime, by definition!"},
		{"one", 1, false, "1 is not prime, by definition!"},
		{"negative number", -11, false, "Negative numbers are not prime, by definition!"},
	}

	for _, e := range primeTests {
		result, msg := isPrime(e.testNum)
		if e.expected && !result {
			t.Errorf("%s: expected true but got false", e.name)
		}

		if !e.expected && result {
			t.Errorf("%s: expected false but got true", e.name)
		}

		if e.msg != msg {
			t.Errorf("%s: expected %s but got %s", e.name, e.msg, msg)
		}
	}
}

func TestCheckNumbers(t *testing.T) {
	tests := []struct {
		input string
		want  string
		done  bool
	}{
		{"q\n", "", true},
		{"hello\n", "Please enter a whole number!", false},
		{"0\n", "0 is not prime, by definition!", false},
		{"1\n", "1 is not prime, by definition!", false},
		{"-42\n", "Negative numbers are not prime, by definition!", false},
		{"2\n", "2 is a prime number!", false},
		{"3\n", "3 is a prime number!", false},
		{"4\n", "4 is not a prime number because it is divisible by 2!", false},
	}

	for _, tt := range tests {
		scanner := bufio.NewScanner(strings.NewReader(tt.input))
		got, done := checkNumbers(scanner)
		if got != tt.want {
			t.Errorf("checkNumbers(%q) = %q, want %q", tt.input, got, tt.want)
		}
		if done != tt.done {
			t.Errorf("checkNumbers(%q) returned done = %v, want %v", tt.input, done, tt.done)
		}
	}
}

func TestPrompt(t *testing.T) {
	// create a temporary file to capture the output
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	// replace stdout with the temporary file
	old := os.Stdout
	os.Stdout = tmpfile
	defer func() {
		os.Stdout = old
	}()

	// call the function being tested
	prompt()

	// read the captured output from the temporary file
	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}
	output, err := ioutil.ReadAll(tmpfile)
	if err != nil {
		t.Fatal(err)
	}

	// assert that the output matches the expected prompt
	expected := "-> "
	if string(output) != expected {
		t.Errorf("expected prompt %q, but got %q", expected, string(output))
	}
}

func TestReadUserInput(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedRes  string
		expectedDone bool
	}{
		{
			name:         "valid input",
			input:        "7\n",
			expectedRes:  "7 is a prime number!",
			expectedDone: false,
		},
		{
			name:         "invalid input",
			input:        "hello\n",
			expectedRes:  "Please enter a whole number!",
			expectedDone: false,
		},
		{
			name:         "quit command",
			input:        "q\n",
			expectedRes:  "",
			expectedDone: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, w, _ := os.Pipe()
			oldStdin := os.Stdin
			defer func() {
				os.Stdin = oldStdin
			}()
			os.Stdin = r

			doneChan := make(chan bool)
			go func() {
				res, done := checkNumbers(bufio.NewScanner(w))
				if done != tt.expectedDone {
					t.Errorf("expected done chan to be %v, but got %v", tt.expectedDone, done)
				}
				if res != tt.expectedRes {
					t.Errorf("expected result to be %q, but got %q", tt.expectedRes, res)
				}
			}()

			if _, err := fmt.Fprint(w, tt.input); err != nil {
				t.Fatal(err)
			}

			// close the writer so the scanner gets an EOF
			if err := w.Close(); err != nil {
				t.Fatal(err)
			}

			// wait for the goroutine to finish
			<-doneChan
			close(doneChan)
		})
	}
}
