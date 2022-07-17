package tool

import "testing"

func TestF2b(t *testing.T) {
	tests := []struct {
		name string
	}{
		{"successful case"},
	}
	for range tests {
		F2b()
	}
}

func Test_fileToByte(t *testing.T) {
	type args struct {
		inPath  string
		outPath string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"successful", args{"D:\\Software\\GoProjects\\src\\chess\\client\\img", "D:\\Software\\GoProjects\\src\\chess\\client\\chess"}, false},
		{"failed", args{"", ""}, true},
	}
	for _, tt := range tests {
		if err := fileToByte(tt.args.inPath, tt.args.outPath); (err != nil) != tt.wantErr {
			t.Errorf("%q. fileToByte() error = %v, wantErr %v", tt.name, err, tt.wantErr)
		}
	}
}
