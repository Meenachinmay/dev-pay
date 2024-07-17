### Blazing fast gRPC server setup with Tigerbeetle
* This repo has the code for a gRPC server which can handle 5 million client gRPC calls efficiently.
* Motive of this repo was to try TB txns database to check it's capabilities.
* Tigerbeetle is superfast transactional database which can handle 5,000,000 accounts creation in just 5 seconds and 1,000,000 in just 0.19s.
* It took me time to setup the concurrency with golang and it was not a easy project.