package appui

import (
	"context"
	"io/ioutil"
	"path/filepath"
	"strings"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/moncho/dry/ui"
	"github.com/moncho/dry/ui/termui"
)

type testVolumesService struct {
	volumes []*types.Volume
}

func (v testVolumesService) VolumeList(context.Context) ([]*types.Volume, error) {
	return v.volumes, nil
}

func TestVolumesWidget(t *testing.T) {
	type args struct {
		service volumesService
		s       Screen
	}
	tests := []struct {
		name       string
		args       args
		operations func(*VolumesWidget)
	}{
		{
			"TestVolumesWidget unmounted widget",
			args{
				testVolumesService{},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     5,
				},
			},
			func(*VolumesWidget) {},
		},
		{
			"TestVolumesWidget mounted widget no volumes",
			args{
				testVolumesService{},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     10},
			},
			func(w *VolumesWidget) {
				w.Mount()
			},
		},
		{
			"TestVolumesWidget mounted widget two volumes",
			args{
				testVolumesService{
					volumes: []*types.Volume{
						&types.Volume{
							Driver: "local",
							Name:   "volume1",
						},
						&types.Volume{
							Driver: "local",
							Name:   "volume2",
						},
					},
				},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     10},
			},
			func(w *VolumesWidget) {
				w.Mount()
			},
		},
		{
			"TestVolumesWidget show first 4 volumes",
			args{
				testVolumesService{
					volumes: []*types.Volume{
						&types.Volume{
							Driver: "local",
							Name:   "volume1",
						},
						&types.Volume{
							Driver: "local",
							Name:   "volume2",
						},
						&types.Volume{
							Driver: "local",
							Name:   "volume3",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume4",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume5",
						},
					},
				},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     8},
			},
			func(w *VolumesWidget) {
				w.Mount()
			},
		},
		{
			"TestVolumesWidget show last 4 volumes",
			args{
				testVolumesService{
					volumes: []*types.Volume{
						&types.Volume{
							Driver: "local",
							Name:   "volume1",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume2",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume3",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume4",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume5",
						},
					},
				},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     8},
			},
			func(w *VolumesWidget) {
				w.screen.Cursor().Max(5)
				w.screen.Cursor().Bottom()
				w.Mount()
			},
		},
		{
			"TestVolumesWidget sort volumes",
			args{
				testVolumesService{
					volumes: []*types.Volume{
						&types.Volume{
							Driver: "local",
							Name:   "volume5",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume4",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume3",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume2",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume1",
						},
					},
				},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     8},
			},
			func(w *VolumesWidget) {
				w.Sort()
				w.Mount()
			},
		},
		{
			"TestVolumesWidget double sort volumes",
			args{
				testVolumesService{
					volumes: []*types.Volume{
						&types.Volume{
							Driver: "local1",
							Name:   "volume5",
						}, &types.Volume{
							Driver: "local1",
							Name:   "volume4",
						}, &types.Volume{
							Driver: "local2",
							Name:   "volume3",
						}, &types.Volume{
							Driver: "local2",
							Name:   "volume2",
						}, &types.Volume{
							Driver: "local2",
							Name:   "volume1",
						},
					},
				},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     8},
			},
			func(w *VolumesWidget) {
				w.Sort()
				w.Sort()
				w.Mount()
			},
		},
		{
			"TestVolumesWidget filter volumes",
			args{
				testVolumesService{
					volumes: []*types.Volume{
						&types.Volume{
							Driver: "local",
							Name:   "volume5",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume4",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume3",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume2",
						}, &types.Volume{
							Driver: "local",
							Name:   "volume1",
						},
					},
				},
				&testScreen{
					cursor: &ui.Cursor{},
					x1:     30,
					y1:     8},
			},
			func(w *VolumesWidget) {
				w.Filter("volume3")
				w.Mount()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			widget := NewVolumesWidget(tt.args.service, tt.args.s)
			tt.operations(widget)
			got, err := termui.String(widget)
			if err != nil {
				t.Fatalf("could not stringify buffer content: %s", err.Error())
			}

			golden := filepath.Join("testdata", strings.ReplaceAll(tt.name, " ", "_")+".golden")
			if *update {
				ioutil.WriteFile(golden, []byte(got), 0644)
			}
			want, err := ioutil.ReadFile(golden)

			if got != string(want) {
				t.Errorf("NewVolumesWidget() = %v, want %v", got, string(want))
			}
		})
	}
}
