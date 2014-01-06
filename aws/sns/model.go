package sns

type CreateResponse struct {
	CreatePlatformEndpointResult Endpoint
}

type Endpoint struct {
	EndpointArn string
}
