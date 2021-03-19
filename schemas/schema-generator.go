package types

//go:generate mkdir -p ../models/generated/submission
//go:generate schema-generate -i ./submissionmodel-schema-v1.0.json -o ../models/generated/submission/SubmissionModel.go -p submission

//go:generate mkdir -p ../models/generated/initiation
//go:generate schema-generate -i ./initiationmodel-schema-v1.0.json -o ../models/generated/initiation/InitiationModel.go -p initiation
