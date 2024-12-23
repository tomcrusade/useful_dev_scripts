package entity

type CloudServiceTechStackName string

var (
	CloudServiceTechStackNameMySQL      = CloudServiceTechStackName("mysql")
	CloudServiceTechStackNamePostgreSQL = CloudServiceTechStackName("postgres")
	CloudServiceTechStackNameVault      = CloudServiceTechStackName("vault")
)

var CloudServiceTechStackMap = map[CloudServiceTechStackName]CloudServiceTechStackName{
	CloudServiceTechStackNameMySQL:      CloudServiceTechStackNameMySQL,
	CloudServiceTechStackNamePostgreSQL: CloudServiceTechStackNamePostgreSQL,
	CloudServiceTechStackNameVault:      CloudServiceTechStackNameVault,
}

type CloudServiceEnvName string

var (
	CloudServiceEnvNameDevelopment = CloudServiceEnvName("dev")
	CloudServiceEnvNameStaging     = CloudServiceEnvName("stg")
	CloudServiceEnvNameProduction  = CloudServiceEnvName("prod")
)
