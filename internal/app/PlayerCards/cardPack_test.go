package playercards

import (
	"fmt"
	"testing"

	"github.com/nazandr/fantasy_api/internal/app/store"
	"github.com/stretchr/testify/assert"
)

func Test_OpenCommonPack(t *testing.T) {
	s, _ := store.TestStore(t, store.NewConfig().Database_url)

	pack := NewPack()
	freqs := make(map[int]int)
	pack.OpenCommonPack(s)

	for i := 0; i < 200; i++ {
		err := pack.OpenCommonPack(s)
		assert.NoError(t, err)

		for _, f := range pack.Cards {
			freqs[f.Rarity]++
		}
	}

	fmt.Printf("\nCommon: %d\tRare: %d\tEpic: %d\tLegend: %d\n",
		freqs[0], freqs[1], freqs[2], freqs[3])

	assert.NotEmpty(t, pack)
}

func Test_randomWithProbabilitis(t *testing.T) {
	it := []int{0, 1, 2, 3}
	w := []float32{0.7, 0.23, 0.05, 0.02}

	freqs := make(map[int]int)

	for i := 0; i < 1000; i++ {
		r := randomWithProbabilitis(it, w)
		freqs[r]++
	}

	fmt.Printf("\nCommon: %d\tRare: %d\tEpic: %d\tLegend: %d\n",
		freqs[0], freqs[1], freqs[2], freqs[3])
}
