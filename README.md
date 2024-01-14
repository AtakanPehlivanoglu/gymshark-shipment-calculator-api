# gymshark-shipment-calculator-api
Coding challenge for solving items packing problem.

## Deployments
Pushing a new commit will trigger [Push Docker Image to ECR Github Action Workflow](https://github.com/AtakanPehlivanoglu/gymshark-shipment-calculator-api/actions/workflows/main.yml) which will update private ECR repository with the latest Docker image. 

Automatic Deployments are enabled on AWS AppRunner whenever new Docker image is being pushed to ECR repositry.  
 
## Swagger UI
[Swagger](https://exz4e5um5a.eu-central-1.awsapprunner.com/swagger/index.html#/default/get_calculate__itemCount_) 

~~Currently in disabled status, it could be activated on AWS App Runner anytime.~~

## Run Locally
### Option 1
`SPEC_FILE_PATH=config go run ./shipment-calculator-api`

### Option 2
Add `SPEC_FILE_PATH=config` ENV variable in anyway and run `main.go` file.



