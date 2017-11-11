package main

//LogOption is the CLI file targets
type LogOption []string

func (i *LogOption) String() string {
	return "Representation"
}

func (i *LogOption) Set(value string) error {
	*i = append(*i, value)
	return nil
}
