package bili

import "testing"

func TestGetBVID(t *testing.T) {
	type args struct {
		rawURL string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "1. ",
			args: args{
				rawURL: "https://www.bilibili.com/video/BV1NP411o7MB/?spm_id_from=333.1007.top_right_bar_window_custom_collection.content.click",
			},
			want:    "BV1NP411o7MB",
			wantErr: false,
		},

		{
			name: "2. ",
			args: args{
				rawURL: "https://www.bilibili.com/video/BV1NP411o7MB/?spm_id_from=333.1007.top_right_bar_window_custom_collection.content.click",
			},
			want:    "BV1NP411o7MB",
			wantErr: false,
		},
		{
			name: "3. ",
			args: args{
				rawURL: "BV1NP411o7MB",
			},
			want:    "BV1NP411o7MB",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBVIDFromURL(tt.args.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBVID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBVID() = %v, want %v", got, tt.want)
			}
		})
	}
}
