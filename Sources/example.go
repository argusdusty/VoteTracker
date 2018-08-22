package Sources

import (
	. "VoteTracker/Utils"
)

type ExampleSource struct {
	Param string
}

func (U *ExampleSource) Load(Data Summary) (result Summary, err error) {
	// Clear the existing data so you don't duplicate candidates
	result = Data
	result.Candidates = nil
	result.Regions = make(map[string]SummaryRegion)

	// Here you would load from an API or other online website to get the current results using `LoadURL` (which caches for a few seconds and sends valid If-Modified-Since headers)

	// Do results
	result.Candidates = append(result.Candidates, SummaryCandidate{Candidate: "argusdusty", Votes: 100, Winner: true, PartyLetter: "A"})
	result.Candidates = append(result.Candidates, SummaryCandidate{Candidate: "badguy", Votes: 1, Winner: false, PartyLetter: "B"})
	for _, c := range result.Candidates {
		result.TotalVotes += c.Votes
	}
	result.PortionComplete = 1.0 // 100% complete

	// Do a region
	region := result.Regions["fips"]
	region.Candidates = append(region.Candidates, RegionCandidate{Candidate: "argusdusty", Votes: 100})
	region.Candidates = append(region.Candidates, RegionCandidate{Candidate: "badguy", Votes: 1})
	for _, c := range region.Candidates {
		region.TotalVotes += c.Votes
	}
	region.PortionComplete = 1.0 // 100% complete

	// Sort the data to get the order of candidates right
	result.Sort()
	return
}
