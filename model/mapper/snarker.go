package mapper

import (
	"github.com/figment-networks/mina-indexer/coda"
	"github.com/figment-networks/mina-indexer/model"
	"github.com/figment-networks/mina-indexer/model/util"
)

func Snarker(block *coda.Block, job *coda.CompletedWork) (*model.Snarker, error) {
	height := BlockHeight(block)
	time := BlockTime(block)

	snarker := &model.Snarker{
		Account:     job.Prover,
		Fee:         util.MustUInt64(job.Fee),
		JobsCount:   1,
		WorksCount:  len(job.WorkIds),
		StartHeight: height,
		StartTime:   time,
		LastHeight:  height,
		LastTime:    time,
	}

	if err := snarker.Validate(); err != nil {
		return nil, err
	}
	return snarker, nil
}

func Snarkers(block *coda.Block) ([]model.Snarker, error) {
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
