package log

import (
	"testing"
)

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll_OneField/Phuslu_OneField/G1-16      	 8853261	       135.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G2-16      	 8776542	       135.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G3-16      	 8791803	       136.5 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Phuslu_OneField/G4-16      	 8733136	       136.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G1-16     	 7888216	       151.3 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G2-16     	 7875188	       152.4 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G3-16     	 7851661	       152.9 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Zerolog_OneField/G4-16     	 7870828	       152.1 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G1-16    	18455356	        64.79 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G2-16    	11042335	       108.0 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G3-16    	10664072	       110.3 ns/op	       6 B/op	       0 allocs/op
// BenchmarkAll_OneField/Arenalog_OneField/G4-16    	10783002	       109.1 ns/op	       5 B/op	       0 allocs/op

func BenchmarkAll_OneField(b *testing.B) {
	b.Run("Phuslu_OneField", BenchmarkPhuslu_OneField)
	b.Run("Zerolog_OneField", BenchmarkZerolog_OneField)
	b.Run("Arenalog_OneField", BenchmarkArenalog_OneField)
}

// cpu: AMD Ryzen 7 5800H with Radeon Graphics
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=1-16         	 4706324	       255.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=2-16         	 8641028	       138.6 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=3-16         	12812220	        94.44 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Phuslu_SeveralFields/gomaxprocs=4-16         	16818895	        72.68 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=1-16        	 4494318	       267.8 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=2-16        	 8168961	       145.2 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=3-16        	12094114	        99.91 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Zerolog_SeveralFields/gomaxprocs=4-16        	15629527	        77.51 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=1-16       	15442278	        76.86 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=2-16       	16374577	        71.50 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=3-16       	14127042	        88.08 ns/op	       0 B/op	       0 allocs/op
// BenchmarkAll_SeveralFields/Arenalog_SeveralFields/gomaxprocs=4-16       	13617920	        86.94 ns/op	       0 B/op	       0 allocs/op

func BenchmarkAll_SeveralFields(b *testing.B) {
	b.Run("Phuslu_SeveralFields", BenchmarkPhuslu_WithFields_Parallel)
	b.Run("Zerolog_SeveralFields", BenchmarkZerolog_WithFields_Parallel)
	b.Run("Arenalog_SeveralFields", BenchmarkArenalog_MultipleFields_Parallel)
}
