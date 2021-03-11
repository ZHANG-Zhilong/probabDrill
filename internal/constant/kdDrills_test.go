package constant

import (
	"probabDrill/internal/entity"
	"reflect"
	"testing"
)

func TestNearKDrills(t *testing.T) {
	drills := GetDrillSet()
	nears1 := drills[0].NearDrills(drills, 1)
	nears2 := drills[1].NearDrills(drills, 1)
	nears3 := drills[2].NearDrills(drills, 1)
	nears4 := drills[3].NearDrills(drills, 1)

	type args struct {
		vdrill entity.Drill
		k      int
	}
	tests := []struct {
		name       string
		args       args
		wantDrills []entity.Drill
	}{
		// TODO: Add test cases.
		{"fuck", args{drills[0], 1}, nears1},
		{"fuck", args{drills[1], 1}, nears2},
		{"fuck", args{drills[2], 1}, nears3},
		{"fuck", args{drills[3], 1}, nears4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDrills := NearKDrills(tt.args.vdrill, tt.args.k); !reflect.DeepEqual(gotDrills, tt.wantDrills) {
				t.Errorf("%v\n", tt.args)
				t.Errorf("\nNearHelpDrills() =\t%v,\n want \t\t\t\t%v", gotDrills, tt.wantDrills)
			}
		})
	}
}

func TestNearHelpDrills(t *testing.T) {
	drills := GetHelpDrillSet()
	nears1 := drills[0].NearDrills(drills, 1)
	nears2 := drills[999].NearDrills(drills, 1)
	nears3 := drills[299].NearDrills(drills, 1)
	nears4 := drills[399].NearDrills(drills, 1)

	type args struct {
		vdrill entity.Drill
		k      int
	}
	tests := []struct {
		name       string
		args       args
		wantDrills []entity.Drill
	}{
		// TODO: Add test cases.
		{"fuck", args{drills[0], 1}, nears1},
		{"fuck", args{drills[999], 1}, nears2},
		{"fuck", args{drills[299], 1}, nears3},
		{"fuck", args{drills[399], 1}, nears4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotDrills := NearHelpDrills(tt.args.vdrill, tt.args.k); !reflect.DeepEqual(gotDrills, tt.wantDrills) {
				t.Errorf("%v\n", tt.args)
				t.Errorf("\nNearHelpDrills() =\t%v,\n want \t\t\t\t%v", gotDrills, tt.wantDrills)
			}
		})
	}
}

func TestDemo(t *testing.T) {
	Demo()
}
