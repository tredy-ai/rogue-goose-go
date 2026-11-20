package orbitalwiki

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseOrbitClass(t *testing.T) {
	tests := []struct {
		name, in string
		want     OrbitClass
		wantErr  bool
	}{
		{name: "canonical", in: "LEO", want: OrbitClassLEO},
		{name: "lowercase", in: "vleo", want: OrbitClassVLEO},
		{name: "surrounding whitespace", in: "  GEO\t", want: OrbitClassGEO},
		{name: "empty is unset", in: "", want: OrbitClassUnset},
		{name: "unknown value", in: "LUNAR", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseOrbitClass(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name, in string
		want     Status
		wantErr  bool
	}{
		{name: "canonical", in: "OPERATIONAL", want: StatusOperational},
		{name: "lowercase", in: "extended_mission", want: StatusExtendedMission},
		{name: "surrounding whitespace", in: " DECAYED\n", want: StatusDecayed},
		{name: "empty is unset", in: "", want: StatusUnset},
		{name: "unknown value", in: "SIDEWAYS", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseStatus(tt.in)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestOrbitClassStringRoundTrip(t *testing.T) {
	for _, v := range OrbitClassValues() {
		got, err := ParseOrbitClass(v)
		require.NoError(t, err)
		assert.Equal(t, v, got.String())
	}
}

func TestStatusStringRoundTrip(t *testing.T) {
	for _, v := range StatusValues() {
		got, err := ParseStatus(v)
		require.NoError(t, err)
		assert.Equal(t, v, got.String())
	}
}
