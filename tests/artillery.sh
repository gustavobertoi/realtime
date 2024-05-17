docker-compose up -d

artillery run "artillery.yaml" --output tests/artifacts/report.json

docker-compose down