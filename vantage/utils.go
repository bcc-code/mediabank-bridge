package vantage

func mapToList(m map[workflowParamName]*vantageWorkflowVariable) []vantageWorkflowVariable {
	l := make([]vantageWorkflowVariable, len(m))
	i := 0
	for _, v := range m {
		l[i] = *v
		i++
	}
	return l
}
