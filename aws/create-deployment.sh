#!/bin/bash

aws apigateway create-deployment --rest-api-id 7ld4iwn31l --profile pipeline-kafka-deploy --stage-name development | jq