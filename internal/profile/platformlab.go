package profile

type Dependency struct {
	Name string
	Host string
	Port int
}

type Service struct {
	Name         string
	Container    string
	Port         int
	HealthPath   string
	ReadyPath    string
	Dependencies []Dependency
}

var PlatformLabServices = map[string]Service{
	"web-bff": {
		Name:       "web-bff",
		Container:  "platform-lab-web-bff",
		Port:       8080,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",
	},

	"catalog-service": {
		Name:       "catalog-service",
		Container:  "platform-lab-catalog-service",
		Port:       8081,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",

		Dependencies: []Dependency{
			{
				Name: "PostgreSQL",
				Host: "postgres",
				Port: 5432,
			},
		},
	},

	"cart-service": {
		Name:       "cart-service",
		Container:  "platform-lab-cart-service",
		Port:       8082,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",

		Dependencies: []Dependency{
			{
				Name: "Redis",
				Host: "redis",
				Port: 6379,
			},
		},
	},

	"order-service": {
		Name:       "order-service",
		Container:  "platform-lab-order-service",
		Port:       8083,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",

		Dependencies: []Dependency{
			{
				Name: "PostgreSQL",
				Host: "postgres",
				Port: 5432,
			},
			{
				Name: "Kafka",
				Host: "kafka",
				Port: 9092,
			},
		},
	},

	"inventory-service": {
		Name:       "inventory-service",
		Container:  "platform-lab-inventory-service",
		Port:       8084,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",

		Dependencies: []Dependency{
			{
				Name: "PostgreSQL",
				Host: "postgres",
				Port: 5432,
			},
			{
				Name: "Kafka",
				Host: "kafka",
				Port: 9092,
			},
		},
	},

	"payment-service": {
		Name:       "payment-service",
		Container:  "platform-lab-payment-service",
		Port:       8085,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",

		Dependencies: []Dependency{
			{
				Name: "PostgreSQL",
				Host: "postgres",
				Port: 5432,
			},
			{
				Name: "Kafka",
				Host: "kafka",
				Port: 9092,
			},
		},
	},

	"notification-service": {
		Name:       "notification-service",
		Container:  "platform-lab-notification-service",
		Port:       8086,
		HealthPath: "/healthz",
		ReadyPath:  "/readyz",

		Dependencies: []Dependency{
			{
				Name: "PostgreSQL",
				Host: "postgres",
				Port: 5432,
			},
			{
				Name: "Kafka",
				Host: "kafka",
				Port: 9092,
			},
		},
	},
}

func FindService(
	name string,
) (Service, bool) {

	if service, ok :=
		PlatformLabServices[name]; ok {

		return service, true
	}

	for _, service := range PlatformLabServices {

		if service.Container == name {
			return service, true
		}
	}

	return Service{}, false
}
