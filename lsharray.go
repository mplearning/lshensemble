package lshensemble

import (
	"math"
	"sync"
)

// LshForestArray represents a MinHash LSH implemented using an array of LshForest.
// It allows a wider range for the K and L parameters.
type LshForestArray struct {
	maxK    int
	numHash int
	array   []*LshForest
}

// Initialize with parameters:
// maxK is the maximum value for the MinHash parameter K - the number of hash functions per "band".
// numHash is the number of hash functions in MinHash.
func NewLshForestArray(maxK, numHash int) *LshForestArray {
	array := make([]*LshForest, maxK)
	for k := 1; k <= maxK; k++ {
		array[k-1] = NewLshForest(k, numHash/k)
	}
	return &LshForestArray{
		maxK:    maxK,
		numHash: numHash,
		array:   array,
	}
}

// Add a key with MinHash signature into the index.
// The key won't be searchable until Index() is called.
func (a *LshForestArray) Add(key string, sig Signature) {
	var wg sync.WaitGroup
	wg.Add(len(a.array))
	for i := range a.array {
		go func(lsh *LshForest) {
			lsh.Add(key, sig)
			wg.Done()
		}(a.array[i])
	}
	wg.Wait()
}

// Makes all the keys added searchable.
func (a *LshForestArray) Index() {
	var wg sync.WaitGroup
	wg.Add(len(a.array))
	for i := range a.array {
		go func(lsh *LshForest) {
			lsh.Index()
			wg.Done()
		}(a.array[i])
	}
	wg.Wait()
}

// Return candidate keys given the query signature and parameters.
func (a *LshForestArray) Query(sig Signature, K, L int, out chan<- string, done <-chan struct{}) {
	a.array[K-1].Query(sig, -1, L, out, done)
}

// OptimalKL returns the optimal K and L for containment search,
// and the false positive and negative probabilities.
// where x is the indexed domain size, q is the query domain size,
// and t is the containment threshold.
func (a *LshForestArray) OptimalKL(x, q int, t float64) (optK, optL int, fp, fn float64) {
	minError := math.MaxFloat64
	for l := 1; l <= a.numHash; l++ {
		for k := 1; k <= a.maxK; k++ {
			if k*l > a.numHash {
				continue
			}
			currFp := probFalsePositive(x, q, l, k, t, integrationPrecision)
			currFn := probFalseNegative(x, q, l, k, t, integrationPrecision)
			currErr := currFn + currFp
			if minError > currErr {
				minError = currErr
				optK = k
				optL = l
				fp = currFp
				fn = currFn
			}
		}
	}
	return
}
