package app

import (
	"context"
	"log"
	"revel-dynamodb-v2/app/repositories"
	"revel-dynamodb-v2/app/services"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	_ "github.com/revel/modules"
	"github.com/revel/revel"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string

	Service *services.MovieService
)

func InitDynamoDB() *dynamodb.Client {
	fakeAccessKeyId, _ := revel.Config.String("aws.accesKeyId")
	fakeSecretAccessKey, _ := revel.Config.String("aws.secretAccessKey")
	region, _ := revel.Config.String("aws.region")
	port, _ := revel.Config.String("db.port")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			fakeAccessKeyId,      // Dummy AWS Access Key ID
			fakeSecretAccessKey, // Dummy AWS Secret Access Key
			"",                 // Session token, leave empty for local DynamoDB
		)),
	)
	if err != nil {
		log.Fatalf("failed to load AWS configuration: %v", err)
	}

    client := dynamodb.NewFromConfig(cfg, func (o *dynamodb.Options) {
        o.BaseEndpoint = aws.String(port)
    })
	log.Println("DynamoDB Client initialized successfully")
	return client
}

func InitRepositories() repositories.Repository {
	tableName, _ := revel.Config.String("db.name")
	client := InitDynamoDB()
	repository := repositories.NewDynamoDBRepository(client, tableName)
	log.Println("Repositories initialized successfully")
	return repository
}

func InitServices() {
	repo := InitRepositories()
	Service = services.NewMovieService(repo)
	log.Println("Services initialized successfully")
}

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	// Register startup functions with OnAppStart
	// revel.DevMode and revel.RunMode only work inside of OnAppStart. See Example Startup Script
	// ( order dependent )
	// revel.OnAppStart(ExampleStartupScript)
	// revel.OnAppStart(InitDB)
	// revel.OnAppStart(FillCache)
	revel.OnAppStart(InitServices)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

//func ExampleStartupScript() {
//	// revel.DevMod and revel.RunMode work here
//	// Use this script to check for dev mode and set dev/prod startup scripts here!
//	if revel.DevMode == true {
//		// Dev mode
//	}
//}