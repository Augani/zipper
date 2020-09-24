package zipper

import "testing"

func TestZip(t *testing.T) {
	var destination, err = ZipIt("C:\\Users\\augustus\\Desktop\\projects\\goWorkSpace\\src\\github.com\\augani\\musapp", "", "")
	if err != nil {
		t.Errorf("There was an error")
	}
	if destination != "C:\\Users\\augustus\\Desktop\\projects\\goWorkSpace\\src\\github.com\\augani\\musapp.zip" {
		t.Errorf("Destination was wrong, got: %s, want: %s.", destination, "C:\\Users\\augan\\Desktop\\newWestern.zip")
	}
}

func TestUnZip(t *testing.T) {
	var done, err = UnZipIt("C:\\Users\\augustus\\Desktop\\projects\\goWorkSpace\\src\\github.com\\augani\\musapp.zip", "C:\\Users\\augustus\\Music")
	if err != nil {
		t.Errorf("There was an error")
	}
	if !done {
		t.Errorf("Destination was wrong, got: %t, want: %t.", done, true)
	}
}
