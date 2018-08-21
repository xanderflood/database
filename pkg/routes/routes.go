package routes

//Info struct type for the resolver
type Info struct {
	CreateTable Route
	Index       Route
	Insert      Route
}

//Resolver provides path builders for all routes
var Resolver = Info{
	CreateTable: PlainRoute{
		route: "/v1/createTable",
	},
	Index: StringRoute{
		route: "/v1/index/{table}",
		keys:  []string{"{table}"},
	},
	Insert: StringRoute{
		route: "/v1/insert/{table}",
		keys:  []string{"{table}"},
	},
}
