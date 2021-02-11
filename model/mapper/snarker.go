package mapper

import (
	"github.com/figment-networks/mina-indexer/client/graph"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/util"
)

// Snarker returns a snarker model constructed from the graph input
func Snarker(block *graph.Block, job *graph.CompletedWork) (*model.Snarker, error) {
	height := BlockHeight(block)
	time := BlockTime(block)

	snarker := &model.Snarker{
		PublicKey:   job.Prover,
		Fee:         util.MustUInt64(job.Fee),
		JobsCount:   1,
		WorksCount:  len(job.WorkIds),
		StartHeight: height,
		StartTime:   time,
		LastHeight:  height,
		LastTime:    time,
	}

	return snarker, snarker.Validate()
}

// Snarkers returns a collection of snarker models constructed from the graph input
func Snarkers(block *graph.Block) ([]model.Snarker, error) {
	if block == nil {
		return nil, nil
	}

	snarkers := map[string]*model.Snarker{}

	for _, job := range block.SnarkJobs {
		if snarkers[job.Prover] == nil {
			s, err := Snarker(block, job)
			if err != nil {
				return nil, err
			}
			snarkers[job.Prover] = s
		} else {
			snarkers[job.Prover].JobsCount++
			snarkers[job.Prover].WorksCount += len(job.WorkIds)
		}
	}

	result := []model.Snarker{}
	for _, s := range snarkers {
		result = append(result, *s)
	}

	return result, nil
}
