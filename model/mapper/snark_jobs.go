package mapper

import (
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/types"
)

// SnarkJob returns a job model constructed from the graph input
func SnarkJob(block *graph.Block, w *graph.CompletedWork) (*model.SnarkJob, error) {
	j := &model.SnarkJob{
		Height:     BlockHeight(block),
		Time:       BlockTime(block),
		Prover:     w.Prover,
		Fee:        types.NewAmount(w.Fee),
		WorksCount: len(w.WorkIds),
	}
	return j, j.Validate()
}

// SnarkJobs returns list of jobs constructed from the graph input
func SnarkJobs(block *graph.Block) ([]model.SnarkJob, error) {
	if block == nil {
		return nil, nil
	}

	result := []model.SnarkJob{}

	for _, w := range block.SnarkJobs {
		j, err := SnarkJob(block, w)
		if err != nil {
			return nil, err
		}
		result = append(result, *j)
	}

	return result, nil
}
