package self

type AgentName string

const (
	AgentNameDefault             AgentName = "default"
	AgentNameAgentWriteParseData AgentName = "agent_write_parse_data"
	AgentNameAgentWriteStoreData AgentName = "agent_write_store_data"
	AgentNameAgentReadParseData  AgentName = "agent_read_parse_data"
	AgentNameAgentReadResponse   AgentName = "agent_read_response"
)
