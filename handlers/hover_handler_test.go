package handlers

// func TestFindWord(t *testing.T) {
// 	t.Run("returns error if file cannot be found", func(t *testing.T) {
// 		_, err := WordFinder(0, 0, "")
// 		if err == nil {
// 			t.Fatalf("expected: error, got: nil")
// 		}
//
// 		want := "open : no such file or directory"
// 		if err.Error() != want {
// 			t.Errorf("expected: %s, got: %s", want, err.Error())
// 		}
// 	})
// 	t.Run("returns the work at the character and line number specified", func(t *testing.T) {
// 		tempF, err := os.CreateTemp("", "hover_handler_test")
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		defer os.Remove(tempF.Name())
//
// 		_, err = tempF.Write([]byte("some stuff"))
// 		if err != nil {
// 			t.Fatal(err)
// 		}
//
// 		// word, err := WordFinder(0, 0, "file://"+tempF.Name())
// 		// word, err := WordFinder(0, 0, tempF.Name())
// 		if err != nil {
// 			t.Fatalf("error: %s", err)
// 		}
//
// 		want := "I"
// 		if word != want {
// 			t.Errorf("expected: %s, got: '%s'", want, word)
// 		}
// 	})
// }
