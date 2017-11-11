package main

type UseridIndex map[string][]ApacheLogEntry

func NewuseridIndex(parseIndex []LogIndex) *UseridIndex {
	index := make(UseridIndex)
	for _, logIndex := range parseIndex {
		for userid, log := range logIndex {
			index[userid] = append(index[userid], log...)
		}
	}
	return &index
}
