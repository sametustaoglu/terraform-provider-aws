package appsync_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/appsync"
	"github.com/hashicorp/aws-sdk-go-base/tfawserr"
	sdkacctest "github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	tfappsync "github.com/hashicorp/terraform-provider-aws/internal/service/appsync"
)

func testAccAppSyncResolver_basic(t *testing.T) {
	var resolver1 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver1),
					acctest.MatchResourceAttrRegionalARN(resourceName, "arn", "appsync", regexp.MustCompile("apis/.+/types/.+/resolvers/.+")),
					resource.TestCheckResourceAttr(resourceName, "data_source", rName),
					resource.TestCheckResourceAttrSet(resourceName, "request_template"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppSyncResolver_disappears(t *testing.T) {
	var api1 appsync.GraphqlApi
	var resolver1 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	appsyncGraphqlApiResourceName := "aws_appsync_graphql_api.test"
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_basic(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGraphQLAPIExists(appsyncGraphqlApiResourceName, &api1),
					testAccCheckResolverExists(resourceName, &resolver1),
					acctest.CheckResourceDisappears(acctest.Provider, tfappsync.ResourceResolver(), resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccAppSyncResolver_dataSource(t *testing.T) {
	var resolver1, resolver2 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_DataSource(rName, "test_ds_1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver1),
					resource.TestCheckResourceAttr(resourceName, "data_source", "test_ds_1"),
				),
			},
			{
				Config: testAccAppsyncResolver_DataSource(rName, "test_ds_2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver2),
					resource.TestCheckResourceAttr(resourceName, "data_source", "test_ds_2"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppSyncResolver_DataSource_lambda(t *testing.T) {
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_DataSource_lambda(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver),
					resource.TestCheckResourceAttr(resourceName, "data_source", rName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppSyncResolver_requestTemplate(t *testing.T) {
	var resolver1, resolver2 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_RequestTemplate(rName, "/"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver1),
					resource.TestMatchResourceAttr(resourceName, "request_template", regexp.MustCompile("resourcePath\": \"/\"")),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppsyncResolver_RequestTemplate(rName, "/test"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver2),
					resource.TestMatchResourceAttr(resourceName, "request_template", regexp.MustCompile("resourcePath\": \"/test\"")),
				),
			},
		},
	})
}

func testAccAppSyncResolver_responseTemplate(t *testing.T) {
	var resolver1, resolver2 appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_ResponseTemplate(rName, 200),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver1),
					resource.TestMatchResourceAttr(resourceName, "response_template", regexp.MustCompile(`ctx\.result\.statusCode == 200`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAppsyncResolver_ResponseTemplate(rName, 201),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver2),
					resource.TestMatchResourceAttr(resourceName, "response_template", regexp.MustCompile(`ctx\.result\.statusCode == 201`)),
				),
			},
		},
	})
}

func testAccAppSyncResolver_multipleResolvers(t *testing.T) {
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_multipleResolvers(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName+"1", &resolver),
					testAccCheckResolverExists(resourceName+"2", &resolver),
					testAccCheckResolverExists(resourceName+"3", &resolver),
					testAccCheckResolverExists(resourceName+"4", &resolver),
					testAccCheckResolverExists(resourceName+"5", &resolver),
					testAccCheckResolverExists(resourceName+"6", &resolver),
					testAccCheckResolverExists(resourceName+"7", &resolver),
					testAccCheckResolverExists(resourceName+"8", &resolver),
					testAccCheckResolverExists(resourceName+"9", &resolver),
					testAccCheckResolverExists(resourceName+"10", &resolver),
				),
			},
		},
	})
}

func testAccAppSyncResolver_pipeline(t *testing.T) {
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_pipelineConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver),
					resource.TestCheckResourceAttr(resourceName, "pipeline_config.0.functions.#", "1"),
					resource.TestCheckResourceAttrPair(resourceName, "pipeline_config.0.functions.0", "aws_appsync_function.test", "function_id"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppSyncResolver_caching(t *testing.T) {
	var resolver appsync.Resolver
	rName := fmt.Sprintf("tfacctest%d", sdkacctest.RandInt())
	resourceName := "aws_appsync_resolver.test"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { acctest.PreCheck(t); acctest.PreCheckPartitionHasService(appsync.EndpointsID, t) },
		ErrorCheck:   acctest.ErrorCheck(t, appsync.EndpointsID),
		Providers:    acctest.Providers,
		CheckDestroy: testAccCheckResolverDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAppsyncResolver_cachingConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResolverExists(resourceName, &resolver),
					resource.TestCheckResourceAttr(resourceName, "caching_config.0.caching_keys.#", "2"),
					resource.TestCheckResourceAttr(resourceName, "caching_config.0.ttl", "60"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckResolverDestroy(s *terraform.State) error {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncConn
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_appsync_resolver" {
			continue
		}

		apiID, typeName, fieldName, err := tfappsync.DecodeResolverID(rs.Primary.ID)

		if err != nil {
			return err
		}

		input := &appsync.GetResolverInput{
			ApiId:     aws.String(apiID),
			TypeName:  aws.String(typeName),
			FieldName: aws.String(fieldName),
		}

		_, err = conn.GetResolver(input)

		if tfawserr.ErrCodeEquals(err, appsync.ErrCodeNotFoundException) {
			continue
		}

		if err != nil {
			return err
		}
	}
	return nil
}

func testAccCheckResolverExists(name string, resolver *appsync.Resolver) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource has no ID: %s", name)
		}

		apiID, typeName, fieldName, err := tfappsync.DecodeResolverID(rs.Primary.ID)

		if err != nil {
			return err
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncConn

		input := &appsync.GetResolverInput{
			ApiId:     aws.String(apiID),
			TypeName:  aws.String(typeName),
			FieldName: aws.String(fieldName),
		}

		output, err := conn.GetResolver(input)

		if err != nil {
			return err
		}

		*resolver = *output.Resolver

		return nil
	}
}

func testAccAppsyncResolver_basic(rName string) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = %q
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, rName, rName)
}

func testAccAppsyncResolver_DataSource(rName, dataSource string) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test_ds_1" {
  api_id = aws_appsync_graphql_api.test.id
  name   = "test_ds_1"
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_datasource" "test_ds_2" {
  api_id = aws_appsync_graphql_api.test.id
  name   = "test_ds_2"
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.%s.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, rName, dataSource)
}

func testAccAppsyncResolver_DataSource_lambda(rName string) string {
	return testAccAppsyncDatasourceConfig_base_Lambda(rName) + fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  name             = %q
  service_role_arn = aws_iam_role.test.arn
  type             = "AWS_LAMBDA"

  lambda_config {
    function_arn = aws_lambda_function.test.arn
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name
}
`, rName, rName)
}

func testAccAppsyncResolver_RequestTemplate(rName, resourcePath string) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = %q
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": %q,
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, rName, rName, resourcePath)
}

func testAccAppsyncResolver_ResponseTemplate(rName string, statusCode int) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = %q
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id      = aws_appsync_graphql_api.test.id
  field       = "singlePost"
  type        = "Query"
  data_source = aws_appsync_datasource.test.name

  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        ## you can forward the headers using the below utility
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == %d)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, rName, rName, statusCode)
}

func testAccAppsyncResolver_multipleResolvers(rName string) string {
	var queryFields string
	var resolverResources string
	for i := 1; i <= 10; i++ {
		queryFields = queryFields + fmt.Sprintf(`
	singlePost%d(id: ID!): Post
`, i)
		resolverResources = resolverResources + fmt.Sprintf(`
resource "aws_appsync_resolver" "test%d" {
  api_id           = "${aws_appsync_graphql_api.test.id}"
  field            = "singlePost%d"
  type             = "Query"
  data_source      = "${aws_appsync_datasource.test.name}"
  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF
  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}
`, i, i)
	}

	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = %q

  schema = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
%s
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = %q
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

%s

`, rName, queryFields, rName, resolverResources)
}

func testAccAppsyncResolver_pipelineConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = "%[1]s"
  schema              = <<EOF
type Mutation {
		putPost(id: ID!, title: String!): Post
}

type Post {
		id: ID!
		title: String!
}

type Query {
		singlePost(id: ID!): Post
}

schema {
		query: Query
		mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = "%[1]s"
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_function" "test" {
  api_id                   = aws_appsync_graphql_api.test.id
  data_source              = aws_appsync_datasource.test.name
  name                     = "%[1]s"
  request_mapping_template = <<EOF
{
		"version": "2018-05-29",
		"method": "GET",
		"resourcePath": "/",
		"params":{
				"headers": $utils.http.copyheaders($ctx.request.headers)
		}
}
EOF

  response_mapping_template = <<EOF
#if($ctx.result.statusCode == 200)
		$ctx.result.body
#else
		$utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF
}

resource "aws_appsync_resolver" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  field            = "singlePost"
  type             = "Query"
  kind             = "PIPELINE"
  request_template = <<EOF
{
		"version": "2018-05-29",
		"method": "GET",
		"resourcePath": "/",
		"params":{
				"headers": $utils.http.copyheaders($ctx.request.headers)
		}
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
		$ctx.result.body
#else
		$utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF

  pipeline_config {
    functions = [aws_appsync_function.test.function_id]
  }
}

`, rName)
}

func testAccAppsyncResolver_cachingConfig(rName string) string {
	return fmt.Sprintf(`
resource "aws_appsync_graphql_api" "test" {
  authentication_type = "API_KEY"
  name                = "%[1]s"
  schema              = <<EOF
type Mutation {
	putPost(id: ID!, title: String!): Post
}

type Post {
	id: ID!
	title: String!
}

type Query {
	singlePost(id: ID!): Post
}

schema {
	query: Query
	mutation: Mutation
}
EOF
}

resource "aws_appsync_datasource" "test" {
  api_id = aws_appsync_graphql_api.test.id
  name   = "%[1]s"
  type   = "HTTP"

  http_config {
    endpoint = "http://example.com"
  }
}

resource "aws_appsync_resolver" "test" {
  api_id           = aws_appsync_graphql_api.test.id
  field            = "singlePost"
  type             = "Query"
  kind             = "UNIT"
  data_source      = aws_appsync_datasource.test.name
  request_template = <<EOF
{
    "version": "2018-05-29",
    "method": "GET",
    "resourcePath": "/",
    "params":{
        "headers": $utils.http.copyheaders($ctx.request.headers)
    }
}
EOF

  response_template = <<EOF
#if($ctx.result.statusCode == 200)
    $ctx.result.body
#else
    $utils.appendError($ctx.result.body, $ctx.result.statusCode)
#end
EOF

  caching_config {
    caching_keys = [
      "$context.identity.sub",
      "$context.arguments.id",
    ]
    ttl = 60
  }
}
`, rName)
}
