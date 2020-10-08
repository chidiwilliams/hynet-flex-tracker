package main

import "testing"

func Test_usageFromResponse(t *testing.T) {
	tests := []struct {
		name     string
		response map[string]string
		want     int
		wantErr  bool
	}{
		{
			response: map[string]string{
				"ussd_action": "0",
				"ussd_dcs":    "72",
				"ussd_data":   "004400650061007200200043007500730074006F006D00650072002C00200079006F00750020006800610076006500200063006F006E00730075006D006500640020003100360032003700390038004D0042003B002000620075006E0064006C006500200065007800700069007200650073002000310037002F00310030002F0032003000320030002E00200045006E006A006F007900200079006F0075007200200075006E006C0069006D0069007400650064002000620072006F007700730069006E006700200065007800700065007200690065006E00630065002E0020005400680061006E006B0079006F0075002E000A",
			},
			want: 162798,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := usageFromResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("usageFromResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("usageFromResponse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
