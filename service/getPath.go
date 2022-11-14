package service

type MetadataRequestStruct struct {
	APIs ApiSlice
}

type ApiSlice []ApiStruct

type ApiStruct struct {
	Services   ServicesStruct
	Path       string
	Protocol   string
	Method     string
	Parameter  []ParametersStruct
	Headers    []HeaderStruct
	BatchIndex int
	ParentApi  []ParentApiStruct
	Name       string
}

type ParentApiStruct struct {
	ParentName   string
	Key          string
	DefaultValue interface{}
	ToKey        string
	ToType       string
}

type ServicesStruct struct {
	Host string
	Port string
}

type HeaderStruct struct {
	Key          string
	DefaultValue string
}

type ParametersStruct struct {
	Type         string
	Key          string
	DefaultValue string
	Require      bool
	Body         map[string]interface{}
}

func (api *ApiSlice) Len() int {
	return len(*api)
}

func (api *ApiSlice) Less(i, j int) bool {
	return (*api)[i].BatchIndex < (*api)[j].BatchIndex
}

func (api *ApiSlice) Swap(i, j int) {
	(*api)[i], (*api)[j] = (*api)[j], (*api)[i]
}

func GetPathFromMetadata(path string, method string) MetadataRequestStruct {
	api := ApiStruct{
		Services: ServicesStruct{
			Host: "127.0.0.1",
			Port: "8081",
		},
		Path:     "/ping",
		Protocol: "http",
		Method:   "GET",
		Parameter: []ParametersStruct{{
			Key:          "id",
			DefaultValue: "string",
			Require:      true,
			Type:         "query",
			Body:         nil,
		}, {
			Key:          "id",
			DefaultValue: "string",
			Require:      true,
			Type:         "body",
			Body: map[string]interface{}{
				"test": nil,
			},
		}},
		BatchIndex: 0,
		Name:       "ping",
	}
	api2 := ApiStruct{
		Services: ServicesStruct{
			Host: "127.0.0.1",
			Port: "8081",
		},
		Path:     "/ping1",
		Protocol: "http",
		Method:   "GET",
		Parameter: []ParametersStruct{{
			Key:          "id",
			Type:         "query",
			DefaultValue: "string",
			Require:      true,
			Body:         nil,
		}},
		Headers: []HeaderStruct{{
			Key:          "jwt",
			DefaultValue: "token",
		}, {
			Key:          "Content-Type",
			DefaultValue: "application/pdf",
		}},
		BatchIndex: 1,
		ParentApi: []ParentApiStruct{{
			ParentName:   "ping",
			Key:          "rewq",
			DefaultValue: 0,
			ToType:       "query",
			ToKey:        "ewqe",
		}},
		Name: "ping1",
	}

	api3 := ApiStruct{
		Services: ServicesStruct{
			Host: "127.0.0.1",
			Port: "8081",
		},
		Path:     "/postPing",
		Protocol: "http",
		Method:   "POST",
		Parameter: []ParametersStruct{{
			Key:          "postIds",
			DefaultValue: "12324 ",
			Require:      true,
			Type:         "query",
			Body:         nil,
		}, {
			Key:          "",
			DefaultValue: "",
			Require:      true,
			Type:         "json",
			Body: map[string]interface{}{
				"rerqwe": nil,
			},
		}},
		BatchIndex: 0,
		Name:       "postPing",
	}
	return MetadataRequestStruct{APIs: []ApiStruct{api, api2, api3}}
}
